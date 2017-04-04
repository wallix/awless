package triplestore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Encoder interface {
	Encode(tris ...Triple) error
}

type Decoder interface {
	Decode() ([]Triple, error)
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

type datasetDecoder struct {
	newDecoderFunc func(io.Reader) Decoder
	rs             []io.Reader
}

func NewDatasetDecoder(fn func(io.Reader) Decoder, readers ...io.Reader) Decoder {
	return &datasetDecoder{newDecoderFunc: fn, rs: readers}
}

func (dec *datasetDecoder) Decode() ([]Triple, error) {
	type result struct {
		err    error
		tris   []Triple
		reader io.Reader
	}

	results := make(chan *result, len(dec.rs))
	done := make(chan struct{})
	defer close(done)

	var wg sync.WaitGroup
	for _, reader := range dec.rs {
		wg.Add(1)
		go func(r io.Reader) {
			defer wg.Done()
			tris, err := dec.newDecoderFunc(r).Decode()
			select {
			case results <- &result{tris: tris, err: err, reader: r}:
			case <-done:
				return
			}
		}(reader)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var all []Triple
	for r := range results {
		if r.err != nil {
			switch rr := r.reader.(type) {
			case *os.File:
				return all, fmt.Errorf("file '%s': %s", rr.Name(), r.err)
			default:
				return all, r.err
			}
		}
		all = append(all, r.tris...)
	}

	return all, nil
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
		resID, _ := obj.ResourceID()
		binary.Write(&buff, binary.BigEndian, wordLength(len(resID)))
		buff.WriteString(resID)
	}

	return buff.Bytes(), nil
}

type binaryDecoder struct {
	r       io.Reader
	triples []Triple
}

func NewBinaryDecoder(r io.Reader) Decoder {
	return &binaryDecoder{r: r}
}

func (dec *binaryDecoder) Decode() ([]Triple, error) {
	for {
		done, err := dec.decodeTriple()
		if done {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return dec.triples, nil
}

func (dec *binaryDecoder) decodeTriple() (bool, error) {
	sub, err := dec.readWord()
	if err == io.EOF {
		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("subject: %s", err)
	}

	pred, err := dec.readWord()
	if err != nil {
		return false, fmt.Errorf("predicate: %s", err)
	}

	var objType uint8
	if err := binary.Read(dec.r, binary.BigEndian, &objType); err != nil {
		return false, fmt.Errorf("object type: %s", err)
	}

	var decodedObj object
	if objType == resourceTypeEncoding {
		resource, err := dec.readWord()
		if err != nil {
			return false, fmt.Errorf("resource: %s", err)
		}
		decodedObj.resourceID = string(resource)

	} else {
		decodedObj.isLit = true
		var decodedLiteral literal

		litType, err := dec.readWord()
		if err != nil {
			return false, fmt.Errorf("literate type: %s", err)
		}
		decodedLiteral.typ = XsdType(litType)

		val, err := dec.readWord()
		if err != nil {
			return false, fmt.Errorf("literate: %s", err)
		}

		decodedLiteral.val = string(val)
		decodedObj.lit = decodedLiteral
	}

	dec.triples = append(dec.triples, &triple{
		sub:  subject(string(sub)),
		pred: predicate(string(pred)),
		obj:  decodedObj,
	})

	return false, nil
}

func (dec *binaryDecoder) readWord() ([]byte, error) {
	var len wordLength
	if err := binary.Read(dec.r, binary.BigEndian, &len); err != nil {
		return nil, err
	}

	word := make([]byte, len)
	if _, err := io.ReadFull(dec.r, word); err != nil {
		return nil, fmt.Errorf("triplestore: binary: cannot decode word of length %d bytes: %s", len, err)
	}

	return word, nil
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
		if rid, ok := t.Object().ResourceID(); ok {
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

			buff.WriteString(fmt.Sprintf("\"%s\"%s", lit.Value(), namespace))
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
					id = uri + strings.TrimLeft(id, prefix)
					continue
				}
			}
		}
		if !strings.HasPrefix(id, "http") && enc.c.Base != "" {
			id = enc.c.Base + id
		}
	}
	return id
}
