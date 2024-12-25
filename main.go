package main

import (
	"os"

	"fardjad.com/dqlite-vip/cmd"
)

func main() {
	root := &cmd.Root{}
	if err := root.Command().Execute(); err != nil {
		os.Exit(1)
	}
}
