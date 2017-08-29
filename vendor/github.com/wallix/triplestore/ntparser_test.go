package triplestore

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNTParser(t *testing.T) {
	tcases := []struct {
		input    string
		expected []Triple
	}{
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
			input:    `<sub> <pred> "dquote:\"" .\n`,
			expected: []Triple{SubjPred("sub", "pred").StringLiteral(`dquote:\"`)},
		},
		{
			input:    "<sub> <pred> _:anon.\n",
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
		tris, err := p.parse()
		if err != nil {
			t.Fatalf("input=[%s]: %s", tcase.input, err)
		}
		if got, want := len(tris), len(tcase.expected); got != want {
			t.Fatalf("triples size (case %d): got %d, want %d", j+1, got, want)
		}
		for i, tri := range tris {
			if got, want := tri, tcase.expected[i]; !got.Equal(want) {
				t.Fatalf("triple (%d)\ngot %v\n\nwant %v", i+1, got, want)
			}
		}
	}
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
		tris, err := newLenientNTParser(strings.NewReader(tcase.input)).parse()
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

func TestLexer(t *testing.T) {
	tcases := []struct {
		input    string
		expected []ntToken
	}{
		// simple
		{"<node>", []ntToken{nodeTok("node")}},
		{"_:bnode .", []ntToken{bnodeTok("bnode"), wspaceTok, fullstopTok}},
		{"_:bnode <pred>", []ntToken{bnodeTok("bnode"), wspaceTok, nodeTok("pred")}},
		{"#comment", []ntToken{commentTok("comment")}},
		{"# comment", []ntToken{commentTok(" comment")}},
		{"\"lit\"", []ntToken{litTok("lit")}},
		{"^^<xsd:float>", []ntToken{datatypeTok("xsd:float")}},
		{" ", []ntToken{wspaceTok}},
		{".", []ntToken{fullstopTok}},
		{"\n", []ntToken{lineFeedTok}},
		{"# comment\n", []ntToken{commentTok(" comment"), lineFeedTok}},
		{"@en .", []ntToken{langtagTok("en"), wspaceTok, fullstopTok}},
		{"@en.", []ntToken{langtagTok("en"), fullstopTok}},
		{"@en .\n", []ntToken{langtagTok("en"), wspaceTok, fullstopTok, lineFeedTok}},

		{"#", []ntToken{commentTok("")}}, // fixed with go-fuzz

		// escaped
		{`<no>de>`, []ntToken{nodeTok("no>de")}},
		{`<no\>de>`, []ntToken{nodeTok("no\\>de")}},
		{`<node\\>`, []ntToken{nodeTok("node\\\\")}},
		{`"\\"`, []ntToken{litTok(`\\`)}},
		{`"quot"ed"`, []ntToken{litTok(`quot"ed`)}},
		{`"quot\"ed"`, []ntToken{litTok("quot\\\"ed")}},

		// triples
		{"<sub> <pred> \"3\"^^<xsd:integer> .", []ntToken{
			nodeTok("sub"), wspaceTok, nodeTok("pred"), wspaceTok, litTok("3"),
			datatypeTok("xsd:integer"), wspaceTok, fullstopTok,
		}},
		{"<sub><pred>\"3\"^^<xsd:integer>.", []ntToken{
			nodeTok("sub"), nodeTok("pred"), litTok("3"), datatypeTok("xsd:integer"), fullstopTok,
		}},
		{"<sub> <pred> \"lit\" . # commenting", []ntToken{
			nodeTok("sub"), wspaceTok, nodeTok("pred"), wspaceTok, litTok("lit"),
			wspaceTok, fullstopTok, wspaceTok, commentTok(" commenting"),
		}},
		{"<sub><pred>\"lit\".#commenting", []ntToken{
			nodeTok("sub"), nodeTok("pred"), litTok("lit"), fullstopTok, commentTok("commenting"),
		}},

		// triple with bnodes
		{"_:sub <pred>\"lit\".#commenting", []ntToken{
			bnodeTok("sub"), wspaceTok, nodeTok("pred"), litTok("lit"), fullstopTok, commentTok("commenting"),
		}},
		{"<sub> <pred> _:lit . #commenting", []ntToken{
			nodeTok("sub"), wspaceTok, nodeTok("pred"), wspaceTok, bnodeTok("lit"), wspaceTok, fullstopTok, wspaceTok, commentTok("commenting"),
		}},
		{"_:sub<pred>_:lit.#commenting", []ntToken{
			bnodeTok("sub"), nodeTok("pred"), bnodeTok("lit"), fullstopTok, commentTok("commenting"),
		}},

		// triples with langtag
		{`<sub> <pred> "lit"@russ . # commenting`, []ntToken{
			nodeTok("sub"), wspaceTok, nodeTok("pred"), wspaceTok, litTok("lit"),
			langtagTok("russ"), wspaceTok, fullstopTok, wspaceTok, commentTok(" commenting"),
		}},
	}

	for i, tcase := range tcases {
		l := newNTLexer(strings.NewReader(tcase.input))
		var toks []ntToken
		for tok, _ := l.nextToken(); tok.kind != EOF_TOK; tok, _ = l.nextToken() {
			toks = append(toks, tok)
		}
		if got, want := toks, tcase.expected; !reflect.DeepEqual(got, want) {
			t.Fatalf("case %d input=[%s]\ngot %#v\n\nwant %#v", i+1, tcase.input, got, want)
		}
	}
}

func TestLexerReadNode(t *testing.T) {
	tcases := []struct {
		input string
		node  string
	}{
		{"<", ""},
		{">", ""},
		{" >", " "},
		{"", ""},
		{"z", ""},
		{`\z>`, "\\z"},
		{"\n>", "\n"},

		{"subject>", "subject"},
		{"subject> ", "subject"},
		{"subject> .", "subject"},
		{"s  ubject>", "s  ubject"},
		{"subject>   <", "subject"},
		{"subject>  	 <", "subject"}, // with tabs
		{"    subject>   <", "    subject"},
		{"subject><", "subject"},
		{"subje   ct><", "subje   ct"},
		{"sub>ject>", "sub>ject"},
		{"sub > ject>", "sub > ject"},
		{"sub>ject>      ", "sub>ject"},
		{"subject", ""},

		{"pred>   \"", "pred"},
		{"pred>\"", "pred"},

		{"resource>.", "resource"},
		{"resource> .", "resource"},
		{"resource>> .", "resource>"},
		{"resource>  .   ", "resource"},
	}

	for i, tcase := range tcases {
		l := newNTLexer(strings.NewReader(tcase.input))
		n, err := l.readNode()
		if err != nil {
			t.Fatalf("case %d: '%s': %s", i+1, tcase.input, err)
		}
		if got, want := n, tcase.node; got != want {
			t.Fatalf("case %d '%s': got '%s', want '%s'", i+1, tcase.input, got, want)
		}
	}
}

func TestLexerReadBnode(t *testing.T) {
	tcases := []struct {
		input string
		node  string
	}{
		{"a .", "a"},
		{"a<", "a"},
		{"a    <", "a"},
		{"a <", "a"},
		{"a .", "a"},
		{"a.", "a"},
		{"a     .", "a"},
		{"a.\n", "a"},
	}

	for i, tcase := range tcases {
		l := newNTLexer(strings.NewReader(tcase.input))
		n, err := l.readBnode()
		if err != nil {
			t.Fatalf("case %d: '%s': %s", i+1, tcase.input, err)
		}
		if got, want := n, tcase.node; got != want {
			t.Fatalf("case %d '%s': got '%s', want '%s'", i+1, tcase.input, got, want)
		}
	}
}
func TestLexerReadStringLiteral(t *testing.T) {
	tcases := []struct {
		input string
		node  string
	}{
		{"", ""},
		{`"`, ""},
		{`  "`, "  "},
		{"z", ""},
		{"\u00E9\" .", "\u00E9"},
		{`\n"`, "\\n"},
		{`lit"`, "lit"},
		{`l it"`, "l it"},
		{"li\"t\"", "li\"t"},
		{"li \"t\"", "li \"t"},
		{"li\"t\" .", "li\"t"},
		{"li\"t\".", "li\"t"},
		{"li\"t\" .", "li\"t"},
		{"li\"t\"  .  ", "li\"t"},
		{"li\"t\"^", "li\"t"},
		{"li\"t\"^^", "li\"t"},
		{"li\"t\" ^", "li\"t"},
		{"li\"t\" ^^", "li\"t"},
		{"li\"t\"   ^", "li\"t"},
		{"li\"t\"     ^^", "li\"t"},
	}

	for i, tcase := range tcases {
		l := newNTLexer(strings.NewReader(tcase.input))
		n, err := l.readStringLiteral()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := n, tcase.node; got != want {
			t.Fatalf("case %d: got '%s', want '%s'", i+1, got, want)
		}
	}
}

func TestLexerReadComment(t *testing.T) {
	tcases := []struct {
		input string
		node  string
	}{
		{"", ""},
		{"#", "#"},
		{" comment \n", " comment "},
		{"\n", ""},
	}

	for i, tcase := range tcases {
		l := newNTLexer(strings.NewReader(tcase.input))
		n, err := l.readComment()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := n, tcase.node; got != want {
			t.Fatalf("case %d: got '%s', want '%s'", i+1, got, want)
		}
	}
}

func TestLexerPeekRunes(t *testing.T) {
	tcases := []struct {
		input string
		found rune
	}{
		{input: "     t", found: 't'},
		{input: "t", found: 't'},
		{input: "", found: rune(0)},
		{input: " ", found: rune(0)},
	}

	for _, tcase := range tcases {
		l := newNTLexer(strings.NewReader(tcase.input))
		found, _ := l.peekNextNonWithespaceRune()
		if got, want := found, tcase.found; got != want {
			t.Fatalf("input [%s]: got %q, want %q", tcase.input, got, want)
		}
	}
}

func TestLexerReadUnreadRunes(t *testing.T) {
	t.Run("sentence", func(t *testing.T) {
		l := newNTLexer(strings.NewReader("tron"))
		l.readRune()
		if got, want := l.current, 't'; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		l.readRune()
		if got, want := l.current, 'r'; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		l.unreadRune()
		if got, want := l.current, 't'; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		if got, want := l.width, 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
	})

	t.Run("read empty", func(t *testing.T) {
		l := newNTLexer(strings.NewReader(""))
		l.readRune()
		if got, want := l.current, rune(0); got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		if got, want := l.width, 0; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		l.readRune()
		if got, want := l.current, rune(0); got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		if got, want := l.width, 0; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
	})

	t.Run("unread empty", func(t *testing.T) {
		l := newNTLexer(strings.NewReader(""))
		l.unreadRune()
		if got, want := l.current, rune(0); got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		if got, want := l.width, 0; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		l.unreadRune()
		if got, want := l.current, rune(0); got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		if got, want := l.width, 0; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
	})

	t.Run("read length one", func(t *testing.T) {
		l := newNTLexer(strings.NewReader("s"))
		l.readRune()
		if got, want := l.current, 's'; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		l.readRune()
		if got, want := l.current, rune(0); got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		if got, want := l.width, 0; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		l.readRune()
		if got, want := l.current, rune(0); got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		l.unreadRune()
		if got, want := l.current, 's'; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("unread length one", func(t *testing.T) {
		l := newNTLexer(strings.NewReader("s"))
		l.unreadRune()
		if got, want := l.current, rune(0); got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		l.unreadRune()
		if got, want := l.current, rune(0); got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		if got, want := l.width, 0; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		l.readRune()
		if got, want := l.current, 's'; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

func mustTriple(s, p string, i interface{}) Triple {
	t, err := SubjPredLit(s, p, i)
	if err != nil {
		panic(err)
	}
	return t
}
