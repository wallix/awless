package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

var StatsServerPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuUK69ARmXV0Xsj30+6S7
+oqDPmfIwQ0FxhlI6fcqlZ57mmURuZIJ4nnXxZrx5LXmbKGjDRgWtFLNQ2JFUGZB
y/vzBIxA64cEKE7Hkbh7MW6nQayoDOnb9ZPBqK5IjoGvnF0BsYoaKdP4Jy7Nbx9o
oBCBtu8q6WeCMBlGMnmFtjRCHPgpIf9/3vylFlNn6LRRG/DLO2xY4Is/wj2KM98O
6XWhU5PO7gC6ZX0BQpCOvTR1DOIXAW2JAyJHtJM4jFR5kBYY03dNqrKOIUyOWELo
pi+Bfy5FDK42Q/uJfUOJ5f6Ae/qIxxzKH7ixeXdCFvdzPvv4M4gGkqBAhpnFwLeX
SwIDAQAB
-----END PUBLIC KEY-----`

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
