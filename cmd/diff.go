package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/api"
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
		var awsInfra *api.AwsInfra
		var awsAccess *api.AwsAccess

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			infra, err := infraApi.FetchAwsInfra()
			exitOn(err)
			awsInfra = infra
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			access, err := accessApi.FetchAwsAccess()
			exitOn(err)
			awsAccess = access
		}()

		wg.Wait()

		localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.Dir, config.InfraFilename))
		if err != nil {
			return err
		}

		remoteInfra, err := rdf.BuildAwsInfraGraph(viper.GetString("region"), awsInfra)
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
		fmt.Fprintln(w, extras.MustMarshal())
		fmt.Fprintln(w, "Missings:")
		fmt.Fprintln(w, missings.MustMarshal())

		localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.Dir, config.AccessFilename))
		if err != nil {
			return err
		}

		remoteAccess, err := rdf.BuildAwsAccessGraph(viper.GetString("region"), awsAccess)
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
		fmt.Fprintln(w, extras.MustMarshal())
		fmt.Fprintln(w, "Missings:")
		fmt.Fprintln(w, missings.MustMarshal())

		w.Flush()

		return nil
	},
}
