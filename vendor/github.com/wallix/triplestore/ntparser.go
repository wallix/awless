package triplestore

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

type ntParser struct {
	lex *lexer
}

func newNTParser(s string) *ntParser {
	return &ntParser{
		lex: newLexer(s),
	}
}

func (p *ntParser) parse() ([]Triple, error) {
	var tris []Triple
	var tok ntToken
	var nodeCount int
	var sub, pred, lit, datatype string
	var isLit, isResource, hasDatatype, fullStopped bool
	var obj object

	reset := func() {
		sub, pred, lit, datatype = "", "", "", ""
		obj = object{}
		isLit, isResource, hasDatatype, fullStopped = false, false, false, false
		nodeCount = 0
	}

	for tok.kind != EOF_TOK {
		var err error
		tok, err = p.lex.nextToken()
		if err != nil {
			return tris, err
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
		case LIT_TOK:
			if nodeCount != 2 {
				return tris, errors.New("reaching literate but missing element")
			}
			nodeCount++
			isLit = true
			lit = tok.lit
		case DATATYPE_TOK:
			hasDatatype = true
			datatype = tok.lit
		case FULLSTOP_TOK:
			if nodeCount != 3 {
				return tris, errors.New("reaching full stop but missing element")
			}
			fullStopped = true
			if isResource {
				tris = append(tris, SubjPred(sub, pred).Resource(lit))
			} else if isLit {
				if hasDatatype {
					obj = object{
						isLit: true,
						lit: literal{
							typ: XsdType(datatype),
							val: lit,
						},
					}
					tris = append(tris, SubjPred(sub, pred).Object(obj))
				} else {
					tris = append(tris, SubjPred(sub, pred).StringLiteral(lit))
				}
			}
			reset()
		case UNKNOWN_TOK:
			return tris, fmt.Errorf("unknown token '%s'", tok.lit)
		case LINEFEED_TOK:
			continue
		}
	}

	if nodeCount > 0 {
		return tris, errors.New("cannot parse line")
	}

	if nodeCount != 0 && !fullStopped {
		return tris, errors.New("wrong number of elements")
	}

	return tris, nil
}

type ntTokenType int

const (
	UNKNOWN_TOK ntTokenType = iota
	IRI_TOK
	EOF_TOK
	WHITESPACE_TOK
	FULLSTOP_TOK
	LIT_TOK
	DATATYPE_TOK
	COMMENT_TOK
	LINEFEED_TOK
)

type ntToken struct {
	kind ntTokenType
	lit  string
}

func iriTok(s string) ntToken      { return ntToken{kind: IRI_TOK, lit: s} }
func litTok(s string) ntToken      { return ntToken{kind: LIT_TOK, lit: s} }
func datatypeTok(s string) ntToken { return ntToken{kind: DATATYPE_TOK, lit: s} }
func commentTok(s string) ntToken  { return ntToken{kind: COMMENT_TOK, lit: s} }
func unknownTok(s string) ntToken  { return ntToken{kind: UNKNOWN_TOK, lit: s} }

var (
	wspaceTok   = ntToken{kind: WHITESPACE_TOK, lit: " "}
	fullstopTok = ntToken{kind: FULLSTOP_TOK, lit: "."}
	lineFeedTok = ntToken{kind: LINEFEED_TOK, lit: "\n"}
	eofTok      = ntToken{kind: EOF_TOK}
)

type lexer struct {
	input                  string
	position, readPosition int
	char                   rune
}

func newLexer(s string) *lexer {
	return &lexer{
		input: s,
	}
}

func (l *lexer) nextToken() (ntToken, error) {
	if err := l.readChar(); err != nil {
		return ntToken{}, err
	}
	switch l.char {
	case '<':
		n, err := l.readIRI()
		return iriTok(n), err
	case ' ':
		return wspaceTok, nil
	case '.':
		return fullstopTok, nil
	case '\n':
		return lineFeedTok, nil
	case '"':
		n, err := l.readStringLiteral()
		return litTok(n), err
	case '^':
		if err := l.readChar(); err != nil {
			return ntToken{}, err
		}
		if l.char == 0 {
			return eofTok, nil
		}
		if l.char != '^' {
			panic(fmt.Sprintf("invalid datatype: expecting '^', got '%c': input [%s]", l.char, l.input))
		}
		if err := l.readChar(); err != nil {
			return ntToken{}, err
		}
		if l.char == 0 {
			return eofTok, nil
		}
		if l.char != '<' {
			panic(fmt.Sprintf("invalid datatype: expecting '<', got '%c'. Input: [%s]", l.char, l.input))
		}
		n, err := l.readIRI()
		return datatypeTok(n), err
	case '#':
		l.readChar()
		n, err := l.readComment()
		return commentTok(n), err
	case 0:
		return eofTok, nil
	default:
		return unknownTok(string(l.char)), nil
	}
}

func (l *lexer) readChar() error {
	var width int
	var err error
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char, width, err = decodeRune(l.input[l.readPosition:], l.readPosition)
		if err != nil {
			return err
		}
	}
	l.position = l.readPosition
	l.readPosition += width
	return nil
}

func (l *lexer) peekNextNonWithespaceChar() (found rune, count int, err error) {
	pos := l.readPosition
	if pos >= len(l.input) {
		return
	}
	var width int
	for {
		found, width, err = decodeRune(l.input[pos:], pos)
		if err != nil {
			return
		}
		count++
		if found == ' ' {
			pos = pos + width
			continue
		} else {
			return
		}
	}
}

func (l *lexer) readIRI() (string, error) {
	start := l.readPosition
	for {
		if err := l.readChar(); err != nil {
			return "", err
		}
		if l.char == '>' {
			peek, _, err := l.peekNextNonWithespaceChar()
			if err != nil {
				return "", err
			}
			if peek == 0 || peek == '<' || peek == '"' || peek == '.' {
				return l.input[start:l.position], nil
			}
		}
		if l.char == 0 {
			return "", nil
		}
	}
}

func (l *lexer) readStringLiteral() (string, error) {
	start := l.readPosition
	for {
		if err := l.readChar(); err != nil {
			return "", err
		}
		if l.char == '"' {
			peek, _, err := l.peekNextNonWithespaceChar()
			if err != nil {
				return "", err
			}
			if peek == 0 || peek == '.' || peek == '^' {
				return l.input[start:l.position], nil
			}
		}
		if l.char == 0 {
			return "", nil
		}
	}
}

func (l *lexer) readComment() (string, error) {
	pos := l.position
	for untilLineEnd(l.char) {
		if err := l.readChar(); err != nil {
			return "", err
		}
	}
	return l.input[pos:l.position], nil
}

func untilLineEnd(c rune) bool {
	return c != '\n' && c != 0
}

func decodeRune(s string, pos int) (r rune, width int, err error) {
	if s == "" {
		return 0, 0, nil
	}
	r, width = utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		switch width {
		case 0:
			err = fmt.Errorf("empty utf8 char starting at position %d", pos)
			return
		case 1:
			err = fmt.Errorf("invalid utf8 encoding starting at position %d", pos)
			return
		}
	}
	return
}
