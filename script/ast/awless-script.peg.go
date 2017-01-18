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
		/* 2 Action <- <(('c' 'r' 'e' 'a' 't' 'e') / ('d' 'e' 'l' 'e' 't' 'e'))> */
		nil,
		/* 3 Entity <- <((&('i') ('i' 'n' 's' 't' 'a' 'n' 'c' 'e')) | (&('s') ('s' 'u' 'b' 'n' 'e' 't')) | (&('v') ('v' 'p' 'c')))> */
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
					position33 := position
					{
						position34 := position
						{
							switch buffer[position] {
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

						add(ruleEntity, position34)
					}
					add(rulePegText, position33)
				}
				{
					add(ruleAction2, position)
				}
				{
					position37, tokenIndex37 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l37
					}
					{
						position39 := position
						{
							position42 := position
							{
								position43 := position
								if !_rules[ruleIdentifier]() {
									goto l37
								}
								add(rulePegText, position43)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l37
							}
							{
								position45 := position
								{
									position46, tokenIndex46 := position, tokenIndex
									{
										position48 := position
										{
											position49 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l47
											}
											position++
										l50:
											{
												position51, tokenIndex51 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l51
												}
												position++
												goto l50
											l51:
												position, tokenIndex = position51, tokenIndex51
											}
											if !matchDot() {
												goto l47
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l47
											}
											position++
										l52:
											{
												position53, tokenIndex53 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l53
												}
												position++
												goto l52
											l53:
												position, tokenIndex = position53, tokenIndex53
											}
											if !matchDot() {
												goto l47
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l47
											}
											position++
										l54:
											{
												position55, tokenIndex55 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l55
												}
												position++
												goto l54
											l55:
												position, tokenIndex = position55, tokenIndex55
											}
											if !matchDot() {
												goto l47
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l47
											}
											position++
										l56:
											{
												position57, tokenIndex57 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l57
												}
												position++
												goto l56
											l57:
												position, tokenIndex = position57, tokenIndex57
											}
											if buffer[position] != rune('/') {
												goto l47
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l47
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
											add(ruleCidrValue, position49)
										}
										add(rulePegText, position48)
									}
									{
										add(ruleAction8, position)
									}
									goto l46
								l47:
									position, tokenIndex = position46, tokenIndex46
									{
										position62 := position
										{
											position63 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l61
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
											if !matchDot() {
												goto l61
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l61
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
											if !matchDot() {
												goto l61
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l61
											}
											position++
										l68:
											{
												position69, tokenIndex69 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l69
												}
												position++
												goto l68
											l69:
												position, tokenIndex = position69, tokenIndex69
											}
											if !matchDot() {
												goto l61
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l61
											}
											position++
										l70:
											{
												position71, tokenIndex71 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l71
												}
												position++
												goto l70
											l71:
												position, tokenIndex = position71, tokenIndex71
											}
											add(ruleIpValue, position63)
										}
										add(rulePegText, position62)
									}
									{
										add(ruleAction9, position)
									}
									goto l46
								l61:
									position, tokenIndex = position46, tokenIndex46
									{
										position74 := position
										{
											position75 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l73
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
											add(ruleIntValue, position75)
										}
										add(rulePegText, position74)
									}
									{
										add(ruleAction10, position)
									}
									goto l46
								l73:
									position, tokenIndex = position46, tokenIndex46
									{
										switch buffer[position] {
										case '$':
											{
												position80 := position
												if buffer[position] != rune('$') {
													goto l37
												}
												position++
												{
													position81 := position
													if !_rules[ruleIdentifier]() {
														goto l37
													}
													add(rulePegText, position81)
												}
												add(ruleRefValue, position80)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position83 := position
												if buffer[position] != rune('@') {
													goto l37
												}
												position++
												{
													position84 := position
													if !_rules[ruleIdentifier]() {
														goto l37
													}
													add(rulePegText, position84)
												}
												add(ruleAliasValue, position83)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position86 := position
												if buffer[position] != rune('{') {
													goto l37
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l37
												}
												{
													position87 := position
													if !_rules[ruleIdentifier]() {
														goto l37
													}
													add(rulePegText, position87)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l37
												}
												if buffer[position] != rune('}') {
													goto l37
												}
												position++
												add(ruleHoleValue, position86)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position89 := position
												{
													position90 := position
													{
														switch buffer[position] {
														case '_':
															if buffer[position] != rune('_') {
																goto l37
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l37
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l37
															}
															position++
															break
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l37
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l37
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l37
															}
															position++
															break
														}
													}

												l91:
													{
														position92, tokenIndex92 := position, tokenIndex
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l92
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l92
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l92
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l92
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l92
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l92
																}
																position++
																break
															}
														}

														goto l91
													l92:
														position, tokenIndex = position92, tokenIndex92
													}
													add(ruleStringValue, position90)
												}
												add(rulePegText, position89)
											}
											{
												add(ruleAction11, position)
											}
											break
										}
									}

								}
							l46:
								add(ruleValue, position45)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l37
							}
							add(ruleParam, position42)
						}
					l40:
						{
							position41, tokenIndex41 := position, tokenIndex
							{
								position96 := position
								{
									position97 := position
									if !_rules[ruleIdentifier]() {
										goto l41
									}
									add(rulePegText, position97)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l41
								}
								{
									position99 := position
									{
										position100, tokenIndex100 := position, tokenIndex
										{
											position102 := position
											{
												position103 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l101
												}
												position++
											l104:
												{
													position105, tokenIndex105 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l105
													}
													position++
													goto l104
												l105:
													position, tokenIndex = position105, tokenIndex105
												}
												if !matchDot() {
													goto l101
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l101
												}
												position++
											l106:
												{
													position107, tokenIndex107 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l107
													}
													position++
													goto l106
												l107:
													position, tokenIndex = position107, tokenIndex107
												}
												if !matchDot() {
													goto l101
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l101
												}
												position++
											l108:
												{
													position109, tokenIndex109 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l109
													}
													position++
													goto l108
												l109:
													position, tokenIndex = position109, tokenIndex109
												}
												if !matchDot() {
													goto l101
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l101
												}
												position++
											l110:
												{
													position111, tokenIndex111 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l111
													}
													position++
													goto l110
												l111:
													position, tokenIndex = position111, tokenIndex111
												}
												if buffer[position] != rune('/') {
													goto l101
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l101
												}
												position++
											l112:
												{
													position113, tokenIndex113 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l113
													}
													position++
													goto l112
												l113:
													position, tokenIndex = position113, tokenIndex113
												}
												add(ruleCidrValue, position103)
											}
											add(rulePegText, position102)
										}
										{
											add(ruleAction8, position)
										}
										goto l100
									l101:
										position, tokenIndex = position100, tokenIndex100
										{
											position116 := position
											{
												position117 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l115
												}
												position++
											l118:
												{
													position119, tokenIndex119 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l119
													}
													position++
													goto l118
												l119:
													position, tokenIndex = position119, tokenIndex119
												}
												if !matchDot() {
													goto l115
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l115
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
													goto l115
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l115
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
													goto l115
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l115
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
												add(ruleIpValue, position117)
											}
											add(rulePegText, position116)
										}
										{
											add(ruleAction9, position)
										}
										goto l100
									l115:
										position, tokenIndex = position100, tokenIndex100
										{
											position128 := position
											{
												position129 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l127
												}
												position++
											l130:
												{
													position131, tokenIndex131 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l131
													}
													position++
													goto l130
												l131:
													position, tokenIndex = position131, tokenIndex131
												}
												add(ruleIntValue, position129)
											}
											add(rulePegText, position128)
										}
										{
											add(ruleAction10, position)
										}
										goto l100
									l127:
										position, tokenIndex = position100, tokenIndex100
										{
											switch buffer[position] {
											case '$':
												{
													position134 := position
													if buffer[position] != rune('$') {
														goto l41
													}
													position++
													{
														position135 := position
														if !_rules[ruleIdentifier]() {
															goto l41
														}
														add(rulePegText, position135)
													}
													add(ruleRefValue, position134)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position137 := position
													if buffer[position] != rune('@') {
														goto l41
													}
													position++
													{
														position138 := position
														if !_rules[ruleIdentifier]() {
															goto l41
														}
														add(rulePegText, position138)
													}
													add(ruleAliasValue, position137)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position140 := position
													if buffer[position] != rune('{') {
														goto l41
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l41
													}
													{
														position141 := position
														if !_rules[ruleIdentifier]() {
															goto l41
														}
														add(rulePegText, position141)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l41
													}
													if buffer[position] != rune('}') {
														goto l41
													}
													position++
													add(ruleHoleValue, position140)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position143 := position
													{
														position144 := position
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l41
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l41
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l41
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l41
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l41
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l41
																}
																position++
																break
															}
														}

													l145:
														{
															position146, tokenIndex146 := position, tokenIndex
															{
																switch buffer[position] {
																case '_':
																	if buffer[position] != rune('_') {
																		goto l146
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
																		goto l146
																	}
																	position++
																	break
																case '-':
																	if buffer[position] != rune('-') {
																		goto l146
																	}
																	position++
																	break
																case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l146
																	}
																	position++
																	break
																case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																	if c := buffer[position]; c < rune('A') || c > rune('Z') {
																		goto l146
																	}
																	position++
																	break
																default:
																	if c := buffer[position]; c < rune('a') || c > rune('z') {
																		goto l146
																	}
																	position++
																	break
																}
															}

															goto l145
														l146:
															position, tokenIndex = position146, tokenIndex146
														}
														add(ruleStringValue, position144)
													}
													add(rulePegText, position143)
												}
												{
													add(ruleAction11, position)
												}
												break
											}
										}

									}
								l100:
									add(ruleValue, position99)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l41
								}
								add(ruleParam, position96)
							}
							goto l40
						l41:
							position, tokenIndex = position41, tokenIndex41
						}
						{
							add(ruleAction3, position)
						}
						add(ruleParams, position39)
					}
					goto l38
				l37:
					position, tokenIndex = position37, tokenIndex37
				}
			l38:
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
			position153, tokenIndex153 := position, tokenIndex
			{
				position154 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l153
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l153
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l153
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l153
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l153
						}
						position++
						break
					}
				}

			l155:
				{
					position156, tokenIndex156 := position, tokenIndex
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

					goto l155
				l156:
					position, tokenIndex = position156, tokenIndex156
				}
				add(ruleIdentifier, position154)
			}
			return true
		l153:
			position, tokenIndex = position153, tokenIndex153
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
				position168 := position
			l169:
				{
					position170, tokenIndex170 := position, tokenIndex
					{
						position171 := position
						{
							position172, tokenIndex172 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l173
							}
							goto l172
						l173:
							position, tokenIndex = position172, tokenIndex172
							if !_rules[ruleEndOfLine]() {
								goto l170
							}
						}
					l172:
						add(ruleSpace, position171)
					}
					goto l169
				l170:
					position, tokenIndex = position170, tokenIndex170
				}
				add(ruleSpacing, position168)
			}
			return true
		},
		/* 18 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position175 := position
			l176:
				{
					position177, tokenIndex177 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l177
					}
					goto l176
				l177:
					position, tokenIndex = position177, tokenIndex177
				}
				add(ruleWhiteSpacing, position175)
			}
			return true
		},
		/* 19 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position178, tokenIndex178 := position, tokenIndex
			{
				position179 := position
				if !_rules[ruleWhitespace]() {
					goto l178
				}
			l180:
				{
					position181, tokenIndex181 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l181
					}
					goto l180
				l181:
					position, tokenIndex = position181, tokenIndex181
				}
				add(ruleMustWhiteSpacing, position179)
			}
			return true
		l178:
			position, tokenIndex = position178, tokenIndex178
			return false
		},
		/* 20 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position182, tokenIndex182 := position, tokenIndex
			{
				position183 := position
				if !_rules[ruleSpacing]() {
					goto l182
				}
				if buffer[position] != rune('=') {
					goto l182
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l182
				}
				add(ruleEqual, position183)
			}
			return true
		l182:
			position, tokenIndex = position182, tokenIndex182
			return false
		},
		/* 21 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 22 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position185, tokenIndex185 := position, tokenIndex
			{
				position186 := position
				{
					position187, tokenIndex187 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l188
					}
					position++
					goto l187
				l188:
					position, tokenIndex = position187, tokenIndex187
					if buffer[position] != rune('\t') {
						goto l185
					}
					position++
				}
			l187:
				add(ruleWhitespace, position186)
			}
			return true
		l185:
			position, tokenIndex = position185, tokenIndex185
			return false
		},
		/* 23 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position189, tokenIndex189 := position, tokenIndex
			{
				position190 := position
				{
					position191, tokenIndex191 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l192
					}
					position++
					if buffer[position] != rune('\n') {
						goto l192
					}
					position++
					goto l191
				l192:
					position, tokenIndex = position191, tokenIndex191
					if buffer[position] != rune('\n') {
						goto l193
					}
					position++
					goto l191
				l193:
					position, tokenIndex = position191, tokenIndex191
					if buffer[position] != rune('\r') {
						goto l189
					}
					position++
				}
			l191:
				add(ruleEndOfLine, position190)
			}
			return true
		l189:
			position, tokenIndex = position189, tokenIndex189
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
