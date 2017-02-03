package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var labelCopyCmd = &cobra.Command{
	Use:     "copy [<from-repo>]",
	Aliases: []string{"c"},
	Short:   "Copy labels into a repository",
	Long: `Copy labels into a repository.

If given without an argument, it will copy global labels. If <from-repo>
is specified, it will copy all labels from that repository. This argument
can be a repo name as in the config file or a path (e.g. 'myuser/myrepo').
In the later case it will use --url, without its path.`,
	Example: `  $ gitlab label copy -U https://gitlab.com/user/myrepo -t <TOKEN>
  # = copy global labels into 'user/myrepo'

  $ gitlab label copy -r myrepo myotherrepo
  # = copy labels from 'myotherrepo' into 'myrepo'`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Fprintln(os.Stderr, "error: invalid arguments")
			cmd.Usage()
			os.Exit(2)
		}
		var (
			from, to *Repo
			err      error
		)
		if len(args) == 0 {
			if to, err = LoadFromConfig(repo); err != nil {
				fmt.Fprintf(os.Stderr, "error: invalid target repository: %v\n", err.Error())
				os.Exit(1)
			}
		} else if len(args) == 1 {
			if to, err = LoadFromConfig(repo); err != nil {
				fmt.Fprintf(os.Stderr, "error: invalid target repository: %v\n", err.Error())
				os.Exit(1)
			}
			if from, err = LoadFromConfig(args[0]); err != nil {
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
}
