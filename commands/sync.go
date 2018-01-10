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
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws/services"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
)

var (
	servicesToSyncFlags map[string]*bool
	profileSyncFlag     bool
)

func init() {
	RootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVar(&profileSyncFlag, "profile-sync", false, "Will dump a cpu and mem profiling file")

	servicesToSyncFlags = make(map[string]*bool)
	for _, service := range awsservices.ServiceNames {
		servicesToSyncFlags[service] = new(bool)
		syncCmd.Flags().BoolVar(servicesToSyncFlags[service], service, false, fmt.Sprintf("Sync '%s' service only", service))
	}
}

var syncCmd = &cobra.Command{
	Use:               "sync",
	Short:             "Manual sync of remote resources to the local store (ex: when autosync is unset)",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, firstInstallDoneHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),

	RunE: func(cmd *cobra.Command, args []string) error {
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
		localGraphs := make(map[string]cloud.GraphAPI)
		for _, service := range services {
			localGraphs[service.Name()] = sync.LoadLocalGraphForService(service.Name(), config.GetAWSProfile(), config.GetAWSRegion())
		}
		logger.Infof("running sync for region '%s'", config.GetAWSRegion())

		var syncErr error
		var graphs map[string]cloud.GraphAPI
		syncFn := func() {
			graphs, syncErr = sync.DefaultSyncer.Sync(services...)
		}

		start := time.Now()
		if profileSyncFlag {
			withProfiling(syncFn)
		} else {
			syncFn()
		}
		if syncErr != nil {
			logger.Verbose(syncErr)
		}

		for k, g := range graphs {
			displaySyncStats(k, g)
		}
		logger.Infof("sync took %s", time.Since(start))

		return nil
	},
}

func withProfiling(fn func()) {
	logger.Infof("sync profiling on")
	mem, err := os.Create("mem-sync.prof")
	if err != nil {
		log.Fatal("could not create mem profile: ", err)
	}
	logger.Infof("running garbage collection before profiling")
	runtime.GC() // cleaned up memeory before running function
	defer mem.Close()

	cpu, err := os.Create("cpu-sync.prof")
	if err != nil {
		log.Fatal("could not create cpu profile: ", err)
	}
	if err := pprof.StartCPUProfile(cpu); err != nil {
		log.Fatal("could not start cpu profile: ", err)
	}

	fn()

	pprof.StopCPUProfile()
	if err := pprof.WriteHeapProfile(mem); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	logger.Infof("Generated profiling files %s and %s", cpu.Name(), mem.Name())
}

func displaySyncStats(serviceName string, g cloud.GraphAPI) {
	var strs []string
	for rt, service := range awsservices.ServicePerResourceType {
		if service == serviceName {
			res, err := g.Find(cloud.NewQuery(rt))
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
