package triplestore

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

type lenientNTParser struct {
	lex       *ntLexer
	lineCount int
}

func newLenientNTParser(r io.Reader) *lenientNTParser {
	return &lenientNTParser{
		lex: newNTLexer(r),
	}
}

func (p *lenientNTParser) parse() ([]Triple, error) {
	var triples []Triple
	var tok ntToken
	var nodeCount int
	var sub, pred, lit, datatype, langtag string
	var isLit, isResource, isSubBnode, isObjBnode, hasLangtag, hasDatatype, fullStopped bool
	var obj object

	reset := func() {
		sub, pred, lit, datatype, langtag = "", "", "", "", ""
		obj = object{}
		isLit, isResource, isSubBnode, isObjBnode, hasDatatype, hasLangtag, fullStopped = false, false, false, false, false, false, false
		nodeCount = 0
	}

	for tok.kind != EOF_TOK {
		var err error
		tok, err = p.lex.nextToken()
		if err != nil {
			return nil, err
		}
		switch tok.kind {
		case COMMENT_TOK:
			continue
		case IRI_TOK:
			nodeCount++
			switch nodeCount {
			case 1:
				sub = tok.lit
			case 2:
				pred = tok.lit
			case 3:
				isResource = true
				lit = tok.lit
			}
		case BNODE_TOK:
			nodeCount++
			switch nodeCount {
			case 1:
				sub = tok.lit
				isSubBnode = true
			case 2:
				return triples, errors.New("blank node can only be subject or object")
			case 3:
				isObjBnode = true
				lit = tok.lit
			}
		case LANGTAG_TOK:
			if nodeCount != 3 {
				return triples, errors.New("langtag misplaced")
			}
			hasLangtag = true
			langtag = tok.lit
		case LIT_TOK:
			if nodeCount != 2 {
				return triples, fmt.Errorf("tok '%s':reaching literate but missing element (node count %d)", tok.lit, nodeCount)
			}
			nodeCount++
			isLit = true
			lit = tok.lit
		case DATATYPE_TOK:
			hasDatatype = true
			datatype = tok.lit
		case FULLSTOP_TOK:
			if nodeCount != 3 {
				return triples, fmt.Errorf("reaching full stop but missing element (node count %d)", nodeCount)
			}
			fullStopped = true
			var tBuilder *tripleBuilder
			if isSubBnode {
				tBuilder = BnodePred(sub, pred)
			} else {
				tBuilder = SubjPred(sub, pred)
			}

			if isResource {
				triples = append(triples, tBuilder.Resource(lit))
			} else if isObjBnode {
				triples = append(triples, tBuilder.Bnode(lit))
			} else if isLit {
				if hasDatatype {
					obj = object{
						isLit: true,
						lit: literal{
							typ: XsdType(datatype),
							val: lit,
						},
					}
					triples = append(triples, tBuilder.Object(obj))
				} else if hasLangtag {
					triples = append(triples, tBuilder.StringLiteralWithLang(lit, langtag))
				} else {
					triples = append(triples, tBuilder.StringLiteral(lit))
				}
			}
			reset()
		case UNKNOWN_TOK:
			continue
		case LINEFEED_TOK:
			continue
		}
	}

	if nodeCount > 0 {
		return nil, fmt.Errorf("line %d: cannot parse at token '%s' (node count: %d)", p.lineCount, tok.lit, nodeCount)
	}

	if nodeCount != 0 && !fullStopped {
		return nil, errors.New("wrong number of elements")
	}

	return triples, nil
}

type ntTokenType int

const (
	UNKNOWN_TOK ntTokenType = iota
	IRI_TOK
	BNODE_TOK
	EOF_TOK
	WHITESPACE_TOK
	FULLSTOP_TOK
	LIT_TOK
	DATATYPE_TOK
	LANGTAG_TOK
	COMMENT_TOK
	LINEFEED_TOK
)

type ntToken struct {
	kind ntTokenType
	lit  string
}

func nodeTok(s string) ntToken     { return ntToken{kind: IRI_TOK, lit: s} }
func bnodeTok(s string) ntToken    { return ntToken{kind: BNODE_TOK, lit: s} }
func litTok(s string) ntToken      { return ntToken{kind: LIT_TOK, lit: s} }
func datatypeTok(s string) ntToken { return ntToken{kind: DATATYPE_TOK, lit: s} }
func langtagTok(s string) ntToken  { return ntToken{kind: LANGTAG_TOK, lit: s} }
func commentTok(s string) ntToken  { return ntToken{kind: COMMENT_TOK, lit: s} }
func unknownTok(s string) ntToken  { return ntToken{kind: UNKNOWN_TOK, lit: s} }

var (
	wspaceTok   = ntToken{kind: WHITESPACE_TOK, lit: " "}
	fullstopTok = ntToken{kind: FULLSTOP_TOK, lit: "."}
	lineFeedTok = ntToken{kind: LINEFEED_TOK, lit: "\n"}
	eofTok      = ntToken{kind: EOF_TOK}
)

type ntLexer struct {
	reader       *bufio.Reader
	buff         []byte
	current      rune
	width, index int
}

func newNTLexer(r io.Reader) *ntLexer {
	return &ntLexer{reader: bufio.NewReader(r)}
}

func (l *ntLexer) reset(input []byte) {
	l.current, l.width, l.index = 0, 0, 0
}

