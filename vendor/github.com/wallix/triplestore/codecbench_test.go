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

//BenchmarkAllEncoding/binary-4      	                2000	    609710 ns/op	  295264 B/op	   11012 allocs/op
//BenchmarkAllEncoding/binary_streaming-4         	    2000	   1046498 ns/op	  295269 B/op	   11012 allocs/op
//BenchmarkAllEncoding/ntriples-4                 	    3000	    530518 ns/op	  529136 B/op	    4014 allocs/op
//BenchmarkAllEncoding/ntriples_streaming-4       	    2000	    988511 ns/op	  529170 B/op	    4015 allocs/op
//BenchmarkAllEncoding/ntriples_with_context-4    	    1000	   1272959 ns/op	  764839 B/op	    8015 allocs/op
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
			if err := NewLenientNTEncoder(&buff).Encode(triples...); err != nil {
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
			if err := NewLenientNTStreamEncoder(&buff).StreamEncode(context.Background(), triC); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ntriples with context", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buff bytes.Buffer
			if err := NewLenientNTEncoderWithContext(&buff, RDFContext).Encode(triples...); err != nil {
				b.Fatal(err)
			}
		}
	})
}

//BenchmarkAllDecoding/binary-4                   	 3000000	       493 ns/op	      72 B/op	       3 allocs/op
//BenchmarkAllDecoding/binary_streaming-4         	  300000	      4519 ns/op	     168 B/op	       4 allocs/op
//BenchmarkAllDecoding/ntriples-4                 	 1000000	      1079 ns/op	    4112 B/op	       2 allocs/op
//BenchmarkAllDecoding/ntriples_streaming-4       	 1000000	      1928 ns/op	    4212 B/op	       3 allocs/op
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
			if _, err := NewLenientNTDecoder(ntFile).Decode(); err != nil {
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
			results := NewLenientNTStreamDecoder(ntFile).StreamDecode(context.Background())
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
