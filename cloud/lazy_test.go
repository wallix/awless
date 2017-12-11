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
package cloud

import (
	"io"
	"testing"
)

func TestLazyLoadingGraph(t *testing.T) {
	var nbCalls int
	loadingFunc := func() GraphAPI {
		nbCalls++
		return &StubGraph{}
	}

	lazy := &LazyGraph{LoadingFunc: loadingFunc}
	lazy.FindOne(NewQuery(""))
	if got, want := nbCalls, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	lazy.FindOne(NewQuery(""))
	if got, want := nbCalls, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

type StubGraph struct {
}

func (g *StubGraph) Find(Query) ([]Resource, error) {
	return nil, nil
}

func (g *StubGraph) FindWithProperties(props map[string]interface{}) ([]Resource, error) {
	return nil, nil
}

func (g *StubGraph) FilterGraph(Query) (GraphAPI, error) {
	return nil, nil
}

func (g *StubGraph) FindOne(Query) (Resource, error) {
	return nil, nil
}

func (g *StubGraph) MarshalTo(io.Writer) error {
	return nil
}

func (g *StubGraph) ResourceRelations(Resource, string, bool) ([]Resource, error) {
	return nil, nil
}

func (g *StubGraph) VisitRelations(Resource, string, bool, func(Resource, int) error) error {
	return nil
}

func (g *StubGraph) ResourceSiblings(Resource) ([]Resource, error) {
	return nil, nil
}

func (g *StubGraph) Merge(GraphAPI) error {
	return nil
}
