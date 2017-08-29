package triplestore

import (
	"bufio"
	"bytes"
	"regexp"
)

func contains(arr [][]byte, s []byte) bool {
	for _, a := range arr {
		if bytes.Equal(s, a) {
			return true
		}
	}
	return false
}

var (
	endOfLineComments = regexp.MustCompile(`(.*\.)\s+(#.*)`)
)

func cleanupNTriplesForComparison(b []byte) []byte {
	scn := bufio.NewScanner(bytes.NewReader(b))
	var cleaned bytes.Buffer
	for scn.Scan() {
		line := scn.Text()
		if empty, _ := regexp.MatchString(`^\s*$`, line); empty {
			continue
		}
		if comment, _ := regexp.MatchString(`^\s*#`, line); comment {
			continue
		}
		l := endOfLineComments.ReplaceAll([]byte(line), []byte("$1"))
		cleaned.Write(l)
		cleaned.WriteByte('\n')
	}

	return cleaned.Bytes()
}
