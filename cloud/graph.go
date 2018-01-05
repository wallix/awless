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
)

type GraphAPI interface {
	Find(Query) ([]Resource, error)
	FindWithProperties(map[string]interface{}) ([]Resource, error)
	FilterGraph(Query) (GraphAPI, error)
	FindOne(Query) (Resource, error)
	MarshalTo(w io.Writer) error
	ResourceRelations(r Resource, relation string, recursive bool) ([]Resource, error)
	VisitRelations(Resource, string, bool, func(Resource, int) error) error
	ResourceSiblings(Resource) ([]Resource, error)
	Merge(GraphAPI) error
}

type Resource interface {
	Type() string
	Id() string
	String() string
	Format(string) string
	Properties() map[string]interface{}
	Property(string) (interface{}, bool)
	Meta(string) (interface{}, bool)
	Same(Resource) bool
}

type Resources []Resource

func (res Resources) Map(f func(Resource) string) (out []string) {
	for _, r := range res {
		out = append(out, f(r))
	}
	return
}
