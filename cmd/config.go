package cmd

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config actions",
	Long:  `Perform config-related actions.`,
}

func init() {
	RootCmd.AddCommand(configCmd)
}
