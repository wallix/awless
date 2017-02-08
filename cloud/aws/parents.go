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

type addParentFn func(*graph.Graph, interface{}) error

var addParentsFns = map[string][]addParentFn{
	graph.Vpc.String():             {addRegionParent},
	graph.Subnet.String():          {addParentWithFieldIfExists(graph.Vpc, "VpcId")},
	graph.Instance.String():        {addParentWithFieldIfExists(graph.Subnet, "SubnetId"), addParentListWithField(graph.SecurityGroup, "SecurityGroups", "GroupId")},
	graph.SecurityGroup.String():   {addParentWithFieldMustExist(graph.Vpc, "VpcId")},
	graph.Keypair.String():         {addRegionParent},
	graph.InternetGateway.String(): {addRegionParent, addParentListWithField(graph.Vpc, "Attachments", "VpcId")},
	graph.RouteTable.String():      {addParentListWithField(graph.Subnet, "Associations", "SubnetId"), addParentWithFieldMustExist(graph.Vpc, "VpcId")},
	graph.User.String():            {addRegionParent, userAddGroupsParents, addManagedPoliciesParents},
	graph.Role.String():            {addRegionParent, addManagedPoliciesParents},
	graph.Group.String():           {addRegionParent, addManagedPoliciesParents},
	graph.Policy.String():          {addRegionParent},
	graph.Bucket.String():          {addRegionParent},
}

var ParentNotFound = errors.New("empty field to add parent")

func addParentWithFieldIfExists(parentT graph.ResourceType, field string) addParentFn {
	return func(g *graph.Graph, i interface{}) error {
		err := addParentWithFieldMustExist(parentT, field)(g, i)
		if err != ParentNotFound {
			return err
		}
		return nil
	}
}

func addParentWithFieldMustExist(parentT graph.ResourceType, field string) addParentFn {
	return func(g *graph.Graph, i interface{}) error {
		res, err := initResource(i)
		if err != nil {
			return err
		}
		value := reflect.ValueOf(i)
		if value.Kind() != reflect.Ptr {
			return fmt.Errorf("add parent to %s: %T not a pointer", res.Id(), i)
		}
		struc := value.Elem()
		if struc.Kind() != reflect.Struct {
			return fmt.Errorf("add parent to %s: %T not a stuct pointer", res.Id(), i)
		}

		structField := struc.FieldByName(field)
		if !structField.IsValid() {
			return fmt.Errorf("add parent to %s: ", field, i)
		}
		str, ok := structField.Interface().(*string)
		if !ok {
			return fmt.Errorf("add parent to %s: %T not a string pointer", field, structField.Interface())
		}

		if awssdk.StringValue(str) == "" {
			return ParentNotFound
		}
		parent, err := g.GetResource(parentT, awssdk.StringValue(str))
		if err != nil {
			return err
		}
		g.AddParent(parent, res)
		return nil
	}
}

func addParentListWithField(parentT graph.ResourceType, listField, parentField string) addParentFn {
	return func(g *graph.Graph, i interface{}) error {
		res, err := initResource(i)
		if err != nil {
			return err
		}
		value := reflect.ValueOf(i)
		if value.Kind() != reflect.Ptr {
			return fmt.Errorf("add parent to %s: %T not a pointer", res.Id(), i)
		}
		struc := value.Elem()
		if struc.Kind() != reflect.Struct {
			return fmt.Errorf("add parent to %s: %T not a struct pointer", res.Id(), i)
		}

		structField := struc.FieldByName(listField)
		if !structField.IsValid() {
			return fmt.Errorf("add parent to %s: %T not a struct field", structField, i)
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
			listStructField := listStruc.FieldByName(parentField)
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
			parent, err := g.GetResource(parentT, awssdk.StringValue(str))
			if err != nil {
				return err
			}
			g.AddParent(parent, res)
		}
		return nil
	}
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
	g.AddParent(regionN, res)
	return nil
}

func addManagedPoliciesParents(g *graph.Graph, i interface{}) error {
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
		g.AddParent(parent, res)
	}
	return nil
}

func userAddGroupsParents(g *graph.Graph, i interface{}) error {
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
		g.AddParent(parent, n)
	}
	return nil
}
