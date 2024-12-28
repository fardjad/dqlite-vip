package vip

import (
	"context"

	"fardjad.com/dqlite-vip/cluster"
	"fardjad.com/dqlite-vip/cluster/cluster_kv"
)

type Manager interface {
	Start(ctx context.Context) error
	Stop()
}

type ManagerFactoryFunc func(cluster.ClusterNode, cluster_kv.Watcher, Configurator, string) Manager

type Configurator interface {
	AddVIP(iface string, address string) error
	RemoveVIP(iface string, address string) error
}

type ConfiguratorFactoryFunc func(gratuitousArpCount int) Configurator
