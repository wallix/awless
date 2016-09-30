package groups

import (
	"sort"

	"github.com/apex/log"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/users"
)

// AddParam input parameters of Create
type AddParam struct {
	GroupName string
	UserName  string
}

func (c *Commands) AddUser(p AddParam) {
	l := c.config.Log.WithFields(log.Fields{
		"group-name": p.GroupName,
		"user-name":  p.UserName,
	})
	l.Infof("adding user to group")
	group := Group{}
	err := c.config.Load(prefix+p.GroupName, &group)
	if err == config.ErrNotExist {
		l.Errorf("group doesn't exists")
		return
	}
	uc := users.NewAPI(c.config)
	if !uc.Exists(p.UserName) {
		l.Errorf("user doesn't exists")
		return
	}
	group.Users = append(group.Users, p.UserName)
	sort.Strings(group.Users)
	c.config.Save(prefix+p.GroupName, group)
}
