package cmd

import (
	"bytes"
	"fmt"
)
import "github.com/aws/aws-sdk-go/service/iam"

func display(format string, item interface{}, err error) {
	if err != nil {
		fmt.Println(err.Error())
	}

	switch format {
	case "raw":
		fmt.Println(item)
	default:
		lineDisplay(item)
	}
}

func lineDisplay(item interface{}) {
	var buf bytes.Buffer

	switch item.(type) {
	case *iam.ListUsersOutput:
		for _, user := range item.(*iam.ListUsersOutput).Users {
			buf.WriteString(fmt.Sprintf("id:%s, name:%s\n", *user.UserId, *user.UserName))
		}
	case *iam.ListGroupsOutput:
		for _, group := range item.(*iam.ListGroupsOutput).Groups {
			buf.WriteString(fmt.Sprintf("id:%s, name:%s\n", *group.GroupId, *group.GroupName))
		}
	case *iam.ListPoliciesOutput:
		for _, policy := range item.(*iam.ListPoliciesOutput).Policies {
			buf.WriteString(fmt.Sprintf("id:%s, name:%s\n", *policy.PolicyId, *policy.PolicyName))
		}
	default:
		fmt.Println("unknown entities to display")
	}

	fmt.Println(buf.String())
}
