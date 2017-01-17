package shell

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func GetTerminalWidth() (int, error) {
	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, err
	}
	return w, nil
}
