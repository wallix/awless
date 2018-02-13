package triplestore

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Decoder interface {
	Decode() ([]Triple, error)
}

type DecodeResult struct {
	Tri Triple
	Err error
}

type StreamDecoder interface {
	StreamDecode(context.Context) <-chan DecodeResult
}

// Use for retro compatibilty when changing file format on existing stores
func NewAutoDecoder(r io.Reader) Decoder {
	ok, newR := IsNTFormat(r)
	if ok {
		return NewLenientNTDecoder(newR)
	}
	return NewBinaryDecoder(newR)
}

// Loosely detect if a ntriples format contrary to a binary format
// Used for retro compatibilty when changing file format on existing stores
// Detecttion work with ntriples format flushed by this library (i.e. no comment, no spaces, ...)
func IsNTFormat(r io.Reader) (bool, io.Reader) {
	firstChar := make([]byte, 1)
	multi := io.MultiReader(bytes.NewReader(firstChar), r)
	if _, err := r.Read(firstChar); err != nil {
		return false, multi
	}
	return bytes.Equal(firstChar, []byte{'<'}), multi
}

func NewLenientNTDecoder(r io.Reader) Decoder {
	return &ntDecoder{r: r}
}

func NewLenientNTStreamDecoder(r io.Reader) StreamDecoder {
	return &ntDecoder{r: r}
}

type ntDecoder struct {
	r io.Reader
}

func (d *ntDecoder) Decode() ([]Triple, error) {
	return newLenientNTParser(d.r).Parse()
}

func (d *ntDecoder) StreamDecode(ctx context.Context) <-chan DecodeResult {
	decC := make(chan DecodeResult)

	go func() {
		defer close(decC)

		scanner := bufio.NewScanner(d.r)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if scanner.Scan() {
					tris, err := newLenientNTParser(strings.NewReader(scanner.Text())).Parse()
					if err != nil {
						decC <- DecodeResult{Err: err}
					} else if len(tris) == 1 {
						decC <- DecodeResult{Tri: tris[0]}
					}
				} else {
					if err := scanner.Err(); err != nil {
						decC <- DecodeResult{Err: err}
					}
					return
				}
			}
		}
	}()

	return decC
}

type binaryDecoder struct {
	r       io.Reader
	rc      io.ReadCloser // for stream decoding
	triples []Triple
}

func NewBinaryStreamDecoder(r io.ReadCloser) StreamDecoder {
	return &binaryDecoder{rc: r}
}

func (dec *binaryDecoder) StreamDecode(ctx context.Context) <-chan DecodeResult {
	decC := make(chan DecodeResult)

	go func() {
		defer close(decC)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				tri, done, err := decodeTriple(dec.rc)
				if done {
					return
				}
				decC <- DecodeResult{Tri: tri, Err: err}
			}
		}
	}()

	return decC
}

func NewBinaryDecoder(r io.Reader) Decoder {
	return &binaryDecoder{r: r}
}

func (dec *binaryDecoder) Decode() ([]Triple, error) {
	var out []Triple
	for {
		tri, done, err := decodeTriple(dec.r)
		if tri != nil {
			out = append(out, tri)
		}
		if done {
			break
		} else if err != nil {
			return out, err
		}
	}

	return out, nil
}

func decodeTriple(r io.Reader) (Triple, bool, error) {
	var isSubBNode bool
	err := binary.Read(r, binary.BigEndian, &isSubBNode)
	if err == io.EOF {
		return nil, true, nil
	} else if err != nil {
		return nil, false, fmt.Errorf("is subject bnode: %s", err)
	}

	sub, err := readWord(r)
	if err != nil {
		return nil, false, fmt.Errorf("subject: %s", err)
	}

	pred, err := readWord(r)
	if err != nil {
		return nil, false, fmt.Errorf("predicate: %s", err)
	}

	var objType uint8
	if err := binary.Read(r, binary.BigEndian, &objType); err != nil {
		return nil, false, fmt.Errorf("object type: %s", err)
	}

	var decodedObj object
	if objType == resourceTypeEncoding {
		resource, err := readWord(r)
		if err != nil {
			return nil, false, fmt.Errorf("resource: %s", err)
		}
		decodedObj.resource = string(resource)
	} else if objType == bnodeTypeEncoding {
		bnode, err := readWord(r)
		if err != nil {
			return nil, false, fmt.Errorf("bnode object: %s", err)
		}
		decodedObj.bnode = string(bnode)
		decodedObj.isBnode = true
	} else {
		decodedObj.isLit = true
		var decodedLiteral literal

		if objType == literalWithLangEncoding {
			lang, err := readWord(r)
			if err != nil {
				return nil, false, fmt.Errorf("lang: %s", err)
			}
			decodedLiteral.langtag = string(lang)
		} else {
			litType, err := readWord(r)
			if err != nil {
				return nil, false, fmt.Errorf("literate type: %s", err)
			}
			decodedLiteral.typ = XsdType(litType)
		}

		val, err := readWord(r)
		if err != nil {
			return nil, false, fmt.Errorf("literate: %s", err)
		}
		if decodedLiteral.typ == XsdString || objType == literalWithLangEncoding {
			decodedLiteral.val = unescapeStringLiteral(string(val))
		} else {
			decodedLiteral.val = string(val)
		}

		decodedObj.lit = decodedLiteral
	}

	return &triple{
		isSubBnode: isSubBNode,
		sub:        string(sub),
		pred:       string(pred),
		obj:        decodedObj,
	}, false, nil
}

func readWord(r io.Reader) ([]byte, error) {
	var len wordLength
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}

	word := make([]byte, len)
	if _, err := io.ReadFull(r, word); err != nil {
		return nil, fmt.Errorf("triplestore: binary: cannot decode word of length %d bytes: %s", len, err)
	}

	return word, nil
}

type datasetDecoder struct {
	newDecoderFunc func(io.Reader) Decoder
	rs             []io.Reader
}

// NewDatasetDecoder - a dataset is a basically a collection of RDFGraph.
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

var unescaper = strings.NewReplacer("\\n", "\n", "\\r", "\r")

func unescapeStringLiteral(s string) string {
	return unescaper.Replace(s)
}
