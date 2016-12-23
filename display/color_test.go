package display

import (
	"testing"

	"github.com/fatih/color"
)

func TestStringToColor(t *testing.T) {
	data := map[string]color.Attribute{
		"red":    color.FgRed,
		"yellow": color.FgYellow,
		"blue":   color.FgBlue,
		"green":  color.FgGreen,
		"cyan":   color.FgCyan,
		"white":  color.FgWhite,
		"black":  color.FgBlack,
	}
	for k, expected := range data {
		if got, want := stringToColor(k), expected; got != want {
			t.Errorf("got %v; want %v\n", got, want)
		}
	}

}
