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
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws/services"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
)

func applyHooks(funcs ...func(*cobra.Command, []string) error) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		for _, fn := range funcs {
			if err := fn(cmd, args); err != nil {
				exitOn(err)
			}
		}
	}
}

func initAwlessEnvHook(cmd *cobra.Command, args []string) error {
	if err := config.InitAwlessEnv(); err != nil {
		return fmt.Errorf("cannot init awless environment: %s", err)
	}
	if awsRegionGlobalFlag != "" {
		if err := config.SetVolatile(config.RegionConfigKey, awsRegionGlobalFlag); err != nil {
			return err
		}
	} else if envRegion := os.Getenv("AWS_DEFAULT_REGION"); envRegion != "" {
		if err := config.SetVolatile(config.RegionConfigKey, envRegion); err != nil {
			return err
		}
	}
	if awsProfileGlobalFlag != "" {
		if err := config.SetVolatile(config.ProfileConfigKey, awsProfileGlobalFlag); err != nil {
			return err
		}
	} else if envProfile := os.Getenv("AWS_DEFAULT_PROFILE"); envProfile != "" {
		if err := config.SetVolatile(config.ProfileConfigKey, envProfile); err != nil {
			return err
		}
	}

	switch awsColorGlobalFlag {
	case "never":
		color.NoColor = true
	case "always":
		color.NoColor = false
	}

	return nil
}

func initCloudServicesHook(cmd *cobra.Command, args []string) error {
	if localGlobalFlag {
		return nil
	}
	awsConf := config.GetConfigWithPrefix("aws.")
	logger.Verbosef("awless %s - loading AWS session with profile '%v' and region '%v'", config.Version, awsConf[config.ProfileConfigKey], awsConf[config.RegionConfigKey])

	if err := awsservices.Init(awsConf, logger.DefaultLogger, config.SetProfileCallback, networkMonitorFlag); err != nil {
		return err
	}

	if config.TriggerSyncOnConfigUpdate && !strings.HasPrefix(cmd.Name(), "sync") {
		var services []cloud.Service
		for _, s := range cloud.ServiceRegistry {
			services = append(services, s)
		}
		if !noSyncGlobalFlag {
			logger.Infof("Syncing new region '%s'... (disable with --no-sync global flag)", awsConf[config.RegionConfigKey])
			sync.NewSyncer(logger.DefaultLogger).Sync(services...)
		}
	}

	return nil
}

func includeHookIf(cond *bool, hook func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		if *cond {
			return hook(c, args)
		}
		return nil
	}
}

func initSyncerHook(cmd *cobra.Command, args []string) error {
	if noSyncGlobalFlag {
		sync.DefaultSyncer = sync.NoOpSyncer()
	} else {
		sync.DefaultSyncer = sync.NewSyncer(logger.DefaultLogger)
	}
	return nil
}

func initLoggerHook(cmd *cobra.Command, args []string) error {
	var flag int
	if verboseGlobalFlag {
		flag = logger.VerboseF
	}
	if extraVerboseGlobalFlag {
		flag = flag | logger.ExtraVerboseF
	}

	logger.DefaultLogger.SetVerbose(flag)
	if silentGlobalFlag {
		logger.DefaultLogger = logger.DiscardLogger
	}
	return nil
}

func onVersionUpgrade(cmd *cobra.Command, args []string) error {
	var lastVersion string
	if derr := database.Execute(func(db *database.DB) (err error) {
		lastVersion, err = db.GetStringValue("current.version")
		return
	}); derr != nil {
		fmt.Printf("cannot verify stored version in db: %s\n", derr)
	}

	if config.IsSemverUpgrade(lastVersion, config.Version) {
		if err := database.Execute(func(db *database.DB) error {
			return db.SetStringValue("current.version", config.Version)
		}); err != nil {
			fmt.Printf("cannot store upgraded version in db: %s\n", err)
		}
		logger.Infof("You have just upgraded awless from %s to %s", lastVersion, config.Version)
		migrationActionsAndExtraMessages(config.Version)
	}

	return nil
}

func verifyNewVersionHook(cmd *cobra.Command, args []string) error {
	if localGlobalFlag {
		return nil
	}
	config.VerifyNewVersionAvailable("https://updates.awless.io", os.Stderr)
	return nil
}

func networkMonitorHook(cmd *cobra.Command, args []string) error {
	if networkMonitorFlag {
		awsservices.DefaultNetworkMonitor.DisplayStats(os.Stderr)
	}
	return nil
}

func firstInstallDoneHook(cmd *cobra.Command, args []string) error {
	if config.TriggerSyncOnConfigUpdate {
		fmt.Fprintln(os.Stderr, "\nAll done. Enjoy!")
		fmt.Fprintln(os.Stderr, "You can review and configure awless with `awless config`")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Now running: `%s`\n", cmd.CommandPath())
	}
	return nil
}

func migrationActionsAndExtraMessages(current string) {
	switch current {
	case "v0.1.7":
		config.Set("instance.distro", "amazonlinux")
		logger.Info("In v0.1.7, the default template config value 'instance.image' has been deprecated in favor of 'instance.distro'")

		ami, _ := config.Get("instance.image")
		if isNotAwlessFormerDefaultAMI(fmt.Sprint(ami)) {
			logger.Warningf("\tYou had a customized value of '%s' for the now deprecated 'instance.image'", fmt.Sprint(ami))
			logger.Warning("\tThis value will not be taken into account anymore as default when running templates")
		} else {
			logger.Info("\tMigrated correctly the deprecated 'instance.image' to 'instance.distro'")
		}
		config.Unset("instance.image")
		logger.Info("\tYou can check your config values with 'awless config'")
	}
}

func isNotAwlessFormerDefaultAMI(s string) bool {
	amis := []string{"ami-c58c1dd3", "ami-4191b524", "ami-7a85a01a", "ami-4836a428", "ami-0bd66a6f", "ami-d3c0c4b5", "ami-b6daced2", "ami-b968bad6", "ami-fc5ae39f", "ami-762a2315", "ami-923d12f5", "ami-9d15c7f3", "ami-52c7b43d", "ami-2bccae47"}
	for _, e := range amis {
		if e == s {
			return false
		}
	}
	return true
}
