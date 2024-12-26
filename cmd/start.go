package cmd

import (
	"context"

	"fardjad.com/dqlite-vip/api"
	"fardjad.com/dqlite-vip/cluster"
	"github.com/spf13/cobra"
)

type start struct {
	waiter                  Waiter
	clusterNodeFactory      cluster.ClusterNodeFactory
	backgroundServerFactory api.BackgroundServerFactory

	// flags
	dataDir     string
	bindCluster string
	bindHttp    string
	join        string
}

func (c *start) runE(cmd *cobra.Command, args []string) error {
	join := []string{}
	if c.join != "" {
		join = append(join, c.join)
	}
	clusterNode, err := c.clusterNodeFactory.NewClusterNode(c.dataDir, c.bindCluster, join)
	if err != nil {
		return err
	}

	handlers := api.NewHandlers(clusterNode)
	server := c.backgroundServerFactory.NewServer(c.bindHttp, handlers.Mux())

	server.ListenAndServeInBackground()
	defer server.Shutdown(context.Background())

	clusterNode.Start()
	defer clusterNode.Close()

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

	return cmd
}
