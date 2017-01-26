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
		/* 2 Action <- <(('s' 't' 'a' 'r' 't') / ((&('s') ('s' 't' 'o' 'p')) | (&('d') ('d' 'e' 'l' 'e' 't' 'e')) | (&('c') ('c' 'r' 'e' 'a' 't' 'e'))))> */
		nil,
		/* 3 Entity <- <((&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('r') ('r' 'o' 'l' 'e')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('t') ('t' 'a' 'g' 's')) | (&('i') ('i' 'n' 's' 't' 'a' 'n' 'c' 'e')) | (&('s') ('s' 'u' 'b' 'n' 'e' 't')) | (&('v') ('v' 'p' 'c')))> */
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
							case 's':
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
							default:
								if buffer[position] != rune('v') {
									goto l26
								}
								position++
								if buffer[position] != rune('p') {
									goto l26
								}
								position++
								if buffer[position] != rune('c') {
									goto l26
								}
								position++
								break
							}
						}

						add(ruleEntity, position35)
					}
					add(rulePegText, position34)
				}
				{
					add(ruleAction2, position)
				}
				{
					position38, tokenIndex38 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l38
					}
					{
						position40 := position
						{
							position43 := position
							{
								position44 := position
								if !_rules[ruleIdentifier]() {
									goto l38
								}
								add(rulePegText, position44)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l38
							}
							{
								position46 := position
								{
									position47, tokenIndex47 := position, tokenIndex
									{
										position49 := position
										{
											position50 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l48
											}
											position++
										l51:
											{
												position52, tokenIndex52 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l52
												}
												position++
												goto l51
											l52:
												position, tokenIndex = position52, tokenIndex52
											}
											if !matchDot() {
												goto l48
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l48
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
												goto l48
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l48
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
												goto l48
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l48
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
											if buffer[position] != rune('/') {
												goto l48
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l48
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
											add(ruleCidrValue, position50)
										}
										add(rulePegText, position49)
									}
									{
										add(ruleAction8, position)
									}
									goto l47
								l48:
									position, tokenIndex = position47, tokenIndex47
									{
										position63 := position
										{
											position64 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l62
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
											if !matchDot() {
												goto l62
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l62
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
												goto l62
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l62
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
												goto l62
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l62
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
											add(ruleIpValue, position64)
										}
										add(rulePegText, position63)
									}
									{
										add(ruleAction9, position)
									}
									goto l47
								l62:
									position, tokenIndex = position47, tokenIndex47
									{
										position75 := position
										{
											position76 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l74
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
											add(ruleIntValue, position76)
										}
										add(rulePegText, position75)
									}
									{
										add(ruleAction10, position)
									}
									goto l47
								l74:
									position, tokenIndex = position47, tokenIndex47
									{
										switch buffer[position] {
										case '$':
											{
												position81 := position
												if buffer[position] != rune('$') {
													goto l38
												}
												position++
												{
													position82 := position
													if !_rules[ruleIdentifier]() {
														goto l38
													}
													add(rulePegText, position82)
												}
												add(ruleRefValue, position81)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position84 := position
												if buffer[position] != rune('@') {
													goto l38
												}
												position++
												{
													position85 := position
													if !_rules[ruleIdentifier]() {
														goto l38
													}
													add(rulePegText, position85)
												}
												add(ruleAliasValue, position84)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position87 := position
												if buffer[position] != rune('{') {
													goto l38
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l38
												}
												{
													position88 := position
													if !_rules[ruleIdentifier]() {
														goto l38
													}
													add(rulePegText, position88)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l38
												}
												if buffer[position] != rune('}') {
													goto l38
												}
												position++
												add(ruleHoleValue, position87)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position90 := position
												{
													position91 := position
													{
														switch buffer[position] {
														case '_':
															if buffer[position] != rune('_') {
																goto l38
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l38
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l38
															}
															position++
															break
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l38
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l38
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l38
															}
															position++
															break
														}
													}

												l92:
													{
														position93, tokenIndex93 := position, tokenIndex
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l93
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l93
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l93
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l93
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l93
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l93
																}
																position++
																break
															}
														}

														goto l92
													l93:
														position, tokenIndex = position93, tokenIndex93
													}
													add(ruleStringValue, position91)
												}
												add(rulePegText, position90)
											}
											{
												add(ruleAction11, position)
											}
											break
										}
									}

								}
							l47:
								add(ruleValue, position46)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l38
							}
							add(ruleParam, position43)
						}
					l41:
						{
							position42, tokenIndex42 := position, tokenIndex
							{
								position97 := position
								{
									position98 := position
									if !_rules[ruleIdentifier]() {
										goto l42
									}
									add(rulePegText, position98)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l42
								}
								{
									position100 := position
									{
										position101, tokenIndex101 := position, tokenIndex
										{
											position103 := position
											{
												position104 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l102
												}
												position++
											l105:
												{
													position106, tokenIndex106 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l106
													}
													position++
													goto l105
												l106:
													position, tokenIndex = position106, tokenIndex106
												}
												if !matchDot() {
													goto l102
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l102
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
													goto l102
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l102
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
													goto l102
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l102
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
												if buffer[position] != rune('/') {
													goto l102
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l102
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
												add(ruleCidrValue, position104)
											}
											add(rulePegText, position103)
										}
										{
											add(ruleAction8, position)
										}
										goto l101
									l102:
										position, tokenIndex = position101, tokenIndex101
										{
											position117 := position
											{
												position118 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l116
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
													goto l116
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l116
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
													goto l116
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l116
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
													goto l116
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l116
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
												add(ruleIpValue, position118)
											}
											add(rulePegText, position117)
										}
										{
											add(ruleAction9, position)
										}
										goto l101
									l116:
										position, tokenIndex = position101, tokenIndex101
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
												add(ruleIntValue, position130)
											}
											add(rulePegText, position129)
										}
										{
											add(ruleAction10, position)
										}
										goto l101
									l128:
										position, tokenIndex = position101, tokenIndex101
										{
											switch buffer[position] {
											case '$':
												{
													position135 := position
													if buffer[position] != rune('$') {
														goto l42
													}
													position++
													{
														position136 := position
														if !_rules[ruleIdentifier]() {
															goto l42
														}
														add(rulePegText, position136)
													}
													add(ruleRefValue, position135)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position138 := position
													if buffer[position] != rune('@') {
														goto l42
													}
													position++
													{
														position139 := position
														if !_rules[ruleIdentifier]() {
															goto l42
														}
														add(rulePegText, position139)
													}
													add(ruleAliasValue, position138)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position141 := position
													if buffer[position] != rune('{') {
														goto l42
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l42
													}
													{
														position142 := position
														if !_rules[ruleIdentifier]() {
															goto l42
														}
														add(rulePegText, position142)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l42
													}
													if buffer[position] != rune('}') {
														goto l42
													}
													position++
													add(ruleHoleValue, position141)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position144 := position
													{
														position145 := position
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

													l146:
														{
															position147, tokenIndex147 := position, tokenIndex
															{
																switch buffer[position] {
																case '_':
																	if buffer[position] != rune('_') {
																		goto l147
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
																		goto l147
																	}
																	position++
																	break
																case '-':
																	if buffer[position] != rune('-') {
																		goto l147
																	}
																	position++
																	break
																case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l147
																	}
																	position++
																	break
																case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																	if c := buffer[position]; c < rune('A') || c > rune('Z') {
																		goto l147
																	}
																	position++
																	break
																default:
																	if c := buffer[position]; c < rune('a') || c > rune('z') {
																		goto l147
																	}
																	position++
																	break
																}
															}

															goto l146
														l147:
															position, tokenIndex = position147, tokenIndex147
														}
														add(ruleStringValue, position145)
													}
													add(rulePegText, position144)
												}
												{
													add(ruleAction11, position)
												}
												break
											}
										}

									}
								l101:
									add(ruleValue, position100)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l42
								}
								add(ruleParam, position97)
							}
							goto l41
						l42:
							position, tokenIndex = position42, tokenIndex42
						}
						{
							add(ruleAction3, position)
						}
						add(ruleParams, position40)
					}
					goto l39
				l38:
					position, tokenIndex = position38, tokenIndex38
				}
			l39:
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
			position154, tokenIndex154 := position, tokenIndex
			{
				position155 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l154
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l154
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l154
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l154
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l154
						}
						position++
						break
					}
				}

			l156:
				{
					position157, tokenIndex157 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l157
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l157
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l157
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l157
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l157
							}
							position++
							break
						}
					}

					goto l156
				l157:
					position, tokenIndex = position157, tokenIndex157
				}
				add(ruleIdentifier, position155)
			}
			return true
		l154:
			position, tokenIndex = position154, tokenIndex154
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
				position169 := position
			l170:
				{
					position171, tokenIndex171 := position, tokenIndex
					{
						position172 := position
						{
							position173, tokenIndex173 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l174
							}
							goto l173
						l174:
							position, tokenIndex = position173, tokenIndex173
							if !_rules[ruleEndOfLine]() {
								goto l171
							}
						}
					l173:
						add(ruleSpace, position172)
					}
					goto l170
				l171:
					position, tokenIndex = position171, tokenIndex171
				}
				add(ruleSpacing, position169)
			}
			return true
		},
		/* 18 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position176 := position
			l177:
				{
					position178, tokenIndex178 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l178
					}
					goto l177
				l178:
					position, tokenIndex = position178, tokenIndex178
				}
				add(ruleWhiteSpacing, position176)
			}
			return true
		},
		/* 19 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position179, tokenIndex179 := position, tokenIndex
			{
				position180 := position
				if !_rules[ruleWhitespace]() {
					goto l179
				}
			l181:
				{
					position182, tokenIndex182 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l182
					}
					goto l181
				l182:
					position, tokenIndex = position182, tokenIndex182
				}
				add(ruleMustWhiteSpacing, position180)
			}
			return true
		l179:
			position, tokenIndex = position179, tokenIndex179
			return false
		},
		/* 20 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position183, tokenIndex183 := position, tokenIndex
			{
				position184 := position
				if !_rules[ruleSpacing]() {
					goto l183
				}
				if buffer[position] != rune('=') {
					goto l183
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l183
				}
				add(ruleEqual, position184)
			}
			return true
		l183:
			position, tokenIndex = position183, tokenIndex183
			return false
		},
		/* 21 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 22 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position186, tokenIndex186 := position, tokenIndex
			{
				position187 := position
				{
					position188, tokenIndex188 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex = position188, tokenIndex188
					if buffer[position] != rune('\t') {
						goto l186
					}
					position++
				}
			l188:
				add(ruleWhitespace, position187)
			}
			return true
		l186:
			position, tokenIndex = position186, tokenIndex186
			return false
		},
		/* 23 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position190, tokenIndex190 := position, tokenIndex
			{
				position191 := position
				{
					position192, tokenIndex192 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l193
					}
					position++
					if buffer[position] != rune('\n') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex = position192, tokenIndex192
					if buffer[position] != rune('\n') {
						goto l194
					}
					position++
					goto l192
				l194:
					position, tokenIndex = position192, tokenIndex192
					if buffer[position] != rune('\r') {
						goto l190
					}
					position++
				}
			l192:
				add(ruleEndOfLine, position191)
			}
			return true
		l190:
			position, tokenIndex = position190, tokenIndex190
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
