package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/database"
)

var (
	AwlessHome                          = filepath.Join(os.Getenv("HOME"), ".awless")
	RepoDir                             = filepath.Join(AwlessHome, "aws", "rdf")
	Dir                                 = filepath.Join(AwlessHome, "aws")
	KeysDir                             = filepath.Join(AwlessHome, "keys")
	InfraFilename                       = "infra.rdf"
	AccessFilename                      = "access.rdf"
	AwlessFirstInstall, AwlessFirstSync bool
)

func InitAwlessEnv() error {
	os.Setenv("__AWLESS_HOME", AwlessHome)
	_, err := os.Stat(AwlessHome)
	_, ierr := os.Stat(filepath.Join(RepoDir, InfraFilename))
	_, aerr := os.Stat(filepath.Join(RepoDir, AccessFilename))
	AwlessFirstSync = os.IsNotExist(ierr) || os.IsNotExist(aerr)

	AwlessFirstInstall = os.IsNotExist(err)

	os.MkdirAll(RepoDir, 0700)
	os.MkdirAll(KeysDir, 0700)

	var region string

	if AwlessFirstInstall {
		fmt.Println("First install. Welcome!\n")
		region, err = resolveAndSetDefaults()
		if err != nil {
			return err
		}
	} else {
		db, err, close := database.Current()
		if err != nil {
			return fmt.Errorf("init env: database error: ", err)
		}
		defer close()
		region = db.MustGetDefaultRegion()
	}

	os.Setenv("__AWLESS_REGION", region)

	return nil
}

func resolveAndSetDefaults() (string, error) {
	var region, ami string

	if sess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}); err == nil {
		region = awssdk.StringValue(sess.Config.Region)
	}

	if aws.IsValidRegion(region) {
		fmt.Printf("Found existing AWS region '%s'\n", region)
		fmt.Println("Setting it as your default region.")
		fmt.Println("Show config with `awless config list`. Change region with `awless config set region`")
		fmt.Println()
	} else {
		fmt.Println("Could not find any AWS region in your environment.")
		region = askRegion()
	}

	var hasAMI bool
	if ami, hasAMI = amiPerRegion[region]; !hasAMI {
		fmt.Printf("Could not find a default ami for your region %s\n. Set it manually with `awless config set region ...`", region)
	}

	defaults := map[string]interface{}{
		database.SyncAuto:         true,
		database.RegionKey:        region,
		database.InstanceTypeKey:  "t1.micro",
		database.InstanceCountKey: 1,
	}

	if hasAMI {
		defaults[database.InstanceImageKey] = ami
	}

	db, err, close := database.Current()
	if err != nil {
		return region, fmt.Errorf("database error: ", err)
	}
	defer close()
	for k, v := range defaults {
		err := db.SetDefault(k, v)
		if err != nil {
			return region, err
		}
	}

	fmt.Println("\nThose defaults have been set in your config (manage them with `awless config`):")
	for k, v := range defaults {
		fmt.Printf("\t%s = %v\n", k, v)
	}
	fmt.Println("\nAll done. Enjoy!\n")

	return region, nil
}

func askRegion() string {
	var region string
	fmt.Println("Please choose one region:")

	fmt.Println(strings.Join(aws.AllRegions(), ", "))
	fmt.Println()
	fmt.Print("Copy/paste one region > ")
	fmt.Scan(&region)
	for !aws.IsValidRegion(region) {
		fmt.Print("Invalid! Copy/paste one valid region > ")
		fmt.Scan(&region)
	}
	return region
}

var amiPerRegion = map[string]string{
	"us-east-1":      "ami-1b814f72",
	"us-west-2":      "ami-30fe7300",
	"us-west-1":      "ami-11d68a54",
	"eu-west-1":      "ami-973b06e3",
	"ap-southeast-1": "ami-b4b0cae6",
	"ap-southeast-2": "ami-b3990e89",
	"ap-northeast-1": "ami-0644f007",
	"sa-east-1":      "ami-3e3be423",
}
