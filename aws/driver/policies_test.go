package awsdriver

import "testing"

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
