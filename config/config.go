package config

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	awsconfig "github.com/wallix/awless/aws/config"
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
	RegionConfigKey                = "aws.region"
	ProfileConfigKey               = "aws.profile"

	//Config prefix
	awsCloudPrefix = "aws."

	//Defaults
	instanceTypeDefaultsKey    = "instance.type"
	instanceImageDefaultsKey   = "instance.image"
	instanceCountDefaultsKey   = "instance.count"
	instanceTimeoutDefaultsKey = "instance.timeout"
)

var configDefinitions = map[string]*Definition{
	autosyncConfigKey:                {help: "Automatically synchronize your cloud locally", defaultValue: "true", parseParamFn: parseBool},
	RegionConfigKey:                  {help: "AWS region", defaultValue: "us-east-1", parseParamFn: awsconfig.ParseRegion, stdinParamProviderFn: awsconfig.StdinRegionSelector, onUpdateFn: awsconfig.WarningChangeRegion},
	ProfileConfigKey:                 {help: "AWS profile", defaultValue: "default"},
	"aws.infra.sync":                 {help: "Sync AWS EC2/ELBv2 service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.access.sync":                {help: "Sync AWS IAM service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.storage.sync":               {help: "Sync AWS S3 service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.storage.storageobject.sync": {help: "Sync AWS S3/storageobject (when empty: true)", defaultValue: "false", parseParamFn: parseBool},
	"aws.notification.sync":          {help: "Sync AWS SNS service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.queue.sync":                 {help: "Sync AWS SQS service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	"aws.dns.sync":                   {help: "Sync Route53 service (when empty: true)", defaultValue: "true", parseParamFn: parseBool},
	checkUpgradeFrequencyConfigKey:   {help: "Upgrade check frequency (hours); a negative value disables check", defaultValue: "8", parseParamFn: parseInt},
}

var defaultsDefinitions = map[string]*Definition{
	instanceTypeDefaultsKey:    {defaultValue: "t2.micro", help: "AWS EC2 instance type", stdinParamProviderFn: awsconfig.StdinInstanceTypeSelector, parseParamFn: awsconfig.ParseInstanceType},
	instanceImageDefaultsKey:   {help: "AWS EC2 AMI"},
	instanceCountDefaultsKey:   {defaultValue: "1", help: "Number of instances to create on AWS EC2", parseParamFn: parseInt},
	instanceTimeoutDefaultsKey: {defaultValue: "180", help: "Time to wait when checking instance states on AWS EC2", parseParamFn: parseInt},
}

var deprecated = map[string]string{
	"sync.auto": autosyncConfigKey,
	"region":    RegionConfigKey,
}

type Definition struct {
	help                 string
	parseParamFn         func(string) (interface{}, error)
	stdinParamProviderFn func() string
	onUpdateFn           func(interface{})
	defaultValue         string
}

func LoadAll() error {
	db, err, dbclose := database.Current()
	if err != nil {
		return fmt.Errorf("load config: %s", err)
	}
	defer dbclose()

	Config, err = db.GetConfigs(configDatabaseKey)
	if err != nil {
		return fmt.Errorf("config: load config: %s", err)
	}

	Defaults, err = db.GetConfigs(defaultsDatabaseKey)
	if err != nil {
		return fmt.Errorf("config: load defaults: %s", err)
	}
	return nil
}

func Display() string {
	return fmt.Sprintf("%s\n%s", displayConfig(), displayDefaults())
}

func InitConfigAndDefaults() error {
	for k, v := range configDefinitions {
		if err := Set(k, v.defaultValue); err != nil {
			return err
		}
	}
	for k, v := range defaultsDefinitions {
		if err := Set(k, v.defaultValue); err != nil {
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

	db, err, dbclose := database.Current()
	if err != nil {
		return fmt.Errorf("set config: %s", err)
	}
	defer dbclose()
	if err := db.SetConfig(databaseKey, key, v); err != nil {
		return err
	}
	if def != nil && def.onUpdateFn != nil {
		def.onUpdateFn(v)
	}

	return nil
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
		db, err, dbclose := database.Current()
		if err != nil {
			return fmt.Errorf("unset config: %s", err)
		}
		err = db.UnsetConfig(dbKey, key)
		if err != nil {
			return fmt.Errorf("unset config: %s", err)
		}
		dbclose()
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
