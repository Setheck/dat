package cmd

import (
	"fmt"
	"github.com/go-delve/delve/cmd/dlv/cmds"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:                        "dat",
	Short:                      "dat is a simple tool for formatting dates and times",
	Long:                       `Format dem dates and times.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(args)
		return nil
	},
}

func Execute() {
	if err := cmds.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}