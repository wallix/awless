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

package awsspecmeta

func Lookuper(action, entity string, paramKeys []string) interface{} {
	switch action + "." + entity {
	case "create.internetgateway":
		m := &CreateInternetgatewayMeta{}
		if m.Match(action, entity, paramKeys) {
			return m
		}
	case "attach.policy":
		m := &AttachPolicyMeta{}
		if m.Match(action, entity, paramKeys) {
			return m
		}
	}
	return nil
}

func contains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}
