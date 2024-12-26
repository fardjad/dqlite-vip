package cluster

import "context"

type ClusterMemberInfo struct {
	ID      uint64 `json:"id"`
	Address string `json:"address"`
	Role    string `json:"role"`
}

type ClusterNode interface {
	Start() error
	Ready(ctx context.Context) error
	ID() uint64
	LeaderID(ctx context.Context) (uint64, error)
	ClusterMembers(ctx context.Context) ([]*ClusterMemberInfo, error)
	Close() error

	SetString(key string, value string) error
	GetString(key string) (string, error)
}

type ClusterNodeFactory interface {
	NewClusterNode(dataDir string, bindCluster string, join []string) (ClusterNode, error)
}
