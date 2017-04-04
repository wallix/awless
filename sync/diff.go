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

package sync

import (
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/sync/repo"
)

// Diff represents the deleted/inserted RDF triples of a revision
type Diff struct {
	From       *repo.Rev
	To         *repo.Rev
	InfraDiff  *graph.Diff
	AccessDiff *graph.Diff
}

func BuildDiff(from, to *repo.Rev, root string) (*Diff, error) {
	infraDiff, err := graph.DefaultDiffer.Run(root, from.Infra, to.Infra)
	if err != nil {
		return nil, err
	}

	accessDiff, err := graph.DefaultDiffer.Run(root, from.Access, to.Access)
	if err != nil {
		return nil, err
	}

	res := &Diff{
		From:       from,
		To:         to,
		InfraDiff:  infraDiff,
		AccessDiff: accessDiff,
	}

	return res, nil
}
