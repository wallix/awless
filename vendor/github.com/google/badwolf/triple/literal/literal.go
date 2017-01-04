// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package literal provides an abstraction to manipulate BadWolf literals.
package literal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/pborman/uuid"
)

// Type represents the type contained in a literal.
type Type uint8

const (
	// Bool indicates that the type contained in the literal is a bool.
	Bool Type = iota
	// Int64 indicates that the type contained in the literal is an int64.
	Int64
	// Float64 indicates that the type contained in the literal is a float64.
	Float64
	// Text indicates that the type contained in the literal is a string.
	Text
	// Blob indicates that the type contained in the literal is a []byte.
	Blob
)

// Strings returns the pretty printing version of the type
func (t Type) String() string {
	switch t {
	case Bool:
		return "bool"
	case Int64:
		return "int64"
	case Float64:
		return "float64"
	case Text:
		return "text"
	case Blob:
		return "blob"
	default:
		return "UNKNOWN"
	}
}

// Literal represents the type and value boxed in the literal.
type Literal struct {
	t Type
	v interface{}
}

// Type returns the type of a literal.
func (l *Literal) Type() Type {
	return l.t
}

// String returns a string representation of the literal.
func (l *Literal) String() string {
	return fmt.Sprintf("\"%v\"^^type:%v", l.Interface(), l.Type())
}

// ToComparableString returns a string that can be directly compared.
func (l *Literal) ToComparableString() string {
	s := ""
	switch l.t {
	case Int64:
		s = fmt.Sprintf("\"%032d\"^^type:%v", l.Interface(), l.Type())
	case Float64:
		s = fmt.Sprintf("\"%032f\"^^type:%v", l.Interface(), l.Type())
	default:
		s = l.String()
	}
	return s
}

// Bool returns the value of a literal as a boolean.
func (l *Literal) Bool() (bool, error) {
	if l.t != Bool {
		return false, fmt.Errorf("literal.Bool: literal is of type %v; cannot be converted to a bool", l.t)
	}
	return l.v.(bool), nil
}

// Int64 returns the value of a literal as an int64.
func (l *Literal) Int64() (int64, error) {
	if l.t != Int64 {
		return 0, fmt.Errorf("literal.Int64: literal is of type %v; cannot be converted to a int64", l.t)
	}
	return l.v.(int64), nil
}

// Float64 returns the value of a literal as a float64.
func (l *Literal) Float64() (float64, error) {
	if l.t != Float64 {
		return 0, fmt.Errorf("literal.Float64: literal is of type %v; cannot be converted to a float64", l.t)
	}
	return l.v.(float64), nil
}

// Text returns the value of a literal as a string.
func (l *Literal) Text() (string, error) {
	if l.t != Text {
		return "", fmt.Errorf("literal.Text: literal is of type %v; cannot be converted to a string", l.t)
	}
	return l.v.(string), nil
}

// Blob returns the value of a literal as a []byte.
func (l *Literal) Blob() ([]byte, error) {
	if l.t != Blob {
		return nil, fmt.Errorf("literal.Blob: literal is of type %v; cannot be converted to a []byte", l.t)
	}
	return l.v.([]byte), nil
}

// Interface returns the value as a simple interface{}.
func (l *Literal) Interface() interface{} {
	return l.v
}

// Builder interface provides a standard way to build literals given a type and
// a given value.
type Builder interface {
	Build(t Type, v interface{}) (*Literal, error)
	Parse(s string) (*Literal, error)
}

// A singleton used to build all literals.
var defaultBuilder Builder

func init() {
	defaultBuilder = &unboundBuilder{}
}

// The default builder is unbound. This allows to create a literal arbitrarily
// long.
type unboundBuilder struct{}

// Build creates a new unbound literal from a type and a value.
func (b *unboundBuilder) Build(t Type, v interface{}) (*Literal, error) {
	switch v.(type) {
	case bool:
		if t != Bool {
			return nil, fmt.Errorf("literal.Build: type %v does not match type of value %v", t, v)
		}
	case int64:
		if t != Int64 {
			return nil, fmt.Errorf("literal.Build: type %v does not match type of value %v", t, v)
		}
	case float64:
		if t != Float64 {
			return nil, fmt.Errorf("literal.Build: type %v does not match type of value %v", t, v)
		}
	case string:
		if t != Text {
			return nil, fmt.Errorf("literal.Build: type %v does not match type of value %v", t, v)
		}
	case []byte:
		if t != Blob {
			return nil, fmt.Errorf("literal.Build: type %v does not match type of value %v", t, v)
		}
	default:
		return nil, fmt.Errorf("literal.Build: type %T is not supported when building literals", v)
	}
	return &Literal{
		t: t,
		v: v,
	}, nil
}

