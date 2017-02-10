/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
