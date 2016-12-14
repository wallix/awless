package mockup

import (
	"crypto/sha256"
	"fmt"

	"github.com/wallix/awless/cloud"
)

type Mockup struct {
}

func (m *Mockup) GetUserId() (string, error) {
	return fmt.Sprintf("%x", sha256.Sum256([]byte("user"))), nil
}

func (m *Mockup) GetAccountId() (string, error) {
	return fmt.Sprintf("%x", sha256.Sum256([]byte("account"))), nil
}

func InitMockup() {
	mockup := &Mockup{}
	cloud.Current = mockup
}
