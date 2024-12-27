package cluster_kv

import (
	"log"

	"fardjad.com/dqlite-vip/cluster"
	"fardjad.com/dqlite-vip/cluster/cluster_events"
	"fardjad.com/dqlite-vip/time"
)

type Watcher struct {
	clusterNode   cluster.ClusterNode
	kv            map[string]string
	changeEmitter *cluster_events.ChangeEmitter
	ticker        time.Ticker
}

func NewWatcher(clusterNode cluster.ClusterNode, ticker time.Ticker) *Watcher {
	return &Watcher{
		clusterNode: clusterNode,
		ticker:      ticker,

		kv:            make(map[string]string),
		changeEmitter: cluster_events.NewChangeEmitter(),
	}
}

func (w *Watcher) Watch(key string) (chan string, cluster_events.CancelFunc, error) {
	subcription := w.changeEmitter.Subscribe(key)

	go func() {
		for range w.ticker.C() {
			value, err := w.clusterNode.GetString(key)
			if err != nil {
				log.Printf("Failed to get value for key %s: %v", key, err)
			}
			if value != w.kv[key] {
				w.kv[key] = value
				w.changeEmitter.Publish(key, value)
			}
		}
	}()

	return subcription.Ch, func() {
		w.ticker.Stop()
		subcription.Cancel()
	}, nil
}
