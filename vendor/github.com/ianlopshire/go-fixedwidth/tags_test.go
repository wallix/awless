package fixedwidth

import (
	"testing"
)

func TestParseTag(t *testing.T) {
	for _, tt := range []struct {
		name     string
		tag      string
		startPos int
		endPos   int
		ok       bool
	}{
		{"Valid Tag", "0,10", 0, 10, true},
		{"Valid Tag Single position", "5,5", 5, 5, true},
		{"Tag Empty", "", 0, 0, false},
		{"Tag Too short", "0", 0, 0, false},
		{"Tag Too Long", "2,10,11", 0, 0, false},
		{"StartPos Not Integer", "hello,3", 0, 0, false},
		{"EndPos Not Integer", "3,hello", 0, 0, false},
		{"Tag Contains a Space", "4, 11", 0, 0, false},
		{"Tag Interval Invalid", "14,5", 0, 0, false},
		{"Tag Both Positions Zero", "0,0", 0, 0, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			startPos, endPos, ok := parseTag(tt.tag)
			if tt.ok != ok {
				t.Errorf("parseTag() ok want %v, have %v", tt.ok, ok)
			}

			// only check startPos and endPos if valid tags are expected
			if tt.ok {
				if tt.startPos != startPos {
					t.Errorf("parseTag() startPos want %v, have %v", tt.startPos, startPos)
				}

				if tt.endPos != endPos {
					t.Errorf("parseTag() endPos want %v, have %v", tt.endPos, endPos)
				}
			}
		})
	}
}
