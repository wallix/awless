package reflectline

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/chzyer/readline"
)

type Shell struct {
	node ShellNode
}

func NewShell(node ShellNode) *Shell {
	return &Shell{node: node}
}

func (s *Shell) Run() error {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    s.node,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		args, err := s.parseArgs([]rune(line))
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		err = s.node.Call(args)
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
	}
	return nil
}

func (s *Shell) parseArgs(line []rune) (args []string, err error) {
	start := -1
	i := 0
	for i < len(line) {
		r := line[i]
		if unicode.IsSpace(r) {
			if start != -1 {
				args = append(args, string(line[start:i]))
				start = -1
			}
			i = i + 1
			continue
		}
		if start == -1 {
			start = i
		}
		if r == '\'' {
			j := strings.IndexRune(string(line[i+1:len(line)]), '\'')
			if j == -1 {
				err = fmt.Errorf("unterminated simple quoted string")
				return
			}
			i = j + i + 1
			continue
		}
		if r == '"' {
			j := indexEndOfDoubleQuote(line[i+1 : len(line)])
			if j == -1 {
				err = fmt.Errorf("unterminated double quote string")
				return
			}
			i = j + i + 2
			continue
		}
		i = i + 1
	}
	if start != -1 {
		args = append(args, string(line[start:i]))
	}
	return
}
