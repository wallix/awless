/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"path/filepath"
	"testing"

	p "github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestInstanceCredentialsFromName(t *testing.T) {
	inst_1 := resourcetest.Instance("inst_1").Prop(p.KeyPair, "my-key-name").Prop(p.PublicIP, "1.2.3.4").Build()
	inst_2 := resourcetest.Instance("inst_2").Prop(p.PublicIP, "2.3.4.5").Build()
	inst_3 := resourcetest.Instance("inst_3").Build()
	inst_12 := resourcetest.Instance("inst_12").Build()
	g := graph.NewGraph()
	g.AddResource(inst_1, inst_2)

	keypath, IP, err := instanceCredentialsFromGraph(g, inst_1, "")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := IP, "1.2.3.4"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := filepath.Base(keypath), "my-key-name"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	keypath, IP, err = instanceCredentialsFromGraph(g, inst_1, "/path/toward/myinst.key")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := IP, "1.2.3.4"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := keypath, "/path/toward/myinst.key"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	keypath, IP, err = instanceCredentialsFromGraph(g, inst_2, "/path/toward/inst2.key")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := IP, "2.3.4.5"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := keypath, "/path/toward/inst2.key"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	if _, _, err = instanceCredentialsFromGraph(g, inst_12, ""); err == nil {
		t.Fatal("expected error got none")
	}
	if _, _, err := instanceCredentialsFromGraph(g, inst_3, ""); err == nil {
		t.Fatal("expected error got none")
	}
	if _, _, err := instanceCredentialsFromGraph(g, inst_2, ""); err == nil {
		t.Fatal("expected error got none")
	}
}
