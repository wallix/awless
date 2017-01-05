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
	*Script

	Buffer string
	buffer []rune
	rules  [31]func() bool
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
			p.AddParamValue(text)
		case ruleAction6:
			p.AddParamIntValue(text)
		case ruleAction7:
			p.AddParamHoleValue(text)

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
		/* 5 Expr <- <(<Action> Action1 MustWhiteSpacing (<Entity> MustWhiteSpacing Action2)? Params*)> */
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
					position33, tokenIndex33 := position, tokenIndex
					{
						position35 := position
						{
							position36 := position
							{
								switch buffer[position] {
								case 'i':
									if buffer[position] != rune('i') {
										goto l33
									}
									position++
									if buffer[position] != rune('n') {
										goto l33
									}
									position++
									if buffer[position] != rune('s') {
										goto l33
									}
									position++
									if buffer[position] != rune('t') {
										goto l33
									}
									position++
									if buffer[position] != rune('a') {
										goto l33
									}
									position++
									if buffer[position] != rune('n') {
										goto l33
									}
									position++
									if buffer[position] != rune('c') {
										goto l33
									}
									position++
									if buffer[position] != rune('e') {
										goto l33
									}
									position++
									break
								case 's':
									if buffer[position] != rune('s') {
										goto l33
									}
									position++
									if buffer[position] != rune('u') {
										goto l33
									}
									position++
									if buffer[position] != rune('b') {
										goto l33
									}
									position++
									if buffer[position] != rune('n') {
										goto l33
									}
									position++
									if buffer[position] != rune('e') {
										goto l33
									}
									position++
									if buffer[position] != rune('t') {
										goto l33
									}
									position++
									break
								default:
									if buffer[position] != rune('v') {
										goto l33
									}
									position++
									if buffer[position] != rune('p') {
										goto l33
									}
									position++
									if buffer[position] != rune('c') {
										goto l33
									}
									position++
									break
								}
							}

							add(ruleEntity, position36)
						}
						add(rulePegText, position35)
					}
					if !_rules[ruleMustWhiteSpacing]() {
						goto l33
					}
					{
						add(ruleAction2, position)
					}
					goto l34
				l33:
					position, tokenIndex = position33, tokenIndex33
				}
			l34:
			l39:
				{
					position40, tokenIndex40 := position, tokenIndex
					{
						position41 := position
						{
							position44 := position
							{
								position45 := position
								if !_rules[ruleIdentifier]() {
									goto l40
								}
								add(rulePegText, position45)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l40
							}
							{
								position47 := position
								{
									switch buffer[position] {
									case '{':
										{
											position49 := position
											if buffer[position] != rune('{') {
												goto l40
											}
											position++
											{
												position50 := position
												if !_rules[ruleIdentifier]() {
													goto l40
												}
												add(rulePegText, position50)
											}
											if buffer[position] != rune('}') {
												goto l40
											}
											position++
											add(ruleHoleValue, position49)
										}
										{
											add(ruleAction7, position)
										}
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										{
											position52 := position
											{
												position53 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l40
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
												add(ruleIntValue, position53)
											}
											add(rulePegText, position52)
										}
										{
											add(ruleAction6, position)
										}
										break
									default:
										{
											position57 := position
											{
												position58 := position
												{
													position61, tokenIndex61 := position, tokenIndex
													if c := buffer[position]; c < rune('a') || c > rune('z') {
														goto l62
													}
													position++
													goto l61
												l62:
													position, tokenIndex = position61, tokenIndex61
													if c := buffer[position]; c < rune('A') || c > rune('Z') {
														goto l40
													}
													position++
												}
											l61:
											l59:
												{
													position60, tokenIndex60 := position, tokenIndex
													{
														position63, tokenIndex63 := position, tokenIndex
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l64
														}
														position++
														goto l63
													l64:
														position, tokenIndex = position63, tokenIndex63
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l60
														}
														position++
													}
												l63:
													goto l59
												l60:
													position, tokenIndex = position60, tokenIndex60
												}
												add(ruleStringValue, position58)
											}
											add(rulePegText, position57)
										}
										{
											add(ruleAction5, position)
										}
										break
									}
								}

								add(ruleValue, position47)
							}
							{
								position66 := position
							l67:
								{
									position68, tokenIndex68 := position, tokenIndex
									if !_rules[ruleWhitespace]() {
										goto l68
									}
									goto l67
								l68:
									position, tokenIndex = position68, tokenIndex68
								}
								add(ruleWhiteSpacing, position66)
							}
							add(ruleParam, position44)
						}
					l42:
						{
							position43, tokenIndex43 := position, tokenIndex
							{
								position69 := position
								{
									position70 := position
									if !_rules[ruleIdentifier]() {
										goto l43
									}
									add(rulePegText, position70)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l43
								}
								{
									position72 := position
									{
										switch buffer[position] {
										case '{':
											{
												position74 := position
												if buffer[position] != rune('{') {
													goto l43
												}
												position++
												{
													position75 := position
													if !_rules[ruleIdentifier]() {
														goto l43
													}
													add(rulePegText, position75)
												}
												if buffer[position] != rune('}') {
													goto l43
												}
												position++
												add(ruleHoleValue, position74)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
											{
												position77 := position
												{
													position78 := position
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l43
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
												add(ruleAction6, position)
											}
											break
										default:
											{
												position82 := position
												{
													position83 := position
													{
														position86, tokenIndex86 := position, tokenIndex
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l87
														}
														position++
														goto l86
													l87:
														position, tokenIndex = position86, tokenIndex86
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l43
														}
														position++
													}
												l86:
												l84:
													{
														position85, tokenIndex85 := position, tokenIndex
														{
															position88, tokenIndex88 := position, tokenIndex
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l89
															}
															position++
															goto l88
														l89:
															position, tokenIndex = position88, tokenIndex88
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l85
															}
															position++
														}
													l88:
														goto l84
													l85:
														position, tokenIndex = position85, tokenIndex85
													}
													add(ruleStringValue, position83)
												}
												add(rulePegText, position82)
											}
											{
												add(ruleAction5, position)
											}
											break
										}
									}

									add(ruleValue, position72)
								}
								{
									position91 := position
								l92:
									{
										position93, tokenIndex93 := position, tokenIndex
										if !_rules[ruleWhitespace]() {
											goto l93
										}
										goto l92
									l93:
										position, tokenIndex = position93, tokenIndex93
									}
									add(ruleWhiteSpacing, position91)
								}
								add(ruleParam, position69)
							}
							goto l42
						l43:
							position, tokenIndex = position43, tokenIndex43
						}
						{
							add(ruleAction3, position)
						}
						add(ruleParams, position41)
					}
					goto l39
				l40:
					position, tokenIndex = position40, tokenIndex40
				}
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
		/* 8 Identifier <- <([a-z] / [A-Z])+> */
		func() bool {
			position97, tokenIndex97 := position, tokenIndex
			{
				position98 := position
				{
					position101, tokenIndex101 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l102
					}
					position++
					goto l101
				l102:
					position, tokenIndex = position101, tokenIndex101
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l97
					}
					position++
				}
			l101:
			l99:
				{
					position100, tokenIndex100 := position, tokenIndex
					{
						position103, tokenIndex103 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l104
						}
						position++
						goto l103
					l104:
						position, tokenIndex = position103, tokenIndex103
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l100
						}
						position++
					}
				l103:
					goto l99
				l100:
					position, tokenIndex = position100, tokenIndex100
				}
				add(ruleIdentifier, position98)
			}
			return true
		l97:
			position, tokenIndex = position97, tokenIndex97
			return false
		},
		/* 9 Value <- <((&('{') (HoleValue Action7)) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (<IntValue> Action6)) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action5)))> */
		nil,
		/* 10 StringValue <- <([a-z] / [A-Z])+> */
		nil,
		/* 11 IntValue <- <[0-9]+> */
		nil,
		/* 12 HoleValue <- <('{' <Identifier> '}')> */
		nil,
		/* 13 Spacing <- <Space*> */
		func() bool {
			{
				position110 := position
			l111:
				{
					position112, tokenIndex112 := position, tokenIndex
					{
						position113 := position
						{
							position114, tokenIndex114 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l115
							}
							goto l114
						l115:
							position, tokenIndex = position114, tokenIndex114
							if !_rules[ruleEndOfLine]() {
								goto l112
							}
						}
					l114:
						add(ruleSpace, position113)
					}
					goto l111
				l112:
					position, tokenIndex = position112, tokenIndex112
				}
				add(ruleSpacing, position110)
			}
			return true
		},
		/* 14 WhiteSpacing <- <Whitespace*> */
		nil,
		/* 15 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position117, tokenIndex117 := position, tokenIndex
			{
				position118 := position
				if !_rules[ruleWhitespace]() {
					goto l117
				}
			l119:
				{
					position120, tokenIndex120 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l120
					}
					goto l119
				l120:
					position, tokenIndex = position120, tokenIndex120
				}
				add(ruleMustWhiteSpacing, position118)
			}
			return true
		l117:
			position, tokenIndex = position117, tokenIndex117
			return false
		},
		/* 16 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position121, tokenIndex121 := position, tokenIndex
			{
				position122 := position
				if !_rules[ruleSpacing]() {
					goto l121
				}
				if buffer[position] != rune('=') {
					goto l121
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l121
				}
				add(ruleEqual, position122)
			}
			return true
		l121:
			position, tokenIndex = position121, tokenIndex121
			return false
		},
		/* 17 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 18 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position124, tokenIndex124 := position, tokenIndex
			{
				position125 := position
				{
					position126, tokenIndex126 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex = position126, tokenIndex126
					if buffer[position] != rune('\t') {
						goto l124
					}
					position++
				}
			l126:
				add(ruleWhitespace, position125)
			}
			return true
		l124:
			position, tokenIndex = position124, tokenIndex124
			return false
		},
		/* 19 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position128, tokenIndex128 := position, tokenIndex
			{
				position129 := position
				{
					position130, tokenIndex130 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l131
					}
					position++
					if buffer[position] != rune('\n') {
						goto l131
					}
					position++
					goto l130
				l131:
					position, tokenIndex = position130, tokenIndex130
					if buffer[position] != rune('\n') {
						goto l132
					}
					position++
					goto l130
				l132:
					position, tokenIndex = position130, tokenIndex130
					if buffer[position] != rune('\r') {
						goto l128
					}
					position++
				}
			l130:
				add(ruleEndOfLine, position129)
			}
			return true
		l128:
			position, tokenIndex = position128, tokenIndex128
			return false
		},
		/* 20 EndOfFile <- <!.> */
		nil,
		nil,
		/* 23 Action0 <- <{ p.AddDeclarationIdentifier(text) }> */
		nil,
		/* 24 Action1 <- <{ p.AddAction(text) }> */
		nil,
		/* 25 Action2 <- <{ p.AddEntity(text) }> */
		nil,
		/* 26 Action3 <- <{ p.EndOfParams() }> */
		nil,
		/* 27 Action4 <- <{ p.AddParamKey(text) }> */
		nil,
		/* 28 Action5 <- <{ p.AddParamValue(text) }> */
		nil,
		/* 29 Action6 <- <{ p.AddParamIntValue(text) }> */
		nil,
		/* 30 Action7 <- <{  p.AddParamHoleValue(text) }> */
		nil,
	}
	p.rules = _rules
}
