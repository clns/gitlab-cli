package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	gogitlab "github.com/xanzy/go-gitlab"
)

var matchLabel string
var replaceLabel string
var colorLabel string

var labelUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "Update labels in a repository",
	Long: `Update labels in a repository.

The --match flag is required and is a Go regex that will be used to match the label
name. At least one of --replace or --color is required to update the label(s).`,
	Example: `  $ gitlab label update -r myrepo --match "(.*):(.*)" --replace "${1}/${2}"`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			to  *Repo
			err error
		)
		if to, err = LoadFromConfig(repo); err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid repository: %v\n", err.Error())
			os.Exit(1)
		}

		if err := to.Client.Labels.UpdateWithRegex(*to.Project.ID, &gogitlab.UpdateLabelOptions{
			Name:    matchLabel,
			NewName: replaceLabel,
			Color:   colorLabel,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	labelCmd.AddCommand(labelUpdateCmd)

	labelUpdateCmd.Flags().StringVar(&matchLabel, "match", "", "Label name to match, as a Go regex (https://golang.org/pkg/regexp/syntax)")
	labelUpdateCmd.Flags().StringVar(&replaceLabel, "replace", "", "Label name replacement (https://golang.org/pkg/regexp/#Regexp.FindAllString)")
	labelUpdateCmd.Flags().StringVar(&colorLabel, "color", "", "Label color (e.g. '#000000')")
}
