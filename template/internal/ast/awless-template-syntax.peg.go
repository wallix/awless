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
	ruleValueExpr
	ruleCmdExpr
	ruleParams
	ruleParam
	ruleIdentifier
	ruleNoRefValue
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
	ruleAction18
	ruleAction19
)

var rul3s = [...]string{
	"Unknown",
	"Script",
	"Statement",
	"Action",
	"Entity",
	"Declaration",
	"ValueExpr",
	"CmdExpr",
	"Params",
	"Param",
	"Identifier",
	"NoRefValue",
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
	"Action18",
	"Action19",
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
	rules  [56]func() bool
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
			p.addValue()
		case ruleAction2:
			p.LineDone()
		case ruleAction3:
			p.addAction(text)
		case ruleAction4:
			p.addEntity(text)
		case ruleAction5:
			p.LineDone()
		case ruleAction6:
			p.addParamKey(text)
		case ruleAction7:
			p.addParamHoleValue(text)
		case ruleAction8:
			p.addAliasParam(text)
		case ruleAction9:
			p.addParamValue(text)
		case ruleAction10:
			p.addParamValue(text)
		case ruleAction11:
			p.addParamIntValue(text)
		case ruleAction12:
			p.addParamValue(text)
		case ruleAction13:
			p.addParamRefValue(text)
		case ruleAction14:
			p.addParamCidrValue(text)
		case ruleAction15:
			p.addParamIpValue(text)
		case ruleAction16:
			p.addCsvValue(text)
		case ruleAction17:
			p.addParamValue(text)
		case ruleAction18:
			p.LineDone()
		case ruleAction19:
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
						if !_rules[ruleCmdExpr]() {
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
							{
								position13, tokenIndex13 := position, tokenIndex
								if !_rules[ruleCmdExpr]() {
									goto l14
								}
								goto l13
							l14:
								position, tokenIndex = position13, tokenIndex13
								{
									position15 := position
									{
										add(ruleAction1, position)
									}
									if !_rules[ruleNoRefValue]() {
										goto l9
									}
									{
										add(ruleAction2, position)
									}
									add(ruleValueExpr, position15)
								}
							}
						l13:
							add(ruleDeclaration, position10)
						}
						goto l7
					l9:
						position, tokenIndex = position7, tokenIndex7
						{
							position18 := position
							{
								position19, tokenIndex19 := position, tokenIndex
								if buffer[position] != rune('#') {
									goto l20
								}
								position++
							l21:
								{
									position22, tokenIndex22 := position, tokenIndex
									{
										position23, tokenIndex23 := position, tokenIndex
										if !_rules[ruleEndOfLine]() {
											goto l23
										}
										goto l22
									l23:
										position, tokenIndex = position23, tokenIndex23
									}
									if !matchDot() {
										goto l22
									}
									goto l21
								l22:
									position, tokenIndex = position22, tokenIndex22
								}
								goto l19
							l20:
								position, tokenIndex = position19, tokenIndex19
								if buffer[position] != rune('/') {
									goto l0
								}
								position++
								if buffer[position] != rune('/') {
									goto l0
								}
								position++
							l24:
								{
									position25, tokenIndex25 := position, tokenIndex
									{
										position26, tokenIndex26 := position, tokenIndex
										if !_rules[ruleEndOfLine]() {
											goto l26
										}
										goto l25
									l26:
										position, tokenIndex = position26, tokenIndex26
									}
									if !matchDot() {
										goto l25
									}
									goto l24
								l25:
									position, tokenIndex = position25, tokenIndex25
								}
								{
									add(ruleAction18, position)
								}
							}
						l19:
							add(ruleComment, position18)
						}
					}
				l7:
					if !_rules[ruleWhiteSpacing]() {
						goto l0
					}
				l28:
					{
						position29, tokenIndex29 := position, tokenIndex
						if !_rules[ruleEndOfLine]() {
							goto l29
						}
						goto l28
					l29:
						position, tokenIndex = position29, tokenIndex29
					}
					add(ruleStatement, position6)
				}
			l30:
				{
					position31, tokenIndex31 := position, tokenIndex
					if !_rules[ruleBlankLine]() {
						goto l31
					}
					goto l30
				l31:
					position, tokenIndex = position31, tokenIndex31
				}
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
				l32:
					{
						position33, tokenIndex33 := position, tokenIndex
						if !_rules[ruleBlankLine]() {
							goto l33
						}
						goto l32
					l33:
						position, tokenIndex = position33, tokenIndex33
					}
					{
						position34 := position
						if !_rules[ruleWhiteSpacing]() {
							goto l3
						}
						{
							position35, tokenIndex35 := position, tokenIndex
							if !_rules[ruleCmdExpr]() {
								goto l36
							}
							goto l35
						l36:
							position, tokenIndex = position35, tokenIndex35
							{
								position38 := position
								{
									position39 := position
									if !_rules[ruleIdentifier]() {
										goto l37
									}
									add(rulePegText, position39)
								}
								{
									add(ruleAction0, position)
								}
								if !_rules[ruleEqual]() {
									goto l37
								}
								{
									position41, tokenIndex41 := position, tokenIndex
									if !_rules[ruleCmdExpr]() {
										goto l42
									}
									goto l41
								l42:
									position, tokenIndex = position41, tokenIndex41
									{
										position43 := position
										{
											add(ruleAction1, position)
										}
										if !_rules[ruleNoRefValue]() {
											goto l37
										}
										{
											add(ruleAction2, position)
										}
										add(ruleValueExpr, position43)
									}
								}
							l41:
								add(ruleDeclaration, position38)
							}
							goto l35
						l37:
							position, tokenIndex = position35, tokenIndex35
							{
								position46 := position
								{
									position47, tokenIndex47 := position, tokenIndex
									if buffer[position] != rune('#') {
										goto l48
									}
									position++
								l49:
									{
										position50, tokenIndex50 := position, tokenIndex
										{
											position51, tokenIndex51 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l51
											}
											goto l50
										l51:
											position, tokenIndex = position51, tokenIndex51
										}
										if !matchDot() {
											goto l50
										}
										goto l49
									l50:
										position, tokenIndex = position50, tokenIndex50
									}
									goto l47
								l48:
									position, tokenIndex = position47, tokenIndex47
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
								l52:
									{
										position53, tokenIndex53 := position, tokenIndex
										{
											position54, tokenIndex54 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l54
											}
											goto l53
										l54:
											position, tokenIndex = position54, tokenIndex54
										}
										if !matchDot() {
											goto l53
										}
										goto l52
									l53:
										position, tokenIndex = position53, tokenIndex53
									}
									{
										add(ruleAction18, position)
									}
								}
							l47:
								add(ruleComment, position46)
							}
						}
					l35:
						if !_rules[ruleWhiteSpacing]() {
							goto l3
						}
					l56:
						{
							position57, tokenIndex57 := position, tokenIndex
							if !_rules[ruleEndOfLine]() {
								goto l57
							}
							goto l56
						l57:
							position, tokenIndex = position57, tokenIndex57
						}
						add(ruleStatement, position34)
					}
				l58:
					{
						position59, tokenIndex59 := position, tokenIndex
						if !_rules[ruleBlankLine]() {
							goto l59
						}
						goto l58
					l59:
						position, tokenIndex = position59, tokenIndex59
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				if !_rules[ruleWhiteSpacing]() {
					goto l0
				}
				{
					position60 := position
					{
						position61, tokenIndex61 := position, tokenIndex
						if !matchDot() {
							goto l61
						}
						goto l0
					l61:
						position, tokenIndex = position61, tokenIndex61
					}
					add(ruleEndOfFile, position60)
				}
				add(ruleScript, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Statement <- <(WhiteSpacing (CmdExpr / Declaration / Comment) WhiteSpacing EndOfLine*)> */
		nil,
		/* 2 Action <- <[a-z]+> */
		nil,
		/* 3 Entity <- <([a-z] / [0-9])+> */
		nil,
		/* 4 Declaration <- <(<Identifier> Action0 Equal (CmdExpr / ValueExpr))> */
		nil,
		/* 5 ValueExpr <- <(Action1 NoRefValue Action2)> */
		nil,
		/* 6 CmdExpr <- <(<Action> Action3 MustWhiteSpacing <Entity> Action4 (MustWhiteSpacing Params)? Action5)> */
		func() bool {
			position67, tokenIndex67 := position, tokenIndex
			{
				position68 := position
				{
					position69 := position
					{
						position70 := position
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l67
						}
						position++
					l71:
						{
							position72, tokenIndex72 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l72
							}
							position++
							goto l71
						l72:
							position, tokenIndex = position72, tokenIndex72
						}
						add(ruleAction, position70)
					}
					add(rulePegText, position69)
				}
				{
					add(ruleAction3, position)
				}
				if !_rules[ruleMustWhiteSpacing]() {
					goto l67
				}
				{
					position74 := position
					{
						position75 := position
						{
							position78, tokenIndex78 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l79
							}
							position++
							goto l78
						l79:
							position, tokenIndex = position78, tokenIndex78
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l67
							}
							position++
						}
					l78:
					l76:
						{
							position77, tokenIndex77 := position, tokenIndex
							{
								position80, tokenIndex80 := position, tokenIndex
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l81
								}
								position++
								goto l80
							l81:
								position, tokenIndex = position80, tokenIndex80
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l77
								}
								position++
							}
						l80:
							goto l76
						l77:
							position, tokenIndex = position77, tokenIndex77
						}
						add(ruleEntity, position75)
					}
					add(rulePegText, position74)
				}
				{
					add(ruleAction4, position)
				}
				{
					position83, tokenIndex83 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l83
					}
					{
						position85 := position
						{
							position88 := position
							{
								position89 := position
								if !_rules[ruleIdentifier]() {
									goto l83
								}
								add(rulePegText, position89)
							}
							{
								add(ruleAction6, position)
							}
							if !_rules[ruleEqual]() {
								goto l83
							}
							{
								position91 := position
								{
									position92, tokenIndex92 := position, tokenIndex
									{
										position94 := position
										if buffer[position] != rune('$') {
											goto l93
										}
										position++
										{
											position95 := position
											if !_rules[ruleIdentifier]() {
												goto l93
											}
											add(rulePegText, position95)
										}
										add(ruleRefValue, position94)
									}
									{
										add(ruleAction13, position)
									}
									goto l92
								l93:
									position, tokenIndex = position92, tokenIndex92
									if !_rules[ruleNoRefValue]() {
										goto l83
									}
								}
							l92:
								add(ruleValue, position91)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l83
							}
							add(ruleParam, position88)
						}
					l86:
						{
							position87, tokenIndex87 := position, tokenIndex
							{
								position97 := position
								{
									position98 := position
									if !_rules[ruleIdentifier]() {
										goto l87
									}
									add(rulePegText, position98)
								}
								{
									add(ruleAction6, position)
								}
								if !_rules[ruleEqual]() {
									goto l87
								}
								{
									position100 := position
									{
										position101, tokenIndex101 := position, tokenIndex
										{
											position103 := position
											if buffer[position] != rune('$') {
												goto l102
											}
											position++
											{
												position104 := position
												if !_rules[ruleIdentifier]() {
													goto l102
												}
												add(rulePegText, position104)
											}
											add(ruleRefValue, position103)
										}
										{
											add(ruleAction13, position)
										}
										goto l101
									l102:
										position, tokenIndex = position101, tokenIndex101
										if !_rules[ruleNoRefValue]() {
											goto l87
										}
									}
								l101:
									add(ruleValue, position100)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l87
								}
								add(ruleParam, position97)
							}
							goto l86
						l87:
							position, tokenIndex = position87, tokenIndex87
						}
						add(ruleParams, position85)
					}
					goto l84
				l83:
					position, tokenIndex = position83, tokenIndex83
				}
			l84:
				{
					add(ruleAction5, position)
				}
				add(ruleCmdExpr, position68)
			}
			return true
		l67:
			position, tokenIndex = position67, tokenIndex67
			return false
		},
		/* 7 Params <- <Param+> */
		nil,
		/* 8 Param <- <(<Identifier> Action6 Equal Value WhiteSpacing)> */
		nil,
		/* 9 Identifier <- <((&('.') '.') | (&('_') '_') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position109, tokenIndex109 := position, tokenIndex
			{
				position110 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l109
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l109
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l109
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l109
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l109
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l109
						}
						position++
						break
					}
				}

			l111:
				{
					position112, tokenIndex112 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l112
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l112
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l112
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l112
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l112
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l112
							}
							position++
							break
						}
					}

					goto l111
				l112:
					position, tokenIndex = position112, tokenIndex112
				}
				add(ruleIdentifier, position110)
			}
			return true
		l109:
			position, tokenIndex = position109, tokenIndex109
			return false
		},
		/* 10 NoRefValue <- <((AliasValue Action8) / (DoubleQuote CustomTypedValue DoubleQuote) / (SingleQuote CustomTypedValue SingleQuote) / CustomTypedValue / (<IntValue> Action11) / ((&('\'') (SingleQuote <SingleQuotedValue> Action10 SingleQuote)) | (&('"') (DoubleQuote <DoubleQuotedValue> Action9 DoubleQuote)) | (&('{') (HoleValue Action7)) | (&('+' | '-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | ';' | '<' | '>' | '@' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | '~') (<StringValue> Action12))))> */
		func() bool {
			position115, tokenIndex115 := position, tokenIndex
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
						add(ruleAction8, position)
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
						add(ruleAction11, position)
					}
					goto l117
				l130:
					position, tokenIndex = position117, tokenIndex117
					{
						switch buffer[position] {
						case '\'':
							if !_rules[ruleSingleQuote]() {
								goto l115
							}
							{
								position137 := position
								if !_rules[ruleSingleQuotedValue]() {
									goto l115
								}
								add(rulePegText, position137)
							}
							{
								add(ruleAction10, position)
							}
							if !_rules[ruleSingleQuote]() {
								goto l115
							}
							break
						case '"':
							if !_rules[ruleDoubleQuote]() {
								goto l115
							}
							{
								position139 := position
								if !_rules[ruleDoubleQuotedValue]() {
									goto l115
								}
								add(rulePegText, position139)
							}
							{
								add(ruleAction9, position)
							}
							if !_rules[ruleDoubleQuote]() {
								goto l115
							}
							break
						case '{':
							{
								position141 := position
								if buffer[position] != rune('{') {
									goto l115
								}
								position++
								if !_rules[ruleWhiteSpacing]() {
									goto l115
								}
								{
									position142 := position
									if !_rules[ruleIdentifier]() {
										goto l115
									}
									add(rulePegText, position142)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l115
								}
								if buffer[position] != rune('}') {
									goto l115
								}
								position++
								add(ruleHoleValue, position141)
							}
							{
								add(ruleAction7, position)
							}
							break
						default:
							{
								position144 := position
								if !_rules[ruleStringValue]() {
									goto l115
								}
								add(rulePegText, position144)
							}
							{
								add(ruleAction12, position)
							}
							break
						}
					}

				}
			l117:
				add(ruleNoRefValue, position116)
			}
			return true
		l115:
			position, tokenIndex = position115, tokenIndex115
			return false
		},
		/* 11 Value <- <((RefValue Action13) / NoRefValue)> */
		nil,
		/* 12 CustomTypedValue <- <((<CidrValue> Action14) / (<IpValue> Action15) / (<CSVValue> Action16) / (<IntRangeValue> Action17))> */
		func() bool {
			position147, tokenIndex147 := position, tokenIndex
			{
				position148 := position
				{
					position149, tokenIndex149 := position, tokenIndex
					{
						position151 := position
						{
							position152 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l150
							}
							position++
						l153:
							{
								position154, tokenIndex154 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l154
								}
								position++
								goto l153
							l154:
								position, tokenIndex = position154, tokenIndex154
							}
							if !matchDot() {
								goto l150
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l150
							}
							position++
						l155:
							{
								position156, tokenIndex156 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l156
								}
								position++
								goto l155
							l156:
								position, tokenIndex = position156, tokenIndex156
							}
							if !matchDot() {
								goto l150
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l150
							}
							position++
						l157:
							{
								position158, tokenIndex158 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l158
								}
								position++
								goto l157
							l158:
								position, tokenIndex = position158, tokenIndex158
							}
							if !matchDot() {
								goto l150
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l150
							}
							position++
						l159:
							{
								position160, tokenIndex160 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l160
								}
								position++
								goto l159
							l160:
								position, tokenIndex = position160, tokenIndex160
							}
							if buffer[position] != rune('/') {
								goto l150
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l150
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
							add(ruleCidrValue, position152)
						}
						add(rulePegText, position151)
					}
					{
						add(ruleAction14, position)
					}
					goto l149
				l150:
					position, tokenIndex = position149, tokenIndex149
					{
						position165 := position
						{
							position166 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l164
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
								goto l164
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l164
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
								goto l164
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l164
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
							if !matchDot() {
								goto l164
							}
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l164
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
							add(ruleIpValue, position166)
						}
						add(rulePegText, position165)
					}
					{
						add(ruleAction15, position)
					}
					goto l149
				l164:
					position, tokenIndex = position149, tokenIndex149
					{
						position177 := position
						{
							position178 := position
							if !_rules[ruleStringValue]() {
								goto l176
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l176
							}
							if buffer[position] != rune(',') {
								goto l176
							}
							position++
							if !_rules[ruleWhiteSpacing]() {
								goto l176
							}
						l179:
							{
								position180, tokenIndex180 := position, tokenIndex
								if !_rules[ruleStringValue]() {
									goto l180
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l180
								}
								if buffer[position] != rune(',') {
									goto l180
								}
								position++
								if !_rules[ruleWhiteSpacing]() {
									goto l180
								}
								goto l179
							l180:
								position, tokenIndex = position180, tokenIndex180
							}
							if !_rules[ruleStringValue]() {
								goto l176
							}
							add(ruleCSVValue, position178)
						}
						add(rulePegText, position177)
					}
					{
						add(ruleAction16, position)
					}
					goto l149
				l176:
					position, tokenIndex = position149, tokenIndex149
					{
						position182 := position
						{
							position183 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l147
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
							if buffer[position] != rune('-') {
								goto l147
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l147
							}
							position++
						l186:
							{
								position187, tokenIndex187 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l187
								}
								position++
								goto l186
							l187:
								position, tokenIndex = position187, tokenIndex187
							}
							add(ruleIntRangeValue, position183)
						}
						add(rulePegText, position182)
					}
					{
						add(ruleAction17, position)
					}
				}
			l149:
				add(ruleCustomTypedValue, position148)
			}
			return true
		l147:
			position, tokenIndex = position147, tokenIndex147
			return false
		},
		/* 13 StringValue <- <((&('>') '>') | (&('<') '<') | (&('@') '@') | (&('~') '~') | (&(';') ';') | (&('+') '+') | (&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position189, tokenIndex189 := position, tokenIndex
			{
				position190 := position
				{
					switch buffer[position] {
					case '>':
						if buffer[position] != rune('>') {
							goto l189
						}
						position++
						break
					case '<':
						if buffer[position] != rune('<') {
							goto l189
						}
						position++
						break
					case '@':
						if buffer[position] != rune('@') {
							goto l189
						}
						position++
						break
					case '~':
						if buffer[position] != rune('~') {
							goto l189
						}
						position++
						break
					case ';':
						if buffer[position] != rune(';') {
							goto l189
						}
						position++
						break
					case '+':
						if buffer[position] != rune('+') {
							goto l189
						}
						position++
						break
					case '/':
						if buffer[position] != rune('/') {
							goto l189
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l189
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l189
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l189
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l189
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l189
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l189
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l189
						}
						position++
						break
					}
				}

			l191:
				{
					position192, tokenIndex192 := position, tokenIndex
					{
						switch buffer[position] {
						case '>':
							if buffer[position] != rune('>') {
								goto l192
							}
							position++
							break
						case '<':
							if buffer[position] != rune('<') {
								goto l192
							}
							position++
							break
						case '@':
							if buffer[position] != rune('@') {
								goto l192
							}
							position++
							break
						case '~':
							if buffer[position] != rune('~') {
								goto l192
							}
							position++
							break
						case ';':
							if buffer[position] != rune(';') {
								goto l192
							}
							position++
							break
						case '+':
							if buffer[position] != rune('+') {
								goto l192
							}
							position++
							break
						case '/':
							if buffer[position] != rune('/') {
								goto l192
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l192
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l192
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l192
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l192
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l192
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l192
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l192
							}
							position++
							break
						}
					}

					goto l191
				l192:
					position, tokenIndex = position192, tokenIndex192
				}
				add(ruleStringValue, position190)
			}
			return true
		l189:
			position, tokenIndex = position189, tokenIndex189
			return false
		},
		/* 14 DoubleQuotedValue <- <(!'"' .)*> */
		func() bool {
			{
				position196 := position
			l197:
				{
					position198, tokenIndex198 := position, tokenIndex
					{
						position199, tokenIndex199 := position, tokenIndex
						if buffer[position] != rune('"') {
							goto l199
						}
						position++
						goto l198
					l199:
						position, tokenIndex = position199, tokenIndex199
					}
					if !matchDot() {
						goto l198
					}
					goto l197
				l198:
					position, tokenIndex = position198, tokenIndex198
				}
				add(ruleDoubleQuotedValue, position196)
			}
			return true
		},
		/* 15 SingleQuotedValue <- <(!'\'' .)*> */
		func() bool {
			{
				position201 := position
			l202:
				{
					position203, tokenIndex203 := position, tokenIndex
					{
						position204, tokenIndex204 := position, tokenIndex
						if buffer[position] != rune('\'') {
							goto l204
						}
						position++
						goto l203
					l204:
						position, tokenIndex = position204, tokenIndex204
					}
					if !matchDot() {
						goto l203
					}
					goto l202
				l203:
					position, tokenIndex = position203, tokenIndex203
				}
				add(ruleSingleQuotedValue, position201)
			}
			return true
		},
		/* 16 CSVValue <- <((StringValue WhiteSpacing ',' WhiteSpacing)+ StringValue)> */
		nil,
		/* 17 CidrValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+ '/' [0-9]+)> */
		nil,
		/* 18 IpValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+)> */
		nil,
		/* 19 IntValue <- <[0-9]+> */
		nil,
		/* 20 IntRangeValue <- <([0-9]+ '-' [0-9]+)> */
		nil,
		/* 21 RefValue <- <('$' <Identifier>)> */
		nil,
		/* 22 AliasValue <- <(('@' <StringValue>) / ('@' DoubleQuote <DoubleQuotedValue> DoubleQuote) / ('@' SingleQuote <SingleQuotedValue> SingleQuote))> */
		nil,
		/* 23 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 24 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)* Action18))> */
		nil,
		/* 25 SingleQuote <- <'\''> */
		func() bool {
			position214, tokenIndex214 := position, tokenIndex
			{
				position215 := position
				if buffer[position] != rune('\'') {
					goto l214
				}
				position++
				add(ruleSingleQuote, position215)
			}
			return true
		l214:
			position, tokenIndex = position214, tokenIndex214
			return false
		},
		/* 26 DoubleQuote <- <'"'> */
		func() bool {
			position216, tokenIndex216 := position, tokenIndex
			{
				position217 := position
				if buffer[position] != rune('"') {
					goto l216
				}
				position++
				add(ruleDoubleQuote, position217)
			}
			return true
		l216:
			position, tokenIndex = position216, tokenIndex216
			return false
		},
		/* 27 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position219 := position
			l220:
				{
					position221, tokenIndex221 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l221
					}
					goto l220
				l221:
					position, tokenIndex = position221, tokenIndex221
				}
				add(ruleWhiteSpacing, position219)
			}
			return true
		},
		/* 28 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position222, tokenIndex222 := position, tokenIndex
			{
				position223 := position
				if !_rules[ruleWhitespace]() {
					goto l222
				}
			l224:
				{
					position225, tokenIndex225 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l225
					}
					goto l224
				l225:
					position, tokenIndex = position225, tokenIndex225
				}
				add(ruleMustWhiteSpacing, position223)
			}
			return true
		l222:
			position, tokenIndex = position222, tokenIndex222
			return false
		},
		/* 29 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position226, tokenIndex226 := position, tokenIndex
			{
				position227 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l226
				}
				if buffer[position] != rune('=') {
					goto l226
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l226
				}
				add(ruleEqual, position227)
			}
			return true
		l226:
			position, tokenIndex = position226, tokenIndex226
			return false
		},
		/* 30 BlankLine <- <(WhiteSpacing EndOfLine Action19)> */
		func() bool {
			position228, tokenIndex228 := position, tokenIndex
			{
				position229 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l228
				}
				if !_rules[ruleEndOfLine]() {
					goto l228
				}
				{
					add(ruleAction19, position)
				}
				add(ruleBlankLine, position229)
			}
			return true
		l228:
			position, tokenIndex = position228, tokenIndex228
			return false
		},
		/* 31 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position231, tokenIndex231 := position, tokenIndex
			{
				position232 := position
				{
					position233, tokenIndex233 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l234
					}
					position++
					goto l233
				l234:
					position, tokenIndex = position233, tokenIndex233
					if buffer[position] != rune('\t') {
						goto l231
					}
					position++
				}
			l233:
				add(ruleWhitespace, position232)
			}
			return true
		l231:
			position, tokenIndex = position231, tokenIndex231
			return false
		},
		/* 32 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position235, tokenIndex235 := position, tokenIndex
			{
				position236 := position
				{
					position237, tokenIndex237 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l238
					}
					position++
					if buffer[position] != rune('\n') {
						goto l238
					}
					position++
					goto l237
				l238:
					position, tokenIndex = position237, tokenIndex237
					if buffer[position] != rune('\n') {
						goto l239
					}
					position++
					goto l237
				l239:
					position, tokenIndex = position237, tokenIndex237
					if buffer[position] != rune('\r') {
						goto l235
					}
					position++
				}
			l237:
				add(ruleEndOfLine, position236)
			}
			return true
		l235:
			position, tokenIndex = position235, tokenIndex235
			return false
		},
		/* 33 EndOfFile <- <!.> */
		nil,
		nil,
		/* 36 Action0 <- <{ p.addDeclarationIdentifier(text) }> */
		nil,
		/* 37 Action1 <- <{ p.addValue() }> */
		nil,
		/* 38 Action2 <- <{ p.LineDone() }> */
		nil,
		/* 39 Action3 <- <{ p.addAction(text) }> */
		nil,
		/* 40 Action4 <- <{ p.addEntity(text) }> */
		nil,
		/* 41 Action5 <- <{ p.LineDone() }> */
		nil,
		/* 42 Action6 <- <{ p.addParamKey(text) }> */
		nil,
		/* 43 Action7 <- <{  p.addParamHoleValue(text) }> */
		nil,
		/* 44 Action8 <- <{  p.addAliasParam(text) }> */
		nil,
		/* 45 Action9 <- <{ p.addParamValue(text) }> */
		nil,
		/* 46 Action10 <- <{ p.addParamValue(text) }> */
		nil,
		/* 47 Action11 <- <{ p.addParamIntValue(text) }> */
		nil,
		/* 48 Action12 <- <{ p.addParamValue(text) }> */
		nil,
		/* 49 Action13 <- <{  p.addParamRefValue(text) }> */
		nil,
		/* 50 Action14 <- <{ p.addParamCidrValue(text) }> */
		nil,
		/* 51 Action15 <- <{ p.addParamIpValue(text) }> */
		nil,
		/* 52 Action16 <- <{p.addCsvValue(text)}> */
		nil,
		/* 53 Action17 <- <{ p.addParamValue(text) }> */
		nil,
		/* 54 Action18 <- <{ p.LineDone() }> */
		nil,
		/* 55 Action19 <- <{ p.LineDone() }> */
		nil,
	}
	p.rules = _rules
}
