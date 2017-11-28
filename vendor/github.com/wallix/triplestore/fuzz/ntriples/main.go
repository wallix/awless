package ntriples

import "github.com/wallix/triplestore"
import "bytes"

func Fuzz(data []byte) int {
	dec := triplestore.NewLenientNTDecoder(bytes.NewReader(data))
	if _, err := dec.Decode(); err != nil {
		return 0
	}
	return 1
}
