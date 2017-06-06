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
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
)

var (
	servicesToSyncFlags map[string]*bool
)

func init() {
	RootCmd.AddCommand(syncCmd)

	servicesToSyncFlags = make(map[string]*bool)
	for _, service := range aws.ServiceNames {
		servicesToSyncFlags[service] = new(bool)
		syncCmd.Flags().BoolVar(servicesToSyncFlags[service], service, false, fmt.Sprintf("Sync '%s' service only", service))
	}
}

var syncCmd = &cobra.Command{
	Use:               "sync",
	Short:             "Manual sync of your remote resources to your local rdf store. For example when auto sync unset",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade),

	RunE: func(cmd *cobra.Command, args []string) error {
		if extraVerboseGlobalFlag {
			logger.DefaultLogger.SetVerbose(logger.ExtraVerboseF)
		} else {
			logger.DefaultLogger.SetVerbose(logger.VerboseF) //Forcing verbose to display sync info
		}

		var services []cloud.Service
		displayAllServices := true
		for _, srv := range cloud.ServiceRegistry {
			if *servicesToSyncFlags[srv.Name()] {
				displayAllServices = false
			}
		}
		for _, srv := range cloud.ServiceRegistry {
			if displayAllServices || *servicesToSyncFlags[srv.Name()] {
				services = append(services, srv)
			}
		}
		localGraphs := make(map[string]*graph.Graph)
		for _, service := range services {
			localGraphs[service.Name()] = sync.LoadCurrentLocalGraph(service.Name(), config.GetAWSRegion())
		}
		logger.Info("running sync: fetching remote resources for local store")
		start := time.Now()

		graphs, err := sync.DefaultSyncer.Sync(services...)
		if err != nil {
			logger.Verbose(err)
		}

		for k, g := range graphs {
			displaySyncStats(k, g)
		}
		logger.Infof("sync took %s", time.Since(start))

		return nil
	},
}

func displaySyncStats(serviceName string, g *graph.Graph) {
	var strs []string
	for rt, service := range aws.ServicePerResourceType {
		if service == serviceName {
			res, err := g.GetAllResources(rt)
			if err != nil {
				continue
			}
			nbRes := len(res)
			if nbRes > 1 {
				strs = append(strs, fmt.Sprintf("%d %s", nbRes, cloud.PluralizeResource(rt)))
			} else {
				strs = append(strs, fmt.Sprintf("%d %s", nbRes, rt))
			}
		}
	}
	logger.Infof("-> %s: %s", serviceName, strings.Join(strs, ", "))
}
