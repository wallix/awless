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

func (access *Access) FetchAccess() ([]*iam.Group, []*iam.User, map[string][]string, error) {
	var fetchErr error
	var groups []*iam.Group
	var users []*iam.User
	usersByGroup := make(map[string][]string)

	type fetchFn func() (interface{}, error)

	allFetch := []fetchFn{access.Groups, access.Users}
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
				groups = append(groups, r.(*iam.ListGroupsOutput).Groups...)
			case *iam.ListUsersOutput:
				users = append(users, r.(*iam.ListUsersOutput).Users...)
			}
		case fetchErr = <-errc:
			return groups, users, usersByGroup, fetchErr
		}
	}

	for _, group := range groups {
		groupName := aws.StringValue(group.GroupName)
		groupId := aws.StringValue(group.GroupId)
		groupUsers, err := access.UsersForGroup(groupName)
		if err != nil {
			return groups, users, usersByGroup, err
		}

		for _, groupUser := range groupUsers.([]*iam.User) {
			usersByGroup[groupId] = append(usersByGroup[groupId], aws.StringValue(groupUser.UserId))
		}
	}

	return groups, users, usersByGroup, fetchErr
}
