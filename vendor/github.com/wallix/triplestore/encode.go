package triplestore

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/url"
	"strings"
)

type Encoder interface {
	Encode(tris ...Triple) error
}

type StreamEncoder interface {
	StreamEncode(context.Context, <-chan Triple) error
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

type wordLength uint32

const (
	resourceTypeEncoding    = uint8(0)
	literalTypeEncoding     = uint8(1)
	bnodeTypeEncoding       = uint8(2)
	literalWithLangEncoding = uint8(3)
)

type binaryEncoder struct {
	w io.Writer
}

func NewBinaryStreamEncoder(w io.Writer) StreamEncoder {
	return &binaryEncoder{w}
}

func NewBinaryEncoder(w io.Writer) Encoder {
	return &binaryEncoder{w}
}

func (enc *binaryEncoder) StreamEncode(ctx context.Context, triples <-chan Triple) error {
	if triples == nil {
		return nil
	}
	var buf bytes.Buffer
	for {
		select {
		case tri, ok := <-triples:
			if !ok {
				return nil
			}
			if err := enc.writeTriple(tri, &buf); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (enc *binaryEncoder) Encode(tris ...Triple) error {
	var buf bytes.Buffer
	for _, t := range tris {
		if err := enc.writeTriple(t, &buf); err != nil {
			return err
		}
	}
	return nil
}

func (enc *binaryEncoder) writeTriple(t Triple, buf *bytes.Buffer) error {
	if err := encodeBinTriple(t, buf); err != nil {
		return err
	}
	if _, err := enc.w.Write(buf.Bytes()); err != nil {
		return err
	}
	buf.Reset()
	return nil
}

func encodeBinTriple(t Triple, buff *bytes.Buffer) error {
	sub, pred := t.Subject(), t.Predicate()

	binary.Write(buff, binary.BigEndian, t.(*triple).isSubBnode)

	binary.Write(buff, binary.BigEndian, wordLength(len(sub)))
	buff.WriteString(sub)

	binary.Write(buff, binary.BigEndian, wordLength(len(pred)))
	buff.WriteString(pred)

	obj := t.Object()
	if lit, isLit := obj.Literal(); isLit {
		if lang := lit.Lang(); len(lang) > 0 {
			binary.Write(buff, binary.BigEndian, literalWithLangEncoding)
			binary.Write(buff, binary.BigEndian, wordLength(len(lang)))
			buff.WriteString(string(lang))
		} else {
			binary.Write(buff, binary.BigEndian, literalTypeEncoding)
			typ := lit.Type()
			binary.Write(buff, binary.BigEndian, wordLength(len(typ)))
			buff.WriteString(string(typ))
		}

		litVal := lit.Value()
		if lit.Type() == XsdString {
			litVal = escapeStringLiteral(litVal)
		}
		binary.Write(buff, binary.BigEndian, wordLength(len(litVal)))
		buff.WriteString(litVal)
	} else if bnode, isBnode := obj.Bnode(); isBnode {
		binary.Write(buff, binary.BigEndian, bnodeTypeEncoding)
		binary.Write(buff, binary.BigEndian, wordLength(len(bnode)))
		buff.WriteString(bnode)
	} else {
		binary.Write(buff, binary.BigEndian, resourceTypeEncoding)
		res, _ := obj.Resource()
		binary.Write(buff, binary.BigEndian, wordLength(len(res)))
		buff.WriteString(res)
	}

	return nil
}

type ntriplesEncoder struct {
	w io.Writer
	c *Context
}

func NewLenientNTStreamEncoder(w io.Writer) StreamEncoder {
	return &ntriplesEncoder{w: w}
}

func NewLenientNTEncoder(w io.Writer) Encoder {
	return &ntriplesEncoder{w: w}
}

func NewLenientNTEncoderWithContext(w io.Writer, c *Context) Encoder {
	return &ntriplesEncoder{w: w, c: c}
}

func (enc *ntriplesEncoder) StreamEncode(ctx context.Context, triples <-chan Triple) error {
	if triples == nil {
		return nil
	}
	var buf bytes.Buffer
	finalWrite := func() error {
		_, err := enc.w.Write(buf.Bytes())
		return err
	}
	for {
		select {
		case tri, ok := <-triples:
			if !ok {
				return finalWrite()
			}
			encodeNTriple(tri, enc.c, &buf)
		case <-ctx.Done():
			return finalWrite()
		}
	}
}

func (enc *ntriplesEncoder) Encode(tris ...Triple) error {
	var buff bytes.Buffer

	for _, t := range tris {
		encodeNTriple(t, enc.c, &buff)
	}
	_, err := enc.w.Write(buff.Bytes())
	return err
}

func encodeNTriple(t Triple, ctx *Context, buff *bytes.Buffer) {
	var sub string
	if tt := t.(*triple); tt.isSubBnode {
		sub = "_:" + buildIRI(ctx, t.Subject())
	} else {
		sub = "<" + buildIRI(ctx, t.Subject()) + ">"
	}
	buff.WriteString(sub + " <" + buildIRI(ctx, t.Predicate()) + "> ")

	if bnode, isBnode := t.Object().Bnode(); isBnode {
		buff.WriteString("_:" + bnode)
	} else {
		if rid, ok := t.Object().Resource(); ok {
			buff.WriteString("<" + buildIRI(ctx, rid) + ">")
		} else if lit, ok := t.Object().Literal(); ok {
			if lit.Lang() != "" {
				buff.WriteString("\"" + escapeStringLiteral(lit.Value()) + "\"@" + lit.Lang())
			} else {
				switch lit.Type() {
				case XsdString:
					// namespace empty as per spec
					buff.WriteString("\"" + escapeStringLiteral(lit.Value()) + "\"")
				default:
					if ctx != nil {
						if _, ok := ctx.Prefixes["xsd"]; ok {
							buff.WriteString("\"" + lit.Value() + "\"^^<" + lit.Type().NTriplesNamespaced() + ">")
						}
					} else {
						buff.WriteString("\"" + lit.Value() + "\"^^<" + string(lit.Type()) + ">")
					}
				}
			}
		}
	}
	buff.Write([]byte(" .\n"))
}

func buildIRI(ctx *Context, id string) string {
	if ctx != nil {
		if ctx.Prefixes != nil {
			for k, uri := range ctx.Prefixes {
				prefix := k + ":"
				if strings.HasPrefix(id, prefix) {
					id = uri + url.QueryEscape(strings.TrimPrefix(id, prefix))
					continue
				}
			}
		}
		if !strings.HasPrefix(id, "http") && ctx.Base != "" {
			id = ctx.Base + url.QueryEscape(id)
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

var escaper = strings.NewReplacer("\n", "\\n", "\r", "\\r")

func escapeStringLiteral(s string) string {
	return escaper.Replace(s)
}
