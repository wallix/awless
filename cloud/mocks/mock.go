package mocks

import (
	"crypto/sha256"
	"fmt"

	"github.com/wallix/awless/cloud"
)

type Mock struct{}

func (m *Mock) GetUserId() (string, error) {
	return fmt.Sprintf("%x", sha256.Sum256([]byte("user"))), nil
}

func (m *Mock) GetAccountId() (string, error) {
	return fmt.Sprintf("%x", sha256.Sum256([]byte("account"))), nil
}

func InitServices() {
	cloud.Current = &Mock{}
}
