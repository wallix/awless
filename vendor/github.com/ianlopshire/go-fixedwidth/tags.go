package fixedwidth

import (
	"strconv"
	"strings"
)

// parseTag splits a struct fields fixed tag into its start and end positions.
// If the tag is not valid, ok will be false.
func parseTag(tag string) (startPos, endPos int, ok bool) {
	parts := strings.Split(tag, ",")
	if len(parts) != 2 {
		return startPos, endPos, false
	}

	var err error
	if startPos, err = strconv.Atoi(parts[0]); err != nil {
		return startPos, endPos, false
	}
	if endPos, err = strconv.Atoi(parts[1]); err != nil {
		return startPos, endPos, false
	}
	if startPos > endPos || (startPos == 0 && endPos == 0) {
		return startPos, endPos, false
	}

	return startPos, endPos, true
}
