package reflectline

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

// ShellCommand is a structure that describe a command
type ShellCommand struct {
	completer readline.AutoCompleter
	inputTy   reflect.Type
	fn        interface{}
}

// NewShellCommand creates a new ShellNode that describe a command
func NewShellCommand(fn interface{}) (ShellNode, error) {
	sc := &ShellCommand{fn: fn}
	ty := reflect.TypeOf(fn)
	switch ty.Kind() {
	case reflect.Func:
		switch ty.NumIn() {
		case 0:
		case 1:
			sc.inputTy = ty.In(0)
			sc.completer = NewOptionsCompleter(sc.inputTy)
		default:
			return nil, fmt.Errorf("The fn argument MUST have only one argument")
		}
		return sc, nil
	default:
		return nil, fmt.Errorf("The fn argument MUST be a function")
	}
}

// MustShellCommand as NewShellCommand but panics if an error occurs
func MustShellCommand(fn interface{}) ShellNode {
	cmd, err := NewShellCommand(fn)
	if err != nil {
		panic("Cannot create shell command: " + err.Error())
	}
	return cmd
}

func (s *ShellCommand) buildInput(args []string) (reflect.Value, error) {
	input := reflect.New(s.inputTy).Elem()
	for _, arg := range args {
		if arg[0] != '-' {
			return reflect.ValueOf(0), fmt.Errorf("'%s' is not an option", arg)
		}
		i := strings.Index(arg, "=")
		val := ""
		if i != -1 {
			val = arg[i+1 : len(arg)]
		} else {
			i = len(arg)
		}
		flds := strings.Split(arg[1:i], ".")
		if err := fillFlds(input, s.inputTy, flds, val); err != nil {
			return reflect.ValueOf(0), err
		}
	}
	return input, nil
}

func fillFlds(val reflect.Value, typ reflect.Type, fields []string, str string) error {
	if len(fields) == 0 {
		v, err := parseInput(typ, str)
		if err != nil {
			return err
		}
		val.Set(v)
		return nil
	}
	switch typ.Kind() {
	case reflect.Struct:
		fname := optionToField(fields[0])
		sf, b := typ.FieldByName(fname)
		if b {
			return fillFlds(val.FieldByName(fname), sf.Type, fields[1:], str)
		}
	}
	return fmt.Errorf("option(%s) is not expected", fields[0])
}

func parseInput(typ reflect.Type, str string) (reflect.Value, error) {
	switch typ.Kind() {
	case reflect.Bool:
		if str == "" {
			return reflect.ValueOf(true), nil
		}
		r, err := strconv.ParseBool(str)
		return reflect.ValueOf(r), err
	case reflect.Float64:
		r, err := strconv.ParseFloat(str, 64)
		return reflect.ValueOf(r), err
	case reflect.Int:
		r, err := strconv.ParseInt(str, 10, 0)
		return reflect.ValueOf(int(r)), err
	case reflect.String:
		return reflect.ValueOf(str), nil
	}
	return reflect.ValueOf(struct{}{}), fmt.Errorf("Cannot parse input %s", typ)
}

// Call build the input parameter of its function and call it
func (s *ShellCommand) Call(args []string) error {
	v, err := s.buildInput(args)
	if err != nil {
		return err
	}
	reflect.ValueOf(s.fn).Call([]reflect.Value{v})
	return nil
}

// Do is part of AutoCompleter interface
func (s *ShellCommand) Do(line []rune, pos int) ([][]rune, int) {
	return s.completer.Do(line, pos)
}
