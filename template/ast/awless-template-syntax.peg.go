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
		/* 3 Entity <- <(('v' 'p' 'c') / ('s' 'u' 'b' 'n' 'e' 't') / ('i' 'n' 's' 't' 'a' 'n' 'c' 'e') / ('t' 'a' 'g') / ('s' 'e' 'c' 'u' 'r' 'i' 't' 'y' 'g' 'r' 'o' 'u' 'p') / ('r' 'o' 'u' 't' 'e' 't' 'a' 'b' 'l' 'e') / ('r' 'o' 'u' 't' 'e') / ('l' 'o' 'a' 'd' 'b' 'a' 'l' 'a' 'n' 'c' 'e' 'r') / ('t' 'a' 'r' 'g' 'e' 't' 'g' 'r' 'o' 'u' 'p') / ('d' 'a' 't' 'a' 'b' 'a' 's' 'e') / ('r' 'o' 'l' 'e') / ('s' 't' 'o' 'r' 'a' 'g' 'e' 'o' 'b' 'j' 'e' 'c' 't') / ((&('r') ('r' 'e' 'c' 'o' 'r' 'd')) | (&('z') ('z' 'o' 'n' 'e')) | (&('q') ('q' 'u' 'e' 'u' 'e')) | (&('t') ('t' 'o' 'p' 'i' 'c')) | (&('s') ('s' 'u' 'b' 's' 'c' 'r' 'i' 'p' 't' 'i' 'o' 'n')) | (&('b') ('b' 'u' 'c' 'k' 'e' 't')) | (&('a') ('a' 'c' 'c' 'e' 's' 's' 'k' 'e' 'y')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('d') ('d' 'b' 's' 'u' 'b' 'n' 'e' 't' 'g' 'r' 'o' 'u' 'p')) | (&('l') ('l' 'i' 's' 't' 'e' 'n' 'e' 'r')) | (&('i') ('i' 'n' 't' 'e' 'r' 'n' 'e' 't' 'g' 'a' 't' 'e' 'w' 'a' 'y')) | (&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e')) | (&('n') ('n' 'o' 'n' 'e'))))> */
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
							if buffer[position] != rune('s') {
								goto l80
							}
							position++
							if buffer[position] != rune('t') {
								goto l80
							}
							position++
							if buffer[position] != rune('o') {
								goto l80
							}
							position++
							if buffer[position] != rune('r') {
								goto l80
							}
							position++
							if buffer[position] != rune('a') {
								goto l80
							}
							position++
							if buffer[position] != rune('g') {
								goto l80
							}
							position++
							if buffer[position] != rune('e') {
								goto l80
							}
							position++
							if buffer[position] != rune('o') {
								goto l80
							}
							position++
							if buffer[position] != rune('b') {
								goto l80
							}
							position++
							if buffer[position] != rune('j') {
								goto l80
							}
							position++
							if buffer[position] != rune('e') {
								goto l80
							}
							position++
							if buffer[position] != rune('c') {
								goto l80
							}
							position++
							if buffer[position] != rune('t') {
								goto l80
							}
							position++
							goto l68
						l80:
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
					position83, tokenIndex83 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l83
					}
					{
						position85 := position
						{
							position88 := position
							{
								position89 := position
								if !_rules[ruleIdentifier]() {
									goto l83
								}
								add(rulePegText, position89)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l83
							}
							{
								position91 := position
								{
									position92, tokenIndex92 := position, tokenIndex
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
											if buffer[position] != rune('/') {
												goto l93
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l93
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
											add(ruleCidrValue, position95)
										}
										add(rulePegText, position94)
									}
									{
										add(ruleAction8, position)
									}
									goto l92
								l93:
									position, tokenIndex = position92, tokenIndex92
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
											if !matchDot() {
												goto l107
											}
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
											if !matchDot() {
												goto l107
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l107
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
											if !matchDot() {
												goto l107
											}
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l107
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
											add(ruleIpValue, position109)
										}
										add(rulePegText, position108)
									}
									{
										add(ruleAction9, position)
									}
									goto l92
								l107:
									position, tokenIndex = position92, tokenIndex92
									{
										position120 := position
										{
											position121 := position
											if !_rules[ruleStringValue]() {
												goto l119
											}
											if !_rules[ruleWhiteSpacing]() {
												goto l119
											}
											if buffer[position] != rune(',') {
												goto l119
											}
											position++
											if !_rules[ruleWhiteSpacing]() {
												goto l119
											}
										l122:
											{
												position123, tokenIndex123 := position, tokenIndex
												if !_rules[ruleStringValue]() {
													goto l123
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l123
												}
												if buffer[position] != rune(',') {
													goto l123
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l123
												}
												goto l122
											l123:
												position, tokenIndex = position123, tokenIndex123
											}
											if !_rules[ruleStringValue]() {
												goto l119
											}
											add(ruleCSVValue, position121)
										}
										add(rulePegText, position120)
									}
									{
										add(ruleAction10, position)
									}
									goto l92
								l119:
									position, tokenIndex = position92, tokenIndex92
									{
										position126 := position
										{
											position127 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l125
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
											if buffer[position] != rune('-') {
												goto l125
											}
											position++
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l125
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
											add(ruleIntRangeValue, position127)
										}
										add(rulePegText, position126)
									}
									{
										add(ruleAction11, position)
									}
									goto l92
								l125:
									position, tokenIndex = position92, tokenIndex92
									{
										position134 := position
										{
											position135 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l133
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
											add(ruleIntValue, position135)
										}
										add(rulePegText, position134)
									}
									{
										add(ruleAction12, position)
									}
									goto l92
								l133:
									position, tokenIndex = position92, tokenIndex92
									{
										switch buffer[position] {
										case '$':
											{
												position140 := position
												if buffer[position] != rune('$') {
													goto l83
												}
												position++
												{
													position141 := position
													if !_rules[ruleIdentifier]() {
														goto l83
													}
													add(rulePegText, position141)
												}
												add(ruleRefValue, position140)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position143 := position
												{
													position144 := position
													if buffer[position] != rune('@') {
														goto l83
													}
													position++
													if !_rules[ruleStringValue]() {
														goto l83
													}
													add(rulePegText, position144)
												}
												add(ruleAliasValue, position143)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position146 := position
												if buffer[position] != rune('{') {
													goto l83
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l83
												}
												{
													position147 := position
													if !_rules[ruleIdentifier]() {
														goto l83
													}
													add(rulePegText, position147)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l83
												}
												if buffer[position] != rune('}') {
													goto l83
												}
												position++
												add(ruleHoleValue, position146)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position149 := position
												if !_rules[ruleStringValue]() {
													goto l83
												}
												add(rulePegText, position149)
											}
											{
												add(ruleAction13, position)
											}
											break
										}
									}

								}
							l92:
								add(ruleValue, position91)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l83
							}
							add(ruleParam, position88)
						}
					l86:
						{
							position87, tokenIndex87 := position, tokenIndex
							{
								position151 := position
								{
									position152 := position
									if !_rules[ruleIdentifier]() {
										goto l87
									}
									add(rulePegText, position152)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l87
								}
								{
									position154 := position
									{
										position155, tokenIndex155 := position, tokenIndex
										{
											position157 := position
											{
												position158 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l156
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
													goto l156
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l156
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
													goto l156
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l156
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
												if !matchDot() {
													goto l156
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l156
												}
												position++
											l165:
												{
													position166, tokenIndex166 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l166
													}
													position++
													goto l165
												l166:
													position, tokenIndex = position166, tokenIndex166
												}
												if buffer[position] != rune('/') {
													goto l156
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l156
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
												add(ruleCidrValue, position158)
											}
											add(rulePegText, position157)
										}
										{
											add(ruleAction8, position)
										}
										goto l155
									l156:
										position, tokenIndex = position155, tokenIndex155
										{
											position171 := position
											{
												position172 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l170
												}
												position++
											l173:
												{
													position174, tokenIndex174 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l174
													}
													position++
													goto l173
												l174:
													position, tokenIndex = position174, tokenIndex174
												}
												if !matchDot() {
													goto l170
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l170
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
												if !matchDot() {
													goto l170
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l170
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
												if !matchDot() {
													goto l170
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l170
												}
												position++
											l179:
												{
													position180, tokenIndex180 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l180
													}
													position++
													goto l179
												l180:
													position, tokenIndex = position180, tokenIndex180
												}
												add(ruleIpValue, position172)
											}
											add(rulePegText, position171)
										}
										{
											add(ruleAction9, position)
										}
										goto l155
									l170:
										position, tokenIndex = position155, tokenIndex155
										{
											position183 := position
											{
												position184 := position
												if !_rules[ruleStringValue]() {
													goto l182
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l182
												}
												if buffer[position] != rune(',') {
													goto l182
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l182
												}
											l185:
												{
													position186, tokenIndex186 := position, tokenIndex
													if !_rules[ruleStringValue]() {
														goto l186
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l186
													}
													if buffer[position] != rune(',') {
														goto l186
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l186
													}
													goto l185
												l186:
													position, tokenIndex = position186, tokenIndex186
												}
												if !_rules[ruleStringValue]() {
													goto l182
												}
												add(ruleCSVValue, position184)
											}
											add(rulePegText, position183)
										}
										{
											add(ruleAction10, position)
										}
										goto l155
									l182:
										position, tokenIndex = position155, tokenIndex155
										{
											position189 := position
											{
												position190 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l188
												}
												position++
											l191:
												{
													position192, tokenIndex192 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l192
													}
													position++
													goto l191
												l192:
													position, tokenIndex = position192, tokenIndex192
												}
												if buffer[position] != rune('-') {
													goto l188
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l188
												}
												position++
											l193:
												{
													position194, tokenIndex194 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l194
													}
													position++
													goto l193
												l194:
													position, tokenIndex = position194, tokenIndex194
												}
												add(ruleIntRangeValue, position190)
											}
											add(rulePegText, position189)
										}
										{
											add(ruleAction11, position)
										}
										goto l155
									l188:
										position, tokenIndex = position155, tokenIndex155
										{
											position197 := position
											{
												position198 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l196
												}
												position++
											l199:
												{
													position200, tokenIndex200 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l200
													}
													position++
													goto l199
												l200:
													position, tokenIndex = position200, tokenIndex200
												}
												add(ruleIntValue, position198)
											}
											add(rulePegText, position197)
										}
										{
											add(ruleAction12, position)
										}
										goto l155
									l196:
										position, tokenIndex = position155, tokenIndex155
										{
											switch buffer[position] {
											case '$':
												{
													position203 := position
													if buffer[position] != rune('$') {
														goto l87
													}
													position++
													{
														position204 := position
														if !_rules[ruleIdentifier]() {
															goto l87
														}
														add(rulePegText, position204)
													}
													add(ruleRefValue, position203)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position206 := position
													{
														position207 := position
														if buffer[position] != rune('@') {
															goto l87
														}
														position++
														if !_rules[ruleStringValue]() {
															goto l87
														}
														add(rulePegText, position207)
													}
													add(ruleAliasValue, position206)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position209 := position
													if buffer[position] != rune('{') {
														goto l87
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l87
													}
													{
														position210 := position
														if !_rules[ruleIdentifier]() {
															goto l87
														}
														add(rulePegText, position210)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l87
													}
													if buffer[position] != rune('}') {
														goto l87
													}
													position++
													add(ruleHoleValue, position209)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position212 := position
													if !_rules[ruleStringValue]() {
														goto l87
													}
													add(rulePegText, position212)
												}
												{
													add(ruleAction13, position)
												}
												break
											}
										}

									}
								l155:
									add(ruleValue, position154)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l87
								}
								add(ruleParam, position151)
							}
							goto l86
						l87:
							position, tokenIndex = position87, tokenIndex87
						}
						add(ruleParams, position85)
					}
					goto l84
				l83:
					position, tokenIndex = position83, tokenIndex83
				}
			l84:
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
			position217, tokenIndex217 := position, tokenIndex
			{
				position218 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l217
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l217
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l217
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l217
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l217
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l217
						}
						position++
						break
					}
				}

			l219:
				{
					position220, tokenIndex220 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l220
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l220
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l220
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l220
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l220
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l220
							}
							position++
							break
						}
					}

					goto l219
				l220:
					position, tokenIndex = position220, tokenIndex220
				}
				add(ruleIdentifier, position218)
			}
			return true
		l217:
			position, tokenIndex = position217, tokenIndex217
			return false
		},
		/* 9 Value <- <((<CidrValue> Action8) / (<IpValue> Action9) / (<CSVValue> Action10) / (<IntRangeValue> Action11) / (<IntValue> Action12) / ((&('$') (RefValue Action7)) | (&('@') (AliasValue Action6)) | (&('{') (HoleValue Action5)) | (&('-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<StringValue> Action13))))> */
		nil,
		/* 10 StringValue <- <((&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position224, tokenIndex224 := position, tokenIndex
			{
				position225 := position
				{
					switch buffer[position] {
					case '/':
						if buffer[position] != rune('/') {
							goto l224
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l224
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l224
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l224
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l224
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l224
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l224
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l224
						}
						position++
						break
					}
				}

			l226:
				{
					position227, tokenIndex227 := position, tokenIndex
					{
						switch buffer[position] {
						case '/':
							if buffer[position] != rune('/') {
								goto l227
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l227
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l227
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l227
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l227
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l227
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l227
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l227
							}
							position++
							break
						}
					}

					goto l226
				l227:
					position, tokenIndex = position227, tokenIndex227
				}
				add(ruleStringValue, position225)
			}
			return true
		l224:
			position, tokenIndex = position224, tokenIndex224
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
				position240 := position
			l241:
				{
					position242, tokenIndex242 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l242
					}
					goto l241
				l242:
					position, tokenIndex = position242, tokenIndex242
				}
				add(ruleWhiteSpacing, position240)
			}
			return true
		},
		/* 21 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position243, tokenIndex243 := position, tokenIndex
			{
				position244 := position
				if !_rules[ruleWhitespace]() {
					goto l243
				}
			l245:
				{
					position246, tokenIndex246 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l246
					}
					goto l245
				l246:
					position, tokenIndex = position246, tokenIndex246
				}
				add(ruleMustWhiteSpacing, position244)
			}
			return true
		l243:
			position, tokenIndex = position243, tokenIndex243
			return false
		},
		/* 22 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position247, tokenIndex247 := position, tokenIndex
			{
				position248 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l247
				}
				if buffer[position] != rune('=') {
					goto l247
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l247
				}
				add(ruleEqual, position248)
			}
			return true
		l247:
			position, tokenIndex = position247, tokenIndex247
			return false
		},
		/* 23 BlankLine <- <(WhiteSpacing EndOfLine Action15)> */
		func() bool {
			position249, tokenIndex249 := position, tokenIndex
			{
				position250 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l249
				}
				if !_rules[ruleEndOfLine]() {
					goto l249
				}
				{
					add(ruleAction15, position)
				}
				add(ruleBlankLine, position250)
			}
			return true
		l249:
			position, tokenIndex = position249, tokenIndex249
			return false
		},
		/* 24 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position252, tokenIndex252 := position, tokenIndex
			{
				position253 := position
				{
					position254, tokenIndex254 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l255
					}
					position++
					goto l254
				l255:
					position, tokenIndex = position254, tokenIndex254
					if buffer[position] != rune('\t') {
						goto l252
					}
					position++
				}
			l254:
				add(ruleWhitespace, position253)
			}
			return true
		l252:
			position, tokenIndex = position252, tokenIndex252
			return false
		},
		/* 25 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position256, tokenIndex256 := position, tokenIndex
			{
				position257 := position
				{
					position258, tokenIndex258 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l259
					}
					position++
					if buffer[position] != rune('\n') {
						goto l259
					}
					position++
					goto l258
				l259:
					position, tokenIndex = position258, tokenIndex258
					if buffer[position] != rune('\n') {
						goto l260
					}
					position++
					goto l258
				l260:
					position, tokenIndex = position258, tokenIndex258
					if buffer[position] != rune('\r') {
						goto l256
					}
					position++
				}
			l258:
				add(ruleEndOfLine, position257)
			}
			return true
		l256:
			position, tokenIndex = position256, tokenIndex256
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
