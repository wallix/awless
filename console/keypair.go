package console

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

func GenerateSSHKeyPair(size int) ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, nil, err
	}

	privPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	sshPub, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return ssh.MarshalAuthorizedKey(sshPub), privPem, nil
}
