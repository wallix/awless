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
	"github.com/spf13/cobra"
	"github.com/wallix/awless/web"
)

var (
	webPortFlag string
)

func init() {
	RootCmd.AddCommand(webCmd)

	webCmd.Flags().StringVarP(&webPortFlag, "port", "p", ":8080", "Web UI port to listen on")
}

var webCmd = &cobra.Command{
	Use:    "web",
	Hidden: true,
	Short:  "Browse your cloud data through a web ui",

	Run: func(cmd *cobra.Command, args []string) {
		server := web.New(webPortFlag)
		exitOn(server.Start())
	},
}
