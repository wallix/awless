package fixedwidth

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"encoding"
)

func ExampleUnmarshal() {
	// define the format
	var people []struct {
		ID        int     `fixed:"1,5"`
		FirstName string  `fixed:"6,15"`
		LastName  string  `fixed:"16,25"`
		Grade     float64 `fixed:"26,30"`
	}

	// define some fixed-with data to parse
	data := []byte("" +
		"1    Ian       Lopshire  99.50" + "\n" +
		"2    John      Doe       89.50" + "\n" +
		"3    Jane      Doe       79.50" + "\n")

	err := Unmarshal(data, &people)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", people[0])
	fmt.Printf("%+v\n", people[1])
	fmt.Printf("%+v\n", people[2])
	// Output:
	//{ID:1 FirstName:Ian LastName:Lopshire Grade:99.5}
	//{ID:2 FirstName:John LastName:Doe Grade:89.5}
	//{ID:3 FirstName:Jane LastName:Doe Grade:79.5}
}

func TestUnmarshal(t *testing.T) {
	// allTypes contains a field with all current supported types.
	type allTypes struct {
		String          string          `fixed:"1,5"`
		Int             int             `fixed:"6,10"`
		Float           float64         `fixed:"11,15"`
		TextUnmarshaler EncodableString `fixed:"16,20"`
	}
	for _, tt := range []struct {
		name      string
		rawValue  []byte
		target    interface{}
		expected  interface{}
		shouldErr bool
	}{
		{
			name:     "Basic Slice Case",
			rawValue: []byte("foo  123  1.2  bar" + "\n" + "bar  321  2.1  foo"),
			target:   &[]allTypes{},
			expected: &[]allTypes{
				{"foo", 123, 1.2, EncodableString{"bar", nil}},
				{"bar", 321, 2.1, EncodableString{"foo", nil}},
			},
			shouldErr: false,
		},
		{
			name:      "Basic Struct Case",
			rawValue:  []byte("foo  123  1.2  bar"),
			target:    &allTypes{},
			expected:  &allTypes{"foo", 123, 1.2, EncodableString{"bar", nil}},
			shouldErr: false,
		},
		{
			name:      "Unmarshal Error",
			rawValue:  []byte("foo  nan  ddd  bar"),
			target:    &allTypes{},
			expected:  &allTypes{},
			shouldErr: true,
		},
		{
			name:      "Empty Line",
			rawValue:  []byte(""),
			target:    &allTypes{},
			expected:  &allTypes{"", 0, 0, EncodableString{"", nil}},
			shouldErr: false,
		},
		{
			name:      "Invalid Target",
			rawValue:  []byte("foo  123  1.2  bar"),
			target:    allTypes{},
			expected:  allTypes{},
			shouldErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal(tt.rawValue, tt.target)
			if tt.shouldErr != (err != nil) {
				t.Errorf("Unmarshal() err want %v, have %v (%v)", tt.shouldErr, err != nil, err.Error())
			}
			if !tt.shouldErr && !reflect.DeepEqual(tt.target, tt.expected) {
				t.Errorf("Unmarshal() want %+v, have %+v", tt.expected, tt.target)
			}

		})
	}

	t.Run("Invalid Unmarshal Errors", func(t *testing.T) {
		for _, tt := range []struct {
			name      string
			v         interface{}
			shouldErr bool
		}{
			{"Invalid Unmarshal Nil", nil, true},
			{"Invalid Unmarshal Not Pointer 1", struct{}{}, true},
			{"Invalid Unmarshal Not Pointer 2", []struct{}{}, true},
			{"Valid Unmarshal slice", &[]struct{}{}, false},
			{"Valid Unmarshal struct", &struct{}{}, false},
		} {
			t.Run(tt.name, func(t *testing.T) {
				err := Unmarshal([]byte{}, tt.v)
				if tt.shouldErr != (err != nil) {
					t.Errorf("Unmarshal() err want %v, have %v (%v)", tt.shouldErr, err != nil, err)
				}
			})
		}
	})
}

func TestNewValueSetter(t *testing.T) {
	for _, tt := range []struct {
		name      string
		raw       []byte
		expected  interface{}
		shouldErr bool
	}{
		{"invalid type", []byte("foo"), true, true},

		{"textUnmarshaler implementation", []byte("foo"), &EncodableString{"foo", nil}, false},
		{"textUnmarshaler implementation if addressed", []byte("foo"), EncodableString{"foo", nil}, false},
		{"textUnmarshaler implementation as interface", []byte("foo"), encoding.TextUnmarshaler(&EncodableString{"foo", nil}), false},
		{"textUnmarshaler implementation in interface", []byte("foo"), interface{}(&EncodableString{"foo", nil}), false},
		{"textUnmarshaler implementation if addressed in interface", []byte("foo"), interface{}(EncodableString{"foo", nil}), false},

		{"string", []byte("foo"), string("foo"), false},
		{"string empty", []byte(""), string(""), false},
		{"string interface", []byte("foo"), interface{}(string("foo")), false},
		{"string interface empty", []byte(""), interface{}(string("")), false},
		{"*string", []byte("foo"), stringp("foo"), false},
		{"*string empty", []byte(""), (*string)(nil), false},

		{"int", []byte("1"), int(1), false},
		{"int zero", []byte("0"), int(0), false},
		{"int empty", []byte(""), int(0), false},
		{"*int", []byte("1"), intp(1), false},
		{"*int zero", []byte("0"), intp(0), false},
		{"*int empty", []byte(""), (*int)(nil), false},
		{"int Invalid", []byte("foo"), int(0), true},

		{"float64", []byte("1.23"), float64(1.23), false},
		{"*float64", []byte("1.23"), float64p(1.23), false},
		{"*float64 zero", []byte("0"), float64p(0), false},
		{"*float64 empty", []byte(""), (*float64)(nil), false},
		{"float64 Invalid", []byte("foo"), float64(0), true},

		{"float32", []byte("1.23"), float32(1.23), false},
		{"float32 Invalid", []byte("foo"), float32(0), true},

		{"int8", []byte("1"), int8(1), false},
		{"int16", []byte("1"), int16(1), false},
		{"int32", []byte("1"), int32(1), false},
		{"int64", []byte("1"), int64(1), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// ensure we have an addressable target
			var i = reflect.Indirect(reflect.New(reflect.TypeOf(tt.expected)))

			err := newValueSetter(i.Type())(i, tt.raw)
			if tt.shouldErr != (err != nil) {
				t.Errorf("newValueSetter(%s)() err want %v, have %v (%v)", reflect.TypeOf(tt.expected).Name(), tt.shouldErr, err != nil, err.Error())
			}
			if !tt.shouldErr && !reflect.DeepEqual(tt.expected, i.Interface()) {
				t.Errorf("newValueSetter(%s)() want %s, have %s", reflect.TypeOf(tt.expected).Name(), tt.expected, i)
			}
		})
	}
}
