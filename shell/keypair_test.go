package shell

import (
	"strings"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	size := 1024
	pub, private, err := GenerateSSHKeyPair(size)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := strings.HasPrefix(string(pub), "ssh-rsa"), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := strings.HasPrefix(string(private), "-----BEGIN RSA PRIVATE KEY-----"), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := strings.HasSuffix(string(private), "-----END RSA PRIVATE KEY-----\n"), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
}
