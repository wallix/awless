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
		/* 3 Entity <- <(('v' 'p' 'c') / ('s' 'u' 'b' 'n' 'e' 't') / ((&('s') ('s' 'e' 'c' 'u' 'r' 'i' 't' 'y' 'g' 'r' 'o' 'u' 'p')) | (&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('r') ('r' 'o' 'l' 'e')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('t') ('t' 'a' 'g' 's')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e')) | (&('i') ('i' 'n' 's' 't' 'a' 'n' 'c' 'e'))))> */
		nil,
		/* 4 Declaration <- <(<Identifier> Action0 Equal Expr)> */
		nil,
		/* 5 Expr <- <(<Action> Action1 MustWhiteSpacing <Entity> Action2 (MustWhiteSpacing Params)?)> */
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
							{
								switch buffer[position] {
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
								case 'r':
									if buffer[position] != rune('r') {
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
									if buffer[position] != rune('e') {
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
								case 'v':
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
								default:
									if buffer[position] != rune('i') {
										goto l26
									}
									position++
									if buffer[position] != rune('n') {
										goto l26
									}
									position++
									if buffer[position] != rune('s') {
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
									if buffer[position] != rune('n') {
										goto l26
									}
									position++
									if buffer[position] != rune('c') {
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
					position42, tokenIndex42 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l42
					}
					{
						position44 := position
						{
							position47 := position
							{
								position48 := position
								if !_rules[ruleIdentifier]() {
									goto l42
								}
								add(rulePegText, position48)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l42
							}
							{
								position50 := position
								{
									position51, tokenIndex51 := position, tokenIndex
									{
										position53 := position
										{
											position54 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l52
											}
											position++
										l55:
											{
												position56, tokenIndex56 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l56
												}
												position++
												goto l55
											l56:
												position, tokenIndex = position56, tokenIndex56
											}
											if !matchDot() {
												goto l52
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l52
											}
											position++
										l57:
											{
												position58, tokenIndex58 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l58
												}
												position++
												goto l57
											l58:
												position, tokenIndex = position58, tokenIndex58
											}
											if !matchDot() {
												goto l52
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l52
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
												goto l52
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l52
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
											if buffer[position] != rune('/') {
												goto l52
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l52
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
											add(ruleCidrValue, position54)
										}
										add(rulePegText, position53)
									}
									{
										add(ruleAction8, position)
									}
									goto l51
								l52:
									position, tokenIndex = position51, tokenIndex51
									{
										position67 := position
										{
											position68 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l66
											}
											position++
										l69:
											{
												position70, tokenIndex70 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l70
												}
												position++
												goto l69
											l70:
												position, tokenIndex = position70, tokenIndex70
											}
											if !matchDot() {
												goto l66
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l66
											}
											position++
										l71:
											{
												position72, tokenIndex72 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l72
												}
												position++
												goto l71
											l72:
												position, tokenIndex = position72, tokenIndex72
											}
											if !matchDot() {
												goto l66
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l66
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
												goto l66
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l66
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
											add(ruleIpValue, position68)
										}
										add(rulePegText, position67)
									}
									{
										add(ruleAction9, position)
									}
									goto l51
								l66:
									position, tokenIndex = position51, tokenIndex51
									{
										position79 := position
										{
											position80 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l78
											}
											position++
										l81:
											{
												position82, tokenIndex82 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l82
												}
												position++
												goto l81
											l82:
												position, tokenIndex = position82, tokenIndex82
											}
											if buffer[position] != rune('-') {
												goto l78
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l78
											}
											position++
										l83:
											{
												position84, tokenIndex84 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l84
												}
												position++
												goto l83
											l84:
												position, tokenIndex = position84, tokenIndex84
											}
											add(ruleIntRangeValue, position80)
										}
										add(rulePegText, position79)
									}
									{
										add(ruleAction10, position)
									}
									goto l51
								l78:
									position, tokenIndex = position51, tokenIndex51
									{
										position87 := position
										{
											position88 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l86
											}
											position++
										l89:
											{
												position90, tokenIndex90 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l90
												}
												position++
												goto l89
											l90:
												position, tokenIndex = position90, tokenIndex90
											}
											add(ruleIntValue, position88)
										}
										add(rulePegText, position87)
									}
									{
										add(ruleAction11, position)
									}
									goto l51
								l86:
									position, tokenIndex = position51, tokenIndex51
									{
										switch buffer[position] {
										case '$':
											{
												position93 := position
												if buffer[position] != rune('$') {
													goto l42
												}
												position++
												{
													position94 := position
													if !_rules[ruleIdentifier]() {
														goto l42
													}
													add(rulePegText, position94)
												}
												add(ruleRefValue, position93)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position96 := position
												if buffer[position] != rune('@') {
													goto l42
												}
												position++
												{
													position97 := position
													if !_rules[ruleIdentifier]() {
														goto l42
													}
													add(rulePegText, position97)
												}
												add(ruleAliasValue, position96)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position99 := position
												if buffer[position] != rune('{') {
													goto l42
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l42
												}
												{
													position100 := position
													if !_rules[ruleIdentifier]() {
														goto l42
													}
													add(rulePegText, position100)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l42
												}
												if buffer[position] != rune('}') {
													goto l42
												}
												position++
												add(ruleHoleValue, position99)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position102 := position
												{
													position103 := position
													{
														switch buffer[position] {
														case '_':
															if buffer[position] != rune('_') {
																goto l42
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l42
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l42
															}
															position++
															break
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l42
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l42
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l42
															}
															position++
															break
														}
													}

												l104:
													{
														position105, tokenIndex105 := position, tokenIndex
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l105
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l105
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l105
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l105
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l105
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l105
																}
																position++
																break
															}
														}

														goto l104
													l105:
														position, tokenIndex = position105, tokenIndex105
													}
													add(ruleStringValue, position103)
												}
												add(rulePegText, position102)
											}
											{
												add(ruleAction12, position)
											}
											break
										}
									}

								}
							l51:
								add(ruleValue, position50)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l42
							}
							add(ruleParam, position47)
						}
					l45:
						{
							position46, tokenIndex46 := position, tokenIndex
							{
								position109 := position
								{
									position110 := position
									if !_rules[ruleIdentifier]() {
										goto l46
									}
									add(rulePegText, position110)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l46
								}
								{
									position112 := position
									{
										position113, tokenIndex113 := position, tokenIndex
										{
											position115 := position
											{
												position116 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l114
												}
												position++
											l117:
												{
													position118, tokenIndex118 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l118
													}
													position++
													goto l117
												l118:
													position, tokenIndex = position118, tokenIndex118
												}
												if !matchDot() {
													goto l114
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l114
												}
												position++
											l119:
												{
													position120, tokenIndex120 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l120
													}
													position++
													goto l119
												l120:
													position, tokenIndex = position120, tokenIndex120
												}
												if !matchDot() {
													goto l114
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l114
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
													goto l114
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l114
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
												if buffer[position] != rune('/') {
													goto l114
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l114
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
												add(ruleCidrValue, position116)
											}
											add(rulePegText, position115)
										}
										{
											add(ruleAction8, position)
										}
										goto l113
									l114:
										position, tokenIndex = position113, tokenIndex113
										{
											position129 := position
											{
												position130 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l128
												}
												position++
											l131:
												{
													position132, tokenIndex132 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l132
													}
													position++
													goto l131
												l132:
													position, tokenIndex = position132, tokenIndex132
												}
												if !matchDot() {
													goto l128
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l128
												}
												position++
											l133:
												{
													position134, tokenIndex134 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l134
													}
													position++
													goto l133
												l134:
													position, tokenIndex = position134, tokenIndex134
												}
												if !matchDot() {
													goto l128
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l128
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
													goto l128
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l128
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
												add(ruleIpValue, position130)
											}
											add(rulePegText, position129)
										}
										{
											add(ruleAction9, position)
										}
										goto l113
									l128:
										position, tokenIndex = position113, tokenIndex113
										{
											position141 := position
											{
												position142 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l140
												}
												position++
											l143:
												{
													position144, tokenIndex144 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l144
													}
													position++
													goto l143
												l144:
													position, tokenIndex = position144, tokenIndex144
												}
												if buffer[position] != rune('-') {
													goto l140
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l140
												}
												position++
											l145:
												{
													position146, tokenIndex146 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l146
													}
													position++
													goto l145
												l146:
													position, tokenIndex = position146, tokenIndex146
												}
												add(ruleIntRangeValue, position142)
											}
											add(rulePegText, position141)
										}
										{
											add(ruleAction10, position)
										}
										goto l113
									l140:
										position, tokenIndex = position113, tokenIndex113
										{
											position149 := position
											{
												position150 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l148
												}
												position++
											l151:
												{
													position152, tokenIndex152 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l152
													}
													position++
													goto l151
												l152:
													position, tokenIndex = position152, tokenIndex152
												}
												add(ruleIntValue, position150)
											}
											add(rulePegText, position149)
										}
										{
											add(ruleAction11, position)
										}
										goto l113
									l148:
										position, tokenIndex = position113, tokenIndex113
										{
											switch buffer[position] {
											case '$':
												{
													position155 := position
													if buffer[position] != rune('$') {
														goto l46
													}
													position++
													{
														position156 := position
														if !_rules[ruleIdentifier]() {
															goto l46
														}
														add(rulePegText, position156)
													}
													add(ruleRefValue, position155)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position158 := position
													if buffer[position] != rune('@') {
														goto l46
													}
													position++
													{
														position159 := position
														if !_rules[ruleIdentifier]() {
															goto l46
														}
														add(rulePegText, position159)
													}
													add(ruleAliasValue, position158)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position161 := position
													if buffer[position] != rune('{') {
														goto l46
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l46
													}
													{
														position162 := position
														if !_rules[ruleIdentifier]() {
															goto l46
														}
														add(rulePegText, position162)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l46
													}
													if buffer[position] != rune('}') {
														goto l46
													}
													position++
													add(ruleHoleValue, position161)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position164 := position
													{
														position165 := position
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

													l166:
														{
															position167, tokenIndex167 := position, tokenIndex
															{
																switch buffer[position] {
																case '_':
																	if buffer[position] != rune('_') {
																		goto l167
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
																		goto l167
																	}
																	position++
																	break
																case '-':
																	if buffer[position] != rune('-') {
																		goto l167
																	}
																	position++
																	break
																case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l167
																	}
																	position++
																	break
																case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																	if c := buffer[position]; c < rune('A') || c > rune('Z') {
																		goto l167
																	}
																	position++
																	break
																default:
																	if c := buffer[position]; c < rune('a') || c > rune('z') {
																		goto l167
																	}
																	position++
																	break
																}
															}

															goto l166
														l167:
															position, tokenIndex = position167, tokenIndex167
														}
														add(ruleStringValue, position165)
													}
													add(rulePegText, position164)
												}
												{
													add(ruleAction12, position)
												}
												break
											}
										}

									}
								l113:
									add(ruleValue, position112)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l46
								}
								add(ruleParam, position109)
							}
							goto l45
						l46:
							position, tokenIndex = position46, tokenIndex46
						}
						{
							add(ruleAction3, position)
						}
						add(ruleParams, position44)
					}
					goto l43
				l42:
					position, tokenIndex = position42, tokenIndex42
				}
			l43:
				add(ruleExpr, position27)
			}
			return true
		l26:
			position, tokenIndex = position26, tokenIndex26
			return false
		},
		/* 6 Params <- <(Param+ Action3)> */
		nil,
		/* 7 Param <- <(<Identifier> Action4 Equal Value WhiteSpacing)> */
		nil,
		/* 8 Identifier <- <((&('.') '.') | (&('_') '_') | (&('-') '-') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position174, tokenIndex174 := position, tokenIndex
			{
				position175 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l174
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l174
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l174
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l174
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l174
						}
						position++
						break
					}
				}

			l176:
				{
					position177, tokenIndex177 := position, tokenIndex
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

					goto l176
				l177:
					position, tokenIndex = position177, tokenIndex177
				}
				add(ruleIdentifier, position175)
			}
			return true
		l174:
			position, tokenIndex = position174, tokenIndex174
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
				position190 := position
			l191:
				{
					position192, tokenIndex192 := position, tokenIndex
					{
						position193 := position
						{
							position194, tokenIndex194 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l195
							}
							goto l194
						l195:
							position, tokenIndex = position194, tokenIndex194
							if !_rules[ruleEndOfLine]() {
								goto l192
							}
						}
					l194:
						add(ruleSpace, position193)
					}
					goto l191
				l192:
					position, tokenIndex = position192, tokenIndex192
				}
				add(ruleSpacing, position190)
			}
			return true
		},
		/* 19 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position197 := position
			l198:
				{
					position199, tokenIndex199 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l199
					}
					goto l198
				l199:
					position, tokenIndex = position199, tokenIndex199
				}
				add(ruleWhiteSpacing, position197)
			}
			return true
		},
		/* 20 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position200, tokenIndex200 := position, tokenIndex
			{
				position201 := position
				if !_rules[ruleWhitespace]() {
					goto l200
				}
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
				add(ruleMustWhiteSpacing, position201)
			}
			return true
		l200:
			position, tokenIndex = position200, tokenIndex200
			return false
		},
		/* 21 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position204, tokenIndex204 := position, tokenIndex
			{
				position205 := position
				if !_rules[ruleSpacing]() {
					goto l204
				}
				if buffer[position] != rune('=') {
					goto l204
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l204
				}
				add(ruleEqual, position205)
			}
			return true
		l204:
			position, tokenIndex = position204, tokenIndex204
			return false
		},
		/* 22 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 23 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position207, tokenIndex207 := position, tokenIndex
			{
				position208 := position
				{
					position209, tokenIndex209 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l210
					}
					position++
					goto l209
				l210:
					position, tokenIndex = position209, tokenIndex209
					if buffer[position] != rune('\t') {
						goto l207
					}
					position++
				}
			l209:
				add(ruleWhitespace, position208)
			}
			return true
		l207:
			position, tokenIndex = position207, tokenIndex207
			return false
		},
		/* 24 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position211, tokenIndex211 := position, tokenIndex
			{
				position212 := position
				{
					position213, tokenIndex213 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l214
					}
					position++
					if buffer[position] != rune('\n') {
						goto l214
					}
					position++
					goto l213
				l214:
					position, tokenIndex = position213, tokenIndex213
					if buffer[position] != rune('\n') {
						goto l215
					}
					position++
					goto l213
				l215:
					position, tokenIndex = position213, tokenIndex213
					if buffer[position] != rune('\r') {
						goto l211
					}
					position++
				}
			l213:
				add(ruleEndOfLine, position212)
			}
			return true
		l211:
			position, tokenIndex = position211, tokenIndex211
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
