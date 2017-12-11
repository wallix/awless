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

package inspectors

import (
	"fmt"
	"io"
	"strings"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
)

type OpenBuckets struct {
	openToAny     []string
	openToAnyAuth []string
}

func (*OpenBuckets) Name() string {
	return "open_buckets"
}

func (a *OpenBuckets) Inspect(g cloud.GraphAPI) error {
	buckets, err := g.Find(cloud.NewQuery(cloud.Bucket))
	if err != nil {
		return err
	}

	openToAuthUsers := make(map[string]bool)
	openToUsers := make(map[string]bool)

	for _, buck := range buckets {
		grants, ok := buck.Properties()["Grants"].([]*graph.Grant)
		if ok {
			for _, g := range grants {
				if strings.Contains(g.Grantee.GranteeID, "AllUsers") {
					openToUsers[fmt.Sprint(buck.Properties()["ID"])] = true
				}
				if strings.Contains(g.Grantee.GranteeID, "AuthenticatedUsers") {
					openToAuthUsers[fmt.Sprint(buck.Properties()["ID"])] = true
				}
			}
		}
	}

	for k := range openToUsers {
		a.openToAny = append(a.openToAny, k)
	}

	for k := range openToAuthUsers {
		a.openToAnyAuth = append(a.openToAnyAuth, k)
	}

	return nil
}

func (a *OpenBuckets) Print(w io.Writer) {
	if len(a.openToAny) > 0 {
		fmt.Fprintf(w, "Buckets open to anybody: %s\n", strings.Join(a.openToAny, ", "))
	}
	if len(a.openToAnyAuth) > 0 {
		fmt.Fprintf(w, "Buckets open to anyone with an AWS account: %s\n", strings.Join(a.openToAnyAuth, ", "))
	}

	if len(a.openToAnyAuth) == 0 && len(a.openToAny) == 0 {
		fmt.Fprintln(w, "none found")
	}
}
