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
)

var (
	configFilename       = "config.yaml"
	databaseFilename     = "database.db"
	Dir                  = filepath.Join(os.Getenv("HOME"), ".awless", "aws")
	Path                 = filepath.Join(Dir, configFilename)
	DatabasePath         = filepath.Join(os.Getenv("HOME"), ".awless", databaseFilename)
	StatsServerUrl       = "http://127.0.0.1:8080"
	StatsServerPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuUK69ARmXV0Xsj30+6S7
+oqDPmfIwQ0FxhlI6fcqlZ57mmURuZIJ4nnXxZrx5LXmbKGjDRgWtFLNQ2JFUGZB
y/vzBIxA64cEKE7Hkbh7MW6nQayoDOnb9ZPBqK5IjoGvnF0BsYoaKdP4Jy7Nbx9o
oBCBtu8q6WeCMBlGMnmFtjRCHPgpIf9/3vylFlNn6LRRG/DLO2xY4Is/wj2KM98O
6XWhU5PO7gC6ZX0BQpCOvTR1DOIXAW2JAyJHtJM4jFR5kBYY03dNqrKOIUyOWELo
pi+Bfy5FDK42Q/uJfUOJ5f6Ae/qIxxzKH7ixeXdCFvdzPvv4M4gGkqBAhpnFwLeX
SwIDAQAB
-----END PUBLIC KEY-----`
	StatsExpirationDuration = 2 * time.Minute

	InfraFilename  = "infra.rdf"
	AccessFilename = "access.rdf"
)

func CreateDefaultConf() {
	os.MkdirAll(Dir, 0700)

	if _, err := os.Stat(Path); os.IsNotExist(err) {
		ioutil.WriteFile(Path, []byte("region: \"eu-west-1\"\n"), 0400)
		fmt.Printf("Creating default config file at %s\n", Path)
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
