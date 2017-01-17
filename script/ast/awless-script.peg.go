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
	rules  [35]func() bool
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
			p.AddParamCidrValue(text)
		case ruleAction7:
			p.AddParamIpValue(text)
		case ruleAction8:
			p.AddParamIntValue(text)
		case ruleAction9:
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
										if buffer[position] != rune('{') {
											goto l47
										}
										position++
										if !_rules[ruleWhiteSpacing]() {
											goto l47
										}
										{
											position49 := position
											if !_rules[ruleIdentifier]() {
												goto l47
											}
											add(rulePegText, position49)
										}
										if !_rules[ruleWhiteSpacing]() {
											goto l47
										}
										if buffer[position] != rune('}') {
											goto l47
										}
										position++
										add(ruleHoleValue, position48)
									}
									{
										add(ruleAction5, position)
									}
									goto l46
								l47:
									position, tokenIndex = position46, tokenIndex46
									{
										position52 := position
										{
											position53 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l51
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
												goto l51
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l51
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
											if !matchDot() {
												goto l51
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l51
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
												goto l51
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l51
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
											if buffer[position] != rune('/') {
												goto l51
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l51
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
											add(ruleCidrValue, position53)
										}
										add(rulePegText, position52)
									}
									{
										add(ruleAction6, position)
									}
									goto l46
								l51:
									position, tokenIndex = position46, tokenIndex46
									{
										position66 := position
										{
											position67 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l65
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
												goto l65
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l65
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
											if !matchDot() {
												goto l65
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l65
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
												goto l65
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l65
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
											add(ruleIpValue, position67)
										}
										add(rulePegText, position66)
									}
									{
										add(ruleAction7, position)
									}
									goto l46
								l65:
									position, tokenIndex = position46, tokenIndex46
									{
										position78 := position
										{
											position79 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l77
											}
											position++
										l80:
											{
												position81, tokenIndex81 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l81
												}
												position++
												goto l80
											l81:
												position, tokenIndex = position81, tokenIndex81
											}
											add(ruleIntValue, position79)
										}
										add(rulePegText, position78)
									}
									{
										add(ruleAction8, position)
									}
									goto l46
								l77:
									position, tokenIndex = position46, tokenIndex46
									{
										position83 := position
										{
											position84 := position
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

										l85:
											{
												position86, tokenIndex86 := position, tokenIndex
												{
													switch buffer[position] {
													case '_':
														if buffer[position] != rune('_') {
															goto l86
														}
														position++
														break
													case '.':
														if buffer[position] != rune('.') {
															goto l86
														}
														position++
														break
													case '-':
														if buffer[position] != rune('-') {
															goto l86
														}
														position++
														break
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l86
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l86
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l86
														}
														position++
														break
													}
												}

												goto l85
											l86:
												position, tokenIndex = position86, tokenIndex86
											}
											add(ruleStringValue, position84)
										}
										add(rulePegText, position83)
									}
									{
										add(ruleAction9, position)
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
								position90 := position
								{
									position91 := position
									if !_rules[ruleIdentifier]() {
										goto l41
									}
									add(rulePegText, position91)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l41
								}
								{
									position93 := position
									{
										position94, tokenIndex94 := position, tokenIndex
										{
											position96 := position
											if buffer[position] != rune('{') {
												goto l95
											}
											position++
											if !_rules[ruleWhiteSpacing]() {
												goto l95
											}
											{
												position97 := position
												if !_rules[ruleIdentifier]() {
													goto l95
												}
												add(rulePegText, position97)
											}
											if !_rules[ruleWhiteSpacing]() {
												goto l95
											}
											if buffer[position] != rune('}') {
												goto l95
											}
											position++
											add(ruleHoleValue, position96)
										}
										{
											add(ruleAction5, position)
										}
										goto l94
									l95:
										position, tokenIndex = position94, tokenIndex94
										{
											position100 := position
											{
												position101 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l99
												}
												position++
											l102:
												{
													position103, tokenIndex103 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l103
													}
													position++
													goto l102
												l103:
													position, tokenIndex = position103, tokenIndex103
												}
												if !matchDot() {
													goto l99
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l99
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
													goto l99
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l99
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
													goto l99
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l99
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
												if buffer[position] != rune('/') {
													goto l99
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l99
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
												add(ruleCidrValue, position101)
											}
											add(rulePegText, position100)
										}
										{
											add(ruleAction6, position)
										}
										goto l94
									l99:
										position, tokenIndex = position94, tokenIndex94
										{
											position114 := position
											{
												position115 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l113
												}
												position++
											l116:
												{
													position117, tokenIndex117 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l117
													}
													position++
													goto l116
												l117:
													position, tokenIndex = position117, tokenIndex117
												}
												if !matchDot() {
													goto l113
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l113
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
													goto l113
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l113
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
													goto l113
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l113
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
												add(ruleIpValue, position115)
											}
											add(rulePegText, position114)
										}
										{
											add(ruleAction7, position)
										}
										goto l94
									l113:
										position, tokenIndex = position94, tokenIndex94
										{
											position126 := position
											{
												position127 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l125
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
												add(ruleIntValue, position127)
											}
											add(rulePegText, position126)
										}
										{
											add(ruleAction8, position)
										}
										goto l94
									l125:
										position, tokenIndex = position94, tokenIndex94
										{
											position131 := position
											{
												position132 := position
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

											l133:
												{
													position134, tokenIndex134 := position, tokenIndex
													{
														switch buffer[position] {
														case '_':
															if buffer[position] != rune('_') {
																goto l134
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l134
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l134
															}
															position++
															break
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l134
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l134
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l134
															}
															position++
															break
														}
													}

													goto l133
												l134:
													position, tokenIndex = position134, tokenIndex134
												}
												add(ruleStringValue, position132)
											}
											add(rulePegText, position131)
										}
										{
											add(ruleAction9, position)
										}
									}
								l94:
									add(ruleValue, position93)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l41
								}
								add(ruleParam, position90)
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
			position141, tokenIndex141 := position, tokenIndex
			{
				position142 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l141
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l141
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l141
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l141
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l141
						}
						position++
						break
					}
				}

			l143:
				{
					position144, tokenIndex144 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l144
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l144
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l144
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l144
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l144
							}
							position++
							break
						}
					}

					goto l143
				l144:
					position, tokenIndex = position144, tokenIndex144
				}
				add(ruleIdentifier, position142)
			}
			return true
		l141:
			position, tokenIndex = position141, tokenIndex141
			return false
		},
		/* 9 Value <- <((HoleValue Action5) / (<CidrValue> Action6) / (<IpValue> Action7) / (<IntValue> Action8) / (<StringValue> Action9))> */
		nil,
		/* 10 StringValue <- <((&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		nil,
		/* 11 CidrValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+ '/' [0-9]+)> */
		nil,
		/* 12 IpValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+)> */
		nil,
		/* 13 IntValue <- <[0-9]+> */
		nil,
		/* 14 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 15 Spacing <- <Space*> */
		func() bool {
			{
				position154 := position
			l155:
				{
					position156, tokenIndex156 := position, tokenIndex
					{
						position157 := position
						{
							position158, tokenIndex158 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l159
							}
							goto l158
						l159:
							position, tokenIndex = position158, tokenIndex158
							if !_rules[ruleEndOfLine]() {
								goto l156
							}
						}
					l158:
						add(ruleSpace, position157)
					}
					goto l155
				l156:
					position, tokenIndex = position156, tokenIndex156
				}
				add(ruleSpacing, position154)
			}
			return true
		},
		/* 16 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position161 := position
			l162:
				{
					position163, tokenIndex163 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l163
					}
					goto l162
				l163:
					position, tokenIndex = position163, tokenIndex163
				}
				add(ruleWhiteSpacing, position161)
			}
			return true
		},
		/* 17 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position164, tokenIndex164 := position, tokenIndex
			{
				position165 := position
				if !_rules[ruleWhitespace]() {
					goto l164
				}
			l166:
				{
					position167, tokenIndex167 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l167
					}
					goto l166
				l167:
					position, tokenIndex = position167, tokenIndex167
				}
				add(ruleMustWhiteSpacing, position165)
			}
			return true
		l164:
			position, tokenIndex = position164, tokenIndex164
			return false
		},
		/* 18 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position168, tokenIndex168 := position, tokenIndex
			{
				position169 := position
				if !_rules[ruleSpacing]() {
					goto l168
				}
				if buffer[position] != rune('=') {
					goto l168
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l168
				}
				add(ruleEqual, position169)
			}
			return true
		l168:
			position, tokenIndex = position168, tokenIndex168
			return false
		},
		/* 19 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 20 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position171, tokenIndex171 := position, tokenIndex
			{
				position172 := position
				{
					position173, tokenIndex173 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l174
					}
					position++
					goto l173
				l174:
					position, tokenIndex = position173, tokenIndex173
					if buffer[position] != rune('\t') {
						goto l171
					}
					position++
				}
			l173:
				add(ruleWhitespace, position172)
			}
			return true
		l171:
			position, tokenIndex = position171, tokenIndex171
			return false
		},
		/* 21 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position175, tokenIndex175 := position, tokenIndex
			{
				position176 := position
				{
					position177, tokenIndex177 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l178
					}
					position++
					if buffer[position] != rune('\n') {
						goto l178
					}
					position++
					goto l177
				l178:
					position, tokenIndex = position177, tokenIndex177
					if buffer[position] != rune('\n') {
						goto l179
					}
					position++
					goto l177
				l179:
					position, tokenIndex = position177, tokenIndex177
					if buffer[position] != rune('\r') {
						goto l175
					}
					position++
				}
			l177:
				add(ruleEndOfLine, position176)
			}
			return true
		l175:
			position, tokenIndex = position175, tokenIndex175
			return false
		},
		/* 22 EndOfFile <- <!.> */
		nil,
		nil,
		/* 25 Action0 <- <{ p.AddDeclarationIdentifier(text) }> */
		nil,
		/* 26 Action1 <- <{ p.AddAction(text) }> */
		nil,
		/* 27 Action2 <- <{ p.AddEntity(text) }> */
		nil,
		/* 28 Action3 <- <{ p.EndOfParams() }> */
		nil,
		/* 29 Action4 <- <{ p.AddParamKey(text) }> */
		nil,
		/* 30 Action5 <- <{  p.AddParamHoleValue(text) }> */
		nil,
		/* 31 Action6 <- <{ p.AddParamCidrValue(text) }> */
		nil,
		/* 32 Action7 <- <{ p.AddParamIpValue(text) }> */
		nil,
		/* 33 Action8 <- <{ p.AddParamIntValue(text) }> */
		nil,
		/* 34 Action9 <- <{ p.AddParamValue(text) }> */
		nil,
	}
	p.rules = _rules
}
