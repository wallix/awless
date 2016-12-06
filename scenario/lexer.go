//go:generate stringer -type=Token

package scenario

import (
	"bufio"
	"strings"

	"github.com/wallix/awless/scenario/driver"
)

type Lexer struct {
	sc  *bufio.Scanner
	pos int
}

func (l *Lexer) ParseScenario(raw string) *driver.Scenario {
	scanner := bufio.NewScanner(strings.NewReader(raw))

	scen := &driver.Scenario{}

	for scanner.Scan() {
		line := l.parseLine(scanner.Text())
		scen.Lines = append(scen.Lines, line)
	}

	return scen
}

func (l *Lexer) parseLine(raw string) *driver.Line {
	l.sc = bufio.NewScanner(strings.NewReader(raw))
	l.sc.Split(bufio.ScanWords)
	l.pos = 0

	lin := &driver.Line{
		Params: make(map[driver.Token]interface{}),
	}

	var lastKey driver.Token

	for l.sc.Scan() {
		token := l.sc.Text()

		switch l.pos {
		case 0:
			lin.Action = driver.TokenFromString(token)
		case 1:
			lin.Resource = driver.TokenFromString(token)
		default:
			switch {
			case l.pos%2 == 0:
				lastKey = driver.TokenFromString(token)
			default:
				lin.Params[lastKey] = token
			}
		}

		l.pos++
	}

	return lin
}
