package ast

import (
	"reflect"
	"testing"
)

func TestCompositeValues(t *testing.T) {
	tcases := []struct {
		val          CompositeValue
		holesFillers map[string]interface{}
		refsFillers  map[string]interface{}
		expHoles     map[string][]string
		expRefs      []string
		expAliases   []string
		expValue     interface{}
	}{
		{val: &interfaceValue{val: "test"}, expValue: "test"},
		{val: &interfaceValue{val: "test with spe${cial ch'ars"}, expValue: "test with spe${cial ch'ars"},
		{val: &interfaceValue{val: 10}, expValue: 10},
		{val: &holeValue{hole: "myhole"}, expHoles: map[string][]string{"myhole": nil}},
		{val: &referenceValue{ref: "myref"}, expRefs: []string{"myref"}},
		{val: &aliasValue{alias: "myalias"}, expAliases: []string{"myalias"}},
		{
			val: newCompositeValue(
				&interfaceValue{val: "test"},
				&interfaceValue{val: 10},
				&holeValue{hole: "myhole"},
				&referenceValue{ref: "myref"},
				&aliasValue{alias: "myalias"},
			),
			expRefs:    []string{"myref"},
			expHoles:   map[string][]string{"myhole": nil},
			expValue:   []interface{}{"test", 10},
			expAliases: []string{"myalias"},
		},
		{val: &holeValue{hole: "myhole"}, holesFillers: map[string]interface{}{"myhole": "my-value"}, expValue: "my-value"},
		{val: &holeValue{hole: "myhole"}, holesFillers: map[string]interface{}{"myhole": &aliasValue{alias: "myalias"}}, expValue: nil, expAliases: []string{"myalias"}},
		{val: &holeValue{hole: "myhole"}, holesFillers: map[string]interface{}{"myhole": newCompositeValue(&aliasValue{alias: "myalias1"}, &aliasValue{alias: "myalias2"})}, expValue: nil, expAliases: []string{"myalias1", "myalias2"}},
		{
			val: newCompositeValue(
				&interfaceValue{val: "test"},
				&interfaceValue{val: 10},
				&holeValue{hole: "myhole"},
				&referenceValue{ref: "myref"},
			),
			refsFillers:  map[string]interface{}{"myref": "refvalue"},
			holesFillers: map[string]interface{}{"myhole": "my-value"},
			expValue:     []interface{}{"test", 10, "my-value", "refvalue"},
		},
		{
			val: &concatenationValue{
				vals: []CompositeValue{&interfaceValue{val: "prefix-"}, &holeValue{hole: "hole1"}, &interfaceValue{val: "middle1-"}, &holeValue{hole: "hole2"}, &interfaceValue{val: "-middle2-"}, &holeValue{hole: "hole3"}, &interfaceValue{val: "suffix"}},
			},
			expHoles: map[string][]string{"hole1": nil, "hole2": nil, "hole3": nil},
			expValue: "prefix-middle1--middle2-suffix",
		},
		{
			val: &concatenationValue{
				vals: []CompositeValue{&interfaceValue{val: "prefix-"}, &holeValue{hole: "hole1"}, &interfaceValue{val: "middle1-"}, &holeValue{hole: "hole2.name"}, &interfaceValue{val: "-middle2-"}, &holeValue{hole: "hole3"}, &interfaceValue{val: "suffix"}},
			},
			holesFillers: map[string]interface{}{"hole1": "value1", "hole2.name": 2, "hole3": "value3"},
			expValue:     "prefix-value1middle1-2-middle2-value3suffix",
		},
		{
			val: &concatenationValue{
				vals: []CompositeValue{&interfaceValue{val: "ins$\\ta{nce}-"}, &holeValue{hole: "instance.name"}, &holeValue{hole: "version"}},
			},
			expHoles: map[string][]string{"instance.name": nil, "version": nil},
			expValue: "ins$\\ta{nce}-",
		},
		{
			val: &concatenationValue{
				vals: []CompositeValue{&interfaceValue{val: "ins$\\ta{nce}-"}, &holeValue{hole: "instance.name"}, &holeValue{hole: "version"}},
			},
			holesFillers: map[string]interface{}{"instance.name": "toto", "version": 2},
			expValue:     "ins$\\ta{nce}-toto2",
		},
	}

	for i, tcase := range tcases {
		if withHoles, ok := tcase.val.(WithHoles); ok {
			withHoles.ProcessHoles(tcase.holesFillers)
		}
		if withRefs, ok := tcase.val.(WithRefs); ok {
			withRefs.ProcessRefs(tcase.refsFillers)
		}
		if len(tcase.expHoles) > 0 {
			withHoles, ok := tcase.val.(WithHoles)
			if !ok {
				t.Fatalf("%d: holes: expect value to implement `WithHoles`", i+1)
			}
			if got, want := withHoles.GetHoles(), tcase.expHoles; !reflect.DeepEqual(got, want) {
				t.Fatalf("%d: holes: got %#v, want %#v", i+1, got, want)
			}
		} else if withHoles, ok := tcase.val.(WithHoles); ok && len(withHoles.GetHoles()) > 0 {
			t.Fatalf("%d: holes: expect to have no holes, got %v", i, withHoles.GetHoles())
		}
		if len(tcase.expRefs) > 0 {
			withRefs, ok := tcase.val.(WithRefs)
			if !ok {
				t.Fatalf("%d: refs: expect value to implement `WithRefs`", i+1)
			}
			if got, want := withRefs.GetRefs(), tcase.expRefs; !reflect.DeepEqual(got, want) {
				t.Fatalf("%d: refs: got %#v, want %#v", i+1, got, want)
			}
		} else if withRefs, ok := tcase.val.(WithRefs); ok && len(withRefs.GetRefs()) > 0 {
			t.Fatalf("%d: refs: expect to have no refs, got %v", i, withRefs.GetRefs())
		}

		if len(tcase.expAliases) > 0 {
			aliasVal, ok := tcase.val.(WithAlias)
			if !ok {
				t.Fatalf("%d: aliases: expect value to implement `WithAlias`", i+1)
			}
			if got, want := aliasVal.GetAliases(), tcase.expAliases; !reflect.DeepEqual(got, want) {
				t.Fatalf("%d: aliases: got %#v, want %#v", i+1, got, want)
			}
		} else if withAliases, ok := tcase.val.(WithAlias); ok && len(withAliases.GetAliases()) > 0 {
			t.Fatalf("%d: aliases: expect to have no aliases, got %v", i, withAliases.GetAliases())
		}

		if got, want := tcase.val.Value(), tcase.expValue; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: value: got %#v, want %#v", i+1, got, want)
		}
	}
}

