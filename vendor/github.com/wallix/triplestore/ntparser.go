package triplestore

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

type lenientNTParser struct {
	r io.Reader
}

func newLenientNTParser(r io.Reader) *lenientNTParser {
	return &lenientNTParser{r: r}
}

func (p *lenientNTParser) Parse() (out []Triple, err error) {
	var count int
	scanner := bufio.NewScanner(p.r)
	for scanner.Scan() {
		count++
		line := bytes.TrimLeft(scanner.Bytes(), " \t")
		if len(line) < 1 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		t, terr := parseTriple(line)
		if terr != nil {
			return out, fmt.Errorf("lenient parsing: line %d: %s", count, terr)
		}
		out = append(out, t)
	}

	err = scanner.Err()
	return
}

func parseTriple(b []byte) (Triple, error) {
	tBuilder := new(tripleBuilder)
	var err error
	if bytes.HasPrefix(b, []byte("_:")) {
		if tBuilder.sub, b, err = parseBNodeSubject(b[2:]); err != nil {
			return nil, err
		}
		tBuilder.isSubBnode = true
	} else if bytes.HasPrefix(b, []byte("<")) {
		if tBuilder.sub, b, err = parseIRISubject(b[1:]); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid subject in %s", b)
	}

	if bytes.HasPrefix(b, []byte{'<'}) {
		if tBuilder.pred, b, err = parsePredicate(b[1:]); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid predicate in %s", b)
	}

	if bytes.HasPrefix(b, []byte{'<'}) {
		obj, _, err := parseIRIObject(b[1:])
		return tBuilder.Resource(obj), err
	} else if bytes.HasPrefix(b, []byte("_:")) {
		obj, _, err := parseBNodeObject(b[2:])
		return tBuilder.Bnode(obj), err
	} else if bytes.HasPrefix(b, []byte{'"'}) {
		lit, b, err := parseLiteralObject(b[1:])
		if err != nil {
			return nil, err
		}
		if bytes.HasPrefix(b, []byte("^^<")) {
			dtype, _, err := parseIRIObject(b[3:])
			obj := object{
				isLit: true,
				lit: literal{
					typ: XsdType(dtype),
					val: lit,
				},
			}
			return tBuilder.Object(obj), err
		} else if bytes.HasPrefix(b, []byte{'@'}) {
			lang, _, err := parseLangtag(b[1:])
			return tBuilder.StringLiteralWithLang(unescapeStringLiteral(lit), lang), err
		} else {
			return tBuilder.StringLiteral(unescapeStringLiteral(lit)), err
		}
	} else {
		return nil, errors.New("invalid object")
	}
}

func parseLangtag(b []byte) (string, []byte, error) {
	var index int
	for {
		r, size, err, eol := decode(b[index:])
		if err != nil {
			return "", nil, err
		}
		if eol {
			return "", nil, errors.New("invalid language tag")
		}
		index += size

		if r == '.' {
			if found, advance := peekNext(b[index:]); found == '#' || found == 0 {
				return string(b[:index-1]), b[index-1+advance:], nil
			}
		}

		if r == ' ' {
			if found, advance := peekNext(b[index:]); found == '.' {
				return string(b[:index-1]), b[index+advance:], nil
			}
		}
	}
}

func parseLiteralObject(b []byte) (string, []byte, error) {
	var index int
	for {
		r, size, err, eol := decode(b[index:])
		if err != nil {
			return "", nil, err
		}
		if eol {
			return "", nil, errors.New("invalid literal object")
		}
		index += size

		if r == '"' {
			if found, advance, other := doublePeekNext(b[index:]); (found == '.' && other == '#') || (found == '.' && other == 0) || (found == '^' && other == '^') || found == '@' {
				return string(b[:index-1]), b[index+advance:], nil
			}
		}
	}
}

func parsePredicate(b []byte) (string, []byte, error) {
	var index int
	for {
		r, size, err, eol := decode(b[index:])
		if err != nil {
			return "", nil, err
		}
		if eol {
			return "", nil, errors.New("invalid predicate")
		}
		index += size

		if r == '>' {
			if found, advance := peekNext(b[index:]); found == '<' || found == '"' || found == '_' {
				return string(b[:index-1]), b[index+advance:], nil
			}
		}
	}
}

func parseIRIObject(b []byte) (string, []byte, error) {
	var index int
	for {
		r, size, err, eol := decode(b[index:])
		if err != nil {
			return "", nil, err
		}
		if eol {
			return "", nil, errors.New("invalid IRI object")
		}
		index += size

		if r == '>' {
			if found, advance := peekNext(b[index:]); found == '.' {
				return string(b[:index-1]), b[index+advance:], nil
			}
		}
	}
}

func parseIRISubject(b []byte) (string, []byte, error) {
	var index int
	for {
		r, size, err, eol := decode(b[index:])
		if err != nil {
			return "", nil, err
		}
		if eol {
			return "", nil, errors.New("invalid IRI subject")
		}
		index += size

		if r == '>' {
			if found, advance := peekNext(b[index:]); found == '<' {
				return string(b[:index-1]), b[index+advance:], nil
			}
		}
	}
}

func parseBNodeObject(b []byte) (string, []byte, error) {
	var index int
	for {
		r, size, err, eol := decode(b[index:])
		if err != nil {
			return "", nil, err
		}
		if eol {
			return "", nil, errors.New("invalid bnode object")
		}
		index += size

		if r == '.' {
			if found, advance := peekNext(b[index:]); found == '#' || found == 0 {
				return string(b[:index-1]), b[index-1+advance:], nil
			}
		}

		if r == ' ' || r == '\t' {
			if found, advance := peekNext(b[index:]); found == '.' {
				return string(b[:index-1]), b[index+advance:], nil
			}
		}
	}
}

func parseBNodeSubject(b []byte) (string, []byte, error) {
	var index int
	for {
		r, size, err, eol := decode(b[index:])
		if err != nil {
			return "", nil, err
		}
		if eol {
			return "", nil, errors.New("invalid bnode subject")
		}
		index += size

		if r == '<' {
			return string(b[:index-1]), b[index-1:], nil
		}
		if r == ' ' || r == '\t' {
			if found, advance := peekNext(b[index:]); found == '<' {
				return string(b[:index-1]), b[index+advance:], nil
			}
		}
	}
}

func decode(b []byte) (rune, int, error, bool) {
	r, size := utf8.DecodeRune(b)
	if r == utf8.RuneError && size == 1 {
		return r, 0, errors.New("invalid utf8 encoding"), false
	}
	if r == utf8.RuneError && size == 0 {
		return 0, 0, nil, true
	}
	return r, size, nil, false
}

func peekNext(b []byte) (found rune, advance int) {
	for {
		r, size := utf8.DecodeRune(b[advance:])
		if r == utf8.RuneError {
			return 0, 0
		}

		if r != ' ' && r != '\t' {
			found = r
			return
		}
		advance += size
	}
}

func doublePeekNext(b []byte) (first rune, advance int, second rune) {
	first, advance = peekNext(b)
	if n := advance + utf8.RuneLen(first); len(b) >= n {
		second, _ = peekNext(b[n:])
	}
	return first, advance, second
}
