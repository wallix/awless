package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/api"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/stats"
)

var (
	accessApi *api.Access
	infraApi  *api.Infra
	statsDB   *stats.DB

	verboseFlag bool
)

var RootCmd = &cobra.Command{
	Use:   "awless",
	Short: "Manage your cloud",
	Long:  "Awless is a powerful command line tool to inspect, sync and manage your infrastructure",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if statsDB != nil {
			defer statsDB.Close()
			statsDB.AddHistoryCommand(append(strings.Split(cmd.CommandPath(), " "), args...))
		}
	},
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
	RootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Turn on verbose mode for all commands")

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

	statsDB, err = stats.OpenDB(config.DatabasePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can not save history:", err)
	} else if statsDB.CheckStatsToSend(config.StatsExpirationDuration) {
		publicKey, err := config.LoadPublicKey()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			if !config.AwlessFirstSync {
				if err := statsDB.SendStats(config.StatsServerUrl, *publicKey); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}
	}
}
