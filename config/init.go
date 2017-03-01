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
	"strings"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/database"
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
		fmt.Println("First install. Welcome!\n")
		_, err = resolveAndSetDefaults()
		if err != nil {
			return err
		}
	}

	return nil
}

func resolveAndSetDefaults() (string, error) {
	var region, ami string

	if sess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}); err == nil {
		region = awssdk.StringValue(sess.Config.Region)
	}

	if aws.IsValidRegion(region) {
		fmt.Printf("Found existing AWS region '%s'. Setting it as your default region.\n", region)
	} else {
		fmt.Println("Could not find any AWS region in your environment. Please choose one region:")
		region = askRegion()
	}

	var hasAMI bool
	if ami, hasAMI = amiPerRegion[region]; !hasAMI {
		fmt.Printf("Could not find a default ami for your region %s\n. Set it manually with `awless config set instance.image ...`", region)
	}

	defaults := map[string]interface{}{
		database.SyncAuto:         true,
		database.RegionKey:        region,
		database.InstanceTypeKey:  "t2.micro",
		database.InstanceCountKey: 1,
		database.ProfileKey:       "default",
	}

	if hasAMI {
		defaults[database.InstanceImageKey] = ami
	}

	db, err, close := database.Current()
	if err != nil {
		return region, fmt.Errorf("database error: %s", err)
	}
	defer close()
	for k, v := range defaults {
		err := db.SetDefault(k, v)
		if err != nil {
			return region, err
		}
	}

	fmt.Println("\nThose defaults have been set in your config:")
	for k, v := range defaults {
		fmt.Printf("\t%s = %v\n", k, v)
	}
	fmt.Println("\nShow and update config with `awless config`. Ex: `awless config set region`")
	fmt.Println("\nAll done. Enjoy!\n")

	return region, nil
}

func askRegion() string {
	var region string

	fmt.Println(strings.Join(aws.AllRegions(), ", "))
	fmt.Println()
	fmt.Print("Copy/paste one region > ")
	fmt.Scan(&region)
	for !aws.IsValidRegion(region) {
		fmt.Print("Invalid! Copy/paste one valid region > ")
		fmt.Scan(&region)
	}
	return region
}

var amiPerRegion = map[string]string{
	"us-east-1":      "ami-0b33d91d",
	"us-east-2":      "ami-c55673a0",
	"us-west-1":      "ami-165a0876",
	"us-west-2":      "ami-f173cc91",
	"ca-central-1":   "ami-ebed508f",
	"eu-west-1":      "ami-70edb016",
	"eu-west-2":      "ami-f1949e95",
	"eu-central-1":   "ami-af0fc0c0",
	"ap-southeast-1": "ami-dc9339bf",
	"ap-southeast-2": "ami-1c47407f",
	"ap-northeast-1": "ami-56d4ad31",
	"ap-northeast-2": "ami-dac312b4",
	"ap-south-1":     "ami-f9daac96",
	"sa-east-1":      "ami-80086dec",
}
