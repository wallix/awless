package main

import (
	"fmt"

	"github.com/wallix/awless/reflectline"
)

type Test struct {
	Multi struct {
		Aaa int
		Bbb int
	}
	Multpass string
}

func Command(f interface{}) reflectline.ShellNode {
	cmd, err := reflectline.NewShellCommand(f)
	if err != nil {
		panic("Cannot create shell command: " + err.Error())
	}
	return cmd
}

func main() {
	// primary group
	pnode := reflectline.NewShellGroup()

	// user group
	{
		unode := reflectline.NewShellGroup()
		unode.AddNode("add", Command(func(t Test) {
			fmt.Printf("ADD:: %v\n", t)
		}))
		unode.AddNode("delete", Command(func(t Test) {
			fmt.Printf("DELETE:: %v\n", t)
		}))
		pnode.AddNode("users", unode)
	}

	shell := reflectline.NewShell(pnode)
	shell.Run()
}
