package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
)

func init() {
	RootCmd.AddCommand(diffCmd)
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show diff between your local and remote infra",

	RunE: func(cmd *cobra.Command, args []string) error {
		localInfra, err := rdf.NewNamedGraphFromFile("localInfra", filepath.Join(config.Dir, config.InfraFilename))
		if err != nil {
			return err
		}

		awsInfra, err := infraApi.FetchAwsInfra()
		if err != nil {
			return err
		}
		remoteInfra, err := rdf.BuildAwsInfraGraph("remoteInfra", viper.GetString("region"), awsInfra)
		if err != nil {
			return err
		}

		extras, missings, err := rdf.Compare(viper.GetString("region"), localInfra, remoteInfra)
		if err != nil {
			return err
		}

		const padding = 5
		w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
		fmt.Fprintln(w, "INFRA")
		fmt.Fprintln(w, "Extras:")
		fmt.Fprintln(w, extras.FlushString())
		fmt.Fprintln(w, "Missings:")
		fmt.Fprintln(w, missings.FlushString())

		localAccess, err := rdf.NewNamedGraphFromFile("localAccess", filepath.Join(config.Dir, config.AccessFilename))
		if err != nil {
			return err
		}

		access, err := accessApi.FetchAwsAccess()
		if err != nil {
			return err
		}
		remoteAccess, err := rdf.BuildAwsAccessGraph("remoteAccess", viper.GetString("region"), access)
		if err != nil {
			return err
		}

		extras, missings, err = rdf.Compare(viper.GetString("region"), localAccess, remoteAccess)
		if err != nil {
			return err
		}

		fmt.Fprintln(w)
		fmt.Fprintln(w, "ACCESS:")
		fmt.Fprintln(w, "Extras:")
		fmt.Fprintln(w, extras.FlushString())
		fmt.Fprintln(w, "Missings:")
		fmt.Fprintln(w, missings.FlushString())

		w.Flush()

		return nil
	},
}
