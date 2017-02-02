package database

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/wallix/awless/cloud/aws"
)

type secuMock struct {
	stsiface.STSAPI
}

func (m *secuMock) GetUserId() (string, error) {
	return fmt.Sprintf("%x", sha256.Sum256([]byte("user"))), nil
}

func (m *secuMock) GetAccountId() (string, error) {
	return fmt.Sprintf("%x", sha256.Sum256([]byte("account"))), nil
}

func init() {
	aws.SecuAPI = &secuMock{}
}

func newTestDb() (*DB, func()) {
	f, e := ioutil.TempDir(".", "test")
	if e != nil {
		panic(e)
	}

	os.Setenv("__AWLESS_HOME", f)

	InitDB(true)
	db, closing := Current()

	todefer := func() {
		closing()
		os.RemoveAll(f)
	}

	return db, todefer
}
