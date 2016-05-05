package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.2.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version of this tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gitlab-cli %v\n", Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
