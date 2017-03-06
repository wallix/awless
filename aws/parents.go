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
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/graph"
)

const (
	PARENT_OF = iota // default
	APPLIES_ON
	DEPENDING_ON
)

type funcBuilder struct {
	parent                              graph.ResourceType
	fieldName, listName, stringListName string
	relation                            int
}

type addParentFn func(*graph.Graph, interface{}) error

var addParentsFns = map[string][]addParentFn{
	// Infra
	graph.Subnet.String(): {
		funcBuilder{parent: graph.Vpc, fieldName: "VpcId"}.build(),
	},
	graph.Instance.String(): {
		funcBuilder{parent: graph.Subnet, fieldName: "SubnetId"}.build(),
		funcBuilder{parent: graph.SecurityGroup, fieldName: "GroupId", listName: "SecurityGroups", relation: APPLIES_ON}.build(),
		funcBuilder{parent: graph.Keypair, fieldName: "KeyName", relation: APPLIES_ON}.build(),
	},
	graph.SecurityGroup.String(): {
		funcBuilder{parent: graph.Vpc, fieldName: "VpcId"}.build(),
	},
	graph.InternetGateway.String(): {
		addRegionParent,
		funcBuilder{parent: graph.Vpc, fieldName: "VpcId", listName: "Attachments", relation: DEPENDING_ON}.build(),
	},
	graph.RouteTable.String(): {
		funcBuilder{parent: graph.Subnet, fieldName: "SubnetId", listName: "Associations", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: graph.Vpc, fieldName: "VpcId"}.build(),
	},
	graph.Volume.String(): {
		funcBuilder{parent: graph.AvailabilityZone, fieldName: "AvailabilityZone"}.build(),
		funcBuilder{parent: graph.Instance, fieldName: "InstanceId", listName: "Attachments", relation: DEPENDING_ON}.build(),
	},
	graph.LoadBalancer.String(): {
		funcBuilder{parent: graph.Vpc, fieldName: "VpcId"}.build(),
		funcBuilder{parent: graph.Subnet, fieldName: "SubnetId", listName: "AvailabilityZones", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: graph.AvailabilityZone, fieldName: "ZoneName", listName: "AvailabilityZones", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: graph.SecurityGroup, stringListName: "SecurityGroups", relation: APPLIES_ON}.build(),
	},
	graph.Listener.String(): {
		funcBuilder{parent: graph.LoadBalancer, fieldName: "LoadBalancerArn"}.build(),
	},
	graph.TargetGroup.String(): {
		funcBuilder{parent: graph.Vpc, fieldName: "VpcId"}.build(),
		funcBuilder{parent: graph.LoadBalancer, stringListName: "LoadBalancerArns", relation: APPLIES_ON}.build(),
		fetchTargetsAndAddRelations,
	},
	graph.Vpc.String():              {addRegionParent},
	graph.AvailabilityZone.String(): {addRegionParent},
	graph.Keypair.String():          {addRegionParent},
	graph.User.String():             {addRegionParent, userAddGroupsRelations, addManagedPoliciesRelations},
	graph.Role.String():             {addRegionParent, addManagedPoliciesRelations},
	graph.Group.String():            {addRegionParent, addManagedPoliciesRelations},
	graph.Policy.String():           {addRegionParent},
	graph.Bucket.String():           {addRegionParent},
}

func (fb funcBuilder) build() addParentFn {
	switch {
	case fb.listName != "":
		return fb.addRelationListWithField()
	case fb.stringListName != "":
		return fb.addRelationListWithStringField()
	default:
		return fb.addRelationWithField()
	}
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
			return nil
		}

		parent, err := g.GetResource(fb.parent, awssdk.StringValue(str))
		if err != nil {
			return err
		}

		return addRelation(g, parent, res, fb.relation)
	}
}

func (fb funcBuilder) addRelationListWithStringField() addParentFn {
	return func(g *graph.Graph, i interface{}) error {
		structField, err := verifyValidStructField(i, fb.stringListName)
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
			str, ok := structField.Index(i).Interface().(*string)
			if !ok {
				return fmt.Errorf("add parent to %s: not a string pointer: %T", res.Id(), str)
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
				return fmt.Errorf("add parent to %s: unknown field %s in %d", res.Id(), listStructField, i)
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
	case DEPENDING_ON:
		g.AddAppliesOnRelation(other, first)
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
		return fmt.Errorf("add parent to %s: unknown field %s in %d", res.Id(), structField, i)
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
	n, err := initResource(user)
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

func fetchTargetsAndAddRelations(g *graph.Graph, i interface{}) error {
	group, ok := i.(*elbv2.TargetGroup)
	if !ok {
		return fmt.Errorf("add targets relation: not a target group, but a %T", i)
	}
	parent, err := initResource(group)
	if err != nil {
		return err
	}

	targets, err := InfraService.(*Infra).DescribeTargetHealth(&elbv2.DescribeTargetHealthInput{TargetGroupArn: group.TargetGroupArn})
	if err != nil {
		return err
	}

	for _, t := range targets.TargetHealthDescriptions {
		n, err := g.GetResource(graph.Instance, awssdk.StringValue(t.Target.Id))
		if err != nil {
			return err
		}
		g.AddAppliesOnRelation(parent, n)
	}
	return nil
}
