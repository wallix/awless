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

package aws

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/graph"
)

const (
	PARENT_OF = iota // default
	APPLIES_ON
)

type funcBuilder struct {
	parent              graph.ResourceType
	fieldName, listName string
	relation            int
	must                bool
}

type addParentFn func(*graph.Graph, interface{}) error

var addParentsFns = map[string][]addParentFn{
	graph.Subnet.String(): {
		funcBuilder{parent: graph.Vpc, fieldName: "VpcId"}.build(),
	},
	graph.Instance.String(): {
		funcBuilder{parent: graph.Subnet, fieldName: "SubnetId"}.build(),
		funcBuilder{parent: graph.SecurityGroup, fieldName: "GroupId", listName: "SecurityGroups", relation: APPLIES_ON}.build(),
		funcBuilder{parent: graph.Keypair, fieldName: "KeyName", relation: APPLIES_ON}.build(),
	},
	graph.SecurityGroup.String(): {
		funcBuilder{parent: graph.Vpc, fieldName: "VpcId", must: true}.build(),
	},
	graph.InternetGateway.String(): {
		addRegionParent,
		funcBuilder{parent: graph.Vpc, fieldName: "VpcId", listName: "Attachments"}.build(),
	},
	graph.RouteTable.String(): {
		funcBuilder{parent: graph.Subnet, fieldName: "SubnetId", listName: "Associations"}.build(),
		funcBuilder{parent: graph.Vpc, fieldName: "VpcId", must: true}.build(),
	},
	graph.Vpc.String():     {addRegionParent},
	graph.Keypair.String(): {addRegionParent},
	graph.User.String():    {addRegionParent, userAddGroupsRelations, addManagedPoliciesRelations},
	graph.Role.String():    {addRegionParent, addManagedPoliciesRelations},
	graph.Group.String():   {addRegionParent, addManagedPoliciesRelations},
	graph.Policy.String():  {addRegionParent},
	graph.Bucket.String():  {addRegionParent},
}

var ParentNotFound = errors.New("empty field to add parent")

func (fb funcBuilder) build() addParentFn {
	if fb.listName != "" {
		return fb.addRelationListWithField()
	}

	return fb.addRelationWithField()
}

func (fb funcBuilder) addRelationWithField() addParentFn {
	return func(g *graph.Graph, i interface{}) error {
		structField, err := verifyValidStructField(i, fb.fieldName)
		if err != nil {
			return err
		}

		str, ok := structField.Interface().(*string)
		if !ok {
			return fmt.Errorf("add parent to %s: %T not a string pointer", fb.fieldName, structField.Interface())
		}

		res, err := initResource(i)
		if err != nil {
			return err
		}

		if awssdk.StringValue(str) == "" {
			if fb.must {
				return ParentNotFound
			} else {
				return nil
			}
		}

		parent, err := g.GetResource(fb.parent, awssdk.StringValue(str))
		if err != nil {
			return err
		}

		return addRelation(g, parent, res, fb.relation)
	}
}

func (fb funcBuilder) addRelationListWithField() addParentFn {
	return func(g *graph.Graph, i interface{}) error {
		structField, err := verifyValidStructField(i, fb.listName)
		if err != nil {
			return err
		}

		res, err := initResource(i)
		if err != nil {
			return err
		}

		if !structField.IsValid() || structField.Kind() != reflect.Slice {
			return fmt.Errorf("add parent to %s: field not a slice: %T", res.Id(), structField.Kind())
		}

		for i := 0; i < structField.Len(); i++ {
			listValue := structField.Index(i)
			if listValue.Kind() != reflect.Ptr {
				return fmt.Errorf("add parent to %s: not a pointer: %s", res.Id(), listValue.Kind())
			}
			listStruc := listValue.Elem()
			if listStruc.Kind() != reflect.Struct {
				return fmt.Errorf("add parent to %s: not a struct: %s", res.Id(), listStruc.Kind())
			}
			listStructField := listStruc.FieldByName(fb.fieldName)
			if !listStructField.IsValid() {
				return fmt.Errorf("add parent to %s: unknown field %s in %T", listStructField, i)
			}
			str, ok := listStructField.Interface().(*string)
			if !ok {
				return fmt.Errorf("add parent to %s: %T is not a string pointer", listStructField, listStructField.Interface())
			}

			if awssdk.StringValue(str) == "" {
				continue
			}
			parent, err := g.GetResource(fb.parent, awssdk.StringValue(str))
			if err != nil {
				return err
			}

			if err = addRelation(g, parent, res, fb.relation); err != nil {
				return err
			}
		}
		return nil
	}
}

