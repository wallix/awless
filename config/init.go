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

	"strconv"

	"github.com/wallix/awless/aws/services"
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
	os.Setenv("__AWLESS_CACHE", filepath.Join(AwlessHome, "cache"))
	os.Setenv("__AWLESS_KEYS_DIR", KeysDir)
}

func InitAwlessEnv() error {
	_, err := os.Stat(DBPath)

	AwlessFirstInstall = os.IsNotExist(err)
	os.Setenv("__AWLESS_FIRST_INSTALL", strconv.FormatBool(AwlessFirstInstall))

	os.MkdirAll(KeysDir, 0700)

	if AwlessFirstInstall {
		fmt.Fprintln(os.Stderr, AWLESS_ASCII_LOGO)
		fmt.Fprintln(os.Stderr, "Welcome! Resolving environment data...")
		fmt.Fprintln(os.Stderr)

		if err = InitConfig(resolveRequiredConfigFromEnv()); err != nil {
			return err
		}

		err = database.Execute(func(db *database.DB) error {
			return db.SetStringValue("current.version", Version)
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot store current version in db: %s\n", err)
		}
	}

	if err = LoadConfig(); err != nil {
		return err
	}

	return nil
}

func resolveRequiredConfigFromEnv() map[string]string {
	region := awsservices.ResolveRegionFromEnv()

	resolved := make(map[string]string)
	if region != "" {
		resolved[RegionConfigKey] = region
	}

	return resolved
}
