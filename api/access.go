package api

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type Access struct {
	*iam.IAM
}

func NewAccess() (*Access, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &Access{iam.New(sess)}, nil
}

func (a *Access) Users() string {
	if resp, err := a.ListUsers(&iam.ListUsersInput{}); err != nil {
		return err.Error()
	} else {
		return fmt.Sprint(resp)
	}
}

func (a *Access) Groups() string {
	if resp, err := a.ListGroups(&iam.ListGroupsInput{}); err != nil {
		return err.Error()
	} else {
		return fmt.Sprint(resp)
	}
}

func (a *Access) Roles() string {
	if resp, err := a.ListRoles(&iam.ListRolesInput{}); err != nil {
		return err.Error()
	} else {
		return fmt.Sprint(resp)
	}
}

func (a *Access) Policies() string {
	if resp, err := a.ListPolicies(&iam.ListPoliciesInput{}); err != nil {
		return err.Error()
	} else {
		return fmt.Sprint(resp)
	}
}
