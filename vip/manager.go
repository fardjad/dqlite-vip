package vip

import (
	"context"
	"log"

	"fardjad.com/dqlite-vip/cluster"
	"fardjad.com/dqlite-vip/cluster/cluster_events"
	"fardjad.com/dqlite-vip/cluster/cluster_kv"
)

// Implements [Manager]
type manager struct {
	clusterNode  cluster.ClusterNode
	watcher      cluster_kv.Watcher
	cancelFunc   cluster_events.CancelFunc
	configurator Configurator

	iface   string
	address string
}

func (m *manager) Start(ctx context.Context) error {
	ch, cancel, err := m.watcher.Watch("vip")
	if err != nil {
		return err
	}
	m.cancelFunc = cancel

	go func() {
		for change := range ch {
			if !m.clusterNode.IsLeader(ctx) {
				return
			}

			err := m.configurator.RemoveVIP(m.iface, change.Previous)
			if err != nil {
				log.Printf("failed to remove VIP: %v", err)
			}

			err = m.configurator.AddVIP(m.iface, change.Current)
			if err != nil {
				log.Printf("failed to add VIP: %v", err)
			}
		}
	}()

	return nil
}

func (m *manager) Stop() {
	if m.cancelFunc != nil {
		m.cancelFunc()
	}
}

// Implements [ManagerFactoryFunc]
func NewManager(clusterNode cluster.ClusterNode, watcher cluster_kv.Watcher, configurator Configurator, iface string) Manager {
	return &manager{
		clusterNode:  clusterNode,
		watcher:      watcher,
		configurator: configurator,

		iface: iface,
	}
}
