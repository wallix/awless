package console

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func GetTerminalWidth() int {
	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0
	}
	return w
}
