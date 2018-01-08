package triplestore

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestMultilineEmptyAndCommentLine(t *testing.T) {
	p := newLenientNTParser(strings.NewReader(`  # my triples

# starting
<sub><pred>"obj"@en .

# ending

`))
	tris, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	src := NewSource()
	src.Add(tris...)
	snap := src.Snapshot()
	if got, want := snap.Count(), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if !snap.Contains(SubjPred("sub", "pred").StringLiteralWithLang("obj", "en")) {
		t.Fatal("expected to contain triple")
	}
}

func TestParsing(t *testing.T) {
	tcases := []struct {
		input    string
		expected []Triple
	}{
		{
			input: `<sub> <pred> "quoting "anything".".`,
			expected: []Triple{
				SubjPred("sub", "pred").StringLiteral(`quoting "anything".`),
			},
		},
		{
			input: `<sub> <pred> "quoting 'anything'.".`,
			expected: []Triple{
				SubjPred("sub", "pred").StringLiteral("quoting 'anything'."),
			},
		},
		{
			input: "	<sub>	<pred> <lol> .\n<sub2> <pred2> \"lol2\" .",
			expected: []Triple{
				SubjPred("sub", "pred").Resource("lol"),
				mustTriple("sub2", "pred2", "lol2"),
			},
		},
		{
			input: "<sub> <pred> \"2\"^^<myinteger> .\n<sub2> <pred2> <lol2> .",
			expected: []Triple{
				SubjPred("sub", "pred").Object(object{isLit: true, lit: literal{typ: "myinteger", val: "2"}}),
				SubjPred("sub2", "pred2").Resource("lol2"),
			},
		},
		{
			input: "<sub><pred> \"2\"^^<myinteger> .\n<sub2> <pred2> \"lol2\"@en.",
			expected: []Triple{
				SubjPred("sub", "pred").Object(object{isLit: true, lit: literal{typ: "myinteger", val: "2"}}),
				SubjPred("sub2", "pred2").StringLiteralWithLang("lol2", "en"),
			},
		},
		{
			input:    "_:sub<pred><obj>. # comment",
			expected: []Triple{BnodePred("sub", "pred").Resource("obj")},
		},
		{
			input:    "_:sub <pred><obj>. # comment",
			expected: []Triple{BnodePred("sub", "pred").Resource("obj")},
		},
		{
			input:    "<sub> <pred> \"dquote:\"\" .\n",
			expected: []Triple{SubjPred("sub", "pred").StringLiteral(`dquote:"`)},
		},
		{
			input:    "<sub><pred><obj>.\n",
			expected: []Triple{SubjPred("sub", "pred").Resource("obj")},
		},
		{
			input:    "<sub> <pred> _:anon.\n",
			expected: []Triple{SubjPred("sub", "pred").Bnode("anon")},
		},
		{
			input:    "<sub><pred>_:anon.\n",
			expected: []Triple{SubjPred("sub", "pred").Bnode("anon")},
		},
		{
			input:    `<sub> <pred> _:anon.`,
			expected: []Triple{SubjPred("sub", "pred").Bnode("anon")},
		},
		{
			input:    "<sub> <pred> \"\u00E9\".\n",
			expected: []Triple{SubjPred("sub", "pred").StringLiteral("é")},
		},
		{
			input:    "<sub> <pred> \"\u00E9\".",
			expected: []Triple{SubjPred("sub", "pred").StringLiteral("é")},
		},
		{
			input:    "<sub> <pred> \"\032\".",
			expected: []Triple{SubjPred("sub", "pred").StringLiteral("\032")},
		},
		{
			input:    "<sub> <pred> \"\x1A\".",
			expected: []Triple{SubjPred("sub", "pred").StringLiteral("\x1A")},
		},
	}

	for j, tcase := range tcases {
		p := newLenientNTParser(strings.NewReader(tcase.input))
		tris, err := p.Parse()
		if err != nil {
			t.Fatalf("input=[%s]: %s", tcase.input, err)
		}
		if got, want := len(tris), len(tcase.expected); got != want {
			t.Fatalf("triples size (case %d): got %d, want %d", j+1, got, want)
		}
		for i, tri := range tris {
			if got, want := tri, tcase.expected[i]; !got.Equal(want) {
				t.Fatalf("case %d: input [%s]: triple (%d)\ngot %#v\n\nwant %#v", j+1, tcase.input, i+1, got, want)
			}
		}
	}
}

