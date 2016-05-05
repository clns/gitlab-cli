package cmd

import (
	"fmt"
	"os"

	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile              string
	repo, repourl, token string
	user, password       string
	verbose              bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gitlab-cli",
	Short: "GitLab CLI tool",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	helpFunc := RootCmd.HelpFunc()
	RootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		initVerbose()
		CheckUpdate()
		helpFunc(cmd, args)
	})
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initVerbose)
	cobra.OnInitialize(CheckUpdate)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gitlab-cli.yaml)")

	// Repository flags
	RootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "repo name (as in the config file)")
	RootCmd.PersistentFlags().StringVarP(&repourl, "url", "U", "", "repository URL, including the path (e.g. https://mygitlab.com/group/repo)")
	RootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "GitLab token (see http://doc.gitlab.com/ce/api/#authentication)")
	RootCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "GitLab login (user or email), if no token provided")
	RootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "GitLab password, if no token provided (if empty, will prompt)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "print logs")

	viper.BindPFlag("_url", RootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("_token", RootCmd.PersistentFlags().Lookup("token"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".gitlab-cli") // name of config file (without extension)
	viper.AddConfigPath("$HOME")       // adding home directory as first search path
	viper.AutomaticEnv()               // read in environment variables that match
	viper.ConfigFileUsed()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initVerbose() {
	if !verbose {
		log.SetOutput(ioutil.Discard)
	}
}
