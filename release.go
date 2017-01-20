// +build ignore

package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var builds = map[string][]string{
	"darwin":  []string{"amd64"},
	"linux":   []string{"386", "amd64"},
	"windows": []string{"386", "amd64"},
}

func main() {
	var wg sync.WaitGroup

	for osname, archs := range builds {
		for _, arch := range archs {
			wg.Add(1)
			go func(o, a string) {
				defer wg.Done()
				if err := buildAndZip(o, a); err != nil {
					fmt.Fprintln(os.Stderr, "%s", err)
					return
				}
			}(osname, arch)
		}
	}

	wg.Wait()
}

func buildAndZip(osname, arch string) error {
	env := []string{
		fmt.Sprintf("GOPATH=%s", os.Getenv("GOPATH")),
		fmt.Sprintf("GOARCH=%s", arch),
		fmt.Sprintf("GOOS=%s", osname),
	}

	builddir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(builddir)

	printInfo("Building artefact for %s %s\n", osname, arch)

	var binName string

	switch osname {
	case "windows":
		binName = "awless.exe"
	default:
		binName = "awless"
	}

	artefactPath := filepath.Join(builddir, binName)

	if err := run(env, "go", "build", "-o", artefactPath, "-ldflags", "-s -w"); err != nil {
		return err
	}

	zipFile, err := os.OpenFile(fmt.Sprintf("%s-%s-%s.zip", strings.Split(binName, ".")[0], osname, arch), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	w := zip.NewWriter(zipFile)

	f, err := w.Create(binName)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(artefactPath)
	if err != nil {
		return err
	}

	if _, err = f.Write(content); err != nil {
		return err
	}

	return w.Close()
}

type environment []string

func run(env environment, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = env

	_, err := cmd.Output()
	if err != nil {
		printKo("error running [%s %s] with env %v\n", name, strings.Join(args, " "), env)

		if e, ok := err.(*exec.ExitError); ok {
			fmt.Println()
			fmt.Printf("%s\n", e.Stderr)
			fmt.Println()
		}

		return err
	}

	printOk("%s %s\n", name, strings.Join(args, " "))

	return nil
}

func printOk(s string, a ...interface{}) {
	fmt.Printf("\033[32m[OK]\033[m %s", fmt.Sprintf(s, a...))
}

func printKo(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "\033[31m[KO]\033[m %s", fmt.Sprintf(s, a...))
}

func printInfo(s string, a ...interface{}) {
	fmt.Printf("[+] %s", fmt.Sprintf(s, a...))
}
