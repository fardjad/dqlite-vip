package vip

import (
	"context"
	"database/sql"
	"log"
	"time"

	"fardjad.com/dqlite-vip/cluster"
	"fardjad.com/dqlite-vip/utils"
)

// Implements [Manager]
type manager struct {
	clusterNode       cluster.ClusterNode
	ticker            utils.BetterTicker
	configurator      Configurator
	iface             string
	dbAccessTimeout   time.Duration
	findLeaderTimeout time.Duration

	// after start
	vip        string
	isLeader   bool
	cancelFunc context.CancelFunc
}

func (m *manager) readVIP(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, m.dbAccessTimeout)
	defer cancel()

	val, err := m.clusterNode.GetString(ctx, "vip")

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
	ctx, cancel := context.WithCancel(ctx)
	m.cancelFunc = cancel

	go func() {
		for {
			select {
			case <-m.ticker.C():
				m.isLeader = func() bool {
					ctx, cancel := context.WithTimeout(ctx, m.findLeaderTimeout)
					defer cancel()
					return m.clusterNode.IsLeader(ctx)
				}()

				newVIP, err := func() (string, error) {
					ctx, cancel := context.WithTimeout(ctx, m.dbAccessTimeout)
					defer cancel()
					return m.readVIP(ctx)
				}()
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
		// TODO: make these configurable
		dbAccessTimeout:   10 * time.Second,
		findLeaderTimeout: 5 * time.Second,
	}
}
