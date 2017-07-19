package triplestore

import (
	"bytes"
	"io"
	"text/scanner"
)

func NewNTriplesDecoder(r io.Reader) Decoder {
	return &ntDecoder{r: r}
}

type ntDecoder struct {
	r  io.Reader
	sc scanner.Scanner
}

const (
	LT  = '<'
	GT  = '>'
	DOT = '.'
	LF  = '\n'
	CR  = '\r'
)

func (d *ntDecoder) Decode() ([]Triple, error) {
	d.sc.Init(d.r)

	d.sc.Mode = scanner.ScanStrings | scanner.ScanRawStrings
	d.sc.Whitespace = 0

	var tris []Triple
	var nodeCount int
	var tok rune
	var sub, pred, obj string

	for tok != scanner.EOF {
		tok = d.sc.Scan()
		if tok == LT {
			nodeCount++
			if nodeCount == 1 {
				sub = d.parseNode()
			}
			if nodeCount == 2 {
				pred = d.parseNode()
			}
			if nodeCount == 3 {
				obj = d.parseNode()
			}
		}

		if tok == DOT {
			tris = append(tris, SubjPredRes(sub, pred, obj))
			nodeCount = 0
			sub, pred, obj = "", "", ""
		}
	}

	return tris, nil
}

func (d *ntDecoder) parseNode() string {
	var buff bytes.Buffer
	for tok := d.sc.Scan(); tok != GT; tok = d.sc.Scan() {
		buff.WriteRune(tok)
	}

	return buff.String()
}
