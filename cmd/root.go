package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
	"github.com/wallix/awless/stats"
)

var (
	statsDB *stats.DB

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

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Turn on verbose mode for all commands")

	var err error

	statsDB, err = stats.OpenDB(config.DatabasePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can not save history:", err)
	} else if statsDB.CheckStatsToSend(config.StatsExpirationDuration) {
		publicKey, err := config.LoadPublicKey()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			if !config.AwlessFirstSync {
				localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.InfraFilename))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				if err := statsDB.SendStats(config.StatsServerUrl, *publicKey, localInfra, localAccess); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}
	}
}
