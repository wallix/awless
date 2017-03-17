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
		/* 0 Script <- <(BlankLine* Statement+ BlankLine* WhiteSpacing EndOfFile)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					if !_rules[ruleBlankLine]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
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
			l4:
				{
					position5, tokenIndex5 := position, tokenIndex
					{
						position25 := position
						if !_rules[ruleWhiteSpacing]() {
							goto l5
						}
						{
							position26, tokenIndex26 := position, tokenIndex
							if !_rules[ruleExpr]() {
								goto l27
							}
							goto l26
						l27:
							position, tokenIndex = position26, tokenIndex26
							{
								position29 := position
								{
									position30 := position
									if !_rules[ruleIdentifier]() {
										goto l28
									}
									add(rulePegText, position30)
								}
								{
									add(ruleAction0, position)
								}
								if !_rules[ruleEqual]() {
									goto l28
								}
								if !_rules[ruleExpr]() {
									goto l28
								}
								add(ruleDeclaration, position29)
							}
							goto l26
						l28:
							position, tokenIndex = position26, tokenIndex26
							{
								position32 := position
								{
									position33, tokenIndex33 := position, tokenIndex
									if buffer[position] != rune('#') {
										goto l34
									}
									position++
								l35:
									{
										position36, tokenIndex36 := position, tokenIndex
										{
											position37, tokenIndex37 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l37
											}
											goto l36
										l37:
											position, tokenIndex = position37, tokenIndex37
										}
										if !matchDot() {
											goto l36
										}
										goto l35
									l36:
										position, tokenIndex = position36, tokenIndex36
									}
									goto l33
								l34:
									position, tokenIndex = position33, tokenIndex33
									if buffer[position] != rune('/') {
										goto l5
									}
									position++
									if buffer[position] != rune('/') {
										goto l5
									}
									position++
								l38:
									{
										position39, tokenIndex39 := position, tokenIndex
										{
											position40, tokenIndex40 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l40
											}
											goto l39
										l40:
											position, tokenIndex = position40, tokenIndex40
										}
										if !matchDot() {
											goto l39
										}
										goto l38
									l39:
										position, tokenIndex = position39, tokenIndex39
									}
									{
										add(ruleAction14, position)
									}
								}
							l33:
								add(ruleComment, position32)
							}
						}
					l26:
						if !_rules[ruleWhiteSpacing]() {
							goto l5
						}
					l42:
						{
							position43, tokenIndex43 := position, tokenIndex
							if !_rules[ruleEndOfLine]() {
								goto l43
							}
							goto l42
						l43:
							position, tokenIndex = position43, tokenIndex43
						}
						add(ruleStatement, position25)
					}
					goto l4
				l5:
					position, tokenIndex = position5, tokenIndex5
				}
			l44:
				{
					position45, tokenIndex45 := position, tokenIndex
					if !_rules[ruleBlankLine]() {
						goto l45
					}
					goto l44
				l45:
					position, tokenIndex = position45, tokenIndex45
				}
				if !_rules[ruleWhiteSpacing]() {
					goto l0
				}
				{
					position46 := position
					{
						position47, tokenIndex47 := position, tokenIndex
						if !matchDot() {
							goto l47
						}
						goto l0
					l47:
						position, tokenIndex = position47, tokenIndex47
					}
					add(ruleEndOfFile, position46)
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
		/* 2 Action <- <(('c' 'r' 'e' 'a' 't' 'e') / ('d' 'e' 'l' 'e' 't' 'e') / ('s' 't' 'a' 'r' 't') / ((&('d') ('d' 'e' 't' 'a' 'c' 'h')) | (&('c') ('c' 'h' 'e' 'c' 'k')) | (&('a') ('a' 't' 't' 'a' 'c' 'h')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('s') ('s' 't' 'o' 'p')) | (&('n') ('n' 'o' 'n' 'e'))))> */
		nil,
		/* 3 Entity <- <(('v' 'p' 'c') / ('s' 'u' 'b' 'n' 'e' 't') / ('i' 'n' 's' 't' 'a' 'n' 'c' 'e') / ('t' 'a' 'g') / ('s' 'e' 'c' 'u' 'r' 'i' 't' 'y' 'g' 'r' 'o' 'u' 'p') / ('r' 'o' 'u' 't' 'e' 't' 'a' 'b' 'l' 'e') / ('r' 'o' 'u' 't' 'e') / ('r' 'o' 'l' 'e') / ('s' 't' 'o' 'r' 'a' 'g' 'e' 'o' 'b' 'j' 'e' 'c' 't') / ('t' 'o' 'p' 'i' 'c') / ('l' 'o' 'a' 'd' 'b' 'a' 'l' 'a' 'n' 'c' 'e' 'r') / ((&('r') ('r' 'e' 'c' 'o' 'r' 'd')) | (&('z') ('z' 'o' 'n' 'e')) | (&('t') ('t' 'a' 'r' 'g' 'e' 't' 'g' 'r' 'o' 'u' 'p')) | (&('l') ('l' 'i' 's' 't' 'e' 'n' 'e' 'r')) | (&('q') ('q' 'u' 'e' 'u' 'e')) | (&('s') ('s' 'u' 'b' 's' 'c' 'r' 'i' 'p' 't' 'i' 'o' 'n')) | (&('b') ('b' 'u' 'c' 'k' 'e' 't')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('i') ('i' 'n' 't' 'e' 'r' 'n' 'e' 't' 'g' 'a' 't' 'e' 'w' 'a' 'y')) | (&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e')) | (&('n') ('n' 'o' 'n' 'e'))))> */
		nil,
		/* 4 Declaration <- <(<Identifier> Action0 Equal Expr)> */
		nil,
		/* 5 Expr <- <(<Action> Action1 MustWhiteSpacing <Entity> Action2 (MustWhiteSpacing Params)? Action3)> */
		func() bool {
			position52, tokenIndex52 := position, tokenIndex
			{
				position53 := position
				{
					position54 := position
					{
						position55 := position
						{
							position56, tokenIndex56 := position, tokenIndex
							if buffer[position] != rune('c') {
								goto l57
							}
							position++
							if buffer[position] != rune('r') {
								goto l57
							}
							position++
							if buffer[position] != rune('e') {
								goto l57
							}
							position++
							if buffer[position] != rune('a') {
								goto l57
							}
							position++
							if buffer[position] != rune('t') {
								goto l57
							}
							position++
							if buffer[position] != rune('e') {
								goto l57
							}
							position++
							goto l56
						l57:
							position, tokenIndex = position56, tokenIndex56
							if buffer[position] != rune('d') {
								goto l58
							}
							position++
							if buffer[position] != rune('e') {
								goto l58
							}
							position++
							if buffer[position] != rune('l') {
								goto l58
							}
							position++
							if buffer[position] != rune('e') {
								goto l58
							}
							position++
							if buffer[position] != rune('t') {
								goto l58
							}
							position++
							if buffer[position] != rune('e') {
								goto l58
							}
							position++
							goto l56
						l58:
							position, tokenIndex = position56, tokenIndex56
							if buffer[position] != rune('s') {
								goto l59
							}
							position++
							if buffer[position] != rune('t') {
								goto l59
							}
							position++
							if buffer[position] != rune('a') {
								goto l59
							}
							position++
							if buffer[position] != rune('r') {
								goto l59
							}
							position++
							if buffer[position] != rune('t') {
								goto l59
							}
							position++
							goto l56
						l59:
							position, tokenIndex = position56, tokenIndex56
							{
								switch buffer[position] {
								case 'd':
									if buffer[position] != rune('d') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('a') {
										goto l52
									}
									position++
									if buffer[position] != rune('c') {
										goto l52
									}
									position++
									if buffer[position] != rune('h') {
										goto l52
									}
									position++
									break
								case 'c':
									if buffer[position] != rune('c') {
										goto l52
									}
									position++
									if buffer[position] != rune('h') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('c') {
										goto l52
									}
									position++
									if buffer[position] != rune('k') {
										goto l52
									}
									position++
									break
								case 'a':
									if buffer[position] != rune('a') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('a') {
										goto l52
									}
									position++
									if buffer[position] != rune('c') {
										goto l52
									}
									position++
									if buffer[position] != rune('h') {
										goto l52
									}
									position++
									break
								case 'u':
									if buffer[position] != rune('u') {
										goto l52
									}
									position++
									if buffer[position] != rune('p') {
										goto l52
									}
									position++
									if buffer[position] != rune('d') {
										goto l52
									}
									position++
									if buffer[position] != rune('a') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									break
								case 's':
									if buffer[position] != rune('s') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('o') {
										goto l52
									}
									position++
									if buffer[position] != rune('p') {
										goto l52
									}
									position++
									break
								default:
									if buffer[position] != rune('n') {
										goto l52
									}
									position++
									if buffer[position] != rune('o') {
										goto l52
									}
									position++
									if buffer[position] != rune('n') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									break
								}
							}

						}
					l56:
						add(ruleAction, position55)
					}
					add(rulePegText, position54)
				}
				{
					add(ruleAction1, position)
				}
				if !_rules[ruleMustWhiteSpacing]() {
					goto l52
				}
				{
					position62 := position
					{
						position63 := position
						{
							position64, tokenIndex64 := position, tokenIndex
							if buffer[position] != rune('v') {
								goto l65
							}
							position++
							if buffer[position] != rune('p') {
								goto l65
							}
							position++
							if buffer[position] != rune('c') {
								goto l65
							}
							position++
							goto l64
						l65:
							position, tokenIndex = position64, tokenIndex64
							if buffer[position] != rune('s') {
								goto l66
							}
							position++
							if buffer[position] != rune('u') {
								goto l66
							}
							position++
							if buffer[position] != rune('b') {
								goto l66
							}
							position++
							if buffer[position] != rune('n') {
								goto l66
							}
							position++
							if buffer[position] != rune('e') {
								goto l66
							}
							position++
							if buffer[position] != rune('t') {
								goto l66
							}
							position++
							goto l64
						l66:
							position, tokenIndex = position64, tokenIndex64
							if buffer[position] != rune('i') {
								goto l67
							}
							position++
							if buffer[position] != rune('n') {
								goto l67
							}
							position++
							if buffer[position] != rune('s') {
								goto l67
							}
							position++
							if buffer[position] != rune('t') {
								goto l67
							}
							position++
							if buffer[position] != rune('a') {
								goto l67
							}
							position++
							if buffer[position] != rune('n') {
								goto l67
							}
							position++
							if buffer[position] != rune('c') {
								goto l67
							}
							position++
							if buffer[position] != rune('e') {
								goto l67
							}
							position++
							goto l64
						l67:
							position, tokenIndex = position64, tokenIndex64
							if buffer[position] != rune('t') {
								goto l68
							}
							position++
							if buffer[position] != rune('a') {
								goto l68
							}
							position++
							if buffer[position] != rune('g') {
								goto l68
							}
							position++
							goto l64
						l68:
							position, tokenIndex = position64, tokenIndex64
							if buffer[position] != rune('s') {
								goto l69
							}
							position++
							if buffer[position] != rune('e') {
								goto l69
							}
							position++
							if buffer[position] != rune('c') {
								goto l69
							}
							position++
							if buffer[position] != rune('u') {
								goto l69
							}
							position++
							if buffer[position] != rune('r') {
								goto l69
							}
							position++
							if buffer[position] != rune('i') {
								goto l69
							}
							position++
							if buffer[position] != rune('t') {
								goto l69
							}
							position++
							if buffer[position] != rune('y') {
								goto l69
							}
							position++
							if buffer[position] != rune('g') {
								goto l69
							}
							position++
							if buffer[position] != rune('r') {
								goto l69
							}
							position++
							if buffer[position] != rune('o') {
								goto l69
							}
							position++
							if buffer[position] != rune('u') {
								goto l69
							}
							position++
							if buffer[position] != rune('p') {
								goto l69
							}
							position++
							goto l64
						l69:
							position, tokenIndex = position64, tokenIndex64
							if buffer[position] != rune('r') {
								goto l70
							}
							position++
							if buffer[position] != rune('o') {
								goto l70
							}
							position++
							if buffer[position] != rune('u') {
								goto l70
							}
							position++
							if buffer[position] != rune('t') {
								goto l70
							}
							position++
							if buffer[position] != rune('e') {
								goto l70
							}
							position++
							if buffer[position] != rune('t') {
								goto l70
							}
							position++
							if buffer[position] != rune('a') {
								goto l70
							}
							position++
							if buffer[position] != rune('b') {
								goto l70
							}
							position++
							if buffer[position] != rune('l') {
								goto l70
							}
							position++
							if buffer[position] != rune('e') {
								goto l70
							}
							position++
							goto l64
						l70:
							position, tokenIndex = position64, tokenIndex64
							if buffer[position] != rune('r') {
								goto l71
							}
							position++
							if buffer[position] != rune('o') {
								goto l71
							}
							position++
							if buffer[position] != rune('u') {
								goto l71
							}
							position++
							if buffer[position] != rune('t') {
								goto l71
							}
							position++
							if buffer[position] != rune('e') {
								goto l71
							}
							position++
							goto l64
						l71:
							position, tokenIndex = position64, tokenIndex64
							if buffer[position] != rune('r') {
								goto l72
							}
							position++
							if buffer[position] != rune('o') {
								goto l72
							}
							position++
							if buffer[position] != rune('l') {
								goto l72
							}
							position++
							if buffer[position] != rune('e') {
								goto l72
							}
							position++
							goto l64
						l72:
							position, tokenIndex = position64, tokenIndex64
							if buffer[position] != rune('s') {
								goto l73
							}
							position++
							if buffer[position] != rune('t') {
								goto l73
							}
							position++
							if buffer[position] != rune('o') {
								goto l73
							}
							position++
							if buffer[position] != rune('r') {
								goto l73
							}
							position++
							if buffer[position] != rune('a') {
								goto l73
							}
							position++
							if buffer[position] != rune('g') {
								goto l73
							}
							position++
							if buffer[position] != rune('e') {
								goto l73
							}
							position++
							if buffer[position] != rune('o') {
								goto l73
							}
							position++
							if buffer[position] != rune('b') {
								goto l73
							}
							position++
							if buffer[position] != rune('j') {
								goto l73
							}
							position++
							if buffer[position] != rune('e') {
								goto l73
							}
							position++
							if buffer[position] != rune('c') {
								goto l73
							}
							position++
							if buffer[position] != rune('t') {
								goto l73
							}
							position++
							goto l64
						l73:
							position, tokenIndex = position64, tokenIndex64
							if buffer[position] != rune('t') {
								goto l74
							}
							position++
							if buffer[position] != rune('o') {
								goto l74
							}
							position++
							if buffer[position] != rune('p') {
								goto l74
							}
							position++
							if buffer[position] != rune('i') {
								goto l74
							}
							position++
							if buffer[position] != rune('c') {
								goto l74
							}
							position++
							goto l64
						l74:
							position, tokenIndex = position64, tokenIndex64
							if buffer[position] != rune('l') {
								goto l75
							}
							position++
							if buffer[position] != rune('o') {
								goto l75
							}
							position++
							if buffer[position] != rune('a') {
								goto l75
							}
							position++
							if buffer[position] != rune('d') {
								goto l75
							}
							position++
							if buffer[position] != rune('b') {
								goto l75
							}
							position++
							if buffer[position] != rune('a') {
								goto l75
							}
							position++
							if buffer[position] != rune('l') {
								goto l75
							}
							position++
							if buffer[position] != rune('a') {
								goto l75
							}
							position++
							if buffer[position] != rune('n') {
								goto l75
							}
							position++
							if buffer[position] != rune('c') {
								goto l75
							}
							position++
							if buffer[position] != rune('e') {
								goto l75
							}
							position++
							if buffer[position] != rune('r') {
								goto l75
							}
							position++
							goto l64
						l75:
							position, tokenIndex = position64, tokenIndex64
							{
								switch buffer[position] {
								case 'r':
									if buffer[position] != rune('r') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('c') {
										goto l52
									}
									position++
									if buffer[position] != rune('o') {
										goto l52
									}
									position++
									if buffer[position] != rune('r') {
										goto l52
									}
									position++
									if buffer[position] != rune('d') {
										goto l52
									}
									position++
									break
								case 'z':
									if buffer[position] != rune('z') {
										goto l52
									}
									position++
									if buffer[position] != rune('o') {
										goto l52
									}
									position++
									if buffer[position] != rune('n') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									break
								case 't':
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('a') {
										goto l52
									}
									position++
									if buffer[position] != rune('r') {
										goto l52
									}
									position++
									if buffer[position] != rune('g') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('g') {
										goto l52
									}
									position++
									if buffer[position] != rune('r') {
										goto l52
									}
									position++
									if buffer[position] != rune('o') {
										goto l52
									}
									position++
									if buffer[position] != rune('u') {
										goto l52
									}
									position++
									if buffer[position] != rune('p') {
										goto l52
									}
									position++
									break
								case 'l':
									if buffer[position] != rune('l') {
										goto l52
									}
									position++
									if buffer[position] != rune('i') {
										goto l52
									}
									position++
									if buffer[position] != rune('s') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('n') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('r') {
										goto l52
									}
									position++
									break
								case 'q':
									if buffer[position] != rune('q') {
										goto l52
									}
									position++
									if buffer[position] != rune('u') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('u') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									break
								case 's':
									if buffer[position] != rune('s') {
										goto l52
									}
									position++
									if buffer[position] != rune('u') {
										goto l52
									}
									position++
									if buffer[position] != rune('b') {
										goto l52
									}
									position++
									if buffer[position] != rune('s') {
										goto l52
									}
									position++
									if buffer[position] != rune('c') {
										goto l52
									}
									position++
									if buffer[position] != rune('r') {
										goto l52
									}
									position++
									if buffer[position] != rune('i') {
										goto l52
									}
									position++
									if buffer[position] != rune('p') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('i') {
										goto l52
									}
									position++
									if buffer[position] != rune('o') {
										goto l52
									}
									position++
									if buffer[position] != rune('n') {
										goto l52
									}
									position++
									break
								case 'b':
									if buffer[position] != rune('b') {
										goto l52
									}
									position++
									if buffer[position] != rune('u') {
										goto l52
									}
									position++
									if buffer[position] != rune('c') {
										goto l52
									}
									position++
									if buffer[position] != rune('k') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									break
								case 'p':
									if buffer[position] != rune('p') {
										goto l52
									}
									position++
									if buffer[position] != rune('o') {
										goto l52
									}
									position++
									if buffer[position] != rune('l') {
										goto l52
									}
									position++
									if buffer[position] != rune('i') {
										goto l52
									}
									position++
									if buffer[position] != rune('c') {
										goto l52
									}
									position++
									if buffer[position] != rune('y') {
										goto l52
									}
									position++
									break
								case 'g':
									if buffer[position] != rune('g') {
										goto l52
									}
									position++
									if buffer[position] != rune('r') {
										goto l52
									}
									position++
									if buffer[position] != rune('o') {
										goto l52
									}
									position++
									if buffer[position] != rune('u') {
										goto l52
									}
									position++
									if buffer[position] != rune('p') {
										goto l52
									}
									position++
									break
								case 'u':
									if buffer[position] != rune('u') {
										goto l52
									}
									position++
									if buffer[position] != rune('s') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('r') {
										goto l52
									}
									position++
									break
								case 'i':
									if buffer[position] != rune('i') {
										goto l52
									}
									position++
									if buffer[position] != rune('n') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('r') {
										goto l52
									}
									position++
									if buffer[position] != rune('n') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('g') {
										goto l52
									}
									position++
									if buffer[position] != rune('a') {
										goto l52
									}
									position++
									if buffer[position] != rune('t') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('w') {
										goto l52
									}
									position++
									if buffer[position] != rune('a') {
										goto l52
									}
									position++
									if buffer[position] != rune('y') {
										goto l52
									}
									position++
									break
								case 'k':
									if buffer[position] != rune('k') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									if buffer[position] != rune('y') {
										goto l52
									}
									position++
									if buffer[position] != rune('p') {
										goto l52
									}
									position++
									if buffer[position] != rune('a') {
										goto l52
									}
									position++
									if buffer[position] != rune('i') {
										goto l52
									}
									position++
									if buffer[position] != rune('r') {
										goto l52
									}
									position++
									break
								case 'v':
									if buffer[position] != rune('v') {
										goto l52
									}
									position++
									if buffer[position] != rune('o') {
										goto l52
									}
									position++
									if buffer[position] != rune('l') {
										goto l52
									}
									position++
									if buffer[position] != rune('u') {
										goto l52
									}
									position++
									if buffer[position] != rune('m') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									break
								default:
									if buffer[position] != rune('n') {
										goto l52
									}
									position++
									if buffer[position] != rune('o') {
										goto l52
									}
									position++
									if buffer[position] != rune('n') {
										goto l52
									}
									position++
									if buffer[position] != rune('e') {
										goto l52
									}
									position++
									break
								}
							}

						}
					l64:
						add(ruleEntity, position63)
					}
					add(rulePegText, position62)
				}
				{
					add(ruleAction2, position)
				}
				{
					position78, tokenIndex78 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l78
					}
					{
						position80 := position
						{
							position83 := position
							{
								position84 := position
								if !_rules[ruleIdentifier]() {
									goto l78
								}
								add(rulePegText, position84)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l78
							}
							{
								position86 := position
								{
									position87, tokenIndex87 := position, tokenIndex
									{
										position89 := position
										{
											position90 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l88
											}
											position++
										l91:
											{
												position92, tokenIndex92 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l92
												}
												position++
												goto l91
											l92:
												position, tokenIndex = position92, tokenIndex92
											}
											if !matchDot() {
												goto l88
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l88
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
											if !matchDot() {
												goto l88
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l88
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
												goto l88
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l88
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
											if buffer[position] != rune('/') {
												goto l88
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l88
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
											add(ruleCidrValue, position90)
										}
										add(rulePegText, position89)
									}
									{
										add(ruleAction8, position)
									}
									goto l87
								l88:
									position, tokenIndex = position87, tokenIndex87
									{
										position103 := position
										{
											position104 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l102
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
												goto l102
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l102
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
												goto l102
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l102
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
												goto l102
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l102
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
											add(ruleIpValue, position104)
										}
										add(rulePegText, position103)
									}
									{
										add(ruleAction9, position)
									}
									goto l87
								l102:
									position, tokenIndex = position87, tokenIndex87
									{
										position115 := position
										{
											position116 := position
											if !_rules[ruleStringValue]() {
												goto l114
											}
											if !_rules[ruleWhiteSpacing]() {
												goto l114
											}
											if buffer[position] != rune(',') {
												goto l114
											}
											position++
											if !_rules[ruleWhiteSpacing]() {
												goto l114
											}
										l117:
											{
												position118, tokenIndex118 := position, tokenIndex
												if !_rules[ruleStringValue]() {
													goto l118
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l118
												}
												if buffer[position] != rune(',') {
													goto l118
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l118
												}
												goto l117
											l118:
												position, tokenIndex = position118, tokenIndex118
											}
											if !_rules[ruleStringValue]() {
												goto l114
											}
											add(ruleCSVValue, position116)
										}
										add(rulePegText, position115)
									}
									{
										add(ruleAction10, position)
									}
									goto l87
								l114:
									position, tokenIndex = position87, tokenIndex87
									{
										position121 := position
										{
											position122 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l120
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
											if buffer[position] != rune('-') {
												goto l120
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l120
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
											add(ruleIntRangeValue, position122)
										}
										add(rulePegText, position121)
									}
									{
										add(ruleAction11, position)
									}
									goto l87
								l120:
									position, tokenIndex = position87, tokenIndex87
									{
										position129 := position
										{
											position130 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l128
											}
											position++
										l131:
											{
												position132, tokenIndex132 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l132
												}
												position++
												goto l131
											l132:
												position, tokenIndex = position132, tokenIndex132
											}
											add(ruleIntValue, position130)
										}
										add(rulePegText, position129)
									}
									{
										add(ruleAction12, position)
									}
									goto l87
								l128:
									position, tokenIndex = position87, tokenIndex87
									{
										switch buffer[position] {
										case '$':
											{
												position135 := position
												if buffer[position] != rune('$') {
													goto l78
												}
												position++
												{
													position136 := position
													if !_rules[ruleIdentifier]() {
														goto l78
													}
													add(rulePegText, position136)
												}
												add(ruleRefValue, position135)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position138 := position
												{
													position139 := position
													if buffer[position] != rune('@') {
														goto l78
													}
													position++
													if !_rules[ruleStringValue]() {
														goto l78
													}
													add(rulePegText, position139)
												}
												add(ruleAliasValue, position138)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position141 := position
												if buffer[position] != rune('{') {
													goto l78
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l78
												}
												{
													position142 := position
													if !_rules[ruleIdentifier]() {
														goto l78
													}
													add(rulePegText, position142)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l78
												}
												if buffer[position] != rune('}') {
													goto l78
												}
												position++
												add(ruleHoleValue, position141)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position144 := position
												if !_rules[ruleStringValue]() {
													goto l78
												}
												add(rulePegText, position144)
											}
											{
												add(ruleAction13, position)
											}
											break
										}
									}

								}
							l87:
								add(ruleValue, position86)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l78
							}
							add(ruleParam, position83)
						}
					l81:
						{
							position82, tokenIndex82 := position, tokenIndex
							{
								position146 := position
								{
									position147 := position
									if !_rules[ruleIdentifier]() {
										goto l82
									}
									add(rulePegText, position147)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l82
								}
								{
									position149 := position
									{
										position150, tokenIndex150 := position, tokenIndex
										{
											position152 := position
											{
												position153 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l151
												}
												position++
											l154:
												{
													position155, tokenIndex155 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l155
													}
													position++
													goto l154
												l155:
													position, tokenIndex = position155, tokenIndex155
												}
												if !matchDot() {
													goto l151
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l151
												}
												position++
											l156:
												{
													position157, tokenIndex157 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l157
													}
													position++
													goto l156
												l157:
													position, tokenIndex = position157, tokenIndex157
												}
												if !matchDot() {
													goto l151
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l151
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
													goto l151
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l151
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
												if buffer[position] != rune('/') {
													goto l151
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l151
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
												add(ruleCidrValue, position153)
											}
											add(rulePegText, position152)
										}
										{
											add(ruleAction8, position)
										}
										goto l150
									l151:
										position, tokenIndex = position150, tokenIndex150
										{
											position166 := position
											{
												position167 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l165
												}
												position++
											l168:
												{
													position169, tokenIndex169 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l169
													}
													position++
													goto l168
												l169:
													position, tokenIndex = position169, tokenIndex169
												}
												if !matchDot() {
													goto l165
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l165
												}
												position++
											l170:
												{
													position171, tokenIndex171 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l171
													}
													position++
													goto l170
												l171:
													position, tokenIndex = position171, tokenIndex171
												}
												if !matchDot() {
													goto l165
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l165
												}
												position++
											l172:
												{
													position173, tokenIndex173 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l173
													}
													position++
													goto l172
												l173:
													position, tokenIndex = position173, tokenIndex173
												}
												if !matchDot() {
													goto l165
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l165
												}
												position++
											l174:
												{
													position175, tokenIndex175 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l175
													}
													position++
													goto l174
												l175:
													position, tokenIndex = position175, tokenIndex175
												}
												add(ruleIpValue, position167)
											}
											add(rulePegText, position166)
										}
										{
											add(ruleAction9, position)
										}
										goto l150
									l165:
										position, tokenIndex = position150, tokenIndex150
										{
											position178 := position
											{
												position179 := position
												if !_rules[ruleStringValue]() {
													goto l177
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l177
												}
												if buffer[position] != rune(',') {
													goto l177
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l177
												}
											l180:
												{
													position181, tokenIndex181 := position, tokenIndex
													if !_rules[ruleStringValue]() {
														goto l181
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l181
													}
													if buffer[position] != rune(',') {
														goto l181
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l181
													}
													goto l180
												l181:
													position, tokenIndex = position181, tokenIndex181
												}
												if !_rules[ruleStringValue]() {
													goto l177
												}
												add(ruleCSVValue, position179)
											}
											add(rulePegText, position178)
										}
										{
											add(ruleAction10, position)
										}
										goto l150
									l177:
										position, tokenIndex = position150, tokenIndex150
										{
											position184 := position
											{
												position185 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l183
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
												if buffer[position] != rune('-') {
													goto l183
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l183
												}
												position++
											l188:
												{
													position189, tokenIndex189 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l189
													}
													position++
													goto l188
												l189:
													position, tokenIndex = position189, tokenIndex189
												}
												add(ruleIntRangeValue, position185)
											}
											add(rulePegText, position184)
										}
										{
											add(ruleAction11, position)
										}
										goto l150
									l183:
										position, tokenIndex = position150, tokenIndex150
										{
											position192 := position
											{
												position193 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l191
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
												add(ruleIntValue, position193)
											}
											add(rulePegText, position192)
										}
										{
											add(ruleAction12, position)
										}
										goto l150
									l191:
										position, tokenIndex = position150, tokenIndex150
										{
											switch buffer[position] {
											case '$':
												{
													position198 := position
													if buffer[position] != rune('$') {
														goto l82
													}
													position++
													{
														position199 := position
														if !_rules[ruleIdentifier]() {
															goto l82
														}
														add(rulePegText, position199)
													}
													add(ruleRefValue, position198)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position201 := position
													{
														position202 := position
														if buffer[position] != rune('@') {
															goto l82
														}
														position++
														if !_rules[ruleStringValue]() {
															goto l82
														}
														add(rulePegText, position202)
													}
													add(ruleAliasValue, position201)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position204 := position
													if buffer[position] != rune('{') {
														goto l82
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l82
													}
													{
														position205 := position
														if !_rules[ruleIdentifier]() {
															goto l82
														}
														add(rulePegText, position205)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l82
													}
													if buffer[position] != rune('}') {
														goto l82
													}
													position++
													add(ruleHoleValue, position204)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position207 := position
													if !_rules[ruleStringValue]() {
														goto l82
													}
													add(rulePegText, position207)
												}
												{
													add(ruleAction13, position)
												}
												break
											}
										}

									}
								l150:
									add(ruleValue, position149)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l82
								}
								add(ruleParam, position146)
							}
							goto l81
						l82:
							position, tokenIndex = position82, tokenIndex82
						}
						add(ruleParams, position80)
					}
					goto l79
				l78:
					position, tokenIndex = position78, tokenIndex78
				}
			l79:
				{
					add(ruleAction3, position)
				}
				add(ruleExpr, position53)
			}
			return true
		l52:
			position, tokenIndex = position52, tokenIndex52
			return false
		},
		/* 6 Params <- <Param+> */
		nil,
		/* 7 Param <- <(<Identifier> Action4 Equal Value WhiteSpacing)> */
		nil,
		/* 8 Identifier <- <((&('.') '.') | (&('_') '_') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position212, tokenIndex212 := position, tokenIndex
			{
				position213 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
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

			l214:
				{
					position215, tokenIndex215 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l215
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l215
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l215
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l215
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l215
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l215
							}
							position++
							break
						}
					}

					goto l214
				l215:
					position, tokenIndex = position215, tokenIndex215
				}
				add(ruleIdentifier, position213)
			}
			return true
		l212:
			position, tokenIndex = position212, tokenIndex212
			return false
		},
		/* 9 Value <- <((<CidrValue> Action8) / (<IpValue> Action9) / (<CSVValue> Action10) / (<IntRangeValue> Action11) / (<IntValue> Action12) / ((&('$') (RefValue Action7)) | (&('@') (AliasValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action13))))> */
		nil,
		/* 10 StringValue <- <((&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position219, tokenIndex219 := position, tokenIndex
			{
				position220 := position
				{
					switch buffer[position] {
					case '/':
						if buffer[position] != rune('/') {
							goto l219
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l219
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l219
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l219
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l219
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l219
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l219
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l219
						}
						position++
						break
					}
				}

			l221:
				{
					position222, tokenIndex222 := position, tokenIndex
					{
						switch buffer[position] {
						case '/':
							if buffer[position] != rune('/') {
								goto l222
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l222
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l222
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l222
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l222
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l222
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l222
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l222
							}
							position++
							break
						}
					}

					goto l221
				l222:
					position, tokenIndex = position222, tokenIndex222
				}
				add(ruleStringValue, position220)
			}
			return true
		l219:
			position, tokenIndex = position219, tokenIndex219
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
				position235 := position
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
				add(ruleWhiteSpacing, position235)
			}
			return true
		},
		/* 21 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position238, tokenIndex238 := position, tokenIndex
			{
				position239 := position
				if !_rules[ruleWhitespace]() {
					goto l238
				}
			l240:
				{
					position241, tokenIndex241 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l241
					}
					goto l240
				l241:
					position, tokenIndex = position241, tokenIndex241
				}
				add(ruleMustWhiteSpacing, position239)
			}
			return true
		l238:
			position, tokenIndex = position238, tokenIndex238
			return false
		},
		/* 22 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position242, tokenIndex242 := position, tokenIndex
			{
				position243 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l242
				}
				if buffer[position] != rune('=') {
					goto l242
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l242
				}
				add(ruleEqual, position243)
			}
			return true
		l242:
			position, tokenIndex = position242, tokenIndex242
			return false
		},
		/* 23 BlankLine <- <(WhiteSpacing EndOfLine Action15)> */
		func() bool {
			position244, tokenIndex244 := position, tokenIndex
			{
				position245 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l244
				}
				if !_rules[ruleEndOfLine]() {
					goto l244
				}
				{
					add(ruleAction15, position)
				}
				add(ruleBlankLine, position245)
			}
			return true
		l244:
			position, tokenIndex = position244, tokenIndex244
			return false
		},
		/* 24 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position247, tokenIndex247 := position, tokenIndex
			{
				position248 := position
				{
					position249, tokenIndex249 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l250
					}
					position++
					goto l249
				l250:
					position, tokenIndex = position249, tokenIndex249
					if buffer[position] != rune('\t') {
						goto l247
					}
					position++
				}
			l249:
				add(ruleWhitespace, position248)
			}
			return true
		l247:
			position, tokenIndex = position247, tokenIndex247
			return false
		},
		/* 25 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position251, tokenIndex251 := position, tokenIndex
			{
				position252 := position
				{
					position253, tokenIndex253 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l254
					}
					position++
					if buffer[position] != rune('\n') {
						goto l254
					}
					position++
					goto l253
				l254:
					position, tokenIndex = position253, tokenIndex253
					if buffer[position] != rune('\n') {
						goto l255
					}
					position++
					goto l253
				l255:
					position, tokenIndex = position253, tokenIndex253
					if buffer[position] != rune('\r') {
						goto l251
					}
					position++
				}
			l253:
				add(ruleEndOfLine, position252)
			}
			return true
		l251:
			position, tokenIndex = position251, tokenIndex251
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
