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
	ruleCSVValue
	ruleCidrValue
	ruleIpValue
	ruleIntValue
	ruleIntRangeValue
	ruleRefValue
	ruleAliasValue
	ruleHoleValue
	ruleComment
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
	"CSVValue",
	"CidrValue",
	"IpValue",
	"IntValue",
	"IntRangeValue",
	"RefValue",
	"AliasValue",
	"HoleValue",
	"Comment",
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
	rules  [45]func() bool
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
			p.addParamValue(text)
		case ruleAction7:
			p.addParamRefValue(text)
		case ruleAction8:
			p.addParamCidrValue(text)
		case ruleAction9:
			p.addParamIpValue(text)
		case ruleAction10:
			p.addCsvValue(text)
		case ruleAction11:
			p.addParamValue(text)
		case ruleAction12:
			p.addParamIntValue(text)
		case ruleAction13:
			p.addParamValue(text)
		case ruleAction14:
			p.LineDone()
		case ruleAction15:
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
									add(ruleAction14, position)
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
										add(ruleAction14, position)
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
											if !matchDot() {
												goto l78
											}
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
											if !matchDot() {
												goto l78
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l78
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
											if !matchDot() {
												goto l78
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l78
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
											if buffer[position] != rune('/') {
												goto l78
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l78
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
											add(ruleCidrValue, position80)
										}
										add(rulePegText, position79)
									}
									{
										add(ruleAction8, position)
									}
									goto l77
								l78:
									position, tokenIndex = position77, tokenIndex77
									{
										position93 := position
										{
											position94 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l92
											}
											position++
										l95:
											{
												position96, tokenIndex96 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l96
												}
												position++
												goto l95
											l96:
												position, tokenIndex = position96, tokenIndex96
											}
											if !matchDot() {
												goto l92
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l92
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
											if !matchDot() {
												goto l92
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l92
											}
											position++
										l99:
											{
												position100, tokenIndex100 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l100
												}
												position++
												goto l99
											l100:
												position, tokenIndex = position100, tokenIndex100
											}
											if !matchDot() {
												goto l92
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l92
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
											add(ruleIpValue, position94)
										}
										add(rulePegText, position93)
									}
									{
										add(ruleAction9, position)
									}
									goto l77
								l92:
									position, tokenIndex = position77, tokenIndex77
									{
										position105 := position
										{
											position106 := position
											if !_rules[ruleStringValue]() {
												goto l104
											}
											if !_rules[ruleWhiteSpacing]() {
												goto l104
											}
											if buffer[position] != rune(',') {
												goto l104
											}
											position++
											if !_rules[ruleWhiteSpacing]() {
												goto l104
											}
										l107:
											{
												position108, tokenIndex108 := position, tokenIndex
												if !_rules[ruleStringValue]() {
													goto l108
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l108
												}
												if buffer[position] != rune(',') {
													goto l108
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l108
												}
												goto l107
											l108:
												position, tokenIndex = position108, tokenIndex108
											}
											if !_rules[ruleStringValue]() {
												goto l104
											}
											add(ruleCSVValue, position106)
										}
										add(rulePegText, position105)
									}
									{
										add(ruleAction10, position)
									}
									goto l77
								l104:
									position, tokenIndex = position77, tokenIndex77
									{
										position111 := position
										{
											position112 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l110
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
											if buffer[position] != rune('-') {
												goto l110
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l110
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
											add(ruleIntRangeValue, position112)
										}
										add(rulePegText, position111)
									}
									{
										add(ruleAction11, position)
									}
									goto l77
								l110:
									position, tokenIndex = position77, tokenIndex77
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
											add(ruleIntValue, position120)
										}
										add(rulePegText, position119)
									}
									{
										add(ruleAction12, position)
									}
									goto l77
								l118:
									position, tokenIndex = position77, tokenIndex77
									{
										switch buffer[position] {
										case '$':
											{
												position125 := position
												if buffer[position] != rune('$') {
													goto l68
												}
												position++
												{
													position126 := position
													if !_rules[ruleIdentifier]() {
														goto l68
													}
													add(rulePegText, position126)
												}
												add(ruleRefValue, position125)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position128 := position
												{
													position129 := position
													if buffer[position] != rune('@') {
														goto l68
													}
													position++
													if !_rules[ruleStringValue]() {
														goto l68
													}
													add(rulePegText, position129)
												}
												add(ruleAliasValue, position128)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position131 := position
												if buffer[position] != rune('{') {
													goto l68
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l68
												}
												{
													position132 := position
													if !_rules[ruleIdentifier]() {
														goto l68
													}
													add(rulePegText, position132)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l68
												}
												if buffer[position] != rune('}') {
													goto l68
												}
												position++
												add(ruleHoleValue, position131)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position134 := position
												if !_rules[ruleStringValue]() {
													goto l68
												}
												add(rulePegText, position134)
											}
											{
												add(ruleAction13, position)
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
								position136 := position
								{
									position137 := position
									if !_rules[ruleIdentifier]() {
										goto l72
									}
									add(rulePegText, position137)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l72
								}
								{
									position139 := position
									{
										position140, tokenIndex140 := position, tokenIndex
										{
											position142 := position
											{
												position143 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l141
												}
												position++
											l144:
												{
													position145, tokenIndex145 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l145
													}
													position++
													goto l144
												l145:
													position, tokenIndex = position145, tokenIndex145
												}
												if !matchDot() {
													goto l141
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l141
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
												if !matchDot() {
													goto l141
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l141
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
												if !matchDot() {
													goto l141
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l141
												}
												position++
											l150:
												{
													position151, tokenIndex151 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l151
													}
													position++
													goto l150
												l151:
													position, tokenIndex = position151, tokenIndex151
												}
												if buffer[position] != rune('/') {
													goto l141
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l141
												}
												position++
											l152:
												{
													position153, tokenIndex153 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l153
													}
													position++
													goto l152
												l153:
													position, tokenIndex = position153, tokenIndex153
												}
												add(ruleCidrValue, position143)
											}
											add(rulePegText, position142)
										}
										{
											add(ruleAction8, position)
										}
										goto l140
									l141:
										position, tokenIndex = position140, tokenIndex140
										{
											position156 := position
											{
												position157 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l155
												}
												position++
											l158:
												{
													position159, tokenIndex159 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l159
													}
													position++
													goto l158
												l159:
													position, tokenIndex = position159, tokenIndex159
												}
												if !matchDot() {
													goto l155
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l155
												}
												position++
											l160:
												{
													position161, tokenIndex161 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l161
													}
													position++
													goto l160
												l161:
													position, tokenIndex = position161, tokenIndex161
												}
												if !matchDot() {
													goto l155
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l155
												}
												position++
											l162:
												{
													position163, tokenIndex163 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l163
													}
													position++
													goto l162
												l163:
													position, tokenIndex = position163, tokenIndex163
												}
												if !matchDot() {
													goto l155
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l155
												}
												position++
											l164:
												{
													position165, tokenIndex165 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l165
													}
													position++
													goto l164
												l165:
													position, tokenIndex = position165, tokenIndex165
												}
												add(ruleIpValue, position157)
											}
											add(rulePegText, position156)
										}
										{
											add(ruleAction9, position)
										}
										goto l140
									l155:
										position, tokenIndex = position140, tokenIndex140
										{
											position168 := position
											{
												position169 := position
												if !_rules[ruleStringValue]() {
													goto l167
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l167
												}
												if buffer[position] != rune(',') {
													goto l167
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l167
												}
											l170:
												{
													position171, tokenIndex171 := position, tokenIndex
													if !_rules[ruleStringValue]() {
														goto l171
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l171
													}
													if buffer[position] != rune(',') {
														goto l171
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l171
													}
													goto l170
												l171:
													position, tokenIndex = position171, tokenIndex171
												}
												if !_rules[ruleStringValue]() {
													goto l167
												}
												add(ruleCSVValue, position169)
											}
											add(rulePegText, position168)
										}
										{
											add(ruleAction10, position)
										}
										goto l140
									l167:
										position, tokenIndex = position140, tokenIndex140
										{
											position174 := position
											{
												position175 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l173
												}
												position++
											l176:
												{
													position177, tokenIndex177 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l177
													}
													position++
													goto l176
												l177:
													position, tokenIndex = position177, tokenIndex177
												}
												if buffer[position] != rune('-') {
													goto l173
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l173
												}
												position++
											l178:
												{
													position179, tokenIndex179 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l179
													}
													position++
													goto l178
												l179:
													position, tokenIndex = position179, tokenIndex179
												}
												add(ruleIntRangeValue, position175)
											}
											add(rulePegText, position174)
										}
										{
											add(ruleAction11, position)
										}
										goto l140
									l173:
										position, tokenIndex = position140, tokenIndex140
										{
											position182 := position
											{
												position183 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l181
												}
												position++
											l184:
												{
													position185, tokenIndex185 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l185
													}
													position++
													goto l184
												l185:
													position, tokenIndex = position185, tokenIndex185
												}
												add(ruleIntValue, position183)
											}
											add(rulePegText, position182)
										}
										{
											add(ruleAction12, position)
										}
										goto l140
									l181:
										position, tokenIndex = position140, tokenIndex140
										{
											switch buffer[position] {
											case '$':
												{
													position188 := position
													if buffer[position] != rune('$') {
														goto l72
													}
													position++
													{
														position189 := position
														if !_rules[ruleIdentifier]() {
															goto l72
														}
														add(rulePegText, position189)
													}
													add(ruleRefValue, position188)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position191 := position
													{
														position192 := position
														if buffer[position] != rune('@') {
															goto l72
														}
														position++
														if !_rules[ruleStringValue]() {
															goto l72
														}
														add(rulePegText, position192)
													}
													add(ruleAliasValue, position191)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position194 := position
													if buffer[position] != rune('{') {
														goto l72
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l72
													}
													{
														position195 := position
														if !_rules[ruleIdentifier]() {
															goto l72
														}
														add(rulePegText, position195)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l72
													}
													if buffer[position] != rune('}') {
														goto l72
													}
													position++
													add(ruleHoleValue, position194)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position197 := position
													if !_rules[ruleStringValue]() {
														goto l72
													}
													add(rulePegText, position197)
												}
												{
													add(ruleAction13, position)
												}
												break
											}
										}

									}
								l140:
									add(ruleValue, position139)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l72
								}
								add(ruleParam, position136)
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
			position202, tokenIndex202 := position, tokenIndex
			{
				position203 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l202
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l202
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l202
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l202
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l202
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l202
						}
						position++
						break
					}
				}

			l204:
				{
					position205, tokenIndex205 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l205
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l205
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l205
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l205
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l205
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l205
							}
							position++
							break
						}
					}

					goto l204
				l205:
					position, tokenIndex = position205, tokenIndex205
				}
				add(ruleIdentifier, position203)
			}
			return true
		l202:
			position, tokenIndex = position202, tokenIndex202
			return false
		},
		/* 9 Value <- <((<CidrValue> Action8) / (<IpValue> Action9) / (<CSVValue> Action10) / (<IntRangeValue> Action11) / (<IntValue> Action12) / ((&('$') (RefValue Action7)) | (&('@') (AliasValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action13))))> */
		nil,
		/* 10 StringValue <- <((&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position209, tokenIndex209 := position, tokenIndex
			{
				position210 := position
				{
					switch buffer[position] {
					case '/':
						if buffer[position] != rune('/') {
							goto l209
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l209
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l209
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l209
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l209
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l209
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l209
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l209
						}
						position++
						break
					}
				}

			l211:
				{
					position212, tokenIndex212 := position, tokenIndex
					{
						switch buffer[position] {
						case '/':
							if buffer[position] != rune('/') {
								goto l212
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l212
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l212
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l212
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l212
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l212
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l212
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l212
							}
							position++
							break
						}
					}

					goto l211
				l212:
					position, tokenIndex = position212, tokenIndex212
				}
				add(ruleStringValue, position210)
			}
			return true
		l209:
			position, tokenIndex = position209, tokenIndex209
			return false
		},
		/* 11 CSVValue <- <((StringValue WhiteSpacing ',' WhiteSpacing)+ StringValue)> */
		nil,
		/* 12 CidrValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+ '/' [0-9]+)> */
		nil,
		/* 13 IpValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+)> */
		nil,
		/* 14 IntValue <- <[0-9]+> */
		nil,
		/* 15 IntRangeValue <- <([0-9]+ '-' [0-9]+)> */
		nil,
		/* 16 RefValue <- <('$' <Identifier>)> */
		nil,
		/* 17 AliasValue <- <<('@' StringValue)>> */
		nil,
		/* 18 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 19 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)* Action14))> */
		nil,
		/* 20 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position225 := position
			l226:
				{
					position227, tokenIndex227 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l227
					}
					goto l226
				l227:
					position, tokenIndex = position227, tokenIndex227
				}
				add(ruleWhiteSpacing, position225)
			}
			return true
		},
		/* 21 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position228, tokenIndex228 := position, tokenIndex
			{
				position229 := position
				if !_rules[ruleWhitespace]() {
					goto l228
				}
			l230:
				{
					position231, tokenIndex231 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l231
					}
					goto l230
				l231:
					position, tokenIndex = position231, tokenIndex231
				}
				add(ruleMustWhiteSpacing, position229)
			}
			return true
		l228:
			position, tokenIndex = position228, tokenIndex228
			return false
		},
		/* 22 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position232, tokenIndex232 := position, tokenIndex
			{
				position233 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l232
				}
				if buffer[position] != rune('=') {
					goto l232
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l232
				}
				add(ruleEqual, position233)
			}
			return true
		l232:
			position, tokenIndex = position232, tokenIndex232
			return false
		},
		/* 23 BlankLine <- <(WhiteSpacing EndOfLine Action15)> */
		func() bool {
			position234, tokenIndex234 := position, tokenIndex
			{
				position235 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l234
				}
				if !_rules[ruleEndOfLine]() {
					goto l234
				}
				{
					add(ruleAction15, position)
				}
				add(ruleBlankLine, position235)
			}
			return true
		l234:
			position, tokenIndex = position234, tokenIndex234
			return false
		},
		/* 24 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position237, tokenIndex237 := position, tokenIndex
			{
				position238 := position
				{
					position239, tokenIndex239 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l240
					}
					position++
					goto l239
				l240:
					position, tokenIndex = position239, tokenIndex239
					if buffer[position] != rune('\t') {
						goto l237
					}
					position++
				}
			l239:
				add(ruleWhitespace, position238)
			}
			return true
		l237:
			position, tokenIndex = position237, tokenIndex237
			return false
		},
		/* 25 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position241, tokenIndex241 := position, tokenIndex
			{
				position242 := position
				{
					position243, tokenIndex243 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l244
					}
					position++
					if buffer[position] != rune('\n') {
						goto l244
					}
					position++
					goto l243
				l244:
					position, tokenIndex = position243, tokenIndex243
					if buffer[position] != rune('\n') {
						goto l245
					}
					position++
					goto l243
				l245:
					position, tokenIndex = position243, tokenIndex243
					if buffer[position] != rune('\r') {
						goto l241
					}
					position++
				}
			l243:
				add(ruleEndOfLine, position242)
			}
			return true
		l241:
			position, tokenIndex = position241, tokenIndex241
			return false
		},
		/* 26 EndOfFile <- <!.> */
		nil,
		nil,
		/* 29 Action0 <- <{ p.addDeclarationIdentifier(text) }> */
		nil,
		/* 30 Action1 <- <{ p.addAction(text) }> */
		nil,
		/* 31 Action2 <- <{ p.addEntity(text) }> */
		nil,
		/* 32 Action3 <- <{ p.LineDone() }> */
		nil,
		/* 33 Action4 <- <{ p.addParamKey(text) }> */
		nil,
		/* 34 Action5 <- <{  p.addParamHoleValue(text) }> */
		nil,
		/* 35 Action6 <- <{  p.addParamValue(text) }> */
		nil,
		/* 36 Action7 <- <{  p.addParamRefValue(text) }> */
		nil,
		/* 37 Action8 <- <{ p.addParamCidrValue(text) }> */
		nil,
		/* 38 Action9 <- <{ p.addParamIpValue(text) }> */
		nil,
		/* 39 Action10 <- <{p.addCsvValue(text)}> */
		nil,
		/* 40 Action11 <- <{ p.addParamValue(text) }> */
		nil,
		/* 41 Action12 <- <{ p.addParamIntValue(text) }> */
		nil,
		/* 42 Action13 <- <{ p.addParamValue(text) }> */
		nil,
		/* 43 Action14 <- <{ p.LineDone() }> */
		nil,
		/* 44 Action15 <- <{ p.LineDone() }> */
		nil,
	}
	p.rules = _rules
}
