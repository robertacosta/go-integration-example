package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{}

func init() {
	RootCmd.AddCommand(serverCmd)
	RootCmd.AddCommand(testServerCmd)
	RootCmd.AddCommand(workerCmd)
}
