package cluster

type ClusterNode interface {
	Start() error
	IsLeader() bool
	SetString(key string, value string) error
	GetString(key string) (string, error)
	Close() error
}

type ClusterNodeFactory interface {
	NewClusterNode(dataDir string, bindCluster string, bindHttp string, join []string) (ClusterNode, error)
}
