package aws

import (
	"regexp"
	"sort"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

var DefaultAMIUsers = []string{"ec2-user", "ubuntu", "centos", "bitnami", "admin", "root"}

func AllRegions() []string {
	var regions sort.StringSlice
	partitions := endpoints.DefaultResolver().(endpoints.EnumPartitions).Partitions()
	for _, p := range partitions {
		for id := range p.Regions() {
			regions = append(regions, id)
		}
	}
	sort.Sort(regions)
	return regions
}

func IsValidRegion(given string) bool {
	reg, _ := regexp.Compile("^(us|eu|ap|sa|ca)\\-\\w+\\-\\d+$")
	regChina, _ := regexp.Compile("^cn\\-\\w+\\-\\d+$")
	regUsGov, _ := regexp.Compile("^us\\-gov\\-\\w+\\-\\d+$")

	return reg.MatchString(given) || regChina.MatchString(given) || regUsGov.MatchString(given)
}

type Secu struct {
	*sts.STS
}

func NewSecu(sess *session.Session) *Secu {
	return &Secu{sts.New(sess)}
}

func (s *Secu) CallerIdentity() (interface{}, error) {
	return s.GetCallerIdentity(&sts.GetCallerIdentityInput{})
}

func (s *Secu) GetUserId() (string, error) {
	output, err := s.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return awssdk.StringValue(output.Arn), nil
}

func (s *Secu) GetAccountId() (string, error) {
	output, err := s.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return awssdk.StringValue(output.Account), nil
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
	Users []*iam.User

	GroupsDetail []*iam.GroupDetail
	UsersDetail  []*iam.UserDetail
	RolesDetail  []*iam.RoleDetail
	Policies     []*iam.ManagedPolicyDetail

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

func (access *Access) global_fetch() (*AwsAccess, error) {
	resultc, errc := multiFetch(access.AccountDetails, access.fetch_all_user)

	awsAccess := NewAwsAccess()

	for r := range resultc {
		switch rr := r.(type) {
		case *iam.ListUsersOutput:
			awsAccess.Users = append(awsAccess.Users, rr.Users...)

		case *iam.GetAccountAuthorizationDetailsOutput:
			for _, user := range rr.UserDetailList {
				awsAccess.UsersDetail = append(awsAccess.UsersDetail, user)

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

			for _, group := range rr.GroupDetailList {
				awsAccess.GroupsDetail = append(awsAccess.GroupsDetail, group)

				policies := []string{}
				for _, policy := range group.GroupPolicyList {
					policies = append(policies, awssdk.StringValue(policy.PolicyName))
				}
				for _, policy := range group.AttachedManagedPolicies {
					policies = append(policies, awssdk.StringValue(policy.PolicyName))
				}
				awsAccess.GroupPolicies[awssdk.StringValue(group.GroupId)] = policies
			}

			for _, role := range rr.RoleDetailList {
				awsAccess.RolesDetail = append(awsAccess.RolesDetail, role)

				policies := []string{}
				for _, policy := range role.RolePolicyList {
					policies = append(policies, awssdk.StringValue(policy.PolicyName))
				}
				for _, policy := range role.AttachedManagedPolicies {
					policies = append(policies, awssdk.StringValue(policy.PolicyName))
				}
				awsAccess.RolePolicies[awssdk.StringValue(role.RoleId)] = policies
			}

			for _, policy := range rr.Policies {
				awsAccess.Policies = append(awsAccess.Policies, policy)
			}
		}
	}

	return awsAccess, <-errc
}
