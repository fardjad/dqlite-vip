package cluster

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/canonical/go-dqlite/v2/app"
)

type DqliteClusterNodeFactory struct{}

func (f *DqliteClusterNodeFactory) NewClusterNode(dataDir string, bindCluster string, join []string) (ClusterNode, error) {
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

func (n *DqliteClusterNode) Start() error {
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

	n.app = app
	return nil
}

func (n *DqliteClusterNode) Ready(ctx context.Context) error {
	return n.app.Ready(ctx)
}

func (n *DqliteClusterNode) ID() uint64 {
	return n.app.ID()
}

func (n *DqliteClusterNode) LeaderID(ctx context.Context) (uint64, error) {
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

func (n *DqliteClusterNode) ClusterMembers(ctx context.Context) ([]*ClusterMemberInfo, error) {
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

func (n *DqliteClusterNode) Close() error {
	return n.app.Close()
}

func (n *DqliteClusterNode) SetString(key string, value string) error {
	return errors.ErrUnsupported
}

func (n *DqliteClusterNode) GetString(key string) (string, error) {
	return "", errors.ErrUnsupported
}
