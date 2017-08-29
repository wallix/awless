package triplestore

import (
	"bytes"
	"io/ioutil"
	"os"
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

			tris, err := NewLenientNTDecoder(bytes.NewReader(b)).Decode()
			if err != nil {
				t.Fatalf("file %s: %s", filename, err)
			}

			var buf bytes.Buffer
			if err := NewLenientNTEncoder(&buf).Encode(tris...); err != nil {
				t.Fatalf("file %s: re-encoding error: %s", filename, err)
			}

			expected := cleanupNTriplesForComparison(b)
			expectedFilepath := filename + ".expected"
			if _, err := os.Stat(expectedFilepath); !os.IsNotExist(err) {
				expected, err = ioutil.ReadFile(expectedFilepath)
				if err != nil {
					t.Fatal(err)
				}
			}

			if got, want := cleanupNTriplesForComparison(buf.Bytes()), expected; !bytes.Equal(got, want) {
				t.Fatalf("file %s: re-encoding mismatch\n\ngot\n%s\n\nwant\n%s\n", filename, got, want)
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

			if _, err := NewLenientNTDecoder(bytes.NewReader(b)).Decode(); err == nil {
				t.Fatalf("filename '%s': expected err, got none", filename)
			}
		}
	})
}
