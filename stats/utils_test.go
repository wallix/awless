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
	f, e := ioutil.TempDir(".", "test")
	if e != nil {
		panic(e)
	}

	os.Setenv("__AWLESS_HOME", f)

	database.InitDB(true)
	db, closing := database.Current()

	todefer := func() {
		closing()
		os.RemoveAll(f)
	}

	defaults := map[string]interface{}{
		config.RegionKey:        "eu-west-1",
		config.InstanceTypeKey:  "t2.micro",
		config.InstanceImageKey: "ami-9398d3e0",
		config.InstanceCountKey: 1,
	}
	for k, v := range defaults {
		err := db.SetDefault(k, v)
		if err != nil {
			panic(err)
		}
	}

	return db, todefer
}
