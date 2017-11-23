package config

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/database"
)

var (
	Config   = map[string]interface{}{}
	Defaults = map[string]interface{}{}
)

const (
	configDatabaseKey   = "userconfig"
	defaultsDatabaseKey = "defaults"

	//Config
	autosyncConfigKey              = "autosync"
	checkUpgradeFrequencyConfigKey = "upgrade.checkfrequency"
	schedulerURL                   = "scheduler.url"
	RegionConfigKey                = "aws.region"
	ProfileConfigKey               = "aws.profile"

	//Config prefix
	awsCloudPrefix = "aws."
)

var configDefinitions = map[string]*Definition{
	autosyncConfigKey:              {help: "Automatically synchronize your cloud locally", defaultValue: "true", parseParamFn: parseBool},
	RegionConfigKey:                {help: "AWS region", parseParamFn: awsconfig.ParseRegion, stdinParamProviderFn: awsconfig.StdinRegionSelector, onUpdateFns: []onUpdateFunc{runSyncWithUpdatedRegion}},
	ProfileConfigKey:               {help: "AWS profile", defaultValue: "default"},
	"aws.infra.sync":               {help: "Enable/disable sync of infra services (EC2, RDS, etc.) (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.access.sync":              {help: "Enable/disable sync of IAM service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.storage.sync":             {help: "Enable/disable sync of S3 service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.storage.s3object.sync":    {help: "Enable/disable sync of S3/s3object (when empty: true)", defaultValue: "false", parseParamFn: parseBool},
	"aws.dns.sync":                 {help: "Enable/disable sync of DNS service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.dns.record.sync":          {help: "Enable/disable sync of DNS/record (when empty: true)", defaultValue: "false", parseParamFn: parseBool},
	"aws.notification.sync":        {help: "Enable/disable sync of SNS service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.monitoring.sync":          {help: "Enable/disable sync of CloudWatch service (when empty: true)", defaultValue: "false", parseParamFn: parseBool},
	"aws.lambda.sync":              {help: "Enable/disable sync of Lambda service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.messaging.sync":           {help: "Enable/disable sync of SQS/SNS service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.cdn.sync":                 {help: "Enable/disable sync of CloudFront service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.cloudformation.sync":      {help: "Enable/disable sync of CloudFormation service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	checkUpgradeFrequencyConfigKey: {help: "Upgrade check frequency (hours); a negative value disables check", defaultValue: "8", parseParamFn: parseInt},
	schedulerURL:                   {help: "URL used by awless CLI to interact with pre-installed https://github.com/wallix/awless-scheduler", defaultValue: "http://localhost:8082"},
}

var defaultsDefinitions = map[string]*Definition{
	"instance.type":          {defaultValue: "t2.micro", help: "AWS EC2 instance type", stdinParamProviderFn: awsconfig.StdinInstanceTypeSelector, parseParamFn: awsconfig.ParseInstanceType},
	"instance.distro":        {defaultValue: "amazonlinux", help: "Query to fetch latest community bare distro image id (see awless search images -h)", parseParamFn: parseDistroQuery},
	"instance.count":         {defaultValue: "1", help: "Number of instances to create on AWS EC2", parseParamFn: parseInt},
	"instance.timeout":       {defaultValue: "180", help: "Time to wait when checking instance states on AWS EC2", parseParamFn: parseInt},
	"securitygroup.protocol": {defaultValue: "tcp", help: "The IP protocol to authorize on the security group"},
	"volume.device":          {defaultValue: "/dev/sdh", help: "Device name to expose to an EC2 instance"},
	"elasticip.domain":       {defaultValue: "vpc", help: "The domain of elastic IP addresses (standard or vpc)"},
	"image.delete-snapshots": {defaultValue: "true", help: "Delete linked snapshots when deleting an image"},
	"database.type":          {defaultValue: "db.t2.micro", help: "Default RDS database type"},
}

var deprecated = map[string]string{
	"sync.auto": autosyncConfigKey,
	"region":    RegionConfigKey,
}

var TriggerSyncOnConfigUpdate bool

type onUpdateFunc func(interface{})

type Definition struct {
	help                 string
	parseParamFn         func(string) (interface{}, error)
	stdinParamProviderFn func() string
	onUpdateFns          []onUpdateFunc
	defaultValue         string
}

func LoadConfig() error {
	err := database.Execute(func(db *database.DB) (dberr error) {
		Config, dberr = db.GetConfigs(configDatabaseKey)
		if dberr != nil {
			return fmt.Errorf("config: load config: %s", dberr)
		}

		Defaults, dberr = db.GetConfigs(defaultsDatabaseKey)
		if dberr != nil {
			return fmt.Errorf("config: load defaults: %s", dberr)
		}
		return
	})

	return err
}

func DisplayConfig() string {
	return fmt.Sprintf("%s\n%s", displayConfig(), displayDefaults())
}

func InitConfig(fromEnv map[string]string) error {
	for k, v := range configDefinitions {
		val := v.defaultValue
		if vv, ok := fromEnv[k]; ok {
			val = vv
		}
		if err := Set(k, val); err != nil {
			return err
		}
	}
	for k, v := range defaultsDefinitions {
		val := v.defaultValue
		if vv, ok := fromEnv[k]; ok {
			val = vv
		}
		if err := Set(k, val); err != nil {
			return err
		}
	}
	return nil
}

func Set(key, value string) error {
	v, def, isConf, err := setVolatile(key, value)
	if err != nil {
		return err
	}
	var databaseKey string
	if isConf {
		databaseKey = configDatabaseKey
	} else {
		databaseKey = defaultsDatabaseKey
	}

	if err := database.Execute(func(db *database.DB) error {
		return db.SetConfig(databaseKey, key, v)
	}); err != nil {
		return err
	}

	if def != nil {
		for _, fn := range def.onUpdateFns {
			fn(v)
		}
	}

	return nil
}

func SetProfileCallback(value string) error {
	return Set(ProfileConfigKey, value)
}

func Unset(key string) error {
	var dbKey string
	if _, ok := Config[key]; ok {
		delete(Config, key)
		dbKey = configDatabaseKey
	}
	if _, ok := Defaults[key]; ok {
		delete(Defaults, key)
		dbKey = defaultsDatabaseKey
	}
	if dbKey != "" {
		if err := database.Execute(func(db *database.DB) error {
			return db.UnsetConfig(dbKey, key)
		}); err != nil {
			return fmt.Errorf("unset config: %s", err)
		}
	}

	return nil
}

func Get(key string) (interface{}, bool) {
	if v, ok := Config[key]; ok {
		return v, ok
	}
	v, ok := Defaults[key]
	return v, ok
}

func SetVolatile(key, value string) error {
	_, _, _, err := setVolatile(key, value)
	return err
}

func InteractiveSet(key string) error {
	var val string
	if def, ok := configDefinitions[key]; ok && def.stdinParamProviderFn != nil {
		val = def.stdinParamProviderFn()
	} else if def, ok := defaultsDefinitions[key]; ok && def.stdinParamProviderFn != nil {
		val = def.stdinParamProviderFn()
	} else {
		val = defaultStdinParamProvider()
	}
	return Set(key, val)
}

func parseBool(i string) (interface{}, error) {
	b, err := strconv.ParseBool(i)
	if err != nil {
		return b, fmt.Errorf("invalid value, expected a boolean, got '%s'", i)
	}
	return b, nil
}

func parseInt(a string) (interface{}, error) {
	i, err := strconv.Atoi(a)
	if err != nil {
		return i, fmt.Errorf("invalid value, expected an int, got '%s'", a)
	}
	return i, nil
}

func defaultParser(value string) (interface{}, error) {
	if num, err := strconv.Atoi(value); err == nil {
		return num, nil
	}
	if b, err := strconv.ParseBool(value); err == nil {
		return b, nil
	}
	return value, nil
}

func parseDistroQuery(v string) (interface{}, error) {
	_, err := awsspec.ParseImageQuery(v)
	return v, err
}

func defaultStdinParamProvider() string {
	var value string
	for value == "" {
		fmt.Print("Value ? > ")
		fmt.Scan(&value)
	}
	return value
}

func setVolatile(key, value string) (interface{}, *Definition, bool, error) {
	var isConf bool
	confDef, confOk := configDefinitions[key]
	defDef, defOk := defaultsDefinitions[key]
	var def *Definition
	switch {
	case confOk && defOk:
		return nil, def, isConf, fmt.Errorf("%s can not be in both config and defaults", key)
	case confOk:
		isConf = true
		def = confDef
	case defOk:
		def = defDef
	default:
		if strings.Contains(key, awsCloudPrefix) {
			isConf = true
		}
	}
	var v interface{}
	var err error
	if def != nil && def.parseParamFn != nil {
		if v, err = def.parseParamFn(value); err != nil {
			return nil, def, isConf, err
		}
	} else {
		if v, err = defaultParser(value); err != nil {
			return nil, def, isConf, err
		}
	}
	if isConf {
		Config[key] = v
	} else {
		Defaults[key] = v
	}
	return v, def, isConf, nil
}

func displayConfig() string {
	var b bytes.Buffer
	b.WriteString("# Config parameters\n")
	t := tabwriter.NewWriter(&b, 0, 0, 3, ' ', 0)
	var keys []string
	for k := range Config {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(t, "\t%s:\t%v\t(%[2]T)", k, Config[k])
		if def, ok := configDefinitions[k]; ok && def.help != "" {
			fmt.Fprintf(t, "\t# %s\n", def.help)
		} else {
			fmt.Fprintln(t)
		}
	}
	for k := range configDefinitions {
		if _, ok := Config[k]; !ok {
			fmt.Fprintf(t, "\t%s:\t\t", k)
			if def, ok := configDefinitions[k]; ok && def.help != "" {
				fmt.Fprintf(t, "\t# %s\n", def.help)
			} else {
				fmt.Fprintln(t)
			}
		}
	}
	t.Flush()
	return b.String()
}

func displayDefaults() string {
	var b bytes.Buffer
	b.WriteString("# Template defaults\n")
	b.WriteString("   ## Predefined\n")
	t := tabwriter.NewWriter(&b, 0, 0, 3, ' ', 0)
	var keys []string
	for k := range Defaults {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if def, ok := defaultsDefinitions[k]; ok {
			if def.help != "" {
				fmt.Fprintf(t, "\t%s:\t%v\t(%[2]T)\t# %s\n", k, Defaults[k], def.help)
			} else {
				fmt.Fprintf(t, "\t%s:\t%v\t(%[2]T)\n", k, Defaults[k])
			}
		}
	}
	for k := range defaultsDefinitions {
		if _, ok := Defaults[k]; !ok {
			fmt.Fprintf(t, "\t%s:\t \t(UNSET)", k)
			if def, ok := defaultsDefinitions[k]; ok && def.help != "" {
				fmt.Fprintf(t, "\t# %s\n", def.help)
			} else {
				fmt.Fprintln(t)
			}
		}
	}
	t.Flush()
	count := 0
	t = tabwriter.NewWriter(&b, 0, 0, 3, ' ', 0)
	for _, k := range keys {
		if _, ok := defaultsDefinitions[k]; !ok {
			count++
			fmt.Fprintf(t, "\t%s:\t%v\t(%[2]T)", k, Defaults[k])
			if newKey, ok := deprecated[k]; ok {
				fmt.Fprintf(t, "\t# DEPRECATED, update with `awless config set %s` `awless config unset %s`", newKey, k)
			}
			fmt.Fprintln(t)
		}
	}
	if count > 0 {
		b.WriteString("\n   ## User defined\n")
		t.Flush()
	}
	return b.String()
}

func runSyncWithUpdatedRegion(i interface{}) {
	if !GetAutosync() {
		return
	}

	if !awsconfig.IsValidRegion(fmt.Sprint(i)) {
		return
	}

	TriggerSyncOnConfigUpdate = true
}
