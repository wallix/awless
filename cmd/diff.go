package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/google/badwolf/triple"
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
		localInfra, err := triplesFromFile(config.InfraFilename)
		if err != nil {
			return err
		}
		vpcs, subnets, instances, err := infraApi.FetchInfra()
		if err != nil {
			return err
		}

		remoteInfra, err := rdf.BuildInfraRdfTriples(viper.GetString("region"), vpcs, subnets, instances)
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
		fmt.Fprintln(w, rdf.MarshalTriples(extras))
		fmt.Fprintln(w, "Missings:")
		fmt.Fprintln(w, rdf.MarshalTriples(missings))

		localAccess, err := triplesFromFile(config.AccessFilename)
		if err != nil {
			return err
		}

		groups, users, usersByGroup, err := accessApi.FetchAccess()
		if err != nil {
			return err
		}

		remoteAccess, err := rdf.BuildAccessRdfTriples(viper.GetString("region"), groups, users, usersByGroup)
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
		fmt.Fprintln(w, rdf.MarshalTriples(extras))
		fmt.Fprintln(w, "Missings:")
		fmt.Fprintln(w, rdf.MarshalTriples(missings))

		w.Flush()

		return nil
	},
}

func triplesFromFile(filename string) ([]*triple.Triple, error) {
	if content, err := ioutil.ReadFile(filepath.Join(config.Dir, filename)); err != nil {
		return nil, err
	} else {
		return rdf.UnmarshalTriples(string(content))
	}
}