func TestCompositeValuesStringer(t *testing.T) {
	tcases := []struct {
		val    CompositeValue
		expect string
	}{
		{val: &interfaceValue{val: "test"}, expect: "test"},
		{val: &interfaceValue{val: "te\"st"}, expect: "'te\"st'"},
		{val: &interfaceValue{val: "te'st"}, expect: "\"te'st\""},
		{val: &interfaceValue{val: 10}, expect: "10"},
		{val: &interfaceValue{val: "10"}, expect: "'10'"},
		{val: &holeValue{hole: "myhole"}, expect: "{myhole}"},
		{val: &referenceValue{ref: "myref"}, expect: "$myref"},
		{val: &aliasValue{alias: "myalias"}, expect: "@myalias"},
		{
			val: newCompositeValue(
				&interfaceValue{val: "test"},
				&interfaceValue{val: 10},
				&holeValue{hole: "myhole"},
				&referenceValue{ref: "myref"},
				&aliasValue{alias: "myalias"},
			),
			expect: "[test,10,{myhole},$myref,@myalias]",
		},
		{
			val: &concatenationValue{
				vals: []CompositeValue{&interfaceValue{val: "prefix-"}, &holeValue{hole: "hole1"}, &interfaceValue{val: "middle1-"}, &holeValue{hole: "hole2"}, &interfaceValue{val: "-middle2-"}, &holeValue{hole: "hole3"}, &interfaceValue{val: "suffix"}},
			},
			expect: "'prefix-'+{hole1}+'middle1-'+{hole2}+'-middle2-'+{hole3}+'suffix'",
		},
		{val: &holeValue{val: []interface{}{"val1", "val2"}}, expect: "[val1,val2]"},
		{val: &referenceValue{val: []interface{}{"val1", 12}}, expect: "[val1,12]"},
		{
			val: &concatenationValue{
				vals: []CompositeValue{&interfaceValue{val: "ins$\\ta{nce}-"}, &holeValue{hole: "instance.name"}, &holeValue{hole: "version"}},
			},
			expect: "'ins$\\ta{nce}-'+{instance.name}+{version}",
		},
		{
			val: &concatenationValue{
				vals: []CompositeValue{&interfaceValue{val: "ins$\\ta{nce}-"}, &holeValue{hole: "instance.name", val: "toto"}, &holeValue{hole: "version"}},
			},
			expect: "'ins$\\ta{nce}-'+'toto'+{version}",
		},
		{
			val: &concatenationValue{
				vals: []CompositeValue{&interfaceValue{val: "ins$\\ta{nce}-"}, &holeValue{hole: "instance.name", val: "toto"}, &holeValue{hole: "version", val: 2}},
			},
			expect: "'ins$\\ta{nce}-toto2'",
		},
		{
			val: &concatenationValue{
				vals: []CompositeValue{&interfaceValue{val: "instance-"}, &holeValue{hole: "instance.name", val: "toto"}, &holeValue{hole: "version", val: 2}},
			},
			expect: "instance-toto2",
		},
	}

	for i, tcase := range tcases {
		if got, want := tcase.val.String(), tcase.expect; got != want {
			t.Fatalf("%d: got %s, want %s", i+1, got, want)
		}
	}
}

