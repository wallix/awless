package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
)
import "github.com/aws/aws-sdk-go/service/iam"

func display(item interface{}, err error, format ...string) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(format) < 1 {
		jsonDisplay(item)
		return
	}

	switch format[0] {
	case "raw":
		fmt.Println(item)
	default:
		lineDisplay(item)
	}
}

func jsonDisplay(item interface{}) {
	j, err := json.MarshalIndent(item, "", " ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%s\n", j)
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
		fmt.Println(item)
		return
	}

	fmt.Println(buf.String())
}
