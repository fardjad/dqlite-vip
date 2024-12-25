package cmd

import (
	"testing"

	"fardjad.com/dqlite-vip/mocks"
	"github.com/stretchr/testify/suite"
)

type StartTestSuite struct {
	suite.Suite

	waiter      *mocks.Waiter
	clusterNode *mocks.ClusterNode
}

func (s *StartTestSuite) SetupTest() {
	s.waiter = mocks.NewWaiter(s.T())
	s.waiter.EXPECT().Wait().Return()

	s.clusterNode = mocks.NewClusterNode(s.T())
	s.clusterNode.EXPECT().Start().Return(nil)
	s.clusterNode.EXPECT().Close().Return(nil)
}

func (s *StartTestSuite) TestStart() {
	clusterNodeFactory := mocks.NewClusterNodeFactory(s.T())
	clusterNodeFactory.EXPECT().NewClusterNode(
		"/tmp/data",
		"127.0.0.1:9000",
		"127.0.0.1:8080",
		[]string{},
	).Return(s.clusterNode, nil)

	startCmd := &start{
		waiter:             s.waiter,
		clusterNodeFactory: clusterNodeFactory,
	}

	cmd := startCmd.command()
	cmd.SetArgs([]string{
		"--data-dir", "/tmp/data",
		"--bind-cluster", "127.0.0.1:9000",
		"--bind-http", "127.0.0.1:8080",
	})

	err := cmd.Execute()
	s.NoError(err, "Expected no error, but got: %v", err)

	clusterNodeFactory.AssertExpectations(s.T())
	s.waiter.AssertExpectations(s.T())
	s.clusterNode.AssertExpectations(s.T())
}

func (s *StartTestSuite) TestStartWithJoin() {
	clusterNodeFactory := mocks.NewClusterNodeFactory(s.T())
	clusterNodeFactory.EXPECT().NewClusterNode(
		"/tmp/data",
		"127.0.0.1:9000",
		"127.0.0.1:8080",
		[]string{"127.0.0.1:9001"},
	).Return(s.clusterNode, nil)

	startCmd := &start{
		waiter:             s.waiter,
		clusterNodeFactory: clusterNodeFactory,
	}

	cmd := startCmd.command()
	cmd.SetArgs([]string{
		"--data-dir", "/tmp/data",
		"--bind-cluster", "127.0.0.1:9000",
		"--bind-http", "127.0.0.1:8080",
		"--join", "127.0.0.1:9001",
	})

	err := cmd.Execute()
	s.NoError(err, "Expected no error, but got: %v", err)

	clusterNodeFactory.AssertExpectations(s.T())
	s.waiter.AssertExpectations(s.T())
	s.clusterNode.AssertExpectations(s.T())
}

func TestStartTestSuite(t *testing.T) {
	suite.Run(t, new(StartTestSuite))
}
