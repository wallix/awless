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
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
)

const (
	PARENT_OF = iota // default
	APPLIES_ON
	DEPENDING_ON
)

type funcBuilder struct {
	parent                              string
	fieldName, listName, stringListName string
	relation                            int
}

type addParentFn func(*graph.Graph, interface{}) error

var addParentsFns = map[string][]addParentFn{
	// Infra
	cloud.Subnet: {
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId"}.build(),
	},
	cloud.Instance: {
		funcBuilder{parent: cloud.Subnet, fieldName: "SubnetId"}.build(),
		funcBuilder{parent: cloud.SecurityGroup, fieldName: "GroupId", listName: "SecurityGroups", relation: APPLIES_ON}.build(),
		funcBuilder{parent: cloud.Keypair, fieldName: "KeyName", relation: APPLIES_ON}.build(),
	},
	cloud.SecurityGroup: {
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId"}.build(),
	},
	cloud.InternetGateway: {
		addRegionParent,
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId", listName: "Attachments", relation: DEPENDING_ON}.build(),
	},
	cloud.RouteTable: {
		funcBuilder{parent: cloud.Subnet, fieldName: "SubnetId", listName: "Associations", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId"}.build(),
	},
	cloud.Volume: {
		funcBuilder{parent: cloud.AvailabilityZone, fieldName: "AvailabilityZone"}.build(),
		funcBuilder{parent: cloud.Instance, fieldName: "InstanceId", listName: "Attachments", relation: DEPENDING_ON}.build(),
	},
	// Loadbalancer
	cloud.LoadBalancer: {
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId"}.build(),
		funcBuilder{parent: cloud.Subnet, fieldName: "SubnetId", listName: "AvailabilityZones", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: cloud.AvailabilityZone, fieldName: "ZoneName", listName: "AvailabilityZones", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: cloud.SecurityGroup, stringListName: "SecurityGroups", relation: APPLIES_ON}.build(),
	},
	cloud.Listener: {
		funcBuilder{parent: cloud.LoadBalancer, fieldName: "LoadBalancerArn"}.build(),
	},
	cloud.TargetGroup: {
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId"}.build(),
		funcBuilder{parent: cloud.LoadBalancer, stringListName: "LoadBalancerArns", relation: APPLIES_ON}.build(),
		fetchTargetsAndAddRelations,
	},
	// Database
	cloud.Database: {
		funcBuilder{parent: cloud.AvailabilityZone, fieldName: "AvailabilityZone"}.build(),
		funcBuilder{parent: cloud.SecurityGroup, listName: "VpcSecurityGroups", fieldName: "VpcSecurityGroupId", relation: APPLIES_ON}.build(),
	},
	cloud.Vpc:              {addRegionParent},
	cloud.AvailabilityZone: {addRegionParent},
	cloud.Keypair:          {addRegionParent},
	cloud.User:             {userAddGroupsRelations, addManagedPoliciesRelations},
	cloud.Role:             {addManagedPoliciesRelations},
	cloud.Group:            {addManagedPoliciesRelations},
	cloud.Bucket:           {addRegionParent},
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

		parent := graph.InitResource(fb.parent, awssdk.StringValue(str))
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
			parent := graph.InitResource(fb.parent, awssdk.StringValue(str))

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
			parent := graph.InitResource(fb.parent, awssdk.StringValue(str))

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
	resources, err := g.GetAllResources(cloud.Region)
	if err != nil {
		return err
	}
	if len(resources) != 1 {
		return fmt.Errorf("aws fetch: expect exactly one region in cloud. but got %d", len(resources))
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
		pid, ok := a.ResolveToId(g, cloud.Policy)
		if !ok {
			fmt.Fprintf(os.Stderr, "add parent to '%s/%s': unknown policy named '%s'. Ignoring it.\n", res.Type(), res.Id(), awssdk.StringValue(policy.PolicyName))
			return nil
		}
		parent := graph.InitResource(cloud.Policy, pid)
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
		groupName := awssdk.StringValue(group)
		resources, err := g.ResolveResources(&graph.And{Resolvers: []graph.Resolver{
			&graph.ByProperty{Name: "Name", Val: groupName},
			&graph.ByType{Typ: cloud.Group},
		}})
		if err != nil {
			return err
		}
		switch len(resources) {
		case 0:
			fmt.Fprintf(os.Stderr, "no group with name %s found for user %s\n", groupName, n.Id())
		case 1:
			g.AddAppliesOnRelation(resources[0], n)
		default:
			fmt.Fprintf(os.Stderr, "multiple groups with name %s found for user %s:%v\n", groupName, n.Id(), resources)
		}
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
		n := graph.InitResource(cloud.Instance, awssdk.StringValue(t.Target.Id))
		g.AddAppliesOnRelation(parent, n)
	}
	return nil
}
