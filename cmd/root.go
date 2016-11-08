package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/api"
)

var (
	configPath string

	accessApi *api.Access
	infraApi  *api.Infra
)

var RootCmd = &cobra.Command{
	Use:   "awless",
	Short: "Manage your cloud",
	Long:  "Awless is a CLI to ....:",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("config")        // name of config file (without extension)
	viper.AddConfigPath("$HOME/.awless") // adding home directory as first search path

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	var err error

	if accessApi, err = api.NewAccess(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to init the access api: %s\n", err)
		os.Exit(-1)
	}
	if infraApi, err = api.NewInfra(viper.GetString("region")); err != nil {
		fmt.Fprintf(os.Stderr, "unable to init the infra api: %s\n", err)
		os.Exit(-1)
	}
}
