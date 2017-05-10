/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/logger"
)

var (
	showIdsOnlyFlag, showIdOnlyFlag bool
)

func init() {
	RootCmd.AddCommand(searchCmd)

	awsImagesCmd.Flags().BoolVar(&showIdsOnlyFlag, "ids-only", false, "Returns only the id of all AMIs matching your query")
	awsImagesCmd.Flags().BoolVar(&showIdOnlyFlag, "id-only", false, "Returns only one (the latest) AMI id matching your query")

	searchCmd.AddCommand(awsImagesCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Perform various searches, finding and resolution",
}

var awsImagesCmd = &cobra.Command{
	Use:              "awsimages",
	PersistentPreRun: applyHooks(initAwlessEnvHook, initLoggerHook, initCloudServicesHook),
	Short:            "Find corresponding images according to an image query",
	Long:             "Find corresponding images according to an image query. Query string specification:\n\n\t\tOWNER:DISTRO[VARIANT]:ARCH:VIRTUALIZATION:STORE\n\nEverything optional expect for the OWNER",
	Example:          "  awless search awsimages redhat:rhel[7.2]\n  awless search awsimages canonical --id-only\n  awless search awsimages amazonlinux::::instance-store --ids-only",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			exitOn(errors.New("malformed image query string. Expecting: OWNER:DISTRO[VARIANT]:ARCH:VIRTUALIZATION:STORE (with everything optional expect for the OWNER)"))
		}

		resolver := &aws.ImageResolver{aws.InfraService.(*aws.Infra)}

		query := args[0]
		logger.Infof("launching search for image query '%s'", query)
		imgs, err := resolver.Resolve(query)
		exitOn(err)

		var ids []string
		for _, img := range imgs {
			ids = append(ids, img.Id)
		}

		if showIdsOnlyFlag {
			for _, id := range ids {
				fmt.Println(id)
			}
			return
		}

		if showIdOnlyFlag {
			for i, id := range ids {
				fmt.Println(id)
				if i == 0 {
					break
				}
			}
			return
		}

		b, err := json.MarshalIndent(imgs, "", " ")
		exitOn(err)

		fmt.Fprintln(os.Stdout, string(b))
	},
}
