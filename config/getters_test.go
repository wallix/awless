package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestGetSyncEnabled(t *testing.T) {
	f, e := ioutil.TempDir(".", "test")
	if e != nil {
		t.Fatal(e)
	}
	defer os.RemoveAll(f)

	os.Setenv("__AWLESS_HOME", f)

	t.Run("Resource configuration", func(t *testing.T) {
		configDefinitions = map[string]*Definition{
			"aws.region": {help: "AWS region", defaultValue: ""},
		}

		if err := InitConfig(map[string]string{RegionConfigKey: "eu-west-1"}); err != nil {
			t.Fatal(err)
		}
		if err := LoadConfig(); err != nil {
			t.Fatal(err)
		}
		Set("aws.ec2.sync", "true")
		Set("aws.ec2.subnet.sync", "true")
		Set("aws.ec2.instance.sync", "false")
		Set("aws.iam.group.sync", "true")
		Set("aws.iam.user.sync", "false")
		Set("other.iam.user.sync", "false")
		expect := map[string]interface{}{
			"aws.region":            "eu-west-1",
			"aws.ec2.sync":          true,
			"aws.ec2.subnet.sync":   true,
			"aws.ec2.instance.sync": false,
			"aws.iam.group.sync":    true,
			"aws.iam.user.sync":     false,
		}
		if got, want := GetConfigWithPrefix("aws."), expect; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %+v, want %+v", got, want)
		}
	})
}
