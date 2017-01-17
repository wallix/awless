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
	rules  [37]func() bool
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
			p.AddParamRefValue(text)
		case ruleAction7:
			p.AddParamCidrValue(text)
		case ruleAction8:
			p.AddParamIpValue(text)
		case ruleAction9:
			p.AddParamIntValue(text)
		case ruleAction10:
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
										add(ruleAction7, position)
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
										add(ruleAction8, position)
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
										add(ruleAction9, position)
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
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position83 := position
												if buffer[position] != rune('{') {
													goto l37
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l37
												}
												{
													position84 := position
													if !_rules[ruleIdentifier]() {
														goto l37
													}
													add(rulePegText, position84)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l37
												}
												if buffer[position] != rune('}') {
													goto l37
												}
												position++
												add(ruleHoleValue, position83)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position86 := position
												{
													position87 := position
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

												l88:
													{
														position89, tokenIndex89 := position, tokenIndex
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l89
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l89
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l89
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l89
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l89
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l89
																}
																position++
																break
															}
														}

														goto l88
													l89:
														position, tokenIndex = position89, tokenIndex89
													}
													add(ruleStringValue, position87)
												}
												add(rulePegText, position86)
											}
											{
												add(ruleAction10, position)
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
								position93 := position
								{
									position94 := position
									if !_rules[ruleIdentifier]() {
										goto l41
									}
									add(rulePegText, position94)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l41
								}
								{
									position96 := position
									{
										position97, tokenIndex97 := position, tokenIndex
										{
											position99 := position
											{
												position100 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l98
												}
												position++
											l101:
												{
													position102, tokenIndex102 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l102
													}
													position++
													goto l101
												l102:
													position, tokenIndex = position102, tokenIndex102
												}
												if !matchDot() {
													goto l98
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l98
												}
												position++
											l103:
												{
													position104, tokenIndex104 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l104
													}
													position++
													goto l103
												l104:
													position, tokenIndex = position104, tokenIndex104
												}
												if !matchDot() {
													goto l98
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l98
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
													goto l98
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l98
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
												if buffer[position] != rune('/') {
													goto l98
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l98
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
												add(ruleCidrValue, position100)
											}
											add(rulePegText, position99)
										}
										{
											add(ruleAction7, position)
										}
										goto l97
									l98:
										position, tokenIndex = position97, tokenIndex97
										{
											position113 := position
											{
												position114 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l112
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
												if !matchDot() {
													goto l112
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l112
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
													goto l112
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l112
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
													goto l112
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l112
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
												add(ruleIpValue, position114)
											}
											add(rulePegText, position113)
										}
										{
											add(ruleAction8, position)
										}
										goto l97
									l112:
										position, tokenIndex = position97, tokenIndex97
										{
											position125 := position
											{
												position126 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l124
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
												add(ruleIntValue, position126)
											}
											add(rulePegText, position125)
										}
										{
											add(ruleAction9, position)
										}
										goto l97
									l124:
										position, tokenIndex = position97, tokenIndex97
										{
											switch buffer[position] {
											case '$':
												{
													position131 := position
													if buffer[position] != rune('$') {
														goto l41
													}
													position++
													{
														position132 := position
														if !_rules[ruleIdentifier]() {
															goto l41
														}
														add(rulePegText, position132)
													}
													add(ruleRefValue, position131)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position134 := position
													if buffer[position] != rune('{') {
														goto l41
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l41
													}
													{
														position135 := position
														if !_rules[ruleIdentifier]() {
															goto l41
														}
														add(rulePegText, position135)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l41
													}
													if buffer[position] != rune('}') {
														goto l41
													}
													position++
													add(ruleHoleValue, position134)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position137 := position
													{
														position138 := position
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

													l139:
														{
															position140, tokenIndex140 := position, tokenIndex
															{
																switch buffer[position] {
																case '_':
																	if buffer[position] != rune('_') {
																		goto l140
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
																		goto l140
																	}
																	position++
																	break
																case '-':
																	if buffer[position] != rune('-') {
																		goto l140
																	}
																	position++
																	break
																case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l140
																	}
																	position++
																	break
																case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																	if c := buffer[position]; c < rune('A') || c > rune('Z') {
																		goto l140
																	}
																	position++
																	break
																default:
																	if c := buffer[position]; c < rune('a') || c > rune('z') {
																		goto l140
																	}
																	position++
																	break
																}
															}

															goto l139
														l140:
															position, tokenIndex = position140, tokenIndex140
														}
														add(ruleStringValue, position138)
													}
													add(rulePegText, position137)
												}
												{
													add(ruleAction10, position)
												}
												break
											}
										}

									}
								l97:
									add(ruleValue, position96)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l41
								}
								add(ruleParam, position93)
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
			position147, tokenIndex147 := position, tokenIndex
			{
				position148 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l147
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
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

			l149:
				{
					position150, tokenIndex150 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l150
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l150
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l150
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l150
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l150
							}
							position++
							break
						}
					}

					goto l149
				l150:
					position, tokenIndex = position150, tokenIndex150
				}
				add(ruleIdentifier, position148)
			}
			return true
		l147:
			position, tokenIndex = position147, tokenIndex147
			return false
		},
		/* 9 Value <- <((<CidrValue> Action7) / (<IpValue> Action8) / (<IntValue> Action9) / ((&('$') (RefValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action10))))> */
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
		/* 15 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 16 Spacing <- <Space*> */
		func() bool {
			{
				position161 := position
			l162:
				{
					position163, tokenIndex163 := position, tokenIndex
					{
						position164 := position
						{
							position165, tokenIndex165 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l166
							}
							goto l165
						l166:
							position, tokenIndex = position165, tokenIndex165
							if !_rules[ruleEndOfLine]() {
								goto l163
							}
						}
					l165:
						add(ruleSpace, position164)
					}
					goto l162
				l163:
					position, tokenIndex = position163, tokenIndex163
				}
				add(ruleSpacing, position161)
			}
			return true
		},
		/* 17 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position168 := position
			l169:
				{
					position170, tokenIndex170 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l170
					}
					goto l169
				l170:
					position, tokenIndex = position170, tokenIndex170
				}
				add(ruleWhiteSpacing, position168)
			}
			return true
		},
		/* 18 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position171, tokenIndex171 := position, tokenIndex
			{
				position172 := position
				if !_rules[ruleWhitespace]() {
					goto l171
				}
			l173:
				{
					position174, tokenIndex174 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l174
					}
					goto l173
				l174:
					position, tokenIndex = position174, tokenIndex174
				}
				add(ruleMustWhiteSpacing, position172)
			}
			return true
		l171:
			position, tokenIndex = position171, tokenIndex171
			return false
		},
		/* 19 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position175, tokenIndex175 := position, tokenIndex
			{
				position176 := position
				if !_rules[ruleSpacing]() {
					goto l175
				}
				if buffer[position] != rune('=') {
					goto l175
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l175
				}
				add(ruleEqual, position176)
			}
			return true
		l175:
			position, tokenIndex = position175, tokenIndex175
			return false
		},
		/* 20 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 21 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position178, tokenIndex178 := position, tokenIndex
			{
				position179 := position
				{
					position180, tokenIndex180 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l181
					}
					position++
					goto l180
				l181:
					position, tokenIndex = position180, tokenIndex180
					if buffer[position] != rune('\t') {
						goto l178
					}
					position++
				}
			l180:
				add(ruleWhitespace, position179)
			}
			return true
		l178:
			position, tokenIndex = position178, tokenIndex178
			return false
		},
		/* 22 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position182, tokenIndex182 := position, tokenIndex
			{
				position183 := position
				{
					position184, tokenIndex184 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l185
					}
					position++
					if buffer[position] != rune('\n') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex = position184, tokenIndex184
					if buffer[position] != rune('\n') {
						goto l186
					}
					position++
					goto l184
				l186:
					position, tokenIndex = position184, tokenIndex184
					if buffer[position] != rune('\r') {
						goto l182
					}
					position++
				}
			l184:
				add(ruleEndOfLine, position183)
			}
			return true
		l182:
			position, tokenIndex = position182, tokenIndex182
			return false
		},
		/* 23 EndOfFile <- <!.> */
		nil,
		nil,
		/* 26 Action0 <- <{ p.AddDeclarationIdentifier(text) }> */
		nil,
		/* 27 Action1 <- <{ p.AddAction(text) }> */
		nil,
		/* 28 Action2 <- <{ p.AddEntity(text) }> */
		nil,
		/* 29 Action3 <- <{ p.EndOfParams() }> */
		nil,
		/* 30 Action4 <- <{ p.AddParamKey(text) }> */
		nil,
		/* 31 Action5 <- <{  p.AddParamHoleValue(text) }> */
		nil,
		/* 32 Action6 <- <{  p.AddParamRefValue(text) }> */
		nil,
		/* 33 Action7 <- <{ p.AddParamCidrValue(text) }> */
		nil,
		/* 34 Action8 <- <{ p.AddParamIpValue(text) }> */
		nil,
		/* 35 Action9 <- <{ p.AddParamIntValue(text) }> */
		nil,
		/* 36 Action10 <- <{ p.AddParamValue(text) }> */
		nil,
	}
	p.rules = _rules
}
