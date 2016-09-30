package reflectline

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func fieldToOption(field string) string {
	if len(field) <= 0 {
		return field
	}
	var writer bytes.Buffer
	reader := strings.NewReader(field)
	r, _, _ := reader.ReadRune()
	writer.WriteRune(unicode.ToLower(r))
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if unicode.IsUpper(r) {
			writer.WriteRune('-')
			writer.WriteRune(unicode.ToLower(r))
		} else {
			writer.WriteRune(r)
		}
	}
	return writer.String()
}

func optionToField(option string) string {
	if len(option) <= 0 {
		return option
	}
	var writer bytes.Buffer
	reader := strings.NewReader(option)
	r, _, _ := reader.ReadRune()
	writer.WriteRune(unicode.ToUpper(r))
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if r == '-' {
			r1, _, err := reader.ReadRune()
			if err == io.EOF {
				writer.WriteRune(r)
				break
			}
			writer.WriteRune(unicode.ToUpper(r1))
		} else {
			writer.WriteRune(r)
		}
	}
	return writer.String()
}

func indexEndOfDoubleQuote(line []rune) int {
	if len(line) == 0 {
		os.Exit(1)
	}
	for i, r := range line {
		if r == '"' && (i == 0 || line[i-1] != '\\') {
			fmt.Printf("eod line(%s) i=%d\n", string(line), i)
			return i
		}
	}
	return -1
}
