package triplestore

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// BenchmarkEncodingMemallocation-4   	   20000	     71052 ns/op	   27488 B/op	    1209 allocs/op
func BenchmarkEncodingMemallocation(b *testing.B) {
	var triples []Triple

	for i := 0; i < 100; i++ {
		triples = append(triples, SubjPred(fmt.Sprint(i), "digit").IntegerLiteral(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buff bytes.Buffer
		err := NewBinaryEncoder(&buff).Encode(triples...)
		if err != nil {
			b.Fatal(err)
		}
	}

}

//BenchmarkAllEncoding/binary-4         	                2000	    658993 ns/op
//BenchmarkAllEncoding/binary_streaming-4         	    1000	   1346275 ns/op
//BenchmarkAllEncoding/ntriples-4                 	    5000	    427614 ns/op
//BenchmarkAllEncoding/ntriples_streaming-4       	    2000	    988594 ns/op
//BenchmarkAllEncoding/ntriples_with_context-4    	    2000	   1346434 ns/op
func BenchmarkAllEncoding(b *testing.B) {
	var triples []Triple

	for i := 0; i < 1000; i++ {
		triples = append(triples, SubjPred(fmt.Sprint(i), "digit").IntegerLiteral(i))
	}

	b.ResetTimer()

	b.Run("binary", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buff bytes.Buffer
			if err := NewBinaryEncoder(&buff).Encode(triples...); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("binary streaming", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			triC := make(chan Triple)
			go tripleChan(triples, triC)
			b.StartTimer()
			var buff bytes.Buffer
			if err := NewBinaryStreamEncoder(&buff).StreamEncode(context.Background(), triC); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ntriples", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buff bytes.Buffer
			if err := NewNTriplesEncoder(&buff).Encode(triples...); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ntriples streaming", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			triC := make(chan Triple)
			go tripleChan(triples, triC)
			b.StartTimer()
			var buff bytes.Buffer
			if err := NewNTriplesStreamEncoder(&buff).StreamEncode(context.Background(), triC); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ntriples with context", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buff bytes.Buffer
			if err := NewNTriplesEncoderWithContext(&buff, RDFContext).Encode(triples...); err != nil {
				b.Fatal(err)
			}
		}
	})
}

//BenchmarkAllDecoding/binary-4                   	 3000000	       419 ns/op
//BenchmarkAllDecoding/binary_streaming-4         	  300000	      5607 ns/op
//BenchmarkAllDecoding/ntriples-4                 	 2000000	       596 ns/op
//BenchmarkAllDecoding/ntriples_streaming-4       	 1000000	      2089 ns/op
func BenchmarkAllDecoding(b *testing.B) {
	binaryFile, err := os.Open(filepath.Join("testdata", "bench", "decode_1.bin"))
	if err != nil {
		b.Fatal(err)
	}
	defer binaryFile.Close()

	b.ResetTimer()

	b.Run("binary", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := NewBinaryDecoder(binaryFile).Decode(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("binary streaming", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			results := NewBinaryStreamDecoder(binaryFile).StreamDecode(context.Background())
			for r := range results {
				if r.Err != nil {
					b.Fatal(r.Err)
				}
			}
		}
	})

	b.Run("ntriples", func(b *testing.B) {
		ntFile, err := os.Open(filepath.Join("testdata", "bench", "decode_1.nt"))
		if err != nil {
			b.Fatal(err)
		}
		defer ntFile.Close()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if _, err := NewNTriplesDecoder(ntFile).Decode(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ntriples streaming", func(b *testing.B) {
		ntFile, err := os.Open(filepath.Join("testdata", "bench", "decode_1.nt"))
		if err != nil {
			b.Fatal(err)
		}
		defer ntFile.Close()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			results := NewNTriplesStreamDecoder(ntFile).StreamDecode(context.Background())
			for r := range results {
				if r.Err != nil {
					b.Fatal(r.Err)
				}
			}
		}
	})
}

func tripleChan(triples []Triple, triC chan<- Triple) {
	for _, t := range triples {
		triC <- t
	}
	close(triC)
}
