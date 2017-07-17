package triplestore

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
)

type Decoder interface {
	Decode() ([]Triple, error)
}

// Use for retro compatibilty when changing file format on existing stores
func NewAutoDecoder(r io.Reader) Decoder {
	if IsBinaryFormat(r) {
		return NewBinaryDecoder(r)
	}
	return NewNTriplesDecoder(r)
}

func IsBinaryFormat(r io.Reader) bool {
	begin, err := bufio.NewReader(r).Peek(binary.Size(wordLength(0)) + 1)
	if err != nil {
		return false
	}
	dec := &binaryDecoder{r: bytes.NewReader(begin)}
	_, err = dec.readWord()
	return err == nil
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
		decodedObj.resource = string(resource)

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

type datasetDecoder struct {
	newDecoderFunc func(io.Reader) Decoder
	rs             []io.Reader
}

// A dataset is a basically a collection of RDFGraph.
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
