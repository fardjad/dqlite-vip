package cmd

import (
	"context"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	apiMocks "fardjad.com/dqlite-vip/mocks/api"
	clusterMocks "fardjad.com/dqlite-vip/mocks/cluster"
	cmdMocks "fardjad.com/dqlite-vip/mocks/cmd"
	utilsMocks "fardjad.com/dqlite-vip/mocks/utils"
	vipMocks "fardjad.com/dqlite-vip/mocks/vip"
	"fardjad.com/dqlite-vip/utils"
)

const DATA_DIR = "/tmp/data"
const BIND_CLUSTER = "127.0.0.1:8000"
const BIND_HTTP = "127.0.0.1:9000"
const IFACE = "iface"

type StartTestSuite struct {
	suite.Suite

	clusterNode                 *clusterMocks.ClusterNode
	clusterNodeFactoryFunc      *clusterMocks.ClusterNodeFactoryFunc
	backgroundServerFactoryFunc *apiMocks.BackgroundServerFactoryFunc
	vipManagerFactoryFunc       *vipMocks.ManagerFactoryFunc

	cmd *cobra.Command
}

func (s *StartTestSuite) SetupTest() {
	waiter := cmdMocks.NewWaiter(s.T())
	waiter.EXPECT().Wait().Return()

	s.clusterNode = clusterMocks.NewClusterNode(s.T())
	s.clusterNode.EXPECT().Start(mock.Anything).Return(nil)
	s.clusterNode.EXPECT().Close(mock.Anything).Return(nil)
	s.clusterNodeFactoryFunc = clusterMocks.NewClusterNodeFactoryFunc(s.T())

	backgroundServer := apiMocks.NewBackgroundServer(s.T())
	backgroundServer.EXPECT().ListenAndServeInBackground().Return(nil)
	backgroundServer.EXPECT().Shutdown(context.Background()).Return(nil)
	s.backgroundServerFactoryFunc = apiMocks.NewBackgroundServerFactoryFunc(s.T())
	s.backgroundServerFactoryFunc.EXPECT().Execute(BIND_HTTP, mock.Anything).Return(backgroundServer)

	ticker := utils.NewFakeTicker()
	tickerFactoryFunc := utilsMocks.NewBetterTickerFactoryFunc(s.T())
	tickerFactoryFunc.EXPECT().Execute(time.Second).Return(ticker)

	configurator := vipMocks.NewConfigurator(s.T())
	configuratorFactoryFunc := vipMocks.NewConfiguratorFactoryFunc(s.T())
	configuratorFactoryFunc.EXPECT().Execute(5).Return(configurator)

	vipManager := vipMocks.NewManager(s.T())
	vipManager.EXPECT().Start(mock.Anything).Return()
	vipManager.EXPECT().Stop().Return()
	s.vipManagerFactoryFunc = vipMocks.NewManagerFactoryFunc(s.T())
	s.vipManagerFactoryFunc.EXPECT().Execute(s.clusterNode, ticker, configurator, IFACE).Return(vipManager)

	startCmd := &start{
		waiter:                      waiter,
		clusterNodeFactoryFunc:      s.clusterNodeFactoryFunc.Execute,
		backgroundServerFactoryFunc: s.backgroundServerFactoryFunc.Execute,
		tickerFactoryFunc:           tickerFactoryFunc.Execute,
		configuratorFactoryFunc:     configuratorFactoryFunc.Execute,
		vipManagerFactoryFunc:       s.vipManagerFactoryFunc.Execute,
	}
	s.cmd = startCmd.command()
}

func (s *StartTestSuite) TestStart() {
	s.clusterNodeFactoryFunc.EXPECT().Execute(
		DATA_DIR,
		BIND_CLUSTER,
		[]string{},
	).Return(s.clusterNode, nil)

	s.cmd.SetArgs([]string{
		"--data-dir", DATA_DIR,
		"--bind-cluster", BIND_CLUSTER,
		"--bind-http", BIND_HTTP,
		"--iface", IFACE,
	})

	s.NoError(s.cmd.Execute())
}

func (s *StartTestSuite) TestStartWithJoin() {
	s.clusterNodeFactoryFunc.EXPECT().Execute(
		DATA_DIR,
		BIND_CLUSTER,
		[]string{"127.0.0.1:9001"},
	).Return(s.clusterNode, nil)

	s.cmd.SetArgs([]string{
		"--data-dir", DATA_DIR,
		"--bind-cluster", BIND_CLUSTER,
		"--bind-http", BIND_HTTP,
		"--join", "127.0.0.1:9001",
		"--iface", IFACE,
	})

	s.NoError(s.cmd.Execute())
}

func TestStartTestSuite(t *testing.T) {
	suite.Run(t, new(StartTestSuite))
}
