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
	rules  [39]func() bool
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
			p.AddParamIntValue(text)
		case ruleAction11:
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
		/* 2 Action <- <(('s' 't' 'a' 'r' 't') / ((&('a') ('a' 't' 't' 'a' 'c' 'h')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('s') ('s' 't' 'o' 'p')) | (&('d') ('d' 'e' 'l' 'e' 't' 'e')) | (&('c') ('c' 'r' 'e' 'a' 't' 'e'))))> */
		nil,
		/* 3 Entity <- <(('v' 'p' 'c') / ((&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('r') ('r' 'o' 'l' 'e')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('t') ('t' 'a' 'g' 's')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e')) | (&('i') ('i' 'n' 's' 't' 'a' 'n' 'c' 'e')) | (&('s') ('s' 'u' 'b' 'n' 'e' 't'))))> */
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
							if buffer[position] != rune('s') {
								goto l31
							}
							position++
							if buffer[position] != rune('t') {
								goto l31
							}
							position++
							if buffer[position] != rune('a') {
								goto l31
							}
							position++
							if buffer[position] != rune('r') {
								goto l31
							}
							position++
							if buffer[position] != rune('t') {
								goto l31
							}
							position++
							goto l30
						l31:
							position, tokenIndex = position30, tokenIndex30
							{
								switch buffer[position] {
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
								case 'd':
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
								default:
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
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
					position34 := position
					{
						position35 := position
						{
							position36, tokenIndex36 := position, tokenIndex
							if buffer[position] != rune('v') {
								goto l37
							}
							position++
							if buffer[position] != rune('p') {
								goto l37
							}
							position++
							if buffer[position] != rune('c') {
								goto l37
							}
							position++
							goto l36
						l37:
							position, tokenIndex = position36, tokenIndex36
							{
								switch buffer[position] {
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
								case 'i':
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
								default:
									if buffer[position] != rune('s') {
										goto l26
									}
									position++
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('b') {
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
									break
								}
							}

						}
					l36:
						add(ruleEntity, position35)
					}
					add(rulePegText, position34)
				}
				{
					add(ruleAction2, position)
				}
				{
					position40, tokenIndex40 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l40
					}
					{
						position42 := position
						{
							position45 := position
							{
								position46 := position
								if !_rules[ruleIdentifier]() {
									goto l40
								}
								add(rulePegText, position46)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l40
							}
							{
								position48 := position
								{
									position49, tokenIndex49 := position, tokenIndex
									{
										position51 := position
										{
											position52 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l50
											}
											position++
										l53:
											{
												position54, tokenIndex54 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l54
												}
												position++
												goto l53
											l54:
												position, tokenIndex = position54, tokenIndex54
											}
											if !matchDot() {
												goto l50
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l50
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
												goto l50
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l50
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
												goto l50
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l50
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
											if buffer[position] != rune('/') {
												goto l50
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l50
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
											add(ruleCidrValue, position52)
										}
										add(rulePegText, position51)
									}
									{
										add(ruleAction8, position)
									}
									goto l49
								l50:
									position, tokenIndex = position49, tokenIndex49
									{
										position65 := position
										{
											position66 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l64
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
											if !matchDot() {
												goto l64
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l64
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
												goto l64
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l64
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
												goto l64
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l64
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
											add(ruleIpValue, position66)
										}
										add(rulePegText, position65)
									}
									{
										add(ruleAction9, position)
									}
									goto l49
								l64:
									position, tokenIndex = position49, tokenIndex49
									{
										position77 := position
										{
											position78 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l76
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
											add(ruleIntValue, position78)
										}
										add(rulePegText, position77)
									}
									{
										add(ruleAction10, position)
									}
									goto l49
								l76:
									position, tokenIndex = position49, tokenIndex49
									{
										switch buffer[position] {
										case '$':
											{
												position83 := position
												if buffer[position] != rune('$') {
													goto l40
												}
												position++
												{
													position84 := position
													if !_rules[ruleIdentifier]() {
														goto l40
													}
													add(rulePegText, position84)
												}
												add(ruleRefValue, position83)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position86 := position
												if buffer[position] != rune('@') {
													goto l40
												}
												position++
												{
													position87 := position
													if !_rules[ruleIdentifier]() {
														goto l40
													}
													add(rulePegText, position87)
												}
												add(ruleAliasValue, position86)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position89 := position
												if buffer[position] != rune('{') {
													goto l40
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l40
												}
												{
													position90 := position
													if !_rules[ruleIdentifier]() {
														goto l40
													}
													add(rulePegText, position90)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l40
												}
												if buffer[position] != rune('}') {
													goto l40
												}
												position++
												add(ruleHoleValue, position89)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position92 := position
												{
													position93 := position
													{
														switch buffer[position] {
														case '_':
															if buffer[position] != rune('_') {
																goto l40
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l40
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l40
															}
															position++
															break
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l40
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l40
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l40
															}
															position++
															break
														}
													}

												l94:
													{
														position95, tokenIndex95 := position, tokenIndex
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l95
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l95
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l95
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l95
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l95
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l95
																}
																position++
																break
															}
														}

														goto l94
													l95:
														position, tokenIndex = position95, tokenIndex95
													}
													add(ruleStringValue, position93)
												}
												add(rulePegText, position92)
											}
											{
												add(ruleAction11, position)
											}
											break
										}
									}

								}
							l49:
								add(ruleValue, position48)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l40
							}
							add(ruleParam, position45)
						}
					l43:
						{
							position44, tokenIndex44 := position, tokenIndex
							{
								position99 := position
								{
									position100 := position
									if !_rules[ruleIdentifier]() {
										goto l44
									}
									add(rulePegText, position100)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l44
								}
								{
									position102 := position
									{
										position103, tokenIndex103 := position, tokenIndex
										{
											position105 := position
											{
												position106 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l104
												}
												position++
											l107:
												{
													position108, tokenIndex108 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l108
													}
													position++
													goto l107
												l108:
													position, tokenIndex = position108, tokenIndex108
												}
												if !matchDot() {
													goto l104
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l104
												}
												position++
											l109:
												{
													position110, tokenIndex110 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l110
													}
													position++
													goto l109
												l110:
													position, tokenIndex = position110, tokenIndex110
												}
												if !matchDot() {
													goto l104
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l104
												}
												position++
											l111:
												{
													position112, tokenIndex112 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l112
													}
													position++
													goto l111
												l112:
													position, tokenIndex = position112, tokenIndex112
												}
												if !matchDot() {
													goto l104
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l104
												}
												position++
											l113:
												{
													position114, tokenIndex114 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l114
													}
													position++
													goto l113
												l114:
													position, tokenIndex = position114, tokenIndex114
												}
												if buffer[position] != rune('/') {
													goto l104
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l104
												}
												position++
											l115:
												{
													position116, tokenIndex116 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l116
													}
													position++
													goto l115
												l116:
													position, tokenIndex = position116, tokenIndex116
												}
												add(ruleCidrValue, position106)
											}
											add(rulePegText, position105)
										}
										{
											add(ruleAction8, position)
										}
										goto l103
									l104:
										position, tokenIndex = position103, tokenIndex103
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
												add(ruleIpValue, position120)
											}
											add(rulePegText, position119)
										}
										{
											add(ruleAction9, position)
										}
										goto l103
									l118:
										position, tokenIndex = position103, tokenIndex103
										{
											position131 := position
											{
												position132 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l130
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
												add(ruleIntValue, position132)
											}
											add(rulePegText, position131)
										}
										{
											add(ruleAction10, position)
										}
										goto l103
									l130:
										position, tokenIndex = position103, tokenIndex103
										{
											switch buffer[position] {
											case '$':
												{
													position137 := position
													if buffer[position] != rune('$') {
														goto l44
													}
													position++
													{
														position138 := position
														if !_rules[ruleIdentifier]() {
															goto l44
														}
														add(rulePegText, position138)
													}
													add(ruleRefValue, position137)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position140 := position
													if buffer[position] != rune('@') {
														goto l44
													}
													position++
													{
														position141 := position
														if !_rules[ruleIdentifier]() {
															goto l44
														}
														add(rulePegText, position141)
													}
													add(ruleAliasValue, position140)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position143 := position
													if buffer[position] != rune('{') {
														goto l44
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l44
													}
													{
														position144 := position
														if !_rules[ruleIdentifier]() {
															goto l44
														}
														add(rulePegText, position144)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l44
													}
													if buffer[position] != rune('}') {
														goto l44
													}
													position++
													add(ruleHoleValue, position143)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position146 := position
													{
														position147 := position
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l44
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l44
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l44
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l44
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l44
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l44
																}
																position++
																break
															}
														}

													l148:
														{
															position149, tokenIndex149 := position, tokenIndex
															{
																switch buffer[position] {
																case '_':
																	if buffer[position] != rune('_') {
																		goto l149
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
																		goto l149
																	}
																	position++
																	break
																case '-':
																	if buffer[position] != rune('-') {
																		goto l149
																	}
																	position++
																	break
																case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l149
																	}
																	position++
																	break
																case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																	if c := buffer[position]; c < rune('A') || c > rune('Z') {
																		goto l149
																	}
																	position++
																	break
																default:
																	if c := buffer[position]; c < rune('a') || c > rune('z') {
																		goto l149
																	}
																	position++
																	break
																}
															}

															goto l148
														l149:
															position, tokenIndex = position149, tokenIndex149
														}
														add(ruleStringValue, position147)
													}
													add(rulePegText, position146)
												}
												{
													add(ruleAction11, position)
												}
												break
											}
										}

									}
								l103:
									add(ruleValue, position102)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l44
								}
								add(ruleParam, position99)
							}
							goto l43
						l44:
							position, tokenIndex = position44, tokenIndex44
						}
						{
							add(ruleAction3, position)
						}
						add(ruleParams, position42)
					}
					goto l41
				l40:
					position, tokenIndex = position40, tokenIndex40
				}
			l41:
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
			position156, tokenIndex156 := position, tokenIndex
			{
				position157 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l156
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l156
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l156
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l156
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l156
						}
						position++
						break
					}
				}

			l158:
				{
					position159, tokenIndex159 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l159
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l159
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l159
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l159
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l159
							}
							position++
							break
						}
					}

					goto l158
				l159:
					position, tokenIndex = position159, tokenIndex159
				}
				add(ruleIdentifier, position157)
			}
			return true
		l156:
			position, tokenIndex = position156, tokenIndex156
			return false
		},
		/* 9 Value <- <((<CidrValue> Action8) / (<IpValue> Action9) / (<IntValue> Action10) / ((&('$') (RefValue Action7)) | (&('@') (AliasValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action11))))> */
		nil,
		/* 10 StringValue <- <((&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		nil,
		/* 11 CidrValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+ '/' [0-9]+)> */
		nil,
		/* 12 IpValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+)> */
		nil,
		/* 13 IntValue <- <[0-9]+> */
		nil,
		/* 14 RefValue <- <('$' <Identifier>)> */
		nil,
		/* 15 AliasValue <- <('@' <Identifier>)> */
		nil,
		/* 16 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 17 Spacing <- <Space*> */
		func() bool {
			{
				position171 := position
			l172:
				{
					position173, tokenIndex173 := position, tokenIndex
					{
						position174 := position
						{
							position175, tokenIndex175 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l176
							}
							goto l175
						l176:
							position, tokenIndex = position175, tokenIndex175
							if !_rules[ruleEndOfLine]() {
								goto l173
							}
						}
					l175:
						add(ruleSpace, position174)
					}
					goto l172
				l173:
					position, tokenIndex = position173, tokenIndex173
				}
				add(ruleSpacing, position171)
			}
			return true
		},
		/* 18 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position178 := position
			l179:
				{
					position180, tokenIndex180 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l180
					}
					goto l179
				l180:
					position, tokenIndex = position180, tokenIndex180
				}
				add(ruleWhiteSpacing, position178)
			}
			return true
		},
		/* 19 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position181, tokenIndex181 := position, tokenIndex
			{
				position182 := position
				if !_rules[ruleWhitespace]() {
					goto l181
				}
			l183:
				{
					position184, tokenIndex184 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l184
					}
					goto l183
				l184:
					position, tokenIndex = position184, tokenIndex184
				}
				add(ruleMustWhiteSpacing, position182)
			}
			return true
		l181:
			position, tokenIndex = position181, tokenIndex181
			return false
		},
		/* 20 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position185, tokenIndex185 := position, tokenIndex
			{
				position186 := position
				if !_rules[ruleSpacing]() {
					goto l185
				}
				if buffer[position] != rune('=') {
					goto l185
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l185
				}
				add(ruleEqual, position186)
			}
			return true
		l185:
			position, tokenIndex = position185, tokenIndex185
			return false
		},
		/* 21 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 22 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position188, tokenIndex188 := position, tokenIndex
			{
				position189 := position
				{
					position190, tokenIndex190 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex = position190, tokenIndex190
					if buffer[position] != rune('\t') {
						goto l188
					}
					position++
				}
			l190:
				add(ruleWhitespace, position189)
			}
			return true
		l188:
			position, tokenIndex = position188, tokenIndex188
			return false
		},
		/* 23 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position192, tokenIndex192 := position, tokenIndex
			{
				position193 := position
				{
					position194, tokenIndex194 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l195
					}
					position++
					if buffer[position] != rune('\n') {
						goto l195
					}
					position++
					goto l194
				l195:
					position, tokenIndex = position194, tokenIndex194
					if buffer[position] != rune('\n') {
						goto l196
					}
					position++
					goto l194
				l196:
					position, tokenIndex = position194, tokenIndex194
					if buffer[position] != rune('\r') {
						goto l192
					}
					position++
				}
			l194:
				add(ruleEndOfLine, position193)
			}
			return true
		l192:
			position, tokenIndex = position192, tokenIndex192
			return false
		},
		/* 24 EndOfFile <- <!.> */
		nil,
		nil,
		/* 27 Action0 <- <{ p.AddDeclarationIdentifier(text) }> */
		nil,
		/* 28 Action1 <- <{ p.AddAction(text) }> */
		nil,
		/* 29 Action2 <- <{ p.AddEntity(text) }> */
		nil,
		/* 30 Action3 <- <{ p.EndOfParams() }> */
		nil,
		/* 31 Action4 <- <{ p.AddParamKey(text) }> */
		nil,
		/* 32 Action5 <- <{  p.AddParamHoleValue(text) }> */
		nil,
		/* 33 Action6 <- <{  p.AddParamAliasValue(text) }> */
		nil,
		/* 34 Action7 <- <{  p.AddParamRefValue(text) }> */
		nil,
		/* 35 Action8 <- <{ p.AddParamCidrValue(text) }> */
		nil,
		/* 36 Action9 <- <{ p.AddParamIpValue(text) }> */
		nil,
		/* 37 Action10 <- <{ p.AddParamIntValue(text) }> */
		nil,
		/* 38 Action11 <- <{ p.AddParamValue(text) }> */
		nil,
	}
	p.rules = _rules
}
