package provider

import "github.com/apex/log"

const (
	prefix = "/provider/"
)

// ConfigureParam intput params of Configure
type ConfigureParam struct {
	Name string
	Kind string
}

// Configure a provider
func (c *Commands) Configure(p ConfigureParam) {
	l := c.config.Log.WithFields(&log.Fields{
		"name": p.Name,
		"kind": p.Kind,
	})
	cp, ok := cprov[p.Kind]
	if !ok {
		l.Errorf("unknown provider kind")
		return
	}
	i, e := cp.ReadConfig()
	if e != nil {
		l.WithError(e).Errorf("cannot read provider configuration")
		return
	}
	e = cp.CheckConfig(i)
	if e != nil {
		l.WithError(e).Errorf("config check error")
		return
	}
	l.Infof("success")
	c.config.Save(prefix+p.Name+"/config", p)
	c.config.Save(prefix+p.Name+"/custom", i)
	c.config.Save(prefix+".default", p)
}
