package params_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/wallix/awless/template/params"
)

func TestListParams(t *testing.T) {
	tcases := []struct {
		rules  params.Rule
		exp    []string
		expOpt []string
		expSug []string
	}{
		{rules: params.Opt("b", "c", "a"), expOpt: []string{"a", "b", "c"}},
		{rules: params.AllOf(params.Key("f"), params.Opt("a", params.Suggested("b"))), exp: []string{"f"}, expOpt: []string{"a", "b"}, expSug: []string{"b"}},
		{rules: params.OnlyOneOf(params.Key("b"), params.Key("c"), params.Opt("a")), exp: []string{"b", "c"}, expOpt: []string{"a"}},
		{rules: params.AtLeastOneOf(params.Key("b"), params.Key("c"), params.Opt(params.Suggested("a", "f"))), exp: []string{"b", "c"}, expOpt: []string{"a", "f"}, expSug: []string{"a", "f"}},
	}

	for i, tcase := range tcases {
		actual, actualOpt, actualSug := params.List(tcase.rules)
		if got, want := actual, tcase.exp; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d. actual: got %v, want %v", i+1, got, want)
		}
		if got, want := actualOpt, tcase.expOpt; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d. actual opt: got %v, want %v", i+1, got, want)
		}
		if got, want := actualSug, tcase.expSug; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d. suggested: got %v, want %v", i+1, got, want)
		}
	}

}

func TestPrintRules(t *testing.T) {
	tcases := []struct {
		rules params.Rule
		out   string
	}{
		{rules: params.AllOf(params.Key("1"), params.Key("2"), params.OnlyOneOf(params.Key("3"), params.Key("4")), params.AtLeastOneOf(params.Key("5"), params.Key("6"))),
			out: "1 + 2 + (3 | 4) + (5 / 6)"},
		{rules: params.AllOf(params.OnlyOneOf(params.Key("user"), params.Key("group"), params.Key("role")), params.OnlyOneOf(params.Key("arn"), params.AllOf(params.Key("service"), params.Key("access")))),
			out: "(user | group | role) + (arn | service + access)"},
	}

	for _, tcase := range tcases {
		if got, want := tcase.rules.String(), tcase.out; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	}
}

func TestRuleMissing(t *testing.T) {
	tcases := []struct {
		rules   params.Rule
		in      []string
		missing []string
	}{
		{rules: params.AllOf()},
		{rules: params.OnlyOneOf()},
		{rules: params.AtLeastOneOf()},

		{rules: params.AllOf(params.Key("1")), missing: []string{"1"}},
		{rules: params.OnlyOneOf(params.Key("1")), missing: []string{"1"}},
		{rules: params.AtLeastOneOf(params.Key("1")), missing: []string{"1"}},

		{rules: params.AllOf(params.Key("2"), params.Key("1")), missing: []string{"2", "1"}},
		{rules: params.OnlyOneOf(params.Key("2"), params.Key("1")), missing: []string{"2"}},
		{rules: params.AtLeastOneOf(params.Key("2"), params.Key("1")), missing: []string{"2"}},

		{rules: params.AllOf(params.Key("1"), params.Key("2")), in: []string{"1", "2"}},
		{rules: params.AllOf(params.Key("1"), params.Key("2")), in: []string{"1"}, missing: []string{"2"}},
		{rules: params.OnlyOneOf(params.Key("1"), params.Key("2")), in: []string{"1"}},
		{rules: params.AtLeastOneOf(params.Key("2"), params.Key("1")), in: []string{"2"}},

		{rules: params.OnlyOneOf(
			params.Key("5"),
			params.OnlyOneOf(params.Key("1"), params.Key("2")),
			params.AtLeastOneOf(params.Key("3"), params.Key("4"))), missing: []string{"5"}},
		{rules: params.OnlyOneOf(
			params.OnlyOneOf(params.Key("1"), params.Key("2")),
			params.Key("5"),
			params.AtLeastOneOf(params.Key("3"), params.Key("4"))), missing: []string{"1"}},
		{rules: params.OnlyOneOf(
			params.OnlyOneOf(params.Key("1"), params.Key("2")),
			params.Key("5")), missing: []string{"1"}},
		{rules: params.OnlyOneOf(
			params.AtLeastOneOf(params.Key("3"), params.Key("4")),
			params.Key("5"),
			params.OnlyOneOf(params.Key("1"), params.Key("2")),
			params.AtLeastOneOf(params.Key("3"), params.Key("4"))), in: []string{"3"}},

		{rules: params.AllOf(
			params.Key("5"),
			params.OnlyOneOf(params.Key("1"), params.Key("2")),
			params.AtLeastOneOf(params.Key("3"), params.Key("4"))), missing: []string{"5", "1", "3"}},
		{rules: params.AllOf(
			params.Key("5"),
			params.OnlyOneOf(params.Key("1"), params.Key("2")),
			params.AtLeastOneOf(params.Key("3"), params.Key("4"))), in: []string{"5"}, missing: []string{"1", "3"}},
		{rules: params.AllOf(
			params.Key("5"),
			params.OnlyOneOf(params.Key("1"), params.Key("2")),
			params.AtLeastOneOf(params.Key("3"), params.Key("4"))), in: []string{"5", "1"}, missing: []string{"3"}},
		{rules: params.AllOf(
			params.Key("5"),
			params.OnlyOneOf(params.Key("1"), params.Key("2")),
			params.AtLeastOneOf(params.Key("3"), params.Key("4"))), in: []string{"5", "3"}, missing: []string{"1"}},
	}

	for i, tcase := range tcases {
		if got, want := tcase.rules.Missing(tcase.in), tcase.missing; !reflect.DeepEqual(got, want) {
			t.Fatalf("missing: %d: got %v, want %v", i+1, got, want)
		}
	}
}

