package triplestore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type Encoder interface {
	Encode(tris ...Triple) error
}

func NewContext() *Context {
	return &Context{Prefixes: make(map[string]string)}
}

type Context struct {
	Base     string
	Prefixes map[string]string
}

var RDFContext = &Context{
	Prefixes: map[string]string{
		"xsd":  "http://www.w3.org/2001/XMLSchema#",
		"rdf":  "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
		"rdfs": "http://www.w3.org/2000/01/rdf-schema#",
	},
}

type binaryEncoder struct {
	w io.Writer
}

type wordLength uint32

const (
	resourceTypeEncoding = uint8(0)
	literalTypeEncoding  = uint8(1)
)

func NewBinaryEncoder(w io.Writer) Encoder {
	return &binaryEncoder{w}
}

func (enc *binaryEncoder) Encode(tris ...Triple) error {
	for _, t := range tris {
		b, err := encodeTriple(t)
		if err != nil {
			return err
		}

		if _, err := enc.w.Write(b); err != nil {
			return err
		}
	}

	return nil
}

func encodeTriple(t Triple) ([]byte, error) {
	sub, pred := t.Subject(), t.Predicate()

	var buff bytes.Buffer

	binary.Write(&buff, binary.BigEndian, wordLength(len(sub)))
	buff.WriteString(sub)

	binary.Write(&buff, binary.BigEndian, wordLength(len(pred)))
	buff.WriteString(pred)

	obj := t.Object()
	if lit, isLit := obj.Literal(); isLit {
		binary.Write(&buff, binary.BigEndian, literalTypeEncoding)
		typ := lit.Type()
		binary.Write(&buff, binary.BigEndian, wordLength(len(typ)))
		buff.WriteString(string(typ))

		litVal := lit.Value()
		binary.Write(&buff, binary.BigEndian, wordLength(len(litVal)))
		buff.WriteString(litVal)
	} else {
		binary.Write(&buff, binary.BigEndian, resourceTypeEncoding)
		res, _ := obj.Resource()
		binary.Write(&buff, binary.BigEndian, wordLength(len(res)))
		buff.WriteString(res)
	}

	return buff.Bytes(), nil
}

type ntriplesEncoder struct {
	w io.Writer
	c *Context
}

func NewNTriplesEncoder(w io.Writer) Encoder {
	return &ntriplesEncoder{w: w}
}

func NewNTriplesEncoderWithContext(w io.Writer, c *Context) Encoder {
	return &ntriplesEncoder{w: w, c: c}
}

func (enc *ntriplesEncoder) Encode(tris ...Triple) error {
	var buff bytes.Buffer
	for _, t := range tris {
		buff.WriteString(fmt.Sprintf("<%s> <%s> ", enc.buildIRI(t.Subject()), enc.buildIRI(t.Predicate())))
		if rid, ok := t.Object().Resource(); ok {
			buff.WriteString(fmt.Sprintf("<%s>", enc.buildIRI(rid)))
		}
		if lit, ok := t.Object().Literal(); ok {
			var namespace string
			switch lit.Type() {
			case XsdString:
				// namespace empty as per spec
			default:
				namespace = lit.Type().NTriplesNamespaced()
			}

			buff.WriteString(fmt.Sprintf("%s%s", strconv.QuoteToASCII(lit.Value()), namespace))
		}
		buff.WriteString(" .\n")
	}

	_, err := enc.w.Write(buff.Bytes())
	return err
}

func (enc *ntriplesEncoder) buildIRI(id string) string {
	if enc.c != nil {
		if enc.c.Prefixes != nil {
			for k, uri := range enc.c.Prefixes {
				prefix := k + ":"
				if strings.HasPrefix(id, prefix) {
					id = uri + url.QueryEscape(strings.TrimPrefix(id, prefix))
					continue
				}
			}
		}
		if !strings.HasPrefix(id, "http") && enc.c.Base != "" {
			id = enc.c.Base + url.QueryEscape(id)
		}
	}
	return id
}

type dotGraphEncoder struct {
	pred string
	w    io.Writer
}

func NewDotGraphEncoder(w io.Writer, predicate string) Encoder {
	return &dotGraphEncoder{w: w, pred: predicate}
}

func (dg *dotGraphEncoder) Encode(tris ...Triple) error {
	src := NewSource()
	src.Add(tris...)

	snap := src.Snapshot()
	all := snap.WithPredicate(dg.pred)

	queryDone := make(map[string][]string)

	getTypes := func(ref string) ([]string, bool) {
		if all, ok := queryDone[ref]; ok {
			return all, true
		} else {
			fresh := snap.WithSubjPred(ref, "rdf:type")
			for _, typ := range fresh {
				val, _ := typ.Object().Resource()
				queryDone[ref] = append(queryDone[ref], val)
			}
			return queryDone[ref], false
		}
	}

	fmt.Fprintf(dg.w, "digraph \"%s\" {\n", dg.pred)
	for _, tri := range all {
		sub := tri.Subject()
		res, ok := tri.Object().Resource()
		if ok {
			fmt.Fprintf(dg.w, "\"%s\" -> \"%s\";\n", sub, res)

			subTypes, done := getTypes(sub)
			if !done {
				for _, typ := range subTypes {
					fmt.Fprintf(dg.w, "\"%s\" [label=\"%s<%s>\"];\n", sub, sub, typ)
				}
			}

			resTypes, done := getTypes(res)
			if !done {
				for _, typ := range resTypes {
					fmt.Fprintf(dg.w, "\"%s\" [label=\"%s<%s>\"];\n", res, res, typ)
				}
			}
		}
	}

	fmt.Fprintf(dg.w, "}")

	return nil
}
