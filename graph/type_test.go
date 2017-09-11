package graph

import (
	"encoding/json"
	"net"
	"reflect"
	"testing"
)

func TestFirewallRuleContainsIP(t *testing.T) {
	tcases := []struct {
		nets   []string
		ip     string
		result bool
	}{
		{[]string{}, "89.87.189.250", false},
		{[]string{"89.0.0.0/8"}, "89.87.189.250", true},
		{[]string{"89.0.0.0/16"}, "89.87.189.250", false},
		{[]string{"89.87.0.0/16"}, "89.87.189.250", true},
		{[]string{"89.0.0.0/0"}, "89.87.1", false},
	}

	for i, tcase := range tcases {
		rule := &FirewallRule{}
		for _, n := range tcase.nets {
			_, ipnet, _ := net.ParseCIDR(n)
			rule.IPRanges = append(rule.IPRanges, ipnet)
		}
		if rule.Contains(tcase.ip) != tcase.result {
			t.Fatalf("%d. case %s in %v expected %t", i+1, tcase.ip, tcase.nets, tcase.result)
		}
	}
}

func TestPortRangeContainsPort(t *testing.T) {
	tcases := []struct {
		prange PortRange
		port   int64
		result bool
	}{
		{PortRange{Any: true}, 2373, true},
		{PortRange{Any: false}, 2373, false},
		{PortRange{FromPort: 22}, 22, true},
		{PortRange{ToPort: 22}, 22, true},
		{PortRange{FromPort: 22, ToPort: 22}, 22, true},
		{PortRange{FromPort: 20, ToPort: 22}, 22, true},
		{PortRange{FromPort: 22, ToPort: 25}, 22, true},
		{PortRange{FromPort: 20, ToPort: 25}, 22, true},
		{PortRange{FromPort: 23, ToPort: 25}, 22, false},
		{PortRange{}, 22, false},
	}

	for i, tcase := range tcases {
		if tcase.prange.Contains(tcase.port) != tcase.result {
			t.Fatalf("%d. case %d in %v expected %t", i+1, tcase.port, tcase.prange, tcase.result)
		}
	}
}

