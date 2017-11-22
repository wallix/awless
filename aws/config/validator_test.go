package awsconfig

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestRegionsValid(t *testing.T) {
	if got, want := stringInSlice("eu-west-1", allRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("us-east-1", allRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("us-west-1", allRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("eu-test-1", allRegions()), false; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	for _, k := range allRegions() {
		if got, want := IsValidRegion(k), true; got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	}
	if got, want := IsValidRegion("aa-test-10"), false; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}

func TestProfileValid(t *testing.T) {
	awsHomeTmp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(awsHomeTmp)
	}()

	awsHomeFunc = func() string {
		return awsHomeTmp
	}

	ioutil.WriteFile(filepath.Join(awsHomeTmp, "config"), []byte(`[profile mfa]
region = us-west-1
role_arn = arn:aws:iam::1234567890:role/my-role
source_profile = default
[janedoe-mfa]
source_profile = jdoe
mfa_serial = arn:aws:iam::1234567890:mfa/janedoe
role_arn = arn:aws:iam::1234567890:role/my-role
`), 0600)
	ioutil.WriteFile(filepath.Join(awsHomeTmp, "credentials"), []byte(`[default]
aws_access_key_id = ABCDEXAMPLE01234
aws_secret_access_key = aSecretKeyInMycredentials

[readonly]
aws_access_key_id =  ABCDEXAMPLE01234567
aws_secret_access_key = anotherSecretKeyInMycredentials
`), 0600)

	tcases := []struct {
		profile string
		expect  bool
	}{
		{profile: "", expect: false},
		{profile: "nothere", expect: false},
		{profile: "default", expect: true},
		{profile: "readonly", expect: true},
		{profile: "mfa", expect: true},
		{profile: "janedoe-mfa", expect: true},
	}
	for i, tcase := range tcases {
		if got, want := IsValidProfile(tcase.profile), tcase.expect; got != want {
			t.Fatalf("%d: '%s': got %t, want %t", i+1, tcase.profile, got, want)
		}
	}
}

func TestInstanceTypeValid(t *testing.T) {
	tcases := []struct {
		str    string
		expect bool
	}{
		{"t2.micro", true},
		{"m3.large", true},
		{"t.", false},
		{".", false},
		{"a.", false},
	}
	for _, tcase := range tcases {
		if got, want := isValidInstanceType(tcase.str), tcase.expect; got != want {
			t.Errorf("%s: got %t, want %t", tcase.str, got, want)
		}
	}
}
