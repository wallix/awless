package awsservices

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/wallix/awless/logger"
)

func TestFileCacheProvider(t *testing.T) {
	name, err := ioutil.TempDir(".", "cache")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(name)
	os.Setenv("__AWLESS_CACHE", name)

	mock := &mockCredWithExpirationProvider{value: credentials.Value{SecretAccessKey: "my valid secret string", ProviderName: stscreds.ProviderName}}
	creds := credentials.NewCredentials(mock)
	stscreds.DefaultDuration = 30 * time.Millisecond //Force cached credential expiration after 20 millisecond

	provider := fileCacheProvider{creds: creds, profile: "default", log: logger.DiscardLogger}
	retrievedCredential, err := provider.Retrieve()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := retrievedCredential.SecretAccessKey, "my valid secret string"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := mock.accessCount, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	retrievedCredential, err = provider.Retrieve()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := retrievedCredential.SecretAccessKey, "my valid secret string"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := mock.accessCount, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	time.Sleep(stscreds.DefaultDuration)

	retrievedCredential, err = provider.Retrieve()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := retrievedCredential.SecretAccessKey, "my valid secret string"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := mock.accessCount, 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

type mockCredWithExpirationProvider struct {
	accessCount int
	value       credentials.Value
}

func (m *mockCredWithExpirationProvider) Retrieve() (credentials.Value, error) {
	m.accessCount++
	return m.value, nil
}

func (m *mockCredWithExpirationProvider) IsExpired() bool {
	return false
}
