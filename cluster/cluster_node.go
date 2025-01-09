package cluster

import "context"

type ClusterMemberInfo struct {
	ID      uint64 `json:"id"`
	Address string `json:"address"`
	Role    string `json:"role"`
}

type ClusterNode interface {
	Start(ctx context.Context) error
	Ready(ctx context.Context) error
	ID() uint64
	LeaderID(ctx context.Context) (uint64, error)
	IsLeader(ctx context.Context) bool
	ClusterMembers(ctx context.Context) ([]*ClusterMemberInfo, error)
	Close(ctx context.Context) error

	SetString(ctx context.Context, key string, value string) error
	GetString(ctx context.Context, key string) (string, error)
}

type ClusterNodeFactoryFunc func(dataDir string, bindCluster string, join []string) (ClusterNode, error)
