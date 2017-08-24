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
	rules  [58]func() bool
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
			p.addHolesStringParam(text)
		case ruleAction10:
			p.addParamHoleValue(text)
		case ruleAction11:
			p.addHolesStringParam(text)
		case ruleAction12:
			p.addAliasParam(text)
		case ruleAction13:
			p.addStringValue(text)
		case ruleAction14:
			p.addStringValue(text)
		case ruleAction15:
			p.addParamValue(text)
		case ruleAction16:
			p.addParamRefValue(text)
		case ruleAction17:
			p.addParamCidrValue(text)
		case ruleAction18:
			p.addParamIpValue(text)
		case ruleAction19:
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
		/* 12 NoRefValue <- <((HoleWithSuffixValue Action9) / (HoleValue Action10) / (HolesStringValue Action11) / (AliasValue Action12) / (DoubleQuote CustomTypedValue DoubleQuote) / (SingleQuote CustomTypedValue SingleQuote) / CustomTypedValue / ((&('\'') (SingleQuote <SingleQuotedValue> Action14 SingleQuote)) | (&('"') (DoubleQuote <DoubleQuotedValue> Action13 DoubleQuote)) | (&('*' | '+' | '-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | ';' | '<' | '>' | '@' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | '~') (<OtherParamValue> Action15))))> */
		nil,
		/* 13 Value <- <((RefValue Action16) / NoRefValue)> */
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
						add(ruleAction16, position)
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
									position124 := position
									if !_rules[ruleHoleValue]() {
										goto l122
									}
									if !_rules[ruleOtherParamValue]() {
										goto l122
									}
								l125:
									{
										position126, tokenIndex126 := position, tokenIndex
										if !_rules[ruleOtherParamValue]() {
											goto l126
										}
										goto l125
									l126:
										position, tokenIndex = position126, tokenIndex126
									}
								l127:
									{
										position128, tokenIndex128 := position, tokenIndex
										{
											position129, tokenIndex129 := position, tokenIndex
											if !_rules[ruleOtherParamValue]() {
												goto l129
											}
											goto l130
										l129:
											position, tokenIndex = position129, tokenIndex129
										}
									l130:
										if !_rules[ruleHoleValue]() {
											goto l128
										}
										{
											position131, tokenIndex131 := position, tokenIndex
											if !_rules[ruleOtherParamValue]() {
												goto l131
											}
											goto l132
										l131:
											position, tokenIndex = position131, tokenIndex131
										}
									l132:
										goto l127
									l128:
										position, tokenIndex = position128, tokenIndex128
									}
									add(rulePegText, position124)
								}
								add(ruleHoleWithSuffixValue, position123)
							}
							{
								add(ruleAction9, position)
							}
							goto l121
						l122:
							position, tokenIndex = position121, tokenIndex121
							if !_rules[ruleHoleValue]() {
								goto l134
							}
							{
								add(ruleAction10, position)
							}
							goto l121
						l134:
							position, tokenIndex = position121, tokenIndex121
							{
								position137 := position
								{
									position138 := position
									{
										position141, tokenIndex141 := position, tokenIndex
										if !_rules[ruleOtherParamValue]() {
											goto l141
										}
										goto l142
									l141:
										position, tokenIndex = position141, tokenIndex141
									}
								l142:
									if !_rules[ruleHoleValue]() {
										goto l136
									}
									{
										position143, tokenIndex143 := position, tokenIndex
										if !_rules[ruleOtherParamValue]() {
											goto l143
										}
										goto l144
									l143:
										position, tokenIndex = position143, tokenIndex143
									}
								l144:
								l139:
									{
										position140, tokenIndex140 := position, tokenIndex
										{
											position145, tokenIndex145 := position, tokenIndex
											if !_rules[ruleOtherParamValue]() {
												goto l145
											}
											goto l146
										l145:
											position, tokenIndex = position145, tokenIndex145
										}
									l146:
										if !_rules[ruleHoleValue]() {
											goto l140
										}
										{
											position147, tokenIndex147 := position, tokenIndex
											if !_rules[ruleOtherParamValue]() {
												goto l147
											}
											goto l148
										l147:
											position, tokenIndex = position147, tokenIndex147
										}
									l148:
										goto l139
									l140:
										position, tokenIndex = position140, tokenIndex140
									}
									add(rulePegText, position138)
								}
								add(ruleHolesStringValue, position137)
							}
							{
								add(ruleAction11, position)
							}
							goto l121
						l136:
							position, tokenIndex = position121, tokenIndex121
							{
								position151 := position
								{
									position152, tokenIndex152 := position, tokenIndex
									if buffer[position] != rune('@') {
										goto l153
									}
									position++
									{
										position154 := position
										if !_rules[ruleOtherParamValue]() {
											goto l153
										}
										add(rulePegText, position154)
									}
									goto l152
								l153:
									position, tokenIndex = position152, tokenIndex152
									if buffer[position] != rune('@') {
										goto l155
									}
									position++
									if !_rules[ruleDoubleQuote]() {
										goto l155
									}
									{
										position156 := position
										if !_rules[ruleDoubleQuotedValue]() {
											goto l155
										}
										add(rulePegText, position156)
									}
									if !_rules[ruleDoubleQuote]() {
										goto l155
									}
									goto l152
								l155:
									position, tokenIndex = position152, tokenIndex152
									if buffer[position] != rune('@') {
										goto l150
									}
									position++
									if !_rules[ruleSingleQuote]() {
										goto l150
									}
									{
										position157 := position
										if !_rules[ruleSingleQuotedValue]() {
											goto l150
										}
										add(rulePegText, position157)
									}
									if !_rules[ruleSingleQuote]() {
										goto l150
									}
								}
							l152:
								add(ruleAliasValue, position151)
							}
							{
								add(ruleAction12, position)
							}
							goto l121
						l150:
							position, tokenIndex = position121, tokenIndex121
							if !_rules[ruleDoubleQuote]() {
								goto l159
							}
							if !_rules[ruleCustomTypedValue]() {
								goto l159
							}
							if !_rules[ruleDoubleQuote]() {
								goto l159
							}
							goto l121
						l159:
							position, tokenIndex = position121, tokenIndex121
							if !_rules[ruleSingleQuote]() {
								goto l160
							}
							if !_rules[ruleCustomTypedValue]() {
								goto l160
							}
							if !_rules[ruleSingleQuote]() {
								goto l160
							}
							goto l121
						l160:
							position, tokenIndex = position121, tokenIndex121
							if !_rules[ruleCustomTypedValue]() {
								goto l161
							}
							goto l121
						l161:
							position, tokenIndex = position121, tokenIndex121
							{
								switch buffer[position] {
								case '\'':
									if !_rules[ruleSingleQuote]() {
										goto l113
									}
									{
										position163 := position
										if !_rules[ruleSingleQuotedValue]() {
											goto l113
										}
										add(rulePegText, position163)
									}
									{
										add(ruleAction14, position)
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
										position165 := position
										if !_rules[ruleDoubleQuotedValue]() {
											goto l113
										}
										add(rulePegText, position165)
									}
									{
										add(ruleAction13, position)
									}
									if !_rules[ruleDoubleQuote]() {
										goto l113
									}
									break
								default:
									{
										position167 := position
										if !_rules[ruleOtherParamValue]() {
											goto l113
										}
										add(rulePegText, position167)
									}
									{
										add(ruleAction15, position)
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
		/* 14 CustomTypedValue <- <((<CidrValue> Action17) / (<IpValue> Action18) / (<IntRangeValue> Action19))> */
		func() bool {
			position169, tokenIndex169 := position, tokenIndex
			{
				position170 := position
				{
					position171, tokenIndex171 := position, tokenIndex
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
							if buffer[position] != rune('.') {
								goto l172
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l172
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
							if buffer[position] != rune('.') {
								goto l172
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l172
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
							if buffer[position] != rune('.') {
								goto l172
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l172
							}
							position++
						l181:
							{
								position182, tokenIndex182 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l182
								}
								position++
								goto l181
							l182:
								position, tokenIndex = position182, tokenIndex182
							}
							if buffer[position] != rune('/') {
								goto l172
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l172
							}
							position++
						l183:
							{
								position184, tokenIndex184 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l184
								}
								position++
								goto l183
							l184:
								position, tokenIndex = position184, tokenIndex184
							}
							add(ruleCidrValue, position174)
						}
						add(rulePegText, position173)
					}
					{
						add(ruleAction17, position)
					}
					goto l171
				l172:
					position, tokenIndex = position171, tokenIndex171
					{
						position187 := position
						{
							position188 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l186
							}
							position++
						l189:
							{
								position190, tokenIndex190 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l190
								}
								position++
								goto l189
							l190:
								position, tokenIndex = position190, tokenIndex190
							}
							if buffer[position] != rune('.') {
								goto l186
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l186
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
							if buffer[position] != rune('.') {
								goto l186
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l186
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
							if buffer[position] != rune('.') {
								goto l186
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l186
							}
							position++
						l195:
							{
								position196, tokenIndex196 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l196
								}
								position++
								goto l195
							l196:
								position, tokenIndex = position196, tokenIndex196
							}
							add(ruleIpValue, position188)
						}
						add(rulePegText, position187)
					}
					{
						add(ruleAction18, position)
					}
					goto l171
				l186:
					position, tokenIndex = position171, tokenIndex171
					{
						position198 := position
						{
							position199 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l169
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
							if buffer[position] != rune('-') {
								goto l169
							}
							position++
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l169
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
							add(ruleIntRangeValue, position199)
						}
						add(rulePegText, position198)
					}
					{
						add(ruleAction19, position)
					}
				}
			l171:
				add(ruleCustomTypedValue, position170)
			}
			return true
		l169:
			position, tokenIndex = position169, tokenIndex169
			return false
		},
		/* 15 OtherParamValue <- <((&('*') '*') | (&('>') '>') | (&('<') '<') | (&('@') '@') | (&('~') '~') | (&(';') ';') | (&('+') '+') | (&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position205, tokenIndex205 := position, tokenIndex
			{
				position206 := position
				{
					switch buffer[position] {
					case '*':
						if buffer[position] != rune('*') {
							goto l205
						}
						position++
						break
					case '>':
						if buffer[position] != rune('>') {
							goto l205
						}
						position++
						break
					case '<':
						if buffer[position] != rune('<') {
							goto l205
						}
						position++
						break
					case '@':
						if buffer[position] != rune('@') {
							goto l205
						}
						position++
						break
					case '~':
						if buffer[position] != rune('~') {
							goto l205
						}
						position++
						break
					case ';':
						if buffer[position] != rune(';') {
							goto l205
						}
						position++
						break
					case '+':
						if buffer[position] != rune('+') {
							goto l205
						}
						position++
						break
					case '/':
						if buffer[position] != rune('/') {
							goto l205
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l205
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l205
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l205
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l205
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l205
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l205
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l205
						}
						position++
						break
					}
				}

			l207:
				{
					position208, tokenIndex208 := position, tokenIndex
					{
						switch buffer[position] {
						case '*':
							if buffer[position] != rune('*') {
								goto l208
							}
							position++
							break
						case '>':
							if buffer[position] != rune('>') {
								goto l208
							}
							position++
							break
						case '<':
							if buffer[position] != rune('<') {
								goto l208
							}
							position++
							break
						case '@':
							if buffer[position] != rune('@') {
								goto l208
							}
							position++
							break
						case '~':
							if buffer[position] != rune('~') {
								goto l208
							}
							position++
							break
						case ';':
							if buffer[position] != rune(';') {
								goto l208
							}
							position++
							break
						case '+':
							if buffer[position] != rune('+') {
								goto l208
							}
							position++
							break
						case '/':
							if buffer[position] != rune('/') {
								goto l208
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
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
						case '.':
							if buffer[position] != rune('.') {
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

					goto l207
				l208:
					position, tokenIndex = position208, tokenIndex208
				}
				add(ruleOtherParamValue, position206)
			}
			return true
		l205:
			position, tokenIndex = position205, tokenIndex205
			return false
		},
		/* 16 DoubleQuotedValue <- <(!'"' .)*> */
		func() bool {
			{
				position212 := position
			l213:
				{
					position214, tokenIndex214 := position, tokenIndex
					{
						position215, tokenIndex215 := position, tokenIndex
						if buffer[position] != rune('"') {
							goto l215
						}
						position++
						goto l214
					l215:
						position, tokenIndex = position215, tokenIndex215
					}
					if !matchDot() {
						goto l214
					}
					goto l213
				l214:
					position, tokenIndex = position214, tokenIndex214
				}
				add(ruleDoubleQuotedValue, position212)
			}
			return true
		},
		/* 17 SingleQuotedValue <- <(!'\'' .)*> */
		func() bool {
			{
				position217 := position
			l218:
				{
					position219, tokenIndex219 := position, tokenIndex
					{
						position220, tokenIndex220 := position, tokenIndex
						if buffer[position] != rune('\'') {
							goto l220
						}
						position++
						goto l219
					l220:
						position, tokenIndex = position220, tokenIndex220
					}
					if !matchDot() {
						goto l219
					}
					goto l218
				l219:
					position, tokenIndex = position219, tokenIndex219
				}
				add(ruleSingleQuotedValue, position217)
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
		func() bool {
			position226, tokenIndex226 := position, tokenIndex
			{
				position227 := position
				if buffer[position] != rune('{') {
					goto l226
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l226
				}
				{
					position228 := position
					if !_rules[ruleIdentifier]() {
						goto l226
					}
					add(rulePegText, position228)
				}
				if !_rules[ruleWhiteSpacing]() {
					goto l226
				}
				if buffer[position] != rune('}') {
					goto l226
				}
				position++
				add(ruleHoleValue, position227)
			}
			return true
		l226:
			position, tokenIndex = position226, tokenIndex226
			return false
		},
		/* 24 HolesStringValue <- <<(OtherParamValue? HoleValue OtherParamValue?)+>> */
		nil,
		/* 25 HoleWithSuffixValue <- <<(HoleValue OtherParamValue+ (OtherParamValue? HoleValue OtherParamValue?)*)>> */
		nil,
		/* 26 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)*))> */
		nil,
		/* 27 SingleQuote <- <'\''> */
		func() bool {
			position232, tokenIndex232 := position, tokenIndex
			{
				position233 := position
				if buffer[position] != rune('\'') {
					goto l232
				}
				position++
				add(ruleSingleQuote, position233)
			}
			return true
		l232:
			position, tokenIndex = position232, tokenIndex232
			return false
		},
		/* 28 DoubleQuote <- <'"'> */
		func() bool {
			position234, tokenIndex234 := position, tokenIndex
			{
				position235 := position
				if buffer[position] != rune('"') {
					goto l234
				}
				position++
				add(ruleDoubleQuote, position235)
			}
			return true
		l234:
			position, tokenIndex = position234, tokenIndex234
			return false
		},
		/* 29 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position237 := position
			l238:
				{
					position239, tokenIndex239 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l239
					}
					goto l238
				l239:
					position, tokenIndex = position239, tokenIndex239
				}
				add(ruleWhiteSpacing, position237)
			}
			return true
		},
		/* 30 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position240, tokenIndex240 := position, tokenIndex
			{
				position241 := position
				if !_rules[ruleWhitespace]() {
					goto l240
				}
			l242:
				{
					position243, tokenIndex243 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l243
					}
					goto l242
				l243:
					position, tokenIndex = position243, tokenIndex243
				}
				add(ruleMustWhiteSpacing, position241)
			}
			return true
		l240:
			position, tokenIndex = position240, tokenIndex240
			return false
		},
		/* 31 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position244, tokenIndex244 := position, tokenIndex
			{
				position245 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l244
				}
				if buffer[position] != rune('=') {
					goto l244
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l244
				}
				add(ruleEqual, position245)
			}
			return true
		l244:
			position, tokenIndex = position244, tokenIndex244
			return false
		},
		/* 32 BlankLine <- <(WhiteSpacing EndOfLine)> */
		func() bool {
			position246, tokenIndex246 := position, tokenIndex
			{
				position247 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l246
				}
				if !_rules[ruleEndOfLine]() {
					goto l246
				}
				add(ruleBlankLine, position247)
			}
			return true
		l246:
			position, tokenIndex = position246, tokenIndex246
			return false
		},
		/* 33 Whitespace <- <(' ' / '\t')> */
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
		/* 34 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
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
		/* 35 EndOfFile <- <!.> */
		nil,
		/* 37 Action0 <- <{ p.NewStatement() }> */
		nil,
		/* 38 Action1 <- <{ p.StatementDone() }> */
		nil,
		nil,
		/* 40 Action2 <- <{ p.addDeclarationIdentifier(text) }> */
		nil,
		/* 41 Action3 <- <{ p.addValue() }> */
		nil,
		/* 42 Action4 <- <{ p.addAction(text) }> */
		nil,
		/* 43 Action5 <- <{ p.addEntity(text) }> */
		nil,
		/* 44 Action6 <- <{ p.addParamKey(text) }> */
		nil,
		/* 45 Action7 <- <{  p.addFirstValueInList() }> */
		nil,
		/* 46 Action8 <- <{  p.lastValueInList() }> */
		nil,
		/* 47 Action9 <- <{  p.addHolesStringParam(text) }> */
		nil,
		/* 48 Action10 <- <{  p.addParamHoleValue(text) }> */
		nil,
		/* 49 Action11 <- <{  p.addHolesStringParam(text) }> */
		nil,
		/* 50 Action12 <- <{  p.addAliasParam(text) }> */
		nil,
		/* 51 Action13 <- <{ p.addStringValue(text) }> */
		nil,
		/* 52 Action14 <- <{ p.addStringValue(text) }> */
		nil,
		/* 53 Action15 <- <{ p.addParamValue(text) }> */
		nil,
		/* 54 Action16 <- <{  p.addParamRefValue(text) }> */
		nil,
		/* 55 Action17 <- <{ p.addParamCidrValue(text) }> */
		nil,
		/* 56 Action18 <- <{ p.addParamIpValue(text) }> */
		nil,
		/* 57 Action19 <- <{ p.addParamValue(text) }> */
		nil,
	}
	p.rules = _rules
}
