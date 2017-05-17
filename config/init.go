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

package config

import (
	"fmt"
	"os"
	"path/filepath"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/database"
)

var (
	AwlessHome         = filepath.Join(os.Getenv("HOME"), ".awless")
	DBPath             = filepath.Join(AwlessHome, database.Filename)
	Dir                = filepath.Join(AwlessHome, "aws")
	KeysDir            = filepath.Join(AwlessHome, "keys")
	AwlessFirstInstall bool
)

func init() {
	os.Setenv("__AWLESS_HOME", AwlessHome)
	os.Setenv("__AWLESS_KEYS_DIR", KeysDir)
}

func InitAwlessEnv() error {
	_, err := os.Stat(DBPath)

	AwlessFirstInstall = os.IsNotExist(err)

	os.MkdirAll(KeysDir, 0700)

	if AwlessFirstInstall {
		fmt.Println("First install. Welcome!")
		fmt.Println()

		resolved, err := resolveRequiredConfigFromEnv()
		if err != nil {
			return err
		}

		if err := InitConfig(resolved); err != nil {
			return err
		}

		fmt.Println("\nThose parameters have been set in your config:")
		fmt.Println(DisplayConfig())

		fmt.Println("\nShow and update config with `awless config`. Ex: `awless config set aws.region`")
		fmt.Println("\nAll done. Enjoy!")
		fmt.Println()

		err = database.Execute(func(db *database.DB) error {
			return db.SetStringValue("current.version", Version)
		})
		if err != nil {
			fmt.Printf("cannot store current version in db: %s\n", err)
		}
	}

	if err = LoadConfig(); err != nil {
		return err
	}

	return nil
}

func resolveRequiredConfigFromEnv() (map[string]string, error) {
	var region, ami string
	var sess *session.Session
	var err error

	if sess, err = session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}); err == nil {
		region = awssdk.StringValue(sess.Config.Region)
	}

	if awsconfig.IsValidRegion(region) {
		fmt.Printf("Found existing AWS region '%s'. Setting it as your default region.\n", region)
	} else if sess != nil {
		if r, err := ec2metadata.New(sess).Region(); err == nil {
			fmt.Printf("Found AWS region '%s' from local EC2 instance metadata. Setting it as your default region.\n", region)
			region = r
		}
	}

	if region == "" {
		fmt.Println("Could not find any AWS region in your environment. Please choose one region:")
		region = awsconfig.StdinRegionSelector()
	}

	var hasAMI bool
	if ami, hasAMI = awsconfig.AmiPerRegion[region]; !hasAMI {
		fmt.Printf("Could not find a default ami for your region %s\n. Set it later manually with `awless config set instance.image ...`", region)
	}

	resolved := make(map[string]string)
	resolved[RegionConfigKey] = region
	if hasAMI {
		resolved[instanceImageDefaultsKey] = ami
	}

	return resolved, nil
}
