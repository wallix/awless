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
	"darwin": []string{"amd64"},
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

	fmt.Printf("[+] Building artefacts in tmp dir %s\n", builddir)

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

	zipFile, err := os.OpenFile(fmt.Sprintf("%s-%s-%s.zip", binName, osname, arch), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
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
		fmt.Fprintf(os.Stderr, "error running [%s %s] with env %v\n", name, strings.Join(args, " "), env)

		if e, ok := err.(*exec.ExitError); ok {
			fmt.Println()
			fmt.Printf("%s\n", e.Stderr)
			fmt.Println()
		}

		return err
	}

	fmt.Printf("[OK] %s %s\n", name, strings.Join(args, " "))
	return nil
}
