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
	ruleComment
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
	ruleAction13
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
	"Comment",
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
	"Action13",
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
	rules  [43]func() bool
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
			p.LineDone()
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
		case ruleAction13:
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
							position8 := position
							{
								position9 := position
								if !_rules[ruleIdentifier]() {
									goto l7
								}
								add(rulePegText, position9)
							}
							{
								add(ruleAction0, position)
							}
							if !_rules[ruleEqual]() {
								goto l7
							}
							if !_rules[ruleExpr]() {
								goto l7
							}
							add(ruleDeclaration, position8)
						}
						goto l5
					l7:
						position, tokenIndex = position5, tokenIndex5
						{
							position11 := position
							{
								position12, tokenIndex12 := position, tokenIndex
								if buffer[position] != rune('#') {
									goto l13
								}
								position++
							l14:
								{
									position15, tokenIndex15 := position, tokenIndex
									{
										position16, tokenIndex16 := position, tokenIndex
										if !_rules[ruleEndOfLine]() {
											goto l16
										}
										goto l15
									l16:
										position, tokenIndex = position16, tokenIndex16
									}
									if !matchDot() {
										goto l15
									}
									goto l14
								l15:
									position, tokenIndex = position15, tokenIndex15
								}
								goto l12
							l13:
								position, tokenIndex = position12, tokenIndex12
								if buffer[position] != rune('/') {
									goto l0
								}
								position++
								if buffer[position] != rune('/') {
									goto l0
								}
								position++
							l17:
								{
									position18, tokenIndex18 := position, tokenIndex
									{
										position19, tokenIndex19 := position, tokenIndex
										if !_rules[ruleEndOfLine]() {
											goto l19
										}
										goto l18
									l19:
										position, tokenIndex = position19, tokenIndex19
									}
									if !matchDot() {
										goto l18
									}
									goto l17
								l18:
									position, tokenIndex = position18, tokenIndex18
								}
								{
									add(ruleAction13, position)
								}
							}
						l12:
							add(ruleComment, position11)
						}
					}
				l5:
					if !_rules[ruleSpacing]() {
						goto l0
					}
				l21:
					{
						position22, tokenIndex22 := position, tokenIndex
						if !_rules[ruleEndOfLine]() {
							goto l22
						}
						goto l21
					l22:
						position, tokenIndex = position22, tokenIndex22
					}
					add(ruleStatement, position4)
				}
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					{
						position23 := position
						if !_rules[ruleSpacing]() {
							goto l3
						}
						{
							position24, tokenIndex24 := position, tokenIndex
							if !_rules[ruleExpr]() {
								goto l25
							}
							goto l24
						l25:
							position, tokenIndex = position24, tokenIndex24
							{
								position27 := position
								{
									position28 := position
									if !_rules[ruleIdentifier]() {
										goto l26
									}
									add(rulePegText, position28)
								}
								{
									add(ruleAction0, position)
								}
								if !_rules[ruleEqual]() {
									goto l26
								}
								if !_rules[ruleExpr]() {
									goto l26
								}
								add(ruleDeclaration, position27)
							}
							goto l24
						l26:
							position, tokenIndex = position24, tokenIndex24
							{
								position30 := position
								{
									position31, tokenIndex31 := position, tokenIndex
									if buffer[position] != rune('#') {
										goto l32
									}
									position++
								l33:
									{
										position34, tokenIndex34 := position, tokenIndex
										{
											position35, tokenIndex35 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l35
											}
											goto l34
										l35:
											position, tokenIndex = position35, tokenIndex35
										}
										if !matchDot() {
											goto l34
										}
										goto l33
									l34:
										position, tokenIndex = position34, tokenIndex34
									}
									goto l31
								l32:
									position, tokenIndex = position31, tokenIndex31
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
								l36:
									{
										position37, tokenIndex37 := position, tokenIndex
										{
											position38, tokenIndex38 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l38
											}
											goto l37
										l38:
											position, tokenIndex = position38, tokenIndex38
										}
										if !matchDot() {
											goto l37
										}
										goto l36
									l37:
										position, tokenIndex = position37, tokenIndex37
									}
									{
										add(ruleAction13, position)
									}
								}
							l31:
								add(ruleComment, position30)
							}
						}
					l24:
						if !_rules[ruleSpacing]() {
							goto l3
						}
					l40:
						{
							position41, tokenIndex41 := position, tokenIndex
							if !_rules[ruleEndOfLine]() {
								goto l41
							}
							goto l40
						l41:
							position, tokenIndex = position41, tokenIndex41
						}
						add(ruleStatement, position23)
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				{
					position42 := position
					{
						position43, tokenIndex43 := position, tokenIndex
						if !matchDot() {
							goto l43
						}
						goto l0
					l43:
						position, tokenIndex = position43, tokenIndex43
					}
					add(ruleEndOfFile, position42)
				}
				add(ruleScript, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Statement <- <(Spacing (Expr / Declaration / Comment) Spacing EndOfLine*)> */
		nil,
		/* 2 Action <- <(('c' 'r' 'e' 'a' 't' 'e') / ('d' 'e' 'l' 'e' 't' 'e') / ('s' 't' 'a' 'r' 't') / ((&('d') ('d' 'e' 't' 'a' 'c' 'h')) | (&('c') ('c' 'h' 'e' 'c' 'k')) | (&('a') ('a' 't' 't' 'a' 'c' 'h')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('s') ('s' 't' 'o' 'p'))))> */
		nil,
		/* 3 Entity <- <(('v' 'p' 'c') / ('s' 'u' 'b' 'n' 'e' 't') / ('i' 'n' 's' 't' 'a' 'n' 'c' 'e') / ('r' 'o' 'l' 'e') / ('s' 'e' 'c' 'u' 'r' 'i' 't' 'y' 'g' 'r' 'o' 'u' 'p') / ('r' 'o' 'u' 't' 'e' 't' 'a' 'b' 'l' 'e') / ((&('s') ('s' 't' 'o' 'r' 'a' 'g' 'e' 'o' 'b' 'j' 'e' 'c' 't')) | (&('b') ('b' 'u' 'c' 'k' 'e' 't')) | (&('r') ('r' 'o' 'u' 't' 'e')) | (&('i') ('i' 'n' 't' 'e' 'r' 'n' 'e' 't' 'g' 'a' 't' 'e' 'w' 'a' 'y')) | (&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('t') ('t' 'a' 'g' 's')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e'))))> */
		nil,
		/* 4 Declaration <- <(<Identifier> Action0 Equal Expr)> */
		nil,
		/* 5 Expr <- <(<Action> Action1 MustWhiteSpacing <Entity> Action2 (MustWhiteSpacing Params)? Action3)> */
		func() bool {
			position48, tokenIndex48 := position, tokenIndex
			{
				position49 := position
				{
					position50 := position
					{
						position51 := position
						{
							position52, tokenIndex52 := position, tokenIndex
							if buffer[position] != rune('c') {
								goto l53
							}
							position++
							if buffer[position] != rune('r') {
								goto l53
							}
							position++
							if buffer[position] != rune('e') {
								goto l53
							}
							position++
							if buffer[position] != rune('a') {
								goto l53
							}
							position++
							if buffer[position] != rune('t') {
								goto l53
							}
							position++
							if buffer[position] != rune('e') {
								goto l53
							}
							position++
							goto l52
						l53:
							position, tokenIndex = position52, tokenIndex52
							if buffer[position] != rune('d') {
								goto l54
							}
							position++
							if buffer[position] != rune('e') {
								goto l54
							}
							position++
							if buffer[position] != rune('l') {
								goto l54
							}
							position++
							if buffer[position] != rune('e') {
								goto l54
							}
							position++
							if buffer[position] != rune('t') {
								goto l54
							}
							position++
							if buffer[position] != rune('e') {
								goto l54
							}
							position++
							goto l52
						l54:
							position, tokenIndex = position52, tokenIndex52
							if buffer[position] != rune('s') {
								goto l55
							}
							position++
							if buffer[position] != rune('t') {
								goto l55
							}
							position++
							if buffer[position] != rune('a') {
								goto l55
							}
							position++
							if buffer[position] != rune('r') {
								goto l55
							}
							position++
							if buffer[position] != rune('t') {
								goto l55
							}
							position++
							goto l52
						l55:
							position, tokenIndex = position52, tokenIndex52
							{
								switch buffer[position] {
								case 'd':
									if buffer[position] != rune('d') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('a') {
										goto l48
									}
									position++
									if buffer[position] != rune('c') {
										goto l48
									}
									position++
									if buffer[position] != rune('h') {
										goto l48
									}
									position++
									break
								case 'c':
									if buffer[position] != rune('c') {
										goto l48
									}
									position++
									if buffer[position] != rune('h') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									if buffer[position] != rune('c') {
										goto l48
									}
									position++
									if buffer[position] != rune('k') {
										goto l48
									}
									position++
									break
								case 'a':
									if buffer[position] != rune('a') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('a') {
										goto l48
									}
									position++
									if buffer[position] != rune('c') {
										goto l48
									}
									position++
									if buffer[position] != rune('h') {
										goto l48
									}
									position++
									break
								case 'u':
									if buffer[position] != rune('u') {
										goto l48
									}
									position++
									if buffer[position] != rune('p') {
										goto l48
									}
									position++
									if buffer[position] != rune('d') {
										goto l48
									}
									position++
									if buffer[position] != rune('a') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									break
								default:
									if buffer[position] != rune('s') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('o') {
										goto l48
									}
									position++
									if buffer[position] != rune('p') {
										goto l48
									}
									position++
									break
								}
							}

						}
					l52:
						add(ruleAction, position51)
					}
					add(rulePegText, position50)
				}
				{
					add(ruleAction1, position)
				}
				if !_rules[ruleMustWhiteSpacing]() {
					goto l48
				}
				{
					position58 := position
					{
						position59 := position
						{
							position60, tokenIndex60 := position, tokenIndex
							if buffer[position] != rune('v') {
								goto l61
							}
							position++
							if buffer[position] != rune('p') {
								goto l61
							}
							position++
							if buffer[position] != rune('c') {
								goto l61
							}
							position++
							goto l60
						l61:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('s') {
								goto l62
							}
							position++
							if buffer[position] != rune('u') {
								goto l62
							}
							position++
							if buffer[position] != rune('b') {
								goto l62
							}
							position++
							if buffer[position] != rune('n') {
								goto l62
							}
							position++
							if buffer[position] != rune('e') {
								goto l62
							}
							position++
							if buffer[position] != rune('t') {
								goto l62
							}
							position++
							goto l60
						l62:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('i') {
								goto l63
							}
							position++
							if buffer[position] != rune('n') {
								goto l63
							}
							position++
							if buffer[position] != rune('s') {
								goto l63
							}
							position++
							if buffer[position] != rune('t') {
								goto l63
							}
							position++
							if buffer[position] != rune('a') {
								goto l63
							}
							position++
							if buffer[position] != rune('n') {
								goto l63
							}
							position++
							if buffer[position] != rune('c') {
								goto l63
							}
							position++
							if buffer[position] != rune('e') {
								goto l63
							}
							position++
							goto l60
						l63:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('r') {
								goto l64
							}
							position++
							if buffer[position] != rune('o') {
								goto l64
							}
							position++
							if buffer[position] != rune('l') {
								goto l64
							}
							position++
							if buffer[position] != rune('e') {
								goto l64
							}
							position++
							goto l60
						l64:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('s') {
								goto l65
							}
							position++
							if buffer[position] != rune('e') {
								goto l65
							}
							position++
							if buffer[position] != rune('c') {
								goto l65
							}
							position++
							if buffer[position] != rune('u') {
								goto l65
							}
							position++
							if buffer[position] != rune('r') {
								goto l65
							}
							position++
							if buffer[position] != rune('i') {
								goto l65
							}
							position++
							if buffer[position] != rune('t') {
								goto l65
							}
							position++
							if buffer[position] != rune('y') {
								goto l65
							}
							position++
							if buffer[position] != rune('g') {
								goto l65
							}
							position++
							if buffer[position] != rune('r') {
								goto l65
							}
							position++
							if buffer[position] != rune('o') {
								goto l65
							}
							position++
							if buffer[position] != rune('u') {
								goto l65
							}
							position++
							if buffer[position] != rune('p') {
								goto l65
							}
							position++
							goto l60
						l65:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('r') {
								goto l66
							}
							position++
							if buffer[position] != rune('o') {
								goto l66
							}
							position++
							if buffer[position] != rune('u') {
								goto l66
							}
							position++
							if buffer[position] != rune('t') {
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
							if buffer[position] != rune('a') {
								goto l66
							}
							position++
							if buffer[position] != rune('b') {
								goto l66
							}
							position++
							if buffer[position] != rune('l') {
								goto l66
							}
							position++
							if buffer[position] != rune('e') {
								goto l66
							}
							position++
							goto l60
						l66:
							position, tokenIndex = position60, tokenIndex60
							{
								switch buffer[position] {
								case 's':
									if buffer[position] != rune('s') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('o') {
										goto l48
									}
									position++
									if buffer[position] != rune('r') {
										goto l48
									}
									position++
									if buffer[position] != rune('a') {
										goto l48
									}
									position++
									if buffer[position] != rune('g') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									if buffer[position] != rune('o') {
										goto l48
									}
									position++
									if buffer[position] != rune('b') {
										goto l48
									}
									position++
									if buffer[position] != rune('j') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									if buffer[position] != rune('c') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									break
								case 'b':
									if buffer[position] != rune('b') {
										goto l48
									}
									position++
									if buffer[position] != rune('u') {
										goto l48
									}
									position++
									if buffer[position] != rune('c') {
										goto l48
									}
									position++
									if buffer[position] != rune('k') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									break
								case 'r':
									if buffer[position] != rune('r') {
										goto l48
									}
									position++
									if buffer[position] != rune('o') {
										goto l48
									}
									position++
									if buffer[position] != rune('u') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									break
								case 'i':
									if buffer[position] != rune('i') {
										goto l48
									}
									position++
									if buffer[position] != rune('n') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									if buffer[position] != rune('r') {
										goto l48
									}
									position++
									if buffer[position] != rune('n') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('g') {
										goto l48
									}
									position++
									if buffer[position] != rune('a') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									if buffer[position] != rune('w') {
										goto l48
									}
									position++
									if buffer[position] != rune('a') {
										goto l48
									}
									position++
									if buffer[position] != rune('y') {
										goto l48
									}
									position++
									break
								case 'k':
									if buffer[position] != rune('k') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									if buffer[position] != rune('y') {
										goto l48
									}
									position++
									if buffer[position] != rune('p') {
										goto l48
									}
									position++
									if buffer[position] != rune('a') {
										goto l48
									}
									position++
									if buffer[position] != rune('i') {
										goto l48
									}
									position++
									if buffer[position] != rune('r') {
										goto l48
									}
									position++
									break
								case 'p':
									if buffer[position] != rune('p') {
										goto l48
									}
									position++
									if buffer[position] != rune('o') {
										goto l48
									}
									position++
									if buffer[position] != rune('l') {
										goto l48
									}
									position++
									if buffer[position] != rune('i') {
										goto l48
									}
									position++
									if buffer[position] != rune('c') {
										goto l48
									}
									position++
									if buffer[position] != rune('y') {
										goto l48
									}
									position++
									break
								case 'g':
									if buffer[position] != rune('g') {
										goto l48
									}
									position++
									if buffer[position] != rune('r') {
										goto l48
									}
									position++
									if buffer[position] != rune('o') {
										goto l48
									}
									position++
									if buffer[position] != rune('u') {
										goto l48
									}
									position++
									if buffer[position] != rune('p') {
										goto l48
									}
									position++
									break
								case 'u':
									if buffer[position] != rune('u') {
										goto l48
									}
									position++
									if buffer[position] != rune('s') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									if buffer[position] != rune('r') {
										goto l48
									}
									position++
									break
								case 't':
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('a') {
										goto l48
									}
									position++
									if buffer[position] != rune('g') {
										goto l48
									}
									position++
									if buffer[position] != rune('s') {
										goto l48
									}
									position++
									break
								default:
									if buffer[position] != rune('v') {
										goto l48
									}
									position++
									if buffer[position] != rune('o') {
										goto l48
									}
									position++
									if buffer[position] != rune('l') {
										goto l48
									}
									position++
									if buffer[position] != rune('u') {
										goto l48
									}
									position++
									if buffer[position] != rune('m') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									break
								}
							}

						}
					l60:
						add(ruleEntity, position59)
					}
					add(rulePegText, position58)
				}
				{
					add(ruleAction2, position)
				}
				{
					position69, tokenIndex69 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l69
					}
					{
						position71 := position
						{
							position74 := position
							{
								position75 := position
								if !_rules[ruleIdentifier]() {
									goto l69
								}
								add(rulePegText, position75)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l69
							}
							{
								position77 := position
								{
									position78, tokenIndex78 := position, tokenIndex
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
											if !matchDot() {
												goto l79
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l79
											}
											position++
										l84:
											{
												position85, tokenIndex85 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l85
												}
												position++
												goto l84
											l85:
												position, tokenIndex = position85, tokenIndex85
											}
											if !matchDot() {
												goto l79
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l79
											}
											position++
										l86:
											{
												position87, tokenIndex87 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l87
												}
												position++
												goto l86
											l87:
												position, tokenIndex = position87, tokenIndex87
											}
											if !matchDot() {
												goto l79
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l79
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
											if buffer[position] != rune('/') {
												goto l79
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l79
											}
											position++
										l90:
											{
												position91, tokenIndex91 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l91
												}
												position++
												goto l90
											l91:
												position, tokenIndex = position91, tokenIndex91
											}
											add(ruleCidrValue, position81)
										}
										add(rulePegText, position80)
									}
									{
										add(ruleAction8, position)
									}
									goto l78
								l79:
									position, tokenIndex = position78, tokenIndex78
									{
										position94 := position
										{
											position95 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l93
											}
											position++
										l96:
											{
												position97, tokenIndex97 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l97
												}
												position++
												goto l96
											l97:
												position, tokenIndex = position97, tokenIndex97
											}
											if !matchDot() {
												goto l93
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l93
											}
											position++
										l98:
											{
												position99, tokenIndex99 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l99
												}
												position++
												goto l98
											l99:
												position, tokenIndex = position99, tokenIndex99
											}
											if !matchDot() {
												goto l93
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l93
											}
											position++
										l100:
											{
												position101, tokenIndex101 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l101
												}
												position++
												goto l100
											l101:
												position, tokenIndex = position101, tokenIndex101
											}
											if !matchDot() {
												goto l93
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l93
											}
											position++
										l102:
											{
												position103, tokenIndex103 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l103
												}
												position++
												goto l102
											l103:
												position, tokenIndex = position103, tokenIndex103
											}
											add(ruleIpValue, position95)
										}
										add(rulePegText, position94)
									}
									{
										add(ruleAction9, position)
									}
									goto l78
								l93:
									position, tokenIndex = position78, tokenIndex78
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
											if buffer[position] != rune('-') {
												goto l105
											}
											position++
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
											add(ruleIntRangeValue, position107)
										}
										add(rulePegText, position106)
									}
									{
										add(ruleAction10, position)
									}
									goto l78
								l105:
									position, tokenIndex = position78, tokenIndex78
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
											add(ruleIntValue, position115)
										}
										add(rulePegText, position114)
									}
									{
										add(ruleAction11, position)
									}
									goto l78
								l113:
									position, tokenIndex = position78, tokenIndex78
									{
										switch buffer[position] {
										case '$':
											{
												position120 := position
												if buffer[position] != rune('$') {
													goto l69
												}
												position++
												{
													position121 := position
													if !_rules[ruleIdentifier]() {
														goto l69
													}
													add(rulePegText, position121)
												}
												add(ruleRefValue, position120)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position123 := position
												if buffer[position] != rune('@') {
													goto l69
												}
												position++
												{
													position124 := position
													if !_rules[ruleIdentifier]() {
														goto l69
													}
													add(rulePegText, position124)
												}
												add(ruleAliasValue, position123)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position126 := position
												if buffer[position] != rune('{') {
													goto l69
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l69
												}
												{
													position127 := position
													if !_rules[ruleIdentifier]() {
														goto l69
													}
													add(rulePegText, position127)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l69
												}
												if buffer[position] != rune('}') {
													goto l69
												}
												position++
												add(ruleHoleValue, position126)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position129 := position
												{
													position130 := position
													{
														switch buffer[position] {
														case '/':
															if buffer[position] != rune('/') {
																goto l69
															}
															position++
															break
														case ':':
															if buffer[position] != rune(':') {
																goto l69
															}
															position++
															break
														case '_':
															if buffer[position] != rune('_') {
																goto l69
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l69
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l69
															}
															position++
															break
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l69
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l69
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l69
															}
															position++
															break
														}
													}

												l131:
													{
														position132, tokenIndex132 := position, tokenIndex
														{
															switch buffer[position] {
															case '/':
																if buffer[position] != rune('/') {
																	goto l132
																}
																position++
																break
															case ':':
																if buffer[position] != rune(':') {
																	goto l132
																}
																position++
																break
															case '_':
																if buffer[position] != rune('_') {
																	goto l132
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l132
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l132
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l132
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l132
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l132
																}
																position++
																break
															}
														}

														goto l131
													l132:
														position, tokenIndex = position132, tokenIndex132
													}
													add(ruleStringValue, position130)
												}
												add(rulePegText, position129)
											}
											{
												add(ruleAction12, position)
											}
											break
										}
									}

								}
							l78:
								add(ruleValue, position77)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l69
							}
							add(ruleParam, position74)
						}
					l72:
						{
							position73, tokenIndex73 := position, tokenIndex
							{
								position136 := position
								{
									position137 := position
									if !_rules[ruleIdentifier]() {
										goto l73
									}
									add(rulePegText, position137)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l73
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
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l167
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
												if buffer[position] != rune('-') {
													goto l167
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l167
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
												add(ruleIntRangeValue, position169)
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
											position176 := position
											{
												position177 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l175
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
												add(ruleIntValue, position177)
											}
											add(rulePegText, position176)
										}
										{
											add(ruleAction11, position)
										}
										goto l140
									l175:
										position, tokenIndex = position140, tokenIndex140
										{
											switch buffer[position] {
											case '$':
												{
													position182 := position
													if buffer[position] != rune('$') {
														goto l73
													}
													position++
													{
														position183 := position
														if !_rules[ruleIdentifier]() {
															goto l73
														}
														add(rulePegText, position183)
													}
													add(ruleRefValue, position182)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position185 := position
													if buffer[position] != rune('@') {
														goto l73
													}
													position++
													{
														position186 := position
														if !_rules[ruleIdentifier]() {
															goto l73
														}
														add(rulePegText, position186)
													}
													add(ruleAliasValue, position185)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position188 := position
													if buffer[position] != rune('{') {
														goto l73
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l73
													}
													{
														position189 := position
														if !_rules[ruleIdentifier]() {
															goto l73
														}
														add(rulePegText, position189)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l73
													}
													if buffer[position] != rune('}') {
														goto l73
													}
													position++
													add(ruleHoleValue, position188)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position191 := position
													{
														position192 := position
														{
															switch buffer[position] {
															case '/':
																if buffer[position] != rune('/') {
																	goto l73
																}
																position++
																break
															case ':':
																if buffer[position] != rune(':') {
																	goto l73
																}
																position++
																break
															case '_':
																if buffer[position] != rune('_') {
																	goto l73
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l73
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l73
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l73
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l73
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l73
																}
																position++
																break
															}
														}

													l193:
														{
															position194, tokenIndex194 := position, tokenIndex
															{
																switch buffer[position] {
																case '/':
																	if buffer[position] != rune('/') {
																		goto l194
																	}
																	position++
																	break
																case ':':
																	if buffer[position] != rune(':') {
																		goto l194
																	}
																	position++
																	break
																case '_':
																	if buffer[position] != rune('_') {
																		goto l194
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
																		goto l194
																	}
																	position++
																	break
																case '-':
																	if buffer[position] != rune('-') {
																		goto l194
																	}
																	position++
																	break
																case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l194
																	}
																	position++
																	break
																case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																	if c := buffer[position]; c < rune('A') || c > rune('Z') {
																		goto l194
																	}
																	position++
																	break
																default:
																	if c := buffer[position]; c < rune('a') || c > rune('z') {
																		goto l194
																	}
																	position++
																	break
																}
															}

															goto l193
														l194:
															position, tokenIndex = position194, tokenIndex194
														}
														add(ruleStringValue, position192)
													}
													add(rulePegText, position191)
												}
												{
													add(ruleAction12, position)
												}
												break
											}
										}

									}
								l140:
									add(ruleValue, position139)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l73
								}
								add(ruleParam, position136)
							}
							goto l72
						l73:
							position, tokenIndex = position73, tokenIndex73
						}
						add(ruleParams, position71)
					}
					goto l70
				l69:
					position, tokenIndex = position69, tokenIndex69
				}
			l70:
				{
					add(ruleAction3, position)
				}
				add(ruleExpr, position49)
			}
			return true
		l48:
			position, tokenIndex = position48, tokenIndex48
			return false
		},
		/* 6 Params <- <Param+> */
		nil,
		/* 7 Param <- <(<Identifier> Action4 Equal Value WhiteSpacing)> */
		nil,
		/* 8 Identifier <- <((&('.') '.') | (&('_') '_') | (&('-') '-') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position201, tokenIndex201 := position, tokenIndex
			{
				position202 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
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
					case '-':
						if buffer[position] != rune('-') {
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
						case '.':
							if buffer[position] != rune('.') {
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
						case '-':
							if buffer[position] != rune('-') {
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
				add(ruleIdentifier, position202)
			}
			return true
		l201:
			position, tokenIndex = position201, tokenIndex201
			return false
		},
		/* 9 Value <- <((<CidrValue> Action8) / (<IpValue> Action9) / (<IntRangeValue> Action10) / (<IntValue> Action11) / ((&('$') (RefValue Action7)) | (&('@') (AliasValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action12))))> */
		nil,
		/* 10 StringValue <- <((&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
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
		/* 18 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)* Action13))> */
		nil,
		/* 19 Spacing <- <Space*> */
		func() bool {
			{
				position218 := position
			l219:
				{
					position220, tokenIndex220 := position, tokenIndex
					{
						position221 := position
						{
							position222, tokenIndex222 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l223
							}
							goto l222
						l223:
							position, tokenIndex = position222, tokenIndex222
							if !_rules[ruleEndOfLine]() {
								goto l220
							}
						}
					l222:
						add(ruleSpace, position221)
					}
					goto l219
				l220:
					position, tokenIndex = position220, tokenIndex220
				}
				add(ruleSpacing, position218)
			}
			return true
		},
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
		/* 22 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position232, tokenIndex232 := position, tokenIndex
			{
				position233 := position
				if !_rules[ruleSpacing]() {
					goto l232
				}
				if buffer[position] != rune('=') {
					goto l232
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l232
				}
				add(ruleEqual, position233)
			}
			return true
		l232:
			position, tokenIndex = position232, tokenIndex232
			return false
		},
		/* 23 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 24 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position235, tokenIndex235 := position, tokenIndex
			{
				position236 := position
				{
					position237, tokenIndex237 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l238
					}
					position++
					goto l237
				l238:
					position, tokenIndex = position237, tokenIndex237
					if buffer[position] != rune('\t') {
						goto l235
					}
					position++
				}
			l237:
				add(ruleWhitespace, position236)
			}
			return true
		l235:
			position, tokenIndex = position235, tokenIndex235
			return false
		},
		/* 25 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position239, tokenIndex239 := position, tokenIndex
			{
				position240 := position
				{
					position241, tokenIndex241 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l242
					}
					position++
					if buffer[position] != rune('\n') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex = position241, tokenIndex241
					if buffer[position] != rune('\n') {
						goto l243
					}
					position++
					goto l241
				l243:
					position, tokenIndex = position241, tokenIndex241
					if buffer[position] != rune('\r') {
						goto l239
					}
					position++
				}
			l241:
				add(ruleEndOfLine, position240)
			}
			return true
		l239:
			position, tokenIndex = position239, tokenIndex239
			return false
		},
		/* 26 EndOfFile <- <!.> */
		nil,
		nil,
		/* 29 Action0 <- <{ p.AddDeclarationIdentifier(text) }> */
		nil,
		/* 30 Action1 <- <{ p.AddAction(text) }> */
		nil,
		/* 31 Action2 <- <{ p.AddEntity(text) }> */
		nil,
		/* 32 Action3 <- <{ p.LineDone() }> */
		nil,
		/* 33 Action4 <- <{ p.AddParamKey(text) }> */
		nil,
		/* 34 Action5 <- <{  p.AddParamHoleValue(text) }> */
		nil,
		/* 35 Action6 <- <{  p.AddParamAliasValue(text) }> */
		nil,
		/* 36 Action7 <- <{  p.AddParamRefValue(text) }> */
		nil,
		/* 37 Action8 <- <{ p.AddParamCidrValue(text) }> */
		nil,
		/* 38 Action9 <- <{ p.AddParamIpValue(text) }> */
		nil,
		/* 39 Action10 <- <{ p.AddParamValue(text) }> */
		nil,
		/* 40 Action11 <- <{ p.AddParamIntValue(text) }> */
		nil,
		/* 41 Action12 <- <{ p.AddParamValue(text) }> */
		nil,
		/* 42 Action13 <- <{ p.LineDone() }> */
		nil,
	}
	p.rules = _rules
}
