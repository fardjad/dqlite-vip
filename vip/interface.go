package vip

import (
	"context"

	"fardjad.com/dqlite-vip/cluster"
	"fardjad.com/dqlite-vip/utils"
)

type Manager interface {
	Start(ctx context.Context)
	Stop()
}

type ManagerFactoryFunc func(cluster.ClusterNode, utils.BetterTicker, Configurator, string) Manager

type Configurator interface {
	EnsureVIP(iface string, address string) error
	RemoveVIP(iface string, address string) error
}

type ConfiguratorFactoryFunc func(gratuitousArpCount int) Configurator
