package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/database"
)

var (
	databaseFilename                    = "awless.db"
	AwlessHome                          = filepath.Join(os.Getenv("HOME"), ".awless")
	GitDir                              = filepath.Join(AwlessHome, "aws", "rdf")
	Dir                                 = filepath.Join(AwlessHome, "aws")
	KeysDir                             = filepath.Join(AwlessHome, "keys")
	DatabasePath                        = filepath.Join(AwlessHome, databaseFilename)
	StatsServerUrl                      = "http://52.213.243.16:8080"
	StatsExpirationDuration             = 24 * time.Hour
	Version                             = "0.0.2"
	InfraFilename                       = "infra.rdf"
	AccessFilename                      = "access.rdf"
	AwlessFirstInstall, AwlessFirstSync bool
)

func InitAwlessEnv() error {
	_, err := os.Stat(AwlessHome)
	_, ierr := os.Stat(filepath.Join(GitDir, InfraFilename))
	_, aerr := os.Stat(filepath.Join(GitDir, AccessFilename))
	AwlessFirstSync = os.IsNotExist(ierr) || os.IsNotExist(aerr)

	AwlessFirstInstall = os.IsNotExist(err)
	if AwlessFirstInstall && len(os.Args) > 1 && os.Args[1] == "completion" {
		os.Exit(0)
	}

	os.MkdirAll(GitDir, 0700)
	os.MkdirAll(KeysDir, 0700)

	if err = database.Open(DatabasePath); err != nil {
		return err
	}

	if AwlessFirstInstall {
		fmt.Println("First install. Welcome!")
		fmt.Println()
		region := resolveRegion()
		addDefaults(region)
	}

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
		RegionKey:        region,
		InstanceTypeKey:  "t2.micro",
		InstanceImageKey: "ami-9398d3e0",
		InstanceCountKey: 1,
	}
	for k, v := range defaults {
		err := database.Current.SetDefault(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
