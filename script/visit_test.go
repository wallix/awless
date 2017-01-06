package script

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/wallix/awless/script/ast"
	"github.com/wallix/awless/script/driver"
)

func TestVisitHoles(t *testing.T) {
	s := &ast.Script{}

	n := &ast.ExpressionNode{
		Action: "create", Entity: "vpc",
		Holes: map[string]string{"name": "presidentName", "rank": "presidentRank"},
	}
	s.Statements = append(s.Statements, n)

	fills := map[string]interface{}{
		"presidentName": "trump",
		"presidentRank": 45,
	}

	VisitHoles(s, fills)

	expected := map[string]interface{}{"name": "trump", "rank": 45}
	if got, want := n.Params, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := len(n.Holes), 0; got != want {
		t.Fatalf("length of holes: got %d, want %d", got, want)
	}
}

func TestVisitWithDriver(t *testing.T) {
	s := &ast.Script{}

	n := &ast.ExpressionNode{
		Action: "create", Entity: "vpc",
		Params: map[string]interface{}{"count": 1},
	}
	s.Statements = append(s.Statements, n)

	mDriver := &mockDriver{
		action: "create", entity: "vpc",
		expectedParams: map[string]interface{}{"count": 1},
	}

	if err := Visit(s, mDriver); err != nil {
		t.Fatal(err)
	}
}

type mockDriver struct {
	action, entity string
	expectedParams map[string]interface{}
}

func (r *mockDriver) Lookup(lookups ...string) driver.DriverFn {
	if lookups[0] == r.action && lookups[1] == r.entity {
		return func(params map[string]interface{}) (interface{}, error) {
			if got, want := params, r.expectedParams; !reflect.DeepEqual(got, want) {
				return nil, fmt.Errorf("[%s %s] params mismatch: expected %v, got %v", r.action, r.entity, got, want)
			}
			return nil, nil
		}
	}

	return func(params map[string]interface{}) (interface{}, error) {
		return nil, errors.New("Unexpected lookup fallthrough")
	}
}

func (r *mockDriver) SetLogger(*log.Logger) {}
