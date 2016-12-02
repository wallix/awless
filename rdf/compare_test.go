package rdf

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompareInfraGraph(t *testing.T) {
	cases := []*diffTest{}

	filepath.Walk("testdata/infras", collectTestCases(&cases))

	for _, tcase := range cases {
		localG := NewGraph()
		localG.Unmarshal(tcase.local.Bytes())

		remoteG := NewGraph()
		remoteG.Unmarshal(tcase.remote.Bytes())

		extrasG, missingsG, _, err := Compare("eu-west-1", localG, remoteG)
		if err != nil {
			t.Fatal(err)
		}

		expectExtrasG := NewGraph()
		expectExtrasG.Unmarshal(tcase.extras.Bytes())

		expectMissingsG := NewGraph()
		expectMissingsG.Unmarshal(tcase.missings.Bytes())

		if got, want := extrasG.MustMarshal(), expectExtrasG.MustMarshal(); got != want {
			t.Fatalf("\n[%s] - extras: got\n[%s]\n\nwant\n[%s]\n\n", tcase.filepath, got, want)
		}
		if got, want := missingsG.MustMarshal(), expectMissingsG.MustMarshal(); got != want {
			t.Fatalf("\n[%s] - missings: got\n[%s]\n\nwant\n[%s]\n\n", tcase.filepath, got, want)
		}

		extrasG, missingsG, _, err = Compare("eu-west-1", remoteG, localG)
		if err != nil {
			t.Fatal(err)
		}

		expectInvExtrasG := NewGraph()
		expectInvExtrasG.Unmarshal(tcase.invextras.Bytes())

		expectInvMissingG := NewGraph()
		expectInvMissingG.Unmarshal(tcase.invmissings.Bytes())

		if got, want := extrasG.MustMarshal(), expectInvExtrasG.MustMarshal(); got != want {
			t.Fatalf("\n[%s] - inv extras: got\n[%s]\n\nwant\n[%s]\n\n", tcase.filepath, got, want)
		}
		if got, want := missingsG.MustMarshal(), expectInvMissingG.MustMarshal(); got != want {
			t.Fatalf("\n[%s] - inv missings: got\n[%s]\n\nwant\n[%s]\n\n", tcase.filepath, got, want)
		}
	}
}

func collectTestCases(collect *[]*diffTest) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			*collect = append(*collect, parseTestfile(path))
		}
		return nil
	}
}

type diffTest struct {
	filepath                                                string
	local, remote, extras, missings, invextras, invmissings bytes.Buffer
}

type section int

const (
	Start section = iota
	Local
	Remote
	Extras
	Missings
	InvExtras
	InvMissings
)

func parseTestfile(path string) *diffTest {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	test := &diffTest{filepath: path}
	r := bufio.NewReader(file)

	where := Start

Loop:
	for {
		line, err := r.ReadString('\n')

		switch {
		case err == io.EOF:
			if len(line) > 0 {
				test.invmissings.WriteString(line)
			}
			break Loop
		case err != nil:
			log.Fatal(err)
		}

		switch where {
		case Local:
			test.local.WriteString(line)
		case Remote:
			test.remote.WriteString(line)
		case Extras:
			test.extras.WriteString(line)
		case Missings:
			test.missings.WriteString(line)
		case InvExtras:
			test.invextras.WriteString(line)
		case InvMissings:
			test.invmissings.WriteString(line)
		}

		switch {
		case strings.TrimSpace(line) == "":
		case strings.HasPrefix(line, "#local"):
			where = Local
		case strings.HasPrefix(line, "#remote"):
			where = Remote
		case strings.HasPrefix(line, "#extras"):
			where = Extras
		case strings.HasPrefix(line, "#missings"):
			where = Missings
		case strings.HasPrefix(line, "#invextras"):
			where = InvExtras
		case strings.HasPrefix(line, "#invmissings"):
			where = InvMissings
		}
	}

	return test
}
