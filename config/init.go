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
	"github.com/aws/aws-sdk-go/aws/session"
	awsconfig "github.com/wallix/awless/aws/config"
)

var (
	AwlessHome                          = filepath.Join(os.Getenv("HOME"), ".awless")
	RepoDir                             = filepath.Join(AwlessHome, "aws", "rdf")
	Dir                                 = filepath.Join(AwlessHome, "aws")
	KeysDir                             = filepath.Join(AwlessHome, "keys")
	InfraFilename                       = "infra.rdf"
	AccessFilename                      = "access.rdf"
	AwlessFirstInstall, AwlessFirstSync bool
)

func InitAwlessEnv() error {
	os.Setenv("__AWLESS_HOME", AwlessHome)
	os.Setenv("__AWLESS_KEYS_DIR", KeysDir)
	_, err := os.Stat(AwlessHome)
	_, ierr := os.Stat(filepath.Join(RepoDir, InfraFilename))
	_, aerr := os.Stat(filepath.Join(RepoDir, AccessFilename))
	AwlessFirstSync = os.IsNotExist(ierr) || os.IsNotExist(aerr)

	AwlessFirstInstall = os.IsNotExist(err)

	os.MkdirAll(RepoDir, 0700)
	os.MkdirAll(KeysDir, 0700)

	if AwlessFirstInstall {
		fmt.Println("First install. Welcome!")
		fmt.Println()
		if err = InitConfigAndDefaults(); err != nil {
			return err
		}

		if _, err = overwriteDefaults(); err != nil {
			return err
		}
	}

	if err = LoadAll(); err != nil {
		return err
	}

	return nil
}

func overwriteDefaults() (string, error) {
	var region, ami string

	if sess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}); err == nil {
		region = awssdk.StringValue(sess.Config.Region)
	}

	if awsconfig.IsValidRegion(region) {
		fmt.Printf("Found existing AWS region '%s'. Setting it as your default region.\n", region)
	} else {
		fmt.Println("Could not find any AWS region in your environment. Please choose one region:")
		region = awsconfig.StdinRegionSelector()
	}

	var hasAMI bool
	if ami, hasAMI = awsconfig.AmiPerRegion[region]; !hasAMI {
		fmt.Printf("Could not find a default ami for your region %s\n. Set it manually with `awless config set instance.image ...`", region)
	}

	if err := Set(RegionConfigKey, region); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	if hasAMI {
		if err := Set(instanceImageDefaultsKey, ami); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}

	fmt.Println("\nThose parameters have been set in your config:")
	fmt.Println(Display())

	fmt.Println("\nShow and update config with `awless config`. Ex: `awless config set aws.region`")
	fmt.Println("\nAll done. Enjoy!")
	fmt.Println()

	return region, nil
}
