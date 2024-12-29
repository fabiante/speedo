package main

import (
	"fmt"
	"os"

	"github.com/fabiante/speedo/cmd/speedo/cmds"
	"github.com/spf13/cobra"
)

func main() {
	cli := &cobra.Command{
		Use: "speedo",
	}

	cli.AddCommand(cmds.NewCSV())
	cli.AddCommand(cmds.NewRun())

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
