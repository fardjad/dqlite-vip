package cluster_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"fardjad.com/dqlite-vip/cluster"
	"github.com/stretchr/testify/suite"
)

type DqliteClusterSuite struct {
	suite.Suite
	node1  cluster.ClusterNode
	node2  cluster.ClusterNode
	node3  cluster.ClusterNode
	ctx    context.Context
	cancel context.CancelFunc
}

func getFreePort() (string, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return fmt.Sprintf("127.0.0.1:%d", listener.Addr().(*net.TCPAddr).Port), nil
}

func (s *DqliteClusterSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 10*time.Second)

	node1DataDir := s.T().TempDir()
	node2DataDir := s.T().TempDir()
	node3DataDir := s.T().TempDir()

	port1, err := getFreePort()
	s.Require().NoError(err)
	port2, err := getFreePort()
	s.Require().NoError(err)
	port3, err := getFreePort()
	s.Require().NoError(err)

	node1, err := cluster.NewDqliteClusterNode(node1DataDir, port1, nil)
	s.Require().NoError(err)
	s.Require().NoError(node1.Start(s.ctx))
	s.node1 = node1

	node2, err := cluster.NewDqliteClusterNode(node2DataDir, port2, []string{port1})
	s.Require().NoError(err)
	s.Require().NoError(node2.Start(s.ctx))
	s.node2 = node2

	node3, err := cluster.NewDqliteClusterNode(node3DataDir, port3, []string{port1})
	s.Require().NoError(err)
	s.Require().NoError(node3.Start(s.ctx))
	s.node3 = node3

	s.Require().NoError(s.node1.Ready(s.ctx))
	s.Require().NoError(s.node2.Ready(s.ctx))
	s.Require().NoError(s.node3.Ready(s.ctx))
}

func (s *DqliteClusterSuite) TearDownTest() {
	s.node1.Close(s.ctx)
	s.node2.Close(s.ctx)
	s.node3.Close(s.ctx)
	s.cancel()
}

func (s *DqliteClusterSuite) TestKeyValueOperations() {
	key := "test-key"
	value := "test-value"

	s.ctx, s.cancel = context.WithTimeout(context.Background(), 10*time.Second)

	s.Require().NoError(s.node1.SetString(s.ctx, key, value))

	result, err := s.node1.GetString(s.ctx, key)
	s.Require().NoError(err)
	s.Equal(value, result)

	result, err = s.node2.GetString(s.ctx, key)
	s.Require().NoError(err)
	s.Equal(value, result)

	result, err = s.node3.GetString(s.ctx, key)
	s.Require().NoError(err)
	s.Equal(value, result)
}

func (s *DqliteClusterSuite) TestClusterMembers() {
	members, err := s.node1.ClusterMembers(s.ctx)
	s.Require().NoError(err)
	s.Len(members, 3)
}

func TestDqliteClusterSuite(t *testing.T) {
	suite.Run(t, new(DqliteClusterSuite))
}
