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
		/* 3 Entity <- <[a-z]+> */
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
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l56
						}
						position++
					l65:
						{
							position66, tokenIndex66 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l66
							}
							position++
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
					position68, tokenIndex68 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l68
					}
					{
						position70 := position
						{
							position73 := position
							{
								position74 := position
								if !_rules[ruleIdentifier]() {
									goto l68
								}
								add(rulePegText, position74)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l68
							}
							{
								position76 := position
								{
									position77, tokenIndex77 := position, tokenIndex
									{
										position79 := position
										{
											position80, tokenIndex80 := position, tokenIndex
											if buffer[position] != rune('@') {
												goto l81
											}
											position++
											{
												position82 := position
												if !_rules[ruleStringValue]() {
													goto l81
												}
												add(rulePegText, position82)
											}
											goto l80
										l81:
											position, tokenIndex = position80, tokenIndex80
											if buffer[position] != rune('@') {
												goto l83
											}
											position++
											if !_rules[ruleDoubleQuote]() {
												goto l83
											}
											{
												position84 := position
												if !_rules[ruleDoubleQuotedValue]() {
													goto l83
												}
												add(rulePegText, position84)
											}
											if !_rules[ruleDoubleQuote]() {
												goto l83
											}
											goto l80
										l83:
											position, tokenIndex = position80, tokenIndex80
											if buffer[position] != rune('@') {
												goto l78
											}
											position++
											if !_rules[ruleSingleQuote]() {
												goto l78
											}
											{
												position85 := position
												if !_rules[ruleSingleQuotedValue]() {
													goto l78
												}
												add(rulePegText, position85)
											}
											if !_rules[ruleSingleQuote]() {
												goto l78
											}
										}
									l80:
										add(ruleAliasValue, position79)
									}
									{
										add(ruleAction7, position)
									}
									goto l77
								l78:
									position, tokenIndex = position77, tokenIndex77
									if !_rules[ruleDoubleQuote]() {
										goto l87
									}
									if !_rules[ruleCustomTypedValue]() {
										goto l87
									}
									if !_rules[ruleDoubleQuote]() {
										goto l87
									}
									goto l77
								l87:
									position, tokenIndex = position77, tokenIndex77
									if !_rules[ruleSingleQuote]() {
										goto l88
									}
									if !_rules[ruleCustomTypedValue]() {
										goto l88
									}
									if !_rules[ruleSingleQuote]() {
										goto l88
									}
									goto l77
								l88:
									position, tokenIndex = position77, tokenIndex77
									if !_rules[ruleCustomTypedValue]() {
										goto l89
									}
									goto l77
								l89:
									position, tokenIndex = position77, tokenIndex77
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
										add(ruleAction10, position)
									}
									goto l77
								l90:
									position, tokenIndex = position77, tokenIndex77
									{
										switch buffer[position] {
										case '\'':
											if !_rules[ruleSingleQuote]() {
												goto l68
											}
											{
												position97 := position
												if !_rules[ruleSingleQuotedValue]() {
													goto l68
												}
												add(rulePegText, position97)
											}
											{
												add(ruleAction9, position)
											}
											if !_rules[ruleSingleQuote]() {
												goto l68
											}
											break
										case '"':
											if !_rules[ruleDoubleQuote]() {
												goto l68
											}
											{
												position99 := position
												if !_rules[ruleDoubleQuotedValue]() {
													goto l68
												}
												add(rulePegText, position99)
											}
											{
												add(ruleAction8, position)
											}
											if !_rules[ruleDoubleQuote]() {
												goto l68
											}
											break
										case '$':
											{
												position101 := position
												if buffer[position] != rune('$') {
													goto l68
												}
												position++
												{
													position102 := position
													if !_rules[ruleIdentifier]() {
														goto l68
													}
													add(rulePegText, position102)
												}
												add(ruleRefValue, position101)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position104 := position
												if buffer[position] != rune('{') {
													goto l68
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l68
												}
												{
													position105 := position
													if !_rules[ruleIdentifier]() {
														goto l68
													}
													add(rulePegText, position105)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l68
												}
												if buffer[position] != rune('}') {
													goto l68
												}
												position++
												add(ruleHoleValue, position104)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position107 := position
												if !_rules[ruleStringValue]() {
													goto l68
												}
												add(rulePegText, position107)
											}
											{
												add(ruleAction11, position)
											}
											break
										}
									}

								}
							l77:
								add(ruleValue, position76)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l68
							}
							add(ruleParam, position73)
						}
					l71:
						{
							position72, tokenIndex72 := position, tokenIndex
							{
								position109 := position
								{
									position110 := position
									if !_rules[ruleIdentifier]() {
										goto l72
									}
									add(rulePegText, position110)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l72
								}
								{
									position112 := position
									{
										position113, tokenIndex113 := position, tokenIndex
										{
											position115 := position
											{
												position116, tokenIndex116 := position, tokenIndex
												if buffer[position] != rune('@') {
													goto l117
												}
												position++
												{
													position118 := position
													if !_rules[ruleStringValue]() {
														goto l117
													}
													add(rulePegText, position118)
												}
												goto l116
											l117:
												position, tokenIndex = position116, tokenIndex116
												if buffer[position] != rune('@') {
													goto l119
												}
												position++
												if !_rules[ruleDoubleQuote]() {
													goto l119
												}
												{
													position120 := position
													if !_rules[ruleDoubleQuotedValue]() {
														goto l119
													}
													add(rulePegText, position120)
												}
												if !_rules[ruleDoubleQuote]() {
													goto l119
												}
												goto l116
											l119:
												position, tokenIndex = position116, tokenIndex116
												if buffer[position] != rune('@') {
													goto l114
												}
												position++
												if !_rules[ruleSingleQuote]() {
													goto l114
												}
												{
													position121 := position
													if !_rules[ruleSingleQuotedValue]() {
														goto l114
													}
													add(rulePegText, position121)
												}
												if !_rules[ruleSingleQuote]() {
													goto l114
												}
											}
										l116:
											add(ruleAliasValue, position115)
										}
										{
											add(ruleAction7, position)
										}
										goto l113
									l114:
										position, tokenIndex = position113, tokenIndex113
										if !_rules[ruleDoubleQuote]() {
											goto l123
										}
										if !_rules[ruleCustomTypedValue]() {
											goto l123
										}
										if !_rules[ruleDoubleQuote]() {
											goto l123
										}
										goto l113
									l123:
										position, tokenIndex = position113, tokenIndex113
										if !_rules[ruleSingleQuote]() {
											goto l124
										}
										if !_rules[ruleCustomTypedValue]() {
											goto l124
										}
										if !_rules[ruleSingleQuote]() {
											goto l124
										}
										goto l113
									l124:
										position, tokenIndex = position113, tokenIndex113
										if !_rules[ruleCustomTypedValue]() {
											goto l125
										}
										goto l113
									l125:
										position, tokenIndex = position113, tokenIndex113
										{
											position127 := position
											{
												position128 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l126
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
												add(ruleIntValue, position128)
											}
											add(rulePegText, position127)
										}
										{
											add(ruleAction10, position)
										}
										goto l113
									l126:
										position, tokenIndex = position113, tokenIndex113
										{
											switch buffer[position] {
											case '\'':
												if !_rules[ruleSingleQuote]() {
													goto l72
												}
												{
													position133 := position
													if !_rules[ruleSingleQuotedValue]() {
														goto l72
													}
													add(rulePegText, position133)
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
													position135 := position
													if !_rules[ruleDoubleQuotedValue]() {
														goto l72
													}
													add(rulePegText, position135)
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
													position137 := position
													if buffer[position] != rune('$') {
														goto l72
													}
													position++
													{
														position138 := position
														if !_rules[ruleIdentifier]() {
															goto l72
														}
														add(rulePegText, position138)
													}
													add(ruleRefValue, position137)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position140 := position
													if buffer[position] != rune('{') {
														goto l72
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l72
													}
													{
														position141 := position
														if !_rules[ruleIdentifier]() {
															goto l72
														}
														add(rulePegText, position141)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l72
													}
													if buffer[position] != rune('}') {
														goto l72
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
													if !_rules[ruleStringValue]() {
														goto l72
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
								l113:
									add(ruleValue, position112)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l72
								}
								add(ruleParam, position109)
							}
							goto l71
						l72:
							position, tokenIndex = position72, tokenIndex72
						}
						add(ruleParams, position70)
					}
					goto l69
				l68:
					position, tokenIndex = position68, tokenIndex68
				}
			l69:
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
			position148, tokenIndex148 := position, tokenIndex
			{
				position149 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l148
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l148
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l148
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l148
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l148
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l148
						}
						position++
						break
					}
				}

			l150:
				{
					position151, tokenIndex151 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l151
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l151
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l151
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l151
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l151
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l151
							}
							position++
							break
						}
					}

					goto l150
				l151:
					position, tokenIndex = position151, tokenIndex151
				}
				add(ruleIdentifier, position149)
			}
			return true
		l148:
			position, tokenIndex = position148, tokenIndex148
			return false
		},
		/* 9 Value <- <((AliasValue Action7) / (DoubleQuote CustomTypedValue DoubleQuote) / (SingleQuote CustomTypedValue SingleQuote) / CustomTypedValue / (<IntValue> Action10) / ((&('\'') (SingleQuote <SingleQuotedValue> Action9 SingleQuote)) | (&('"') (DoubleQuote <DoubleQuotedValue> Action8 DoubleQuote)) | (&('$') (RefValue Action6)) | (&('{') (HoleValue Action5)) | (&('+' | '-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | ';' | '<' | '>' | '@' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | '~') (<StringValue> Action11))))> */
		nil,
		/* 10 CustomTypedValue <- <((<CidrValue> Action12) / (<IpValue> Action13) / (<CSVValue> Action14) / (<IntRangeValue> Action15))> */
		func() bool {
			position155, tokenIndex155 := position, tokenIndex
			{
				position156 := position
				{
					position157, tokenIndex157 := position, tokenIndex
					{
						position159 := position
						{
							position160 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l158
							}
							position++
						l161:
							{
								position162, tokenIndex162 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l162
								}
								position++
								goto l161
							l162:
								position, tokenIndex = position162, tokenIndex162
							}
							if !matchDot() {
								goto l158
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l158
							}
							position++
						l163:
							{
								position164, tokenIndex164 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l164
								}
								position++
								goto l163
							l164:
								position, tokenIndex = position164, tokenIndex164
							}
							if !matchDot() {
								goto l158
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l158
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
								goto l158
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l158
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
							if buffer[position] != rune('/') {
								goto l158
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l158
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
							add(ruleCidrValue, position160)
						}
						add(rulePegText, position159)
					}
					{
						add(ruleAction12, position)
					}
					goto l157
				l158:
					position, tokenIndex = position157, tokenIndex157
					{
						position173 := position
						{
							position174 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l172
							}
							position++
						l175:
							{
								position176, tokenIndex176 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l176
								}
								position++
								goto l175
							l176:
								position, tokenIndex = position176, tokenIndex176
							}
							if !matchDot() {
								goto l172
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l172
							}
							position++
						l177:
							{
								position178, tokenIndex178 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l178
								}
								position++
								goto l177
							l178:
								position, tokenIndex = position178, tokenIndex178
							}
							if !matchDot() {
								goto l172
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l172
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
								goto l172
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l172
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
							add(ruleIpValue, position174)
						}
						add(rulePegText, position173)
					}
					{
						add(ruleAction13, position)
					}
					goto l157
				l172:
					position, tokenIndex = position157, tokenIndex157
					{
						position185 := position
						{
							position186 := position
							if !_rules[ruleStringValue]() {
								goto l184
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l184
							}
							if buffer[position] != rune(',') {
								goto l184
							}
							position++
							if !_rules[ruleWhiteSpacing]() {
								goto l184
							}
						l187:
							{
								position188, tokenIndex188 := position, tokenIndex
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
								goto l187
							l188:
								position, tokenIndex = position188, tokenIndex188
							}
							if !_rules[ruleStringValue]() {
								goto l184
							}
							add(ruleCSVValue, position186)
						}
						add(rulePegText, position185)
					}
					{
						add(ruleAction14, position)
					}
					goto l157
				l184:
					position, tokenIndex = position157, tokenIndex157
					{
						position190 := position
						{
							position191 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l155
							}
							position++
						l192:
							{
								position193, tokenIndex193 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l193
								}
								position++
								goto l192
							l193:
								position, tokenIndex = position193, tokenIndex193
							}
							if buffer[position] != rune('-') {
								goto l155
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l155
							}
							position++
						l194:
							{
								position195, tokenIndex195 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l195
								}
								position++
								goto l194
							l195:
								position, tokenIndex = position195, tokenIndex195
							}
							add(ruleIntRangeValue, position191)
						}
						add(rulePegText, position190)
					}
					{
						add(ruleAction15, position)
					}
				}
			l157:
				add(ruleCustomTypedValue, position156)
			}
			return true
		l155:
			position, tokenIndex = position155, tokenIndex155
			return false
		},
		/* 11 StringValue <- <((&('>') '>') | (&('<') '<') | (&('@') '@') | (&('~') '~') | (&(';') ';') | (&('+') '+') | (&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position197, tokenIndex197 := position, tokenIndex
			{
				position198 := position
				{
					switch buffer[position] {
					case '>':
						if buffer[position] != rune('>') {
							goto l197
						}
						position++
						break
					case '<':
						if buffer[position] != rune('<') {
							goto l197
						}
						position++
						break
					case '@':
						if buffer[position] != rune('@') {
							goto l197
						}
						position++
						break
					case '~':
						if buffer[position] != rune('~') {
							goto l197
						}
						position++
						break
					case ';':
						if buffer[position] != rune(';') {
							goto l197
						}
						position++
						break
					case '+':
						if buffer[position] != rune('+') {
							goto l197
						}
						position++
						break
					case '/':
						if buffer[position] != rune('/') {
							goto l197
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l197
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l197
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l197
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l197
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l197
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l197
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l197
						}
						position++
						break
					}
				}

			l199:
				{
					position200, tokenIndex200 := position, tokenIndex
					{
						switch buffer[position] {
						case '>':
							if buffer[position] != rune('>') {
								goto l200
							}
							position++
							break
						case '<':
							if buffer[position] != rune('<') {
								goto l200
							}
							position++
							break
						case '@':
							if buffer[position] != rune('@') {
								goto l200
							}
							position++
							break
						case '~':
							if buffer[position] != rune('~') {
								goto l200
							}
							position++
							break
						case ';':
							if buffer[position] != rune(';') {
								goto l200
							}
							position++
							break
						case '+':
							if buffer[position] != rune('+') {
								goto l200
							}
							position++
							break
						case '/':
							if buffer[position] != rune('/') {
								goto l200
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l200
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l200
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l200
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l200
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l200
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l200
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l200
							}
							position++
							break
						}
					}

					goto l199
				l200:
					position, tokenIndex = position200, tokenIndex200
				}
				add(ruleStringValue, position198)
			}
			return true
		l197:
			position, tokenIndex = position197, tokenIndex197
			return false
		},
		/* 12 DoubleQuotedValue <- <(!'"' .)*> */
		func() bool {
			{
				position204 := position
			l205:
				{
					position206, tokenIndex206 := position, tokenIndex
					{
						position207, tokenIndex207 := position, tokenIndex
						if buffer[position] != rune('"') {
							goto l207
						}
						position++
						goto l206
					l207:
						position, tokenIndex = position207, tokenIndex207
					}
					if !matchDot() {
						goto l206
					}
					goto l205
				l206:
					position, tokenIndex = position206, tokenIndex206
				}
				add(ruleDoubleQuotedValue, position204)
			}
			return true
		},
		/* 13 SingleQuotedValue <- <(!'\'' .)*> */
		func() bool {
			{
				position209 := position
			l210:
				{
					position211, tokenIndex211 := position, tokenIndex
					{
						position212, tokenIndex212 := position, tokenIndex
						if buffer[position] != rune('\'') {
							goto l212
						}
						position++
						goto l211
					l212:
						position, tokenIndex = position212, tokenIndex212
					}
					if !matchDot() {
						goto l211
					}
					goto l210
				l211:
					position, tokenIndex = position211, tokenIndex211
				}
				add(ruleSingleQuotedValue, position209)
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
			position222, tokenIndex222 := position, tokenIndex
			{
				position223 := position
				if buffer[position] != rune('\'') {
					goto l222
				}
				position++
				add(ruleSingleQuote, position223)
			}
			return true
		l222:
			position, tokenIndex = position222, tokenIndex222
			return false
		},
		/* 24 DoubleQuote <- <'"'> */
		func() bool {
			position224, tokenIndex224 := position, tokenIndex
			{
				position225 := position
				if buffer[position] != rune('"') {
					goto l224
				}
				position++
				add(ruleDoubleQuote, position225)
			}
			return true
		l224:
			position, tokenIndex = position224, tokenIndex224
			return false
		},
		/* 25 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position227 := position
			l228:
				{
					position229, tokenIndex229 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l229
					}
					goto l228
				l229:
					position, tokenIndex = position229, tokenIndex229
				}
				add(ruleWhiteSpacing, position227)
			}
			return true
		},
		/* 26 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position230, tokenIndex230 := position, tokenIndex
			{
				position231 := position
				if !_rules[ruleWhitespace]() {
					goto l230
				}
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
				add(ruleMustWhiteSpacing, position231)
			}
			return true
		l230:
			position, tokenIndex = position230, tokenIndex230
			return false
		},
		/* 27 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position234, tokenIndex234 := position, tokenIndex
			{
				position235 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l234
				}
				if buffer[position] != rune('=') {
					goto l234
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l234
				}
				add(ruleEqual, position235)
			}
			return true
		l234:
			position, tokenIndex = position234, tokenIndex234
			return false
		},
		/* 28 BlankLine <- <(WhiteSpacing EndOfLine Action17)> */
		func() bool {
			position236, tokenIndex236 := position, tokenIndex
			{
				position237 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l236
				}
				if !_rules[ruleEndOfLine]() {
					goto l236
				}
				{
					add(ruleAction17, position)
				}
				add(ruleBlankLine, position237)
			}
			return true
		l236:
			position, tokenIndex = position236, tokenIndex236
			return false
		},
		/* 29 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position239, tokenIndex239 := position, tokenIndex
			{
				position240 := position
				{
					position241, tokenIndex241 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex = position241, tokenIndex241
					if buffer[position] != rune('\t') {
						goto l239
					}
					position++
				}
			l241:
				add(ruleWhitespace, position240)
			}
			return true
		l239:
			position, tokenIndex = position239, tokenIndex239
			return false
		},
		/* 30 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position243, tokenIndex243 := position, tokenIndex
			{
				position244 := position
				{
					position245, tokenIndex245 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l246
					}
					position++
					if buffer[position] != rune('\n') {
						goto l246
					}
					position++
					goto l245
				l246:
					position, tokenIndex = position245, tokenIndex245
					if buffer[position] != rune('\n') {
						goto l247
					}
					position++
					goto l245
				l247:
					position, tokenIndex = position245, tokenIndex245
					if buffer[position] != rune('\r') {
						goto l243
					}
					position++
				}
			l245:
				add(ruleEndOfLine, position244)
			}
			return true
		l243:
			position, tokenIndex = position243, tokenIndex243
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
