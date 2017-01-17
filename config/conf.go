package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/wallix/awless/database"
)

var (
	databaseFilename     = "awless.db"
	AwlessHome           = filepath.Join(os.Getenv("HOME"), ".awless")
	GitDir               = filepath.Join(AwlessHome, "aws", "rdf")
	Dir                  = filepath.Join(AwlessHome, "aws")
	KeysDir              = filepath.Join(AwlessHome, "keys")
	DatabasePath         = filepath.Join(AwlessHome, databaseFilename)
	StatsServerUrl       = "http://52.213.243.16:8080"
	StatsServerPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuUK69ARmXV0Xsj30+6S7
+oqDPmfIwQ0FxhlI6fcqlZ57mmURuZIJ4nnXxZrx5LXmbKGjDRgWtFLNQ2JFUGZB
y/vzBIxA64cEKE7Hkbh7MW6nQayoDOnb9ZPBqK5IjoGvnF0BsYoaKdP4Jy7Nbx9o
oBCBtu8q6WeCMBlGMnmFtjRCHPgpIf9/3vylFlNn6LRRG/DLO2xY4Is/wj2KM98O
6XWhU5PO7gC6ZX0BQpCOvTR1DOIXAW2JAyJHtJM4jFR5kBYY03dNqrKOIUyOWELo
pi+Bfy5FDK42Q/uJfUOJ5f6Ae/qIxxzKH7ixeXdCFvdzPvv4M4gGkqBAhpnFwLeX
SwIDAQAB
-----END PUBLIC KEY-----`
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

	err = database.Open(DatabasePath)
	if err != nil {
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

func LoadPublicKey() (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(StatsServerPublicKey))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER encoded public key: " + err.Error())
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, fmt.Errorf("unknown type of public key")
	}
}

func resolveRegion() (region string) {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
	if sess, err := session.NewSession(); err == nil {
		region = aws.StringValue(sess.Config.Region)
	}

	if validRegion(region) {
		fmt.Printf("Found existing AWS region '%s'\n", region)
		fmt.Println("Setting it as your default region. Show config with `awless config`")
		fmt.Println()
		return
	}

	fmt.Println("Could not find any AWS region in your environment. Please choose one:")

	fmt.Println(strings.Join(allRegions(), ", "))
	fmt.Println()
	fmt.Print("Copy/paste one region > ")
	fmt.Scan(&region)
	for !validRegion(region) {
		fmt.Print("Invalid! Copy/paste one valid region > ")
		fmt.Scan(&region)
	}

	return
}

func addDefaults(region string) error {
	defaults := map[string]interface{}{
		RegionKey:        region,
		InstanceTypeKey:  "t2.micro",
		InstanceBaseKey:  "ami-9398d3e0",
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

func allRegions() []string {
	var regions []string
	partitions := endpoints.DefaultResolver().(endpoints.EnumPartitions).Partitions()
	for _, p := range partitions {
		for id := range p.Regions() {
			regions = append(regions, id)
		}
	}
	return regions
}

func validRegion(given string) bool {
	reg, _ := regexp.Compile("^(us|eu|ap|sa|ca)\\-\\w+\\-\\d+$")
	regChina, _ := regexp.Compile("^cn\\-\\w+\\-\\d+$")
	regUsGov, _ := regexp.Compile("^us\\-gov\\-\\w+\\-\\d+$")

	return reg.MatchString(given) || regChina.MatchString(given) || regUsGov.MatchString(given)
}
