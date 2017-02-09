package stats

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/database"
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

func newTestDb() (*database.DB, func()) {
	f, e := ioutil.TempDir(".", "test")
	if e != nil {
		panic(e)
	}

	os.Setenv("__AWLESS_HOME", f)

	database.InitDB(true)
	db, closing := database.MustGetCurrent()

	todefer := func() {
		closing()
		os.RemoveAll(f)
	}

	defaults := map[string]interface{}{
		database.RegionKey:        "eu-west-1",
		database.InstanceTypeKey:  "t2.micro",
		database.InstanceImageKey: "ami-9398d3e0",
		database.InstanceCountKey: 1,
	}
	for k, v := range defaults {
		err := db.SetDefault(k, v)
		if err != nil {
			panic(err)
		}
	}

	return db, todefer
}
