package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wallix/awless/cmd"
)

func main() {
	createAwlessDefaultDir()

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func createAwlessDefaultDir() {
	dir := filepath.Join(os.Getenv("HOME"), ".awless")
	os.MkdirAll(dir, 0700)
}