func verifyValidStructField(i interface{}, name string) (reflect.Value, error) {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("%T not a pointer", i)
	}
	struc := value.Elem()
	if struc.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("%T not a stuct pointer", i)
	}

	structField := struc.FieldByName(name)
	if !structField.IsValid() {
		return reflect.Value{}, fmt.Errorf("invalid field %s: ", name)
	}

	return structField, nil
}

func addRelation(g *graph.Graph, first, other *graph.Resource, relation int) error {
	switch relation {
	case PARENT_OF:
		g.AddParentRelation(first, other)
	case APPLIES_ON:
		g.AddAppliesOnRelation(first, other)
	default:
		return errors.New("unknown relation type")
	}
	return nil
}

func addRegionParent(g *graph.Graph, i interface{}) error {
	resources, err := g.GetAllResources(graph.Region)
	if err != nil {
		return err
	}
	if len(resources) != 1 {
		return fmt.Errorf("aws fetch: expect exactly one region in graph, but got %d", len(resources))
	}
	regionN := resources[0]
	res, err := initResource(i)
	if err != nil {
		return err
	}
	g.AddParentRelation(regionN, res)
	return nil
}

func addManagedPoliciesRelations(g *graph.Graph, i interface{}) error {
	res, err := initResource(i)
	if err != nil {
		return err
	}
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("add parent to %s: unknown type %T", res.Id(), i)
	}
	struc := value.Elem()
	if struc.Kind() != reflect.Struct {
		return fmt.Errorf("add parent to %s: unknown type %T", res.Id(), i)
	}

	structField := struc.FieldByName("AttachedManagedPolicies")
	if !structField.IsValid() {
		return fmt.Errorf("add parent to %s: unknown field %s in %T", structField, i)
	}
	policies, ok := structField.Interface().([]*iam.AttachedPolicy)
	if !ok {
		return fmt.Errorf("add parent to %s: not a valid attached policy list: %T", res.Id(), structField.Interface())
	}

	for _, policy := range policies {
		a := graph.Alias(awssdk.StringValue(policy.PolicyName))
		pid, ok := a.ResolveToId(g, graph.Policy)
		if !ok {
			fmt.Fprintf(os.Stderr, "add parent to '%s/%s': unknown policy named '%s'. Ignoring it.\n", res.Type(), res.Id(), awssdk.StringValue(policy.PolicyName))
			return nil
		}
		parent, err := g.GetResource(graph.Policy, pid)
		if err != nil {
			return err
		}
		g.AddAppliesOnRelation(parent, res)
	}
	return nil
}

func userAddGroupsRelations(g *graph.Graph, i interface{}) error {
	user, ok := i.(*iam.UserDetail)
	if !ok {
		return fmt.Errorf("aws fetch: not a user, but a %T", i)
	}
	n, err := g.GetResource(graph.User, awssdk.StringValue(user.UserId))
	if err != nil {
		return err
	}

	for _, group := range user.GroupList {
		parent, err := g.GetResource(graph.Group, awssdk.StringValue(group))
		if err != nil {
			return err
		}
		g.AddAppliesOnRelation(parent, n)
	}
	return nil
}
