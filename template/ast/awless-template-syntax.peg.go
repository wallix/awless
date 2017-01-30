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
		/* 2 Action <- <(('c' 'r' 'e' 'a' 't' 'e') / ('s' 't' 'a' 'r' 't') / ((&('c') ('c' 'h' 'e' 'c' 'k')) | (&('a') ('a' 't' 't' 'a' 'c' 'h')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('s') ('s' 't' 'o' 'p')) | (&('d') ('d' 'e' 'l' 'e' 't' 'e'))))> */
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
							if buffer[position] != rune('s') {
								goto l32
							}
							position++
							if buffer[position] != rune('t') {
								goto l32
							}
							position++
							if buffer[position] != rune('a') {
								goto l32
							}
							position++
							if buffer[position] != rune('r') {
								goto l32
							}
							position++
							if buffer[position] != rune('t') {
								goto l32
							}
							position++
							goto l30
						l32:
							position, tokenIndex = position30, tokenIndex30
							{
								switch buffer[position] {
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
								case 's':
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
								default:
									if buffer[position] != rune('d') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('l') {
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
									if buffer[position] != rune('e') {
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
					position35 := position
					{
						position36 := position
						{
							position37, tokenIndex37 := position, tokenIndex
							if buffer[position] != rune('v') {
								goto l38
							}
							position++
							if buffer[position] != rune('p') {
								goto l38
							}
							position++
							if buffer[position] != rune('c') {
								goto l38
							}
							position++
							goto l37
						l38:
							position, tokenIndex = position37, tokenIndex37
							if buffer[position] != rune('s') {
								goto l39
							}
							position++
							if buffer[position] != rune('u') {
								goto l39
							}
							position++
							if buffer[position] != rune('b') {
								goto l39
							}
							position++
							if buffer[position] != rune('n') {
								goto l39
							}
							position++
							if buffer[position] != rune('e') {
								goto l39
							}
							position++
							if buffer[position] != rune('t') {
								goto l39
							}
							position++
							goto l37
						l39:
							position, tokenIndex = position37, tokenIndex37
							if buffer[position] != rune('i') {
								goto l40
							}
							position++
							if buffer[position] != rune('n') {
								goto l40
							}
							position++
							if buffer[position] != rune('s') {
								goto l40
							}
							position++
							if buffer[position] != rune('t') {
								goto l40
							}
							position++
							if buffer[position] != rune('a') {
								goto l40
							}
							position++
							if buffer[position] != rune('n') {
								goto l40
							}
							position++
							if buffer[position] != rune('c') {
								goto l40
							}
							position++
							if buffer[position] != rune('e') {
								goto l40
							}
							position++
							goto l37
						l40:
							position, tokenIndex = position37, tokenIndex37
							if buffer[position] != rune('r') {
								goto l41
							}
							position++
							if buffer[position] != rune('o') {
								goto l41
							}
							position++
							if buffer[position] != rune('l') {
								goto l41
							}
							position++
							if buffer[position] != rune('e') {
								goto l41
							}
							position++
							goto l37
						l41:
							position, tokenIndex = position37, tokenIndex37
							if buffer[position] != rune('r') {
								goto l42
							}
							position++
							if buffer[position] != rune('o') {
								goto l42
							}
							position++
							if buffer[position] != rune('u') {
								goto l42
							}
							position++
							if buffer[position] != rune('t') {
								goto l42
							}
							position++
							if buffer[position] != rune('e') {
								goto l42
							}
							position++
							if buffer[position] != rune('t') {
								goto l42
							}
							position++
							if buffer[position] != rune('a') {
								goto l42
							}
							position++
							if buffer[position] != rune('b') {
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
							goto l37
						l42:
							position, tokenIndex = position37, tokenIndex37
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
					l37:
						add(ruleEntity, position36)
					}
					add(rulePegText, position35)
				}
				{
					add(ruleAction2, position)
				}
				{
					position45, tokenIndex45 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l45
					}
					{
						position47 := position
						{
							position50 := position
							{
								position51 := position
								if !_rules[ruleIdentifier]() {
									goto l45
								}
								add(rulePegText, position51)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l45
							}
							{
								position53 := position
								{
									position54, tokenIndex54 := position, tokenIndex
									{
										position56 := position
										{
											position57 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l55
											}
											position++
										l58:
											{
												position59, tokenIndex59 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l59
												}
												position++
												goto l58
											l59:
												position, tokenIndex = position59, tokenIndex59
											}
											if !matchDot() {
												goto l55
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l55
											}
											position++
										l60:
											{
												position61, tokenIndex61 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l61
												}
												position++
												goto l60
											l61:
												position, tokenIndex = position61, tokenIndex61
											}
											if !matchDot() {
												goto l55
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l55
											}
											position++
										l62:
											{
												position63, tokenIndex63 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l63
												}
												position++
												goto l62
											l63:
												position, tokenIndex = position63, tokenIndex63
											}
											if !matchDot() {
												goto l55
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l55
											}
											position++
										l64:
											{
												position65, tokenIndex65 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l65
												}
												position++
												goto l64
											l65:
												position, tokenIndex = position65, tokenIndex65
											}
											if buffer[position] != rune('/') {
												goto l55
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l55
											}
											position++
										l66:
											{
												position67, tokenIndex67 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l67
												}
												position++
												goto l66
											l67:
												position, tokenIndex = position67, tokenIndex67
											}
											add(ruleCidrValue, position57)
										}
										add(rulePegText, position56)
									}
									{
										add(ruleAction8, position)
									}
									goto l54
								l55:
									position, tokenIndex = position54, tokenIndex54
									{
										position70 := position
										{
											position71 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l69
											}
											position++
										l72:
											{
												position73, tokenIndex73 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l73
												}
												position++
												goto l72
											l73:
												position, tokenIndex = position73, tokenIndex73
											}
											if !matchDot() {
												goto l69
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l69
											}
											position++
										l74:
											{
												position75, tokenIndex75 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l75
												}
												position++
												goto l74
											l75:
												position, tokenIndex = position75, tokenIndex75
											}
											if !matchDot() {
												goto l69
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l69
											}
											position++
										l76:
											{
												position77, tokenIndex77 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l77
												}
												position++
												goto l76
											l77:
												position, tokenIndex = position77, tokenIndex77
											}
											if !matchDot() {
												goto l69
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l69
											}
											position++
										l78:
											{
												position79, tokenIndex79 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l79
												}
												position++
												goto l78
											l79:
												position, tokenIndex = position79, tokenIndex79
											}
											add(ruleIpValue, position71)
										}
										add(rulePegText, position70)
									}
									{
										add(ruleAction9, position)
									}
									goto l54
								l69:
									position, tokenIndex = position54, tokenIndex54
									{
										position82 := position
										{
											position83 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l81
											}
											position++
										l84:
											{
												position85, tokenIndex85 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l85
												}
												position++
												goto l84
											l85:
												position, tokenIndex = position85, tokenIndex85
											}
											if buffer[position] != rune('-') {
												goto l81
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l81
											}
											position++
										l86:
											{
												position87, tokenIndex87 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l87
												}
												position++
												goto l86
											l87:
												position, tokenIndex = position87, tokenIndex87
											}
											add(ruleIntRangeValue, position83)
										}
										add(rulePegText, position82)
									}
									{
										add(ruleAction10, position)
									}
									goto l54
								l81:
									position, tokenIndex = position54, tokenIndex54
									{
										position90 := position
										{
											position91 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l89
											}
											position++
										l92:
											{
												position93, tokenIndex93 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l93
												}
												position++
												goto l92
											l93:
												position, tokenIndex = position93, tokenIndex93
											}
											add(ruleIntValue, position91)
										}
										add(rulePegText, position90)
									}
									{
										add(ruleAction11, position)
									}
									goto l54
								l89:
									position, tokenIndex = position54, tokenIndex54
									{
										switch buffer[position] {
										case '$':
											{
												position96 := position
												if buffer[position] != rune('$') {
													goto l45
												}
												position++
												{
													position97 := position
													if !_rules[ruleIdentifier]() {
														goto l45
													}
													add(rulePegText, position97)
												}
												add(ruleRefValue, position96)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position99 := position
												if buffer[position] != rune('@') {
													goto l45
												}
												position++
												{
													position100 := position
													if !_rules[ruleIdentifier]() {
														goto l45
													}
													add(rulePegText, position100)
												}
												add(ruleAliasValue, position99)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position102 := position
												if buffer[position] != rune('{') {
													goto l45
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l45
												}
												{
													position103 := position
													if !_rules[ruleIdentifier]() {
														goto l45
													}
													add(rulePegText, position103)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l45
												}
												if buffer[position] != rune('}') {
													goto l45
												}
												position++
												add(ruleHoleValue, position102)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position105 := position
												{
													position106 := position
													{
														switch buffer[position] {
														case '_':
															if buffer[position] != rune('_') {
																goto l45
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l45
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l45
															}
															position++
															break
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l45
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l45
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l45
															}
															position++
															break
														}
													}

												l107:
													{
														position108, tokenIndex108 := position, tokenIndex
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l108
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l108
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l108
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l108
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l108
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l108
																}
																position++
																break
															}
														}

														goto l107
													l108:
														position, tokenIndex = position108, tokenIndex108
													}
													add(ruleStringValue, position106)
												}
												add(rulePegText, position105)
											}
											{
												add(ruleAction12, position)
											}
											break
										}
									}

								}
							l54:
								add(ruleValue, position53)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l45
							}
							add(ruleParam, position50)
						}
					l48:
						{
							position49, tokenIndex49 := position, tokenIndex
							{
								position112 := position
								{
									position113 := position
									if !_rules[ruleIdentifier]() {
										goto l49
									}
									add(rulePegText, position113)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l49
								}
								{
									position115 := position
									{
										position116, tokenIndex116 := position, tokenIndex
										{
											position118 := position
											{
												position119 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l117
												}
												position++
											l120:
												{
													position121, tokenIndex121 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l121
													}
													position++
													goto l120
												l121:
													position, tokenIndex = position121, tokenIndex121
												}
												if !matchDot() {
													goto l117
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l117
												}
												position++
											l122:
												{
													position123, tokenIndex123 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l123
													}
													position++
													goto l122
												l123:
													position, tokenIndex = position123, tokenIndex123
												}
												if !matchDot() {
													goto l117
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l117
												}
												position++
											l124:
												{
													position125, tokenIndex125 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l125
													}
													position++
													goto l124
												l125:
													position, tokenIndex = position125, tokenIndex125
												}
												if !matchDot() {
													goto l117
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l117
												}
												position++
											l126:
												{
													position127, tokenIndex127 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l127
													}
													position++
													goto l126
												l127:
													position, tokenIndex = position127, tokenIndex127
												}
												if buffer[position] != rune('/') {
													goto l117
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l117
												}
												position++
											l128:
												{
													position129, tokenIndex129 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l129
													}
													position++
													goto l128
												l129:
													position, tokenIndex = position129, tokenIndex129
												}
												add(ruleCidrValue, position119)
											}
											add(rulePegText, position118)
										}
										{
											add(ruleAction8, position)
										}
										goto l116
									l117:
										position, tokenIndex = position116, tokenIndex116
										{
											position132 := position
											{
												position133 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l131
												}
												position++
											l134:
												{
													position135, tokenIndex135 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l135
													}
													position++
													goto l134
												l135:
													position, tokenIndex = position135, tokenIndex135
												}
												if !matchDot() {
													goto l131
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l131
												}
												position++
											l136:
												{
													position137, tokenIndex137 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l137
													}
													position++
													goto l136
												l137:
													position, tokenIndex = position137, tokenIndex137
												}
												if !matchDot() {
													goto l131
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l131
												}
												position++
											l138:
												{
													position139, tokenIndex139 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l139
													}
													position++
													goto l138
												l139:
													position, tokenIndex = position139, tokenIndex139
												}
												if !matchDot() {
													goto l131
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l131
												}
												position++
											l140:
												{
													position141, tokenIndex141 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l141
													}
													position++
													goto l140
												l141:
													position, tokenIndex = position141, tokenIndex141
												}
												add(ruleIpValue, position133)
											}
											add(rulePegText, position132)
										}
										{
											add(ruleAction9, position)
										}
										goto l116
									l131:
										position, tokenIndex = position116, tokenIndex116
										{
											position144 := position
											{
												position145 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l143
												}
												position++
											l146:
												{
													position147, tokenIndex147 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l147
													}
													position++
													goto l146
												l147:
													position, tokenIndex = position147, tokenIndex147
												}
												if buffer[position] != rune('-') {
													goto l143
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l143
												}
												position++
											l148:
												{
													position149, tokenIndex149 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l149
													}
													position++
													goto l148
												l149:
													position, tokenIndex = position149, tokenIndex149
												}
												add(ruleIntRangeValue, position145)
											}
											add(rulePegText, position144)
										}
										{
											add(ruleAction10, position)
										}
										goto l116
									l143:
										position, tokenIndex = position116, tokenIndex116
										{
											position152 := position
											{
												position153 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l151
												}
												position++
											l154:
												{
													position155, tokenIndex155 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l155
													}
													position++
													goto l154
												l155:
													position, tokenIndex = position155, tokenIndex155
												}
												add(ruleIntValue, position153)
											}
											add(rulePegText, position152)
										}
										{
											add(ruleAction11, position)
										}
										goto l116
									l151:
										position, tokenIndex = position116, tokenIndex116
										{
											switch buffer[position] {
											case '$':
												{
													position158 := position
													if buffer[position] != rune('$') {
														goto l49
													}
													position++
													{
														position159 := position
														if !_rules[ruleIdentifier]() {
															goto l49
														}
														add(rulePegText, position159)
													}
													add(ruleRefValue, position158)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position161 := position
													if buffer[position] != rune('@') {
														goto l49
													}
													position++
													{
														position162 := position
														if !_rules[ruleIdentifier]() {
															goto l49
														}
														add(rulePegText, position162)
													}
													add(ruleAliasValue, position161)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position164 := position
													if buffer[position] != rune('{') {
														goto l49
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l49
													}
													{
														position165 := position
														if !_rules[ruleIdentifier]() {
															goto l49
														}
														add(rulePegText, position165)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l49
													}
													if buffer[position] != rune('}') {
														goto l49
													}
													position++
													add(ruleHoleValue, position164)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position167 := position
													{
														position168 := position
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l49
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l49
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l49
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l49
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l49
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l49
																}
																position++
																break
															}
														}

													l169:
														{
															position170, tokenIndex170 := position, tokenIndex
															{
																switch buffer[position] {
																case '_':
																	if buffer[position] != rune('_') {
																		goto l170
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
																		goto l170
																	}
																	position++
																	break
																case '-':
																	if buffer[position] != rune('-') {
																		goto l170
																	}
																	position++
																	break
																case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l170
																	}
																	position++
																	break
																case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																	if c := buffer[position]; c < rune('A') || c > rune('Z') {
																		goto l170
																	}
																	position++
																	break
																default:
																	if c := buffer[position]; c < rune('a') || c > rune('z') {
																		goto l170
																	}
																	position++
																	break
																}
															}

															goto l169
														l170:
															position, tokenIndex = position170, tokenIndex170
														}
														add(ruleStringValue, position168)
													}
													add(rulePegText, position167)
												}
												{
													add(ruleAction12, position)
												}
												break
											}
										}

									}
								l116:
									add(ruleValue, position115)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l49
								}
								add(ruleParam, position112)
							}
							goto l48
						l49:
							position, tokenIndex = position49, tokenIndex49
						}
						add(ruleParams, position47)
					}
					goto l46
				l45:
					position, tokenIndex = position45, tokenIndex45
				}
			l46:
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
			position177, tokenIndex177 := position, tokenIndex
			{
				position178 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l177
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l177
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l177
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l177
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l177
						}
						position++
						break
					}
				}

			l179:
				{
					position180, tokenIndex180 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l180
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l180
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l180
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l180
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l180
							}
							position++
							break
						}
					}

					goto l179
				l180:
					position, tokenIndex = position180, tokenIndex180
				}
				add(ruleIdentifier, position178)
			}
			return true
		l177:
			position, tokenIndex = position177, tokenIndex177
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
				position193 := position
			l194:
				{
					position195, tokenIndex195 := position, tokenIndex
					{
						position196 := position
						{
							position197, tokenIndex197 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l198
							}
							goto l197
						l198:
							position, tokenIndex = position197, tokenIndex197
							if !_rules[ruleEndOfLine]() {
								goto l195
							}
						}
					l197:
						add(ruleSpace, position196)
					}
					goto l194
				l195:
					position, tokenIndex = position195, tokenIndex195
				}
				add(ruleSpacing, position193)
			}
			return true
		},
		/* 19 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position200 := position
			l201:
				{
					position202, tokenIndex202 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l202
					}
					goto l201
				l202:
					position, tokenIndex = position202, tokenIndex202
				}
				add(ruleWhiteSpacing, position200)
			}
			return true
		},
		/* 20 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position203, tokenIndex203 := position, tokenIndex
			{
				position204 := position
				if !_rules[ruleWhitespace]() {
					goto l203
				}
			l205:
				{
					position206, tokenIndex206 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l206
					}
					goto l205
				l206:
					position, tokenIndex = position206, tokenIndex206
				}
				add(ruleMustWhiteSpacing, position204)
			}
			return true
		l203:
			position, tokenIndex = position203, tokenIndex203
			return false
		},
		/* 21 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position207, tokenIndex207 := position, tokenIndex
			{
				position208 := position
				if !_rules[ruleSpacing]() {
					goto l207
				}
				if buffer[position] != rune('=') {
					goto l207
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l207
				}
				add(ruleEqual, position208)
			}
			return true
		l207:
			position, tokenIndex = position207, tokenIndex207
			return false
		},
		/* 22 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 23 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position210, tokenIndex210 := position, tokenIndex
			{
				position211 := position
				{
					position212, tokenIndex212 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex = position212, tokenIndex212
					if buffer[position] != rune('\t') {
						goto l210
					}
					position++
				}
			l212:
				add(ruleWhitespace, position211)
			}
			return true
		l210:
			position, tokenIndex = position210, tokenIndex210
			return false
		},
		/* 24 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position214, tokenIndex214 := position, tokenIndex
			{
				position215 := position
				{
					position216, tokenIndex216 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l217
					}
					position++
					if buffer[position] != rune('\n') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex = position216, tokenIndex216
					if buffer[position] != rune('\n') {
						goto l218
					}
					position++
					goto l216
				l218:
					position, tokenIndex = position216, tokenIndex216
					if buffer[position] != rune('\r') {
						goto l214
					}
					position++
				}
			l216:
				add(ruleEndOfLine, position215)
			}
			return true
		l214:
			position, tokenIndex = position214, tokenIndex214
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
