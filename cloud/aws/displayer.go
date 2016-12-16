package aws

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var simpleDay = "Mon, Jan 2, 2006"

func TabularDisplay(item interface{}, table *tablewriter.Table) {
	switch item.(type) {
	case *iam.ListUsersOutput:
		table.SetHeader([]string{"Name", "Id", "Created"})
		for _, user := range item.(*iam.ListUsersOutput).Users {
			table.Append([]string{*user.UserName, *user.UserId, (*user.CreateDate).Format(simpleDay)})
		}
	case *iam.ListGroupsOutput:
		table.SetHeader([]string{"Name", "Id", "Created"})
		for _, group := range item.(*iam.ListGroupsOutput).Groups {
			table.Append([]string{*group.GroupName, *group.GroupId, (*group.CreateDate).Format(simpleDay)})
		}
	case *iam.ListRolesOutput:
		table.SetHeader([]string{"Name", "Id", "Created"})
		for _, role := range item.(*iam.ListRolesOutput).Roles {
			table.Append([]string{*role.RoleName, *role.RoleId, (*role.CreateDate).Format(simpleDay)})
		}
	case *iam.ListPoliciesOutput:
		table.SetHeader([]string{"Name", "Id", "Created"})
		for _, policy := range item.(*iam.ListPoliciesOutput).Policies {
			table.Append([]string{*policy.PolicyName, *policy.PolicyId, (*policy.CreateDate).Format(simpleDay)})
		}
	case *ec2.DescribeInstancesOutput:
		table.SetHeader([]string{"Id", "Name", "State", "Type", "Priv IP", "Pub IP"})
		for _, reserv := range item.(*ec2.DescribeInstancesOutput).Reservations {
			for _, inst := range reserv.Instances {
				var name string
				for _, t := range inst.Tags {
					if *t.Key == "Name" {
						name = *t.Value
					}
				}
				table.Append([]string{awssdk.StringValue(inst.InstanceId), name, awssdk.StringValue(inst.State.Name), awssdk.StringValue(inst.InstanceType), awssdk.StringValue(inst.PrivateIpAddress), awssdk.StringValue(inst.PublicIpAddress)})
			}
		}
	case *ec2.Reservation:
		table.SetHeader([]string{"Id", "Type", "State", "Priv IP", "Pub IP"})
		for _, inst := range item.(*ec2.Reservation).Instances {
			table.Append([]string{awssdk.StringValue(inst.InstanceId), awssdk.StringValue(inst.State.Name), awssdk.StringValue(inst.InstanceType), awssdk.StringValue(inst.PrivateIpAddress), awssdk.StringValue(inst.PublicIpAddress)})
		}
	case *ec2.DescribeVpcsOutput:
		table.SetHeader([]string{"Id", "Default", "State", "Cidr"})
		for _, vpc := range item.(*ec2.DescribeVpcsOutput).Vpcs {
			table.Append([]string{*vpc.VpcId, printColorIf(*vpc.IsDefault), *vpc.State, *vpc.CidrBlock})
		}
	case *ec2.DescribeSubnetsOutput:
		table.SetHeader([]string{"Id", "Public VMs", "State", "Cidr"})
		for _, subnet := range item.(*ec2.DescribeSubnetsOutput).Subnets {
			table.Append([]string{*subnet.SubnetId, printColorIf(*subnet.MapPublicIpOnLaunch, color.FgRed), *subnet.State, *subnet.CidrBlock})
		}
	default:
		fmt.Println(item)
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
