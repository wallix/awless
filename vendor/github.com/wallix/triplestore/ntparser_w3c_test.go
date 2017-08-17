package triplestore

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestNTriplesW3CTestSuite(t *testing.T) {
	t.Run("positives", func(t *testing.T) {
		path := filepath.Join("testdata", "ntriples", "w3c_suite", "positives", "*.nt")
		filenames, _ := filepath.Glob(path)

		for _, filename := range filenames {
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatalf("cannot read file %s", filename)
			}

			tris, err := NewNTriplesDecoder(bytes.NewReader(b)).Decode()
			if err != nil {
				t.Fatal(err)
			}

			var buf bytes.Buffer
			if err := NewNTriplesEncoder(&buf).Encode(tris...); err != nil {
				t.Fatalf("file %s: re-encoding error: %s", filename, err)
			}

			if got, want := removeNTriplesCommentsAndEmptyLines(buf.Bytes()), removeNTriplesCommentsAndEmptyLines(b); !bytes.Equal(got, want) {
				t.Fatalf("file %s: re-encoding mismatch\n\ngot\n%q\n\nwant\n%q\n", filename, got, want)
			}
		}
	})

	t.Run("negatives", func(t *testing.T) {
		path := filepath.Join("testdata", "ntriples", "w3c_suite", "negatives", "*.nt")
		filenames, _ := filepath.Glob(path)

		for _, filename := range filenames {
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatalf("cannot read file %s", filename)
			}

			if _, err := NewNTriplesDecoder(bytes.NewReader(b)).Decode(); err == nil {
				t.Fatalf("filename '%s': expected err, got none", filename)
			}
		}
	})
}
