package deploy

import (
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/provider"
	"github.com/wallix/awless/users"
)

// API contains method for deploying
type API struct {
	config   *config.Config
	provider provider.Interface
}

// NewAPI creates a new API for the given config
func NewAPI(config *config.Config, p provider.Interface) *API {
	return &API{
		config:   config,
		provider: p,
	}
}

// OPKind describe several kind of operations
type OPKind int

const (
	// Add operation
	Add OPKind = iota

	// Del operation
	Del

	// Noop operation
	Noop
)

// Symbol human readable symbol of an operation
func (k OPKind) Symbol() string {
	switch k {
	case Add:
		return "+"
	case Del:
		return "-"
	case Noop:
		return "#"
	}
	panic("OpKind.Symbol")
}

// UserOP describe an operation on a user
type UserOP struct {
	User users.User
	Kind OPKind
}

// UsersPatch returns the list of operations to complete the diff
func (api *API) UsersPatch(noop bool) ([]UserOP, error) {
	uapi := users.NewAPI(api.config)
	ausers := uapi.List()
	pusers, err := api.provider.ListUserNames()
	if err != nil {
		return nil, err
	}
	ops := make([]UserOP, 0, len(ausers))
	for ia, ip, la, lp := 0, 0, len(ausers), len(pusers); ia < la || ip < lp; {
		if ia >= la || ausers[ia].UserName > pusers[ip] {
			op := UserOP{User: users.User{UserName: pusers[ip]}, Kind: Del}
			ops = append(ops, op)
			ip++
			continue
		}
		if ip >= lp || ausers[ia].UserName < pusers[ip] {
			op := UserOP{User: ausers[ia], Kind: Add}
			ops = append(ops, op)
			ia++
			continue
		}
		if ausers[ia].UserName == pusers[ip] {
			if noop {
				ops = append(ops, UserOP{User: ausers[ia], Kind: Noop})
			}
			ia++
			ip++
			continue
		}
		panic("?")
	}
	return ops, nil
}
