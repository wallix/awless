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
	ruleListWithoutSquareBrackets
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
	ruleHolesStringValue
	ruleHoleWithSuffixValue
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
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
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
	"ListWithoutSquareBrackets",
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
	"HolesStringValue",
	"HoleWithSuffixValue",
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
	"Action18",
	"Action19",
	"Action20",
	"Action21",
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
	rules  [61]func() bool
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
			p.addFirstValueInList()
		case ruleAction10:
			p.lastValueInList()
		case ruleAction11:
			p.addHolesStringParam(text)
		case ruleAction12:
			p.addParamHoleValue(text)
		case ruleAction13:
			p.addHolesStringParam(text)
		case ruleAction14:
			p.addAliasParam(text)
		case ruleAction15:
			p.addStringValue(text)
		case ruleAction16:
			p.addStringValue(text)
		case ruleAction17:
			p.addParamValue(text)
		case ruleAction18:
			p.addParamRefValue(text)
		case ruleAction19:
			p.addParamCidrValue(text)
		case ruleAction20:
			p.addParamIpValue(text)
		case ruleAction21:
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
		/* 10 CompositeValue <- <(ListValue / ListWithoutSquareBrackets / Value)> */
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
						{
							position108, tokenIndex108 := position, tokenIndex
							if !_rules[ruleWhiteSpacing]() {
								goto l108
							}
							if !_rules[ruleValue]() {
								goto l108
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l108
							}
							goto l109
						l108:
							position, tokenIndex = position108, tokenIndex108
						}
					l109:
					l110:
						{
							position111, tokenIndex111 := position, tokenIndex
							if buffer[position] != rune(',') {
								goto l111
							}
							position++
							if !_rules[ruleWhiteSpacing]() {
								goto l111
							}
							if !_rules[ruleValue]() {
								goto l111
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l111
							}
							goto l110
						l111:
							position, tokenIndex = position111, tokenIndex111
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
					{
						position114 := position
						{
							add(ruleAction9, position)
						}
						if !_rules[ruleWhiteSpacing]() {
							goto l113
						}
						if !_rules[ruleValue]() {
							goto l113
						}
						if !_rules[ruleWhiteSpacing]() {
							goto l113
						}
						if buffer[position] != rune(',') {
							goto l113
						}
						position++
						if !_rules[ruleWhiteSpacing]() {
							goto l113
						}
						if !_rules[ruleValue]() {
							goto l113
						}
						if !_rules[ruleWhiteSpacing]() {
							goto l113
						}
					l116:
						{
							position117, tokenIndex117 := position, tokenIndex
							if buffer[position] != rune(',') {
								goto l117
							}
							position++
							if !_rules[ruleWhiteSpacing]() {
								goto l117
							}
							if !_rules[ruleValue]() {
								goto l117
							}
							if !_rules[ruleWhiteSpacing]() {
								goto l117
							}
							goto l116
						l117:
							position, tokenIndex = position117, tokenIndex117
						}
						{
							add(ruleAction10, position)
						}
						add(ruleListWithoutSquareBrackets, position114)
					}
					goto l104
				l113:
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
		/* 11 ListValue <- <(Action7 '[' (WhiteSpacing Value WhiteSpacing)? (',' WhiteSpacing Value WhiteSpacing)* ']' Action8)> */
		nil,
		/* 12 ListWithoutSquareBrackets <- <(Action9 (WhiteSpacing Value WhiteSpacing) (',' WhiteSpacing Value WhiteSpacing)+ Action10)> */
		nil,
		/* 13 NoRefValue <- <((HoleWithSuffixValue Action11) / (HoleValue Action12) / (HolesStringValue Action13) / (AliasValue Action14) / (DoubleQuote CustomTypedValue DoubleQuote) / (SingleQuote CustomTypedValue SingleQuote) / CustomTypedValue / ((&('\'') (SingleQuote <SingleQuotedValue> Action16 SingleQuote)) | (&('"') (DoubleQuote <DoubleQuotedValue> Action15 DoubleQuote)) | (&('*' | '+' | '-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | ';' | '<' | '>' | '@' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | '~') (<OtherParamValue> Action17))))> */
		nil,
		/* 14 Value <- <((RefValue Action18) / NoRefValue)> */
		func() bool {
			position122, tokenIndex122 := position, tokenIndex
			{
				position123 := position
				{
					position124, tokenIndex124 := position, tokenIndex
					{
						position126 := position
						if buffer[position] != rune('$') {
							goto l125
						}
						position++
						{
							position127 := position
							if !_rules[ruleIdentifier]() {
								goto l125
							}
							add(rulePegText, position127)
						}
						add(ruleRefValue, position126)
					}
					{
						add(ruleAction18, position)
					}
					goto l124
				l125:
					position, tokenIndex = position124, tokenIndex124
					{
						position129 := position
						{
							position130, tokenIndex130 := position, tokenIndex
							{
								position132 := position
								{
									position133 := position
									if !_rules[ruleHoleValue]() {
										goto l131
									}
									if !_rules[ruleOtherParamValue]() {
										goto l131
									}
								l134:
									{
										position135, tokenIndex135 := position, tokenIndex
										if !_rules[ruleOtherParamValue]() {
											goto l135
										}
										goto l134
									l135:
										position, tokenIndex = position135, tokenIndex135
									}
								l136:
									{
										position137, tokenIndex137 := position, tokenIndex
										{
											position138, tokenIndex138 := position, tokenIndex
											if !_rules[ruleOtherParamValue]() {
												goto l138
											}
											goto l139
										l138:
											position, tokenIndex = position138, tokenIndex138
										}
									l139:
										if !_rules[ruleHoleValue]() {
											goto l137
										}
										{
											position140, tokenIndex140 := position, tokenIndex
											if !_rules[ruleOtherParamValue]() {
												goto l140
											}
											goto l141
										l140:
											position, tokenIndex = position140, tokenIndex140
										}
									l141:
										goto l136
									l137:
										position, tokenIndex = position137, tokenIndex137
									}
									add(rulePegText, position133)
								}
								add(ruleHoleWithSuffixValue, position132)
							}
							{
								add(ruleAction11, position)
							}
							goto l130
						l131:
							position, tokenIndex = position130, tokenIndex130
							if !_rules[ruleHoleValue]() {
								goto l143
							}
							{
								add(ruleAction12, position)
							}
							goto l130
						l143:
							position, tokenIndex = position130, tokenIndex130
							{
								position146 := position
								{
									position147 := position
									{
										position150, tokenIndex150 := position, tokenIndex
										if !_rules[ruleOtherParamValue]() {
											goto l150
										}
										goto l151
									l150:
										position, tokenIndex = position150, tokenIndex150
									}
								l151:
									if !_rules[ruleHoleValue]() {
										goto l145
									}
									{
										position152, tokenIndex152 := position, tokenIndex
										if !_rules[ruleOtherParamValue]() {
											goto l152
										}
										goto l153
									l152:
										position, tokenIndex = position152, tokenIndex152
									}
								l153:
								l148:
									{
										position149, tokenIndex149 := position, tokenIndex
										{
											position154, tokenIndex154 := position, tokenIndex
											if !_rules[ruleOtherParamValue]() {
												goto l154
											}
											goto l155
										l154:
											position, tokenIndex = position154, tokenIndex154
										}
									l155:
										if !_rules[ruleHoleValue]() {
											goto l149
										}
										{
											position156, tokenIndex156 := position, tokenIndex
											if !_rules[ruleOtherParamValue]() {
												goto l156
											}
											goto l157
										l156:
											position, tokenIndex = position156, tokenIndex156
										}
									l157:
										goto l148
									l149:
										position, tokenIndex = position149, tokenIndex149
									}
									add(rulePegText, position147)
								}
								add(ruleHolesStringValue, position146)
							}
							{
								add(ruleAction13, position)
							}
							goto l130
						l145:
							position, tokenIndex = position130, tokenIndex130
							{
								position160 := position
								{
									position161, tokenIndex161 := position, tokenIndex
									if buffer[position] != rune('@') {
										goto l162
									}
									position++
									{
										position163 := position
										if !_rules[ruleOtherParamValue]() {
											goto l162
										}
										add(rulePegText, position163)
									}
									goto l161
								l162:
									position, tokenIndex = position161, tokenIndex161
									if buffer[position] != rune('@') {
										goto l164
									}
									position++
									if !_rules[ruleDoubleQuote]() {
										goto l164
									}
									{
										position165 := position
										if !_rules[ruleDoubleQuotedValue]() {
											goto l164
										}
										add(rulePegText, position165)
									}
									if !_rules[ruleDoubleQuote]() {
										goto l164
									}
									goto l161
								l164:
									position, tokenIndex = position161, tokenIndex161
									if buffer[position] != rune('@') {
										goto l159
									}
									position++
									if !_rules[ruleSingleQuote]() {
										goto l159
									}
									{
										position166 := position
										if !_rules[ruleSingleQuotedValue]() {
											goto l159
										}
										add(rulePegText, position166)
									}
									if !_rules[ruleSingleQuote]() {
										goto l159
									}
								}
							l161:
								add(ruleAliasValue, position160)
							}
							{
								add(ruleAction14, position)
							}
							goto l130
						l159:
							position, tokenIndex = position130, tokenIndex130
							if !_rules[ruleDoubleQuote]() {
								goto l168
							}
							if !_rules[ruleCustomTypedValue]() {
								goto l168
							}
							if !_rules[ruleDoubleQuote]() {
								goto l168
							}
							goto l130
						l168:
							position, tokenIndex = position130, tokenIndex130
							if !_rules[ruleSingleQuote]() {
								goto l169
							}
							if !_rules[ruleCustomTypedValue]() {
								goto l169
							}
							if !_rules[ruleSingleQuote]() {
								goto l169
							}
							goto l130
						l169:
							position, tokenIndex = position130, tokenIndex130
							if !_rules[ruleCustomTypedValue]() {
								goto l170
							}
							goto l130
						l170:
							position, tokenIndex = position130, tokenIndex130
							{
								switch buffer[position] {
								case '\'':
									if !_rules[ruleSingleQuote]() {
										goto l122
									}
									{
										position172 := position
										if !_rules[ruleSingleQuotedValue]() {
											goto l122
										}
										add(rulePegText, position172)
									}
									{
										add(ruleAction16, position)
									}
									if !_rules[ruleSingleQuote]() {
										goto l122
									}
									break
								case '"':
									if !_rules[ruleDoubleQuote]() {
										goto l122
									}
									{
										position174 := position
										if !_rules[ruleDoubleQuotedValue]() {
											goto l122
										}
										add(rulePegText, position174)
									}
									{
										add(ruleAction15, position)
									}
									if !_rules[ruleDoubleQuote]() {
										goto l122
									}
									break
								default:
									{
										position176 := position
										if !_rules[ruleOtherParamValue]() {
											goto l122
										}
										add(rulePegText, position176)
									}
									{
										add(ruleAction17, position)
									}
									break
								}
							}

						}
					l130:
						add(ruleNoRefValue, position129)
					}
				}
			l124:
				add(ruleValue, position123)
			}
			return true
		l122:
			position, tokenIndex = position122, tokenIndex122
			return false
		},
		/* 15 CustomTypedValue <- <((<CidrValue> Action19) / (<IpValue> Action20) / (<IntRangeValue> Action21))> */
		func() bool {
			position178, tokenIndex178 := position, tokenIndex
			{
				position179 := position
				{
					position180, tokenIndex180 := position, tokenIndex
					{
						position182 := position
						{
							position183 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l181
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
							if buffer[position] != rune('.') {
								goto l181
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l181
							}
							position++
						l186:
							{
								position187, tokenIndex187 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l187
								}
								position++
								goto l186
							l187:
								position, tokenIndex = position187, tokenIndex187
							}
							if buffer[position] != rune('.') {
								goto l181
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l181
							}
							position++
						l188:
							{
								position189, tokenIndex189 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l189
								}
								position++
								goto l188
							l189:
								position, tokenIndex = position189, tokenIndex189
							}
							if buffer[position] != rune('.') {
								goto l181
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l181
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
							if buffer[position] != rune('/') {
								goto l181
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l181
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
							add(ruleCidrValue, position183)
						}
						add(rulePegText, position182)
					}
					{
						add(ruleAction19, position)
					}
					goto l180
				l181:
					position, tokenIndex = position180, tokenIndex180
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
							if buffer[position] != rune('.') {
								goto l195
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l195
							}
							position++
						l200:
							{
								position201, tokenIndex201 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l201
								}
								position++
								goto l200
							l201:
								position, tokenIndex = position201, tokenIndex201
							}
							if buffer[position] != rune('.') {
								goto l195
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l195
							}
							position++
						l202:
							{
								position203, tokenIndex203 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l203
								}
								position++
								goto l202
							l203:
								position, tokenIndex = position203, tokenIndex203
							}
							if buffer[position] != rune('.') {
								goto l195
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l195
							}
							position++
						l204:
							{
								position205, tokenIndex205 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l205
								}
								position++
								goto l204
							l205:
								position, tokenIndex = position205, tokenIndex205
							}
							add(ruleIpValue, position197)
						}
						add(rulePegText, position196)
					}
					{
						add(ruleAction20, position)
					}
					goto l180
				l195:
					position, tokenIndex = position180, tokenIndex180
					{
						position207 := position
						{
							position208 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l178
							}
							position++
						l209:
							{
								position210, tokenIndex210 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l210
								}
								position++
								goto l209
							l210:
								position, tokenIndex = position210, tokenIndex210
							}
							if buffer[position] != rune('-') {
								goto l178
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l178
							}
							position++
						l211:
							{
								position212, tokenIndex212 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l212
								}
								position++
								goto l211
							l212:
								position, tokenIndex = position212, tokenIndex212
							}
							add(ruleIntRangeValue, position208)
						}
						add(rulePegText, position207)
					}
					{
						add(ruleAction21, position)
					}
				}
			l180:
				add(ruleCustomTypedValue, position179)
			}
			return true
		l178:
			position, tokenIndex = position178, tokenIndex178
			return false
		},
		/* 16 OtherParamValue <- <((&('*') '*') | (&('>') '>') | (&('<') '<') | (&('@') '@') | (&('~') '~') | (&(';') ';') | (&('+') '+') | (&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position214, tokenIndex214 := position, tokenIndex
			{
				position215 := position
				{
					switch buffer[position] {
					case '*':
						if buffer[position] != rune('*') {
							goto l214
						}
						position++
						break
					case '>':
						if buffer[position] != rune('>') {
							goto l214
						}
						position++
						break
					case '<':
						if buffer[position] != rune('<') {
							goto l214
						}
						position++
						break
					case '@':
						if buffer[position] != rune('@') {
							goto l214
						}
						position++
						break
					case '~':
						if buffer[position] != rune('~') {
							goto l214
						}
						position++
						break
					case ';':
						if buffer[position] != rune(';') {
							goto l214
						}
						position++
						break
					case '+':
						if buffer[position] != rune('+') {
							goto l214
						}
						position++
						break
					case '/':
						if buffer[position] != rune('/') {
							goto l214
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l214
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l214
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l214
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l214
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l214
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l214
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l214
						}
						position++
						break
					}
				}

			l216:
				{
					position217, tokenIndex217 := position, tokenIndex
					{
						switch buffer[position] {
						case '*':
							if buffer[position] != rune('*') {
								goto l217
							}
							position++
							break
						case '>':
							if buffer[position] != rune('>') {
								goto l217
							}
							position++
							break
						case '<':
							if buffer[position] != rune('<') {
								goto l217
							}
							position++
							break
						case '@':
							if buffer[position] != rune('@') {
								goto l217
							}
							position++
							break
						case '~':
							if buffer[position] != rune('~') {
								goto l217
							}
							position++
							break
						case ';':
							if buffer[position] != rune(';') {
								goto l217
							}
							position++
							break
						case '+':
							if buffer[position] != rune('+') {
								goto l217
							}
							position++
							break
						case '/':
							if buffer[position] != rune('/') {
								goto l217
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
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
						case '.':
							if buffer[position] != rune('.') {
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

					goto l216
				l217:
					position, tokenIndex = position217, tokenIndex217
				}
				add(ruleOtherParamValue, position215)
			}
			return true
		l214:
			position, tokenIndex = position214, tokenIndex214
			return false
		},
		/* 17 DoubleQuotedValue <- <(!'"' .)*> */
		func() bool {
			{
				position221 := position
			l222:
				{
					position223, tokenIndex223 := position, tokenIndex
					{
						position224, tokenIndex224 := position, tokenIndex
						if buffer[position] != rune('"') {
							goto l224
						}
						position++
						goto l223
					l224:
						position, tokenIndex = position224, tokenIndex224
					}
					if !matchDot() {
						goto l223
					}
					goto l222
				l223:
					position, tokenIndex = position223, tokenIndex223
				}
				add(ruleDoubleQuotedValue, position221)
			}
			return true
		},
		/* 18 SingleQuotedValue <- <(!'\'' .)*> */
		func() bool {
			{
				position226 := position
			l227:
				{
					position228, tokenIndex228 := position, tokenIndex
					{
						position229, tokenIndex229 := position, tokenIndex
						if buffer[position] != rune('\'') {
							goto l229
						}
						position++
						goto l228
					l229:
						position, tokenIndex = position229, tokenIndex229
					}
					if !matchDot() {
						goto l228
					}
					goto l227
				l228:
					position, tokenIndex = position228, tokenIndex228
				}
				add(ruleSingleQuotedValue, position226)
			}
			return true
		},
		/* 19 CidrValue <- <([0-9]+ '.' [0-9]+ '.' [0-9]+ '.' [0-9]+ '/' [0-9]+)> */
		nil,
		/* 20 IpValue <- <([0-9]+ '.' [0-9]+ '.' [0-9]+ '.' [0-9]+)> */
		nil,
		/* 21 IntRangeValue <- <([0-9]+ '-' [0-9]+)> */
		nil,
		/* 22 RefValue <- <('$' <Identifier>)> */
		nil,
		/* 23 AliasValue <- <(('@' <OtherParamValue>) / ('@' DoubleQuote <DoubleQuotedValue> DoubleQuote) / ('@' SingleQuote <SingleQuotedValue> SingleQuote))> */
		nil,
		/* 24 HoleValue <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		func() bool {
			position235, tokenIndex235 := position, tokenIndex
			{
				position236 := position
				if buffer[position] != rune('{') {
					goto l235
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l235
				}
				{
					position237 := position
					if !_rules[ruleIdentifier]() {
						goto l235
					}
					add(rulePegText, position237)
				}
				if !_rules[ruleWhiteSpacing]() {
					goto l235
				}
				if buffer[position] != rune('}') {
					goto l235
				}
				position++
				add(ruleHoleValue, position236)
			}
			return true
		l235:
			position, tokenIndex = position235, tokenIndex235
			return false
		},
		/* 25 HolesStringValue <- <<(OtherParamValue? HoleValue OtherParamValue?)+>> */
		nil,
		/* 26 HoleWithSuffixValue <- <<(HoleValue OtherParamValue+ (OtherParamValue? HoleValue OtherParamValue?)*)>> */
		nil,
		/* 27 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)*))> */
		nil,
		/* 28 SingleQuote <- <'\''> */
		func() bool {
			position241, tokenIndex241 := position, tokenIndex
			{
				position242 := position
				if buffer[position] != rune('\'') {
					goto l241
				}
				position++
				add(ruleSingleQuote, position242)
			}
			return true
		l241:
			position, tokenIndex = position241, tokenIndex241
			return false
		},
		/* 29 DoubleQuote <- <'"'> */
		func() bool {
			position243, tokenIndex243 := position, tokenIndex
			{
				position244 := position
				if buffer[position] != rune('"') {
					goto l243
				}
				position++
				add(ruleDoubleQuote, position244)
			}
			return true
		l243:
			position, tokenIndex = position243, tokenIndex243
			return false
		},
		/* 30 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position246 := position
			l247:
				{
					position248, tokenIndex248 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l248
					}
					goto l247
				l248:
					position, tokenIndex = position248, tokenIndex248
				}
				add(ruleWhiteSpacing, position246)
			}
			return true
		},
		/* 31 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position249, tokenIndex249 := position, tokenIndex
			{
				position250 := position
				if !_rules[ruleWhitespace]() {
					goto l249
				}
			l251:
				{
					position252, tokenIndex252 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l252
					}
					goto l251
				l252:
					position, tokenIndex = position252, tokenIndex252
				}
				add(ruleMustWhiteSpacing, position250)
			}
			return true
		l249:
			position, tokenIndex = position249, tokenIndex249
			return false
		},
		/* 32 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position253, tokenIndex253 := position, tokenIndex
			{
				position254 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l253
				}
				if buffer[position] != rune('=') {
					goto l253
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l253
				}
				add(ruleEqual, position254)
			}
			return true
		l253:
			position, tokenIndex = position253, tokenIndex253
			return false
		},
		/* 33 BlankLine <- <(WhiteSpacing EndOfLine)> */
		func() bool {
			position255, tokenIndex255 := position, tokenIndex
			{
				position256 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l255
				}
				if !_rules[ruleEndOfLine]() {
					goto l255
				}
				add(ruleBlankLine, position256)
			}
			return true
		l255:
			position, tokenIndex = position255, tokenIndex255
			return false
		},
		/* 34 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position257, tokenIndex257 := position, tokenIndex
			{
				position258 := position
				{
					position259, tokenIndex259 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l260
					}
					position++
					goto l259
				l260:
					position, tokenIndex = position259, tokenIndex259
					if buffer[position] != rune('\t') {
						goto l257
					}
					position++
				}
			l259:
				add(ruleWhitespace, position258)
			}
			return true
		l257:
			position, tokenIndex = position257, tokenIndex257
			return false
		},
		/* 35 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position261, tokenIndex261 := position, tokenIndex
			{
				position262 := position
				{
					position263, tokenIndex263 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l264
					}
					position++
					if buffer[position] != rune('\n') {
						goto l264
					}
					position++
					goto l263
				l264:
					position, tokenIndex = position263, tokenIndex263
					if buffer[position] != rune('\n') {
						goto l265
					}
					position++
					goto l263
				l265:
					position, tokenIndex = position263, tokenIndex263
					if buffer[position] != rune('\r') {
						goto l261
					}
					position++
				}
			l263:
				add(ruleEndOfLine, position262)
			}
			return true
		l261:
			position, tokenIndex = position261, tokenIndex261
			return false
		},
		/* 36 EndOfFile <- <!.> */
		nil,
		/* 38 Action0 <- <{ p.NewStatement() }> */
		nil,
		/* 39 Action1 <- <{ p.StatementDone() }> */
		nil,
		nil,
		/* 41 Action2 <- <{ p.addDeclarationIdentifier(text) }> */
		nil,
		/* 42 Action3 <- <{ p.addValue() }> */
		nil,
		/* 43 Action4 <- <{ p.addAction(text) }> */
		nil,
		/* 44 Action5 <- <{ p.addEntity(text) }> */
		nil,
		/* 45 Action6 <- <{ p.addParamKey(text) }> */
		nil,
		/* 46 Action7 <- <{  p.addFirstValueInList() }> */
		nil,
		/* 47 Action8 <- <{  p.lastValueInList() }> */
		nil,
		/* 48 Action9 <- <{  p.addFirstValueInList() }> */
		nil,
		/* 49 Action10 <- <{  p.lastValueInList() }> */
		nil,
		/* 50 Action11 <- <{  p.addHolesStringParam(text) }> */
		nil,
		/* 51 Action12 <- <{  p.addParamHoleValue(text) }> */
		nil,
		/* 52 Action13 <- <{  p.addHolesStringParam(text) }> */
		nil,
		/* 53 Action14 <- <{  p.addAliasParam(text) }> */
		nil,
		/* 54 Action15 <- <{ p.addStringValue(text) }> */
		nil,
		/* 55 Action16 <- <{ p.addStringValue(text) }> */
		nil,
		/* 56 Action17 <- <{ p.addParamValue(text) }> */
		nil,
		/* 57 Action18 <- <{  p.addParamRefValue(text) }> */
		nil,
		/* 58 Action19 <- <{ p.addParamCidrValue(text) }> */
		nil,
		/* 59 Action20 <- <{ p.addParamIpValue(text) }> */
		nil,
		/* 60 Action21 <- <{ p.addParamValue(text) }> */
		nil,
	}
	p.rules = _rules
}
