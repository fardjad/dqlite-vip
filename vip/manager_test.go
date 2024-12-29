package vip_test

import (
	"context"
	"testing"
	"time"

	clusterMocks "fardjad.com/dqlite-vip/mocks/cluster"
	vipMocks "fardjad.com/dqlite-vip/mocks/vip"
	"fardjad.com/dqlite-vip/utils"
	"fardjad.com/dqlite-vip/vip"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ManagerTestSuite struct {
	suite.Suite

	clusterNode  *clusterMocks.ClusterNode
	ticker       *utils.FakeTicker
	configurator *vipMocks.Configurator
	manager      *vipMocks.Manager
}

func (s *ManagerTestSuite) SetupTest() {
	s.clusterNode = clusterMocks.NewClusterNode(s.T())
	s.ticker = utils.NewFakeTicker()
	s.configurator = vipMocks.NewConfigurator(s.T())
	s.manager = vipMocks.NewManager(s.T())
}

func (s *ManagerTestSuite) TestManager() {
	m := vip.NewManager(s.clusterNode, s.ticker, s.configurator, "dummy")
	m.Start(context.Background())
	defer m.Stop()

	var wg utils.WaitGroupWithTimeout

	// Node is the leader and VIP is not set
	wg.Add(1)
	s.clusterNode.EXPECT().IsLeader(mock.Anything).RunAndReturn(func(_ context.Context) bool {
		wg.Done()
		return true
	}).Once()
	wg.Add(1)
	s.clusterNode.EXPECT().GetString("vip").RunAndReturn(func(_ string) (string, error) {
		wg.Done()
		return "", nil
	}).Once()
	s.ticker.Tick(time.Now())
	s.True(wg.WaitWithTimeout(1 * time.Second))

	// Set VIP for the first time
	wg.Add(1)
	s.clusterNode.EXPECT().IsLeader(mock.Anything).RunAndReturn(func(_ context.Context) bool {
		wg.Done()
		return true
	}).Once()
	wg.Add(1)
	s.clusterNode.EXPECT().GetString("vip").RunAndReturn(func(_ string) (string, error) {
		wg.Done()
		return "192.168.1.100/24", nil
	}).Once()
	wg.Add(1)
	s.configurator.EXPECT().EnsureVIP("dummy", "192.168.1.100/24").RunAndReturn(func(_, _ string) error {
		wg.Done()
		return nil
	}).Once()
	s.ticker.Tick(time.Now())
	s.True(wg.WaitWithTimeout(1 * time.Second))

	// Change VIP
	wg.Add(1)
	s.clusterNode.EXPECT().IsLeader(mock.Anything).RunAndReturn(func(_ context.Context) bool {
		wg.Done()
		return true
	}).Once()
	wg.Add(1)
	s.clusterNode.EXPECT().GetString("vip").RunAndReturn(func(_ string) (string, error) {
		wg.Done()
		return "192.168.1.101/24", nil
	}).Once()
	wg.Add(1)
	s.configurator.EXPECT().RemoveVIP("dummy", "192.168.1.100/24").RunAndReturn(func(_, _ string) error {
		wg.Done()
		return nil
	}).Once()
	wg.Add(1)
	s.configurator.EXPECT().EnsureVIP("dummy", "192.168.1.101/24").RunAndReturn(func(_, _ string) error {
		wg.Done()
		return nil
	}).Once()
	s.ticker.Tick(time.Now())
	s.True(wg.WaitWithTimeout(1 * time.Second))

	// Node is no longer the leader
	wg.Add(1)
	s.clusterNode.EXPECT().IsLeader(mock.Anything).RunAndReturn(func(_ context.Context) bool {
		wg.Done()
		return false
	}).Once()
	wg.Add(1)
	s.clusterNode.EXPECT().GetString("vip").RunAndReturn(func(_ string) (string, error) {
		wg.Done()
		return "192.168.1.101/24", nil
	}).Once()
	wg.Add(1)
	s.configurator.EXPECT().RemoveVIP("dummy", "192.168.1.101/24").RunAndReturn(func(_, _ string) error {
		wg.Done()
		return nil
	})
	s.ticker.Tick(time.Now())
	s.True(wg.WaitWithTimeout(1 * time.Second))
}

func TestManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}
