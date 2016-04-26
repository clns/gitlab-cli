package cmd

import "github.com/spf13/cobra"

var configRepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Config repos actions",
	Long:  `Perform actions on the repos from the config file.`,
}

func init() {
	configCmd.AddCommand(configRepoCmd)
}
