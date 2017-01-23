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

func (a *Access) CallerIdentity() (interface{}, error) {
	return a.secu.GetCallerIdentity(&sts.GetCallerIdentityInput{})
}

func (a *Access) GetUserId() (string, error) {
	output, err := a.secu.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return awssdk.StringValue(output.Arn), nil
}

func (a *Access) GetAccountId() (string, error) {
	output, err := a.secu.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return awssdk.StringValue(output.Account), nil
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
	Groups   []*iam.GroupDetail
	Users    []*iam.UserDetail
	Roles    []*iam.RoleDetail
	Policies []*iam.ManagedPolicyDetail

	UserGroups map[string][]string

	UserPolicies  map[string][]string
	GroupPolicies map[string][]string
	RolePolicies  map[string][]string
}

func NewAwsAccess() *AwsAccess {
	return &AwsAccess{
		UserGroups:    make(map[string][]string),
		UserPolicies:  make(map[string][]string),
		GroupPolicies: make(map[string][]string),
		RolePolicies:  make(map[string][]string),
	}
}

func (access *Access) FetchAwsAccess() (*AwsAccess, error) {
	result, err := access.AccountDetails()
	if err != nil {
		return nil, err
	}

	account := result.(*iam.GetAccountAuthorizationDetailsOutput)

	awsAccess := NewAwsAccess()

	for _, user := range account.UserDetailList {
		awsAccess.Users = append(awsAccess.Users, user)

		groups := []string{}
		for _, groupId := range user.GroupList {
			groups = append(groups, awssdk.StringValue(groupId))
		}
		awsAccess.UserGroups[awssdk.StringValue(user.UserId)] = groups

		policies := []string{}
		for _, policy := range user.UserPolicyList {
			policies = append(policies, awssdk.StringValue(policy.PolicyName))
		}
		for _, policy := range user.AttachedManagedPolicies {
			policies = append(policies, awssdk.StringValue(policy.PolicyName))
		}
		awsAccess.UserPolicies[awssdk.StringValue(user.UserId)] = policies
	}

	for _, group := range account.GroupDetailList {
		awsAccess.Groups = append(awsAccess.Groups, group)

		policies := []string{}
		for _, policy := range group.GroupPolicyList {
			policies = append(policies, awssdk.StringValue(policy.PolicyName))
		}
		for _, policy := range group.AttachedManagedPolicies {
			policies = append(policies, awssdk.StringValue(policy.PolicyName))
		}
		awsAccess.GroupPolicies[awssdk.StringValue(group.GroupId)] = policies
	}

	for _, role := range account.RoleDetailList {
		awsAccess.Roles = append(awsAccess.Roles, role)

		policies := []string{}
		for _, policy := range role.RolePolicyList {
			policies = append(policies, awssdk.StringValue(policy.PolicyName))
		}
		for _, policy := range role.AttachedManagedPolicies {
			policies = append(policies, awssdk.StringValue(policy.PolicyName))
		}
		awsAccess.RolePolicies[awssdk.StringValue(role.RoleId)] = policies
	}

	for _, policy := range account.Policies {
		awsAccess.Policies = append(awsAccess.Policies, policy)
	}

	return awsAccess, nil
}