func (l *ntLexer) nextToken() (ntToken, error) {
	if err := l.readRune(); err != nil {
		return ntToken{}, err
	}

	switch l.current {
	case '<':
		n, err := l.readNode()
		return nodeTok(n), err
	case '_':
		if err := l.readRune(); err != nil {
			return ntToken{}, err
		}
		if l.current != ':' {
			return ntToken{}, fmt.Errorf("invalid blank node: expecting ':', got '%c'", l.current)
		}
		n, err := l.readBnode()
		return bnodeTok(n), err
	case ' ':
		return wspaceTok, nil
	case '.':
		return fullstopTok, nil
	case '\n':
		return lineFeedTok, nil
	case '"':
		n, err := l.readStringLiteral()
		return litTok(n), err
	case '@':
		n, err := l.readBnode()
		return langtagTok(n), err
	case '^':
		if err := l.readRune(); err != nil {
			return ntToken{}, err
		}
		if l.current == 0 {
			return eofTok, nil
		}
		if l.current != '^' {
			return ntToken{}, fmt.Errorf("invalid datatype: expecting '^', got '%c'", l.current)
		}
		if err := l.readRune(); err != nil {
			return ntToken{}, err
		}
		if l.current == 0 {
			return eofTok, nil
		}
		if l.current != '<' {
			return ntToken{}, fmt.Errorf("invalid datatype: expecting '<', got '%c'", l.current)
		}
		n, err := l.readNode()
		return datatypeTok(n), err
	case '#':
		n, err := l.readComment()
		return commentTok(n), err
	case 0:
		return eofTok, nil
	default:
		return unknownTok(string(l.current)), nil
	}
}

func (l *ntLexer) readRune() (err error) {
	l.current, l.width, err = l.reader.ReadRune()
	if l.current == utf8.RuneError && l.width == 1 {
		return errors.New("lexer read: invalid utf8 encoding")
	}
	if err == io.EOF || l.width == 0 {
		l.current = 0
		return nil
	}
	l.index = l.index + l.width
	l.buff = append(l.buff, []byte(string(l.current))...)
	return nil
}

func (l *ntLexer) unreadRune() {
	for i := 0; i < l.width; i++ {
		l.reader.UnreadByte()
	}
	l.index = l.index - l.width
	l.buff = l.buff[:len(l.buff)-l.width]
	if len(l.buff) > 0 {
		l.current, _ = utf8.DecodeLastRune(l.buff)
	}
}

func (l *ntLexer) peekNextNonWithespaceRune() (found rune, err error) {
	index := 1
	var last byte
	for {
		b, err := l.reader.Peek(index)
		if err == io.EOF {
			return 0, nil
		}
		if err != nil {
			return 0, err
		}
		if l := len(b); l > 0 {
			last = b[l-1]
		} else {
			last = 0
		}
		if last != ' ' && last != '\t' {
			break
		}
		index++
	}

	for {
		b, err := l.reader.Peek(index)
		if err == io.EOF {
			return 0, nil
		}
		if err != nil {
			return 0, err
		}
		r, _ := utf8.DecodeLastRune(b)
		if r == utf8.RuneError {
			index++
			continue
		} else {
			return r, err
		}
	}
}

func (l *ntLexer) readNode() (string, error) {
	l.buff = []byte{}
	for {
		if err := l.readRune(); err != nil {
			return "", err
		}
		if l.current == '>' {
			peek, err := l.peekNextNonWithespaceRune()
			if err != nil {
				return "", err
			}
			if peek == 0 || peek == '<' || peek == '"' || peek == '.' || peek == '_' {
				return l.extractString(), nil
			}
		}
		if l.current == 0 {
			return "", nil
		}
	}
}

func (l *ntLexer) readStringLiteral() (string, error) {
	l.buff = []byte{}
	for {
		if err := l.readRune(); err != nil {
			return "", err
		}
		if l.current == '"' {
			peek, err := l.peekNextNonWithespaceRune()
			if err != nil {
				return "", err
			}
			if peek == 0 || peek == '.' || peek == '^' || peek == '@' {
				return l.extractString(), nil
			}
		}
		if l.current == 0 {
			return "", nil
		}
	}
}

func (l *ntLexer) readBnode() (string, error) {
	l.buff = []byte{}
	for {
		if err := l.readRune(); err != nil {
			return "", err
		}
		if l.current == ' ' {
			peek, err := l.peekNextNonWithespaceRune()
			if err != nil {
				return "", err
			}
			if peek == 0 || peek == '<' || peek == '.' {
				s := l.extractString()
				l.unreadRune()
				return s, nil
			}
		}
		if l.current == '.' {
			peek, err := l.peekNextNonWithespaceRune()
			if err != nil {
				return "", err
			}
			if peek == 0 || peek == '#' || peek == '\n' { // brittle: but handles <sub> <pred> _:bnode.#commenting
				s := l.extractString()
				l.unreadRune()
				return s, nil
			}
		}
		if l.current == 0 {
			return "", nil
		}
		if l.current == '<' {
			s := l.extractString()
			l.unreadRune()
			return s, nil
		}
	}
}

func (l *ntLexer) readComment() (string, error) {
	l.buff = []byte{}
	for {
		if err := l.readRune(); err != nil {
			return "", err
		}
		if l.current == '\n' {
			s := l.extractString()
			l.unreadRune()
			return s, nil
		}
		if l.current == 0 {
			return l.extractString(), nil
		}
	}
}

func (l *ntLexer) extractString() string {
	return string(l.buff[:len(l.buff)-l.width])
}
