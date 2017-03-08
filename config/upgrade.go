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
	db, err, close := database.Current()
	if err != nil {
		return err
	}
	defer close()

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
}

func notifyIfUpgrade(url string, messaging io.Writer) error {
	client := &http.Client{Timeout: 3 * time.Second}
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
		if isSemverUpgrade(Version, latest.Version) {
			var install string
			switch BuildFor {
			case "brew":
				install = "Run `brew upgrade awless`"
			default:
				install = fmt.Sprintf("Run `wget -O awless-%s.zip https://github.com/wallix/awless/releases/download/%s/awless-%s-%s.zip`", latest.Version, latest.Version, runtime.GOOS, runtime.GOARCH)
			}
			fmt.Fprintf(messaging, "New version %s available. Changelog at https://github.com/wallix/awless/blob/master/CHANGELOG.md\n%s\n", latest.Version, install)
		}
	}

	return nil
}

const semverLen = 3

type semver [semverLen]int

func isSemverUpgrade(current, latest string) bool {
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	dot := func(r rune) bool {
		return r == '.'
	}
	cFields := strings.FieldsFunc(current, dot)
	lFields := strings.FieldsFunc(latest, dot)

	if len(cFields) != semverLen || len(lFields) != semverLen {
		return false
	}

	currents := new(semver)
	for i, f := range cFields {
		num, err := strconv.Atoi(f)
		if err != nil {
			return false
		}
		currents[i] = num
	}

	latests := new(semver)
	for i, f := range lFields {
		num, err := strconv.Atoi(f)
		if err != nil {
			return false
		}
		latests[i] = num
	}

	for i := 0; i < semverLen; i++ {
		if latests[i] > currents[i] {
			return true
		} else if latests[i] == currents[i] {
			continue
		} else {
			return false
		}
	}

	return false
}
