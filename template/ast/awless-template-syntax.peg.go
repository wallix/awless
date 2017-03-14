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
	ruleAction14
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
	"Action14",
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
									add(ruleAction14, position)
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
										add(ruleAction14, position)
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
		/* 3 Entity <- <(('v' 'p' 'c') / ('s' 'u' 'b' 'n' 'e' 't') / ('i' 'n' 's' 't' 'a' 'n' 'c' 'e') / ('t' 'a' 'g') / ('s' 'e' 'c' 'u' 'r' 'i' 't' 'y' 'g' 'r' 'o' 'u' 'p') / ('r' 'o' 'u' 't' 'e' 't' 'a' 'b' 'l' 'e') / ('r' 'o' 'u' 't' 'e') / ('r' 'o' 'l' 'e') / ('s' 't' 'o' 'r' 'a' 'g' 'e' 'o' 'b' 'j' 'e' 'c' 't') / ('t' 'o' 'p' 'i' 'c') / ('l' 'o' 'a' 'd' 'b' 'a' 'l' 'a' 'n' 'c' 'e' 'r') / ((&('r') ('r' 'e' 'c' 'o' 'r' 'd')) | (&('z') ('z' 'o' 'n' 'e')) | (&('t') ('t' 'a' 'r' 'g' 'e' 't' 'g' 'r' 'o' 'u' 'p')) | (&('l') ('l' 'i' 's' 't' 'e' 'n' 'e' 'r')) | (&('q') ('q' 'u' 'e' 'u' 'e')) | (&('s') ('s' 'u' 'b' 's' 'c' 'r' 'i' 'p' 't' 'i' 'o' 'n')) | (&('b') ('b' 'u' 'c' 'k' 'e' 't')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('i') ('i' 'n' 't' 'e' 'r' 'n' 'e' 't' 'g' 'a' 't' 'e' 'w' 'a' 'y')) | (&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e')) | (&('n') ('n' 'o' 'n' 'e'))))> */
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
							goto l60
						l67:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('r') {
								goto l68
							}
							position++
							if buffer[position] != rune('o') {
								goto l68
							}
							position++
							if buffer[position] != rune('l') {
								goto l68
							}
							position++
							if buffer[position] != rune('e') {
								goto l68
							}
							position++
							goto l60
						l68:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('s') {
								goto l69
							}
							position++
							if buffer[position] != rune('t') {
								goto l69
							}
							position++
							if buffer[position] != rune('o') {
								goto l69
							}
							position++
							if buffer[position] != rune('r') {
								goto l69
							}
							position++
							if buffer[position] != rune('a') {
								goto l69
							}
							position++
							if buffer[position] != rune('g') {
								goto l69
							}
							position++
							if buffer[position] != rune('e') {
								goto l69
							}
							position++
							if buffer[position] != rune('o') {
								goto l69
							}
							position++
							if buffer[position] != rune('b') {
								goto l69
							}
							position++
							if buffer[position] != rune('j') {
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
							if buffer[position] != rune('t') {
								goto l69
							}
							position++
							goto l60
						l69:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('t') {
								goto l70
							}
							position++
							if buffer[position] != rune('o') {
								goto l70
							}
							position++
							if buffer[position] != rune('p') {
								goto l70
							}
							position++
							if buffer[position] != rune('i') {
								goto l70
							}
							position++
							if buffer[position] != rune('c') {
								goto l70
							}
							position++
							goto l60
						l70:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('l') {
								goto l71
							}
							position++
							if buffer[position] != rune('o') {
								goto l71
							}
							position++
							if buffer[position] != rune('a') {
								goto l71
							}
							position++
							if buffer[position] != rune('d') {
								goto l71
							}
							position++
							if buffer[position] != rune('b') {
								goto l71
							}
							position++
							if buffer[position] != rune('a') {
								goto l71
							}
							position++
							if buffer[position] != rune('l') {
								goto l71
							}
							position++
							if buffer[position] != rune('a') {
								goto l71
							}
							position++
							if buffer[position] != rune('n') {
								goto l71
							}
							position++
							if buffer[position] != rune('c') {
								goto l71
							}
							position++
							if buffer[position] != rune('e') {
								goto l71
							}
							position++
							if buffer[position] != rune('r') {
								goto l71
							}
							position++
							goto l60
						l71:
							position, tokenIndex = position60, tokenIndex60
							{
								switch buffer[position] {
								case 'r':
									if buffer[position] != rune('r') {
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
									if buffer[position] != rune('o') {
										goto l48
									}
									position++
									if buffer[position] != rune('r') {
										goto l48
									}
									position++
									if buffer[position] != rune('d') {
										goto l48
									}
									position++
									break
								case 'z':
									if buffer[position] != rune('z') {
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
								case 't':
									if buffer[position] != rune('t') {
										goto l48
									}
									position++
									if buffer[position] != rune('a') {
										goto l48
									}
									position++
									if buffer[position] != rune('r') {
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
									if buffer[position] != rune('t') {
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
								case 'l':
									if buffer[position] != rune('l') {
										goto l48
									}
									position++
									if buffer[position] != rune('i') {
										goto l48
									}
									position++
									if buffer[position] != rune('s') {
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
									if buffer[position] != rune('n') {
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
					position74, tokenIndex74 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l74
					}
					{
						position76 := position
						{
							position79 := position
							{
								position80 := position
								if !_rules[ruleIdentifier]() {
									goto l74
								}
								add(rulePegText, position80)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l74
							}
							{
								position82 := position
								{
									position83, tokenIndex83 := position, tokenIndex
									{
										position85 := position
										{
											position86 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l84
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
											if !matchDot() {
												goto l84
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l84
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
											if !matchDot() {
												goto l84
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l84
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
												goto l84
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l84
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
											if buffer[position] != rune('/') {
												goto l84
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l84
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
											add(ruleCidrValue, position86)
										}
										add(rulePegText, position85)
									}
									{
										add(ruleAction8, position)
									}
									goto l83
								l84:
									position, tokenIndex = position83, tokenIndex83
									{
										position99 := position
										{
											position100 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l98
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
											if !matchDot() {
												goto l98
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l98
											}
											position++
										l103:
											{
												position104, tokenIndex104 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l104
												}
												position++
												goto l103
											l104:
												position, tokenIndex = position104, tokenIndex104
											}
											if !matchDot() {
												goto l98
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l98
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
												goto l98
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l98
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
											add(ruleIpValue, position100)
										}
										add(rulePegText, position99)
									}
									{
										add(ruleAction9, position)
									}
									goto l83
								l98:
									position, tokenIndex = position83, tokenIndex83
									{
										position111 := position
										{
											position112 := position
											if !_rules[ruleStringValue]() {
												goto l110
											}
											if !_rules[ruleWhiteSpacing]() {
												goto l110
											}
											if buffer[position] != rune(',') {
												goto l110
											}
											position++
											if !_rules[ruleWhiteSpacing]() {
												goto l110
											}
										l113:
											{
												position114, tokenIndex114 := position, tokenIndex
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
												goto l113
											l114:
												position, tokenIndex = position114, tokenIndex114
											}
											if !_rules[ruleStringValue]() {
												goto l110
											}
											add(ruleCSVValue, position112)
										}
										add(rulePegText, position111)
									}
									{
										add(ruleAction10, position)
									}
									goto l83
								l110:
									position, tokenIndex = position83, tokenIndex83
									{
										position117 := position
										{
											position118 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l116
											}
											position++
										l119:
											{
												position120, tokenIndex120 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l120
												}
												position++
												goto l119
											l120:
												position, tokenIndex = position120, tokenIndex120
											}
											if buffer[position] != rune('-') {
												goto l116
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l116
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
											add(ruleIntRangeValue, position118)
										}
										add(rulePegText, position117)
									}
									{
										add(ruleAction11, position)
									}
									goto l83
								l116:
									position, tokenIndex = position83, tokenIndex83
									{
										position125 := position
										{
											position126 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l124
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
											add(ruleIntValue, position126)
										}
										add(rulePegText, position125)
									}
									{
										add(ruleAction12, position)
									}
									goto l83
								l124:
									position, tokenIndex = position83, tokenIndex83
									{
										switch buffer[position] {
										case '$':
											{
												position131 := position
												if buffer[position] != rune('$') {
													goto l74
												}
												position++
												{
													position132 := position
													if !_rules[ruleIdentifier]() {
														goto l74
													}
													add(rulePegText, position132)
												}
												add(ruleRefValue, position131)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position134 := position
												{
													position135 := position
													if buffer[position] != rune('@') {
														goto l74
													}
													position++
													if !_rules[ruleStringValue]() {
														goto l74
													}
													add(rulePegText, position135)
												}
												add(ruleAliasValue, position134)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position137 := position
												if buffer[position] != rune('{') {
													goto l74
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l74
												}
												{
													position138 := position
													if !_rules[ruleIdentifier]() {
														goto l74
													}
													add(rulePegText, position138)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l74
												}
												if buffer[position] != rune('}') {
													goto l74
												}
												position++
												add(ruleHoleValue, position137)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position140 := position
												if !_rules[ruleStringValue]() {
													goto l74
												}
												add(rulePegText, position140)
											}
											{
												add(ruleAction13, position)
											}
											break
										}
									}

								}
							l83:
								add(ruleValue, position82)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l74
							}
							add(ruleParam, position79)
						}
					l77:
						{
							position78, tokenIndex78 := position, tokenIndex
							{
								position142 := position
								{
									position143 := position
									if !_rules[ruleIdentifier]() {
										goto l78
									}
									add(rulePegText, position143)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l78
								}
								{
									position145 := position
									{
										position146, tokenIndex146 := position, tokenIndex
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
												if !matchDot() {
													goto l147
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l147
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
												if !matchDot() {
													goto l147
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l147
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
													goto l147
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l147
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
												if buffer[position] != rune('/') {
													goto l147
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l147
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
												add(ruleCidrValue, position149)
											}
											add(rulePegText, position148)
										}
										{
											add(ruleAction8, position)
										}
										goto l146
									l147:
										position, tokenIndex = position146, tokenIndex146
										{
											position162 := position
											{
												position163 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l161
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
												if !matchDot() {
													goto l161
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l161
												}
												position++
											l166:
												{
													position167, tokenIndex167 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l167
													}
													position++
													goto l166
												l167:
													position, tokenIndex = position167, tokenIndex167
												}
												if !matchDot() {
													goto l161
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l161
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
													goto l161
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l161
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
												add(ruleIpValue, position163)
											}
											add(rulePegText, position162)
										}
										{
											add(ruleAction9, position)
										}
										goto l146
									l161:
										position, tokenIndex = position146, tokenIndex146
										{
											position174 := position
											{
												position175 := position
												if !_rules[ruleStringValue]() {
													goto l173
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l173
												}
												if buffer[position] != rune(',') {
													goto l173
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l173
												}
											l176:
												{
													position177, tokenIndex177 := position, tokenIndex
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
													goto l176
												l177:
													position, tokenIndex = position177, tokenIndex177
												}
												if !_rules[ruleStringValue]() {
													goto l173
												}
												add(ruleCSVValue, position175)
											}
											add(rulePegText, position174)
										}
										{
											add(ruleAction10, position)
										}
										goto l146
									l173:
										position, tokenIndex = position146, tokenIndex146
										{
											position180 := position
											{
												position181 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l179
												}
												position++
											l182:
												{
													position183, tokenIndex183 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l183
													}
													position++
													goto l182
												l183:
													position, tokenIndex = position183, tokenIndex183
												}
												if buffer[position] != rune('-') {
													goto l179
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l179
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
												add(ruleIntRangeValue, position181)
											}
											add(rulePegText, position180)
										}
										{
											add(ruleAction11, position)
										}
										goto l146
									l179:
										position, tokenIndex = position146, tokenIndex146
										{
											position188 := position
											{
												position189 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l187
												}
												position++
											l190:
												{
													position191, tokenIndex191 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l191
													}
													position++
													goto l190
												l191:
													position, tokenIndex = position191, tokenIndex191
												}
												add(ruleIntValue, position189)
											}
											add(rulePegText, position188)
										}
										{
											add(ruleAction12, position)
										}
										goto l146
									l187:
										position, tokenIndex = position146, tokenIndex146
										{
											switch buffer[position] {
											case '$':
												{
													position194 := position
													if buffer[position] != rune('$') {
														goto l78
													}
													position++
													{
														position195 := position
														if !_rules[ruleIdentifier]() {
															goto l78
														}
														add(rulePegText, position195)
													}
													add(ruleRefValue, position194)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position197 := position
													{
														position198 := position
														if buffer[position] != rune('@') {
															goto l78
														}
														position++
														if !_rules[ruleStringValue]() {
															goto l78
														}
														add(rulePegText, position198)
													}
													add(ruleAliasValue, position197)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position200 := position
													if buffer[position] != rune('{') {
														goto l78
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l78
													}
													{
														position201 := position
														if !_rules[ruleIdentifier]() {
															goto l78
														}
														add(rulePegText, position201)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l78
													}
													if buffer[position] != rune('}') {
														goto l78
													}
													position++
													add(ruleHoleValue, position200)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position203 := position
													if !_rules[ruleStringValue]() {
														goto l78
													}
													add(rulePegText, position203)
												}
												{
													add(ruleAction13, position)
												}
												break
											}
										}

									}
								l146:
									add(ruleValue, position145)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l78
								}
								add(ruleParam, position142)
							}
							goto l77
						l78:
							position, tokenIndex = position78, tokenIndex78
						}
						add(ruleParams, position76)
					}
					goto l75
				l74:
					position, tokenIndex = position74, tokenIndex74
				}
			l75:
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
			position208, tokenIndex208 := position, tokenIndex
			{
				position209 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l208
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l208
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l208
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l208
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l208
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l208
						}
						position++
						break
					}
				}

			l210:
				{
					position211, tokenIndex211 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l211
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l211
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l211
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l211
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l211
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l211
							}
							position++
							break
						}
					}

					goto l210
				l211:
					position, tokenIndex = position211, tokenIndex211
				}
				add(ruleIdentifier, position209)
			}
			return true
		l208:
			position, tokenIndex = position208, tokenIndex208
			return false
		},
		/* 9 Value <- <((<CidrValue> Action8) / (<IpValue> Action9) / (<CSVValue> Action10) / (<IntRangeValue> Action11) / (<IntValue> Action12) / ((&('$') (RefValue Action7)) | (&('@') (AliasValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action13))))> */
		nil,
		/* 10 StringValue <- <((&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position215, tokenIndex215 := position, tokenIndex
			{
				position216 := position
				{
					switch buffer[position] {
					case '/':
						if buffer[position] != rune('/') {
							goto l215
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
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
					case '.':
						if buffer[position] != rune('.') {
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

			l217:
				{
					position218, tokenIndex218 := position, tokenIndex
					{
						switch buffer[position] {
						case '/':
							if buffer[position] != rune('/') {
								goto l218
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l218
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l218
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l218
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l218
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l218
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l218
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l218
							}
							position++
							break
						}
					}

					goto l217
				l218:
					position, tokenIndex = position218, tokenIndex218
				}
				add(ruleStringValue, position216)
			}
			return true
		l215:
			position, tokenIndex = position215, tokenIndex215
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
		/* 20 Spacing <- <Space*> */
		func() bool {
			{
				position231 := position
			l232:
				{
					position233, tokenIndex233 := position, tokenIndex
					{
						position234 := position
						{
							position235, tokenIndex235 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l236
							}
							goto l235
						l236:
							position, tokenIndex = position235, tokenIndex235
							if !_rules[ruleEndOfLine]() {
								goto l233
							}
						}
					l235:
						add(ruleSpace, position234)
					}
					goto l232
				l233:
					position, tokenIndex = position233, tokenIndex233
				}
				add(ruleSpacing, position231)
			}
			return true
		},
		/* 21 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position238 := position
			l239:
				{
					position240, tokenIndex240 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l240
					}
					goto l239
				l240:
					position, tokenIndex = position240, tokenIndex240
				}
				add(ruleWhiteSpacing, position238)
			}
			return true
		},
		/* 22 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position241, tokenIndex241 := position, tokenIndex
			{
				position242 := position
				if !_rules[ruleWhitespace]() {
					goto l241
				}
			l243:
				{
					position244, tokenIndex244 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l244
					}
					goto l243
				l244:
					position, tokenIndex = position244, tokenIndex244
				}
				add(ruleMustWhiteSpacing, position242)
			}
			return true
		l241:
			position, tokenIndex = position241, tokenIndex241
			return false
		},
		/* 23 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position245, tokenIndex245 := position, tokenIndex
			{
				position246 := position
				if !_rules[ruleSpacing]() {
					goto l245
				}
				if buffer[position] != rune('=') {
					goto l245
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l245
				}
				add(ruleEqual, position246)
			}
			return true
		l245:
			position, tokenIndex = position245, tokenIndex245
			return false
		},
		/* 24 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 25 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position248, tokenIndex248 := position, tokenIndex
			{
				position249 := position
				{
					position250, tokenIndex250 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l251
					}
					position++
					goto l250
				l251:
					position, tokenIndex = position250, tokenIndex250
					if buffer[position] != rune('\t') {
						goto l248
					}
					position++
				}
			l250:
				add(ruleWhitespace, position249)
			}
			return true
		l248:
			position, tokenIndex = position248, tokenIndex248
			return false
		},
		/* 26 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position252, tokenIndex252 := position, tokenIndex
			{
				position253 := position
				{
					position254, tokenIndex254 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l255
					}
					position++
					if buffer[position] != rune('\n') {
						goto l255
					}
					position++
					goto l254
				l255:
					position, tokenIndex = position254, tokenIndex254
					if buffer[position] != rune('\n') {
						goto l256
					}
					position++
					goto l254
				l256:
					position, tokenIndex = position254, tokenIndex254
					if buffer[position] != rune('\r') {
						goto l252
					}
					position++
				}
			l254:
				add(ruleEndOfLine, position253)
			}
			return true
		l252:
			position, tokenIndex = position252, tokenIndex252
			return false
		},
		/* 27 EndOfFile <- <!.> */
		nil,
		nil,
		/* 30 Action0 <- <{ p.addDeclarationIdentifier(text) }> */
		nil,
		/* 31 Action1 <- <{ p.addAction(text) }> */
		nil,
		/* 32 Action2 <- <{ p.addEntity(text) }> */
		nil,
		/* 33 Action3 <- <{ p.LineDone() }> */
		nil,
		/* 34 Action4 <- <{ p.addParamKey(text) }> */
		nil,
		/* 35 Action5 <- <{  p.addParamHoleValue(text) }> */
		nil,
		/* 36 Action6 <- <{  p.addParamValue(text) }> */
		nil,
		/* 37 Action7 <- <{  p.addParamRefValue(text) }> */
		nil,
		/* 38 Action8 <- <{ p.addParamCidrValue(text) }> */
		nil,
		/* 39 Action9 <- <{ p.addParamIpValue(text) }> */
		nil,
		/* 40 Action10 <- <{p.addCsvValue(text)}> */
		nil,
		/* 41 Action11 <- <{ p.addParamValue(text) }> */
		nil,
		/* 42 Action12 <- <{ p.addParamIntValue(text) }> */
		nil,
		/* 43 Action13 <- <{ p.addParamValue(text) }> */
		nil,
		/* 44 Action14 <- <{ p.LineDone() }> */
		nil,
	}
	p.rules = _rules
}
