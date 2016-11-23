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

	awsAccess := NewAwsAccess()

	for range allFetch {
		select {
		case r := <-resultc:
			switch r.(type) {
			case *iam.ListGroupsOutput:
				awsAccess.Groups = append(awsAccess.Groups, r.(*iam.ListGroupsOutput).Groups...)
			case *iam.ListUsersOutput:
				awsAccess.Users = append(awsAccess.Users, r.(*iam.ListUsersOutput).Users...)
			case *iam.ListRolesOutput:
				awsAccess.Roles = append(awsAccess.Roles, r.(*iam.ListRolesOutput).Roles...)
			}
		case e := <-errc:
			return awsAccess, e
		}
	}

	for _, group := range awsAccess.Groups {
		groupName := aws.StringValue(group.GroupName)
		groupId := aws.StringValue(group.GroupId)
		groupUsers, err := access.UsersForGroup(groupName)
		if err != nil {
			return awsAccess, err
		}

		for _, groupUser := range groupUsers.([]*iam.User) {
			awsAccess.UsersByGroup[groupId] = append(awsAccess.UsersByGroup[groupId], aws.StringValue(groupUser.UserId))
		}
	}

	return awsAccess, nil
}
