package aws

import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

type Access struct {
	*iam.IAM
	secu *sts.STS
}

func NewAccess(sess *session.Session) *Access {
	return &Access{IAM: iam.New(sess), secu: sts.New(sess)}
}

func (a *Access) Users() (interface{}, error) {
	return a.ListUsers(&iam.ListUsersInput{})
}

func (a *Access) Groups() (interface{}, error) {
	return a.ListGroups(&iam.ListGroupsInput{})
}

func (a *Access) UsersForGroup(name string) (interface{}, error) {
	group, err := a.GetGroup(&iam.GetGroupInput{GroupName: awssdk.String(name)})
	if err != nil {
		return nil, err
	}
	return group.Users, nil
}

func (a *Access) Roles() (interface{}, error) {
	return a.ListRoles(&iam.ListRolesInput{})
}

func (a *Access) CallerIdentity() (*sts.GetCallerIdentityOutput, error) {
	return a.secu.GetCallerIdentity(&sts.GetCallerIdentityInput{})
}

func (a *Access) LocalPolicies() (interface{}, error) {
	return a.ListPolicies(&iam.ListPoliciesInput{Scope: awssdk.String(iam.PolicyScopeTypeLocal)})
}

func (a *Access) AttachedToPolicy(policyArn *string) (interface{}, error) {
	params := &iam.ListEntitiesForPolicyInput{
		PolicyArn: policyArn,
	}
	return a.ListEntitiesForPolicy(params)
}

func (a *Access) AccountDetails() (interface{}, error) {
	params := &iam.GetAccountAuthorizationDetailsInput{
		Filter: []*string{
			awssdk.String(iam.EntityTypeUser),
			awssdk.String(iam.EntityTypeRole),
			awssdk.String(iam.EntityTypeGroup),
			awssdk.String(iam.EntityTypeLocalManagedPolicy),
			awssdk.String(iam.EntityTypeAwsmanagedPolicy),
		},
	}
	return a.GetAccountAuthorizationDetails(params)
}

type AwsAccess struct {
	Groups       []*iam.Group
	Users        []*iam.User
	Roles        []*iam.Role
	UsersByGroup map[string][]string

	LocalPolicies         []*iam.Policy
	UsersByLocalPolicies  map[string][]string
	GroupsByLocalPolicies map[string][]string
	RolesByLocalPolicies  map[string][]string
}

func NewAwsAccess() *AwsAccess {
	return &AwsAccess{
		UsersByGroup:          make(map[string][]string),
		UsersByLocalPolicies:  make(map[string][]string),
		GroupsByLocalPolicies: make(map[string][]string),
		RolesByLocalPolicies:  make(map[string][]string),
	}
}

func (access *Access) FetchAwsAccess() (*AwsAccess, error) {
	resultc, errc := multiFetch(access.Groups, access.Users, access.Roles, access.LocalPolicies)

	awsAccess := NewAwsAccess()

	for r := range resultc {
		switch r.(type) {
		case *iam.ListGroupsOutput:
			awsAccess.Groups = append(awsAccess.Groups, r.(*iam.ListGroupsOutput).Groups...)
		case *iam.ListUsersOutput:
			awsAccess.Users = append(awsAccess.Users, r.(*iam.ListUsersOutput).Users...)
		case *iam.ListRolesOutput:
			awsAccess.Roles = append(awsAccess.Roles, r.(*iam.ListRolesOutput).Roles...)
		case *iam.ListPoliciesOutput:
			awsAccess.LocalPolicies = append(awsAccess.LocalPolicies, r.(*iam.ListPoliciesOutput).Policies...)
		}
	}

	if err := <-errc; err != nil {
		return awsAccess, err
	}

	for _, policy := range awsAccess.LocalPolicies {
		resp, err := access.AttachedToPolicy(policy.Arn)
		if err != nil {
			return awsAccess, err
		}

		output := resp.(*iam.ListEntitiesForPolicyOutput)
		for _, group := range output.PolicyGroups {
			awsAccess.GroupsByLocalPolicies[awssdk.StringValue(policy.PolicyId)] = append(awsAccess.GroupsByLocalPolicies[awssdk.StringValue(policy.PolicyId)], awssdk.StringValue(group.GroupId))
		}
		for _, role := range output.PolicyRoles {
			awsAccess.RolesByLocalPolicies[awssdk.StringValue(policy.PolicyId)] = append(awsAccess.RolesByLocalPolicies[awssdk.StringValue(policy.PolicyId)], awssdk.StringValue(role.RoleId))
		}
		for _, user := range output.PolicyUsers {
			awsAccess.UsersByLocalPolicies[awssdk.StringValue(policy.PolicyId)] = append(awsAccess.UsersByLocalPolicies[awssdk.StringValue(policy.PolicyId)], awssdk.StringValue(user.UserId))
		}
	}

	for _, group := range awsAccess.Groups {
		groupName := awssdk.StringValue(group.GroupName)
		groupId := awssdk.StringValue(group.GroupId)
		groupUsers, err := access.UsersForGroup(groupName)
		if err != nil {
			return awsAccess, err
		}

		for _, groupUser := range groupUsers.([]*iam.User) {
			awsAccess.UsersByGroup[groupId] = append(awsAccess.UsersByGroup[groupId], awssdk.StringValue(groupUser.UserId))
		}
	}

	return awsAccess, nil
}
