package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/api"
	"github.com/wallix/awless/config"
)

var (
	accessApi *api.Access
	infraApi  *api.Infra
)

var RootCmd = &cobra.Command{
	Use:   "awless",
	Short: "Manage your cloud",
	Long:  "Awless is a CLI to ....:",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile(config.Path)

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
