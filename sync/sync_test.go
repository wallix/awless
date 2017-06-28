package sync

import (
	"os"
	"testing"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/driver"

	"io/ioutil"

	"path/filepath"

	"github.com/wallix/awless/graph"
)

func TestSyncTripleFiles(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "awlessunittest_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	name1, region1 := "testservice", "paris"
	srv1 := &mockService{
		g:      graph.NewGraph(),
		name:   name1,
		region: region1,
	}

	name2, region2 := "testservice2", "bali"
	srv2 := &mockService{
		g:      graph.NewGraph(),
		name:   name2,
		region: region2,
	}

	os.Setenv("__AWLESS_HOME", tmpDir)

	if _, err := NewSyncer().Sync(srv1, srv2); err != nil {
		t.Fatal(err)
	}

	for _, srv := range []cloud.Service{srv1, srv2} {
		info, err := os.Stat(filepath.Join(tmpDir, "aws", "rdf", srv.Region(), srv.Name()+fileExt))
		if err != nil {
			t.Fatalf("cannot find expected file: %s", err)
		}
		if got, want := info.Name(), srv.Name()+fileExt; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}

	}
}

type mockService struct {
	name, region string
	g            *graph.Graph
}

func (s *mockService) Region() string                             { return s.region }
func (s *mockService) Name() string                               { return s.name }
func (s *mockService) Drivers() []driver.Driver                   { return nil }
func (s *mockService) ResourceTypes() []string                    { return []string{} }
func (s *mockService) FetchResources() (*graph.Graph, error)      { return s.g, nil }
func (s *mockService) IsSyncDisabled() bool                       { return false }
func (s *mockService) FetchByType(t string) (*graph.Graph, error) { return nil, nil }
