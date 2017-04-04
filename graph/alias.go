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

package graph

import (
	cloudrdf "github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

type Alias string

func (a Alias) ResolveToId(g *Graph, resT string) (string, bool) {
	snap := g.store.Snapshot()
	triples := snap.WithPredObj(cloudrdf.Name, tstore.StringLiteral(string(a)))

	for _, t := range triples {
		id := t.Subject()
		rt, err := resolveResourceType(snap, id)
		if err != nil {
			return id, false
		}
		if resT == rt {
			return id, true
		}
	}

	return "", false
}
