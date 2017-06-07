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
	ruleFloatValue
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
	ruleAction20
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
	"FloatValue",
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
	"Action20",
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
	rules  [58]func() bool
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
			p.addParamFloatValue(text)
		case ruleAction12:
			p.addParamIntValue(text)
		case ruleAction13:
			p.addParamValue(text)
		case ruleAction14:
			p.addParamRefValue(text)
		case ruleAction15:
			p.addParamCidrValue(text)
		case ruleAction16:
			p.addParamIpValue(text)
		case ruleAction17:
			p.addCsvValue(text)
		case ruleAction18:
			p.addParamValue(text)
		case ruleAction19:
			p.LineDone()
		case ruleAction20:
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
									add(ruleAction19, position)
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
										add(ruleAction19, position)
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
										add(ruleAction14, position)
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
											add(ruleAction14, position)
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
		/* 10 NoRefValue <- <((AliasValue Action8) / (DoubleQuote CustomTypedValue DoubleQuote) / (SingleQuote CustomTypedValue SingleQuote) / CustomTypedValue / (<FloatValue> Action11) / (<IntValue> Action12) / ((&('\'') (SingleQuote <SingleQuotedValue> Action10 SingleQuote)) | (&('"') (DoubleQuote <DoubleQuotedValue> Action9 DoubleQuote)) | (&('{') (HoleValue Action7)) | (&('*' | '+' | '-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | ';' | '<' | '>' | '@' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | '~') (<StringValue> Action13))))> */
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
							if buffer[position] != rune('.') {
								goto l130
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
							add(ruleFloatValue, position132)
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
						position139 := position
						{
							position140 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l138
							}
							position++
						l141:
							{
								position142, tokenIndex142 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l142
								}
								position++
								goto l141
							l142:
								position, tokenIndex = position142, tokenIndex142
							}
							add(ruleIntValue, position140)
						}
						add(rulePegText, position139)
					}
					{
						add(ruleAction12, position)
					}
					goto l117
				l138:
					position, tokenIndex = position117, tokenIndex117
					{
						switch buffer[position] {
						case '\'':
							if !_rules[ruleSingleQuote]() {
								goto l115
							}
							{
								position145 := position
								if !_rules[ruleSingleQuotedValue]() {
									goto l115
								}
								add(rulePegText, position145)
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
								position147 := position
								if !_rules[ruleDoubleQuotedValue]() {
									goto l115
								}
								add(rulePegText, position147)
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
								position149 := position
								if buffer[position] != rune('{') {
									goto l115
								}
								position++
								if !_rules[ruleWhiteSpacing]() {
									goto l115
								}
								{
									position150 := position
									if !_rules[ruleIdentifier]() {
										goto l115
									}
									add(rulePegText, position150)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l115
								}
								if buffer[position] != rune('}') {
									goto l115
								}
								position++
								add(ruleHoleValue, position149)
							}
							{
								add(ruleAction7, position)
							}
							break
						default:
							{
								position152 := position
								if !_rules[ruleStringValue]() {
									goto l115
								}
								add(rulePegText, position152)
							}
							{
								add(ruleAction13, position)
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
		/* 11 Value <- <((RefValue Action14) / NoRefValue)> */
		nil,
		/* 12 CustomTypedValue <- <((<CidrValue> Action15) / (<IpValue> Action16) / (<CSVValue> Action17) / (<IntRangeValue> Action18))> */
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
							if buffer[position] != rune('.') {
								goto l158
							}
							position++
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
							if buffer[position] != rune('.') {
								goto l158
							}
							position++
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
							if buffer[position] != rune('.') {
								goto l158
							}
							position++
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
						add(ruleAction15, position)
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
							if buffer[position] != rune('.') {
								goto l172
							}
							position++
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
							if buffer[position] != rune('.') {
								goto l172
							}
							position++
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
							if buffer[position] != rune('.') {
								goto l172
							}
							position++
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
						add(ruleAction16, position)
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
						add(ruleAction17, position)
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
						add(ruleAction18, position)
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
		/* 13 StringValue <- <((&('*') '*') | (&('>') '>') | (&('<') '<') | (&('@') '@') | (&('~') '~') | (&(';') ';') | (&('+') '+') | (&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position197, tokenIndex197 := position, tokenIndex
			{
				position198 := position
				{
					switch buffer[position] {
					case '*':
						if buffer[position] != rune('*') {
							goto l197
						}
						position++
						break
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
						case '*':
							if buffer[position] != rune('*') {
								goto l200
							}
							position++
							break
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
		/* 14 DoubleQuotedValue <- <(!'"' .)*> */
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
		/* 15 SingleQuotedValue <- <(!'\'' .)*> */
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
		/* 16 CSVValue <- <((StringValue WhiteSpacing ',' WhiteSpacing)+ StringValue)> */
		nil,
		/* 17 CidrValue <- <([0-9]+ '.' [0-9]+ '.' [0-9]+ '.' [0-9]+ '/' [0-9]+)> */
		nil,
		/* 18 IpValue <- <([0-9]+ '.' [0-9]+ '.' [0-9]+ '.' [0-9]+)> */
		nil,
		/* 19 IntValue <- <[0-9]+> */
		nil,
		/* 20 FloatValue <- <([0-9]+ '.' [0-9]*)> */
		nil,
		/* 21 IntRangeValue <- <([0-9]+ '-' [0-9]+)> */
		nil,
		/* 22 RefValue <- <('$' <Identifier>)> */
		nil,
		/* 23 AliasValue <- <(('@' <StringValue>) / ('@' DoubleQuote <DoubleQuotedValue> DoubleQuote) / ('@' SingleQuote <SingleQuotedValue> SingleQuote))> */
		nil,
		/* 24 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 25 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)* Action19))> */
		nil,
		/* 26 SingleQuote <- <'\''> */
		func() bool {
			position223, tokenIndex223 := position, tokenIndex
			{
				position224 := position
				if buffer[position] != rune('\'') {
					goto l223
				}
				position++
				add(ruleSingleQuote, position224)
			}
			return true
		l223:
			position, tokenIndex = position223, tokenIndex223
			return false
		},
		/* 27 DoubleQuote <- <'"'> */
		func() bool {
			position225, tokenIndex225 := position, tokenIndex
			{
				position226 := position
				if buffer[position] != rune('"') {
					goto l225
				}
				position++
				add(ruleDoubleQuote, position226)
			}
			return true
		l225:
			position, tokenIndex = position225, tokenIndex225
			return false
		},
		/* 28 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position228 := position
			l229:
				{
					position230, tokenIndex230 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l230
					}
					goto l229
				l230:
					position, tokenIndex = position230, tokenIndex230
				}
				add(ruleWhiteSpacing, position228)
			}
			return true
		},
		/* 29 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position231, tokenIndex231 := position, tokenIndex
			{
				position232 := position
				if !_rules[ruleWhitespace]() {
					goto l231
				}
			l233:
				{
					position234, tokenIndex234 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l234
					}
					goto l233
				l234:
					position, tokenIndex = position234, tokenIndex234
				}
				add(ruleMustWhiteSpacing, position232)
			}
			return true
		l231:
			position, tokenIndex = position231, tokenIndex231
			return false
		},
		/* 30 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position235, tokenIndex235 := position, tokenIndex
			{
				position236 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l235
				}
				if buffer[position] != rune('=') {
					goto l235
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l235
				}
				add(ruleEqual, position236)
			}
			return true
		l235:
			position, tokenIndex = position235, tokenIndex235
			return false
		},
		/* 31 BlankLine <- <(WhiteSpacing EndOfLine Action20)> */
		func() bool {
			position237, tokenIndex237 := position, tokenIndex
			{
				position238 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l237
				}
				if !_rules[ruleEndOfLine]() {
					goto l237
				}
				{
					add(ruleAction20, position)
				}
				add(ruleBlankLine, position238)
			}
			return true
		l237:
			position, tokenIndex = position237, tokenIndex237
			return false
		},
		/* 32 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position240, tokenIndex240 := position, tokenIndex
			{
				position241 := position
				{
					position242, tokenIndex242 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex = position242, tokenIndex242
					if buffer[position] != rune('\t') {
						goto l240
					}
					position++
				}
			l242:
				add(ruleWhitespace, position241)
			}
			return true
		l240:
			position, tokenIndex = position240, tokenIndex240
			return false
		},
		/* 33 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position244, tokenIndex244 := position, tokenIndex
			{
				position245 := position
				{
					position246, tokenIndex246 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l247
					}
					position++
					if buffer[position] != rune('\n') {
						goto l247
					}
					position++
					goto l246
				l247:
					position, tokenIndex = position246, tokenIndex246
					if buffer[position] != rune('\n') {
						goto l248
					}
					position++
					goto l246
				l248:
					position, tokenIndex = position246, tokenIndex246
					if buffer[position] != rune('\r') {
						goto l244
					}
					position++
				}
			l246:
				add(ruleEndOfLine, position245)
			}
			return true
		l244:
			position, tokenIndex = position244, tokenIndex244
			return false
		},
		/* 34 EndOfFile <- <!.> */
		nil,
		nil,
		/* 37 Action0 <- <{ p.addDeclarationIdentifier(text) }> */
		nil,
		/* 38 Action1 <- <{ p.addValue() }> */
		nil,
		/* 39 Action2 <- <{ p.LineDone() }> */
		nil,
		/* 40 Action3 <- <{ p.addAction(text) }> */
		nil,
		/* 41 Action4 <- <{ p.addEntity(text) }> */
		nil,
		/* 42 Action5 <- <{ p.LineDone() }> */
		nil,
		/* 43 Action6 <- <{ p.addParamKey(text) }> */
		nil,
		/* 44 Action7 <- <{  p.addParamHoleValue(text) }> */
		nil,
		/* 45 Action8 <- <{  p.addAliasParam(text) }> */
		nil,
		/* 46 Action9 <- <{ p.addParamValue(text) }> */
		nil,
		/* 47 Action10 <- <{ p.addParamValue(text) }> */
		nil,
		/* 48 Action11 <- <{ p.addParamFloatValue(text) }> */
		nil,
		/* 49 Action12 <- <{ p.addParamIntValue(text) }> */
		nil,
		/* 50 Action13 <- <{ p.addParamValue(text) }> */
		nil,
		/* 51 Action14 <- <{  p.addParamRefValue(text) }> */
		nil,
		/* 52 Action15 <- <{ p.addParamCidrValue(text) }> */
		nil,
		/* 53 Action16 <- <{ p.addParamIpValue(text) }> */
		nil,
		/* 54 Action17 <- <{p.addCsvValue(text)}> */
		nil,
		/* 55 Action18 <- <{ p.addParamValue(text) }> */
		nil,
		/* 56 Action19 <- <{ p.LineDone() }> */
		nil,
		/* 57 Action20 <- <{ p.LineDone() }> */
		nil,
	}
	p.rules = _rules
}
