package awsat

import (
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
)

func TestPolicy(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template(
			"create policy name=AwlessInfraReadonlyPolicy effect=Allow action=ec2:Describe*,autoscaling:Describe*,elasticloadbalancing:Describe* "+
				"resource=\"arn:aws:iam::0123456789:mfa/${aws:username}\",\"arn:aws:iam::0123456789:user/${aws:username}\" "+
				"conditions=\"aws:MultiFactorAuthPresent==true\", \"aws:TokenIssueTime!=Null\" description=\"Readonly access to infra resources\"").
			Mock(&iamMock{
				CreatePolicyFunc: func(input *iam.CreatePolicyInput) (*iam.CreatePolicyOutput, error) {
					return &iam.CreatePolicyOutput{Policy: &iam.Policy{Arn: String("new-policy-arn")}}, nil
				},
			}).ExpectInput("CreatePolicy", &iam.CreatePolicyInput{
			PolicyName:  String("AwlessInfraReadonlyPolicy"),
			Description: String("Readonly access to infra resources"),
			PolicyDocument: String(`{
 "Version": "2012-10-17",
 "Statement": [
  {
   "Effect": "Allow",
   "Action": [
    "ec2:Describe*",
    "autoscaling:Describe*",
    "elasticloadbalancing:Describe*"
   ],
   "Resource": [
    "arn:aws:iam::0123456789:mfa/${aws:username}",
    "arn:aws:iam::0123456789:user/${aws:username}"
   ],
   "Condition": {
    "Bool": {
     "aws:MultiFactorAuthPresent": "true"
    },
    "Null": {
     "aws:TokenIssueTime": "false"
    }
   }
  }
 ]
}`)}).ExpectCommandResult("new-policy-arn").ExpectCalls("CreatePolicy").Run(t)
	})

	t.Run("update", func(t *testing.T) {
		Template(
			"update policy arn=arn:my:arn:of:policy:to:update effect=Deny action=ec2:AttachVolume,DescribeVolumeAttribute "+
				"resource=\"arn:aws:ec2:eu-west-1:0123456789:volume/*\" "+
				"conditions=\"aws:MultiFactorAuthPresent==true\"").
			Mock(&iamMock{
				CreatePolicyVersionFunc: func(input *iam.CreatePolicyVersionInput) (*iam.CreatePolicyVersionOutput, error) {
					return nil, nil
				},
				ListPolicyVersionsFunc: func(input *iam.ListPolicyVersionsInput) (*iam.ListPolicyVersionsOutput, error) {
					return &iam.ListPolicyVersionsOutput{Versions: []*iam.PolicyVersion{{VersionId: String("v2"), IsDefaultVersion: Bool(true)}}}, nil
				},
				GetPolicyVersionFunc: func(input *iam.GetPolicyVersionInput) (*iam.GetPolicyVersionOutput, error) {
					return &iam.GetPolicyVersionOutput{PolicyVersion: &iam.PolicyVersion{
						Document: String(url.QueryEscape(`{
 "Version": "2012-10-17",
 "Statement": [
  {
   "Effect": "Allow",
   "Action": [
    "ec2:AttachVolume",
    "ec2:DetachVolume"
   ],
   "Resource": "arn:aws:ec2:eu-west-1:0123456789:instance/*",
   "Condition": {
    "StringEquals": {
     "ec2:ResourceTag/department": "dev"
    },
    "Null": {
     "aws:TokenIssueTime": "false"
    }
   }
  },
  {
   "Effect": "Allow",
   "Action": [
    "ec2:AttachVolume",
    "ec2:DetachVolume"
   ],
   "Resource": "arn:aws:ec2:eu-west-1:0123456789:volume/*",
   "Condition": {
    "StringEquals": {
     "ec2:ResourceTag/volume_user": "${aws:username}"}
    }
  }
 ]
}`)),
						IsDefaultVersion: Bool(true),
						VersionId:        String("v2"),
					}}, nil
				},
			}).ExpectInput("ListPolicyVersions", &iam.ListPolicyVersionsInput{
			PolicyArn: String("arn:my:arn:of:policy:to:update"),
		}).ExpectInput("GetPolicyVersion", &iam.GetPolicyVersionInput{
			PolicyArn: String("arn:my:arn:of:policy:to:update"),
			VersionId: String("v2"),
		}).ExpectInput("CreatePolicyVersion", &iam.CreatePolicyVersionInput{
			PolicyArn:    String("arn:my:arn:of:policy:to:update"),
			SetAsDefault: Bool(true),
			PolicyDocument: String(`{
 "Version": "2012-10-17",
 "Statement": [
  {
   "Effect": "Allow",
   "Action": [
    "ec2:AttachVolume",
    "ec2:DetachVolume"
   ],
   "Resource": "arn:aws:ec2:eu-west-1:0123456789:instance/*",
   "Condition": {
    "StringEquals": {
     "ec2:ResourceTag/department": "dev"
    },
    "Null": {
     "aws:TokenIssueTime": "false"
    }
   }
  },
  {
   "Effect": "Allow",
   "Action": [
    "ec2:AttachVolume",
    "ec2:DetachVolume"
   ],
   "Resource": "arn:aws:ec2:eu-west-1:0123456789:volume/*",
   "Condition": {
    "StringEquals": {
     "ec2:ResourceTag/volume_user": "${aws:username}"
    }
   }
  },
  {
   "Effect": "Deny",
   "Action": [
    "ec2:AttachVolume",
    "DescribeVolumeAttribute"
   ],
   "Resource": [
    "arn:aws:ec2:eu-west-1:0123456789:volume/*"
   ],
   "Condition": {
    "Bool": {
     "aws:MultiFactorAuthPresent": "true"
    }
   }
  }
 ]
}`)}).ExpectCalls("ListPolicyVersions", "GetPolicyVersion", "CreatePolicyVersion").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template(
			"delete policy arn=arn:my:arn:of:policy:to:delete all-versions=true").
			Mock(&iamMock{
				DeletePolicyFunc: func(input *iam.DeletePolicyInput) (*iam.DeletePolicyOutput, error) {
					return nil, nil
				},
				DeletePolicyVersionFunc: func(input *iam.DeletePolicyVersionInput) (*iam.DeletePolicyVersionOutput, error) {
					return nil, nil
				},
				ListPolicyVersionsFunc: func(input *iam.ListPolicyVersionsInput) (*iam.ListPolicyVersionsOutput, error) {
					return &iam.ListPolicyVersionsOutput{Versions: []*iam.PolicyVersion{
						{VersionId: String("v1"), IsDefaultVersion: Bool(false)},
						{VersionId: String("v2"), IsDefaultVersion: Bool(true)},
					}}, nil
				},
			}).ExpectInput("DeletePolicy", &iam.DeletePolicyInput{
			PolicyArn: String("arn:my:arn:of:policy:to:delete"),
		}).ExpectInput("DeletePolicyVersion", &iam.DeletePolicyVersionInput{
			PolicyArn: String("arn:my:arn:of:policy:to:delete"),
			VersionId: String("v1"),
		}).ExpectInput("ListPolicyVersions", &iam.ListPolicyVersionsInput{PolicyArn: String("arn:my:arn:of:policy:to:delete")}).
			ExpectCalls("DeletePolicy", "DeletePolicyVersion", "ListPolicyVersions").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template(
			"attach policy group=administrators access=readonly service=ec2").
			Mock(&iamMock{
				AttachGroupPolicyFunc: func(input *iam.AttachGroupPolicyInput) (*iam.AttachGroupPolicyOutput, error) {
					return nil, nil
				},
			}).ExpectInput("AttachGroupPolicy", &iam.AttachGroupPolicyInput{
			GroupName: String("administrators"),
			PolicyArn: String("arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess"),
		}).ExpectCalls("AttachGroupPolicy").Run(t)

		Template(
			"attach policy user=toto arn=arn:for:my:policy").
			Mock(&iamMock{
				AttachUserPolicyFunc: func(input *iam.AttachUserPolicyInput) (*iam.AttachUserPolicyOutput, error) {
					return nil, nil
				},
			}).ExpectInput("AttachUserPolicy", &iam.AttachUserPolicyInput{
			UserName:  String("toto"),
			PolicyArn: String("arn:for:my:policy"),
		}).ExpectCalls("AttachUserPolicy").Run(t)

		Template(
			"attach policy role=my-role access=full service=ec2").
			Mock(&iamMock{
				AttachRolePolicyFunc: func(input *iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error) {
					return nil, nil
				},
			}).ExpectInput("AttachRolePolicy", &iam.AttachRolePolicyInput{
			RoleName:  String("my-role"),
			PolicyArn: String("arn:aws:iam::aws:policy/AmazonEC2FullAccess"),
		}).ExpectCalls("AttachRolePolicy").Run(t)
	})

	t.Run("attach services meta", func(t *testing.T) {
		Template(
			"attach policy group=administrators access=readonly services=ec2,rds").
			Mock(&iamMock{
				AttachGroupPolicyFunc: func(input *iam.AttachGroupPolicyInput) (*iam.AttachGroupPolicyOutput, error) {
					return nil, nil
				},
			}).ExpectInput("AttachGroupPolicy", &iam.AttachGroupPolicyInput{
			GroupName: String("administrators"),
			PolicyArn: String("arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess"),
		}).ExpectInput("AttachGroupPolicy", &iam.AttachGroupPolicyInput{
			GroupName: String("administrators"),
			PolicyArn: String("arn:aws:iam::aws:policy/AmazonRDSReadOnlyAccess"),
		}).ExpectCalls("AttachGroupPolicy", "AttachGroupPolicy").Run(t)

		Template(
			"attach policy user=fx access=full services=s3,autoscaling").
			Mock(&iamMock{
				AttachUserPolicyFunc: func(input *iam.AttachUserPolicyInput) (*iam.AttachUserPolicyOutput, error) {
					return nil, nil
				},
			}).ExpectInput("AttachUserPolicy", &iam.AttachUserPolicyInput{
			UserName:  String("fx"),
			PolicyArn: String("arn:aws:iam::aws:policy/AmazonS3FullAccess"),
		}).ExpectInput("AttachUserPolicy", &iam.AttachUserPolicyInput{
			UserName:  String("fx"),
			PolicyArn: String("arn:aws:iam::aws:policy/AutoScalingFullAccess"),
		}).ExpectCalls("AttachUserPolicy", "AttachUserPolicy").Run(t)
	})

	t.Run("detach", func(t *testing.T) {
		Template(
			"detach policy group=administrators access=readonly service=ec2").
			Mock(&iamMock{
				DetachGroupPolicyFunc: func(input *iam.DetachGroupPolicyInput) (*iam.DetachGroupPolicyOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DetachGroupPolicy", &iam.DetachGroupPolicyInput{
			GroupName: String("administrators"),
			PolicyArn: String("arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess"),
		}).ExpectCalls("DetachGroupPolicy").Run(t)

		Template(
			"detach policy user=toto arn=arn:for:my:policy").
			Mock(&iamMock{
				DetachUserPolicyFunc: func(input *iam.DetachUserPolicyInput) (*iam.DetachUserPolicyOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DetachUserPolicy", &iam.DetachUserPolicyInput{
			UserName:  String("toto"),
			PolicyArn: String("arn:for:my:policy"),
		}).ExpectCalls("DetachUserPolicy").Run(t)

		Template(
			"detach policy role=my-role access=full service=ec2").
			Mock(&iamMock{
				DetachRolePolicyFunc: func(input *iam.DetachRolePolicyInput) (*iam.DetachRolePolicyOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DetachRolePolicy", &iam.DetachRolePolicyInput{
			RoleName:  String("my-role"),
			PolicyArn: String("arn:aws:iam::aws:policy/AmazonEC2FullAccess"),
		}).ExpectCalls("DetachRolePolicy").Run(t)
	})
}
