package fixedwidth

import (
	"bufio"
	"bytes"
	"encoding"
	"errors"
	"io"
	"reflect"
	"strconv"
)

// Unmarshal parses fixed width encoded data and stores the
// result in the value pointed to by v. If v is nil or not a
// pointer, Unmarshal returns an InvalidUnmarshalError.
func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(data)).Decode(v)
}

// A Decoder reads and decodes fixed width data from an input stream.
type Decoder struct {
	data *bufio.Reader
	done bool
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		data: bufio.NewReader(r),
	}
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "fixedwidth: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "fixedwidth: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "fixedwidth: Unmarshal(nil " + e.Type.String() + ")"
}

// An UnmarshalTypeError describes a value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value  string       // the raw value
	Type   reflect.Type // type of Go value it could not be assigned to
	Struct string       // name of the struct type containing the field
	Field  string       // name of the field holding the Go value
	Cause error // original error
}

func (e *UnmarshalTypeError) Error() string {
	var s string
	if e.Struct != "" || e.Field != "" {
		s =  "fixedwidth: cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
	} else {
		s = "fixedwidth: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
	}
	if e.Cause != nil {
		return s + ":" + e.Cause.Error()
	}
	return s
}

// Decode reads from its input and stores the decoded data to the value
// pointed to by v.
//
// In the case that v points to a struct value, Decode will read a
// single line from the input.
//
// In the case that v points to a slice value, Decode will read until
// the end of its input.
func (d *Decoder) Decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	if reflect.Indirect(reflect.ValueOf(v)).Kind() == reflect.Slice {
		return d.readLines(reflect.ValueOf(v).Elem())
	}
	return d.readLine(reflect.ValueOf(v))
}

func (d *Decoder) readLines(v reflect.Value) (err error) {
	ct := v.Type().Elem()
	for {
		nv := reflect.New(ct).Elem()
		err := d.readLine(nv)
		if err != nil {
			return err
		}
		if d.done {
			break
		}
		v.Set(reflect.Append(v, nv))
	}
	return nil
}

func (d *Decoder) readLine(v reflect.Value) (err error) {
	// TODO: properly handle prefixed lines
	line, _, err := d.data.ReadLine()
	if err == io.EOF {
		d.done = true
		return nil
	} else if err != nil {
		return err
	}

	return newValueSetter(v.Type())(v, line)
}

func rawValueFromLine(line []byte, startPos, endPos int) []byte {
	if len(line) == 0 || startPos >= len(line) {
		return []byte{}
	}
	if endPos > len(line) {
		endPos = len(line)
	}
	return bytes.TrimSpace(line[startPos-1:endPos])
}

type valueSetter func(v reflect.Value, raw []byte) error

var textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

func newValueSetter(t reflect.Type) valueSetter {
	if t.Implements(textUnmarshalerType) {
		return textUnmarshalerSetter(t, false)
	}
	if reflect.PtrTo(t).Implements(textUnmarshalerType) {
		return textUnmarshalerSetter(t, true)
	}

	switch t.Kind() {
	case reflect.Ptr:
		return ptrSetter(t)
	case reflect.Interface:
		return interfaceSetter
	case reflect.Struct:
		return structSetter
	case reflect.String:
		return stringSetter
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return intSetter
	case reflect.Float32:
		return floatSetter(32)
	case reflect.Float64:
		return floatSetter(64)
	}
	return unknownSetter
}

func structSetter (v reflect.Value, raw []byte) error {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		if !fv.IsValid() {
			continue
		}
		sf := t.Field(i)
		startPos, endPos, ok := parseTag(sf.Tag.Get("fixed"))
		if !ok {
			continue
		}
		rawValue := rawValueFromLine(raw, startPos, endPos)
		err := newValueSetter(sf.Type)(fv, rawValue)
		if err != nil {
			return &UnmarshalTypeError{string(rawValue), sf.Type, t.Name(), sf.Name, err}
		}
	}
	return nil
}

func unknownSetter(v reflect.Value, raw []byte) error {
	return errors.New("fixedwidth: unknown type")
}

func nilSetter(v reflect.Value, _ []byte) error {
	v.Set(reflect.Zero(v.Type()))
	return nil
}

func textUnmarshalerSetter(t reflect.Type, shouldAddr bool) valueSetter {
	return func(v reflect.Value, raw []byte) error {
		if shouldAddr {
			v = v.Addr()
		}
		// set to zero value if this is nil
		if t.Kind() == reflect.Ptr && v.IsNil(){
			v.Set(reflect.New(t.Elem()))
		}
		return v.Interface().(encoding.TextUnmarshaler).UnmarshalText(raw)
	}
}

func interfaceSetter(v reflect.Value, raw []byte) error {
	return newValueSetter(v.Elem().Type())(v.Elem(), raw)
}

func ptrSetter(t reflect.Type) valueSetter {
	return func(v reflect.Value, raw []byte) error {
		if len(raw) <= 0 {
			return nilSetter(v, raw)
		}
		if v.IsNil() {
			v.Set(reflect.New(t.Elem()))
		}
		return newValueSetter(v.Elem().Type())(reflect.Indirect(v), raw)
	}
}

func stringSetter(v reflect.Value, raw []byte) error {
	v.SetString(string(raw))
	return nil
}

func intSetter(v reflect.Value, raw []byte) error {
	if len(raw) < 1 {
		return nil
	}
	i, err := strconv.Atoi(string(raw))
	if err != nil {
		return err
	}
	v.SetInt(int64(i))
	return nil
}

func floatSetter(bitSize int) valueSetter {
	return func(v reflect.Value, raw []byte) error {
		if len(raw) < 1 {
			return nil
		}
		f, err := strconv.ParseFloat(string(raw), bitSize)
		if err != nil {
			return err
		}
		v.SetFloat(f)
		return nil
	}
}
