package stats

import (
	"io/ioutil"
	"os"

	"github.com/wallix/awless/cloud/mocks"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
)

func init() {
	mocks.InitServices()
}

func newTestDb() (*database.DB, func()) {
	f, e := ioutil.TempFile(".", "test.db")
	if e != nil {
		panic(e)
	}

	err := database.Open(f.Name())
	if err != nil {
		panic(err)
	}

	todefer := func() {
		os.Remove(f.Name())
		database.Current.Close()
	}

	defaults := map[string]interface{}{
		config.RegionKey:        "eu-west-1",
		config.InstanceTypeKey:  "t2.micro",
		config.InstanceBaseKey:  "ami-9398d3e0",
		config.InstanceCountKey: 1,
	}
	for k, v := range defaults {
		err := database.Current.SetDefault(k, v)
		if err != nil {
			panic(err)
		}
	}

	return database.Current, todefer
}
