package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type Access struct {
	*iam.IAM
}

func NewAccess(sess *session.Session) *Access {
	return &Access{iam.New(sess)}
}

func (a *Access) Users() (interface{}, error) {
	return a.ListUsers(&iam.ListUsersInput{})
}

func (a *Access) Groups() (interface{}, error) {
	return a.ListGroups(&iam.ListGroupsInput{})
}

func (a *Access) UsersForGroup(name string) (interface{}, error) {
	group, err := a.GetGroup(&iam.GetGroupInput{GroupName: aws.String(name)})
	if err != nil {
		return nil, err
	}
	return group.Users, nil
}

func (a *Access) Roles() (interface{}, error) {
	return a.ListRoles(&iam.ListRolesInput{})
}

func (a *Access) Policies() (interface{}, error) {
	return a.ListPolicies(&iam.ListPoliciesInput{})
}

type AwsAccess struct {
	Groups       []*iam.Group
	Users        []*iam.User
	Roles        []*iam.Role
	UsersByGroup map[string][]string
}

func NewAwsAccess() *AwsAccess {
	return &AwsAccess{
		UsersByGroup: make(map[string][]string),
	}
}

func (access *Access) FetchAccess() (*AwsAccess, error) {
	response := NewAwsAccess()

	type fetchFn func() (interface{}, error)

	allFetch := []fetchFn{access.Groups, access.Users, access.Roles}
	resultc := make(chan interface{})
	errc := make(chan error)

	for _, fetch := range allFetch {
		go func(fn fetchFn) {
			if r, err := fn(); err != nil {
				errc <- err
			} else {
				resultc <- r
			}
		}(fetch)
	}

	for range allFetch {
		select {
		case r := <-resultc:
			switch r.(type) {
			case *iam.ListGroupsOutput:
				response.Groups = append(response.Groups, r.(*iam.ListGroupsOutput).Groups...)
			case *iam.ListUsersOutput:
				response.Users = append(response.Users, r.(*iam.ListUsersOutput).Users...)
			case *iam.ListRolesOutput:
				response.Roles = append(response.Roles, r.(*iam.ListRolesOutput).Roles...)
			}
		case e := <-errc:
			return response, e
		}
	}

	for _, group := range response.Groups {
		groupName := aws.StringValue(group.GroupName)
		groupId := aws.StringValue(group.GroupId)
		groupUsers, err := access.UsersForGroup(groupName)
		if err != nil {
			return response, err
		}

		for _, groupUser := range groupUsers.([]*iam.User) {
			response.UsersByGroup[groupId] = append(response.UsersByGroup[groupId], aws.StringValue(groupUser.UserId))
		}
	}

	return response, nil
}
