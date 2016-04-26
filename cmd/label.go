package cmd

import "github.com/spf13/cobra"

var labelCmd = &cobra.Command{
	Use:     "label",
	Aliases: []string{"l"},
	Short:   "Label actions",
	Long:    `Perform actions on labels.`,
}

func init() {
	RootCmd.AddCommand(labelCmd)
}