func TestExtractPolicyDocument(t *testing.T) {
	tcases := []struct {
		document string
		expect   *Policy
	}{
		{
			document: `{
"Version": "2012-10-17",
"Statement": [
	{
		"Effect": "Allow",
		"Action": [
			"ec2:AttachVolume",
			"ec2:DetachVolume"
		],
		"Resource": [
			"arn:aws:ec2:<REGION>:<ACCOUNTNUMBER>:volume/*",
			"arn:aws:ec2:<REGION>:<ACCOUNTNUMBER>:instance/*"
		],
		"Condition": {
			"ArnEquals": {
				"ec2:SourceInstanceARN": "arn:aws:ec2:<REGION>:<ACCOUNTNUMBER>:instance/<INSTANCE-ID>"
			}
		}
	}
]
}`,
			expect: &Policy{
				Version: "2012-10-17",
				Statements: []*PolicyStatement{
					{
						Effect:    "Allow",
						Actions:   []string{"ec2:AttachVolume", "ec2:DetachVolume"},
						Resources: []string{"arn:aws:ec2:<REGION>:<ACCOUNTNUMBER>:volume/*", "arn:aws:ec2:<REGION>:<ACCOUNTNUMBER>:instance/*"},
						Condition: map[string]interface{}{
							"ArnEquals": map[string]interface{}{
								"ec2:SourceInstanceARN": "arn:aws:ec2:<REGION>:<ACCOUNTNUMBER>:instance/<INSTANCE-ID>",
							},
						},
					},
				},
			},
		},
		{
			document: `{
"Version": "2012-10-17",
"Statement": [
	{
		"Effect": "Allow",
		"Action": [
			"ec2:AttachVolume",
			"ec2:DetachVolume"
		],
		"Resource": "arn:aws:ec2:<REGION>:<ACCOUNTNUMBER>:instance/*",
		"Condition": {"StringEquals": {"ec2:ResourceTag/department": "dev"}}
	},
	{
		"Effect": "Allow",
		"Action": [
			"ec2:AttachVolume",
			"ec2:DetachVolume"
		],
		"Resource": "arn:aws:ec2:<REGION>:<ACCOUNTNUMBER>:volume/*",
		"Condition": {"StringEquals": {"ec2:ResourceTag/volume_user": "${aws:username}"}}
	}
]
}`,
			expect: &Policy{
				Version: "2012-10-17",
				Statements: []*PolicyStatement{
					{
						Effect:    "Allow",
						Actions:   []string{"ec2:AttachVolume", "ec2:DetachVolume"},
						Resources: []string{"arn:aws:ec2:<REGION>:<ACCOUNTNUMBER>:instance/*"},
						Condition: map[string]interface{}{
							"StringEquals": map[string]interface{}{
								"ec2:ResourceTag/department": "dev",
							},
						},
					},
					{
						Effect:    "Allow",
						Actions:   []string{"ec2:AttachVolume", "ec2:DetachVolume"},
						Resources: []string{"arn:aws:ec2:<REGION>:<ACCOUNTNUMBER>:volume/*"},
						Condition: map[string]interface{}{
							"StringEquals": map[string]interface{}{
								"ec2:ResourceTag/volume_user": "${aws:username}",
							},
						},
					},
				},
			},
		},
		{
			document: `{
"Version": "2012-10-17",
"Statement": [
	{
		"Effect": "Allow",
		"Principal": {
			"Service": [
				"elasticmapreduce.amazonaws.com",
				"datapipeline.amazonaws.com"
			]
		},
		"Action": "sts:AssumeRole"
	}
]
}`,
			expect: &Policy{
				Version: "2012-10-17",
				Statements: []*PolicyStatement{
					{
						Effect:    "Allow",
						Actions:   []string{"sts:AssumeRole"},
						Principal: &StatementPrincipal{Service: []string{"elasticmapreduce.amazonaws.com", "datapipeline.amazonaws.com"}},
					},
				},
			},
		},
		{
			document: `{
"Version": "2012-10-17",
"Statement": {
"Sid": "AccountBAccess1",
"Effect": "Allow",
"Principal": {"AWS": "111122223333"},
"Action": "s3:*",
"Resource": [
"arn:aws:s3:::mybucket",
"arn:aws:s3:::mybucket/*"
]
}
}`,
			expect: &Policy{
				Version: "2012-10-17",
				Statements: []*PolicyStatement{
					{
						ID:        "AccountBAccess1",
						Effect:    "Allow",
						Actions:   []string{"s3:*"},
						Principal: &StatementPrincipal{AWS: []string{"111122223333"}},
						Resources: []string{"arn:aws:s3:::mybucket", "arn:aws:s3:::mybucket/*"},
					},
				},
			},
		},
		{
			document: `{
"Version": "2012-10-17",
"Statement": {
"Effect": "Deny",
"Principal": "*",
"Action": "*",
"Resource": "*"
}
}`,
			expect: &Policy{
				Version: "2012-10-17",
				Statements: []*PolicyStatement{
					{
						Effect:    "Deny",
						Actions:   []string{"*"},
						Principal: &StatementPrincipal{AWS: []string{"*"}},
						Resources: []string{"*"},
					},
				},
			},
		},
	}
	for i, tcase := range tcases {
		var policy *Policy
		err := json.Unmarshal([]byte(tcase.document), &policy)
		if err != nil {
			t.Fatalf("%d: %s", i+1, err)
		}
		if got, want := policy, tcase.expect; !reflect.DeepEqual(got, want) {
			// fmt.Println("got:")
			// pretty.Print(got)
			// fmt.Println("\nwant:")
			// pretty.Print(want)
			t.Fatalf("%d: got \n%#v\nwant\n%#v\n", i+1, got, want)
		}
	}

}
