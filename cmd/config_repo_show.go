package cmd

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

var configRepoShowCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show repo info from the config file",
	Example: `  $ gitlab config repo show -r myrepo`,
	Run: func(cmd *cobra.Command, args []string) {
		if repo == "" {
			fmt.Fprintf(os.Stderr, "error: no repo name given\n")
			os.Exit(1)
		}
		r, err := LoadFromConfig(repo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid repository: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(r.String())
	},
}

func init() {
	configRepoCmd.AddCommand(configRepoShowCmd)
}
