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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/wallix/awless/database"
)

const (
	lastUpgradeCheckDbKey = "upgrade.lastcheck"
)

func VerifyNewVersionAvailable(url string, messaging io.Writer) error {
	return database.Execute(func(db *database.DB) error {
		last, err := db.GetTimeValue(lastUpgradeCheckDbKey)
		if err != nil {
			return err
		}

		upgradeFreq := getCheckUpgradeFrequency()
		if upgradeFreq < 0 {
			return nil
		}

		if time.Since(last) > upgradeFreq {
			notifyIfUpgrade(url, messaging)
		}

		return db.SetTimeValue(lastUpgradeCheckDbKey, time.Now())
	})
}

func notifyIfUpgrade(url string, messaging io.Writer) error {
	client := &http.Client{Timeout: 1500 * time.Millisecond}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "awless-client-"+Version)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	latest := struct {
		Version, URL string
	}{}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&latest); err == nil {
		if IsSemverUpgrade(Version, latest.Version) {
			var install string
			switch BuildFor {
			case "brew":
				install = "Run `brew upgrade awless`"
			case "zip", "targz":
				ext := "tar.gz"
				if runtime.GOOS == "windows" {
					ext = "zip"
				}
				install = fmt.Sprintf("Run `wget -O awless-%s.%s https://github.com/wallix/awless/releases/download/%s/awless-%s-%s.%s`", latest.Version, ext, latest.Version, runtime.GOOS, runtime.GOARCH, ext)
			default:
				install = "Run `go get -u github.com/wallix/awless`"
			}
			fmt.Fprintf(messaging, "New version %s available. Checkout the latest features at https://github.com/wallix/awless/blob/master/CHANGELOG.md\n%s\n", latest.Version, install)
		}
	}

	return nil
}

const semverLen = 3

type semver [semverLen]int

var SemverInvalidFormatErr = errors.New("semver invalid format")

func IsSemverUpgrade(current, latest string) bool {
	i, err := CompareSemver(current, latest)
	if err != nil {
		return false
	}

	return i < 0
}

func CompareSemver(current, latest string) (int, error) {
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	dot := func(r rune) bool {
		return r == '.'
	}
	cFields := strings.FieldsFunc(current, dot)
	lFields := strings.FieldsFunc(latest, dot)

	if len(cFields) != semverLen || len(lFields) != semverLen {
		return 0, SemverInvalidFormatErr
	}

	currents := new(semver)
	for i, f := range cFields {
		num, err := strconv.Atoi(f)
		if err != nil {
			return 0, SemverInvalidFormatErr
		}
		currents[i] = num
	}

	latests := new(semver)
	for i, f := range lFields {
		num, err := strconv.Atoi(f)
		if err != nil {
			return 0, SemverInvalidFormatErr
		}
		latests[i] = num
	}

	for i := 0; i < semverLen; i++ {
		if latests[i] > currents[i] {
			return -1, nil
		} else if latests[i] == currents[i] {
			continue
		} else {
			return 1, nil
		}
	}

	return 0, nil
}