func TestParsingComponents(t *testing.T) {
	t.Run("literal object", func(t *testing.T) {
		tcases := []struct {
			input, left, exp string
		}{
			{input: `stuff"`},
			{input: `"`},
			{input: `stuff" .`, exp: "stuff", left: "."},
			{input: `stuff" .   `, exp: "stuff", left: ".   "},
			{input: `stuff"	 .# comment`, exp: "stuff", left: ".# comment"},
			{input: `stuff"	 . # comment`, exp: "stuff", left: ". # comment"},
			{input: ` " .`, exp: " ", left: "."},
			{input: `" .`, exp: "", left: "."},
			{input: `stuff"^`},
			{input: `stuff"^^`, exp: "stuff", left: "^^"},
			{input: `stuff" ^^   `, exp: "stuff", left: "^^   "},
			{input: `stuff"@`, exp: "stuff", left: "@"},
			{input: `stuff" @   `, exp: "stuff", left: "@   "},
		}
		for _, tcase := range tcases {
			s, left, _ := parseLiteralObject([]byte(tcase.input))
			if got, want := s, tcase.exp; got != want {
				t.Fatalf("case [%s]: got '%s', want '%s'", tcase.input, got, want)
			}
			if got, want := left, []byte(tcase.left); !bytes.Equal(got, want) {
				t.Fatalf("case [%s]: left: got '%s', want '%s'", tcase.input, got, want)
			}
		}
	})

	t.Run("bnode object", func(t *testing.T) {
		tcases := []struct {
			input, left, exp string
		}{
			{input: "stuff"},
			{input: "stuff.", exp: "stuff", left: "."},
			{input: "stuff .", exp: "stuff", left: "."},
			{input: "stuff	 .   ", exp: "stuff", left: ".   "},
			{input: "stuff	 .# comment", exp: "stuff", left: ".# comment"},
			{input: "stuff	 . # comment", exp: "stuff", left: ". # comment"},
			{input: " .", exp: "", left: "."},
		}
		for _, tcase := range tcases {
			s, left, _ := parseBNodeObject([]byte(tcase.input))
			if got, want := s, tcase.exp; got != want {
				t.Fatalf("case [%s]: got '%s', want '%s'", tcase.input, got, want)
			}
			if got, want := left, []byte(tcase.left); !bytes.Equal(got, want) {
				t.Fatalf("case [%s]: left: got '%s', want '%s'", tcase.input, got, want)
			}
		}
	})

	t.Run("object iri", func(t *testing.T) {
		tcases := []struct {
			input, left, exp string
		}{
			{input: "stuff>"},
			{input: "stuff>.", exp: "stuff", left: "."},
			{input: "stuff> .", exp: "stuff", left: "."},
			{input: "stuff>	 .   ", exp: "stuff", left: ".   "},
			{input: "stuff>	 .# comment", exp: "stuff", left: ".# comment"},
			{input: "stuff>	 . # comment", exp: "stuff", left: ". # comment"},
			{input: ">.", left: "."},
			{input: "> .", left: "."},
		}
		for _, tcase := range tcases {
			s, left, _ := parseIRIObject([]byte(tcase.input))
			if got, want := s, tcase.exp; got != want {
				t.Fatalf("case [%s]: got '%s', want '%s'", tcase.input, got, want)
			}
			if got, want := left, []byte(tcase.left); !bytes.Equal(got, want) {
				t.Fatalf("case [%s]: left: got '%s', want '%s'", tcase.input, got, want)
			}
		}
	})

	t.Run("predicate iri", func(t *testing.T) {
		tcases := []struct {
			input, left, exp string
		}{
			{input: "stuff>"},
			{input: "stuff><", exp: "stuff", left: "<"},
			{input: "stuff> <", exp: "stuff", left: "<"},
			{input: "stuff>	  <obj", exp: "stuff", left: "<obj"},
			{input: `stuff>"`, exp: "stuff", left: "\""},
			{input: `stuff> 	"`, exp: "stuff", left: "\""},
			{input: `stuff>_`, exp: "stuff", left: "_"},
			{input: `stuff> 	_`, exp: "stuff", left: "_"},
			{input: "><", left: "<"},
		}
		for _, tcase := range tcases {
			s, left, _ := parsePredicate([]byte(tcase.input))
			if got, want := s, tcase.exp; got != want {
				t.Fatalf("case [%s]: got '%s', want '%s'", tcase.input, got, want)
			}
			if got, want := left, []byte(tcase.left); !bytes.Equal(got, want) {
				t.Fatalf("case [%s]: left: got '%s', want '%s'", tcase.input, got, want)
			}
		}
	})

	t.Run("subject iri", func(t *testing.T) {
		tcases := []struct {
			input, left, exp string
		}{
			{input: "stuff>"},
			{input: "stuff><", exp: "stuff", left: "<"},
			{input: "stuff> <", exp: "stuff", left: "<"},
			{input: "stuff>	  <pred", exp: "stuff", left: "<pred"},
			{input: "><", left: "<"},
		}
		for _, tcase := range tcases {
			s, left, _ := parseIRISubject([]byte(tcase.input))

			if got, want := s, tcase.exp; got != want {
				t.Fatalf("case [%s]: got '%s', want '%s'", tcase.input, got, want)
			}
			if got, want := left, []byte(tcase.left); !bytes.Equal(got, want) {
				t.Fatalf("case [%s]: left: got '%s', want '%s'", tcase.input, got, want)
			}
		}
	})

	t.Run("subject bnode", func(t *testing.T) {
		tcases := []struct {
			input, left, exp string
		}{
			{input: "stuff"},
			{input: "stuff <", exp: "stuff", left: "<"},
			{input: "stuff <    ", exp: "stuff", left: "<    "},
		}
		for _, tcase := range tcases {
			s, left, _ := parseBNodeSubject([]byte(tcase.input))
			if got, want := s, tcase.exp; got != want {
				t.Fatalf("got '%s', want '%s'", got, want)
			}
			if got, want := left, []byte(tcase.left); !bytes.Equal(got, want) {
				t.Fatalf("case [%s]: left: got '%s', want '%s'", tcase.input, got, want)
			}

		}
	})
}

func TestParserErrorHandling(t *testing.T) {
	tcases := []struct {
		input       string
		errContains string
	}{
		{input: "<sub> <pred> 1 ."},
		//{input: "<one> <two> <three>, <four> ."}, passes
	}

	for _, tcase := range tcases {
		tris, err := newLenientNTParser(strings.NewReader(tcase.input)).Parse()
		if err == nil {
			t.Fatalf("expected err, got none. Triples parsed:\n%#v", Triples(tris).Map(func(tr Triple) string { return fmt.Sprint(tr) }))
		}
		if msg := tcase.errContains; msg != "" {
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected '%s' to contains '%s'", err.Error(), tcase.errContains)
			}
		}
	}
}

func mustTriple(s, p string, i interface{}) Triple {
	t, err := SubjPredLit(s, p, i)
	if err != nil {
		panic(err)
	}
	return t
}
