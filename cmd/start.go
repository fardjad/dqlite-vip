package cmd

import (
	"context"
	"time"

	"fardjad.com/dqlite-vip/api"
	"fardjad.com/dqlite-vip/cluster"
	"fardjad.com/dqlite-vip/utils"
	"fardjad.com/dqlite-vip/vip"
	"github.com/spf13/cobra"
)

type start struct {
	waiter                      Waiter
	clusterNodeFactoryFunc      cluster.ClusterNodeFactoryFunc
	backgroundServerFactoryFunc api.BackgroundServerFactoryFunc
	tickerFactoryFunc           utils.BetterTickerFactoryFunc
	configuratorFactoryFunc     vip.ConfiguratorFactoryFunc
	vipManagerFactoryFunc       vip.ManagerFactoryFunc

	// flags
	dataDir     string
	bindCluster string
	bindHttp    string
	join        string
	iface       string
}

func (c *start) NewVIPManager(clusterNode cluster.ClusterNode) vip.Manager {
	ticker := c.tickerFactoryFunc(time.Second)
	configurator := c.configuratorFactoryFunc(5)
	return c.vipManagerFactoryFunc(clusterNode, ticker, configurator, c.iface)
}

func (c *start) runE(cmd *cobra.Command, args []string) error {
	join := []string{}
	if c.join != "" {
		join = append(join, c.join)
	}
	clusterNode, err := c.clusterNodeFactoryFunc(c.dataDir, c.bindCluster, join)
	if err != nil {
		return err
	}

	handlers := api.NewHandlers(clusterNode)
	server := c.backgroundServerFactoryFunc(c.bindHttp, handlers.Mux())

	err = server.ListenAndServeInBackground()
	if err != nil {
		return err
	}
	defer server.Shutdown(context.Background())

	err = clusterNode.Start(context.Background())
	if err != nil {
		return err
	}
	defer clusterNode.Close(context.Background())

	vipManager := c.NewVIPManager(clusterNode)
	vipManager.Start(context.Background())
	defer vipManager.Stop()

	c.waiter.Wait()

	return nil
}

func (start *start) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a dqlite-vip node",
	}
	cmd.RunE = start.runE

	cmd.Flags().StringVarP(&start.dataDir, "data-dir", "d", "", "[required] directory where the node data is stored")
	cmd.MarkFlagRequired("data-dir")

	cmd.Flags().StringVar(&start.bindCluster, "bind-cluster", "", "[required] address to bind the cluster connection port to")
	cmd.MarkFlagRequired("bind-cluster")

	cmd.Flags().StringVar(&start.bindHttp, "bind-http", "", "[required] address to bind the HTTP API to")
	cmd.MarkFlagRequired("bind-http")

	cmd.Flags().StringVarP(&start.join, "join", "j", "", "address of an existing dqlite-vip node to join")

	cmd.Flags().StringVarP(&start.iface, "iface", "i", "", "Name of the network interface to use for the VIP")
	cmd.MarkFlagRequired("iface")

	return cmd
}
