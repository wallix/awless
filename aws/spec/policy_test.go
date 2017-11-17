package awsspec

import (
	"reflect"
	"strings"
	"testing"
)

func TestBuildPolicyConditions(t *testing.T) {
	tcases := []struct {
		input  string
		output *policyCondition
	}{
		{"aws:UserAgent==Example Corp Java Client", &policyCondition{Type: "StringEquals", Key: "aws:UserAgent", Value: "Example Corp Java Client"}},
		{"aws:UserAgent!=Example Corp Java Client", &policyCondition{Type: "StringNotEquals", Key: "aws:UserAgent", Value: "Example Corp Java Client"}},
		{"s3:prefix=~'home/'", &policyCondition{Type: "StringLike", Key: "s3:prefix", Value: "home/"}},
		{"s3:prefix!~\"\"", &policyCondition{Type: "StringNotLike", Key: "s3:prefix", Value: ""}},
		{"s3:prefix!~", &policyCondition{Type: "StringNotLike", Key: "s3:prefix", Value: ""}},
		{"s3:max-keys==10", &policyCondition{Type: "NumericEquals", Key: "s3:max-keys", Value: "10"}},
		{"s3:max-keys!=10", &policyCondition{Type: "NumericNotEquals", Key: "s3:max-keys", Value: "10"}},
		{"s3:max-keys<10", &policyCondition{Type: "NumericLessThan", Key: "s3:max-keys", Value: "10"}},
		{"s3:max-keys<=10", &policyCondition{Type: "NumericLessThanEquals", Key: "s3:max-keys", Value: "10"}},
		{"s3:max-keys>10", &policyCondition{Type: "NumericGreaterThan", Key: "s3:max-keys", Value: "10"}},
		{"s3:max-keys>=10", &policyCondition{Type: "NumericGreaterThanEquals", Key: "s3:max-keys", Value: "10"}},
		{"aws:CurrentTime==2013-06-30T00:00:00Z", &policyCondition{Type: "DateEquals", Key: "aws:CurrentTime", Value: "2013-06-30T00:00:00Z"}},
		{"aws:CurrentTime!=2013-06-30T00:00:00Z", &policyCondition{Type: "DateNotEquals", Key: "aws:CurrentTime", Value: "2013-06-30T00:00:00Z"}},
		{"aws:CurrentTime<2013-06-30T00:00:00Z", &policyCondition{Type: "DateLessThan", Key: "aws:CurrentTime", Value: "2013-06-30T00:00:00Z"}},
		{"aws:CurrentTime<=2013-06-30T00:00:00Z", &policyCondition{Type: "DateLessThanEquals", Key: "aws:CurrentTime", Value: "2013-06-30T00:00:00Z"}},
		{"aws:CurrentTime>2013-06-30T00:00:00Z", &policyCondition{Type: "DateGreaterThan", Key: "aws:CurrentTime", Value: "2013-06-30T00:00:00Z"}},
		{"aws:CurrentTime>=2013-06-30T00:00:00Z", &policyCondition{Type: "DateGreaterThanEquals", Key: "aws:CurrentTime", Value: "2013-06-30T00:00:00Z"}},
		{"aws:SecureTransport==true", &policyCondition{Type: "Bool", Key: "aws:SecureTransport", Value: "true"}},
		{"aws:SecureTransport!=true", &policyCondition{Type: "Bool", Key: "aws:SecureTransport", Value: "false"}},
		{"aws:binarykey==QmluYXJ5VmFsdWVJbkJhc2U2NA==", &policyCondition{Type: "BinaryEquals", Key: "aws:binarykey", Value: "QmluYXJ5VmFsdWVJbkJhc2U2NA=="}},
		{"aws:SourceIp==203.0.113.0/24", &policyCondition{Type: "IpAddress", Key: "aws:SourceIp", Value: "203.0.113.0/24"}},
		{"aws:SourceIp!=203.0.113.0", &policyCondition{Type: "NotIpAddress", Key: "aws:SourceIp", Value: "203.0.113.0"}},
		{"aws:SourceIp==2001:DB8:1234:5678::/64", &policyCondition{Type: "IpAddress", Key: "aws:SourceIp", Value: "2001:DB8:1234:5678::/64"}},
		{"aws:SourceArn==arn:aws:sns:REGION:123456789012:TOPIC-ID", &policyCondition{Type: "ArnEquals", Key: "aws:SourceArn", Value: "arn:aws:sns:REGION:123456789012:TOPIC-ID"}},
		{"aws:SourceArn!=arn:aws:sns:*:*:TOPIC-ID", &policyCondition{Type: "ArnNotEquals", Key: "aws:SourceArn", Value: "arn:aws:sns:*:*:TOPIC-ID"}},
		{"aws:SourceArn=~arn:aws:sns:*:*:TOPIC-ID", &policyCondition{Type: "ArnLike", Key: "aws:SourceArn", Value: "arn:aws:sns:*:*:TOPIC-ID"}},
		{"aws:SourceArn!~arn:aws:sns:*:*:TOPIC-ID", &policyCondition{Type: "ArnNotLike", Key: "aws:SourceArn", Value: "arn:aws:sns:*:*:TOPIC-ID"}},
		{"aws:TokenIssueTime==Null", &policyCondition{Type: "Null", Key: "aws:TokenIssueTime", Value: "true"}},
		{"aws:TokenIssueTime!=Null", &policyCondition{Type: "Null", Key: "aws:TokenIssueTime", Value: "false"}},
	}
	for i, tcase := range tcases {
		cond, err := parseCondition(tcase.input)
		if err != nil {
			t.Fatalf("%d: %s", i+1, err)
		}
		if got, want := cond, tcase.output; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %#v, want %#v", i+1, got, want)
		}
	}
}

func TestResolvePolicy(t *testing.T) {
	p, err := lookupAWSPolicy("ec2", "readonly")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := p.Name, "AmazonEC2ReadOnlyAccess"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	p, err = lookupAWSPolicy("lambda", "full")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := p.Name, "AWSLambdaFullAccess"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	if _, err = lookupAWSPolicy("lambda", "fully"); err == nil {
		t.Fatal("expecting error got none")
	}
	if _, err = lookupAWSPolicy("lava", "full"); err == nil {
		t.Fatal("expecting error got none")
	}
}

func TestResolvePolicyErrorMessageWithSuggestion(t *testing.T) {
	_, err := lookupAWSPolicy("Administrator", "readonly")
	if err == nil {
		t.Fatal("expected error got none")
	}

	shouldContains := []string{
		"arn:aws:iam::aws:policy/job-function/DatabaseAdministrator",
		"arn:aws:iam::aws:policy/AdministratorAccess",
	}

	for _, e := range shouldContains {
		if msg := err.Error(); !strings.Contains(msg, e) {
			t.Errorf("expect '%s' to contain '%s'", msg, e)
		}
	}
}
