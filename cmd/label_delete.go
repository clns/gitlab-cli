package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var regexpLabel string

var labelDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d"},
	Short:   "Delete labels from a repository",
	Long: `Delete labels from a repository.

The --regex flag can be specified as a Go regexp pattern to delete only
labels that match. If ommitted, all repository labels will be deleted.`,
	Example: `  $ gitlab label delete -r myrepo
  $ gitlab label delete -r myrepo --regexp=".*:.*"`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			to  *Repo
			err error
		)
		if to, err = LoadFromConfig(repo); err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid repository: %v\n", err.Error())
			os.Exit(1)
		}

		if err := to.Client.Labels.DeleteWithRegex(to.Project.ID, regexpLabel); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	labelCmd.AddCommand(labelDeleteCmd)

	labelDeleteCmd.Flags().StringVar(&regexpLabel, "regex", "", "Label name to match, as a Go regex (https://golang.org/pkg/regexp/syntax)")
}
