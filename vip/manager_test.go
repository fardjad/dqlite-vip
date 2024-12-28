package vip_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"fardjad.com/dqlite-vip/cluster/cluster_events"
	"fardjad.com/dqlite-vip/vip"

	mockCluster "fardjad.com/dqlite-vip/mocks/cluster"
	mockClusterEvents "fardjad.com/dqlite-vip/mocks/cluster/cluster_events"
	mockClusterKV "fardjad.com/dqlite-vip/mocks/cluster/cluster_kv"
	mockVip "fardjad.com/dqlite-vip/mocks/vip"
)

type ManagerTestSuite struct {
	suite.Suite
	clusterNode  *mockCluster.ClusterNode
	watcher      *mockClusterKV.Watcher
	cancelFunc   *mockClusterEvents.CancelFunc
	configurator *mockVip.Configurator
	manager      vip.Manager
}

func (s *ManagerTestSuite) SetupTest() {
	s.clusterNode = mockCluster.NewClusterNode(s.T())
	s.watcher = mockClusterKV.NewWatcher(s.T())
	s.cancelFunc = mockClusterEvents.NewCancelFunc(s.T())
	s.configurator = mockVip.NewConfigurator(s.T())
	s.manager = vip.NewManager(s.clusterNode, s.watcher, s.configurator, "dummy")
}

func (s *ManagerTestSuite) TestManager_WhenClusterNodeIsLeader() {
	s.clusterNode.EXPECT().IsLeader(mock.Anything).Return(true)

	changeChannel := make(chan cluster_events.Change)
	s.cancelFunc.EXPECT().Execute().RunAndReturn(func() { close(changeChannel) })
	s.watcher.EXPECT().Watch("vip").Return(changeChannel, s.cancelFunc.Execute, nil)

	s.configurator.EXPECT().RemoveVIP("dummy", "192.168.1.100").Return(nil)
	s.configurator.EXPECT().AddVIP("dummy", "192.168.1.101").Return(nil)

	s.NoError(s.manager.Start(context.Background()))
	defer s.manager.Stop()

	select {
	case changeChannel <- cluster_events.Change{Previous: "192.168.1.100", Current: "192.168.1.101"}:
	case <-time.After(1 * time.Second):
		s.Fail("Timeout waiting for changeChannel")
	}
}

func (s *ManagerTestSuite) TestManager_WhenClusterNodeIsNotLeader() {
	s.clusterNode.EXPECT().IsLeader(mock.Anything).Return(false)

	changeChannel := make(chan cluster_events.Change)
	s.cancelFunc.EXPECT().Execute().RunAndReturn(func() { close(changeChannel) })
	s.watcher.EXPECT().Watch("vip").Return(changeChannel, s.cancelFunc.Execute, nil)

	s.NoError(s.manager.Start(context.Background()))
	defer s.manager.Stop()

	select {
	case changeChannel <- cluster_events.Change{Previous: "192.168.1.100", Current: "192.168.1.101"}:
	case <-time.After(1 * time.Second):
		s.Fail("Timeout waiting for changeChannel")
	}
}

func TestManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}
