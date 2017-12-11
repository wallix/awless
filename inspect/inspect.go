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

package inspect

import (
	"io"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/inspect/inspectors"
)

var InspectorsRegister map[string]Inspector

func init() {
	all := []Inspector{
		&inspectors.Pricer{}, &inspectors.BucketSizer{},
		&inspectors.PortScanner{}, &inspectors.OpenBuckets{},
	}

	InspectorsRegister = make(map[string]Inspector)

	for _, i := range all {
		InspectorsRegister[i.Name()] = i
	}
}

type Inspector interface {
	Name() string
	Inspect(cloud.GraphAPI) error
	Print(io.Writer)
}