func TestUnexpected(t *testing.T) {
	tcases := []struct {
		rules       params.Rule
		in          []string
		errContains []string
	}{
		{rules: params.AllOf(params.Key("1"), params.Key("2")), in: []string{"3"}, errContains: []string{"unexpected", "3"}},
		{rules: params.OnlyOneOf(params.Key("1"), params.Key("2")), in: []string{"3"}, errContains: []string{"unexpected", "3"}},
		{rules: params.AtLeastOneOf(params.Key("1"), params.Key("2")), in: []string{"3"}, errContains: []string{"unexpected", "3"}},
	}

	for i, tcase := range tcases {
		err := params.Run(tcase.rules, tcase.in)

		if len(tcase.errContains) > 0 {
			if err == nil {
				t.Fatalf("%d: expected error got none", i+1)
			}
			msg := err.Error()
			for _, e := range tcase.errContains {
				if !strings.Contains(msg, e) {
					t.Fatalf("%d: expected %s to contain %s", i+1, msg, e)
				}
			}
		}
	}
}

func TestValidateRule(t *testing.T) {
	tcases := []struct {
		rules       params.Rule
		in          []string
		errContains []string
	}{
		{rules: params.AllOf()},
		{rules: params.OnlyOneOf()},
		{rules: params.AtLeastOneOf()},

		{rules: params.AllOf(params.Key("1")), in: []string{"1", "2"}},
		{rules: params.OnlyOneOf(params.Key("1")), in: []string{"1"}},
		{rules: params.AtLeastOneOf(params.Key("1")), in: []string{"1"}},

		{rules: params.AllOf(params.Key("1"), params.Key("2")), in: []string{"1", "2"}},
		{rules: params.OnlyOneOf(params.Key("1"), params.Key("2")), in: []string{"1"}},
		{rules: params.OnlyOneOf(params.Key("1"), params.Key("2")), in: []string{"2"}},
		{rules: params.AtLeastOneOf(params.Key("1"), params.Key("2")), in: []string{"1"}},
		{rules: params.AtLeastOneOf(params.Key("1"), params.Key("2")), in: []string{"2"}},

		{rules: params.AllOf(params.Key("1"), params.Key("2")), in: []string{"1"}, errContains: []string{"2"}},
		{rules: params.AllOf(params.Key("1")), errContains: []string{"1"}},
		{rules: params.OnlyOneOf(params.Key("1")), errContains: []string{"1"}},
		{rules: params.OnlyOneOf(params.Key("1"), params.Key("2")), in: []string{"1", "2"}, errContains: []string{"1", "2"}},
		{rules: params.AtLeastOneOf(params.Key("1"), params.Key("2")), errContains: []string{"1", "2"}},
		{rules: params.AtLeastOneOf(params.Key("1")), errContains: []string{"1"}},

		{rules: params.OnlyOneOf(
			params.AllOf(params.Key("instance"), params.Key("id")),
			params.Key("attachment"),
			params.Opt("force")),
			in: []string{"attachment"}},

		{rules: params.AllOf(
			params.OnlyOneOf(params.Key("distro"), params.Key("image")),
			params.Key("count"), params.Key("type"), params.Key("name"), params.Key("subnet")),
			in: []string{"image", "count", "name", "subnet", "type"}},

		{rules: params.AllOf(
			params.OnlyOneOf(params.Key("1"), params.Key("2")),
			params.Key("3")), errContains: []string{"1", "2", "3"}},

		{rules: params.AllOf(
			params.AtLeastOneOf(params.Key("1"), params.Key("2")),
			params.Key("3")), errContains: []string{"1", "3"}},
	}

	for i, tcase := range tcases {
		err := params.Run(tcase.rules, tcase.in)

		if len(tcase.errContains) > 0 {
			if err == nil {
				t.Fatalf("%d: expected error got none", i+1)
			}
			msg := err.Error()
			for _, e := range tcase.errContains {
				if !strings.Contains(msg, e) {
					t.Fatalf("%d: expected %s to contain %s", i+1, msg, e)
				}
			}
		}
	}
}

func TestRulesRequired(t *testing.T) {
	tcases := []struct {
		rules    params.Rule
		required []string
	}{
		{rules: params.AllOf()},
		{rules: params.OnlyOneOf()},
		{rules: params.AtLeastOneOf()},

		{rules: params.AllOf(params.Key("1")), required: []string{"1"}},
		{rules: params.OnlyOneOf(params.Key("1")), required: []string{"1"}},
		{rules: params.AtLeastOneOf(params.Key("1")), required: []string{"1"}},

		{rules: params.AllOf(params.Key("1"), params.Key("2")), required: []string{"1", "2"}},
		{rules: params.OnlyOneOf(params.Key("1"), params.Key("2")), required: []string{"1"}},
		{rules: params.AtLeastOneOf(params.Key("2"), params.Key("1")), required: []string{"2"}},

		{rules: params.AllOf(
			params.Key("5"),
			params.OnlyOneOf(params.Key("1"), params.Key("2")),
			params.AtLeastOneOf(params.Key("3"), params.Key("4"))), required: []string{"5", "1", "3"}},
		{rules: params.OnlyOneOf(
			params.AllOf(params.Key("1"), params.Key("2")),
			params.AtLeastOneOf(params.Key("3"), params.Key("4"))), required: []string{"1", "2"}},
		{rules: params.OnlyOneOf(
			params.AtLeastOneOf(params.Key("3"), params.Key("4")),
			params.AllOf(params.Key("1"), params.Key("2"))), required: []string{"3"}},
	}

	for _, tcase := range tcases {
		if got, want := tcase.rules.Required(), tcase.required; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
