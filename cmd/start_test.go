package cmd

import (
	"context"
	"testing"

	"fardjad.com/dqlite-vip/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type StartTestSuite struct {
	suite.Suite

	waiter                  *mocks.Waiter
	clusterNode             *mocks.ClusterNode
	clusterNodeFactory      *mocks.ClusterNodeFactory
	backgroundServer        *mocks.BackgroundServer
	backgroundServerFactory *mocks.BackgroundServerFactory
}

func (s *StartTestSuite) SetupTest() {
	s.waiter = mocks.NewWaiter(s.T())
	s.waiter.EXPECT().Wait().Return()

	s.clusterNode = mocks.NewClusterNode(s.T())
	s.clusterNode.EXPECT().Start().Return(nil)
	s.clusterNode.EXPECT().Close().Return(nil)

	s.clusterNodeFactory = mocks.NewClusterNodeFactory(s.T())

	s.backgroundServer = mocks.NewBackgroundServer(s.T())
	s.backgroundServer.EXPECT().ListenAndServeInBackground().Return(nil)
	s.backgroundServer.EXPECT().Shutdown(context.Background()).Return(nil)

	s.backgroundServerFactory = mocks.NewBackgroundServerFactory(s.T())
}

func (s *StartTestSuite) TestStart() {
	s.clusterNodeFactory.EXPECT().NewClusterNode(
		"/tmp/data",
		"127.0.0.1:9000",
		[]string{},
	).Return(s.clusterNode, nil)
	s.backgroundServerFactory.EXPECT().NewServer("127.0.0.1:8080", mock.Anything).Return(s.backgroundServer)

	startCmd := &start{
		waiter:                  s.waiter,
		clusterNodeFactory:      s.clusterNodeFactory,
		backgroundServerFactory: s.backgroundServerFactory,
	}

	cmd := startCmd.command()
	cmd.SetArgs([]string{
		"--data-dir", "/tmp/data",
		"--bind-cluster", "127.0.0.1:9000",
		"--bind-http", "127.0.0.1:8080",
	})

	err := cmd.Execute()
	s.NoError(err, "Expected no error, but got: %v", err)

	s.waiter.AssertExpectations(s.T())
	s.clusterNode.AssertExpectations(s.T())
	s.clusterNodeFactory.AssertExpectations(s.T())
	s.backgroundServer.AssertExpectations(s.T())
	s.backgroundServerFactory.AssertExpectations(s.T())
}

func (s *StartTestSuite) TestStartWithJoin() {
	s.clusterNodeFactory.EXPECT().NewClusterNode(
		"/tmp/data",
		"127.0.0.1:9000",
		[]string{"127.0.0.1:9001"},
	).Return(s.clusterNode, nil)
	s.backgroundServerFactory.EXPECT().NewServer("127.0.0.1:8080", mock.Anything).Return(s.backgroundServer)

	startCmd := &start{
		waiter:                  s.waiter,
		clusterNodeFactory:      s.clusterNodeFactory,
		backgroundServerFactory: s.backgroundServerFactory,
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

	s.waiter.AssertExpectations(s.T())
	s.clusterNode.AssertExpectations(s.T())
	s.clusterNodeFactory.AssertExpectations(s.T())
	s.backgroundServer.AssertExpectations(s.T())
	s.backgroundServerFactory.AssertExpectations(s.T())
}

func TestStartTestSuite(t *testing.T) {
	suite.Run(t, new(StartTestSuite))
}
