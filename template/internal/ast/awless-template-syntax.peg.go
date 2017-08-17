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
	ruleValueExpr
	ruleCmdExpr
	ruleParams
	ruleParam
	ruleIdentifier
	ruleCompositeValue
	ruleListValue
	ruleNoRefValue
	ruleValue
	ruleCustomTypedValue
	ruleOtherParamValue
	ruleDoubleQuotedValue
	ruleSingleQuotedValue
	ruleCidrValue
	ruleIpValue
	ruleIntRangeValue
	ruleRefValue
	ruleAliasValue
	ruleHoleValue
	ruleComment
	ruleSingleQuote
	ruleDoubleQuote
	ruleWhiteSpacing
	ruleMustWhiteSpacing
	ruleEqual
	ruleBlankLine
	ruleWhitespace
	ruleEndOfLine
	ruleEndOfFile
	ruleAction0
	ruleAction1
	rulePegText
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
	ruleAction16
	ruleAction17
)

var rul3s = [...]string{
	"Unknown",
	"Script",
	"Statement",
	"Action",
	"Entity",
	"Declaration",
	"ValueExpr",
	"CmdExpr",
	"Params",
	"Param",
	"Identifier",
	"CompositeValue",
	"ListValue",
	"NoRefValue",
	"Value",
	"CustomTypedValue",
	"OtherParamValue",
	"DoubleQuotedValue",
	"SingleQuotedValue",
	"CidrValue",
	"IpValue",
	"IntRangeValue",
	"RefValue",
	"AliasValue",
	"HoleValue",
	"Comment",
	"SingleQuote",
	"DoubleQuote",
	"WhiteSpacing",
	"MustWhiteSpacing",
	"Equal",
	"BlankLine",
	"Whitespace",
	"EndOfLine",
	"EndOfFile",
	"Action0",
	"Action1",
	"PegText",
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
	"Action16",
	"Action17",
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
	rules  [54]func() bool
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
			p.NewStatement()
		case ruleAction1:
			p.StatementDone()
		case ruleAction2:
			p.addDeclarationIdentifier(text)
		case ruleAction3:
			p.addValue()
		case ruleAction4:
			p.addAction(text)
		case ruleAction5:
			p.addEntity(text)
		case ruleAction6:
			p.addParamKey(text)
		case ruleAction7:
			p.addFirstValueInList()
		case ruleAction8:
			p.lastValueInList()
		case ruleAction9:
			p.addParamHoleValue(text)
		case ruleAction10:
			p.addAliasParam(text)
		case ruleAction11:
			p.addStringValue(text)
		case ruleAction12:
			p.addStringValue(text)
		case ruleAction13:
			p.addParamValue(text)
		case ruleAction14:
			p.addParamRefValue(text)
		case ruleAction15:
			p.addParamCidrValue(text)
		case ruleAction16:
			p.addParamIpValue(text)
		case ruleAction17:
			p.addParamValue(text)

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
					{
						add(ruleAction0, position)
					}
					if !_rules[ruleWhiteSpacing]() {
						goto l0
					}
					{
						position8, tokenIndex8 := position, tokenIndex
						if !_rules[ruleCmdExpr]() {
							goto l9
						}
						goto l8
					l9:
						position, tokenIndex = position8, tokenIndex8
						{
							position11 := position
							{
								position12 := position
								if !_rules[ruleIdentifier]() {
									goto l10
								}
								add(rulePegText, position12)
							}
							{
								add(ruleAction2, position)
							}
							if !_rules[ruleEqual]() {
								goto l10
							}
							{
								position14, tokenIndex14 := position, tokenIndex
								if !_rules[ruleCmdExpr]() {
									goto l15
								}
								goto l14
							l15:
								position, tokenIndex = position14, tokenIndex14
								{
									position16 := position
									{
										add(ruleAction3, position)
									}
									if !_rules[ruleCompositeValue]() {
										goto l10
									}
									add(ruleValueExpr, position16)
								}
							}
						l14:
							add(ruleDeclaration, position11)
						}
						goto l8
					l10:
						position, tokenIndex = position8, tokenIndex8
						{
							position18 := position
							{
								position19, tokenIndex19 := position, tokenIndex
								if buffer[position] != rune('#') {
									goto l20
								}
								position++
							l21:
								{
									position22, tokenIndex22 := position, tokenIndex
									{
										position23, tokenIndex23 := position, tokenIndex
										if !_rules[ruleEndOfLine]() {
											goto l23
										}
										goto l22
									l23:
										position, tokenIndex = position23, tokenIndex23
									}
									if !matchDot() {
										goto l22
									}
									goto l21
								l22:
									position, tokenIndex = position22, tokenIndex22
								}
								goto l19
							l20:
								position, tokenIndex = position19, tokenIndex19
								if buffer[position] != rune('/') {
									goto l0
								}
								position++
								if buffer[position] != rune('/') {
									goto l0
								}
								position++
							l24:
								{
									position25, tokenIndex25 := position, tokenIndex
									{
										position26, tokenIndex26 := position, tokenIndex
										if !_rules[ruleEndOfLine]() {
											goto l26
										}
										goto l25
									l26:
										position, tokenIndex = position26, tokenIndex26
									}
									if !matchDot() {
										goto l25
									}
									goto l24
								l25:
									position, tokenIndex = position25, tokenIndex25
								}
							}
						l19:
							add(ruleComment, position18)
						}
					}
				l8:
					if !_rules[ruleWhiteSpacing]() {
						goto l0
					}
				l27:
					{
						position28, tokenIndex28 := position, tokenIndex
						if !_rules[ruleEndOfLine]() {
							goto l28
						}
						goto l27
					l28:
						position, tokenIndex = position28, tokenIndex28
					}
					{
						add(ruleAction1, position)
					}
					add(ruleStatement, position6)
				}
			l30:
				{
					position31, tokenIndex31 := position, tokenIndex
					if !_rules[ruleBlankLine]() {
						goto l31
					}
					goto l30
				l31:
					position, tokenIndex = position31, tokenIndex31
				}
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
				l32:
					{
						position33, tokenIndex33 := position, tokenIndex
						if !_rules[ruleBlankLine]() {
							goto l33
						}
						goto l32
					l33:
						position, tokenIndex = position33, tokenIndex33
					}
					{
						position34 := position
						{
							add(ruleAction0, position)
						}
						if !_rules[ruleWhiteSpacing]() {
							goto l3
						}
						{
							position36, tokenIndex36 := position, tokenIndex
							if !_rules[ruleCmdExpr]() {
								goto l37
							}
							goto l36
						l37:
							position, tokenIndex = position36, tokenIndex36
							{
								position39 := position
								{
									position40 := position
									if !_rules[ruleIdentifier]() {
										goto l38
									}
									add(rulePegText, position40)
								}
								{
									add(ruleAction2, position)
								}
								if !_rules[ruleEqual]() {
									goto l38
								}
								{
									position42, tokenIndex42 := position, tokenIndex
									if !_rules[ruleCmdExpr]() {
										goto l43
									}
									goto l42
								l43:
									position, tokenIndex = position42, tokenIndex42
									{
										position44 := position
										{
											add(ruleAction3, position)
										}
										if !_rules[ruleCompositeValue]() {
											goto l38
										}
										add(ruleValueExpr, position44)
									}
								}
							l42:
								add(ruleDeclaration, position39)
							}
							goto l36
						l38:
							position, tokenIndex = position36, tokenIndex36
							{
								position46 := position
								{
									position47, tokenIndex47 := position, tokenIndex
									if buffer[position] != rune('#') {
										goto l48
									}
									position++
								l49:
									{
										position50, tokenIndex50 := position, tokenIndex
										{
											position51, tokenIndex51 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l51
											}
											goto l50
										l51:
											position, tokenIndex = position51, tokenIndex51
										}
										if !matchDot() {
											goto l50
										}
										goto l49
									l50:
										position, tokenIndex = position50, tokenIndex50
									}
									goto l47
								l48:
									position, tokenIndex = position47, tokenIndex47
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
									if buffer[position] != rune('/') {
										goto l3
									}
									position++
								l52:
									{
										position53, tokenIndex53 := position, tokenIndex
										{
											position54, tokenIndex54 := position, tokenIndex
											if !_rules[ruleEndOfLine]() {
												goto l54
											}
											goto l53
										l54:
											position, tokenIndex = position54, tokenIndex54
										}
										if !matchDot() {
											goto l53
										}
										goto l52
									l53:
										position, tokenIndex = position53, tokenIndex53
									}
								}
							l47:
								add(ruleComment, position46)
							}
						}
					l36:
						if !_rules[ruleWhiteSpacing]() {
							goto l3
						}
					l55:
						{
							position56, tokenIndex56 := position, tokenIndex
							if !_rules[ruleEndOfLine]() {
								goto l56
							}
							goto l55
						l56:
							position, tokenIndex = position56, tokenIndex56
						}
						{
							add(ruleAction1, position)
						}
						add(ruleStatement, position34)
					}
				l58:
					{
						position59, tokenIndex59 := position, tokenIndex
						if !_rules[ruleBlankLine]() {
							goto l59
						}
						goto l58
					l59:
						position, tokenIndex = position59, tokenIndex59
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				if !_rules[ruleWhiteSpacing]() {
					goto l0
				}
				{
					position60 := position
					{
						position61, tokenIndex61 := position, tokenIndex
						if !matchDot() {
							goto l61
						}
						goto l0
					l61:
						position, tokenIndex = position61, tokenIndex61
					}
					add(ruleEndOfFile, position60)
				}
				add(ruleScript, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Statement <- <(Action0 WhiteSpacing (CmdExpr / Declaration / Comment) WhiteSpacing EndOfLine* Action1)> */
		nil,
		/* 2 Action <- <[a-z]+> */
		nil,
		/* 3 Entity <- <([a-z] / [0-9])+> */
		nil,
		/* 4 Declaration <- <(<Identifier> Action2 Equal (CmdExpr / ValueExpr))> */
		nil,
		/* 5 ValueExpr <- <(Action3 CompositeValue)> */
		nil,
		/* 6 CmdExpr <- <(<Action> Action4 MustWhiteSpacing <Entity> Action5 (MustWhiteSpacing Params)?)> */
		func() bool {
			position67, tokenIndex67 := position, tokenIndex
			{
				position68 := position
				{
					position69 := position
					{
						position70 := position
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l67
						}
						position++
					l71:
						{
							position72, tokenIndex72 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l72
							}
							position++
							goto l71
						l72:
							position, tokenIndex = position72, tokenIndex72
						}
						add(ruleAction, position70)
					}
					add(rulePegText, position69)
				}
				{
					add(ruleAction4, position)
				}
				if !_rules[ruleMustWhiteSpacing]() {
					goto l67
				}
				{
					position74 := position
					{
						position75 := position
						{
							position78, tokenIndex78 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l79
							}
							position++
							goto l78
						l79:
							position, tokenIndex = position78, tokenIndex78
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l67
							}
							position++
						}
					l78:
					l76:
						{
							position77, tokenIndex77 := position, tokenIndex
							{
								position80, tokenIndex80 := position, tokenIndex
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l81
								}
								position++
								goto l80
							l81:
								position, tokenIndex = position80, tokenIndex80
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l77
								}
								position++
							}
						l80:
							goto l76
						l77:
							position, tokenIndex = position77, tokenIndex77
						}
						add(ruleEntity, position75)
					}
					add(rulePegText, position74)
				}
				{
					add(ruleAction5, position)
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
								add(ruleAction6, position)
							}
							if !_rules[ruleEqual]() {
								goto l83
							}
							if !_rules[ruleCompositeValue]() {
								goto l83
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
								position91 := position
								{
									position92 := position
									if !_rules[ruleIdentifier]() {
										goto l87
									}
									add(rulePegText, position92)
								}
								{
									add(ruleAction6, position)
								}
								if !_rules[ruleEqual]() {
									goto l87
								}
								if !_rules[ruleCompositeValue]() {
									goto l87
								}
								if !_rules[ruleWhiteSpacing]() {
									goto l87
								}
								add(ruleParam, position91)
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
				add(ruleCmdExpr, position68)
			}
			return true
		l67:
			position, tokenIndex = position67, tokenIndex67
			return false
		},
		/* 7 Params <- <Param+> */
		nil,
		/* 8 Param <- <(<Identifier> Action6 Equal CompositeValue WhiteSpacing)> */
		nil,
		/* 9 Identifier <- <((&('.') '.') | (&('_') '_') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position96, tokenIndex96 := position, tokenIndex
			{
				position97 := position
				{
					switch buffer[position] {
					case '.':
						if buffer[position] != rune('.') {
							goto l96
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
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

			l98:
				{
					position99, tokenIndex99 := position, tokenIndex
					{
						switch buffer[position] {
						case '.':
							if buffer[position] != rune('.') {
								goto l99
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l99
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l99
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l99
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l99
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l99
							}
							position++
							break
						}
					}

					goto l98
				l99:
					position, tokenIndex = position99, tokenIndex99
				}
				add(ruleIdentifier, position97)
			}
			return true
		l96:
			position, tokenIndex = position96, tokenIndex96
			return false
		},
		/* 10 CompositeValue <- <(ListValue / Value)> */
		func() bool {
			position102, tokenIndex102 := position, tokenIndex
			{
				position103 := position
				{
					position104, tokenIndex104 := position, tokenIndex
					{
						position106 := position
						{
							add(ruleAction7, position)
						}
						if buffer[position] != rune('[') {
							goto l105
						}
						position++
						if !_rules[ruleWhiteSpacing]() {
							goto l105
						}
						if !_rules[ruleValue]() {
							goto l105
						}
						if !_rules[ruleWhiteSpacing]() {
							goto l105
						}
					l108:
						{
							position109, tokenIndex109 := position, tokenIndex
							if buffer[position] != rune(',') {
								goto l109
							}
							position++
							if !_rules[ruleWhiteSpacing]() {
								goto l109
							}
							if !_rules[ruleValue]() {
								goto l109
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l109
							}
							goto l108
						l109:
							position, tokenIndex = position109, tokenIndex109
						}
						if buffer[position] != rune(']') {
							goto l105
						}
						position++
						{
							add(ruleAction8, position)
						}
						add(ruleListValue, position106)
					}
					goto l104
				l105:
					position, tokenIndex = position104, tokenIndex104
					if !_rules[ruleValue]() {
						goto l102
					}
				}
			l104:
				add(ruleCompositeValue, position103)
			}
			return true
		l102:
			position, tokenIndex = position102, tokenIndex102
			return false
		},
		/* 11 ListValue <- <(Action7 '[' WhiteSpacing Value WhiteSpacing (',' WhiteSpacing Value WhiteSpacing)* ']' Action8)> */
		nil,
		/* 12 NoRefValue <- <((AliasValue Action10) / (DoubleQuote CustomTypedValue DoubleQuote) / (SingleQuote CustomTypedValue SingleQuote) / CustomTypedValue / ((&('\'') (SingleQuote <SingleQuotedValue> Action12 SingleQuote)) | (&('"') (DoubleQuote <DoubleQuotedValue> Action11 DoubleQuote)) | (&('{') (HoleValue Action9)) | (&('*' | '+' | '-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | ';' | '<' | '>' | '@' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | '~') (<OtherParamValue> Action13))))> */
		nil,
		/* 13 Value <- <((RefValue Action14) / NoRefValue)> */
		func() bool {
			position113, tokenIndex113 := position, tokenIndex
			{
				position114 := position
				{
					position115, tokenIndex115 := position, tokenIndex
					{
						position117 := position
						if buffer[position] != rune('$') {
							goto l116
						}
						position++
						{
							position118 := position
							if !_rules[ruleIdentifier]() {
								goto l116
							}
							add(rulePegText, position118)
						}
						add(ruleRefValue, position117)
					}
					{
						add(ruleAction14, position)
					}
					goto l115
				l116:
					position, tokenIndex = position115, tokenIndex115
					{
						position120 := position
						{
							position121, tokenIndex121 := position, tokenIndex
							{
								position123 := position
								{
									position124, tokenIndex124 := position, tokenIndex
									if buffer[position] != rune('@') {
										goto l125
									}
									position++
									{
										position126 := position
										if !_rules[ruleOtherParamValue]() {
											goto l125
										}
										add(rulePegText, position126)
									}
									goto l124
								l125:
									position, tokenIndex = position124, tokenIndex124
									if buffer[position] != rune('@') {
										goto l127
									}
									position++
									if !_rules[ruleDoubleQuote]() {
										goto l127
									}
									{
										position128 := position
										if !_rules[ruleDoubleQuotedValue]() {
											goto l127
										}
										add(rulePegText, position128)
									}
									if !_rules[ruleDoubleQuote]() {
										goto l127
									}
									goto l124
								l127:
									position, tokenIndex = position124, tokenIndex124
									if buffer[position] != rune('@') {
										goto l122
									}
									position++
									if !_rules[ruleSingleQuote]() {
										goto l122
									}
									{
										position129 := position
										if !_rules[ruleSingleQuotedValue]() {
											goto l122
										}
										add(rulePegText, position129)
									}
									if !_rules[ruleSingleQuote]() {
										goto l122
									}
								}
							l124:
								add(ruleAliasValue, position123)
							}
							{
								add(ruleAction10, position)
							}
							goto l121
						l122:
							position, tokenIndex = position121, tokenIndex121
							if !_rules[ruleDoubleQuote]() {
								goto l131
							}
							if !_rules[ruleCustomTypedValue]() {
								goto l131
							}
							if !_rules[ruleDoubleQuote]() {
								goto l131
							}
							goto l121
						l131:
							position, tokenIndex = position121, tokenIndex121
							if !_rules[ruleSingleQuote]() {
								goto l132
							}
							if !_rules[ruleCustomTypedValue]() {
								goto l132
							}
							if !_rules[ruleSingleQuote]() {
								goto l132
							}
							goto l121
						l132:
							position, tokenIndex = position121, tokenIndex121
							if !_rules[ruleCustomTypedValue]() {
								goto l133
							}
							goto l121
						l133:
							position, tokenIndex = position121, tokenIndex121
							{
								switch buffer[position] {
								case '\'':
									if !_rules[ruleSingleQuote]() {
										goto l113
									}
									{
										position135 := position
										if !_rules[ruleSingleQuotedValue]() {
											goto l113
										}
										add(rulePegText, position135)
									}
									{
										add(ruleAction12, position)
									}
									if !_rules[ruleSingleQuote]() {
										goto l113
									}
									break
								case '"':
									if !_rules[ruleDoubleQuote]() {
										goto l113
									}
									{
										position137 := position
										if !_rules[ruleDoubleQuotedValue]() {
											goto l113
										}
										add(rulePegText, position137)
									}
									{
										add(ruleAction11, position)
									}
									if !_rules[ruleDoubleQuote]() {
										goto l113
									}
									break
								case '{':
									{
										position139 := position
										if buffer[position] != rune('{') {
											goto l113
										}
										position++
										if !_rules[ruleWhiteSpacing]() {
											goto l113
										}
										{
											position140 := position
											if !_rules[ruleIdentifier]() {
												goto l113
											}
											add(rulePegText, position140)
										}
										if !_rules[ruleWhiteSpacing]() {
											goto l113
										}
										if buffer[position] != rune('}') {
											goto l113
										}
										position++
										add(ruleHoleValue, position139)
									}
									{
										add(ruleAction9, position)
									}
									break
								default:
									{
										position142 := position
										if !_rules[ruleOtherParamValue]() {
											goto l113
										}
										add(rulePegText, position142)
									}
									{
										add(ruleAction13, position)
									}
									break
								}
							}

						}
					l121:
						add(ruleNoRefValue, position120)
					}
				}
			l115:
				add(ruleValue, position114)
			}
			return true
		l113:
			position, tokenIndex = position113, tokenIndex113
			return false
		},
		/* 14 CustomTypedValue <- <((<CidrValue> Action15) / (<IpValue> Action16) / (<IntRangeValue> Action17))> */
		func() bool {
			position144, tokenIndex144 := position, tokenIndex
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
							if buffer[position] != rune('.') {
								goto l147
							}
							position++
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
							if buffer[position] != rune('.') {
								goto l147
							}
							position++
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
							if buffer[position] != rune('.') {
								goto l147
							}
							position++
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
						add(ruleAction15, position)
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
							if buffer[position] != rune('.') {
								goto l161
							}
							position++
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
							if buffer[position] != rune('.') {
								goto l161
							}
							position++
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
							if buffer[position] != rune('.') {
								goto l161
							}
							position++
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
						add(ruleAction16, position)
					}
					goto l146
				l161:
					position, tokenIndex = position146, tokenIndex146
					{
						position173 := position
						{
							position174 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l144
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
							if buffer[position] != rune('-') {
								goto l144
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l144
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
							add(ruleIntRangeValue, position174)
						}
						add(rulePegText, position173)
					}
					{
						add(ruleAction17, position)
					}
				}
			l146:
				add(ruleCustomTypedValue, position145)
			}
			return true
		l144:
			position, tokenIndex = position144, tokenIndex144
			return false
		},
		/* 15 OtherParamValue <- <((&('*') '*') | (&('>') '>') | (&('<') '<') | (&('@') '@') | (&('~') '~') | (&(';') ';') | (&('+') '+') | (&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position180, tokenIndex180 := position, tokenIndex
			{
				position181 := position
				{
					switch buffer[position] {
					case '*':
						if buffer[position] != rune('*') {
							goto l180
						}
						position++
						break
					case '>':
						if buffer[position] != rune('>') {
							goto l180
						}
						position++
						break
					case '<':
						if buffer[position] != rune('<') {
							goto l180
						}
						position++
						break
					case '@':
						if buffer[position] != rune('@') {
							goto l180
						}
						position++
						break
					case '~':
						if buffer[position] != rune('~') {
							goto l180
						}
						position++
						break
					case ';':
						if buffer[position] != rune(';') {
							goto l180
						}
						position++
						break
					case '+':
						if buffer[position] != rune('+') {
							goto l180
						}
						position++
						break
					case '/':
						if buffer[position] != rune('/') {
							goto l180
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l180
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l180
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l180
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l180
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l180
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l180
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l180
						}
						position++
						break
					}
				}

			l182:
				{
					position183, tokenIndex183 := position, tokenIndex
					{
						switch buffer[position] {
						case '*':
							if buffer[position] != rune('*') {
								goto l183
							}
							position++
							break
						case '>':
							if buffer[position] != rune('>') {
								goto l183
							}
							position++
							break
						case '<':
							if buffer[position] != rune('<') {
								goto l183
							}
							position++
							break
						case '@':
							if buffer[position] != rune('@') {
								goto l183
							}
							position++
							break
						case '~':
							if buffer[position] != rune('~') {
								goto l183
							}
							position++
							break
						case ';':
							if buffer[position] != rune(';') {
								goto l183
							}
							position++
							break
						case '+':
							if buffer[position] != rune('+') {
								goto l183
							}
							position++
							break
						case '/':
							if buffer[position] != rune('/') {
								goto l183
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l183
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l183
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l183
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l183
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l183
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l183
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l183
							}
							position++
							break
						}
					}

					goto l182
				l183:
					position, tokenIndex = position183, tokenIndex183
				}
				add(ruleOtherParamValue, position181)
			}
			return true
		l180:
			position, tokenIndex = position180, tokenIndex180
			return false
		},
		/* 16 DoubleQuotedValue <- <(!'"' .)*> */
		func() bool {
			{
				position187 := position
			l188:
				{
					position189, tokenIndex189 := position, tokenIndex
					{
						position190, tokenIndex190 := position, tokenIndex
						if buffer[position] != rune('"') {
							goto l190
						}
						position++
						goto l189
					l190:
						position, tokenIndex = position190, tokenIndex190
					}
					if !matchDot() {
						goto l189
					}
					goto l188
				l189:
					position, tokenIndex = position189, tokenIndex189
				}
				add(ruleDoubleQuotedValue, position187)
			}
			return true
		},
		/* 17 SingleQuotedValue <- <(!'\'' .)*> */
		func() bool {
			{
				position192 := position
			l193:
				{
					position194, tokenIndex194 := position, tokenIndex
					{
						position195, tokenIndex195 := position, tokenIndex
						if buffer[position] != rune('\'') {
							goto l195
						}
						position++
						goto l194
					l195:
						position, tokenIndex = position195, tokenIndex195
					}
					if !matchDot() {
						goto l194
					}
					goto l193
				l194:
					position, tokenIndex = position194, tokenIndex194
				}
				add(ruleSingleQuotedValue, position192)
			}
			return true
		},
		/* 18 CidrValue <- <([0-9]+ '.' [0-9]+ '.' [0-9]+ '.' [0-9]+ '/' [0-9]+)> */
		nil,
		/* 19 IpValue <- <([0-9]+ '.' [0-9]+ '.' [0-9]+ '.' [0-9]+)> */
		nil,
		/* 20 IntRangeValue <- <([0-9]+ '-' [0-9]+)> */
		nil,
		/* 21 RefValue <- <('$' <Identifier>)> */
		nil,
		/* 22 AliasValue <- <(('@' <OtherParamValue>) / ('@' DoubleQuote <DoubleQuotedValue> DoubleQuote) / ('@' SingleQuote <SingleQuotedValue> SingleQuote))> */
		nil,
		/* 23 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 24 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)*))> */
		nil,
		/* 25 SingleQuote <- <'\''> */
		func() bool {
			position203, tokenIndex203 := position, tokenIndex
			{
				position204 := position
				if buffer[position] != rune('\'') {
					goto l203
				}
				position++
				add(ruleSingleQuote, position204)
			}
			return true
		l203:
			position, tokenIndex = position203, tokenIndex203
			return false
		},
		/* 26 DoubleQuote <- <'"'> */
		func() bool {
			position205, tokenIndex205 := position, tokenIndex
			{
				position206 := position
				if buffer[position] != rune('"') {
					goto l205
				}
				position++
				add(ruleDoubleQuote, position206)
			}
			return true
		l205:
			position, tokenIndex = position205, tokenIndex205
			return false
		},
		/* 27 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position208 := position
			l209:
				{
					position210, tokenIndex210 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l210
					}
					goto l209
				l210:
					position, tokenIndex = position210, tokenIndex210
				}
				add(ruleWhiteSpacing, position208)
			}
			return true
		},
		/* 28 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position211, tokenIndex211 := position, tokenIndex
			{
				position212 := position
				if !_rules[ruleWhitespace]() {
					goto l211
				}
			l213:
				{
					position214, tokenIndex214 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l214
					}
					goto l213
				l214:
					position, tokenIndex = position214, tokenIndex214
				}
				add(ruleMustWhiteSpacing, position212)
			}
			return true
		l211:
			position, tokenIndex = position211, tokenIndex211
			return false
		},
		/* 29 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position215, tokenIndex215 := position, tokenIndex
			{
				position216 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l215
				}
				if buffer[position] != rune('=') {
					goto l215
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l215
				}
				add(ruleEqual, position216)
			}
			return true
		l215:
			position, tokenIndex = position215, tokenIndex215
			return false
		},
		/* 30 BlankLine <- <(WhiteSpacing EndOfLine)> */
		func() bool {
			position217, tokenIndex217 := position, tokenIndex
			{
				position218 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l217
				}
				if !_rules[ruleEndOfLine]() {
					goto l217
				}
				add(ruleBlankLine, position218)
			}
			return true
		l217:
			position, tokenIndex = position217, tokenIndex217
			return false
		},
		/* 31 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position219, tokenIndex219 := position, tokenIndex
			{
				position220 := position
				{
					position221, tokenIndex221 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l222
					}
					position++
					goto l221
				l222:
					position, tokenIndex = position221, tokenIndex221
					if buffer[position] != rune('\t') {
						goto l219
					}
					position++
				}
			l221:
				add(ruleWhitespace, position220)
			}
			return true
		l219:
			position, tokenIndex = position219, tokenIndex219
			return false
		},
		/* 32 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position223, tokenIndex223 := position, tokenIndex
			{
				position224 := position
				{
					position225, tokenIndex225 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l226
					}
					position++
					if buffer[position] != rune('\n') {
						goto l226
					}
					position++
					goto l225
				l226:
					position, tokenIndex = position225, tokenIndex225
					if buffer[position] != rune('\n') {
						goto l227
					}
					position++
					goto l225
				l227:
					position, tokenIndex = position225, tokenIndex225
					if buffer[position] != rune('\r') {
						goto l223
					}
					position++
				}
			l225:
				add(ruleEndOfLine, position224)
			}
			return true
		l223:
			position, tokenIndex = position223, tokenIndex223
			return false
		},
		/* 33 EndOfFile <- <!.> */
		nil,
		/* 35 Action0 <- <{ p.NewStatement() }> */
		nil,
		/* 36 Action1 <- <{ p.StatementDone() }> */
		nil,
		nil,
		/* 38 Action2 <- <{ p.addDeclarationIdentifier(text) }> */
		nil,
		/* 39 Action3 <- <{ p.addValue() }> */
		nil,
		/* 40 Action4 <- <{ p.addAction(text) }> */
		nil,
		/* 41 Action5 <- <{ p.addEntity(text) }> */
		nil,
		/* 42 Action6 <- <{ p.addParamKey(text) }> */
		nil,
		/* 43 Action7 <- <{  p.addFirstValueInList() }> */
		nil,
		/* 44 Action8 <- <{  p.lastValueInList() }> */
		nil,
		/* 45 Action9 <- <{  p.addParamHoleValue(text) }> */
		nil,
		/* 46 Action10 <- <{  p.addAliasParam(text) }> */
		nil,
		/* 47 Action11 <- <{ p.addStringValue(text) }> */
		nil,
		/* 48 Action12 <- <{ p.addStringValue(text) }> */
		nil,
		/* 49 Action13 <- <{ p.addParamValue(text) }> */
		nil,
		/* 50 Action14 <- <{  p.addParamRefValue(text) }> */
		nil,
		/* 51 Action15 <- <{ p.addParamCidrValue(text) }> */
		nil,
		/* 52 Action16 <- <{ p.addParamIpValue(text) }> */
		nil,
		/* 53 Action17 <- <{ p.addParamValue(text) }> */
		nil,
	}
	p.rules = _rules
}
