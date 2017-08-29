package triplestore

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
)

func TestStreamBinaryDecoding(t *testing.T) {
	var tris []Triple
	for i := 0; i < 10; i++ {
		tris = append(tris, SubjPred(fmt.Sprint(i), "digit").IntegerLiteral(i))
	}

	t.Run("handles done signal", func(t *testing.T) {
		var buf bytes.Buffer
		ctx, cancel := context.WithCancel(context.Background()) // will stop the decoding
		dec := NewBinaryStreamDecoder(ioutil.NopCloser(&buf))
		results := dec.StreamDecode(ctx)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range results {
				if r.Err != nil {
					t.Fatal(r.Err)
				}
			}
		}()
		cancel()
		wg.Wait()
	})

	t.Run("handles normal stream", func(t *testing.T) {
		var buf bytes.Buffer
		if err := NewBinaryEncoder(&buf).Encode(tris...); err != nil {
			t.Fatal(err)
		}

		dec := NewBinaryStreamDecoder(ioutil.NopCloser(&buf))
		results := dec.StreamDecode(context.Background())

		var all []Triple

		for r := range results {
			if r.Err != nil {
				t.Fatal(r.Err)
			}
			all = append(all, r.Tri)
		}

		if got, want := len(all), 10; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		s := NewSource()
		s.Add(all...)
		snap := s.Snapshot()

		for _, tri := range tris {
			if !snap.Contains(tri) {
				t.Fatalf("end result should contains triple %v", tri)
			}
		}
	})
}

func TestStreamNTriplesDecoding(t *testing.T) {
	var tris []Triple
	for i := 0; i < 10; i++ {
		tris = append(tris, SubjPred(fmt.Sprint(i), "digit").IntegerLiteral(i))
	}

	t.Run("handles done signal", func(t *testing.T) {
		var buf bytes.Buffer
		ctx, cancel := context.WithCancel(context.Background()) // will stop the decoding
		dec := NewLenientNTStreamDecoder(ioutil.NopCloser(&buf))
		results := dec.StreamDecode(ctx)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range results {
				if r.Err != nil {
					t.Fatal(r.Err)
				}
			}
		}()
		cancel()
		wg.Wait()
	})

	t.Run("handles normal stream", func(t *testing.T) {
		var buf bytes.Buffer
		if err := NewLenientNTEncoder(&buf).Encode(tris...); err != nil {
			t.Fatal(err)
		}

		dec := NewLenientNTStreamDecoder(ioutil.NopCloser(&buf))
		results := dec.StreamDecode(context.Background())

		var all []Triple

		for r := range results {
			if r.Err != nil {
				t.Fatal(r.Err)
			}
			all = append(all, r.Tri)
		}

		if got, want := len(all), 10; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		s := NewSource()
		s.Add(all...)
		snap := s.Snapshot()

		for _, tri := range tris {
			if !snap.Contains(tri) {
				t.Fatalf("end result should contains triple %v", tri)
			}
		}
	})
}

func TestStreamBinaryEncoding(t *testing.T) {
	t.Run("handles nil stream", func(t *testing.T) {
		enc := NewBinaryStreamEncoder(bytes.NewBuffer(nil))
		if err := enc.StreamEncode(context.Background(), nil); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("handles done stream", func(t *testing.T) {
		c := make(chan Triple)                                  // will make encoder block
		ctx, cancel := context.WithCancel(context.Background()) // will propagate encoding as done
		enc := NewBinaryStreamEncoder(bytes.NewBuffer(nil))

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := enc.StreamEncode(ctx, c); err != nil {
				t.Fatal(err)
			}
		}()
		cancel()
		wg.Wait()
	})

	var tris []Triple
	for i := 0; i < 10; i++ {
		tris = append(tris, SubjPred(fmt.Sprint(i), "digit").IntegerLiteral(i))
	}

	t.Run("handles normal stream", func(t *testing.T) {
		triC := make(chan Triple)
		go func() {
			for _, tri := range tris {
				triC <- tri
			}
			close(triC)
		}()

		var buf bytes.Buffer

		err := NewBinaryStreamEncoder(&buf).StreamEncode(context.Background(), triC)
		if err != nil {
			t.Fatal(err)
		}

		out, err := NewBinaryDecoder(&buf).Decode()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := len(out), 10; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		s := NewSource()
		s.Add(out...)
		snap := s.Snapshot()

		for _, tri := range tris {
			if !snap.Contains(tri) {
				t.Fatalf("end result should contains triple %v", tri)
			}
		}
	})
}

func TestStreamNTriplesEncoding(t *testing.T) {
	t.Run("handles nil stream", func(t *testing.T) {
		enc := NewLenientNTStreamEncoder(bytes.NewBuffer(nil))
		if err := enc.StreamEncode(context.Background(), nil); err != nil {
			t.Fatal(err)
		}
	})

	var tris []Triple
	for i := 0; i < 10; i++ {
		tris = append(tris, SubjPred(fmt.Sprint(i), "digit").IntegerLiteral(i))
	}

	t.Run("handles normal stream", func(t *testing.T) {
		triC := make(chan Triple)
		go func() {
			for _, tri := range tris {
				triC <- tri
			}
			close(triC)
		}()

		var buf bytes.Buffer

		err := NewLenientNTStreamEncoder(&buf).StreamEncode(context.Background(), triC)
		if err != nil {
			t.Fatal(err)
		}

		out, err := NewLenientNTDecoder(&buf).Decode()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := len(out), 10; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		s := NewSource()
		s.Add(out...)
		snap := s.Snapshot()

		for _, tri := range tris {
			if !snap.Contains(tri) {
				t.Fatalf("end result should contains triple %v", tri)
			}
		}
	})
}
