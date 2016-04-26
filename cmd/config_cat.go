package cmd

import (
	"fmt"

	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCatCmd = &cobra.Command{
	Use:   "cat",
	Short: "Print config file",
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile(viper.ConfigFileUsed())
		if err != nil {
			os.Exit(1)
		}
		fmt.Fprintln(os.Stdout, string(b))
	},
}

func init() {
	configCmd.AddCommand(configCatCmd)
}
