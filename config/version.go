package config

import (
	"bytes"
	"fmt"
)

var (
	Version                                 = "0.0.3"
	buildSha, buildDate, buildArch, buildOS string
)

type BuildInfo struct {
	Version   string
	BuildSha  string
	BuildDate string
	BuildArch string
	BuildOS   string
}

func (b BuildInfo) String() string {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("version=%s", b.Version))

	if b.BuildSha != "" {
		buff.WriteString(fmt.Sprintf(", commit=%s", b.BuildSha))
	}
	if b.BuildDate != "" {
		buff.WriteString(fmt.Sprintf(", build-date=%s", b.BuildDate))
	}
	if b.BuildArch != "" {
		buff.WriteString(fmt.Sprintf(", build-arch=%s", b.BuildArch))
	}
	if b.BuildOS != "" {
		buff.WriteString(fmt.Sprintf(", build-os=%s", b.BuildOS))
	}
	return buff.String()
}

var CurrentBuildInfo = BuildInfo{
	Version:   Version,
	BuildSha:  buildSha,
	BuildDate: buildDate,
	BuildArch: buildArch,
	BuildOS:   buildOS,
}
