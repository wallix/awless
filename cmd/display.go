package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
)
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/fatih/color"
)

func display(item interface{}, err error, format ...string) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(format) > 0 {
		switch format[0] {
		case "raw":
			fmt.Println(item)
		default:
			lineDisplay(item)
		}
	} else {
		lineDisplay(item)
	}
}

var simpleDay = "Mon, Jan 2, 2006"

func lineDisplay(item interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 25, 1, 1, ' ', 0)

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
		fmt.Fprintln(w, "Id\tType\tState\tPriv IP\tPub IP\tLaunched")
		for _, reserv := range item.(*ec2.DescribeInstancesOutput).Reservations {
			for _, inst := range reserv.Instances {
				fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", aws.StringValue(inst.InstanceId), aws.StringValue(inst.State.Name), aws.StringValue(inst.InstanceType), aws.StringValue(inst.PrivateIpAddress), aws.StringValue(inst.PublicIpAddress), (*inst.LaunchTime).Format(simpleDay)))
			}
		}
	case *ec2.Reservation:
		fmt.Fprintln(w, "Id\tType\tState\tPriv IP\tPub IP\tLaunched")
		for _, inst := range item.(*ec2.Reservation).Instances {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", aws.StringValue(inst.InstanceId), aws.StringValue(inst.State.Name), aws.StringValue(inst.InstanceType), aws.StringValue(inst.PrivateIpAddress), aws.StringValue(inst.PublicIpAddress), (*inst.LaunchTime).Format(simpleDay)))
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
		fmt.Println(item)
		return
	}

	w.Flush()
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
