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
			l25:
				{
					position26, tokenIndex26 := position, tokenIndex
					if !_rules[ruleBlankLine]() {
						goto l26
					}
					goto l25
				l26:
					position, tokenIndex = position26, tokenIndex26
				}
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
				l27:
					{
						position28, tokenIndex28 := position, tokenIndex
						if !_rules[ruleBlankLine]() {
							goto l28
						}
						goto l27
					l28:
						position, tokenIndex = position28, tokenIndex28
					}
					{
						position29 := position
						if !_rules[ruleWhiteSpacing]() {
							goto l3
						}
						{
							position30, tokenIndex30 := position, tokenIndex
							if !_rules[ruleExpr]() {
								goto l31
							}
							goto l30
						l31:
							position, tokenIndex = position30, tokenIndex30
							{
								position33 := position
								{
									position34 := position
									if !_rules[ruleIdentifier]() {
										goto l32
									}
									add(rulePegText, position34)
								}
								{
									add(ruleAction0, position)
								}
								if !_rules[ruleEqual]() {
									goto l32
								}
								if !_rules[ruleExpr]() {
									goto l32
								}
								add(ruleDeclaration, position33)
							}
							goto l30
						l32:
							position, tokenIndex = position30, tokenIndex30
							{
								position36 := position
								{
									position37, tokenIndex37 := position, tokenIndex
									if buffer[position] != rune('#') {
										goto l38
									}
									position++
								l39:
									{
										position40, tokenIndex40 := position, tokenIndex
										{
											position41, tokenIndex41 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l41
											}
											goto l40
										l41:
											position, tokenIndex = position41, tokenIndex41
										}
										if !matchDot() {
											goto l40
										}
										goto l39
									l40:
										position, tokenIndex = position40, tokenIndex40
									}
									goto l37
								l38:
									position, tokenIndex = position37, tokenIndex37
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
								l42:
									{
										position43, tokenIndex43 := position, tokenIndex
										{
											position44, tokenIndex44 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l44
											}
											goto l43
										l44:
											position, tokenIndex = position44, tokenIndex44
										}
										if !matchDot() {
											goto l43
										}
										goto l42
									l43:
										position, tokenIndex = position43, tokenIndex43
									}
									{
										add(ruleAction14, position)
									}
								}
							l37:
								add(ruleComment, position36)
							}
						}
					l30:
						if !_rules[ruleWhiteSpacing]() {
							goto l3
						}
					l46:
						{
							position47, tokenIndex47 := position, tokenIndex
							if !_rules[ruleEndOfLine]() {
								goto l47
							}
							goto l46
						l47:
							position, tokenIndex = position47, tokenIndex47
						}
						add(ruleStatement, position29)
					}
				l48:
					{
						position49, tokenIndex49 := position, tokenIndex
						if !_rules[ruleBlankLine]() {
							goto l49
						}
						goto l48
					l49:
						position, tokenIndex = position49, tokenIndex49
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				if !_rules[ruleWhiteSpacing]() {
					goto l0
				}
				{
					position50 := position
					{
						position51, tokenIndex51 := position, tokenIndex
						if !matchDot() {
							goto l51
						}
						goto l0
					l51:
						position, tokenIndex = position51, tokenIndex51
					}
					add(ruleEndOfFile, position50)
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
		/* 3 Entity <- <(('v' 'p' 'c') / ('s' 'u' 'b' 'n' 'e' 't') / ('i' 'n' 's' 't' 'a' 'n' 'c' 'e') / ('t' 'a' 'g') / ('s' 'e' 'c' 'u' 'r' 'i' 't' 'y' 'g' 'r' 'o' 'u' 'p') / ('r' 'o' 'u' 't' 'e' 't' 'a' 'b' 'l' 'e') / ('r' 'o' 'u' 't' 'e') / ('l' 'o' 'a' 'd' 'b' 'a' 'l' 'a' 'n' 'c' 'e' 'r') / ('t' 'a' 'r' 'g' 'e' 't' 'g' 'r' 'o' 'u' 'p') / ('d' 'a' 't' 'a' 'b' 'a' 's' 'e') / ('r' 'o' 'l' 'e') / ((&('r') ('r' 'e' 'c' 'o' 'r' 'd')) | (&('z') ('z' 'o' 'n' 'e')) | (&('q') ('q' 'u' 'e' 'u' 'e')) | (&('t') ('t' 'o' 'p' 'i' 'c')) | (&('s') ('s' 't' 'o' 'r' 'a' 'g' 'e' 'o' 'b' 'j' 'e' 'c' 't' ('s' 'u' 'b' 's' 'c' 'r' 'i' 'p' 't' 'i' 'o' 'n'))) | (&('b') ('b' 'u' 'c' 'k' 'e' 't')) | (&('a') ('a' 'c' 'c' 'e' 's' 's' 'k' 'e' 'y')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('d') ('d' 'b' 's' 'u' 'b' 'n' 'e' 't')) | (&('l') ('l' 'i' 's' 't' 'e' 'n' 'e' 'r')) | (&('i') ('i' 'n' 't' 'e' 'r' 'n' 'e' 't' 'g' 'a' 't' 'e' 'w' 'a' 'y')) | (&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e')) | (&('n') ('n' 'o' 'n' 'e'))))> */
		nil,
		/* 4 Declaration <- <(<Identifier> Action0 Equal Expr)> */
		nil,
		/* 5 Expr <- <(<Action> Action1 MustWhiteSpacing <Entity> Action2 (MustWhiteSpacing Params)? Action3)> */
		func() bool {
			position56, tokenIndex56 := position, tokenIndex
			{
				position57 := position
				{
					position58 := position
					{
						position59 := position
						{
							position60, tokenIndex60 := position, tokenIndex
							if buffer[position] != rune('c') {
								goto l61
							}
							position++
							if buffer[position] != rune('r') {
								goto l61
							}
							position++
							if buffer[position] != rune('e') {
								goto l61
							}
							position++
							if buffer[position] != rune('a') {
								goto l61
							}
							position++
							if buffer[position] != rune('t') {
								goto l61
							}
							position++
							if buffer[position] != rune('e') {
								goto l61
							}
							position++
							goto l60
						l61:
							position, tokenIndex = position60, tokenIndex60
							if buffer[position] != rune('d') {
								goto l62
							}
							position++
							if buffer[position] != rune('e') {
								goto l62
							}
							position++
							if buffer[position] != rune('l') {
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
							if buffer[position] != rune('e') {
								goto l62
							}
							position++
							goto l60
						l62:
							position, tokenIndex = position60, tokenIndex60
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
							if buffer[position] != rune('r') {
								goto l63
							}
							position++
							if buffer[position] != rune('t') {
								goto l63
							}
							position++
							goto l60
						l63:
							position, tokenIndex = position60, tokenIndex60
							{
								switch buffer[position] {
								case 'd':
									if buffer[position] != rune('d') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('a') {
										goto l56
									}
									position++
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									if buffer[position] != rune('h') {
										goto l56
									}
									position++
									break
								case 'c':
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									if buffer[position] != rune('h') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									if buffer[position] != rune('k') {
										goto l56
									}
									position++
									break
								case 'a':
									if buffer[position] != rune('a') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('a') {
										goto l56
									}
									position++
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									if buffer[position] != rune('h') {
										goto l56
									}
									position++
									break
								case 'u':
									if buffer[position] != rune('u') {
										goto l56
									}
									position++
									if buffer[position] != rune('p') {
										goto l56
									}
									position++
									if buffer[position] != rune('d') {
										goto l56
									}
									position++
									if buffer[position] != rune('a') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									break
								case 's':
									if buffer[position] != rune('s') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('p') {
										goto l56
									}
									position++
									break
								default:
									if buffer[position] != rune('n') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('n') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									break
								}
							}

						}
					l60:
						add(ruleAction, position59)
					}
					add(rulePegText, position58)
				}
				{
					add(ruleAction1, position)
				}
				if !_rules[ruleMustWhiteSpacing]() {
					goto l56
				}
				{
					position66 := position
					{
						position67 := position
						{
							position68, tokenIndex68 := position, tokenIndex
							if buffer[position] != rune('v') {
								goto l69
							}
							position++
							if buffer[position] != rune('p') {
								goto l69
							}
							position++
							if buffer[position] != rune('c') {
								goto l69
							}
							position++
							goto l68
						l69:
							position, tokenIndex = position68, tokenIndex68
							if buffer[position] != rune('s') {
								goto l70
							}
							position++
							if buffer[position] != rune('u') {
								goto l70
							}
							position++
							if buffer[position] != rune('b') {
								goto l70
							}
							position++
							if buffer[position] != rune('n') {
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
							goto l68
						l70:
							position, tokenIndex = position68, tokenIndex68
							if buffer[position] != rune('i') {
								goto l71
							}
							position++
							if buffer[position] != rune('n') {
								goto l71
							}
							position++
							if buffer[position] != rune('s') {
								goto l71
							}
							position++
							if buffer[position] != rune('t') {
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
							goto l68
						l71:
							position, tokenIndex = position68, tokenIndex68
							if buffer[position] != rune('t') {
								goto l72
							}
							position++
							if buffer[position] != rune('a') {
								goto l72
							}
							position++
							if buffer[position] != rune('g') {
								goto l72
							}
							position++
							goto l68
						l72:
							position, tokenIndex = position68, tokenIndex68
							if buffer[position] != rune('s') {
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
							if buffer[position] != rune('u') {
								goto l73
							}
							position++
							if buffer[position] != rune('r') {
								goto l73
							}
							position++
							if buffer[position] != rune('i') {
								goto l73
							}
							position++
							if buffer[position] != rune('t') {
								goto l73
							}
							position++
							if buffer[position] != rune('y') {
								goto l73
							}
							position++
							if buffer[position] != rune('g') {
								goto l73
							}
							position++
							if buffer[position] != rune('r') {
								goto l73
							}
							position++
							if buffer[position] != rune('o') {
								goto l73
							}
							position++
							if buffer[position] != rune('u') {
								goto l73
							}
							position++
							if buffer[position] != rune('p') {
								goto l73
							}
							position++
							goto l68
						l73:
							position, tokenIndex = position68, tokenIndex68
							if buffer[position] != rune('r') {
								goto l74
							}
							position++
							if buffer[position] != rune('o') {
								goto l74
							}
							position++
							if buffer[position] != rune('u') {
								goto l74
							}
							position++
							if buffer[position] != rune('t') {
								goto l74
							}
							position++
							if buffer[position] != rune('e') {
								goto l74
							}
							position++
							if buffer[position] != rune('t') {
								goto l74
							}
							position++
							if buffer[position] != rune('a') {
								goto l74
							}
							position++
							if buffer[position] != rune('b') {
								goto l74
							}
							position++
							if buffer[position] != rune('l') {
								goto l74
							}
							position++
							if buffer[position] != rune('e') {
								goto l74
							}
							position++
							goto l68
						l74:
							position, tokenIndex = position68, tokenIndex68
							if buffer[position] != rune('r') {
								goto l75
							}
							position++
							if buffer[position] != rune('o') {
								goto l75
							}
							position++
							if buffer[position] != rune('u') {
								goto l75
							}
							position++
							if buffer[position] != rune('t') {
								goto l75
							}
							position++
							if buffer[position] != rune('e') {
								goto l75
							}
							position++
							goto l68
						l75:
							position, tokenIndex = position68, tokenIndex68
							if buffer[position] != rune('l') {
								goto l76
							}
							position++
							if buffer[position] != rune('o') {
								goto l76
							}
							position++
							if buffer[position] != rune('a') {
								goto l76
							}
							position++
							if buffer[position] != rune('d') {
								goto l76
							}
							position++
							if buffer[position] != rune('b') {
								goto l76
							}
							position++
							if buffer[position] != rune('a') {
								goto l76
							}
							position++
							if buffer[position] != rune('l') {
								goto l76
							}
							position++
							if buffer[position] != rune('a') {
								goto l76
							}
							position++
							if buffer[position] != rune('n') {
								goto l76
							}
							position++
							if buffer[position] != rune('c') {
								goto l76
							}
							position++
							if buffer[position] != rune('e') {
								goto l76
							}
							position++
							if buffer[position] != rune('r') {
								goto l76
							}
							position++
							goto l68
						l76:
							position, tokenIndex = position68, tokenIndex68
							if buffer[position] != rune('t') {
								goto l77
							}
							position++
							if buffer[position] != rune('a') {
								goto l77
							}
							position++
							if buffer[position] != rune('r') {
								goto l77
							}
							position++
							if buffer[position] != rune('g') {
								goto l77
							}
							position++
							if buffer[position] != rune('e') {
								goto l77
							}
							position++
							if buffer[position] != rune('t') {
								goto l77
							}
							position++
							if buffer[position] != rune('g') {
								goto l77
							}
							position++
							if buffer[position] != rune('r') {
								goto l77
							}
							position++
							if buffer[position] != rune('o') {
								goto l77
							}
							position++
							if buffer[position] != rune('u') {
								goto l77
							}
							position++
							if buffer[position] != rune('p') {
								goto l77
							}
							position++
							goto l68
						l77:
							position, tokenIndex = position68, tokenIndex68
							if buffer[position] != rune('d') {
								goto l78
							}
							position++
							if buffer[position] != rune('a') {
								goto l78
							}
							position++
							if buffer[position] != rune('t') {
								goto l78
							}
							position++
							if buffer[position] != rune('a') {
								goto l78
							}
							position++
							if buffer[position] != rune('b') {
								goto l78
							}
							position++
							if buffer[position] != rune('a') {
								goto l78
							}
							position++
							if buffer[position] != rune('s') {
								goto l78
							}
							position++
							if buffer[position] != rune('e') {
								goto l78
							}
							position++
							goto l68
						l78:
							position, tokenIndex = position68, tokenIndex68
							if buffer[position] != rune('r') {
								goto l79
							}
							position++
							if buffer[position] != rune('o') {
								goto l79
							}
							position++
							if buffer[position] != rune('l') {
								goto l79
							}
							position++
							if buffer[position] != rune('e') {
								goto l79
							}
							position++
							goto l68
						l79:
							position, tokenIndex = position68, tokenIndex68
							{
								switch buffer[position] {
								case 'r':
									if buffer[position] != rune('r') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('r') {
										goto l56
									}
									position++
									if buffer[position] != rune('d') {
										goto l56
									}
									position++
									break
								case 'z':
									if buffer[position] != rune('z') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('n') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									break
								case 'q':
									if buffer[position] != rune('q') {
										goto l56
									}
									position++
									if buffer[position] != rune('u') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('u') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									break
								case 't':
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('p') {
										goto l56
									}
									position++
									if buffer[position] != rune('i') {
										goto l56
									}
									position++
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									break
								case 's':
									if buffer[position] != rune('s') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('r') {
										goto l56
									}
									position++
									if buffer[position] != rune('a') {
										goto l56
									}
									position++
									if buffer[position] != rune('g') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('b') {
										goto l56
									}
									position++
									if buffer[position] != rune('j') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('s') {
										goto l56
									}
									position++
									if buffer[position] != rune('u') {
										goto l56
									}
									position++
									if buffer[position] != rune('b') {
										goto l56
									}
									position++
									if buffer[position] != rune('s') {
										goto l56
									}
									position++
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									if buffer[position] != rune('r') {
										goto l56
									}
									position++
									if buffer[position] != rune('i') {
										goto l56
									}
									position++
									if buffer[position] != rune('p') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('i') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('n') {
										goto l56
									}
									position++
									break
								case 'b':
									if buffer[position] != rune('b') {
										goto l56
									}
									position++
									if buffer[position] != rune('u') {
										goto l56
									}
									position++
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									if buffer[position] != rune('k') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									break
								case 'a':
									if buffer[position] != rune('a') {
										goto l56
									}
									position++
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('s') {
										goto l56
									}
									position++
									if buffer[position] != rune('s') {
										goto l56
									}
									position++
									if buffer[position] != rune('k') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('y') {
										goto l56
									}
									position++
									break
								case 'p':
									if buffer[position] != rune('p') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('l') {
										goto l56
									}
									position++
									if buffer[position] != rune('i') {
										goto l56
									}
									position++
									if buffer[position] != rune('c') {
										goto l56
									}
									position++
									if buffer[position] != rune('y') {
										goto l56
									}
									position++
									break
								case 'g':
									if buffer[position] != rune('g') {
										goto l56
									}
									position++
									if buffer[position] != rune('r') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('u') {
										goto l56
									}
									position++
									if buffer[position] != rune('p') {
										goto l56
									}
									position++
									break
								case 'u':
									if buffer[position] != rune('u') {
										goto l56
									}
									position++
									if buffer[position] != rune('s') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('r') {
										goto l56
									}
									position++
									break
								case 'd':
									if buffer[position] != rune('d') {
										goto l56
									}
									position++
									if buffer[position] != rune('b') {
										goto l56
									}
									position++
									if buffer[position] != rune('s') {
										goto l56
									}
									position++
									if buffer[position] != rune('u') {
										goto l56
									}
									position++
									if buffer[position] != rune('b') {
										goto l56
									}
									position++
									if buffer[position] != rune('n') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									break
								case 'l':
									if buffer[position] != rune('l') {
										goto l56
									}
									position++
									if buffer[position] != rune('i') {
										goto l56
									}
									position++
									if buffer[position] != rune('s') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('n') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('r') {
										goto l56
									}
									position++
									break
								case 'i':
									if buffer[position] != rune('i') {
										goto l56
									}
									position++
									if buffer[position] != rune('n') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('r') {
										goto l56
									}
									position++
									if buffer[position] != rune('n') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('g') {
										goto l56
									}
									position++
									if buffer[position] != rune('a') {
										goto l56
									}
									position++
									if buffer[position] != rune('t') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('w') {
										goto l56
									}
									position++
									if buffer[position] != rune('a') {
										goto l56
									}
									position++
									if buffer[position] != rune('y') {
										goto l56
									}
									position++
									break
								case 'k':
									if buffer[position] != rune('k') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									if buffer[position] != rune('y') {
										goto l56
									}
									position++
									if buffer[position] != rune('p') {
										goto l56
									}
									position++
									if buffer[position] != rune('a') {
										goto l56
									}
									position++
									if buffer[position] != rune('i') {
										goto l56
									}
									position++
									if buffer[position] != rune('r') {
										goto l56
									}
									position++
									break
								case 'v':
									if buffer[position] != rune('v') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('l') {
										goto l56
									}
									position++
									if buffer[position] != rune('u') {
										goto l56
									}
									position++
									if buffer[position] != rune('m') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									break
								default:
									if buffer[position] != rune('n') {
										goto l56
									}
									position++
									if buffer[position] != rune('o') {
										goto l56
									}
									position++
									if buffer[position] != rune('n') {
										goto l56
									}
									position++
									if buffer[position] != rune('e') {
										goto l56
									}
									position++
									break
								}
							}

						}
					l68:
						add(ruleEntity, position67)
					}
					add(rulePegText, position66)
				}
				{
					add(ruleAction2, position)
				}
				{
					position82, tokenIndex82 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l82
					}
					{
						position84 := position
						{
							position87 := position
							{
								position88 := position
								if !_rules[ruleIdentifier]() {
									goto l82
								}
								add(rulePegText, position88)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l82
							}
							{
								position90 := position
								{
									position91, tokenIndex91 := position, tokenIndex
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
											if buffer[position] != rune('/') {
												goto l92
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l92
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
											add(ruleCidrValue, position94)
										}
										add(rulePegText, position93)
									}
									{
										add(ruleAction8, position)
									}
									goto l91
								l92:
									position, tokenIndex = position91, tokenIndex91
									{
										position107 := position
										{
											position108 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l106
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
												goto l106
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l106
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
												goto l106
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l106
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
											if !matchDot() {
												goto l106
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l106
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
											add(ruleIpValue, position108)
										}
										add(rulePegText, position107)
									}
									{
										add(ruleAction9, position)
									}
									goto l91
								l106:
									position, tokenIndex = position91, tokenIndex91
									{
										position119 := position
										{
											position120 := position
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
										l121:
											{
												position122, tokenIndex122 := position, tokenIndex
												if !_rules[ruleStringValue]() {
													goto l122
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l122
												}
												if buffer[position] != rune(',') {
													goto l122
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l122
												}
												goto l121
											l122:
												position, tokenIndex = position122, tokenIndex122
											}
											if !_rules[ruleStringValue]() {
												goto l118
											}
											add(ruleCSVValue, position120)
										}
										add(rulePegText, position119)
									}
									{
										add(ruleAction10, position)
									}
									goto l91
								l118:
									position, tokenIndex = position91, tokenIndex91
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
											if buffer[position] != rune('-') {
												goto l124
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l124
											}
											position++
										l129:
											{
												position130, tokenIndex130 := position, tokenIndex
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l130
												}
												position++
												goto l129
											l130:
												position, tokenIndex = position130, tokenIndex130
											}
											add(ruleIntRangeValue, position126)
										}
										add(rulePegText, position125)
									}
									{
										add(ruleAction11, position)
									}
									goto l91
								l124:
									position, tokenIndex = position91, tokenIndex91
									{
										position133 := position
										{
											position134 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l132
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
											add(ruleIntValue, position134)
										}
										add(rulePegText, position133)
									}
									{
										add(ruleAction12, position)
									}
									goto l91
								l132:
									position, tokenIndex = position91, tokenIndex91
									{
										switch buffer[position] {
										case '$':
											{
												position139 := position
												if buffer[position] != rune('$') {
													goto l82
												}
												position++
												{
													position140 := position
													if !_rules[ruleIdentifier]() {
														goto l82
													}
													add(rulePegText, position140)
												}
												add(ruleRefValue, position139)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position142 := position
												{
													position143 := position
													if buffer[position] != rune('@') {
														goto l82
													}
													position++
													if !_rules[ruleStringValue]() {
														goto l82
													}
													add(rulePegText, position143)
												}
												add(ruleAliasValue, position142)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position145 := position
												if buffer[position] != rune('{') {
													goto l82
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l82
												}
												{
													position146 := position
													if !_rules[ruleIdentifier]() {
														goto l82
													}
													add(rulePegText, position146)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l82
												}
												if buffer[position] != rune('}') {
													goto l82
												}
												position++
												add(ruleHoleValue, position145)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position148 := position
												if !_rules[ruleStringValue]() {
													goto l82
												}
												add(rulePegText, position148)
											}
											{
												add(ruleAction13, position)
											}
											break
										}
									}

								}
							l91:
								add(ruleValue, position90)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l82
							}
							add(ruleParam, position87)
						}
					l85:
						{
							position86, tokenIndex86 := position, tokenIndex
							{
								position150 := position
								{
									position151 := position
									if !_rules[ruleIdentifier]() {
										goto l86
									}
									add(rulePegText, position151)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l86
								}
								{
									position153 := position
									{
										position154, tokenIndex154 := position, tokenIndex
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
												if buffer[position] != rune('/') {
													goto l155
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l155
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
												add(ruleCidrValue, position157)
											}
											add(rulePegText, position156)
										}
										{
											add(ruleAction8, position)
										}
										goto l154
									l155:
										position, tokenIndex = position154, tokenIndex154
										{
											position170 := position
											{
												position171 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l169
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
													goto l169
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l169
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
												if !matchDot() {
													goto l169
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l169
												}
												position++
											l176:
												{
													position177, tokenIndex177 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l177
													}
													position++
													goto l176
												l177:
													position, tokenIndex = position177, tokenIndex177
												}
												if !matchDot() {
													goto l169
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l169
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
												add(ruleIpValue, position171)
											}
											add(rulePegText, position170)
										}
										{
											add(ruleAction9, position)
										}
										goto l154
									l169:
										position, tokenIndex = position154, tokenIndex154
										{
											position182 := position
											{
												position183 := position
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
											l184:
												{
													position185, tokenIndex185 := position, tokenIndex
													if !_rules[ruleStringValue]() {
														goto l185
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l185
													}
													if buffer[position] != rune(',') {
														goto l185
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l185
													}
													goto l184
												l185:
													position, tokenIndex = position185, tokenIndex185
												}
												if !_rules[ruleStringValue]() {
													goto l181
												}
												add(ruleCSVValue, position183)
											}
											add(rulePegText, position182)
										}
										{
											add(ruleAction10, position)
										}
										goto l154
									l181:
										position, tokenIndex = position154, tokenIndex154
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
												if buffer[position] != rune('-') {
													goto l187
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l187
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
												add(ruleIntRangeValue, position189)
											}
											add(rulePegText, position188)
										}
										{
											add(ruleAction11, position)
										}
										goto l154
									l187:
										position, tokenIndex = position154, tokenIndex154
										{
											position196 := position
											{
												position197 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l195
												}
												position++
											l198:
												{
													position199, tokenIndex199 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l199
													}
													position++
													goto l198
												l199:
													position, tokenIndex = position199, tokenIndex199
												}
												add(ruleIntValue, position197)
											}
											add(rulePegText, position196)
										}
										{
											add(ruleAction12, position)
										}
										goto l154
									l195:
										position, tokenIndex = position154, tokenIndex154
										{
											switch buffer[position] {
											case '$':
												{
													position202 := position
													if buffer[position] != rune('$') {
														goto l86
													}
													position++
													{
														position203 := position
														if !_rules[ruleIdentifier]() {
															goto l86
														}
														add(rulePegText, position203)
													}
													add(ruleRefValue, position202)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position205 := position
													{
														position206 := position
														if buffer[position] != rune('@') {
															goto l86
														}
														position++
														if !_rules[ruleStringValue]() {
															goto l86
														}
														add(rulePegText, position206)
													}
													add(ruleAliasValue, position205)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position208 := position
													if buffer[position] != rune('{') {
														goto l86
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l86
													}
													{
														position209 := position
														if !_rules[ruleIdentifier]() {
															goto l86
														}
														add(rulePegText, position209)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l86
													}
													if buffer[position] != rune('}') {
														goto l86
													}
													position++
													add(ruleHoleValue, position208)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position211 := position
													if !_rules[ruleStringValue]() {
														goto l86
													}
													add(rulePegText, position211)
												}
												{
													add(ruleAction13, position)
												}
												break
											}
										}

									}
								l154:
									add(ruleValue, position153)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l86
								}
								add(ruleParam, position150)
							}
							goto l85
						l86:
							position, tokenIndex = position86, tokenIndex86
						}
						add(ruleParams, position84)
					}
					goto l83
				l82:
					position, tokenIndex = position82, tokenIndex82
				}
			l83:
				{
					add(ruleAction3, position)
				}
				add(ruleExpr, position57)
			}
			return true
		l56:
			position, tokenIndex = position56, tokenIndex56
			return false
		},
		/* 6 Params <- <Param+> */
		nil,
		/* 7 Param <- <(<Identifier> Action4 Equal Value WhiteSpacing)> */
		nil,
		/* 8 Identifier <- <((&('.') '.') | (&('_') '_') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position216, tokenIndex216 := position, tokenIndex
			{
				position217 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l216
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l216
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l216
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l216
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l216
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l216
						}
						position++
						break
					}
				}

			l218:
				{
					position219, tokenIndex219 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
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

					goto l218
				l219:
					position, tokenIndex = position219, tokenIndex219
				}
				add(ruleIdentifier, position217)
			}
			return true
		l216:
			position, tokenIndex = position216, tokenIndex216
			return false
		},
		/* 9 Value <- <((<CidrValue> Action8) / (<IpValue> Action9) / (<CSVValue> Action10) / (<IntRangeValue> Action11) / (<IntValue> Action12) / ((&('$') (RefValue Action7)) | (&('@') (AliasValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action13))))> */
		nil,
		/* 10 StringValue <- <((&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position223, tokenIndex223 := position, tokenIndex
			{
				position224 := position
				{
					switch buffer[position] {
					case '/':
						if buffer[position] != rune('/') {
							goto l223
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l223
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l223
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l223
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l223
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l223
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l223
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l223
						}
						position++
						break
					}
				}

			l225:
				{
					position226, tokenIndex226 := position, tokenIndex
					{
						switch buffer[position] {
						case '/':
							if buffer[position] != rune('/') {
								goto l226
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l226
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l226
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l226
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l226
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l226
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l226
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l226
							}
							position++
							break
						}
					}

					goto l225
				l226:
					position, tokenIndex = position226, tokenIndex226
				}
				add(ruleStringValue, position224)
			}
			return true
		l223:
			position, tokenIndex = position223, tokenIndex223
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
				position239 := position
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
				add(ruleWhiteSpacing, position239)
			}
			return true
		},
		/* 21 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position242, tokenIndex242 := position, tokenIndex
			{
				position243 := position
				if !_rules[ruleWhitespace]() {
					goto l242
				}
			l244:
				{
					position245, tokenIndex245 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l245
					}
					goto l244
				l245:
					position, tokenIndex = position245, tokenIndex245
				}
				add(ruleMustWhiteSpacing, position243)
			}
			return true
		l242:
			position, tokenIndex = position242, tokenIndex242
			return false
		},
		/* 22 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position246, tokenIndex246 := position, tokenIndex
			{
				position247 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l246
				}
				if buffer[position] != rune('=') {
					goto l246
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l246
				}
				add(ruleEqual, position247)
			}
			return true
		l246:
			position, tokenIndex = position246, tokenIndex246
			return false
		},
		/* 23 BlankLine <- <(WhiteSpacing EndOfLine Action15)> */
		func() bool {
			position248, tokenIndex248 := position, tokenIndex
			{
				position249 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l248
				}
				if !_rules[ruleEndOfLine]() {
					goto l248
				}
				{
					add(ruleAction15, position)
				}
				add(ruleBlankLine, position249)
			}
			return true
		l248:
			position, tokenIndex = position248, tokenIndex248
			return false
		},
		/* 24 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position251, tokenIndex251 := position, tokenIndex
			{
				position252 := position
				{
					position253, tokenIndex253 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l254
					}
					position++
					goto l253
				l254:
					position, tokenIndex = position253, tokenIndex253
					if buffer[position] != rune('\t') {
						goto l251
					}
					position++
				}
			l253:
				add(ruleWhitespace, position252)
			}
			return true
		l251:
			position, tokenIndex = position251, tokenIndex251
			return false
		},
		/* 25 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position255, tokenIndex255 := position, tokenIndex
			{
				position256 := position
				{
					position257, tokenIndex257 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l258
					}
					position++
					if buffer[position] != rune('\n') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex = position257, tokenIndex257
					if buffer[position] != rune('\n') {
						goto l259
					}
					position++
					goto l257
				l259:
					position, tokenIndex = position257, tokenIndex257
					if buffer[position] != rune('\r') {
						goto l255
					}
					position++
				}
			l257:
				add(ruleEndOfLine, position256)
			}
			return true
		l255:
			position, tokenIndex = position255, tokenIndex255
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
