package cluster_kv

import (
	"fardjad.com/dqlite-vip/cluster"
	"fardjad.com/dqlite-vip/cluster/cluster_events"
	"fardjad.com/dqlite-vip/time"
)

type Watcher interface {
	Watch(key string) (chan cluster_events.Change, cluster_events.CancelFunc, error)
}

type WatcherFactoryFunc func(clusterNode cluster.ClusterNode, ticker time.Ticker) Watcher
