package config

import (
	"fmt"
	"strings"
	"time"
)

func GetAWSRegion() string {
	if reg, ok := Config[RegionConfigKey]; ok && reg != "" {
		return fmt.Sprint(reg)
	}
	if reg, ok := Defaults["region"]; ok && reg != "" { // Compatibility with old key
		return fmt.Sprint(reg)
	}
	return ""
}

const defaultAWSSessionProfile = "default"

func GetAWSProfile() string {
	if profile, ok := Config[ProfileConfigKey]; ok && profile != "" {
		return fmt.Sprint(profile)
	}
	if profile, ok := Defaults[ProfileConfigKey]; ok && profile != "" { // Compatibility with old key
		return fmt.Sprint(profile)
	}
	return defaultAWSSessionProfile
}

func GetAutosync() bool {
	if autoSync, ok := Config[autosyncConfigKey].(bool); ok {
		return autoSync
	}
	if autoSync, ok := Defaults["sync.auto"].(bool); ok { //Compatibility with old key
		return autoSync
	}
	return true
}

func GetSchedulerURL() string {
	if u, ok := Config[schedulerURL].(string); ok {
		return u
	}
	return ""
}

func GetConfigWithPrefix(prefix string) map[string]interface{} {
	conf := make(map[string]interface{})
	for k, v := range Config {
		if strings.HasPrefix(k, prefix) {
			conf[k] = v
		}
	}
	return conf
}

func getCheckUpgradeFrequency() time.Duration {
	if frequency, ok := Config[checkUpgradeFrequencyConfigKey].(int); ok {
		return time.Duration(frequency) * time.Hour
	}
	return 8 * time.Hour
}
