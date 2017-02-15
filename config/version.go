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
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var (
	Version                                 = "0.0.12"
	buildSha, buildDate, buildArch, buildOS string
)

type BuildInfo struct {
	Version, Sha, Date, Arch, OS string
}

func (b BuildInfo) String() string {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("version=%s", b.Version))

	if b.Sha != "" {
		buff.WriteString(fmt.Sprintf(", commit=%s", b.Sha))
	}
	if b.Date != "" {
		buff.WriteString(fmt.Sprintf(", build-date=%s", b.Date))
	}
	if b.Arch != "" {
		buff.WriteString(fmt.Sprintf(", build-arch=%s", b.Arch))
	}
	if b.OS != "" {
		buff.WriteString(fmt.Sprintf(", build-os=%s", b.OS))
	}
	return buff.String()
}

var CurrentBuildInfo = BuildInfo{
	Version: Version,
	Sha:     buildSha,
	Date:    buildDate,
	Arch:    buildArch,
	OS:      buildOS,
}

const semverLen = 3

type semver [semverLen]int

func IsUpgrade(current, latest string) bool {
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
