package reflectline

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/chzyer/readline"
)

type ShellNode interface {
	readline.AutoCompleter
	Call(line []string) error
}

// ShellGroup is a set of ShellNode
type ShellGroup struct {
	childs map[string]ShellNode
}

func NewShellGroup() *ShellGroup {
	return &ShellGroup{
		childs: make(map[string]ShellNode),
	}
}

// NewReflectGroup creates a new ShellNode by reflection.
func NewReflectGroup(i interface{}) (*ShellGroup, error) {
	g := NewShellGroup()
	ty := reflect.TypeOf(i)
	vl := reflect.ValueOf(i)
	n := ty.NumMethod()
	for i := 0; i < n; i++ {
		vm := vl.Method(i)
		m := ty.Method(i)
		c, err := NewShellCommand(vm.Interface())
		if err != nil {
			return nil, err
		}
		g.AddNode(fieldToOption(m.Name), c)
	}
	return g, nil
}

// MustReflectGroup as NewReflectGroup but panics in case of error
func MustReflectGroup(i interface{}) *ShellGroup {
	g, e := NewReflectGroup(i)
	if e != nil {
		panic(e)
	}
	return g
}

// AddNode adds a new node to the group
func (s *ShellGroup) AddNode(name string, node ShellNode) {
	s.childs[name] = node
}

// Do is part of AutoCompleter interface
func (s *ShellGroup) Do(line []rune, pos int) ([][]rune, int) {
	if i := strings.IndexFunc(string(line), unicode.IsSpace); i != -1 {
		name := string(line[:i])
		c, ok := s.childs[name]
		if ok {
			return c.Do([]rune(strings.TrimLeftFunc(string(line[i:]), unicode.IsSpace)), pos)
		}
		return [][]rune{}, 0
	}
	newline := make([][]rune, 0, len(s.childs))
	length := len(line)
	str := string(line)
	for k := range s.childs {
		if len(k) > length && k[:length] == str {
			newline = append(newline, []rune(k[length:]))
		}
	}
	return newline, length
}

func (s *ShellGroup) Call(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("command is incomplete")
	}
	name := args[0]
	c, ok := s.childs[name]
	if !ok {
		return fmt.Errorf("command %s not found", name)
	}
	return c.Call(args[1:])
}
