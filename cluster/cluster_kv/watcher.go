package cluster_kv

import (
	"log"

	"fardjad.com/dqlite-vip/cluster"
	"fardjad.com/dqlite-vip/cluster/cluster_events"
	"fardjad.com/dqlite-vip/time"
)

// Implements [Watcher]
type ClusterNodeWatcher struct {
	clusterNode   cluster.ClusterNode
	kv            map[string]string
	changeEmitter *cluster_events.ChangeEmitter
	ticker        time.Ticker
}

func (w *ClusterNodeWatcher) Watch(key string) (chan cluster_events.Change, cluster_events.CancelFunc, error) {
	subcription := w.changeEmitter.Subscribe(key)

	go func() {
		for range w.ticker.C() {
			currentValue, err := w.clusterNode.GetString(key)
			if err != nil {
				log.Printf("Failed to get value for key %s: %v", key, err)
			}
			previousValue := w.kv[key]
			if currentValue != previousValue {
				w.kv[key] = currentValue
				change := cluster_events.Change{Previous: previousValue, Current: currentValue}
				w.changeEmitter.Publish(key, change)
			}
		}
	}()

	return subcription.Ch, func() {
		w.ticker.Stop()
		subcription.Cancel()
	}, nil
}

// Implements [WatcherFactoryFunc]
func NewClusterNodeWatcher(clusterNode cluster.ClusterNode, ticker time.Ticker) Watcher {
	return &ClusterNodeWatcher{
		clusterNode: clusterNode,
		ticker:      ticker,

		kv:            make(map[string]string),
		changeEmitter: cluster_events.NewChangeEmitter(),
	}
}
