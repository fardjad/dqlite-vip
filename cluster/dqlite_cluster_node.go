package cluster

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/canonical/go-dqlite/v2/app"
)

func NewDqliteClusterNode(dataDir string, bindCluster string, join []string) (ClusterNode, error) {
	return &DqliteClusterNode{
		dataDir:     dataDir,
		bindCluster: bindCluster,
		join:        join,
	}, nil
}

type DqliteClusterNode struct {
	app *app.App

	dataDir     string
	bindCluster string
	join        []string
}

func (n *DqliteClusterNode) Start(ctx context.Context) error {
	if err := os.MkdirAll(n.dataDir, 0755); err != nil {
		return fmt.Errorf("can't create %s: %v", n.dataDir, err)
	}

	options := []app.Option{
		app.WithAddress(n.bindCluster),
		app.WithCluster(n.join),
	}
	app, err := app.New(n.dataDir, options...)
	if err != nil {
		return err
	}

	db, err := app.Open(ctx, "db")
	defer db.Close()
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS kv (key TEXT PRIMARY KEY, value TEXT)")
	if err != nil {
		return err
	}

	log.Printf("Dqlite node is listening on %s", n.bindCluster)

	n.app = app
	return nil
}

func (n *DqliteClusterNode) Ready(ctx context.Context) error {
	if n.app == nil {
		return errors.New("node not started")
	}

	return n.app.Ready(ctx)
}

func (n *DqliteClusterNode) ID() uint64 {
	if n.app == nil {
		return 0
	}

	return n.app.ID()
}

func (n *DqliteClusterNode) LeaderID(ctx context.Context) (uint64, error) {
	if n.app == nil {
		return 0, errors.New("node not started")
	}

	client, err := n.app.FindLeader(ctx)
	if err != nil {
		return 0, err
	}

	leader, err := client.Leader(ctx)
	if err != nil {
		return 0, err
	}

	return leader.ID, nil
}

func (n *DqliteClusterNode) IsLeader(ctx context.Context) bool {
	if n.app == nil {
		return false
	}

	leaderID, err := n.LeaderID(ctx)
	if err != nil {
		return false
	}

	return n.app.ID() == leaderID
}

func (n *DqliteClusterNode) ClusterMembers(ctx context.Context) ([]*ClusterMemberInfo, error) {
	if n.app == nil {
		return nil, errors.New("node not started")
	}

	client, err := n.app.Leader(ctx)
	if err != nil {
		return nil, err
	}
	nodeInfos, err := client.Cluster(ctx)
	if err != nil {
		return nil, err
	}

	clusterMembers := make([]*ClusterMemberInfo, 0)
	for _, nodeInfo := range nodeInfos {
		clusterMembers = append(clusterMembers, &ClusterMemberInfo{
			ID:      nodeInfo.ID,
			Address: nodeInfo.Address,
			Role:    nodeInfo.Role.String(),
		})
	}

	return clusterMembers, nil
}

func (n *DqliteClusterNode) Close(ctx context.Context) error {
	if n.app == nil {
		return errors.New("node not started")
	}

	err := n.app.Handover(ctx)
	if err != nil {
		log.Println("Failed to handover leadership:", err)
	}

	return n.app.Close()
}

func (n *DqliteClusterNode) SetString(ctx context.Context, key string, value string) error {
	if n.app == nil {
		return errors.New("node not started")
	}

	db, err := n.app.Open(ctx, "db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("INSERT OR REPLACE INTO kv (key, value) VALUES (?, ?)", key, value)
	return err
}

func (n *DqliteClusterNode) GetString(ctx context.Context, key string) (string, error) {
	if n.app == nil {
		return "", errors.New("node not started")
	}

	db, err := n.app.Open(ctx, "db")
	if err != nil {
		return "", err
	}
	defer db.Close()
	row := db.QueryRow("SELECT value FROM kv WHERE key = ?", key)
	if row.Err() != nil {
		return "", row.Err()
	}

	var value string
	err = row.Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}
