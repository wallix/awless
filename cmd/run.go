package cmd

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	awscloud "github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/script"
	"github.com/wallix/awless/script/driver/aws"
)

func init() {
	RootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run awless scripting. Either as one direct command or a given file",

	RunE: func(cmd *cobra.Command, args []string) error {
		var text string

		if len(args) < 1 {
			return errors.New("missing awless script file path or awless script line")
		}

		if len(args) == 1 {
			content, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}
			text = string(content)
		} else {
			text = strings.Join(args, " ")
		}

		scrpt, err := script.Parse(text)
		if err != nil {
			return err
		}

		defaults, err := database.Current.GetDefaults()
		if err != nil {
			return err
		}
		scrpt.ResolveTemplate(defaults)

		awsDriver := aws.NewDriver(awscloud.InfraService)
		awsDriver.SetLogger(log.New(os.Stdout, "[aws driver] ", log.Ltime))

		if _,err := scrpt.Compile(awsDriver); err != nil {
			return err
		}

		_, err = scrpt.Run(awsDriver)

		return err
	},
}
