package groups

type ListParams struct {
}

// List all users
func (c *Commands) List(_ struct{}) {
	names := c.config.List(prefix)
	for _, name := range names {
		c.config.Log.Infof(name)
	}
}
