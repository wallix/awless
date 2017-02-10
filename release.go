// +build ignore

/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	releaseTag = flag.String("rtag", "", "Git tag to be released")
)

var builds = map[string][]string{
	"darwin":  []string{"amd64"},
	"linux":   []string{"386", "amd64"},
	"windows": []string{"386", "amd64"},
}

func main() {
	flag.Parse()

	allBuild := map[string][]string{
		runtime.GOOS: []string{runtime.GOARCH},
	}

	if *releaseTag != "" {
		allBuild = builds
		printInfo("RELEASING")
	}

	var wg sync.WaitGroup

	for osname, archs := range allBuild {
		for _, arch := range archs {
			wg.Add(1)
			go func(o, a string) {
				defer wg.Done()
				if err := buildAndZip(o, a); err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
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

	printInfo("Building artefact for %s %s", osname, arch)

	var binName string

	switch osname {
	case "windows":
		binName = "awless.exe"
	default:
		binName = "awless"
	}

	artefactPath := filepath.Join(builddir, binName)

	gitRef := "refs/heads/master"
	if *releaseTag != "" {
		if tag, _ := runCmd(nil, "git", "describe", "--exact-match", "--tags"); strings.TrimSpace(tag) != *releaseTag {
			return fmt.Errorf("The git repository is not at tag '%s'", *releaseTag)
		}
		gitRef = fmt.Sprintf("refs/tags/%s", *releaseTag)
	}

	sha, err := runCmd(nil, "git", "show-ref", "-s", gitRef)
	if err != nil {
		return err
	}

	buildInfo := fmt.Sprintf("-X github.com/wallix/awless/config.buildDate=%s -X github.com/wallix/awless/config.buildSha=%s -X github.com/wallix/awless/config.buildOS=%s -X github.com/wallix/awless/config.buildArch=%s",
		time.Now().Format(time.RFC3339),
		strings.TrimSpace(sha),
		osname,
		arch,
	)

	ldflags := fmt.Sprintf("-ldflags=-s -w %s", buildInfo)

	if _, err := runCmd(env, "go", "build", "-o", artefactPath, ldflags); err != nil {
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

func runCmd(env environment, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = env

	out, err := cmd.Output()
	if err != nil {
		printKo("error running command [%s %s] with env %v", name, strings.Join(args, " "), env)

		if e, ok := err.(*exec.ExitError); ok {
			fmt.Println()
			fmt.Printf("%s\n", e.Stderr)
			fmt.Println()
		}

		return string(out), err
	}

	printOk("%s %s", name, strings.Join(args, " "))

	return string(out), nil
}

func printOk(s string, a ...interface{}) {
	fmt.Printf("\033[32m[OK]\033[m %s\n", fmt.Sprintf(s, a...))
}

func printKo(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "\033[31m[KO]\033[m %s\n", fmt.Sprintf(s, a...))
}

func printInfo(s string, a ...interface{}) {
	fmt.Printf("[+] %s\n", fmt.Sprintf(s, a...))
}
