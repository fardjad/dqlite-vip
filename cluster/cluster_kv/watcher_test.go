package cluster_kv

import (
	"testing"
	stdTime "time"

	"github.com/stretchr/testify/suite"

	clusterMocks "fardjad.com/dqlite-vip/mocks/cluster"
	"fardjad.com/dqlite-vip/time"
)

type WatcherTestSuite struct {
	suite.Suite
	clusterNode *clusterMocks.ClusterNode
	ticker      *time.FakeTicker
	watcher     *Watcher
}

func (s *WatcherTestSuite) SetupTest() {
	s.clusterNode = clusterMocks.NewClusterNode(s.T())
	s.ticker = time.NewFakeTicker()
	s.watcher = NewWatcher(s.clusterNode, s.ticker)
}

func (s *WatcherTestSuite) TestWatch() {

	ch, cancel, err := s.watcher.Watch("vip")
	s.NoError(err)
	defer cancel()

	s.clusterNode.EXPECT().GetString("vip").Return("192.168.1.100", nil).Once()
	s.ticker.Tick(stdTime.Now())

	select {
	case change := <-ch:
		s.Equal("", change.Previous)
		s.Equal("192.168.1.100", change.Current)
	case <-stdTime.After(1 * stdTime.Second):
		s.Fail("Timeout waiting for value")
	}

	s.clusterNode.EXPECT().GetString("vip").Return("192.168.1.101", nil).Once()
	s.ticker.Tick(stdTime.Now())

	select {
	case change := <-ch:
		s.Equal("192.168.1.100", change.Previous)
		s.Equal("192.168.1.101", change.Current)
	case <-stdTime.After(1 * stdTime.Second):
		s.Fail("Timeout waiting for value")
	}
}

func TestWatcherTestSuite(t *testing.T) {
	suite.Run(t, new(WatcherTestSuite))
}
