package cmd

import (
	"fardjad.com/dqlite-vip/api"
	"fardjad.com/dqlite-vip/cluster"
	"fardjad.com/dqlite-vip/version"
	"github.com/spf13/cobra"
)

type Root struct{}

func (c *Root) run(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

func (root *Root) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "dqlite-vip",
		Short:        "Dqlite based IP load balancer",
		SilenceUsage: true,
		Version:      version.Version(),
	}
	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.Run = root.run

	start := &start{
		waiter:                      &SigTermWaiter{},
		clusterNodeFactoryFunc:      cluster.NewClusterNode,
		backgroundServerFactoryFunc: api.NewBackgroundHttpServer,
	}
	cmd.AddCommand(start.command())

	return cmd
}
