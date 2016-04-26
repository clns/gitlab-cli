package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configRepoLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List repositories from config file",
	Run: func(cmd *cobra.Command, args []string) {
		for name, _ := range viper.GetStringMap("repos") {
			r := LoadFromConfigNoInit(name)
			fmt.Println(r.String())
		}
	},
}

func init() {
	configRepoCmd.AddCommand(configRepoLsCmd)
}
