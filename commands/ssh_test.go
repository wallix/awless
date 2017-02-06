package commands

import (
	"path/filepath"
	"testing"

	"github.com/wallix/awless/graph"
)

func TestInstanceCredentialsFromName(t *testing.T) {
	g, err := graph.NewGraphFromFile(filepath.Join("testdata", "infra.rdf"))
	if err != nil {
		t.Fatal(err)
	}

	cred, err := instanceCredentialsFromGraph(g, "inst_1")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := cred.IP, "1.2.3.4"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := cred.KeyName, "my-key-name"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := cred.User, ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	_, err = instanceCredentialsFromGraph(g, "inst_12")
	if err == nil {
		t.Fatal("expected error got none")
	}
	if _, err := instanceCredentialsFromGraph(g, "inst_3"); err == nil {
		t.Fatal("expected error got none")
	}
	if _, err := instanceCredentialsFromGraph(g, "inst_2"); err == nil {
		t.Fatal("expected error got none")
	}
}
