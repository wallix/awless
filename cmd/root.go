package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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
	viper.SetConfigName("config")            // name of config file (without extension)
	viper.AddConfigPath("$HOME/.awless/aws") // adding home directory as first search path

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
		return
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("region"))})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
		return
	}
	if _, err = sess.Config.Credentials.Get(); err != nil {
		fmt.Println(err)
		fmt.Fprintln(os.Stderr, "Your AWS credentials seem undefined!")
		os.Exit(-1)
		return
	}

	accessApi = api.NewAccess(sess)
	infraApi = api.NewInfra(sess)
}
