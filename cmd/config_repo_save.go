package cmd

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

var configRepoSaveCmd = &cobra.Command{
	Use:     "save",
	Short:   "Save repo into the config file",
	Example: `  $ gitlab config repo save -r myrepo -u https://gitlan.com/user/repo -t <TOKEN>`,
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
		r.Name = repo
		if err := r.SaveToConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := SaveViperConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	configRepoCmd.AddCommand(configRepoSaveCmd)
}
