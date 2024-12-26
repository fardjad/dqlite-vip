package cluster

import "context"

type ClusterNode interface {
	Start() error
	Ready(ctx context.Context) error
	IsLeader() bool
	SetString(key string, value string) error
	GetString(key string) (string, error)
	Close() error
}

type ClusterNodeFactory interface {
	NewClusterNode(dataDir string, bindCluster string, join []string) (ClusterNode, error)
}