// Parse creates a string out of a prettified representation.
func (b *unboundBuilder) Parse(s string) (*Literal, error) {
	raw := strings.TrimSpace(s)
	if len(raw) == 0 {
		return nil, fmt.Errorf("literal.Parse: cannot parse and empty string into a literal; provided string %q", s)
	}
	if raw[0] != '"' {
		return nil, fmt.Errorf("literal.Parse: text encoded literals must start with \", missing in %s", raw)
	}
	idx := strings.Index(raw, "\"^^type:")
	if idx < 0 {
		return nil, fmt.Errorf("literal.Parse: text encoded literals must have a type; missing in %s", raw)
	}
	v := raw[1:idx]
	t := raw[idx+len("\"^^type:"):]
	switch t {
	case "bool":
		pv, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("literal.Parse: could not convert value %q to bool", v)
		}
		return b.Build(Bool, pv)
	case "int64":
		pv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("literal.Parse: could not convert value %q to int64", v)
		}
		return b.Build(Int64, int64(pv))
	case "float64":
		pv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("literal.Parse: could not convert value %q to float64", v)
		}
		return b.Build(Float64, float64(pv))
	case "text":
		return b.Build(Text, v)
	case "blob":
		values := v[1 : len(v)-1]
		if values == "" {
			return b.Build(Blob, []byte{})
		}
		bs := []byte{}
		for _, s := range strings.Split(values, " ") {
			b, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("literal.Parse: failed to decode byte array on %q with error %v", s, err)
			}
			bs = append(bs, byte(b))
		}
		return b.Build(Blob, bs)
	default:
		return nil, nil
	}
}

// DefaultBuilder returns a builder with no constraints or checks.
func DefaultBuilder() Builder {
	return defaultBuilder
}

// boundedBuilder implements a literal builder where strings and blobs are
// guaranteed of being of bounded size
type boundedBuilder struct {
	max int
}

// Build creates a new literal of bounded size.
func (b *boundedBuilder) Build(t Type, v interface{}) (*Literal, error) {
	switch v.(type) {
	case string:
		if l := len(v.(string)); l > b.max {
			return nil, fmt.Errorf("literal.Build: cannot create literal due to size of %v (%d>%d)", v, l, b.max)
		}
	case []byte:
		if l := len(v.([]byte)); l > b.max {
			return nil, fmt.Errorf("literal.Build: cannot create literal due to size of %v (%d>%d)", v, l, b.max)
		}
	}
	return defaultBuilder.Build(t, v)
}

// Parse creates a string out of a prettyfied representation.
func (b *boundedBuilder) Parse(s string) (*Literal, error) {
	l, err := defaultBuilder.Parse(s)
	if err != nil {
		return nil, err
	}
	t := l.Type()
	switch t {
	case Text:
		if text, err := l.Text(); err != nil || len(text) > b.max {
			return nil, fmt.Errorf("literal.Parse: cannot create literal due to size of %v (%d>%d)", t, len(text), b.max)
		}
	case Blob:
		if blob, err := l.Blob(); err != nil || len(blob) > b.max {
			return nil, fmt.Errorf("literal.Parse: cannot create literal due to size of %v (%d>%d)", t, len(blob), b.max)
		}
	}
	return l, nil
}

// NewBoundedBuilder creates a builder that guarantees that no literal will
// be created if the size of the string or a blob is bigger than the provided
// maximum.
func NewBoundedBuilder(max int) Builder {
	return &boundedBuilder{max: max}
}

// UUID returns a global unique identifier for the given literal. It is
// implemented as the SHA1 UUID of the literal value.
func (l *Literal) UUID() uuid.UUID {
	var buffer bytes.Buffer

	switch v := l.v.(type) {
	case bool:
		if v {
			buffer.WriteString("true")
		} else {
			buffer.WriteString("false")
		}
	case int64:
		b := make([]byte, 8)
		binary.PutVarint(b, v)
		buffer.Write(b)
	case float64:
		bs := math.Float64bits(v)
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, bs)
		buffer.Write(b)
	case string:
		buffer.Write([]byte(v))
	case []byte:
		buffer.Write(v)
	}

	return uuid.NewSHA1(uuid.NIL, buffer.Bytes())
}
