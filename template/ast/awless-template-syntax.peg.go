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
	ruleIntRangeValue
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
	ruleAction12
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
	"IntRangeValue",
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
	"Action12",
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
	rules  [41]func() bool
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
			p.AddParamValue(text)
		case ruleAction11:
			p.AddParamIntValue(text)
		case ruleAction12:
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
											if buffer[position] != rune('-') {
												goto l77
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l77
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
											add(ruleIntRangeValue, position79)
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
										position86 := position
										{
											position87 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l85
											}
											position++
										l88:
											{
												position89, tokenIndex89 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l89
												}
												position++
												goto l88
											l89:
												position, tokenIndex = position89, tokenIndex89
											}
											add(ruleIntValue, position87)
										}
										add(rulePegText, position86)
									}
									{
										add(ruleAction11, position)
									}
									goto l50
								l85:
									position, tokenIndex = position50, tokenIndex50
									{
										switch buffer[position] {
										case '$':
											{
												position92 := position
												if buffer[position] != rune('$') {
													goto l41
												}
												position++
												{
													position93 := position
													if !_rules[ruleIdentifier]() {
														goto l41
													}
													add(rulePegText, position93)
												}
												add(ruleRefValue, position92)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position95 := position
												if buffer[position] != rune('@') {
													goto l41
												}
												position++
												{
													position96 := position
													if !_rules[ruleIdentifier]() {
														goto l41
													}
													add(rulePegText, position96)
												}
												add(ruleAliasValue, position95)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position98 := position
												if buffer[position] != rune('{') {
													goto l41
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l41
												}
												{
													position99 := position
													if !_rules[ruleIdentifier]() {
														goto l41
													}
													add(rulePegText, position99)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l41
												}
												if buffer[position] != rune('}') {
													goto l41
												}
												position++
												add(ruleHoleValue, position98)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position101 := position
												{
													position102 := position
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

												l103:
													{
														position104, tokenIndex104 := position, tokenIndex
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l104
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l104
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l104
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l104
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l104
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l104
																}
																position++
																break
															}
														}

														goto l103
													l104:
														position, tokenIndex = position104, tokenIndex104
													}
													add(ruleStringValue, position102)
												}
												add(rulePegText, position101)
											}
											{
												add(ruleAction12, position)
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
								position108 := position
								{
									position109 := position
									if !_rules[ruleIdentifier]() {
										goto l45
									}
									add(rulePegText, position109)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l45
								}
								{
									position111 := position
									{
										position112, tokenIndex112 := position, tokenIndex
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
												if buffer[position] != rune('/') {
													goto l113
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l113
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
												add(ruleCidrValue, position115)
											}
											add(rulePegText, position114)
										}
										{
											add(ruleAction8, position)
										}
										goto l112
									l113:
										position, tokenIndex = position112, tokenIndex112
										{
											position128 := position
											{
												position129 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l127
												}
												position++
											l130:
												{
													position131, tokenIndex131 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l131
													}
													position++
													goto l130
												l131:
													position, tokenIndex = position131, tokenIndex131
												}
												if !matchDot() {
													goto l127
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l127
												}
												position++
											l132:
												{
													position133, tokenIndex133 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l133
													}
													position++
													goto l132
												l133:
													position, tokenIndex = position133, tokenIndex133
												}
												if !matchDot() {
													goto l127
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l127
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
												if !matchDot() {
													goto l127
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l127
												}
												position++
											l136:
												{
													position137, tokenIndex137 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l137
													}
													position++
													goto l136
												l137:
													position, tokenIndex = position137, tokenIndex137
												}
												add(ruleIpValue, position129)
											}
											add(rulePegText, position128)
										}
										{
											add(ruleAction9, position)
										}
										goto l112
									l127:
										position, tokenIndex = position112, tokenIndex112
										{
											position140 := position
											{
												position141 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l139
												}
												position++
											l142:
												{
													position143, tokenIndex143 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l143
													}
													position++
													goto l142
												l143:
													position, tokenIndex = position143, tokenIndex143
												}
												if buffer[position] != rune('-') {
													goto l139
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l139
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
												add(ruleIntRangeValue, position141)
											}
											add(rulePegText, position140)
										}
										{
											add(ruleAction10, position)
										}
										goto l112
									l139:
										position, tokenIndex = position112, tokenIndex112
										{
											position148 := position
											{
												position149 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l147
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
												add(ruleIntValue, position149)
											}
											add(rulePegText, position148)
										}
										{
											add(ruleAction11, position)
										}
										goto l112
									l147:
										position, tokenIndex = position112, tokenIndex112
										{
											switch buffer[position] {
											case '$':
												{
													position154 := position
													if buffer[position] != rune('$') {
														goto l45
													}
													position++
													{
														position155 := position
														if !_rules[ruleIdentifier]() {
															goto l45
														}
														add(rulePegText, position155)
													}
													add(ruleRefValue, position154)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position157 := position
													if buffer[position] != rune('@') {
														goto l45
													}
													position++
													{
														position158 := position
														if !_rules[ruleIdentifier]() {
															goto l45
														}
														add(rulePegText, position158)
													}
													add(ruleAliasValue, position157)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position160 := position
													if buffer[position] != rune('{') {
														goto l45
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l45
													}
													{
														position161 := position
														if !_rules[ruleIdentifier]() {
															goto l45
														}
														add(rulePegText, position161)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l45
													}
													if buffer[position] != rune('}') {
														goto l45
													}
													position++
													add(ruleHoleValue, position160)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position163 := position
													{
														position164 := position
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

													l165:
														{
															position166, tokenIndex166 := position, tokenIndex
															{
																switch buffer[position] {
																case '_':
																	if buffer[position] != rune('_') {
																		goto l166
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
																		goto l166
																	}
																	position++
																	break
																case '-':
																	if buffer[position] != rune('-') {
																		goto l166
																	}
																	position++
																	break
																case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l166
																	}
																	position++
																	break
																case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																	if c := buffer[position]; c < rune('A') || c > rune('Z') {
																		goto l166
																	}
																	position++
																	break
																default:
																	if c := buffer[position]; c < rune('a') || c > rune('z') {
																		goto l166
																	}
																	position++
																	break
																}
															}

															goto l165
														l166:
															position, tokenIndex = position166, tokenIndex166
														}
														add(ruleStringValue, position164)
													}
													add(rulePegText, position163)
												}
												{
													add(ruleAction12, position)
												}
												break
											}
										}

									}
								l112:
									add(ruleValue, position111)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l45
								}
								add(ruleParam, position108)
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
			position173, tokenIndex173 := position, tokenIndex
			{
				position174 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l173
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l173
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l173
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l173
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l173
						}
						position++
						break
					}
				}

			l175:
				{
					position176, tokenIndex176 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l176
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l176
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l176
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l176
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l176
							}
							position++
							break
						}
					}

					goto l175
				l176:
					position, tokenIndex = position176, tokenIndex176
				}
				add(ruleIdentifier, position174)
			}
			return true
		l173:
			position, tokenIndex = position173, tokenIndex173
			return false
		},
		/* 9 Value <- <((<CidrValue> Action8) / (<IpValue> Action9) / (<IntRangeValue> Action10) / (<IntValue> Action11) / ((&('$') (RefValue Action7)) | (&('@') (AliasValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action12))))> */
		nil,
		/* 10 StringValue <- <((&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		nil,
		/* 11 CidrValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+ '/' [0-9]+)> */
		nil,
		/* 12 IpValue <- <([0-9]+ . [0-9]+ . [0-9]+ . [0-9]+)> */
		nil,
		/* 13 IntValue <- <[0-9]+> */
		nil,
		/* 14 IntRangeValue <- <([0-9]+ '-' [0-9]+)> */
		nil,
		/* 15 RefValue <- <('$' <Identifier>)> */
		nil,
		/* 16 AliasValue <- <('@' <Identifier>)> */
		nil,
		/* 17 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 18 Spacing <- <Space*> */
		func() bool {
			{
				position189 := position
			l190:
				{
					position191, tokenIndex191 := position, tokenIndex
					{
						position192 := position
						{
							position193, tokenIndex193 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l194
							}
							goto l193
						l194:
							position, tokenIndex = position193, tokenIndex193
							if !_rules[ruleEndOfLine]() {
								goto l191
							}
						}
					l193:
						add(ruleSpace, position192)
					}
					goto l190
				l191:
					position, tokenIndex = position191, tokenIndex191
				}
				add(ruleSpacing, position189)
			}
			return true
		},
		/* 19 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position196 := position
			l197:
				{
					position198, tokenIndex198 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l198
					}
					goto l197
				l198:
					position, tokenIndex = position198, tokenIndex198
				}
				add(ruleWhiteSpacing, position196)
			}
			return true
		},
		/* 20 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position199, tokenIndex199 := position, tokenIndex
			{
				position200 := position
				if !_rules[ruleWhitespace]() {
					goto l199
				}
			l201:
				{
					position202, tokenIndex202 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l202
					}
					goto l201
				l202:
					position, tokenIndex = position202, tokenIndex202
				}
				add(ruleMustWhiteSpacing, position200)
			}
			return true
		l199:
			position, tokenIndex = position199, tokenIndex199
			return false
		},
		/* 21 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position203, tokenIndex203 := position, tokenIndex
			{
				position204 := position
				if !_rules[ruleSpacing]() {
					goto l203
				}
				if buffer[position] != rune('=') {
					goto l203
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l203
				}
				add(ruleEqual, position204)
			}
			return true
		l203:
			position, tokenIndex = position203, tokenIndex203
			return false
		},
		/* 22 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 23 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position206, tokenIndex206 := position, tokenIndex
			{
				position207 := position
				{
					position208, tokenIndex208 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex = position208, tokenIndex208
					if buffer[position] != rune('\t') {
						goto l206
					}
					position++
				}
			l208:
				add(ruleWhitespace, position207)
			}
			return true
		l206:
			position, tokenIndex = position206, tokenIndex206
			return false
		},
		/* 24 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position210, tokenIndex210 := position, tokenIndex
			{
				position211 := position
				{
					position212, tokenIndex212 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l213
					}
					position++
					if buffer[position] != rune('\n') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex = position212, tokenIndex212
					if buffer[position] != rune('\n') {
						goto l214
					}
					position++
					goto l212
				l214:
					position, tokenIndex = position212, tokenIndex212
					if buffer[position] != rune('\r') {
						goto l210
					}
					position++
				}
			l212:
				add(ruleEndOfLine, position211)
			}
			return true
		l210:
			position, tokenIndex = position210, tokenIndex210
			return false
		},
		/* 25 EndOfFile <- <!.> */
		nil,
		nil,
		/* 28 Action0 <- <{ p.AddDeclarationIdentifier(text) }> */
		nil,
		/* 29 Action1 <- <{ p.AddAction(text) }> */
		nil,
		/* 30 Action2 <- <{ p.AddEntity(text) }> */
		nil,
		/* 31 Action3 <- <{ p.EndOfParams() }> */
		nil,
		/* 32 Action4 <- <{ p.AddParamKey(text) }> */
		nil,
		/* 33 Action5 <- <{  p.AddParamHoleValue(text) }> */
		nil,
		/* 34 Action6 <- <{  p.AddParamAliasValue(text) }> */
		nil,
		/* 35 Action7 <- <{  p.AddParamRefValue(text) }> */
		nil,
		/* 36 Action8 <- <{ p.AddParamCidrValue(text) }> */
		nil,
		/* 37 Action9 <- <{ p.AddParamIpValue(text) }> */
		nil,
		/* 38 Action10 <- <{ p.AddParamValue(text) }> */
		nil,
		/* 39 Action11 <- <{ p.AddParamIntValue(text) }> */
		nil,
		/* 40 Action12 <- <{ p.AddParamValue(text) }> */
		nil,
	}
	p.rules = _rules
}
