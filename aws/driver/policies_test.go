package awsdriver

import (
	"strings"
	"testing"
)

func TestResolvePolicy(t *testing.T) {
	p, err := LookupAWSPolicy("ec2", "readonly")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := p.Name, "AmazonEC2ReadOnlyAccess"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	p, err = LookupAWSPolicy("lambda", "full")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := p.Name, "AWSLambdaFullAccess"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	if _, err = LookupAWSPolicy("lambda", "fully"); err == nil {
		t.Fatal("expecting error got none")
	}
	if _, err = LookupAWSPolicy("lava", "full"); err == nil {
		t.Fatal("expecting error got none")
	}
}

func TestResolvePolicyErrorMessageWithSuggestion(t *testing.T) {
	_, err := LookupAWSPolicy("Administrator", "readonly")
	if err == nil {
		t.Fatal("expected error got none")
	}

	shouldContains := []string{
		"arn:aws:iam::aws:policy/job-function/DatabaseAdministrator",
		"arn:aws:iam::aws:policy/AdministratorAccess",
	}

	for _, e := range shouldContains {
		if msg := err.Error(); !strings.Contains(msg, e) {
			t.Errorf("expect\n%s\nto contain\n%s\n", msg, e)
		}
	}
}
