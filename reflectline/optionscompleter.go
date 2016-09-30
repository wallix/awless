package reflectline

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/chzyer/readline"
)

// OptionsCompleter is an AutoCompleter for options
type OptionsCompleter struct {
	readline.AutoCompleter
	options map[string]reflect.Kind
}

// NewOptionsCompleter creates a new AutoCompleter for the given input ty
func NewOptionsCompleter(ty reflect.Type) *OptionsCompleter {
	oc := &OptionsCompleter{
		options: make(map[string]reflect.Kind),
	}
	oc.filloptions("", ty)
	return oc
}

func (oc *OptionsCompleter) filloptions(name string, ty reflect.Type) {
	switch ty.Kind() {
	case reflect.Struct:
		n := ty.NumField()
		for i := 0; i < n; i++ {
			sf := ty.Field(i)
			p := fieldToOption(sf.Name)
			if name != "" {
				p = fmt.Sprintf("%s.%s", name, p)
			}
			oc.filloptions(p, sf.Type)
		}
		return
	case reflect.Bool:
		name = fmt.Sprintf("-%s", name)
	default:
		name = fmt.Sprintf("-%s=", name)
	}
	oc.options[name] = ty.Kind()
}

// Do is a part of AutoCompleter
func (oc *OptionsCompleter) Do(line []rune, pos int) ([][]rune, int) {
	str := string(line)
	if len(line) >= 1 {
		for i := len(line) - 1; i >= 0; i-- {
			if unicode.IsSpace(line[i]) {
				str = string(line[i+1 : len(line)])
				break
			}
		}
	}
	length := len(str)
	newLine := make([][]rune, 0, len(oc.options))
	for k := range oc.options {
		if len(k) > length && k[:length] == str {
			newLine = append(newLine, []rune(k[length:]))
		}
	}
	return newLine, length
}
