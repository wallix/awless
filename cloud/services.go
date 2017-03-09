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
	"errors"
	"fmt"
	"strings"

	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/template/driver"
)

var ErrFetchAccessDenied = errors.New("access denied to cloud resource")

type Service interface {
	Name() string
	Drivers() []driver.Driver
	ResourceTypes() []string
	FetchResources() (*graph.Graph, error)
	IsSyncDisabled() bool
	FetchByType(t string) (*graph.Graph, error)
}

var ServiceRegistry = make(map[string]Service)

func GetServiceForType(t string) (Service, error) {
	for _, srv := range ServiceRegistry {
		for _, typ := range srv.ResourceTypes() {
			if typ == t {
				return srv, nil
			}
		}
	}
	return nil, fmt.Errorf("cannot find cloud service for resource type %s", t)
}

func PluralizeResource(singular string) string {
	if strings.HasSuffix(singular, "cy") {
		return strings.TrimSuffix(singular, "y") + "ies"
	}
	return singular + "s"
}
