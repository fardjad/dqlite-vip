package vip

import (
	"context"
	"database/sql"
	"log"

	"fardjad.com/dqlite-vip/cluster"
	"fardjad.com/dqlite-vip/utils"
)

// Implements [Manager]
type manager struct {
	clusterNode  cluster.ClusterNode
	ticker       utils.BetterTicker
	configurator Configurator
	iface        string

	// after start
	vip        string
	isLeader   bool
	cancelFunc context.CancelFunc
}

func (m *manager) readVIP() (string, error) {
	val, err := m.clusterNode.GetString("vip")

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (m *manager) removeVIPs(vips []string) {
	seen := make(map[string]struct{})

	for _, vip := range vips {
		if vip == "" {
			continue
		}

		if _, exists := seen[vip]; !exists {
			seen[vip] = struct{}{}

			if err := m.configurator.RemoveVIP(m.iface, vip); err != nil {
				log.Printf("failed to remove VIP %s from interface %s: %v", vip, m.iface, err)
			}
		}
	}
}

func (m *manager) ensureVIP(vip string) {
	if vip == "" {
		return
	}

	if err := m.configurator.EnsureVIP(m.iface, vip); err != nil {
		log.Printf("failed to add VIP %s to interface %s: %v", vip, m.iface, err)
	}
}

func (m *manager) Start(ctx context.Context) {
	ctx, cancelFunc := context.WithCancel(ctx)
	m.cancelFunc = cancelFunc

	go func() {
		for {
			select {
			case <-m.ticker.C():
				m.isLeader = m.clusterNode.IsLeader(ctx)

				newVIP, err := m.readVIP()
				if err != nil {
					log.Println("failed to read VIP from the database:", err)
					continue
				}

				vipChanged := m.vip != newVIP
				switch {
				case m.isLeader && vipChanged:
					m.removeVIPs([]string{m.vip})
					m.ensureVIP(newVIP)
				case m.isLeader && !vipChanged:
					m.ensureVIP(newVIP)
				default:
					m.removeVIPs([]string{m.vip, newVIP})
				}

				m.vip = newVIP
			case <-ctx.Done():
				m.ticker.Stop()
				m.removeVIPs([]string{m.vip})
				return
			}
		}
	}()
}

func (m *manager) Stop() {
	m.cancelFunc()
}

// Implements [ManagerFactoryFunc]
func NewManager(clusterNode cluster.ClusterNode, ticker utils.BetterTicker, configurator Configurator, iface string) Manager {
	return &manager{
		clusterNode:  clusterNode,
		ticker:       ticker,
		configurator: configurator,
		iface:        iface,
	}
}
