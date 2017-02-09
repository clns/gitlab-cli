package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var fromRepo string

var labelCopyCmd = &cobra.Command{
	Use:     "copy",
	Aliases: []string{"c"},
	Short:   "Copy labels into a repository",
	Long: `Copy labels into a repository.

If --from is omitted, it will copy global labels. If --from is specified,
it will copy all labels from that repository.

The from repo can be a repo name as in the config file or a relative path
as group/repo (e.g. 'myuser/myrepo'). In the later case it will use the url
of the target repo, so the repositories need to be on the same GitLab instance.`,
	Example: `  $ gitlab label copy -U https://gitlab.com/user/myrepo -t <TOKEN>
  $ gitlab label copy --from sourceRepo -r targetRepo
  $ gitlab label copy --from group/repo -r targetRepo`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			from, to *Repo
			err      error
		)
		if to, err = LoadFromConfig(repo); err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid target repository: %v\n", err.Error())
			os.Exit(1)
		}
		if fromRepo != "" {
			if from, err = LoadFromConfig(fromRepo); err != nil {
				fmt.Fprintf(os.Stderr, "error: invalid source repository: %v\n", err.Error())
				os.Exit(1)
			}
		}

		if from == nil {
			// we need to copy the global labels
			if err := to.Client.Labels.CopyGlobalLabelsTo(to.Project.ID); err != nil {
				fmt.Fprintf(os.Stderr, "error: '%s': %v\n",
					to.Project.PathWithNamespace, err)
				os.Exit(1)
			}
		} else {
			// we need to copy labels from one project to another
			if err := to.Client.Labels.CopyLabels(from.Project.ID, to.Project.ID); err != nil {
				fmt.Fprintf(os.Stderr, "error: '%s' to '%s': %v\n",
					from.Project.PathWithNamespace, to.Project.PathWithNamespace, err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	labelCmd.AddCommand(labelCopyCmd)

	labelCopyCmd.Flags().StringVar(&fromRepo, "from", "", "Source repository (optional)")
}