func TestCloneValues(t *testing.T) {
	tcases := []struct {
		from       CompositeValue
		mutationFn func(CompositeValue)
	}{
		{from: &interfaceValue{val: "test"}, mutationFn: func(v CompositeValue) { v.(*interfaceValue).val = "other" }},
		{from: &holeValue{hole: "myhole"}, mutationFn: func(v CompositeValue) { v.(*holeValue).ProcessHoles(map[string]interface{}{"myhole": "myvalue"}) }},
		{from: &referenceValue{ref: "myref"}, mutationFn: func(v CompositeValue) { v.(*referenceValue).ProcessRefs(map[string]interface{}{"myref": "myvalue"}) }},
		{from: &aliasValue{alias: "myalias"}, mutationFn: func(v CompositeValue) {
			v.(*aliasValue).ResolveAlias(func(s string) (string, bool) {
				if s == "myalias" {
					return "myvalue", true
				}
				return "", false
			})
		}},
		{
			from: newCompositeValue(
				&interfaceValue{val: "test"},
				&interfaceValue{val: 10},
				&holeValue{hole: "myhole"},
				&referenceValue{ref: "myref"},
				&aliasValue{alias: "myalias"},
			),
			mutationFn: func(v CompositeValue) { v.(*listValue).ProcessHoles(map[string]interface{}{"myhole": "myvalue"}) },
		},
		{
			from: &concatenationValue{
				vals: []CompositeValue{&interfaceValue{val: "prefix-"}, &holeValue{hole: "hole1"}, &interfaceValue{val: "middle1-"}, &holeValue{hole: "hole2"}, &interfaceValue{val: "-middle2-"}, &holeValue{hole: "hole3"}, &interfaceValue{val: "suffix"}},
			},
			mutationFn: func(v CompositeValue) {
				v.(*concatenationValue).ProcessHoles(map[string]interface{}{"hole1": "myvalue"})
			},
		},
	}
	for i, tcase := range tcases {
		clone := tcase.from.Clone()
		tcase.mutationFn(tcase.from)
		if reflect.DeepEqual(clone.Value(), tcase.from.Value()) {
			t.Fatalf("%d: expect original and clone values to be different, got %#v", i, clone.Value())
		}
		tcase.mutationFn(clone)
		if !reflect.DeepEqual(clone.Value(), tcase.from.Value()) {
			t.Fatalf("%d: expect original and clone values to have the same value, got %#v and %#v", i, clone.Value(), tcase.from.Value())
		}
	}
}

func newCompositeValue(values ...CompositeValue) CompositeValue {
	return &listValue{vals: values}
}
