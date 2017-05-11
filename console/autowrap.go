package console

import (
	"bytes"
	"strings"
)

const wrapingChars = " \n.,?;:/+=*-_°)(\"!@#"

type autoWraper struct {
	maxWidth     int
	wrappingChar string
}

func (aw autoWraper) Wrap(input string) string {
	wrappingChar := aw.wrappingChar
	if wrappingChar == "" {
		wrappingChar = " "
	}
	splits := strings.Split(input, wrappingChar)
	var output bytes.Buffer
	for i, split := range splits {
		if len(split) <= aw.maxWidth {
			output.WriteString(split)
		} else {
			var toWrite bytes.Buffer
			var buff bytes.Buffer
			for _, c := range split {
				if toWrite.Len() != 0 {
					if buff.Len()+toWrite.Len() >= aw.maxWidth {
						output.WriteString(toWrite.String())
						toWrite.Reset()
						output.WriteString(wrappingChar)
					}
				}
				if buff.Len() >= aw.maxWidth {
					output.WriteString(buff.String())
					buff.Reset()
					output.WriteString(wrappingChar)
				}
				buff.WriteRune(c)
				if strings.ContainsRune(wrapingChars, c) {
					toWrite.WriteString(buff.String())
					buff.Reset()
				}
			}
			output.WriteString(toWrite.String())
			output.WriteString(buff.String())
		}
		if i < len(splits)-1 {
			output.WriteString(wrappingChar)
		}
	}
	return output.String()
}
