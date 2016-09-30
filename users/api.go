package users

import "github.com/wallix/awless/config"

// API of users
type API struct {
	config *config.Config
}

// User defined a awless user
type User struct {
	UserName string
	Email    string
}

// NewAPI creates a new users API
func NewAPI(config *config.Config) *API {
	return &API{config: config}
}

// Create a new awless user
func (api *API) Create(u User) error {
	return api.config.Create(prefix+u.UserName, u)
}

// Exists returns true if the awless user exists
func (api *API) Exists(name string) bool {
	return api.config.Exists(prefix + name)
}

// List all awless users
func (api *API) List() []User {
	names := api.config.List(prefix)
	users := make([]User, len(names), len(names))
	for i, name := range names {
		err := api.config.Load(prefix+name, &users[i])
		if err != nil {
			api.config.Log.Fatalf("Internal error, cannot load user[%s]", name)
		}
	}
	return users
}
