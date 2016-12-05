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
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 1, ' ', 0)

	switch item.(type) {
	case *iam.ListUsersOutput:
		for _, user := range item.(*iam.ListUsersOutput).Users {
			fmt.Fprintln(w, fmt.Sprintf("name: %s\tid: %s\tcreated: %s\t", *user.UserName, *user.UserId, (*user.CreateDate).Format(simpleDay)))
		}
	case *iam.ListGroupsOutput:
		for _, group := range item.(*iam.ListGroupsOutput).Groups {
			fmt.Fprintln(w, fmt.Sprintf("name:%s\tid:%s\tcreated: %s\t", *group.GroupName, *group.GroupId, (*group.CreateDate).Format(simpleDay)))
		}
	case *iam.ListRolesOutput:
		for _, role := range item.(*iam.ListRolesOutput).Roles {
			fmt.Fprintln(w, fmt.Sprintf("name: %s\tid: %s\tcreated: %s\t", *role.RoleName, *role.RoleId, (*role.CreateDate).Format(simpleDay)))
		}
	case *iam.ListPoliciesOutput:
		for _, policy := range item.(*iam.ListPoliciesOutput).Policies {
			fmt.Fprintln(w, fmt.Sprintf("name: %s\tid: %s\tcreated: %s\t", *policy.PolicyName, *policy.PolicyId, (*policy.CreateDate).Format(simpleDay)))
		}
	case *ec2.DescribeInstancesOutput:
		for _, reserv := range item.(*ec2.DescribeInstancesOutput).Reservations {
			for _, inst := range reserv.Instances {
				fmt.Fprintln(w, fmt.Sprintf("id: %s\ttype: %s\tstate: %s\tpriv-ip: %s\tpub-ip: %s\tlaunched: %s\t", aws.StringValue(inst.InstanceId), aws.StringValue(inst.State.Name), aws.StringValue(inst.InstanceType), aws.StringValue(inst.PrivateIpAddress), aws.StringValue(inst.PublicIpAddress), (*inst.LaunchTime).Format(simpleDay)))
			}
		}
	case *ec2.Reservation:
		for _, inst := range item.(*ec2.Reservation).Instances {
			fmt.Fprintln(w, fmt.Sprintf("id: %s\ttype: %s\tstate: %s\tpriv-ip: %s\tpub-ip: %s\tlaunched: %s\t", aws.StringValue(inst.InstanceId), aws.StringValue(inst.State.Name), aws.StringValue(inst.InstanceType), aws.StringValue(inst.PrivateIpAddress), aws.StringValue(inst.PublicIpAddress), (*inst.LaunchTime).Format(simpleDay)))
		}
	case *ec2.DescribeVpcsOutput:
		for _, vpc := range item.(*ec2.DescribeVpcsOutput).Vpcs {
			fmt.Fprintln(w, fmt.Sprintf("id: %s\tdefault: %s\tstate: %s\tcidr: %s\t", *vpc.VpcId, printColorIf(*vpc.IsDefault), *vpc.State, *vpc.CidrBlock))
		}
	case *ec2.DescribeSubnetsOutput:
		for _, subnet := range item.(*ec2.DescribeSubnetsOutput).Subnets {
			fmt.Fprintln(w, fmt.Sprintf("id: %s\tpublic-vms: %s\tstate: %s\tcidr: %s\t", *subnet.SubnetId, printColorIf(*subnet.MapPublicIpOnLaunch, color.FgRed), *subnet.State, *subnet.CidrBlock))
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
