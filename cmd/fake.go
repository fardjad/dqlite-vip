package cmd

import (
	"log"

	"fardjad.com/dqlite-vip/cluster"
)

type fakeClusterNode struct{}

func (n *fakeClusterNode) Start() error {
	log.Println("Starting fake cluster node")
	return nil
}

func (n *fakeClusterNode) IsLeader() bool {
	return false
}

func (n *fakeClusterNode) SetString(key string, value string) error {
	return nil
}

func (n *fakeClusterNode) GetString(key string) (string, error) {
	return "", nil
}

func (n *fakeClusterNode) Close() error {
	log.Println("Closing fake cluster node")
	return nil
}

type fakeClusterNodeFactory struct{}

func (f *fakeClusterNodeFactory) NewClusterNode(dataDir, bindCluster, bindHttp string, join []string) (cluster.ClusterNode, error) {
	return &fakeClusterNode{}, nil
}
