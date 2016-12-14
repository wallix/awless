package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
)

var (
	configFilename       = "config.yaml"
	databaseFilename     = "awless.db"
	GitDir               = filepath.Join(os.Getenv("HOME"), ".awless", "aws", "rdf")
	Dir                  = filepath.Join(os.Getenv("HOME"), ".awless", "aws")
	Path                 = filepath.Join(Dir, configFilename)
	DatabasePath         = filepath.Join(os.Getenv("HOME"), ".awless", databaseFilename)
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
	Salt                                = "bg6B8yTTq8chwkN0BqWnEzlP4OkpcQDhO45jUOuXm1zsNGDLj3"
	StatsExpirationDuration             = 2 * time.Minute
	Version                             = "0.2"
	InfraFilename                       = "infra.rdf"
	AccessFilename                      = "access.rdf"
	AwlessFirstInstall, AwlessFirstSync bool
	Session                             *session.Session
)

func InitAwlessEnv() {
	_, ierr := os.Stat(filepath.Join(GitDir, InfraFilename))
	_, aerr := os.Stat(filepath.Join(GitDir, AccessFilename))
	AwlessFirstSync = os.IsNotExist(ierr) || os.IsNotExist(aerr)

	_, err := os.Stat(Path)
	AwlessFirstInstall = os.IsNotExist(err)

	os.MkdirAll(GitDir, 0700)

	if AwlessFirstInstall {
		ioutil.WriteFile(Path, []byte("region: \"eu-west-1\"\n"), 0600)
		fmt.Printf("Creating default config file at %s\n", Path)
	}

	viper.SetConfigFile(Path)

	if err = viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
		return
	}
}

func InitCloudSession() {
	var err error

	Session, err = session.NewSession(&aws.Config{Region: aws.String(viper.GetString("region"))})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
		return
	}
	if _, err = Session.Config.Credentials.Get(); err != nil {
		fmt.Println(err)
		fmt.Fprintln(os.Stderr, "Your AWS credentials seem undefined!")
		os.Exit(-1)
		return
	}
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
