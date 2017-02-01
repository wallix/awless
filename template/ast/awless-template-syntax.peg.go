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
		/* 3 Entity <- <(('v' 'p' 'c') / ('s' 'u' 'b' 'n' 'e' 't') / ('i' 'n' 's' 't' 'a' 'n' 'c' 'e') / ('r' 'o' 'l' 'e') / ('r' 'o' 'u' 't' 'e' 't' 'a' 'b' 'l' 'e') / ((&('r') ('r' 'o' 'u' 't' 'e')) | (&('i') ('i' 'n' 't' 'e' 'r' 'n' 'e' 't' 'g' 'a' 't' 'e' 'w' 'a' 'y')) | (&('s') ('s' 'e' 'c' 'u' 'r' 'i' 't' 'y' 'g' 'r' 'o' 'u' 'p')) | (&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('t') ('t' 'a' 'g' 's')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e'))))> */
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
							if buffer[position] != rune('t') {
								goto l65
							}
							position++
							if buffer[position] != rune('e') {
								goto l65
							}
							position++
							if buffer[position] != rune('t') {
								goto l65
							}
							position++
							if buffer[position] != rune('a') {
								goto l65
							}
							position++
							if buffer[position] != rune('b') {
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
							{
								switch buffer[position] {
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
								case 's':
									if buffer[position] != rune('s') {
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
									if buffer[position] != rune('u') {
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
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('y') {
										goto l48
									}
									position++
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
											if buffer[position] != rune('-') {
												goto l104
											}
											position++
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
											add(ruleIntRangeValue, position106)
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
										position113 := position
										{
											position114 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l112
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
											add(ruleIntValue, position114)
										}
										add(rulePegText, position113)
									}
									{
										add(ruleAction11, position)
									}
									goto l77
								l112:
									position, tokenIndex = position77, tokenIndex77
									{
										switch buffer[position] {
										case '$':
											{
												position119 := position
												if buffer[position] != rune('$') {
													goto l68
												}
												position++
												{
													position120 := position
													if !_rules[ruleIdentifier]() {
														goto l68
													}
													add(rulePegText, position120)
												}
												add(ruleRefValue, position119)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position122 := position
												if buffer[position] != rune('@') {
													goto l68
												}
												position++
												{
													position123 := position
													if !_rules[ruleIdentifier]() {
														goto l68
													}
													add(rulePegText, position123)
												}
												add(ruleAliasValue, position122)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position125 := position
												if buffer[position] != rune('{') {
													goto l68
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l68
												}
												{
													position126 := position
													if !_rules[ruleIdentifier]() {
														goto l68
													}
													add(rulePegText, position126)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l68
												}
												if buffer[position] != rune('}') {
													goto l68
												}
												position++
												add(ruleHoleValue, position125)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position128 := position
												{
													position129 := position
													{
														switch buffer[position] {
														case '_':
															if buffer[position] != rune('_') {
																goto l68
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l68
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l68
															}
															position++
															break
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l68
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l68
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l68
															}
															position++
															break
														}
													}

												l130:
													{
														position131, tokenIndex131 := position, tokenIndex
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l131
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l131
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l131
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l131
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l131
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l131
																}
																position++
																break
															}
														}

														goto l130
													l131:
														position, tokenIndex = position131, tokenIndex131
													}
													add(ruleStringValue, position129)
												}
												add(rulePegText, position128)
											}
											{
												add(ruleAction12, position)
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
								position135 := position
								{
									position136 := position
									if !_rules[ruleIdentifier]() {
										goto l72
									}
									add(rulePegText, position136)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l72
								}
								{
									position138 := position
									{
										position139, tokenIndex139 := position, tokenIndex
										{
											position141 := position
											{
												position142 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l140
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
													goto l140
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l140
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
													goto l140
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l140
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
												if !matchDot() {
													goto l140
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l140
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
												if buffer[position] != rune('/') {
													goto l140
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l140
												}
												position++
											l151:
												{
													position152, tokenIndex152 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l152
													}
													position++
													goto l151
												l152:
													position, tokenIndex = position152, tokenIndex152
												}
												add(ruleCidrValue, position142)
											}
											add(rulePegText, position141)
										}
										{
											add(ruleAction8, position)
										}
										goto l139
									l140:
										position, tokenIndex = position139, tokenIndex139
										{
											position155 := position
											{
												position156 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l154
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
													goto l154
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l154
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
													goto l154
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l154
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
												if !matchDot() {
													goto l154
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l154
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
												add(ruleIpValue, position156)
											}
											add(rulePegText, position155)
										}
										{
											add(ruleAction9, position)
										}
										goto l139
									l154:
										position, tokenIndex = position139, tokenIndex139
										{
											position167 := position
											{
												position168 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l166
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
												if buffer[position] != rune('-') {
													goto l166
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l166
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
												add(ruleIntRangeValue, position168)
											}
											add(rulePegText, position167)
										}
										{
											add(ruleAction10, position)
										}
										goto l139
									l166:
										position, tokenIndex = position139, tokenIndex139
										{
											position175 := position
											{
												position176 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l174
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
												add(ruleIntValue, position176)
											}
											add(rulePegText, position175)
										}
										{
											add(ruleAction11, position)
										}
										goto l139
									l174:
										position, tokenIndex = position139, tokenIndex139
										{
											switch buffer[position] {
											case '$':
												{
													position181 := position
													if buffer[position] != rune('$') {
														goto l72
													}
													position++
													{
														position182 := position
														if !_rules[ruleIdentifier]() {
															goto l72
														}
														add(rulePegText, position182)
													}
													add(ruleRefValue, position181)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position184 := position
													if buffer[position] != rune('@') {
														goto l72
													}
													position++
													{
														position185 := position
														if !_rules[ruleIdentifier]() {
															goto l72
														}
														add(rulePegText, position185)
													}
													add(ruleAliasValue, position184)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position187 := position
													if buffer[position] != rune('{') {
														goto l72
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l72
													}
													{
														position188 := position
														if !_rules[ruleIdentifier]() {
															goto l72
														}
														add(rulePegText, position188)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l72
													}
													if buffer[position] != rune('}') {
														goto l72
													}
													position++
													add(ruleHoleValue, position187)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position190 := position
													{
														position191 := position
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l72
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l72
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l72
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l72
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l72
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l72
																}
																position++
																break
															}
														}

													l192:
														{
															position193, tokenIndex193 := position, tokenIndex
															{
																switch buffer[position] {
																case '_':
																	if buffer[position] != rune('_') {
																		goto l193
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
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

															goto l192
														l193:
															position, tokenIndex = position193, tokenIndex193
														}
														add(ruleStringValue, position191)
													}
													add(rulePegText, position190)
												}
												{
													add(ruleAction12, position)
												}
												break
											}
										}

									}
								l139:
									add(ruleValue, position138)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l72
								}
								add(ruleParam, position135)
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
			position200, tokenIndex200 := position, tokenIndex
			{
				position201 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
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
					case '-':
						if buffer[position] != rune('-') {
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
						case '.':
							if buffer[position] != rune('.') {
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
						case '-':
							if buffer[position] != rune('-') {
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
				add(ruleIdentifier, position201)
			}
			return true
		l200:
			position, tokenIndex = position200, tokenIndex200
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
		/* 18 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)* Action13))> */
		nil,
		/* 19 Spacing <- <Space*> */
		func() bool {
			{
				position217 := position
			l218:
				{
					position219, tokenIndex219 := position, tokenIndex
					{
						position220 := position
						{
							position221, tokenIndex221 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l222
							}
							goto l221
						l222:
							position, tokenIndex = position221, tokenIndex221
							if !_rules[ruleEndOfLine]() {
								goto l219
							}
						}
					l221:
						add(ruleSpace, position220)
					}
					goto l218
				l219:
					position, tokenIndex = position219, tokenIndex219
				}
				add(ruleSpacing, position217)
			}
			return true
		},
		/* 20 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position224 := position
			l225:
				{
					position226, tokenIndex226 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l226
					}
					goto l225
				l226:
					position, tokenIndex = position226, tokenIndex226
				}
				add(ruleWhiteSpacing, position224)
			}
			return true
		},
		/* 21 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position227, tokenIndex227 := position, tokenIndex
			{
				position228 := position
				if !_rules[ruleWhitespace]() {
					goto l227
				}
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
				add(ruleMustWhiteSpacing, position228)
			}
			return true
		l227:
			position, tokenIndex = position227, tokenIndex227
			return false
		},
		/* 22 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position231, tokenIndex231 := position, tokenIndex
			{
				position232 := position
				if !_rules[ruleSpacing]() {
					goto l231
				}
				if buffer[position] != rune('=') {
					goto l231
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l231
				}
				add(ruleEqual, position232)
			}
			return true
		l231:
			position, tokenIndex = position231, tokenIndex231
			return false
		},
		/* 23 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 24 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position234, tokenIndex234 := position, tokenIndex
			{
				position235 := position
				{
					position236, tokenIndex236 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex = position236, tokenIndex236
					if buffer[position] != rune('\t') {
						goto l234
					}
					position++
				}
			l236:
				add(ruleWhitespace, position235)
			}
			return true
		l234:
			position, tokenIndex = position234, tokenIndex234
			return false
		},
		/* 25 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position238, tokenIndex238 := position, tokenIndex
			{
				position239 := position
				{
					position240, tokenIndex240 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l241
					}
					position++
					if buffer[position] != rune('\n') {
						goto l241
					}
					position++
					goto l240
				l241:
					position, tokenIndex = position240, tokenIndex240
					if buffer[position] != rune('\n') {
						goto l242
					}
					position++
					goto l240
				l242:
					position, tokenIndex = position240, tokenIndex240
					if buffer[position] != rune('\r') {
						goto l238
					}
					position++
				}
			l240:
				add(ruleEndOfLine, position239)
			}
			return true
		l238:
			position, tokenIndex = position238, tokenIndex238
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
