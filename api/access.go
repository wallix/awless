package api

import (
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

func (a *Access) Users() (interface{}, error) {
	return a.ListUsers(&iam.ListUsersInput{})
}

func (a *Access) Groups() (interface{}, error) {
	return a.ListGroups(&iam.ListGroupsInput{})
}

func (a *Access) Roles() (interface{}, error) {
	return a.ListRoles(&iam.ListRolesInput{})
}

func (a *Access) Policies() (interface{}, error) {
	return a.ListPolicies(&iam.ListPoliciesInput{})
}
