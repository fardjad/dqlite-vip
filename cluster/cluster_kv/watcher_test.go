package cluster_kv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	clusterMocks "fardjad.com/dqlite-vip/mocks/cluster"
)

type WatcherTestSuite struct {
	suite.Suite
	clusterNode *clusterMocks.ClusterNode
	watcher     *Watcher
}

func (s *WatcherTestSuite) SetupTest() {
	s.clusterNode = clusterMocks.NewClusterNode(s.T())
	ticker := time.NewTicker(1 * time.Millisecond)
	s.watcher = NewWatcher(s.clusterNode, ticker)
}

func (s *WatcherTestSuite) TestWatch() {
	s.clusterNode.EXPECT().GetString("vip").Return("192.168.1.100", nil).Once()

	ch, cancel, err := s.watcher.Watch("vip")
	s.NoError(err)
	defer cancel()

	select {
	case value := <-ch:
		s.Equal("192.168.1.100", value)
	case <-time.After(1 * time.Second):
		s.Fail("Timeout waiting for value")
	}

	s.clusterNode.EXPECT().GetString("vip").Return("192.168.1.101", nil).Once()
	select {
	case value := <-ch:
		s.Equal("192.168.1.101", value)
	case <-time.After(1 * time.Second):
		s.Fail("Timeout waiting for value")
	}
}

func TestWatcherTestSuite(t *testing.T) {
	suite.Run(t, new(WatcherTestSuite))
}
