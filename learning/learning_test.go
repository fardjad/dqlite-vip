package learning

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/canonical/go-dqlite/v2/app"
)

var averageNetworkLatency = 5 * time.Millisecond

func startNode(
	dataPath string,
	nodeAddress string,
	otherNodesAddresses []string,
) (*app.App, error) {
	nodeAddressHash := md5.Sum([]byte(nodeAddress))
	safeFileName := hex.EncodeToString(nodeAddressHash[:])
	dbDir := filepath.Join(dataPath, safeFileName)

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("can't create %s: %v", dbDir, err)
	}

	options := []app.Option{
		app.WithNetworkLatency(averageNetworkLatency),
		app.WithAddress(nodeAddress),
		app.WithCluster(otherNodesAddresses),
	}

	app, err := app.New(dbDir, options...)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func makeAddresses(n int) []string {
	addresses := make([]string, 0)
	for i := 0; i < 3; i++ {
		addresses = append(addresses, "127.0.0.1"+":"+strconv.Itoa(FindFreePort()))
	}

	return addresses
}

func startThreeNodesCluster(dataPath string, addresses []string) []*app.App {
	log.Println("Starting a 3-node cluster")

	nodes := make([]*app.App, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	nodes = append(nodes, Must(startNode(dataPath, addresses[0], []string{})))
	nodes = append(nodes, Must(startNode(dataPath, addresses[1], addresses[0:1])))
	nodes = append(nodes, Must(startNode(dataPath, addresses[2], addresses[0:2])))

	for i, node := range nodes {
		node.Ready(ctx)
		log.Printf("Node %d is ready at %s", i, node.Address())
	}

	cancel()

	return nodes
}

func getLeaderAddress(ctx context.Context, app *app.App) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	leader := Must(Must(app.FindLeader(ctx)).Leader(ctx))
	return leader.Address
}

func TestHealthyCluster(t *testing.T) {
	nodes := startThreeNodesCluster(t.TempDir(), makeAddresses(3))

	for _, node := range nodes {
		node.Close()
	}
}

func TestReElection(t *testing.T) {
	nodes := startThreeNodesCluster(t.TempDir(), makeAddresses(3))

	log.Println("Stopping the first node")
	nodes[0].Close()
	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	for _, node := range nodes[1:] {
		log.Printf("%s: The leader is %s", node.Address(), getLeaderAddress(ctx, node))
	}
	cancel()

	for _, node := range nodes[1:] {
		node.Close()
	}
}

func TestPartialRecovery(t *testing.T) {
	addresses := makeAddresses(3)
	tempDir := t.TempDir()
	nodes := startThreeNodesCluster(tempDir, addresses)

	log.Println("Stopping the first node")
	nodes[0].Close()
	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	for _, node := range nodes[1:] {
		log.Printf("%s: The leader is %s", node.Address(), getLeaderAddress(ctx, node))
	}
	cancel()

	log.Println("Stopping the second node")
	nodes[1].Close()
	time.Sleep(1 * time.Second)

	log.Println("Starting the second node again")
	// After the cluster is formed, it's no longer necessary to pass all the addresses
	nodes[1] = Must(startNode(tempDir, addresses[1], []string{}))

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	for _, node := range nodes[1:] {
		log.Printf("%s: The leader is %s", node.Address(), getLeaderAddress(ctx, node))
	}
	cancel()

	for _, node := range nodes[1:] {
		node.Close()
	}
}

func TestFullRecovery(t *testing.T) {
	addresses := makeAddresses(3)
	tempDir := t.TempDir()
	nodes := startThreeNodesCluster(tempDir, addresses)

	log.Println("Stopping all nodes")
	for _, node := range nodes {
		node.Close()
	}
	time.Sleep(1 * time.Second)

	log.Println("Starting all nodes again")
	for i := 0; i < 3; i++ {
		// After the cluster is formed, it's no longer necessary to pass all the addresses
		nodes[i] = Must(startNode(tempDir, addresses[i], []string{}))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	for _, node := range nodes {
		log.Printf("%s: The leader is %s", node.Address(), getLeaderAddress(ctx, node))
	}
	cancel()

	for _, node := range nodes[1:] {
		node.Close()
	}
}
