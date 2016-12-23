package display

import (
	"strings"

	"github.com/fatih/color"
)

func stringToColor(str string) color.Attribute {
	switch strings.ToLower(str) {
	case "red":
		return color.FgRed
	case "yellow":
		return color.FgYellow
	case "blue":
		return color.FgBlue
	case "green":
		return color.FgGreen
	case "cyan":
		return color.FgCyan
	case "white":
		return color.FgWhite
	default:
		return color.FgBlack
	}
}

func colorDisplay(str string, coloredValues map[string]string) string {
	col := coloredValues[str]
	if col != "" {
		return color.New(stringToColor(col)).SprintFunc()(str)
	}
	return str
}
