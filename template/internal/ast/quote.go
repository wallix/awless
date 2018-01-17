package ast

import (
	"regexp"
	"strconv"
	"strings"
)

var SimpleStringValue = regexp.MustCompile("^[a-zA-Z0-9-._:/+;~@<>*]+$") // in sync with [a-zA-Z0-9-._:/+;~@<>]+ in PEG (with ^ and $ around)

func quoteStringIfNeeded(input string) string {
	if _, err := strconv.Atoi(input); err == nil {
		return "'" + input + "'"
	}
	if _, err := strconv.ParseFloat(input, 64); err == nil {
		return "'" + input + "'"
	}
	if SimpleStringValue.MatchString(input) {
		return input
	} else {
		return Quote(input)
	}
}

func Quote(str string) string {
	if strings.ContainsRune(str, '\'') {
		return "\"" + str + "\""
	} else {
		return "'" + str + "'"
	}
}

func isQuoted(str string) bool {
	if len(str) < 2 {
		return false
	}
	if str[0] == '\'' && str[len(str)-1] == '\'' {
		return true
	}
	if str[0] == '"' && str[len(str)-1] == '"' {
		return true
	}
	return false
}
