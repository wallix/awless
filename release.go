// +build ignore

package main

import (
	"compress/gzip"
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
				err := buildAndZip(o, a)
				if err != nil {
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

	builddir, err := ioutil.TempDir("", fmt.Sprintf("awless-%s-%s-build-", osname, arch))
	if err != nil {
		return err
	}
	defer os.RemoveAll(builddir)

	fmt.Printf("Building artefacts in %s\n", builddir)

	var binName string

	switch osname {
	case "windows":
		binName = "awless.exe"
	default:
		binName = "awless"
	}

	artefactPath := filepath.Join(builddir, binName)

	run(env, "go", "build", "-o", artefactPath, "-ldflags", "-s -w")
	if err != nil {
		return err
	}

	fi, err := os.OpenFile(filepath.Join(builddir, "awless.zip"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(artefactPath)
	if err != nil {
		return err
	}

	fw := gzip.NewWriter(fi)
	defer fw.Close()

	fw.Write(content)
	fw.Flush()

	finalZipname := fmt.Sprintf("awless-%s-%s.zip", osname, arch)

	err = run(nil, "mv", filepath.Join(builddir, "awless.zip"), finalZipname)
	if err != nil {
		return err
	}

	return nil
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
