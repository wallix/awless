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
			p.addParamValue(text)
		case ruleAction11:
			p.addParamIntValue(text)
		case ruleAction12:
			p.addParamValue(text)
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
		/* 2 Action <- <(('c' 'r' 'e' 'a' 't' 'e') / ('d' 'e' 'l' 'e' 't' 'e') / ('s' 't' 'a' 'r' 't') / ((&('d') ('d' 'e' 't' 'a' 'c' 'h')) | (&('c') ('c' 'h' 'e' 'c' 'k')) | (&('a') ('a' 't' 't' 'a' 'c' 'h')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('s') ('s' 't' 'o' 'p')) | (&('n') ('n' 'o' 'n' 'e'))))> */
		nil,
		/* 3 Entity <- <(('v' 'p' 'c') / ('s' 'u' 'b' 'n' 'e' 't') / ('i' 'n' 's' 't' 'a' 'n' 'c' 'e') / ('t' 'a' 'g') / ('r' 'o' 'l' 'e') / ('s' 'e' 'c' 'u' 'r' 'i' 't' 'y' 'g' 'r' 'o' 'u' 'p') / ('r' 'o' 'u' 't' 'e' 't' 'a' 'b' 'l' 'e') / ('s' 't' 'o' 'r' 'a' 'g' 'e' 'o' 'b' 'j' 'e' 'c' 't') / ((&('q') ('q' 'u' 'e' 'u' 'e')) | (&('t') ('t' 'o' 'p' 'i' 'c')) | (&('s') ('s' 'u' 'b' 's' 'c' 'r' 'i' 'p' 't' 'i' 'o' 'n')) | (&('b') ('b' 'u' 'c' 'k' 'e' 't')) | (&('r') ('r' 'o' 'u' 't' 'e')) | (&('i') ('i' 'n' 't' 'e' 'r' 'n' 'e' 't' 'g' 'a' 't' 'e' 'w' 'a' 'y')) | (&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e')) | (&('n') ('n' 'o' 'n' 'e'))))> */
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
									if buffer[position] != rune('p') {
										goto l48
									}
									position++
									break
								default:
									if buffer[position] != rune('n') {
										goto l48
									}
									position++
									if buffer[position] != rune('o') {
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
							if buffer[position] != rune('t') {
								goto l64
							}
							position++
							if buffer[position] != rune('a') {
								goto l64
							}
							position++
							if buffer[position] != rune('g') {
								goto l64
							}
							position++
							goto l60
						l64:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('r') {
								goto l65
							}
							position++
							if buffer[position] != rune('o') {
								goto l65
							}
							position++
							if buffer[position] != rune('l') {
								goto l65
							}
							position++
							if buffer[position] != rune('e') {
								goto l65
							}
							position++
							goto l60
						l65:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('s') {
								goto l66
							}
							position++
							if buffer[position] != rune('e') {
								goto l66
							}
							position++
							if buffer[position] != rune('c') {
								goto l66
							}
							position++
							if buffer[position] != rune('u') {
								goto l66
							}
							position++
							if buffer[position] != rune('r') {
								goto l66
							}
							position++
							if buffer[position] != rune('i') {
								goto l66
							}
							position++
							if buffer[position] != rune('t') {
								goto l66
							}
							position++
							if buffer[position] != rune('y') {
								goto l66
							}
							position++
							if buffer[position] != rune('g') {
								goto l66
							}
							position++
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
							if buffer[position] != rune('p') {
								goto l66
							}
							position++
							goto l60
						l66:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('r') {
								goto l67
							}
							position++
							if buffer[position] != rune('o') {
								goto l67
							}
							position++
							if buffer[position] != rune('u') {
								goto l67
							}
							position++
							if buffer[position] != rune('t') {
								goto l67
							}
							position++
							if buffer[position] != rune('e') {
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
							if buffer[position] != rune('b') {
								goto l67
							}
							position++
							if buffer[position] != rune('l') {
								goto l67
							}
							position++
							if buffer[position] != rune('e') {
								goto l67
							}
							position++
							goto l60
						l67:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('s') {
								goto l68
							}
							position++
							if buffer[position] != rune('t') {
								goto l68
							}
							position++
							if buffer[position] != rune('o') {
								goto l68
							}
							position++
							if buffer[position] != rune('r') {
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
							if buffer[position] != rune('e') {
								goto l68
							}
							position++
							if buffer[position] != rune('o') {
								goto l68
							}
							position++
							if buffer[position] != rune('b') {
								goto l68
							}
							position++
							if buffer[position] != rune('j') {
								goto l68
							}
							position++
							if buffer[position] != rune('e') {
								goto l68
							}
							position++
							if buffer[position] != rune('c') {
								goto l68
							}
							position++
							if buffer[position] != rune('t') {
								goto l68
							}
							position++
							goto l60
						l68:
							position, tokenIndex = position60, tokenIndex60
							{
								switch buffer[position] {
								case 'q':
									if buffer[position] != rune('q') {
										goto l48
									}
									position++
									if buffer[position] != rune('u') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									if buffer[position] != rune('u') {
										goto l48
									}
									position++
									if buffer[position] != rune('e') {
										goto l48
									}
									position++
									break
								case 't':
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
									if buffer[position] != rune('i') {
										goto l48
									}
									position++
									if buffer[position] != rune('c') {
										goto l48
									}
									position++
									break
								case 's':
									if buffer[position] != rune('s') {
										goto l48
									}
									position++
									if buffer[position] != rune('u') {
										goto l48
									}
									position++
									if buffer[position] != rune('b') {
										goto l48
									}
									position++
									if buffer[position] != rune('s') {
										goto l48
									}
									position++
									if buffer[position] != rune('c') {
										goto l48
									}
									position++
									if buffer[position] != rune('r') {
										goto l48
									}
									position++
									if buffer[position] != rune('i') {
										goto l48
									}
									position++
									if buffer[position] != rune('p') {
										goto l48
									}
									position++
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('i') {
										goto l48
									}
									position++
									if buffer[position] != rune('o') {
										goto l48
									}
									position++
									if buffer[position] != rune('n') {
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
								case 'v':
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
								default:
									if buffer[position] != rune('n') {
										goto l48
									}
									position++
									if buffer[position] != rune('o') {
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
					position71, tokenIndex71 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l71
					}
					{
						position73 := position
						{
							position76 := position
							{
								position77 := position
								if !_rules[ruleIdentifier]() {
									goto l71
								}
								add(rulePegText, position77)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l71
							}
							{
								position79 := position
								{
									position80, tokenIndex80 := position, tokenIndex
									{
										position82 := position
										{
											position83 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l81
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
												goto l81
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l81
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
												goto l81
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l81
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
											if !matchDot() {
												goto l81
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l81
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
											if buffer[position] != rune('/') {
												goto l81
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l81
											}
											position++
										l92:
											{
												position93, tokenIndex93 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l93
												}
												position++
												goto l92
											l93:
												position, tokenIndex = position93, tokenIndex93
											}
											add(ruleCidrValue, position83)
										}
										add(rulePegText, position82)
									}
									{
										add(ruleAction8, position)
									}
									goto l80
								l81:
									position, tokenIndex = position80, tokenIndex80
									{
										position96 := position
										{
											position97 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l95
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
												goto l95
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l95
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
												goto l95
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l95
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
											if !matchDot() {
												goto l95
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l95
											}
											position++
										l104:
											{
												position105, tokenIndex105 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l105
												}
												position++
												goto l104
											l105:
												position, tokenIndex = position105, tokenIndex105
											}
											add(ruleIpValue, position97)
										}
										add(rulePegText, position96)
									}
									{
										add(ruleAction9, position)
									}
									goto l80
								l95:
									position, tokenIndex = position80, tokenIndex80
									{
										position108 := position
										{
											position109 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l107
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
											if buffer[position] != rune('-') {
												goto l107
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l107
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
											add(ruleIntRangeValue, position109)
										}
										add(rulePegText, position108)
									}
									{
										add(ruleAction10, position)
									}
									goto l80
								l107:
									position, tokenIndex = position80, tokenIndex80
									{
										position116 := position
										{
											position117 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l115
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
											add(ruleIntValue, position117)
										}
										add(rulePegText, position116)
									}
									{
										add(ruleAction11, position)
									}
									goto l80
								l115:
									position, tokenIndex = position80, tokenIndex80
									{
										switch buffer[position] {
										case '$':
											{
												position122 := position
												if buffer[position] != rune('$') {
													goto l71
												}
												position++
												{
													position123 := position
													if !_rules[ruleIdentifier]() {
														goto l71
													}
													add(rulePegText, position123)
												}
												add(ruleRefValue, position122)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position125 := position
												{
													position126 := position
													if buffer[position] != rune('@') {
														goto l71
													}
													position++
													if !_rules[ruleStringValue]() {
														goto l71
													}
													add(rulePegText, position126)
												}
												add(ruleAliasValue, position125)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position128 := position
												if buffer[position] != rune('{') {
													goto l71
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l71
												}
												{
													position129 := position
													if !_rules[ruleIdentifier]() {
														goto l71
													}
													add(rulePegText, position129)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l71
												}
												if buffer[position] != rune('}') {
													goto l71
												}
												position++
												add(ruleHoleValue, position128)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position131 := position
												if !_rules[ruleStringValue]() {
													goto l71
												}
												add(rulePegText, position131)
											}
											{
												add(ruleAction12, position)
											}
											break
										}
									}

								}
							l80:
								add(ruleValue, position79)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l71
							}
							add(ruleParam, position76)
						}
					l74:
						{
							position75, tokenIndex75 := position, tokenIndex
							{
								position133 := position
								{
									position134 := position
									if !_rules[ruleIdentifier]() {
										goto l75
									}
									add(rulePegText, position134)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l75
								}
								{
									position136 := position
									{
										position137, tokenIndex137 := position, tokenIndex
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
												if !matchDot() {
													goto l138
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l138
												}
												position++
											l143:
												{
													position144, tokenIndex144 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l144
													}
													position++
													goto l143
												l144:
													position, tokenIndex = position144, tokenIndex144
												}
												if !matchDot() {
													goto l138
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l138
												}
												position++
											l145:
												{
													position146, tokenIndex146 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l146
													}
													position++
													goto l145
												l146:
													position, tokenIndex = position146, tokenIndex146
												}
												if !matchDot() {
													goto l138
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l138
												}
												position++
											l147:
												{
													position148, tokenIndex148 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l148
													}
													position++
													goto l147
												l148:
													position, tokenIndex = position148, tokenIndex148
												}
												if buffer[position] != rune('/') {
													goto l138
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l138
												}
												position++
											l149:
												{
													position150, tokenIndex150 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l150
													}
													position++
													goto l149
												l150:
													position, tokenIndex = position150, tokenIndex150
												}
												add(ruleCidrValue, position140)
											}
											add(rulePegText, position139)
										}
										{
											add(ruleAction8, position)
										}
										goto l137
									l138:
										position, tokenIndex = position137, tokenIndex137
										{
											position153 := position
											{
												position154 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l152
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
													goto l152
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l152
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
													goto l152
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l152
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
												if !matchDot() {
													goto l152
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l152
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
												add(ruleIpValue, position154)
											}
											add(rulePegText, position153)
										}
										{
											add(ruleAction9, position)
										}
										goto l137
									l152:
										position, tokenIndex = position137, tokenIndex137
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
												if buffer[position] != rune('-') {
													goto l164
												}
												position++
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
												add(ruleIntRangeValue, position166)
											}
											add(rulePegText, position165)
										}
										{
											add(ruleAction10, position)
										}
										goto l137
									l164:
										position, tokenIndex = position137, tokenIndex137
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
												add(ruleIntValue, position174)
											}
											add(rulePegText, position173)
										}
										{
											add(ruleAction11, position)
										}
										goto l137
									l172:
										position, tokenIndex = position137, tokenIndex137
										{
											switch buffer[position] {
											case '$':
												{
													position179 := position
													if buffer[position] != rune('$') {
														goto l75
													}
													position++
													{
														position180 := position
														if !_rules[ruleIdentifier]() {
															goto l75
														}
														add(rulePegText, position180)
													}
													add(ruleRefValue, position179)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position182 := position
													{
														position183 := position
														if buffer[position] != rune('@') {
															goto l75
														}
														position++
														if !_rules[ruleStringValue]() {
															goto l75
														}
														add(rulePegText, position183)
													}
													add(ruleAliasValue, position182)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position185 := position
													if buffer[position] != rune('{') {
														goto l75
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l75
													}
													{
														position186 := position
														if !_rules[ruleIdentifier]() {
															goto l75
														}
														add(rulePegText, position186)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l75
													}
													if buffer[position] != rune('}') {
														goto l75
													}
													position++
													add(ruleHoleValue, position185)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position188 := position
													if !_rules[ruleStringValue]() {
														goto l75
													}
													add(rulePegText, position188)
												}
												{
													add(ruleAction12, position)
												}
												break
											}
										}

									}
								l137:
									add(ruleValue, position136)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l75
								}
								add(ruleParam, position133)
							}
							goto l74
						l75:
							position, tokenIndex = position75, tokenIndex75
						}
						add(ruleParams, position73)
					}
					goto l72
				l71:
					position, tokenIndex = position71, tokenIndex71
				}
			l72:
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
		/* 8 Identifier <- <((&('.') '.') | (&('_') '_') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position193, tokenIndex193 := position, tokenIndex
			{
				position194 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l193
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l193
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l193
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l193
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l193
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l193
						}
						position++
						break
					}
				}

			l195:
				{
					position196, tokenIndex196 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l196
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l196
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l196
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l196
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l196
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l196
							}
							position++
							break
						}
					}

					goto l195
				l196:
					position, tokenIndex = position196, tokenIndex196
				}
				add(ruleIdentifier, position194)
			}
			return true
		l193:
			position, tokenIndex = position193, tokenIndex193
			return false
		},
		/* 9 Value <- <((<CidrValue> Action8) / (<IpValue> Action9) / (<IntRangeValue> Action10) / (<IntValue> Action11) / ((&('$') (RefValue Action7)) | (&('@') (AliasValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action12))))> */
		nil,
		/* 10 StringValue <- <((&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position200, tokenIndex200 := position, tokenIndex
			{
				position201 := position
				{
					switch buffer[position] {
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

			l202:
				{
					position203, tokenIndex203 := position, tokenIndex
					{
						switch buffer[position] {
						case '/':
							if buffer[position] != rune('/') {
								goto l203
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l203
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l203
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l203
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l203
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l203
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l203
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l203
							}
							position++
							break
						}
					}

					goto l202
				l203:
					position, tokenIndex = position203, tokenIndex203
				}
				add(ruleStringValue, position201)
			}
			return true
		l200:
			position, tokenIndex = position200, tokenIndex200
			return false
		},
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
		/* 16 AliasValue <- <<('@' StringValue)>> */
		nil,
		/* 17 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 18 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)* Action13))> */
		nil,
		/* 19 Spacing <- <Space*> */
		func() bool {
			{
				position215 := position
			l216:
				{
					position217, tokenIndex217 := position, tokenIndex
					{
						position218 := position
						{
							position219, tokenIndex219 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l220
							}
							goto l219
						l220:
							position, tokenIndex = position219, tokenIndex219
							if !_rules[ruleEndOfLine]() {
								goto l217
							}
						}
					l219:
						add(ruleSpace, position218)
					}
					goto l216
				l217:
					position, tokenIndex = position217, tokenIndex217
				}
				add(ruleSpacing, position215)
			}
			return true
		},
		/* 20 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position222 := position
			l223:
				{
					position224, tokenIndex224 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l224
					}
					goto l223
				l224:
					position, tokenIndex = position224, tokenIndex224
				}
				add(ruleWhiteSpacing, position222)
			}
			return true
		},
		/* 21 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position225, tokenIndex225 := position, tokenIndex
			{
				position226 := position
				if !_rules[ruleWhitespace]() {
					goto l225
				}
			l227:
				{
					position228, tokenIndex228 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l228
					}
					goto l227
				l228:
					position, tokenIndex = position228, tokenIndex228
				}
				add(ruleMustWhiteSpacing, position226)
			}
			return true
		l225:
			position, tokenIndex = position225, tokenIndex225
			return false
		},
		/* 22 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position229, tokenIndex229 := position, tokenIndex
			{
				position230 := position
				if !_rules[ruleSpacing]() {
					goto l229
				}
				if buffer[position] != rune('=') {
					goto l229
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l229
				}
				add(ruleEqual, position230)
			}
			return true
		l229:
			position, tokenIndex = position229, tokenIndex229
			return false
		},
		/* 23 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 24 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position232, tokenIndex232 := position, tokenIndex
			{
				position233 := position
				{
					position234, tokenIndex234 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex = position234, tokenIndex234
					if buffer[position] != rune('\t') {
						goto l232
					}
					position++
				}
			l234:
				add(ruleWhitespace, position233)
			}
			return true
		l232:
			position, tokenIndex = position232, tokenIndex232
			return false
		},
		/* 25 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position236, tokenIndex236 := position, tokenIndex
			{
				position237 := position
				{
					position238, tokenIndex238 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l239
					}
					position++
					if buffer[position] != rune('\n') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex = position238, tokenIndex238
					if buffer[position] != rune('\n') {
						goto l240
					}
					position++
					goto l238
				l240:
					position, tokenIndex = position238, tokenIndex238
					if buffer[position] != rune('\r') {
						goto l236
					}
					position++
				}
			l238:
				add(ruleEndOfLine, position237)
			}
			return true
		l236:
			position, tokenIndex = position236, tokenIndex236
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
		/* 39 Action10 <- <{ p.addParamValue(text) }> */
		nil,
		/* 40 Action11 <- <{ p.addParamIntValue(text) }> */
		nil,
		/* 41 Action12 <- <{ p.addParamValue(text) }> */
		nil,
		/* 42 Action13 <- <{ p.LineDone() }> */
		nil,
	}
	p.rules = _rules
}
