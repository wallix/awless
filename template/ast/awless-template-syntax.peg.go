package ast

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleScript
	ruleStatement
	ruleAction
	ruleEntity
	ruleDeclaration
	ruleExpr
	ruleParams
	ruleParam
	ruleIdentifier
	ruleValue
	ruleStringValue
	ruleCidrValue
	ruleIpValue
	ruleIntValue
	ruleIntRangeValue
	ruleRefValue
	ruleAliasValue
	ruleHoleValue
	ruleSpacing
	ruleWhiteSpacing
	ruleMustWhiteSpacing
	ruleEqual
	ruleSpace
	ruleWhitespace
	ruleEndOfLine
	ruleEndOfFile
	rulePegText
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
)

var rul3s = [...]string{
	"Unknown",
	"Script",
	"Statement",
	"Action",
	"Entity",
	"Declaration",
	"Expr",
	"Params",
	"Param",
	"Identifier",
	"Value",
	"StringValue",
	"CidrValue",
	"IpValue",
	"IntValue",
	"IntRangeValue",
	"RefValue",
	"AliasValue",
	"HoleValue",
	"Spacing",
	"WhiteSpacing",
	"MustWhiteSpacing",
	"Equal",
	"Space",
	"Whitespace",
	"EndOfLine",
	"EndOfFile",
	"PegText",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type Peg struct {
	*AST

	Buffer string
	buffer []rune
	rules  [41]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *Peg) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *Peg) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *Peg
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *Peg) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *Peg) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.AddDeclarationIdentifier(text)
		case ruleAction1:
			p.AddAction(text)
		case ruleAction2:
			p.AddEntity(text)
		case ruleAction3:
			p.EndOfParams()
		case ruleAction4:
			p.AddParamKey(text)
		case ruleAction5:
			p.AddParamHoleValue(text)
		case ruleAction6:
			p.AddParamAliasValue(text)
		case ruleAction7:
			p.AddParamRefValue(text)
		case ruleAction8:
			p.AddParamCidrValue(text)
		case ruleAction9:
			p.AddParamIpValue(text)
		case ruleAction10:
			p.AddParamValue(text)
		case ruleAction11:
			p.AddParamIntValue(text)
		case ruleAction12:
			p.AddParamValue(text)

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *Peg) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Script <- <(Spacing Statement+ EndOfFile)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[ruleSpacing]() {
					goto l0
				}
				{
					position4 := position
					if !_rules[ruleSpacing]() {
						goto l0
					}
					{
						position5, tokenIndex5 := position, tokenIndex
						if !_rules[ruleExpr]() {
							goto l6
						}
						goto l5
					l6:
						position, tokenIndex = position5, tokenIndex5
						{
							position7 := position
							{
								position8 := position
								if !_rules[ruleIdentifier]() {
									goto l0
								}
								add(rulePegText, position8)
							}
							{
								add(ruleAction0, position)
							}
							if !_rules[ruleEqual]() {
								goto l0
							}
							if !_rules[ruleExpr]() {
								goto l0
							}
							add(ruleDeclaration, position7)
						}
					}
				l5:
					if !_rules[ruleSpacing]() {
						goto l0
					}
				l10:
					{
						position11, tokenIndex11 := position, tokenIndex
						if !_rules[ruleEndOfLine]() {
							goto l11
						}
						goto l10
					l11:
						position, tokenIndex = position11, tokenIndex11
					}
					add(ruleStatement, position4)
				}
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					{
						position12 := position
						if !_rules[ruleSpacing]() {
							goto l3
						}
						{
							position13, tokenIndex13 := position, tokenIndex
							if !_rules[ruleExpr]() {
								goto l14
							}
							goto l13
						l14:
							position, tokenIndex = position13, tokenIndex13
							{
								position15 := position
								{
									position16 := position
									if !_rules[ruleIdentifier]() {
										goto l3
									}
									add(rulePegText, position16)
								}
								{
									add(ruleAction0, position)
								}
								if !_rules[ruleEqual]() {
									goto l3
								}
								if !_rules[ruleExpr]() {
									goto l3
								}
								add(ruleDeclaration, position15)
							}
						}
					l13:
						if !_rules[ruleSpacing]() {
							goto l3
						}
					l18:
						{
							position19, tokenIndex19 := position, tokenIndex
							if !_rules[ruleEndOfLine]() {
								goto l19
							}
							goto l18
						l19:
							position, tokenIndex = position19, tokenIndex19
						}
						add(ruleStatement, position12)
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				{
					position20 := position
					{
						position21, tokenIndex21 := position, tokenIndex
						if !matchDot() {
							goto l21
						}
						goto l0
					l21:
						position, tokenIndex = position21, tokenIndex21
					}
					add(ruleEndOfFile, position20)
				}
				add(ruleScript, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Statement <- <(Spacing (Expr / Declaration) Spacing EndOfLine*)> */
		nil,
		/* 2 Action <- <(('c' 'r' 'e' 'a' 't' 'e') / ('d' 'e' 'l' 'e' 't' 'e') / ('s' 't' 'a' 'r' 't') / ((&('d') ('d' 'e' 't' 'a' 'c' 'h')) | (&('c') ('c' 'h' 'e' 'c' 'k')) | (&('a') ('a' 't' 't' 'a' 'c' 'h')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('s') ('s' 't' 'o' 'p'))))> */
		nil,
		/* 3 Entity <- <(('v' 'p' 'c') / ('s' 'u' 'b' 'n' 'e' 't') / ('i' 'n' 's' 't' 'a' 'n' 'c' 'e') / ('r' 'o' 'l' 'e') / ('r' 'o' 'u' 't' 'e' 't' 'a' 'b' 'l' 'e') / ((&('r') ('r' 'o' 'u' 't' 'e')) | (&('i') ('i' 'n' 't' 'e' 'r' 'n' 'e' 't' 'g' 'a' 't' 'e' 'w' 'a' 'y')) | (&('s') ('s' 'e' 'c' 'u' 'r' 'i' 't' 'y' 'g' 'r' 'o' 'u' 'p')) | (&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('t') ('t' 'a' 'g' 's')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e'))))> */
		nil,
		/* 4 Declaration <- <(<Identifier> Action0 Equal Expr)> */
		nil,
		/* 5 Expr <- <(<Action> Action1 MustWhiteSpacing <Entity> Action2 (MustWhiteSpacing Params)? Action3)> */
		func() bool {
			position26, tokenIndex26 := position, tokenIndex
			{
				position27 := position
				{
					position28 := position
					{
						position29 := position
						{
							position30, tokenIndex30 := position, tokenIndex
							if buffer[position] != rune('c') {
								goto l31
							}
							position++
							if buffer[position] != rune('r') {
								goto l31
							}
							position++
							if buffer[position] != rune('e') {
								goto l31
							}
							position++
							if buffer[position] != rune('a') {
								goto l31
							}
							position++
							if buffer[position] != rune('t') {
								goto l31
							}
							position++
							if buffer[position] != rune('e') {
								goto l31
							}
							position++
							goto l30
						l31:
							position, tokenIndex = position30, tokenIndex30
							if buffer[position] != rune('d') {
								goto l32
							}
							position++
							if buffer[position] != rune('e') {
								goto l32
							}
							position++
							if buffer[position] != rune('l') {
								goto l32
							}
							position++
							if buffer[position] != rune('e') {
								goto l32
							}
							position++
							if buffer[position] != rune('t') {
								goto l32
							}
							position++
							if buffer[position] != rune('e') {
								goto l32
							}
							position++
							goto l30
						l32:
							position, tokenIndex = position30, tokenIndex30
							if buffer[position] != rune('s') {
								goto l33
							}
							position++
							if buffer[position] != rune('t') {
								goto l33
							}
							position++
							if buffer[position] != rune('a') {
								goto l33
							}
							position++
							if buffer[position] != rune('r') {
								goto l33
							}
							position++
							if buffer[position] != rune('t') {
								goto l33
							}
							position++
							goto l30
						l33:
							position, tokenIndex = position30, tokenIndex30
							{
								switch buffer[position] {
								case 'd':
									if buffer[position] != rune('d') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('h') {
										goto l26
									}
									position++
									break
								case 'c':
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('h') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('k') {
										goto l26
									}
									position++
									break
								case 'a':
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('h') {
										goto l26
									}
									position++
									break
								case 'u':
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									if buffer[position] != rune('d') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									break
								default:
									if buffer[position] != rune('s') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
										goto l26
									}
									position++
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									break
								}
							}

						}
					l30:
						add(ruleAction, position29)
					}
					add(rulePegText, position28)
				}
				{
					add(ruleAction1, position)
				}
				if !_rules[ruleMustWhiteSpacing]() {
					goto l26
				}
				{
					position36 := position
					{
						position37 := position
						{
							position38, tokenIndex38 := position, tokenIndex
							if buffer[position] != rune('v') {
								goto l39
							}
							position++
							if buffer[position] != rune('p') {
								goto l39
							}
							position++
							if buffer[position] != rune('c') {
								goto l39
							}
							position++
							goto l38
						l39:
							position, tokenIndex = position38, tokenIndex38
							if buffer[position] != rune('s') {
								goto l40
							}
							position++
							if buffer[position] != rune('u') {
								goto l40
							}
							position++
							if buffer[position] != rune('b') {
								goto l40
							}
							position++
							if buffer[position] != rune('n') {
								goto l40
							}
							position++
							if buffer[position] != rune('e') {
								goto l40
							}
							position++
							if buffer[position] != rune('t') {
								goto l40
							}
							position++
							goto l38
						l40:
							position, tokenIndex = position38, tokenIndex38
							if buffer[position] != rune('i') {
								goto l41
							}
							position++
							if buffer[position] != rune('n') {
								goto l41
							}
							position++
							if buffer[position] != rune('s') {
								goto l41
							}
							position++
							if buffer[position] != rune('t') {
								goto l41
							}
							position++
							if buffer[position] != rune('a') {
								goto l41
							}
							position++
							if buffer[position] != rune('n') {
								goto l41
							}
							position++
							if buffer[position] != rune('c') {
								goto l41
							}
							position++
							if buffer[position] != rune('e') {
								goto l41
							}
							position++
							goto l38
						l41:
							position, tokenIndex = position38, tokenIndex38
							if buffer[position] != rune('r') {
								goto l42
							}
							position++
							if buffer[position] != rune('o') {
								goto l42
							}
							position++
							if buffer[position] != rune('l') {
								goto l42
							}
							position++
							if buffer[position] != rune('e') {
								goto l42
							}
							position++
							goto l38
						l42:
							position, tokenIndex = position38, tokenIndex38
							if buffer[position] != rune('r') {
								goto l43
							}
							position++
							if buffer[position] != rune('o') {
								goto l43
							}
							position++
							if buffer[position] != rune('u') {
								goto l43
							}
							position++
							if buffer[position] != rune('t') {
								goto l43
							}
							position++
							if buffer[position] != rune('e') {
								goto l43
							}
							position++
							if buffer[position] != rune('t') {
								goto l43
							}
							position++
							if buffer[position] != rune('a') {
								goto l43
							}
							position++
							if buffer[position] != rune('b') {
								goto l43
							}
							position++
							if buffer[position] != rune('l') {
								goto l43
							}
							position++
							if buffer[position] != rune('e') {
								goto l43
							}
							position++
							goto l38
						l43:
							position, tokenIndex = position38, tokenIndex38
							{
								switch buffer[position] {
								case 'r':
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
										goto l26
									}
									position++
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									break
								case 'i':
									if buffer[position] != rune('i') {
										goto l26
									}
									position++
									if buffer[position] != rune('n') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('n') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('g') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('w') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('y') {
										goto l26
									}
									position++
									break
								case 's':
									if buffer[position] != rune('s') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('i') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('y') {
										goto l26
									}
									position++
									if buffer[position] != rune('g') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
										goto l26
									}
									position++
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									break
								case 'k':
									if buffer[position] != rune('k') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('y') {
										goto l26
									}
									position++
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('i') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									break
								case 'p':
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
										goto l26
									}
									position++
									if buffer[position] != rune('l') {
										goto l26
									}
									position++
									if buffer[position] != rune('i') {
										goto l26
									}
									position++
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('y') {
										goto l26
									}
									position++
									break
								case 'g':
									if buffer[position] != rune('g') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
										goto l26
									}
									position++
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									break
								case 'u':
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('s') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									break
								case 't':
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('g') {
										goto l26
									}
									position++
									if buffer[position] != rune('s') {
										goto l26
									}
									position++
									break
								default:
									if buffer[position] != rune('v') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
										goto l26
									}
									position++
									if buffer[position] != rune('l') {
										goto l26
									}
									position++
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('m') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									break
								}
							}

						}
					l38:
						add(ruleEntity, position37)
					}
					add(rulePegText, position36)
				}
				{
					add(ruleAction2, position)
				}
				{
					position46, tokenIndex46 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l46
					}
					{
						position48 := position
						{
							position51 := position
							{
								position52 := position
								if !_rules[ruleIdentifier]() {
									goto l46
								}
								add(rulePegText, position52)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l46
							}
							{
								position54 := position
								{
									position55, tokenIndex55 := position, tokenIndex
									{
										position57 := position
										{
											position58 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l56
											}
											position++
										l59:
											{
												position60, tokenIndex60 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l60
												}
												position++
												goto l59
											l60:
												position, tokenIndex = position60, tokenIndex60
											}
											if !matchDot() {
												goto l56
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l56
											}
											position++
										l61:
											{
												position62, tokenIndex62 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l62
												}
												position++
												goto l61
											l62:
												position, tokenIndex = position62, tokenIndex62
											}
											if !matchDot() {
												goto l56
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l56
											}
											position++
										l63:
											{
												position64, tokenIndex64 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l64
												}
												position++
												goto l63
											l64:
												position, tokenIndex = position64, tokenIndex64
											}
											if !matchDot() {
												goto l56
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l56
											}
											position++
										l65:
											{
												position66, tokenIndex66 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l66
												}
												position++
												goto l65
											l66:
												position, tokenIndex = position66, tokenIndex66
											}
											if buffer[position] != rune('/') {
												goto l56
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l56
											}
											position++
										l67:
											{
												position68, tokenIndex68 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l68
												}
												position++
												goto l67
											l68:
												position, tokenIndex = position68, tokenIndex68
											}
											add(ruleCidrValue, position58)
										}
										add(rulePegText, position57)
									}
									{
										add(ruleAction8, position)
									}
									goto l55
								l56:
									position, tokenIndex = position55, tokenIndex55
									{
										position71 := position
										{
											position72 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l70
											}
											position++
										l73:
											{
												position74, tokenIndex74 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l74
												}
												position++
												goto l73
											l74:
												position, tokenIndex = position74, tokenIndex74
											}
											if !matchDot() {
												goto l70
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l70
											}
											position++
										l75:
											{
												position76, tokenIndex76 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l76
												}
												position++
												goto l75
											l76:
												position, tokenIndex = position76, tokenIndex76
											}
											if !matchDot() {
												goto l70
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l70
											}
											position++
										l77:
											{
												position78, tokenIndex78 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l78
												}
												position++
												goto l77
											l78:
												position, tokenIndex = position78, tokenIndex78
											}
											if !matchDot() {
												goto l70
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l70
											}
											position++
										l79:
											{
												position80, tokenIndex80 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l80
												}
												position++
												goto l79
											l80:
												position, tokenIndex = position80, tokenIndex80
											}
											add(ruleIpValue, position72)
										}
										add(rulePegText, position71)
									}
									{
										add(ruleAction9, position)
									}
									goto l55
								l70:
									position, tokenIndex = position55, tokenIndex55
									{
										position83 := position
										{
											position84 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l82
											}
											position++
										l85:
											{
												position86, tokenIndex86 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l86
												}
												position++
												goto l85
											l86:
												position, tokenIndex = position86, tokenIndex86
											}
											if buffer[position] != rune('-') {
												goto l82
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l82
											}
											position++
										l87:
											{
												position88, tokenIndex88 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l88
												}
												position++
												goto l87
											l88:
												position, tokenIndex = position88, tokenIndex88
											}
											add(ruleIntRangeValue, position84)
										}
										add(rulePegText, position83)
									}
									{
										add(ruleAction10, position)
									}
									goto l55
								l82:
									position, tokenIndex = position55, tokenIndex55
									{
										position91 := position
										{
											position92 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l90
											}
											position++
										l93:
											{
												position94, tokenIndex94 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l94
												}
												position++
												goto l93
											l94:
												position, tokenIndex = position94, tokenIndex94
											}
											add(ruleIntValue, position92)
										}
										add(rulePegText, position91)
									}
									{
										add(ruleAction11, position)
									}
									goto l55
								l90:
									position, tokenIndex = position55, tokenIndex55
									{
										switch buffer[position] {
										case '$':
											{
												position97 := position
												if buffer[position] != rune('$') {
													goto l46
												}
												position++
												{
													position98 := position
													if !_rules[ruleIdentifier]() {
														goto l46
													}
													add(rulePegText, position98)
												}
												add(ruleRefValue, position97)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position100 := position
												if buffer[position] != rune('@') {
													goto l46
												}
												position++
												{
													position101 := position
													if !_rules[ruleIdentifier]() {
														goto l46
													}
													add(rulePegText, position101)
												}
												add(ruleAliasValue, position100)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position103 := position
												if buffer[position] != rune('{') {
													goto l46
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l46
												}
												{
													position104 := position
													if !_rules[ruleIdentifier]() {
														goto l46
													}
													add(rulePegText, position104)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l46
												}
												if buffer[position] != rune('}') {
													goto l46
												}
												position++
												add(ruleHoleValue, position103)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position106 := position
												{
													position107 := position
													{
														switch buffer[position] {
														case '_':
															if buffer[position] != rune('_') {
																goto l46
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l46
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l46
															}
															position++
															break
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l46
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l46
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l46
															}
															position++
															break
														}
													}

												l108:
													{
														position109, tokenIndex109 := position, tokenIndex
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l109
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l109
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l109
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l109
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l109
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l109
																}
																position++
																break
															}
														}

														goto l108
													l109:
														position, tokenIndex = position109, tokenIndex109
													}
													add(ruleStringValue, position107)
												}
												add(rulePegText, position106)
											}
											{
												add(ruleAction12, position)
											}
											break
										}
									}

								}
							l55:
								add(ruleValue, position54)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l46
							}
							add(ruleParam, position51)
						}
					l49:
						{
							position50, tokenIndex50 := position, tokenIndex
							{
								position113 := position
								{
									position114 := position
									if !_rules[ruleIdentifier]() {
										goto l50
									}
									add(rulePegText, position114)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l50
								}
								{
									position116 := position
									{
										position117, tokenIndex117 := position, tokenIndex
										{
											position119 := position
											{
												position120 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l118
												}
												position++
											l121:
												{
													position122, tokenIndex122 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l122
													}
													position++
													goto l121
												l122:
													position, tokenIndex = position122, tokenIndex122
												}
												if !matchDot() {
													goto l118
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l118
												}
												position++
											l123:
												{
													position124, tokenIndex124 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l124
													}
													position++
													goto l123
												l124:
													position, tokenIndex = position124, tokenIndex124
												}
												if !matchDot() {
													goto l118
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l118
												}
												position++
											l125:
												{
													position126, tokenIndex126 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l126
													}
													position++
													goto l125
												l126:
													position, tokenIndex = position126, tokenIndex126
												}
												if !matchDot() {
													goto l118
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l118
												}
												position++
											l127:
												{
													position128, tokenIndex128 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l128
													}
													position++
													goto l127
												l128:
													position, tokenIndex = position128, tokenIndex128
												}
												if buffer[position] != rune('/') {
													goto l118
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l118
												}
												position++
											l129:
												{
													position130, tokenIndex130 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l130
													}
													position++
													goto l129
												l130:
													position, tokenIndex = position130, tokenIndex130
												}
												add(ruleCidrValue, position120)
											}
											add(rulePegText, position119)
										}
										{
											add(ruleAction8, position)
										}
										goto l117
									l118:
										position, tokenIndex = position117, tokenIndex117
										{
											position133 := position
											{
												position134 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l132
												}
												position++
											l135:
												{
													position136, tokenIndex136 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l136
													}
													position++
													goto l135
												l136:
													position, tokenIndex = position136, tokenIndex136
												}
												if !matchDot() {
													goto l132
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l132
												}
												position++
											l137:
												{
													position138, tokenIndex138 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l138
													}
													position++
													goto l137
												l138:
													position, tokenIndex = position138, tokenIndex138
												}
												if !matchDot() {
													goto l132
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l132
												}
												position++
											l139:
												{
													position140, tokenIndex140 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l140
													}
													position++
													goto l139
												l140:
													position, tokenIndex = position140, tokenIndex140
												}
												if !matchDot() {
													goto l132
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l132
												}
												position++
											l141:
												{
													position142, tokenIndex142 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l142
													}
													position++
													goto l141
												l142:
													position, tokenIndex = position142, tokenIndex142
												}
												add(ruleIpValue, position134)
											}
											add(rulePegText, position133)
										}
										{
											add(ruleAction9, position)
										}
										goto l117
									l132:
										position, tokenIndex = position117, tokenIndex117
										{
											position145 := position
											{
												position146 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l144
												}
												position++
											l147:
												{
													position148, tokenIndex148 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l148
													}
													position++
													goto l147
												l148:
													position, tokenIndex = position148, tokenIndex148
												}
												if buffer[position] != rune('-') {
													goto l144
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l144
												}
												position++
											l149:
												{
													position150, tokenIndex150 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l150
													}
													position++
													goto l149
												l150:
													position, tokenIndex = position150, tokenIndex150
												}
												add(ruleIntRangeValue, position146)
											}
											add(rulePegText, position145)
										}
										{
											add(ruleAction10, position)
										}
										goto l117
									l144:
										position, tokenIndex = position117, tokenIndex117
										{
											position153 := position
											{
												position154 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l152
												}
												position++
											l155:
												{
													position156, tokenIndex156 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l156
													}
													position++
													goto l155
												l156:
													position, tokenIndex = position156, tokenIndex156
												}
												add(ruleIntValue, position154)
											}
											add(rulePegText, position153)
										}
										{
											add(ruleAction11, position)
										}
										goto l117
									l152:
										position, tokenIndex = position117, tokenIndex117
										{
											switch buffer[position] {
											case '$':
												{
													position159 := position
													if buffer[position] != rune('$') {
														goto l50
													}
													position++
													{
														position160 := position
														if !_rules[ruleIdentifier]() {
															goto l50
														}
														add(rulePegText, position160)
													}
													add(ruleRefValue, position159)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position162 := position
													if buffer[position] != rune('@') {
														goto l50
													}
													position++
													{
														position163 := position
														if !_rules[ruleIdentifier]() {
															goto l50
														}
														add(rulePegText, position163)
													}
													add(ruleAliasValue, position162)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position165 := position
													if buffer[position] != rune('{') {
														goto l50
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l50
													}
													{
														position166 := position
														if !_rules[ruleIdentifier]() {
															goto l50
														}
														add(rulePegText, position166)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l50
													}
													if buffer[position] != rune('}') {
														goto l50
													}
													position++
													add(ruleHoleValue, position165)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position168 := position
													{
														position169 := position
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l50
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l50
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l50
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l50
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l50
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l50
																}
																position++
																break
															}
														}

													l170:
														{
															position171, tokenIndex171 := position, tokenIndex
															{
																switch buffer[position] {
																case '_':
																	if buffer[position] != rune('_') {
																		goto l171
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
																		goto l171
																	}
																	position++
																	break
																case '-':
																	if buffer[position] != rune('-') {
																		goto l171
																	}
																	position++
																	break
																case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l171
																	}
																	position++
																	break
																case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																	if c := buffer[position]; c < rune('A') || c > rune('Z') {
																		goto l171
																	}
																	position++
																	break
																default:
																	if c := buffer[position]; c < rune('a') || c > rune('z') {
																		goto l171
																	}
																	position++
																	break
																}
															}

															goto l170
														l171:
															position, tokenIndex = position171, tokenIndex171
														}
														add(ruleStringValue, position169)
													}
													add(rulePegText, position168)
												}
												{
													add(ruleAction12, position)
												}
												break
											}
										}

									}
								l117:
									add(ruleValue, position116)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l50
								}
								add(ruleParam, position113)
							}
							goto l49
						l50:
							position, tokenIndex = position50, tokenIndex50
						}
						add(ruleParams, position48)
					}
					goto l47
				l46:
					position, tokenIndex = position46, tokenIndex46
				}
			l47:
				{
					add(ruleAction3, position)
				}
				add(ruleExpr, position27)
			}
			return true
		l26:
			position, tokenIndex = position26, tokenIndex26
			return false
		},
		/* 6 Params <- <Param+> */
		nil,
		/* 7 Param <- <(<Identifier> Action4 Equal Value WhiteSpacing)> */
		nil,
		/* 8 Identifier <- <((&('.') '.') | (&('_') '_') | (&('-') '-') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position178, tokenIndex178 := position, tokenIndex
			{
				position179 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l178
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l178
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l178
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l178
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l178
						}
						position++
						break
					}
				}

			l180:
				{
					position181, tokenIndex181 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l181
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l181
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l181
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l181
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l181
							}
							position++
							break
						}
					}

					goto l180
				l181:
					position, tokenIndex = position181, tokenIndex181
				}
				add(ruleIdentifier, position179)
			}
			return true
		l178:
			position, tokenIndex = position178, tokenIndex178
			return false
		},
		/* 9 Value <- <((<CidrValue> Action8) / (<IpValue> Action9) / (<IntRangeValue> Action10) / (<IntValue> Action11) / ((&('$') (RefValue Action7)) | (&('@') (AliasValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action12))))> */
		nil,
		/* 10 StringValue <- <((&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		nil,
		/* 11 CidrValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+ '/' [0-9]+)> */
		nil,
		/* 12 IpValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+)> */
		nil,
		/* 13 IntValue <- <[0-9]+> */
		nil,
		/* 14 IntRangeValue <- <([0-9]+ '-' [0-9]+)> */
		nil,
		/* 15 RefValue <- <('$' <Identifier>)> */
		nil,
		/* 16 AliasValue <- <('@' <Identifier>)> */
		nil,
		/* 17 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 18 Spacing <- <Space*> */
		func() bool {
			{
				position194 := position
			l195:
				{
					position196, tokenIndex196 := position, tokenIndex
					{
						position197 := position
						{
							position198, tokenIndex198 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l199
							}
							goto l198
						l199:
							position, tokenIndex = position198, tokenIndex198
							if !_rules[ruleEndOfLine]() {
								goto l196
							}
						}
					l198:
						add(ruleSpace, position197)
					}
					goto l195
				l196:
					position, tokenIndex = position196, tokenIndex196
				}
				add(ruleSpacing, position194)
			}
			return true
		},
		/* 19 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position201 := position
			l202:
				{
					position203, tokenIndex203 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l203
					}
					goto l202
				l203:
					position, tokenIndex = position203, tokenIndex203
				}
				add(ruleWhiteSpacing, position201)
			}
			return true
		},
		/* 20 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position204, tokenIndex204 := position, tokenIndex
			{
				position205 := position
				if !_rules[ruleWhitespace]() {
					goto l204
				}
			l206:
				{
					position207, tokenIndex207 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l207
					}
					goto l206
				l207:
					position, tokenIndex = position207, tokenIndex207
				}
				add(ruleMustWhiteSpacing, position205)
			}
			return true
		l204:
			position, tokenIndex = position204, tokenIndex204
			return false
		},
		/* 21 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position208, tokenIndex208 := position, tokenIndex
			{
				position209 := position
				if !_rules[ruleSpacing]() {
					goto l208
				}
				if buffer[position] != rune('=') {
					goto l208
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l208
				}
				add(ruleEqual, position209)
			}
			return true
		l208:
			position, tokenIndex = position208, tokenIndex208
			return false
		},
		/* 22 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 23 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position211, tokenIndex211 := position, tokenIndex
			{
				position212 := position
				{
					position213, tokenIndex213 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l214
					}
					position++
					goto l213
				l214:
					position, tokenIndex = position213, tokenIndex213
					if buffer[position] != rune('\t') {
						goto l211
					}
					position++
				}
			l213:
				add(ruleWhitespace, position212)
			}
			return true
		l211:
			position, tokenIndex = position211, tokenIndex211
			return false
		},
		/* 24 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position215, tokenIndex215 := position, tokenIndex
			{
				position216 := position
				{
					position217, tokenIndex217 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l218
					}
					position++
					if buffer[position] != rune('\n') {
						goto l218
					}
					position++
					goto l217
				l218:
					position, tokenIndex = position217, tokenIndex217
					if buffer[position] != rune('\n') {
						goto l219
					}
					position++
					goto l217
				l219:
					position, tokenIndex = position217, tokenIndex217
					if buffer[position] != rune('\r') {
						goto l215
					}
					position++
				}
			l217:
				add(ruleEndOfLine, position216)
			}
			return true
		l215:
			position, tokenIndex = position215, tokenIndex215
			return false
		},
		/* 25 EndOfFile <- <!.> */
		nil,
		nil,
		/* 28 Action0 <- <{ p.AddDeclarationIdentifier(text) }> */
		nil,
		/* 29 Action1 <- <{ p.AddAction(text) }> */
		nil,
		/* 30 Action2 <- <{ p.AddEntity(text) }> */
		nil,
		/* 31 Action3 <- <{ p.EndOfParams() }> */
		nil,
		/* 32 Action4 <- <{ p.AddParamKey(text) }> */
		nil,
		/* 33 Action5 <- <{  p.AddParamHoleValue(text) }> */
		nil,
		/* 34 Action6 <- <{  p.AddParamAliasValue(text) }> */
		nil,
		/* 35 Action7 <- <{  p.AddParamRefValue(text) }> */
		nil,
		/* 36 Action8 <- <{ p.AddParamCidrValue(text) }> */
		nil,
		/* 37 Action9 <- <{ p.AddParamIpValue(text) }> */
		nil,
		/* 38 Action10 <- <{ p.AddParamValue(text) }> */
		nil,
		/* 39 Action11 <- <{ p.AddParamIntValue(text) }> */
		nil,
		/* 40 Action12 <- <{ p.AddParamValue(text) }> */
		nil,
	}
	p.rules = _rules
}
