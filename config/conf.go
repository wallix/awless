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
	if AwlessFirstInstall && len(os.Args) > 1 && os.Args[1] == "completion" {
		os.Exit(0)
	}

	os.MkdirAll(RepoDir, 0700)
	os.MkdirAll(KeysDir, 0700)

	var region string
	if AwlessFirstInstall {
		fmt.Println("First install. Welcome!")
		fmt.Println()
		region = resolveRegion()
		addDefaults(region)
	} else {
		db, close := database.Current()
		defer close()
		region = db.MustGetDefaultRegion()
	}

	os.Setenv("__AWLESS_REGION", region)

	return nil
}

func resolveRegion() (region string) {
	if sess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}); err == nil {
		region = awssdk.StringValue(sess.Config.Region)
	}

	if aws.IsValidRegion(region) {
		fmt.Printf("Found existing AWS region '%s'\n", region)
		fmt.Println("Setting it as your default region.")
		fmt.Println("Show config with `awless config list`. Change region with `awless config set region`")
		fmt.Println()
		return
	}

	fmt.Println("Could not find any AWS region in your environment.")

	region = askRegion()

	return
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

func addDefaults(region string) error {
	defaults := map[string]interface{}{
		database.RegionKey:        region,
		database.InstanceTypeKey:  "t2.micro",
		database.InstanceImageKey: "ami-9398d3e0",
		database.InstanceCountKey: 1,
	}
	db, close := database.Current()
	defer close()
	for k, v := range defaults {
		err := db.SetDefault(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
