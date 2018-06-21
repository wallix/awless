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
)

const AWLESS_ASCII_LOGO = `
 █████╗  ██╗    ██╗ ██╗     ██████  ██████╗ ██████╗     
██╔══██╗ ██║    ██║ ██║     ██╔══╝  ██╔═══╝ ██╔═══╝    
███████║ ██║ █╗ ██║ ██║     ████╗   ██████  ██████   
██╔══██║ ██║███╗██║ ██║     ██╔═╝       ██╗     ██╗   
██║  ██║ ╚███╔███╔╝ ██████╗ ██████╗ ██████║ ██████║   
╚═╝  ╚═╝  ╚══╝╚══╝  ╚═════╝ ╚═════╝ ╚═════╝ ╚═════╝
`

var (
	Version  = "v0.1.11"
	BuildFor string

	buildSha, buildDate, buildArch, buildOS string
)

type BuildInfo struct {
	Version, Sha, Date, Arch, OS, For string
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
	if b.For != "" {
		buff.WriteString(fmt.Sprintf(", build-for=%s", b.For))
	}
	return buff.String()
}

var CurrentBuildInfo = BuildInfo{
	Version: Version,
	For:     BuildFor,
	Sha:     buildSha,
	Date:    buildDate,
	Arch:    buildArch,
	OS:      buildOS,
}
