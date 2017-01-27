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
	ruleAliasValue
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
	ruleAction11
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
	"AliasValue",
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
	"Action11",
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
	rules  [39]func() bool
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
			p.AddParamAliasValue(text)
		case ruleAction7:
			p.AddParamRefValue(text)
		case ruleAction8:
			p.AddParamCidrValue(text)
		case ruleAction9:
			p.AddParamIpValue(text)
		case ruleAction10:
			p.AddParamIntValue(text)
		case ruleAction11:
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
		/* 2 Action <- <(('s' 't' 'a' 'r' 't') / ((&('a') ('a' 't' 't' 'a' 'c' 'h')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('s') ('s' 't' 'o' 'p')) | (&('d') ('d' 'e' 'l' 'e' 't' 'e')) | (&('c') ('c' 'r' 'e' 'a' 't' 'e'))))> */
		nil,
		/* 3 Entity <- <(('v' 'p' 'c') / ('s' 'u' 'b' 'n' 'e' 't') / ((&('s') ('s' 'e' 'c' 'u' 'r' 'i' 't' 'y' 'g' 'r' 'o' 'u' 'p')) | (&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('r') ('r' 'o' 'l' 'e')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('t') ('t' 'a' 'g' 's')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e')) | (&('i') ('i' 'n' 's' 't' 'a' 'n' 'c' 'e'))))> */
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
							if buffer[position] != rune('s') {
								goto l31
							}
							position++
							if buffer[position] != rune('t') {
								goto l31
							}
							position++
							if buffer[position] != rune('a') {
								goto l31
							}
							position++
							if buffer[position] != rune('r') {
								goto l31
							}
							position++
							if buffer[position] != rune('t') {
								goto l31
							}
							position++
							goto l30
						l31:
							position, tokenIndex = position30, tokenIndex30
							{
								switch buffer[position] {
								case 'a':
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
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
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('h') {
										goto l26
									}
									position++
									break
								case 'u':
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									if buffer[position] != rune('d') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
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
									break
								case 's':
									if buffer[position] != rune('s') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
										goto l26
									}
									position++
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									break
								case 'd':
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
									break
								default:
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
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
									break
								}
							}

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
					position34 := position
					{
						position35 := position
						{
							position36, tokenIndex36 := position, tokenIndex
							if buffer[position] != rune('v') {
								goto l37
							}
							position++
							if buffer[position] != rune('p') {
								goto l37
							}
							position++
							if buffer[position] != rune('c') {
								goto l37
							}
							position++
							goto l36
						l37:
							position, tokenIndex = position36, tokenIndex36
							if buffer[position] != rune('s') {
								goto l38
							}
							position++
							if buffer[position] != rune('u') {
								goto l38
							}
							position++
							if buffer[position] != rune('b') {
								goto l38
							}
							position++
							if buffer[position] != rune('n') {
								goto l38
							}
							position++
							if buffer[position] != rune('e') {
								goto l38
							}
							position++
							if buffer[position] != rune('t') {
								goto l38
							}
							position++
							goto l36
						l38:
							position, tokenIndex = position36, tokenIndex36
							{
								switch buffer[position] {
								case 's':
									if buffer[position] != rune('s') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('i') {
										goto l26
									}
									position++
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('y') {
										goto l26
									}
									position++
									if buffer[position] != rune('g') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
										goto l26
									}
									position++
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									break
								case 'k':
									if buffer[position] != rune('k') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('y') {
										goto l26
									}
									position++
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('i') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									break
								case 'p':
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
										goto l26
									}
									position++
									if buffer[position] != rune('l') {
										goto l26
									}
									position++
									if buffer[position] != rune('i') {
										goto l26
									}
									position++
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('y') {
										goto l26
									}
									position++
									break
								case 'r':
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
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
									break
								case 'g':
									if buffer[position] != rune('g') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
										goto l26
									}
									position++
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('p') {
										goto l26
									}
									position++
									break
								case 'u':
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('s') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									break
								case 't':
									if buffer[position] != rune('t') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('g') {
										goto l26
									}
									position++
									if buffer[position] != rune('s') {
										goto l26
									}
									position++
									break
								case 'v':
									if buffer[position] != rune('v') {
										goto l26
									}
									position++
									if buffer[position] != rune('o') {
										goto l26
									}
									position++
									if buffer[position] != rune('l') {
										goto l26
									}
									position++
									if buffer[position] != rune('u') {
										goto l26
									}
									position++
									if buffer[position] != rune('m') {
										goto l26
									}
									position++
									if buffer[position] != rune('e') {
										goto l26
									}
									position++
									break
								default:
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
								}
							}

						}
					l36:
						add(ruleEntity, position35)
					}
					add(rulePegText, position34)
				}
				{
					add(ruleAction2, position)
				}
				{
					position41, tokenIndex41 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l41
					}
					{
						position43 := position
						{
							position46 := position
							{
								position47 := position
								if !_rules[ruleIdentifier]() {
									goto l41
								}
								add(rulePegText, position47)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l41
							}
							{
								position49 := position
								{
									position50, tokenIndex50 := position, tokenIndex
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
										add(ruleAction8, position)
									}
									goto l50
								l51:
									position, tokenIndex = position50, tokenIndex50
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
										add(ruleAction9, position)
									}
									goto l50
								l65:
									position, tokenIndex = position50, tokenIndex50
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
										add(ruleAction10, position)
									}
									goto l50
								l77:
									position, tokenIndex = position50, tokenIndex50
									{
										switch buffer[position] {
										case '$':
											{
												position84 := position
												if buffer[position] != rune('$') {
													goto l41
												}
												position++
												{
													position85 := position
													if !_rules[ruleIdentifier]() {
														goto l41
													}
													add(rulePegText, position85)
												}
												add(ruleRefValue, position84)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position87 := position
												if buffer[position] != rune('@') {
													goto l41
												}
												position++
												{
													position88 := position
													if !_rules[ruleIdentifier]() {
														goto l41
													}
													add(rulePegText, position88)
												}
												add(ruleAliasValue, position87)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position90 := position
												if buffer[position] != rune('{') {
													goto l41
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l41
												}
												{
													position91 := position
													if !_rules[ruleIdentifier]() {
														goto l41
													}
													add(rulePegText, position91)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l41
												}
												if buffer[position] != rune('}') {
													goto l41
												}
												position++
												add(ruleHoleValue, position90)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position93 := position
												{
													position94 := position
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

												l95:
													{
														position96, tokenIndex96 := position, tokenIndex
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l96
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l96
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l96
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l96
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l96
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l96
																}
																position++
																break
															}
														}

														goto l95
													l96:
														position, tokenIndex = position96, tokenIndex96
													}
													add(ruleStringValue, position94)
												}
												add(rulePegText, position93)
											}
											{
												add(ruleAction11, position)
											}
											break
										}
									}

								}
							l50:
								add(ruleValue, position49)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l41
							}
							add(ruleParam, position46)
						}
					l44:
						{
							position45, tokenIndex45 := position, tokenIndex
							{
								position100 := position
								{
									position101 := position
									if !_rules[ruleIdentifier]() {
										goto l45
									}
									add(rulePegText, position101)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l45
								}
								{
									position103 := position
									{
										position104, tokenIndex104 := position, tokenIndex
										{
											position106 := position
											{
												position107 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l105
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
												if !matchDot() {
													goto l105
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l105
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
												if !matchDot() {
													goto l105
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l105
												}
												position++
											l112:
												{
													position113, tokenIndex113 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l113
													}
													position++
													goto l112
												l113:
													position, tokenIndex = position113, tokenIndex113
												}
												if !matchDot() {
													goto l105
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l105
												}
												position++
											l114:
												{
													position115, tokenIndex115 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l115
													}
													position++
													goto l114
												l115:
													position, tokenIndex = position115, tokenIndex115
												}
												if buffer[position] != rune('/') {
													goto l105
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l105
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
												add(ruleCidrValue, position107)
											}
											add(rulePegText, position106)
										}
										{
											add(ruleAction8, position)
										}
										goto l104
									l105:
										position, tokenIndex = position104, tokenIndex104
										{
											position120 := position
											{
												position121 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l119
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
												if !matchDot() {
													goto l119
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l119
												}
												position++
											l124:
												{
													position125, tokenIndex125 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l125
													}
													position++
													goto l124
												l125:
													position, tokenIndex = position125, tokenIndex125
												}
												if !matchDot() {
													goto l119
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l119
												}
												position++
											l126:
												{
													position127, tokenIndex127 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l127
													}
													position++
													goto l126
												l127:
													position, tokenIndex = position127, tokenIndex127
												}
												if !matchDot() {
													goto l119
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l119
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
												add(ruleIpValue, position121)
											}
											add(rulePegText, position120)
										}
										{
											add(ruleAction9, position)
										}
										goto l104
									l119:
										position, tokenIndex = position104, tokenIndex104
										{
											position132 := position
											{
												position133 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l131
												}
												position++
											l134:
												{
													position135, tokenIndex135 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l135
													}
													position++
													goto l134
												l135:
													position, tokenIndex = position135, tokenIndex135
												}
												add(ruleIntValue, position133)
											}
											add(rulePegText, position132)
										}
										{
											add(ruleAction10, position)
										}
										goto l104
									l131:
										position, tokenIndex = position104, tokenIndex104
										{
											switch buffer[position] {
											case '$':
												{
													position138 := position
													if buffer[position] != rune('$') {
														goto l45
													}
													position++
													{
														position139 := position
														if !_rules[ruleIdentifier]() {
															goto l45
														}
														add(rulePegText, position139)
													}
													add(ruleRefValue, position138)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position141 := position
													if buffer[position] != rune('@') {
														goto l45
													}
													position++
													{
														position142 := position
														if !_rules[ruleIdentifier]() {
															goto l45
														}
														add(rulePegText, position142)
													}
													add(ruleAliasValue, position141)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position144 := position
													if buffer[position] != rune('{') {
														goto l45
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l45
													}
													{
														position145 := position
														if !_rules[ruleIdentifier]() {
															goto l45
														}
														add(rulePegText, position145)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l45
													}
													if buffer[position] != rune('}') {
														goto l45
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
													{
														position148 := position
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l45
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l45
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l45
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l45
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l45
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l45
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
																case '_':
																	if buffer[position] != rune('_') {
																		goto l150
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
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
																case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
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
														add(ruleStringValue, position148)
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
								l104:
									add(ruleValue, position103)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l45
								}
								add(ruleParam, position100)
							}
							goto l44
						l45:
							position, tokenIndex = position45, tokenIndex45
						}
						{
							add(ruleAction3, position)
						}
						add(ruleParams, position43)
					}
					goto l42
				l41:
					position, tokenIndex = position41, tokenIndex41
				}
			l42:
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
			position157, tokenIndex157 := position, tokenIndex
			{
				position158 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l157
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l157
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l157
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l157
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l157
						}
						position++
						break
					}
				}

			l159:
				{
					position160, tokenIndex160 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l160
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l160
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l160
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l160
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l160
							}
							position++
							break
						}
					}

					goto l159
				l160:
					position, tokenIndex = position160, tokenIndex160
				}
				add(ruleIdentifier, position158)
			}
			return true
		l157:
			position, tokenIndex = position157, tokenIndex157
			return false
		},
		/* 9 Value <- <((<CidrValue> Action8) / (<IpValue> Action9) / (<IntValue> Action10) / ((&('$') (RefValue Action7)) | (&('@') (AliasValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action11))))> */
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
		/* 15 AliasValue <- <('@' <Identifier>)> */
		nil,
		/* 16 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 17 Spacing <- <Space*> */
		func() bool {
			{
				position172 := position
			l173:
				{
					position174, tokenIndex174 := position, tokenIndex
					{
						position175 := position
						{
							position176, tokenIndex176 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l177
							}
							goto l176
						l177:
							position, tokenIndex = position176, tokenIndex176
							if !_rules[ruleEndOfLine]() {
								goto l174
							}
						}
					l176:
						add(ruleSpace, position175)
					}
					goto l173
				l174:
					position, tokenIndex = position174, tokenIndex174
				}
				add(ruleSpacing, position172)
			}
			return true
		},
		/* 18 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position179 := position
			l180:
				{
					position181, tokenIndex181 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l181
					}
					goto l180
				l181:
					position, tokenIndex = position181, tokenIndex181
				}
				add(ruleWhiteSpacing, position179)
			}
			return true
		},
		/* 19 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position182, tokenIndex182 := position, tokenIndex
			{
				position183 := position
				if !_rules[ruleWhitespace]() {
					goto l182
				}
			l184:
				{
					position185, tokenIndex185 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l185
					}
					goto l184
				l185:
					position, tokenIndex = position185, tokenIndex185
				}
				add(ruleMustWhiteSpacing, position183)
			}
			return true
		l182:
			position, tokenIndex = position182, tokenIndex182
			return false
		},
		/* 20 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position186, tokenIndex186 := position, tokenIndex
			{
				position187 := position
				if !_rules[ruleSpacing]() {
					goto l186
				}
				if buffer[position] != rune('=') {
					goto l186
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l186
				}
				add(ruleEqual, position187)
			}
			return true
		l186:
			position, tokenIndex = position186, tokenIndex186
			return false
		},
		/* 21 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 22 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position189, tokenIndex189 := position, tokenIndex
			{
				position190 := position
				{
					position191, tokenIndex191 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l192
					}
					position++
					goto l191
				l192:
					position, tokenIndex = position191, tokenIndex191
					if buffer[position] != rune('\t') {
						goto l189
					}
					position++
				}
			l191:
				add(ruleWhitespace, position190)
			}
			return true
		l189:
			position, tokenIndex = position189, tokenIndex189
			return false
		},
		/* 23 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position193, tokenIndex193 := position, tokenIndex
			{
				position194 := position
				{
					position195, tokenIndex195 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l196
					}
					position++
					if buffer[position] != rune('\n') {
						goto l196
					}
					position++
					goto l195
				l196:
					position, tokenIndex = position195, tokenIndex195
					if buffer[position] != rune('\n') {
						goto l197
					}
					position++
					goto l195
				l197:
					position, tokenIndex = position195, tokenIndex195
					if buffer[position] != rune('\r') {
						goto l193
					}
					position++
				}
			l195:
				add(ruleEndOfLine, position194)
			}
			return true
		l193:
			position, tokenIndex = position193, tokenIndex193
			return false
		},
		/* 24 EndOfFile <- <!.> */
		nil,
		nil,
		/* 27 Action0 <- <{ p.AddDeclarationIdentifier(text) }> */
		nil,
		/* 28 Action1 <- <{ p.AddAction(text) }> */
		nil,
		/* 29 Action2 <- <{ p.AddEntity(text) }> */
		nil,
		/* 30 Action3 <- <{ p.EndOfParams() }> */
		nil,
		/* 31 Action4 <- <{ p.AddParamKey(text) }> */
		nil,
		/* 32 Action5 <- <{  p.AddParamHoleValue(text) }> */
		nil,
		/* 33 Action6 <- <{  p.AddParamAliasValue(text) }> */
		nil,
		/* 34 Action7 <- <{  p.AddParamRefValue(text) }> */
		nil,
		/* 35 Action8 <- <{ p.AddParamCidrValue(text) }> */
		nil,
		/* 36 Action9 <- <{ p.AddParamIpValue(text) }> */
		nil,
		/* 37 Action10 <- <{ p.AddParamIntValue(text) }> */
		nil,
		/* 38 Action11 <- <{ p.AddParamValue(text) }> */
		nil,
	}
	p.rules = _rules
}
