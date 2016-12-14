package aws

import (
	"fmt"
	"io"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/fatih/color"
)

var simpleDay = "Mon, Jan 2, 2006"

func TabularDisplay(item interface{}, w io.Writer) {
	switch item.(type) {
	case *iam.ListUsersOutput:
		fmt.Fprintln(w, "Name\tId\tCreated")
		for _, user := range item.(*iam.ListUsersOutput).Users {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s", *user.UserName, *user.UserId, (*user.CreateDate).Format(simpleDay)))
		}
	case *iam.ListGroupsOutput:
		fmt.Fprintln(w, "Name\tId\tCreated")
		for _, group := range item.(*iam.ListGroupsOutput).Groups {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s", *group.GroupName, *group.GroupId, (*group.CreateDate).Format(simpleDay)))
		}
	case *iam.ListRolesOutput:
		fmt.Fprintln(w, "Name\tId\tCreated")
		for _, role := range item.(*iam.ListRolesOutput).Roles {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s", *role.RoleName, *role.RoleId, (*role.CreateDate).Format(simpleDay)))
		}
	case *iam.ListPoliciesOutput:
		fmt.Fprintln(w, "Name\tId\tCreated")
		for _, policy := range item.(*iam.ListPoliciesOutput).Policies {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s", *policy.PolicyName, *policy.PolicyId, (*policy.CreateDate).Format(simpleDay)))
		}
	case *ec2.DescribeInstancesOutput:
		fmt.Fprintln(w, "Id\tName\tState\tType\tPriv IP\tPub IP")
		for _, reserv := range item.(*ec2.DescribeInstancesOutput).Reservations {
			for _, inst := range reserv.Instances {
				var name string
				for _, t := range inst.Tags {
					if *t.Key == "Name" {
						name = *t.Value
					}
				}
				fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", awssdk.StringValue(inst.InstanceId), name, awssdk.StringValue(inst.State.Name), awssdk.StringValue(inst.InstanceType), awssdk.StringValue(inst.PrivateIpAddress), awssdk.StringValue(inst.PublicIpAddress)))
			}
		}
	case *ec2.Reservation:
		fmt.Fprintln(w, "Id\tType\tState\tPriv IP\tPub IP")
		for _, inst := range item.(*ec2.Reservation).Instances {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s", awssdk.StringValue(inst.InstanceId), awssdk.StringValue(inst.State.Name), awssdk.StringValue(inst.InstanceType), awssdk.StringValue(inst.PrivateIpAddress), awssdk.StringValue(inst.PublicIpAddress)))
		}
	case *ec2.DescribeVpcsOutput:
		fmt.Fprintln(w, "Id\tDefault\tState\tCidr")
		for _, vpc := range item.(*ec2.DescribeVpcsOutput).Vpcs {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s", *vpc.VpcId, printColorIf(*vpc.IsDefault), *vpc.State, *vpc.CidrBlock))
		}
	case *ec2.DescribeSubnetsOutput:
		fmt.Fprintln(w, "Id\tPublic VMs\tState\tCidr")
		for _, subnet := range item.(*ec2.DescribeSubnetsOutput).Subnets {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s", *subnet.SubnetId, printColorIf(*subnet.MapPublicIpOnLaunch, color.FgRed), *subnet.State, *subnet.CidrBlock))
		}
	default:
		fmt.Fprintln(w, item)
		return
	}
}

func printColorIf(cond bool, c ...color.Attribute) string {
	col := color.FgGreen
	if len(c) > 0 {
		col = c[0]
	}

	var fn func(string, ...interface{}) string
	if cond {
		fn = color.New(col).SprintfFunc()
	} else {
		fn = color.New().SprintfFunc()
	}

	return fn(fmt.Sprintf("%t", cond))
}
