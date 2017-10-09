/* Copyright 2017 WALLIX

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
package awsspec

import (
	"testing"
)

func TestCreateSubnet(t *testing.T) {
	create := &CreateSubnet{}
	params := map[string]interface{}{
		"cidr": "10.10.10.10/24",
		"vpc":  "vpc-1234",
	}
	t.Run("Validate", func(t *testing.T) {
		params["cidr"] = "10.10.10.10/128"
		errs := create.ValidateCommand(params, nil)
		checkErrs(t, errs, 1, "cidr")

	})
}
