package awsfetch

import (
	"context"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/wallix/awless/fetch"
)

type AccountAuthorizationDetails struct {
	Groups   []*iam.GroupDetail
	Policies []*iam.ManagedPolicyDetail
	Roles    []*iam.RoleDetail
	Users    []*iam.UserDetail
}

func getAccountAuthorizationDetails(ctx context.Context, cache fetch.Cache, api iamiface.IAMAPI) (*AccountAuthorizationDetails, error) {
	var entities []*string
	var cacheKey string
	resourceType, ok := fetch.IsFetchingByType(ctx)
	if ok {
		switch resourceType {
		case "user":
			cacheKey = "usersDetails"
			entities = append(entities, awssdk.String(iam.EntityTypeUser))
		case "group":
			cacheKey = "groupsDetails"
			entities = append(entities, awssdk.String(iam.EntityTypeGroup))
		case "role":
			cacheKey = "rolesDetails"
			entities = append(entities, awssdk.String(iam.EntityTypeRole))
		case "policy":
			cacheKey = "policiesDetails"
			entities = append(entities, awssdk.String(iam.EntityTypeLocalManagedPolicy), awssdk.String(iam.EntityTypeAwsmanagedPolicy))
		}
	} else {
		cacheKey = "accountDetails"
		entities = append(entities, awssdk.String(iam.EntityTypeUser), awssdk.String(iam.EntityTypeGroup), awssdk.String(iam.EntityTypeRole))
		entities = append(entities, awssdk.String(iam.EntityTypeLocalManagedPolicy), awssdk.String(iam.EntityTypeAwsmanagedPolicy))
	}

	if val, err := cache.Get(cacheKey, func() (interface{}, error) {
		return fetchAccountAuthorizationDetails(entities, api)
	}); err != nil {
		return nil, err
	} else if v, ok := val.(*AccountAuthorizationDetails); ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("cannot get account details (val of type %T)", val)
	}
}

func fetchAccountAuthorizationDetails(entities []*string, api iamiface.IAMAPI) (*AccountAuthorizationDetails, error) {
	details := new(AccountAuthorizationDetails)
	err := api.GetAccountAuthorizationDetailsPages(&iam.GetAccountAuthorizationDetailsInput{
		Filter: entities,
	}, func(out *iam.GetAccountAuthorizationDetailsOutput, lastPage bool) (shouldContinue bool) {
		for _, u := range out.UserDetailList {
			details.Users = append(details.Users, u)
		}
		for _, g := range out.GroupDetailList {
			details.Groups = append(details.Groups, g)
		}
		for _, r := range out.RoleDetailList {
			details.Roles = append(details.Roles, r)
		}
		for _, p := range out.Policies {
			details.Policies = append(details.Policies, p)
		}
		return out.Marker != nil
	})

	return details, err
}
