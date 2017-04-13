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
	ruleCustomTypedValue
	ruleStringValue
	ruleDoubleQuotedValue
	ruleSingleQuotedValue
	ruleCSVValue
	ruleCidrValue
	ruleIpValue
	ruleIntValue
	ruleIntRangeValue
	ruleRefValue
	ruleAliasValue
	ruleHoleValue
	ruleComment
	ruleSingleQuote
	ruleDoubleQuote
	ruleWhiteSpacing
	ruleMustWhiteSpacing
	ruleEqual
	ruleBlankLine
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
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
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
	"CustomTypedValue",
	"StringValue",
	"DoubleQuotedValue",
	"SingleQuotedValue",
	"CSVValue",
	"CidrValue",
	"IpValue",
	"IntValue",
	"IntRangeValue",
	"RefValue",
	"AliasValue",
	"HoleValue",
	"Comment",
	"SingleQuote",
	"DoubleQuote",
	"WhiteSpacing",
	"MustWhiteSpacing",
	"Equal",
	"BlankLine",
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
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
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
	rules  [52]func() bool
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
			p.addDeclarationIdentifier(text)
		case ruleAction1:
			p.addAction(text)
		case ruleAction2:
			p.addEntity(text)
		case ruleAction3:
			p.LineDone()
		case ruleAction4:
			p.addParamKey(text)
		case ruleAction5:
			p.addParamHoleValue(text)
		case ruleAction6:
			p.addParamRefValue(text)
		case ruleAction7:
			p.addAliasParam(text)
		case ruleAction8:
			p.addParamValue(text)
		case ruleAction9:
			p.addParamValue(text)
		case ruleAction10:
			p.addParamIntValue(text)
		case ruleAction11:
			p.addParamValue(text)
		case ruleAction12:
			p.addParamCidrValue(text)
		case ruleAction13:
			p.addParamIpValue(text)
		case ruleAction14:
			p.addCsvValue(text)
		case ruleAction15:
			p.addParamValue(text)
		case ruleAction16:
			p.LineDone()
		case ruleAction17:
			p.LineDone()

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
		/* 0 Script <- <((BlankLine* Statement BlankLine*)+ WhiteSpacing EndOfFile)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
			l4:
				{
					position5, tokenIndex5 := position, tokenIndex
					if !_rules[ruleBlankLine]() {
						goto l5
					}
					goto l4
				l5:
					position, tokenIndex = position5, tokenIndex5
				}
				{
					position6 := position
					if !_rules[ruleWhiteSpacing]() {
						goto l0
					}
					{
						position7, tokenIndex7 := position, tokenIndex
						if !_rules[ruleExpr]() {
							goto l8
						}
						goto l7
					l8:
						position, tokenIndex = position7, tokenIndex7
						{
							position10 := position
							{
								position11 := position
								if !_rules[ruleIdentifier]() {
									goto l9
								}
								add(rulePegText, position11)
							}
							{
								add(ruleAction0, position)
							}
							if !_rules[ruleEqual]() {
								goto l9
							}
							if !_rules[ruleExpr]() {
								goto l9
							}
							add(ruleDeclaration, position10)
						}
						goto l7
					l9:
						position, tokenIndex = position7, tokenIndex7
						{
							position13 := position
							{
								position14, tokenIndex14 := position, tokenIndex
								if buffer[position] != rune('#') {
									goto l15
								}
								position++
							l16:
								{
									position17, tokenIndex17 := position, tokenIndex
									{
										position18, tokenIndex18 := position, tokenIndex
										if !_rules[ruleEndOfLine]() {
											goto l18
										}
										goto l17
									l18:
										position, tokenIndex = position18, tokenIndex18
									}
									if !matchDot() {
										goto l17
									}
									goto l16
								l17:
									position, tokenIndex = position17, tokenIndex17
								}
								goto l14
							l15:
								position, tokenIndex = position14, tokenIndex14
								if buffer[position] != rune('/') {
									goto l0
								}
								position++
								if buffer[position] != rune('/') {
									goto l0
								}
								position++
							l19:
								{
									position20, tokenIndex20 := position, tokenIndex
									{
										position21, tokenIndex21 := position, tokenIndex
										if !_rules[ruleEndOfLine]() {
											goto l21
										}
										goto l20
									l21:
										position, tokenIndex = position21, tokenIndex21
									}
									if !matchDot() {
										goto l20
									}
									goto l19
								l20:
									position, tokenIndex = position20, tokenIndex20
								}
								{
									add(ruleAction16, position)
								}
							}
						l14:
							add(ruleComment, position13)
						}
					}
				l7:
					if !_rules[ruleWhiteSpacing]() {
						goto l0
					}
				l23:
					{
						position24, tokenIndex24 := position, tokenIndex
						if !_rules[ruleEndOfLine]() {
							goto l24
						}
						goto l23
					l24:
						position, tokenIndex = position24, tokenIndex24
					}
					add(ruleStatement, position6)
				}
			l25:
				{
					position26, tokenIndex26 := position, tokenIndex
					if !_rules[ruleBlankLine]() {
						goto l26
					}
					goto l25
				l26:
					position, tokenIndex = position26, tokenIndex26
				}
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
				l27:
					{
						position28, tokenIndex28 := position, tokenIndex
						if !_rules[ruleBlankLine]() {
							goto l28
						}
						goto l27
					l28:
						position, tokenIndex = position28, tokenIndex28
					}
					{
						position29 := position
						if !_rules[ruleWhiteSpacing]() {
							goto l3
						}
						{
							position30, tokenIndex30 := position, tokenIndex
							if !_rules[ruleExpr]() {
								goto l31
							}
							goto l30
						l31:
							position, tokenIndex = position30, tokenIndex30
							{
								position33 := position
								{
									position34 := position
									if !_rules[ruleIdentifier]() {
										goto l32
									}
									add(rulePegText, position34)
								}
								{
									add(ruleAction0, position)
								}
								if !_rules[ruleEqual]() {
									goto l32
								}
								if !_rules[ruleExpr]() {
									goto l32
								}
								add(ruleDeclaration, position33)
							}
							goto l30
						l32:
							position, tokenIndex = position30, tokenIndex30
							{
								position36 := position
								{
									position37, tokenIndex37 := position, tokenIndex
									if buffer[position] != rune('#') {
										goto l38
									}
									position++
								l39:
									{
										position40, tokenIndex40 := position, tokenIndex
										{
											position41, tokenIndex41 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l41
											}
											goto l40
										l41:
											position, tokenIndex = position41, tokenIndex41
										}
										if !matchDot() {
											goto l40
										}
										goto l39
									l40:
										position, tokenIndex = position40, tokenIndex40
									}
									goto l37
								l38:
									position, tokenIndex = position37, tokenIndex37
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
								l42:
									{
										position43, tokenIndex43 := position, tokenIndex
										{
											position44, tokenIndex44 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l44
											}
											goto l43
										l44:
											position, tokenIndex = position44, tokenIndex44
										}
										if !matchDot() {
											goto l43
										}
										goto l42
									l43:
										position, tokenIndex = position43, tokenIndex43
									}
									{
										add(ruleAction16, position)
									}
								}
							l37:
								add(ruleComment, position36)
							}
						}
					l30:
						if !_rules[ruleWhiteSpacing]() {
							goto l3
						}
					l46:
						{
							position47, tokenIndex47 := position, tokenIndex
							if !_rules[ruleEndOfLine]() {
								goto l47
							}
							goto l46
						l47:
							position, tokenIndex = position47, tokenIndex47
						}
						add(ruleStatement, position29)
					}
				l48:
					{
						position49, tokenIndex49 := position, tokenIndex
						if !_rules[ruleBlankLine]() {
							goto l49
						}
						goto l48
					l49:
						position, tokenIndex = position49, tokenIndex49
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				if !_rules[ruleWhiteSpacing]() {
					goto l0
				}
				{
					position50 := position
					{
						position51, tokenIndex51 := position, tokenIndex
						if !matchDot() {
							goto l51
						}
						goto l0
					l51:
						position, tokenIndex = position51, tokenIndex51
					}
					add(ruleEndOfFile, position50)
				}
				add(ruleScript, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Statement <- <(WhiteSpacing (Expr / Declaration / Comment) WhiteSpacing EndOfLine*)> */
		nil,
		/* 2 Action <- <[a-z]+> */
		nil,
		/* 3 Entity <- <([a-z] / [0-9])+> */
		nil,
		/* 4 Declaration <- <(<Identifier> Action0 Equal Expr)> */
		nil,
		/* 5 Expr <- <(<Action> Action1 MustWhiteSpacing <Entity> Action2 (MustWhiteSpacing Params)? Action3)> */
		func() bool {
			position56, tokenIndex56 := position, tokenIndex
			{
				position57 := position
				{
					position58 := position
					{
						position59 := position
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l56
						}
						position++
					l60:
						{
							position61, tokenIndex61 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l61
							}
							position++
							goto l60
						l61:
							position, tokenIndex = position61, tokenIndex61
						}
						add(ruleAction, position59)
					}
					add(rulePegText, position58)
				}
				{
					add(ruleAction1, position)
				}
				if !_rules[ruleMustWhiteSpacing]() {
					goto l56
				}
				{
					position63 := position
					{
						position64 := position
						{
							position67, tokenIndex67 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l68
							}
							position++
							goto l67
						l68:
							position, tokenIndex = position67, tokenIndex67
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l56
							}
							position++
						}
					l67:
					l65:
						{
							position66, tokenIndex66 := position, tokenIndex
							{
								position69, tokenIndex69 := position, tokenIndex
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l70
								}
								position++
								goto l69
							l70:
								position, tokenIndex = position69, tokenIndex69
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l66
								}
								position++
							}
						l69:
							goto l65
						l66:
							position, tokenIndex = position66, tokenIndex66
						}
						add(ruleEntity, position64)
					}
					add(rulePegText, position63)
				}
				{
					add(ruleAction2, position)
				}
				{
					position72, tokenIndex72 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l72
					}
					{
						position74 := position
						{
							position77 := position
							{
								position78 := position
								if !_rules[ruleIdentifier]() {
									goto l72
								}
								add(rulePegText, position78)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l72
							}
							{
								position80 := position
								{
									position81, tokenIndex81 := position, tokenIndex
									{
										position83 := position
										{
											position84, tokenIndex84 := position, tokenIndex
											if buffer[position] != rune('@') {
												goto l85
											}
											position++
											{
												position86 := position
												if !_rules[ruleStringValue]() {
													goto l85
												}
												add(rulePegText, position86)
											}
											goto l84
										l85:
											position, tokenIndex = position84, tokenIndex84
											if buffer[position] != rune('@') {
												goto l87
											}
											position++
											if !_rules[ruleDoubleQuote]() {
												goto l87
											}
											{
												position88 := position
												if !_rules[ruleDoubleQuotedValue]() {
													goto l87
												}
												add(rulePegText, position88)
											}
											if !_rules[ruleDoubleQuote]() {
												goto l87
											}
											goto l84
										l87:
											position, tokenIndex = position84, tokenIndex84
											if buffer[position] != rune('@') {
												goto l82
											}
											position++
											if !_rules[ruleSingleQuote]() {
												goto l82
											}
											{
												position89 := position
												if !_rules[ruleSingleQuotedValue]() {
													goto l82
												}
												add(rulePegText, position89)
											}
											if !_rules[ruleSingleQuote]() {
												goto l82
											}
										}
									l84:
										add(ruleAliasValue, position83)
									}
									{
										add(ruleAction7, position)
									}
									goto l81
								l82:
									position, tokenIndex = position81, tokenIndex81
									if !_rules[ruleDoubleQuote]() {
										goto l91
									}
									if !_rules[ruleCustomTypedValue]() {
										goto l91
									}
									if !_rules[ruleDoubleQuote]() {
										goto l91
									}
									goto l81
								l91:
									position, tokenIndex = position81, tokenIndex81
									if !_rules[ruleSingleQuote]() {
										goto l92
									}
									if !_rules[ruleCustomTypedValue]() {
										goto l92
									}
									if !_rules[ruleSingleQuote]() {
										goto l92
									}
									goto l81
								l92:
									position, tokenIndex = position81, tokenIndex81
									if !_rules[ruleCustomTypedValue]() {
										goto l93
									}
									goto l81
								l93:
									position, tokenIndex = position81, tokenIndex81
									{
										position95 := position
										{
											position96 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l94
											}
											position++
										l97:
											{
												position98, tokenIndex98 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l98
												}
												position++
												goto l97
											l98:
												position, tokenIndex = position98, tokenIndex98
											}
											add(ruleIntValue, position96)
										}
										add(rulePegText, position95)
									}
									{
										add(ruleAction10, position)
									}
									goto l81
								l94:
									position, tokenIndex = position81, tokenIndex81
									{
										switch buffer[position] {
										case '\'':
											if !_rules[ruleSingleQuote]() {
												goto l72
											}
											{
												position101 := position
												if !_rules[ruleSingleQuotedValue]() {
													goto l72
												}
												add(rulePegText, position101)
											}
											{
												add(ruleAction9, position)
											}
											if !_rules[ruleSingleQuote]() {
												goto l72
											}
											break
										case '"':
											if !_rules[ruleDoubleQuote]() {
												goto l72
											}
											{
												position103 := position
												if !_rules[ruleDoubleQuotedValue]() {
													goto l72
												}
												add(rulePegText, position103)
											}
											{
												add(ruleAction8, position)
											}
											if !_rules[ruleDoubleQuote]() {
												goto l72
											}
											break
										case '$':
											{
												position105 := position
												if buffer[position] != rune('$') {
													goto l72
												}
												position++
												{
													position106 := position
													if !_rules[ruleIdentifier]() {
														goto l72
													}
													add(rulePegText, position106)
												}
												add(ruleRefValue, position105)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position108 := position
												if buffer[position] != rune('{') {
													goto l72
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l72
												}
												{
													position109 := position
													if !_rules[ruleIdentifier]() {
														goto l72
													}
													add(rulePegText, position109)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l72
												}
												if buffer[position] != rune('}') {
													goto l72
												}
												position++
												add(ruleHoleValue, position108)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position111 := position
												if !_rules[ruleStringValue]() {
													goto l72
												}
												add(rulePegText, position111)
											}
											{
												add(ruleAction11, position)
											}
											break
										}
									}

								}
							l81:
								add(ruleValue, position80)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l72
							}
							add(ruleParam, position77)
						}
					l75:
						{
							position76, tokenIndex76 := position, tokenIndex
							{
								position113 := position
								{
									position114 := position
									if !_rules[ruleIdentifier]() {
										goto l76
									}
									add(rulePegText, position114)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l76
								}
								{
									position116 := position
									{
										position117, tokenIndex117 := position, tokenIndex
										{
											position119 := position
											{
												position120, tokenIndex120 := position, tokenIndex
												if buffer[position] != rune('@') {
													goto l121
												}
												position++
												{
													position122 := position
													if !_rules[ruleStringValue]() {
														goto l121
													}
													add(rulePegText, position122)
												}
												goto l120
											l121:
												position, tokenIndex = position120, tokenIndex120
												if buffer[position] != rune('@') {
													goto l123
												}
												position++
												if !_rules[ruleDoubleQuote]() {
													goto l123
												}
												{
													position124 := position
													if !_rules[ruleDoubleQuotedValue]() {
														goto l123
													}
													add(rulePegText, position124)
												}
												if !_rules[ruleDoubleQuote]() {
													goto l123
												}
												goto l120
											l123:
												position, tokenIndex = position120, tokenIndex120
												if buffer[position] != rune('@') {
													goto l118
												}
												position++
												if !_rules[ruleSingleQuote]() {
													goto l118
												}
												{
													position125 := position
													if !_rules[ruleSingleQuotedValue]() {
														goto l118
													}
													add(rulePegText, position125)
												}
												if !_rules[ruleSingleQuote]() {
													goto l118
												}
											}
										l120:
											add(ruleAliasValue, position119)
										}
										{
											add(ruleAction7, position)
										}
										goto l117
									l118:
										position, tokenIndex = position117, tokenIndex117
										if !_rules[ruleDoubleQuote]() {
											goto l127
										}
										if !_rules[ruleCustomTypedValue]() {
											goto l127
										}
										if !_rules[ruleDoubleQuote]() {
											goto l127
										}
										goto l117
									l127:
										position, tokenIndex = position117, tokenIndex117
										if !_rules[ruleSingleQuote]() {
											goto l128
										}
										if !_rules[ruleCustomTypedValue]() {
											goto l128
										}
										if !_rules[ruleSingleQuote]() {
											goto l128
										}
										goto l117
									l128:
										position, tokenIndex = position117, tokenIndex117
										if !_rules[ruleCustomTypedValue]() {
											goto l129
										}
										goto l117
									l129:
										position, tokenIndex = position117, tokenIndex117
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
										goto l117
									l130:
										position, tokenIndex = position117, tokenIndex117
										{
											switch buffer[position] {
											case '\'':
												if !_rules[ruleSingleQuote]() {
													goto l76
												}
												{
													position137 := position
													if !_rules[ruleSingleQuotedValue]() {
														goto l76
													}
													add(rulePegText, position137)
												}
												{
													add(ruleAction9, position)
												}
												if !_rules[ruleSingleQuote]() {
													goto l76
												}
												break
											case '"':
												if !_rules[ruleDoubleQuote]() {
													goto l76
												}
												{
													position139 := position
													if !_rules[ruleDoubleQuotedValue]() {
														goto l76
													}
													add(rulePegText, position139)
												}
												{
													add(ruleAction8, position)
												}
												if !_rules[ruleDoubleQuote]() {
													goto l76
												}
												break
											case '$':
												{
													position141 := position
													if buffer[position] != rune('$') {
														goto l76
													}
													position++
													{
														position142 := position
														if !_rules[ruleIdentifier]() {
															goto l76
														}
														add(rulePegText, position142)
													}
													add(ruleRefValue, position141)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position144 := position
													if buffer[position] != rune('{') {
														goto l76
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l76
													}
													{
														position145 := position
														if !_rules[ruleIdentifier]() {
															goto l76
														}
														add(rulePegText, position145)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l76
													}
													if buffer[position] != rune('}') {
														goto l76
													}
													position++
													add(ruleHoleValue, position144)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position147 := position
													if !_rules[ruleStringValue]() {
														goto l76
													}
													add(rulePegText, position147)
												}
												{
													add(ruleAction11, position)
												}
												break
											}
										}

									}
								l117:
									add(ruleValue, position116)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l76
								}
								add(ruleParam, position113)
							}
							goto l75
						l76:
							position, tokenIndex = position76, tokenIndex76
						}
						add(ruleParams, position74)
					}
					goto l73
				l72:
					position, tokenIndex = position72, tokenIndex72
				}
			l73:
				{
					add(ruleAction3, position)
				}
				add(ruleExpr, position57)
			}
			return true
		l56:
			position, tokenIndex = position56, tokenIndex56
			return false
		},
		/* 6 Params <- <Param+> */
		nil,
		/* 7 Param <- <(<Identifier> Action4 Equal Value WhiteSpacing)> */
		nil,
		/* 8 Identifier <- <((&('.') '.') | (&('_') '_') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position152, tokenIndex152 := position, tokenIndex
			{
				position153 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l152
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l152
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l152
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l152
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l152
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l152
						}
						position++
						break
					}
				}

			l154:
				{
					position155, tokenIndex155 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l155
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l155
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l155
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l155
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l155
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l155
							}
							position++
							break
						}
					}

					goto l154
				l155:
					position, tokenIndex = position155, tokenIndex155
				}
				add(ruleIdentifier, position153)
			}
			return true
		l152:
			position, tokenIndex = position152, tokenIndex152
			return false
		},
		/* 9 Value <- <((AliasValue Action7) / (DoubleQuote CustomTypedValue DoubleQuote) / (SingleQuote CustomTypedValue SingleQuote) / CustomTypedValue / (<IntValue> Action10) / ((&('\'') (SingleQuote <SingleQuotedValue> Action9 SingleQuote)) | (&('"') (DoubleQuote <DoubleQuotedValue> Action8 DoubleQuote)) | (&('$') (RefValue Action6)) | (&('{') (HoleValue Action5)) | (&('+' | '-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | ';' | '<' | '>' | '@' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | '~') (<StringValue> Action11))))> */
		nil,
		/* 10 CustomTypedValue <- <((<CidrValue> Action12) / (<IpValue> Action13) / (<CSVValue> Action14) / (<IntRangeValue> Action15))> */
		func() bool {
			position159, tokenIndex159 := position, tokenIndex
			{
				position160 := position
				{
					position161, tokenIndex161 := position, tokenIndex
					{
						position163 := position
						{
							position164 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l162
							}
							position++
						l165:
							{
								position166, tokenIndex166 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l166
								}
								position++
								goto l165
							l166:
								position, tokenIndex = position166, tokenIndex166
							}
							if !matchDot() {
								goto l162
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l162
							}
							position++
						l167:
							{
								position168, tokenIndex168 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l168
								}
								position++
								goto l167
							l168:
								position, tokenIndex = position168, tokenIndex168
							}
							if !matchDot() {
								goto l162
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l162
							}
							position++
						l169:
							{
								position170, tokenIndex170 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l170
								}
								position++
								goto l169
							l170:
								position, tokenIndex = position170, tokenIndex170
							}
							if !matchDot() {
								goto l162
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l162
							}
							position++
						l171:
							{
								position172, tokenIndex172 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l172
								}
								position++
								goto l171
							l172:
								position, tokenIndex = position172, tokenIndex172
							}
							if buffer[position] != rune('/') {
								goto l162
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l162
							}
							position++
						l173:
							{
								position174, tokenIndex174 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l174
								}
								position++
								goto l173
							l174:
								position, tokenIndex = position174, tokenIndex174
							}
							add(ruleCidrValue, position164)
						}
						add(rulePegText, position163)
					}
					{
						add(ruleAction12, position)
					}
					goto l161
				l162:
					position, tokenIndex = position161, tokenIndex161
					{
						position177 := position
						{
							position178 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l176
							}
							position++
						l179:
							{
								position180, tokenIndex180 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l180
								}
								position++
								goto l179
							l180:
								position, tokenIndex = position180, tokenIndex180
							}
							if !matchDot() {
								goto l176
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l176
							}
							position++
						l181:
							{
								position182, tokenIndex182 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l182
								}
								position++
								goto l181
							l182:
								position, tokenIndex = position182, tokenIndex182
							}
							if !matchDot() {
								goto l176
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l176
							}
							position++
						l183:
							{
								position184, tokenIndex184 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l184
								}
								position++
								goto l183
							l184:
								position, tokenIndex = position184, tokenIndex184
							}
							if !matchDot() {
								goto l176
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l176
							}
							position++
						l185:
							{
								position186, tokenIndex186 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l186
								}
								position++
								goto l185
							l186:
								position, tokenIndex = position186, tokenIndex186
							}
							add(ruleIpValue, position178)
						}
						add(rulePegText, position177)
					}
					{
						add(ruleAction13, position)
					}
					goto l161
				l176:
					position, tokenIndex = position161, tokenIndex161
					{
						position189 := position
						{
							position190 := position
							if !_rules[ruleStringValue]() {
								goto l188
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l188
							}
							if buffer[position] != rune(',') {
								goto l188
							}
							position++
							if !_rules[ruleWhiteSpacing]() {
								goto l188
							}
						l191:
							{
								position192, tokenIndex192 := position, tokenIndex
								if !_rules[ruleStringValue]() {
									goto l192
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l192
								}
								if buffer[position] != rune(',') {
									goto l192
								}
								position++
								if !_rules[ruleWhiteSpacing]() {
									goto l192
								}
								goto l191
							l192:
								position, tokenIndex = position192, tokenIndex192
							}
							if !_rules[ruleStringValue]() {
								goto l188
							}
							add(ruleCSVValue, position190)
						}
						add(rulePegText, position189)
					}
					{
						add(ruleAction14, position)
					}
					goto l161
				l188:
					position, tokenIndex = position161, tokenIndex161
					{
						position194 := position
						{
							position195 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l159
							}
							position++
						l196:
							{
								position197, tokenIndex197 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l197
								}
								position++
								goto l196
							l197:
								position, tokenIndex = position197, tokenIndex197
							}
							if buffer[position] != rune('-') {
								goto l159
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l159
							}
							position++
						l198:
							{
								position199, tokenIndex199 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l199
								}
								position++
								goto l198
							l199:
								position, tokenIndex = position199, tokenIndex199
							}
							add(ruleIntRangeValue, position195)
						}
						add(rulePegText, position194)
					}
					{
						add(ruleAction15, position)
					}
				}
			l161:
				add(ruleCustomTypedValue, position160)
			}
			return true
		l159:
			position, tokenIndex = position159, tokenIndex159
			return false
		},
		/* 11 StringValue <- <((&('>') '>') | (&('<') '<') | (&('@') '@') | (&('~') '~') | (&(';') ';') | (&('+') '+') | (&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position201, tokenIndex201 := position, tokenIndex
			{
				position202 := position
				{
					switch buffer[position] {
					case '>':
						if buffer[position] != rune('>') {
							goto l201
						}
						position++
						break
					case '<':
						if buffer[position] != rune('<') {
							goto l201
						}
						position++
						break
					case '@':
						if buffer[position] != rune('@') {
							goto l201
						}
						position++
						break
					case '~':
						if buffer[position] != rune('~') {
							goto l201
						}
						position++
						break
					case ';':
						if buffer[position] != rune(';') {
							goto l201
						}
						position++
						break
					case '+':
						if buffer[position] != rune('+') {
							goto l201
						}
						position++
						break
					case '/':
						if buffer[position] != rune('/') {
							goto l201
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l201
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l201
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l201
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l201
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l201
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l201
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l201
						}
						position++
						break
					}
				}

			l203:
				{
					position204, tokenIndex204 := position, tokenIndex
					{
						switch buffer[position] {
						case '>':
							if buffer[position] != rune('>') {
								goto l204
							}
							position++
							break
						case '<':
							if buffer[position] != rune('<') {
								goto l204
							}
							position++
							break
						case '@':
							if buffer[position] != rune('@') {
								goto l204
							}
							position++
							break
						case '~':
							if buffer[position] != rune('~') {
								goto l204
							}
							position++
							break
						case ';':
							if buffer[position] != rune(';') {
								goto l204
							}
							position++
							break
						case '+':
							if buffer[position] != rune('+') {
								goto l204
							}
							position++
							break
						case '/':
							if buffer[position] != rune('/') {
								goto l204
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l204
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l204
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l204
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l204
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l204
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l204
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l204
							}
							position++
							break
						}
					}

					goto l203
				l204:
					position, tokenIndex = position204, tokenIndex204
				}
				add(ruleStringValue, position202)
			}
			return true
		l201:
			position, tokenIndex = position201, tokenIndex201
			return false
		},
		/* 12 DoubleQuotedValue <- <(!'"' .)*> */
		func() bool {
			{
				position208 := position
			l209:
				{
					position210, tokenIndex210 := position, tokenIndex
					{
						position211, tokenIndex211 := position, tokenIndex
						if buffer[position] != rune('"') {
							goto l211
						}
						position++
						goto l210
					l211:
						position, tokenIndex = position211, tokenIndex211
					}
					if !matchDot() {
						goto l210
					}
					goto l209
				l210:
					position, tokenIndex = position210, tokenIndex210
				}
				add(ruleDoubleQuotedValue, position208)
			}
			return true
		},
		/* 13 SingleQuotedValue <- <(!'\'' .)*> */
		func() bool {
			{
				position213 := position
			l214:
				{
					position215, tokenIndex215 := position, tokenIndex
					{
						position216, tokenIndex216 := position, tokenIndex
						if buffer[position] != rune('\'') {
							goto l216
						}
						position++
						goto l215
					l216:
						position, tokenIndex = position216, tokenIndex216
					}
					if !matchDot() {
						goto l215
					}
					goto l214
				l215:
					position, tokenIndex = position215, tokenIndex215
				}
				add(ruleSingleQuotedValue, position213)
			}
			return true
		},
		/* 14 CSVValue <- <((StringValue WhiteSpacing ',' WhiteSpacing)+ StringValue)> */
		nil,
		/* 15 CidrValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+ '/' [0-9]+)> */
		nil,
		/* 16 IpValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+)> */
		nil,
		/* 17 IntValue <- <[0-9]+> */
		nil,
		/* 18 IntRangeValue <- <([0-9]+ '-' [0-9]+)> */
		nil,
		/* 19 RefValue <- <('$' <Identifier>)> */
		nil,
		/* 20 AliasValue <- <(('@' <StringValue>) / ('@' DoubleQuote <DoubleQuotedValue> DoubleQuote) / ('@' SingleQuote <SingleQuotedValue> SingleQuote))> */
		nil,
		/* 21 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 22 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)* Action16))> */
		nil,
		/* 23 SingleQuote <- <'\''> */
		func() bool {
			position226, tokenIndex226 := position, tokenIndex
			{
				position227 := position
				if buffer[position] != rune('\'') {
					goto l226
				}
				position++
				add(ruleSingleQuote, position227)
			}
			return true
		l226:
			position, tokenIndex = position226, tokenIndex226
			return false
		},
		/* 24 DoubleQuote <- <'"'> */
		func() bool {
			position228, tokenIndex228 := position, tokenIndex
			{
				position229 := position
				if buffer[position] != rune('"') {
					goto l228
				}
				position++
				add(ruleDoubleQuote, position229)
			}
			return true
		l228:
			position, tokenIndex = position228, tokenIndex228
			return false
		},
		/* 25 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position231 := position
			l232:
				{
					position233, tokenIndex233 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l233
					}
					goto l232
				l233:
					position, tokenIndex = position233, tokenIndex233
				}
				add(ruleWhiteSpacing, position231)
			}
			return true
		},
		/* 26 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position234, tokenIndex234 := position, tokenIndex
			{
				position235 := position
				if !_rules[ruleWhitespace]() {
					goto l234
				}
			l236:
				{
					position237, tokenIndex237 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l237
					}
					goto l236
				l237:
					position, tokenIndex = position237, tokenIndex237
				}
				add(ruleMustWhiteSpacing, position235)
			}
			return true
		l234:
			position, tokenIndex = position234, tokenIndex234
			return false
		},
		/* 27 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position238, tokenIndex238 := position, tokenIndex
			{
				position239 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l238
				}
				if buffer[position] != rune('=') {
					goto l238
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l238
				}
				add(ruleEqual, position239)
			}
			return true
		l238:
			position, tokenIndex = position238, tokenIndex238
			return false
		},
		/* 28 BlankLine <- <(WhiteSpacing EndOfLine Action17)> */
		func() bool {
			position240, tokenIndex240 := position, tokenIndex
			{
				position241 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l240
				}
				if !_rules[ruleEndOfLine]() {
					goto l240
				}
				{
					add(ruleAction17, position)
				}
				add(ruleBlankLine, position241)
			}
			return true
		l240:
			position, tokenIndex = position240, tokenIndex240
			return false
		},
		/* 29 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position243, tokenIndex243 := position, tokenIndex
			{
				position244 := position
				{
					position245, tokenIndex245 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l246
					}
					position++
					goto l245
				l246:
					position, tokenIndex = position245, tokenIndex245
					if buffer[position] != rune('\t') {
						goto l243
					}
					position++
				}
			l245:
				add(ruleWhitespace, position244)
			}
			return true
		l243:
			position, tokenIndex = position243, tokenIndex243
			return false
		},
		/* 30 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position247, tokenIndex247 := position, tokenIndex
			{
				position248 := position
				{
					position249, tokenIndex249 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l250
					}
					position++
					if buffer[position] != rune('\n') {
						goto l250
					}
					position++
					goto l249
				l250:
					position, tokenIndex = position249, tokenIndex249
					if buffer[position] != rune('\n') {
						goto l251
					}
					position++
					goto l249
				l251:
					position, tokenIndex = position249, tokenIndex249
					if buffer[position] != rune('\r') {
						goto l247
					}
					position++
				}
			l249:
				add(ruleEndOfLine, position248)
			}
			return true
		l247:
			position, tokenIndex = position247, tokenIndex247
			return false
		},
		/* 31 EndOfFile <- <!.> */
		nil,
		nil,
		/* 34 Action0 <- <{ p.addDeclarationIdentifier(text) }> */
		nil,
		/* 35 Action1 <- <{ p.addAction(text) }> */
		nil,
		/* 36 Action2 <- <{ p.addEntity(text) }> */
		nil,
		/* 37 Action3 <- <{ p.LineDone() }> */
		nil,
		/* 38 Action4 <- <{ p.addParamKey(text) }> */
		nil,
		/* 39 Action5 <- <{  p.addParamHoleValue(text) }> */
		nil,
		/* 40 Action6 <- <{  p.addParamRefValue(text) }> */
		nil,
		/* 41 Action7 <- <{  p.addAliasParam(text) }> */
		nil,
		/* 42 Action8 <- <{ p.addParamValue(text) }> */
		nil,
		/* 43 Action9 <- <{ p.addParamValue(text) }> */
		nil,
		/* 44 Action10 <- <{ p.addParamIntValue(text) }> */
		nil,
		/* 45 Action11 <- <{ p.addParamValue(text) }> */
		nil,
		/* 46 Action12 <- <{ p.addParamCidrValue(text) }> */
		nil,
		/* 47 Action13 <- <{ p.addParamIpValue(text) }> */
		nil,
		/* 48 Action14 <- <{p.addCsvValue(text)}> */
		nil,
		/* 49 Action15 <- <{ p.addParamValue(text) }> */
		nil,
		/* 50 Action16 <- <{ p.LineDone() }> */
		nil,
		/* 51 Action17 <- <{ p.LineDone() }> */
		nil,
	}
	p.rules = _rules
}
