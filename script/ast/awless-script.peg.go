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
									position48, tokenIndex48 := position, tokenIndex
									{
										position50 := position
										if buffer[position] != rune('{') {
											goto l49
										}
										position++
										{
											position51 := position
											if !_rules[ruleIdentifier]() {
												goto l49
											}
											add(rulePegText, position51)
										}
										if buffer[position] != rune('}') {
											goto l49
										}
										position++
										add(ruleHoleValue, position50)
									}
									{
										add(ruleAction5, position)
									}
									goto l48
								l49:
									position, tokenIndex = position48, tokenIndex48
									{
										position54 := position
										{
											position55 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l53
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
												goto l53
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l53
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
												goto l53
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l53
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
											if !matchDot() {
												goto l53
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l53
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
											if buffer[position] != rune('/') {
												goto l53
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l53
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
											add(ruleCidrValue, position55)
										}
										add(rulePegText, position54)
									}
									{
										add(ruleAction6, position)
									}
									goto l48
								l53:
									position, tokenIndex = position48, tokenIndex48
									{
										position68 := position
										{
											position69 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l67
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
												goto l67
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l67
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
												goto l67
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l67
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
											if !matchDot() {
												goto l67
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l67
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
											add(ruleIpValue, position69)
										}
										add(rulePegText, position68)
									}
									{
										add(ruleAction7, position)
									}
									goto l48
								l67:
									position, tokenIndex = position48, tokenIndex48
									{
										position80 := position
										{
											position81 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l79
											}
											position++
										l82:
											{
												position83, tokenIndex83 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l83
												}
												position++
												goto l82
											l83:
												position, tokenIndex = position83, tokenIndex83
											}
											add(ruleIntValue, position81)
										}
										add(rulePegText, position80)
									}
									{
										add(ruleAction8, position)
									}
									goto l48
								l79:
									position, tokenIndex = position48, tokenIndex48
									{
										position85 := position
										{
											position86 := position
											{
												switch buffer[position] {
												case '_':
													if buffer[position] != rune('_') {
														goto l40
													}
													position++
													break
												case '.':
													if buffer[position] != rune('.') {
														goto l40
													}
													position++
													break
												case '-':
													if buffer[position] != rune('-') {
														goto l40
													}
													position++
													break
												case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l40
													}
													position++
													break
												case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
													if c := buffer[position]; c < rune('A') || c > rune('Z') {
														goto l40
													}
													position++
													break
												default:
													if c := buffer[position]; c < rune('a') || c > rune('z') {
														goto l40
													}
													position++
													break
												}
											}

										l87:
											{
												position88, tokenIndex88 := position, tokenIndex
												{
													switch buffer[position] {
													case '_':
														if buffer[position] != rune('_') {
															goto l88
														}
														position++
														break
													case '.':
														if buffer[position] != rune('.') {
															goto l88
														}
														position++
														break
													case '-':
														if buffer[position] != rune('-') {
															goto l88
														}
														position++
														break
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l88
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l88
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l88
														}
														position++
														break
													}
												}

												goto l87
											l88:
												position, tokenIndex = position88, tokenIndex88
											}
											add(ruleStringValue, position86)
										}
										add(rulePegText, position85)
									}
									{
										add(ruleAction9, position)
									}
								}
							l48:
								add(ruleValue, position47)
							}
							{
								position92 := position
							l93:
								{
									position94, tokenIndex94 := position, tokenIndex
									if !_rules[ruleWhitespace]() {
										goto l94
									}
									goto l93
								l94:
									position, tokenIndex = position94, tokenIndex94
								}
								add(ruleWhiteSpacing, position92)
							}
							add(ruleParam, position44)
						}
					l42:
						{
							position43, tokenIndex43 := position, tokenIndex
							{
								position95 := position
								{
									position96 := position
									if !_rules[ruleIdentifier]() {
										goto l43
									}
									add(rulePegText, position96)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l43
								}
								{
									position98 := position
									{
										position99, tokenIndex99 := position, tokenIndex
										{
											position101 := position
											if buffer[position] != rune('{') {
												goto l100
											}
											position++
											{
												position102 := position
												if !_rules[ruleIdentifier]() {
													goto l100
												}
												add(rulePegText, position102)
											}
											if buffer[position] != rune('}') {
												goto l100
											}
											position++
											add(ruleHoleValue, position101)
										}
										{
											add(ruleAction5, position)
										}
										goto l99
									l100:
										position, tokenIndex = position99, tokenIndex99
										{
											position105 := position
											{
												position106 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l104
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
													goto l104
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l104
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
													goto l104
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l104
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
												if !matchDot() {
													goto l104
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l104
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
												if buffer[position] != rune('/') {
													goto l104
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l104
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
												add(ruleCidrValue, position106)
											}
											add(rulePegText, position105)
										}
										{
											add(ruleAction6, position)
										}
										goto l99
									l104:
										position, tokenIndex = position99, tokenIndex99
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
												if !matchDot() {
													goto l118
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l118
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
													goto l118
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l118
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
												if !matchDot() {
													goto l118
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l118
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
												add(ruleIpValue, position120)
											}
											add(rulePegText, position119)
										}
										{
											add(ruleAction7, position)
										}
										goto l99
									l118:
										position, tokenIndex = position99, tokenIndex99
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
											add(ruleAction8, position)
										}
										goto l99
									l130:
										position, tokenIndex = position99, tokenIndex99
										{
											position136 := position
											{
												position137 := position
												{
													switch buffer[position] {
													case '_':
														if buffer[position] != rune('_') {
															goto l43
														}
														position++
														break
													case '.':
														if buffer[position] != rune('.') {
															goto l43
														}
														position++
														break
													case '-':
														if buffer[position] != rune('-') {
															goto l43
														}
														position++
														break
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l43
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l43
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l43
														}
														position++
														break
													}
												}

											l138:
												{
													position139, tokenIndex139 := position, tokenIndex
													{
														switch buffer[position] {
														case '_':
															if buffer[position] != rune('_') {
																goto l139
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l139
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l139
															}
															position++
															break
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l139
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l139
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l139
															}
															position++
															break
														}
													}

													goto l138
												l139:
													position, tokenIndex = position139, tokenIndex139
												}
												add(ruleStringValue, position137)
											}
											add(rulePegText, position136)
										}
										{
											add(ruleAction9, position)
										}
									}
								l99:
									add(ruleValue, position98)
								}
								{
									position143 := position
								l144:
									{
										position145, tokenIndex145 := position, tokenIndex
										if !_rules[ruleWhitespace]() {
											goto l145
										}
										goto l144
									l145:
										position, tokenIndex = position145, tokenIndex145
									}
									add(ruleWhiteSpacing, position143)
								}
								add(ruleParam, position95)
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
		/* 8 Identifier <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position149, tokenIndex149 := position, tokenIndex
			{
				position150 := position
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l149
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l149
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l149
						}
						position++
						break
					}
				}

			l151:
				{
					position152, tokenIndex152 := position, tokenIndex
					{
						switch buffer[position] {
						case '_':
							if buffer[position] != rune('_') {
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

					goto l151
				l152:
					position, tokenIndex = position152, tokenIndex152
				}
				add(ruleIdentifier, position150)
			}
			return true
		l149:
			position, tokenIndex = position149, tokenIndex149
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
		/* 14 HoleValue <- <('{' <Identifier> '}')> */
		nil,
		/* 15 Spacing <- <Space*> */
		func() bool {
			{
				position162 := position
			l163:
				{
					position164, tokenIndex164 := position, tokenIndex
					{
						position165 := position
						{
							position166, tokenIndex166 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l167
							}
							goto l166
						l167:
							position, tokenIndex = position166, tokenIndex166
							if !_rules[ruleEndOfLine]() {
								goto l164
							}
						}
					l166:
						add(ruleSpace, position165)
					}
					goto l163
				l164:
					position, tokenIndex = position164, tokenIndex164
				}
				add(ruleSpacing, position162)
			}
			return true
		},
		/* 16 WhiteSpacing <- <Whitespace*> */
		nil,
		/* 17 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position169, tokenIndex169 := position, tokenIndex
			{
				position170 := position
				if !_rules[ruleWhitespace]() {
					goto l169
				}
			l171:
				{
					position172, tokenIndex172 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l172
					}
					goto l171
				l172:
					position, tokenIndex = position172, tokenIndex172
				}
				add(ruleMustWhiteSpacing, position170)
			}
			return true
		l169:
			position, tokenIndex = position169, tokenIndex169
			return false
		},
		/* 18 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position173, tokenIndex173 := position, tokenIndex
			{
				position174 := position
				if !_rules[ruleSpacing]() {
					goto l173
				}
				if buffer[position] != rune('=') {
					goto l173
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l173
				}
				add(ruleEqual, position174)
			}
			return true
		l173:
			position, tokenIndex = position173, tokenIndex173
			return false
		},
		/* 19 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 20 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position176, tokenIndex176 := position, tokenIndex
			{
				position177 := position
				{
					position178, tokenIndex178 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l179
					}
					position++
					goto l178
				l179:
					position, tokenIndex = position178, tokenIndex178
					if buffer[position] != rune('\t') {
						goto l176
					}
					position++
				}
			l178:
				add(ruleWhitespace, position177)
			}
			return true
		l176:
			position, tokenIndex = position176, tokenIndex176
			return false
		},
		/* 21 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position180, tokenIndex180 := position, tokenIndex
			{
				position181 := position
				{
					position182, tokenIndex182 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l183
					}
					position++
					if buffer[position] != rune('\n') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex = position182, tokenIndex182
					if buffer[position] != rune('\n') {
						goto l184
					}
					position++
					goto l182
				l184:
					position, tokenIndex = position182, tokenIndex182
					if buffer[position] != rune('\r') {
						goto l180
					}
					position++
				}
			l182:
				add(ruleEndOfLine, position181)
			}
			return true
		l180:
			position, tokenIndex = position180, tokenIndex180
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
