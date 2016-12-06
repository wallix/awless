package cmd

import (
	"github.com/spf13/cobra"
	awscloud "github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/scenario"
	"github.com/wallix/awless/scenario/driver"
	"github.com/wallix/awless/scenario/driver/aws"
)

func init() {
	createCmd.AddCommand(createInstanceCmd)

	RootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create various type of resources by id: users, groups, instances, vpcs, ...",
}

var createInstanceCmd = &cobra.Command{
	Use:     "instance",
	Aliases: []string{"inst", "i"},
	Short:   "Create an instance",

	RunE: func(cmd *cobra.Command, args []string) error {
		raw := `CREATE VPC CIDR 10.0.0.0/16 REF vpc_1
CREATE SUBNET CIDR 10.0.0.0/16 REF subnet_1 REFERENCES vpc_1
CREATE INSTANCE COUNT 1 REFERENCES subnet_1
`
		lex := &scenario.Lexer{}
		scen := lex.ParseScenario(raw)

		runner := &driver.Runner{aws.NewAwsDriver(awscloud.InfraService)}

		return runner.Run(scen)
	},
}
