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
		/* 2 Action <- <(('c' 'r' 'e' 'a' 't' 'e') / ('s' 't' 'a' 'r' 't') / ((&('c') ('c' 'h' 'e' 'c' 'k')) | (&('a') ('a' 't' 't' 'a' 'c' 'h')) | (&('u') ('u' 'p' 'd' 'a' 't' 'e')) | (&('s') ('s' 't' 'o' 'p')) | (&('d') ('d' 'e' 'l' 'e' 't' 'e'))))> */
		nil,
		/* 3 Entity <- <(('v' 'p' 'c') / ('s' 'u' 'b' 'n' 'e' 't') / ('i' 'n' 's' 't' 'a' 'n' 'c' 'e') / ((&('i') ('i' 'n' 't' 'e' 'r' 'n' 'e' 't' 'g' 'a' 't' 'e' 'w' 'a' 'y')) | (&('s') ('s' 'e' 'c' 'u' 'r' 'i' 't' 'y' 'g' 'r' 'o' 'u' 'p')) | (&('k') ('k' 'e' 'y' 'p' 'a' 'i' 'r')) | (&('p') ('p' 'o' 'l' 'i' 'c' 'y')) | (&('r') ('r' 'o' 'l' 'e')) | (&('g') ('g' 'r' 'o' 'u' 'p')) | (&('u') ('u' 's' 'e' 'r')) | (&('t') ('t' 'a' 'g' 's')) | (&('v') ('v' 'o' 'l' 'u' 'm' 'e'))))> */
		nil,
		/* 4 Declaration <- <(<Identifier> Action0 Equal Expr)> */
		nil,
		/* 5 Expr <- <(<Action> Action1 MustWhiteSpacing <Entity> Action2 (MustWhiteSpacing Params)? Action3)> */
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
							if buffer[position] != rune('s') {
								goto l32
							}
							position++
							if buffer[position] != rune('t') {
								goto l32
							}
							position++
							if buffer[position] != rune('a') {
								goto l32
							}
							position++
							if buffer[position] != rune('r') {
								goto l32
							}
							position++
							if buffer[position] != rune('t') {
								goto l32
							}
							position++
							goto l30
						l32:
							position, tokenIndex = position30, tokenIndex30
							{
								switch buffer[position] {
								case 'c':
									if buffer[position] != rune('c') {
										goto l26
									}
									position++
									if buffer[position] != rune('h') {
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
									if buffer[position] != rune('k') {
										goto l26
									}
									position++
									break
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
								default:
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
					position35 := position
					{
						position36 := position
						{
							position37, tokenIndex37 := position, tokenIndex
							if buffer[position] != rune('v') {
								goto l38
							}
							position++
							if buffer[position] != rune('p') {
								goto l38
							}
							position++
							if buffer[position] != rune('c') {
								goto l38
							}
							position++
							goto l37
						l38:
							position, tokenIndex = position37, tokenIndex37
							if buffer[position] != rune('s') {
								goto l39
							}
							position++
							if buffer[position] != rune('u') {
								goto l39
							}
							position++
							if buffer[position] != rune('b') {
								goto l39
							}
							position++
							if buffer[position] != rune('n') {
								goto l39
							}
							position++
							if buffer[position] != rune('e') {
								goto l39
							}
							position++
							if buffer[position] != rune('t') {
								goto l39
							}
							position++
							goto l37
						l39:
							position, tokenIndex = position37, tokenIndex37
							if buffer[position] != rune('i') {
								goto l40
							}
							position++
							if buffer[position] != rune('n') {
								goto l40
							}
							position++
							if buffer[position] != rune('s') {
								goto l40
							}
							position++
							if buffer[position] != rune('t') {
								goto l40
							}
							position++
							if buffer[position] != rune('a') {
								goto l40
							}
							position++
							if buffer[position] != rune('n') {
								goto l40
							}
							position++
							if buffer[position] != rune('c') {
								goto l40
							}
							position++
							if buffer[position] != rune('e') {
								goto l40
							}
							position++
							goto l37
						l40:
							position, tokenIndex = position37, tokenIndex37
							{
								switch buffer[position] {
								case 'i':
									if buffer[position] != rune('i') {
										goto l26
									}
									position++
									if buffer[position] != rune('n') {
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
									if buffer[position] != rune('r') {
										goto l26
									}
									position++
									if buffer[position] != rune('n') {
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
									if buffer[position] != rune('g') {
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
									if buffer[position] != rune('w') {
										goto l26
									}
									position++
									if buffer[position] != rune('a') {
										goto l26
									}
									position++
									if buffer[position] != rune('y') {
										goto l26
									}
									position++
									break
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
								default:
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
								}
							}

						}
					l37:
						add(ruleEntity, position36)
					}
					add(rulePegText, position35)
				}
				{
					add(ruleAction2, position)
				}
				{
					position43, tokenIndex43 := position, tokenIndex
					if !_rules[ruleMustWhiteSpacing]() {
						goto l43
					}
					{
						position45 := position
						{
							position48 := position
							{
								position49 := position
								if !_rules[ruleIdentifier]() {
									goto l43
								}
								add(rulePegText, position49)
							}
							{
								add(ruleAction4, position)
							}
							if !_rules[ruleEqual]() {
								goto l43
							}
							{
								position51 := position
								{
									position52, tokenIndex52 := position, tokenIndex
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
										add(ruleAction8, position)
									}
									goto l52
								l53:
									position, tokenIndex = position52, tokenIndex52
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
										add(ruleAction9, position)
									}
									goto l52
								l67:
									position, tokenIndex = position52, tokenIndex52
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
											if buffer[position] != rune('-') {
												goto l79
											}
											position++
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
											add(ruleIntRangeValue, position81)
										}
										add(rulePegText, position80)
									}
									{
										add(ruleAction10, position)
									}
									goto l52
								l79:
									position, tokenIndex = position52, tokenIndex52
									{
										position88 := position
										{
											position89 := position
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l87
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
											add(ruleIntValue, position89)
										}
										add(rulePegText, position88)
									}
									{
										add(ruleAction11, position)
									}
									goto l52
								l87:
									position, tokenIndex = position52, tokenIndex52
									{
										switch buffer[position] {
										case '$':
											{
												position94 := position
												if buffer[position] != rune('$') {
													goto l43
												}
												position++
												{
													position95 := position
													if !_rules[ruleIdentifier]() {
														goto l43
													}
													add(rulePegText, position95)
												}
												add(ruleRefValue, position94)
											}
											{
												add(ruleAction7, position)
											}
											break
										case '@':
											{
												position97 := position
												if buffer[position] != rune('@') {
													goto l43
												}
												position++
												{
													position98 := position
													if !_rules[ruleIdentifier]() {
														goto l43
													}
													add(rulePegText, position98)
												}
												add(ruleAliasValue, position97)
											}
											{
												add(ruleAction6, position)
											}
											break
										case '{':
											{
												position100 := position
												if buffer[position] != rune('{') {
													goto l43
												}
												position++
												if !_rules[ruleWhiteSpacing]() {
													goto l43
												}
												{
													position101 := position
													if !_rules[ruleIdentifier]() {
														goto l43
													}
													add(rulePegText, position101)
												}
												if !_rules[ruleWhiteSpacing]() {
													goto l43
												}
												if buffer[position] != rune('}') {
													goto l43
												}
												position++
												add(ruleHoleValue, position100)
											}
											{
												add(ruleAction5, position)
											}
											break
										default:
											{
												position103 := position
												{
													position104 := position
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

												l105:
													{
														position106, tokenIndex106 := position, tokenIndex
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l106
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l106
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l106
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l106
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l106
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l106
																}
																position++
																break
															}
														}

														goto l105
													l106:
														position, tokenIndex = position106, tokenIndex106
													}
													add(ruleStringValue, position104)
												}
												add(rulePegText, position103)
											}
											{
												add(ruleAction12, position)
											}
											break
										}
									}

								}
							l52:
								add(ruleValue, position51)
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l43
							}
							add(ruleParam, position48)
						}
					l46:
						{
							position47, tokenIndex47 := position, tokenIndex
							{
								position110 := position
								{
									position111 := position
									if !_rules[ruleIdentifier]() {
										goto l47
									}
									add(rulePegText, position111)
								}
								{
									add(ruleAction4, position)
								}
								if !_rules[ruleEqual]() {
									goto l47
								}
								{
									position113 := position
									{
										position114, tokenIndex114 := position, tokenIndex
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
												if !matchDot() {
													goto l115
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l115
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
													goto l115
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l115
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
													goto l115
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l115
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
												if buffer[position] != rune('/') {
													goto l115
												}
												position++
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l115
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
												add(ruleCidrValue, position117)
											}
											add(rulePegText, position116)
										}
										{
											add(ruleAction8, position)
										}
										goto l114
									l115:
										position, tokenIndex = position114, tokenIndex114
										{
											position130 := position
											{
												position131 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l129
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
													goto l129
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l129
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
													goto l129
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l129
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
												if !matchDot() {
													goto l129
												}
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l129
												}
												position++
											l138:
												{
													position139, tokenIndex139 := position, tokenIndex
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l139
													}
													position++
													goto l138
												l139:
													position, tokenIndex = position139, tokenIndex139
												}
												add(ruleIpValue, position131)
											}
											add(rulePegText, position130)
										}
										{
											add(ruleAction9, position)
										}
										goto l114
									l129:
										position, tokenIndex = position114, tokenIndex114
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
												if buffer[position] != rune('-') {
													goto l141
												}
												position++
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
												add(ruleIntRangeValue, position143)
											}
											add(rulePegText, position142)
										}
										{
											add(ruleAction10, position)
										}
										goto l114
									l141:
										position, tokenIndex = position114, tokenIndex114
										{
											position150 := position
											{
												position151 := position
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l149
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
												add(ruleIntValue, position151)
											}
											add(rulePegText, position150)
										}
										{
											add(ruleAction11, position)
										}
										goto l114
									l149:
										position, tokenIndex = position114, tokenIndex114
										{
											switch buffer[position] {
											case '$':
												{
													position156 := position
													if buffer[position] != rune('$') {
														goto l47
													}
													position++
													{
														position157 := position
														if !_rules[ruleIdentifier]() {
															goto l47
														}
														add(rulePegText, position157)
													}
													add(ruleRefValue, position156)
												}
												{
													add(ruleAction7, position)
												}
												break
											case '@':
												{
													position159 := position
													if buffer[position] != rune('@') {
														goto l47
													}
													position++
													{
														position160 := position
														if !_rules[ruleIdentifier]() {
															goto l47
														}
														add(rulePegText, position160)
													}
													add(ruleAliasValue, position159)
												}
												{
													add(ruleAction6, position)
												}
												break
											case '{':
												{
													position162 := position
													if buffer[position] != rune('{') {
														goto l47
													}
													position++
													if !_rules[ruleWhiteSpacing]() {
														goto l47
													}
													{
														position163 := position
														if !_rules[ruleIdentifier]() {
															goto l47
														}
														add(rulePegText, position163)
													}
													if !_rules[ruleWhiteSpacing]() {
														goto l47
													}
													if buffer[position] != rune('}') {
														goto l47
													}
													position++
													add(ruleHoleValue, position162)
												}
												{
													add(ruleAction5, position)
												}
												break
											default:
												{
													position165 := position
													{
														position166 := position
														{
															switch buffer[position] {
															case '_':
																if buffer[position] != rune('_') {
																	goto l47
																}
																position++
																break
															case '.':
																if buffer[position] != rune('.') {
																	goto l47
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l47
																}
																position++
																break
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l47
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l47
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l47
																}
																position++
																break
															}
														}

													l167:
														{
															position168, tokenIndex168 := position, tokenIndex
															{
																switch buffer[position] {
																case '_':
																	if buffer[position] != rune('_') {
																		goto l168
																	}
																	position++
																	break
																case '.':
																	if buffer[position] != rune('.') {
																		goto l168
																	}
																	position++
																	break
																case '-':
																	if buffer[position] != rune('-') {
																		goto l168
																	}
																	position++
																	break
																case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l168
																	}
																	position++
																	break
																case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																	if c := buffer[position]; c < rune('A') || c > rune('Z') {
																		goto l168
																	}
																	position++
																	break
																default:
																	if c := buffer[position]; c < rune('a') || c > rune('z') {
																		goto l168
																	}
																	position++
																	break
																}
															}

															goto l167
														l168:
															position, tokenIndex = position168, tokenIndex168
														}
														add(ruleStringValue, position166)
													}
													add(rulePegText, position165)
												}
												{
													add(ruleAction12, position)
												}
												break
											}
										}

									}
								l114:
									add(ruleValue, position113)
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l47
								}
								add(ruleParam, position110)
							}
							goto l46
						l47:
							position, tokenIndex = position47, tokenIndex47
						}
						add(ruleParams, position45)
					}
					goto l44
				l43:
					position, tokenIndex = position43, tokenIndex43
				}
			l44:
				{
					add(ruleAction3, position)
				}
				add(ruleExpr, position27)
			}
			return true
		l26:
			position, tokenIndex = position26, tokenIndex26
			return false
		},
		/* 6 Params <- <Param+> */
		nil,
		/* 7 Param <- <(<Identifier> Action4 Equal Value WhiteSpacing)> */
		nil,
		/* 8 Identifier <- <((&('.') '.') | (&('_') '_') | (&('-') '-') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position175, tokenIndex175 := position, tokenIndex
			{
				position176 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l175
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l175
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l175
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l175
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l175
						}
						position++
						break
					}
				}

			l177:
				{
					position178, tokenIndex178 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l178
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l178
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l178
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l178
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l178
							}
							position++
							break
						}
					}

					goto l177
				l178:
					position, tokenIndex = position178, tokenIndex178
				}
				add(ruleIdentifier, position176)
			}
			return true
		l175:
			position, tokenIndex = position175, tokenIndex175
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
				position191 := position
			l192:
				{
					position193, tokenIndex193 := position, tokenIndex
					{
						position194 := position
						{
							position195, tokenIndex195 := position, tokenIndex
							if !_rules[ruleWhitespace]() {
								goto l196
							}
							goto l195
						l196:
							position, tokenIndex = position195, tokenIndex195
							if !_rules[ruleEndOfLine]() {
								goto l193
							}
						}
					l195:
						add(ruleSpace, position194)
					}
					goto l192
				l193:
					position, tokenIndex = position193, tokenIndex193
				}
				add(ruleSpacing, position191)
			}
			return true
		},
		/* 19 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position198 := position
			l199:
				{
					position200, tokenIndex200 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l200
					}
					goto l199
				l200:
					position, tokenIndex = position200, tokenIndex200
				}
				add(ruleWhiteSpacing, position198)
			}
			return true
		},
		/* 20 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position201, tokenIndex201 := position, tokenIndex
			{
				position202 := position
				if !_rules[ruleWhitespace]() {
					goto l201
				}
			l203:
				{
					position204, tokenIndex204 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l204
					}
					goto l203
				l204:
					position, tokenIndex = position204, tokenIndex204
				}
				add(ruleMustWhiteSpacing, position202)
			}
			return true
		l201:
			position, tokenIndex = position201, tokenIndex201
			return false
		},
		/* 21 Equal <- <(Spacing '=' Spacing)> */
		func() bool {
			position205, tokenIndex205 := position, tokenIndex
			{
				position206 := position
				if !_rules[ruleSpacing]() {
					goto l205
				}
				if buffer[position] != rune('=') {
					goto l205
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l205
				}
				add(ruleEqual, position206)
			}
			return true
		l205:
			position, tokenIndex = position205, tokenIndex205
			return false
		},
		/* 22 Space <- <(Whitespace / EndOfLine)> */
		nil,
		/* 23 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position208, tokenIndex208 := position, tokenIndex
			{
				position209 := position
				{
					position210, tokenIndex210 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex = position210, tokenIndex210
					if buffer[position] != rune('\t') {
						goto l208
					}
					position++
				}
			l210:
				add(ruleWhitespace, position209)
			}
			return true
		l208:
			position, tokenIndex = position208, tokenIndex208
			return false
		},
		/* 24 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position212, tokenIndex212 := position, tokenIndex
			{
				position213 := position
				{
					position214, tokenIndex214 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l215
					}
					position++
					if buffer[position] != rune('\n') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex = position214, tokenIndex214
					if buffer[position] != rune('\n') {
						goto l216
					}
					position++
					goto l214
				l216:
					position, tokenIndex = position214, tokenIndex214
					if buffer[position] != rune('\r') {
						goto l212
					}
					position++
				}
			l214:
				add(ruleEndOfLine, position213)
			}
			return true
		l212:
			position, tokenIndex = position212, tokenIndex212
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
