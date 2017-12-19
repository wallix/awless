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
	ruleUnquotedParamValue
	ruleUnquotedParam
	ruleConcatenationValue
	ruleQuotedStringValue
	ruleQuotedString
	ruleDoubleQuotedValue
	ruleSingleQuotedValue
	ruleIntRangeValue
	ruleRefValue
	ruleAliasValue
	ruleHoleValue
	ruleHole
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
	ruleAction22
	ruleAction23
	ruleAction24
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
	"UnquotedParamValue",
	"UnquotedParam",
	"ConcatenationValue",
	"QuotedStringValue",
	"QuotedString",
	"DoubleQuotedValue",
	"SingleQuotedValue",
	"IntRangeValue",
	"RefValue",
	"AliasValue",
	"HoleValue",
	"Hole",
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
	"Action22",
	"Action23",
	"Action24",
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
	rules  [67]func() bool
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
			p.addAliasParam(text)
		case ruleAction12:
			p.addParamRefValue(text)
		case ruleAction13:
			p.addParamValue(text)
		case ruleAction14:
			p.addParamValue(text)
		case ruleAction15:
			p.addFirstValueInConcatenation()
		case ruleAction16:
			p.lastValueInConcatenation()
		case ruleAction17:
			p.addFirstValueInConcatenation()
		case ruleAction18:
			p.lastValueInConcatenation()
		case ruleAction19:
			p.addStringValue(text)
		case ruleAction20:
			p.addParamHoleValue(text)
		case ruleAction21:
			p.addFirstValueInConcatenation()
		case ruleAction22:
			p.lastValueInConcatenation()
		case ruleAction23:
			p.addFirstValueInConcatenation()
		case ruleAction24:
			p.lastValueInConcatenation()

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
		/* 13 NoRefValue <- <(ConcatenationValue / HoleWithSuffixValue / HoleValue / HolesStringValue / (AliasValue Action11) / (DoubleQuote CustomTypedValue DoubleQuote) / (SingleQuote CustomTypedValue SingleQuote) / CustomTypedValue / QuotedStringValue / UnquotedParamValue)> */
		nil,
		/* 14 Value <- <((RefValue Action12) / NoRefValue)> */
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
						add(ruleAction12, position)
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
									position133, tokenIndex133 := position, tokenIndex
									{
										add(ruleAction15, position)
									}
									if !_rules[ruleHoleValue]() {
										goto l134
									}
									if !_rules[ruleWhiteSpacing]() {
										goto l134
									}
									if buffer[position] != rune('+') {
										goto l134
									}
									position++
									if !_rules[ruleWhiteSpacing]() {
										goto l134
									}
									{
										position138, tokenIndex138 := position, tokenIndex
										if !_rules[ruleQuotedStringValue]() {
											goto l139
										}
										goto l138
									l139:
										position, tokenIndex = position138, tokenIndex138
										if !_rules[ruleHoleValue]() {
											goto l134
										}
									}
								l138:
								l136:
									{
										position137, tokenIndex137 := position, tokenIndex
										if !_rules[ruleWhiteSpacing]() {
											goto l137
										}
										if buffer[position] != rune('+') {
											goto l137
										}
										position++
										if !_rules[ruleWhiteSpacing]() {
											goto l137
										}
										{
											position140, tokenIndex140 := position, tokenIndex
											if !_rules[ruleQuotedStringValue]() {
												goto l141
											}
											goto l140
										l141:
											position, tokenIndex = position140, tokenIndex140
											if !_rules[ruleHoleValue]() {
												goto l137
											}
										}
									l140:
										goto l136
									l137:
										position, tokenIndex = position137, tokenIndex137
									}
									{
										add(ruleAction16, position)
									}
									goto l133
								l134:
									position, tokenIndex = position133, tokenIndex133
									{
										add(ruleAction17, position)
									}
									if !_rules[ruleQuotedStringValue]() {
										goto l131
									}
									if !_rules[ruleWhiteSpacing]() {
										goto l131
									}
									if buffer[position] != rune('+') {
										goto l131
									}
									position++
									if !_rules[ruleWhiteSpacing]() {
										goto l131
									}
									{
										position146, tokenIndex146 := position, tokenIndex
										if !_rules[ruleQuotedStringValue]() {
											goto l147
										}
										goto l146
									l147:
										position, tokenIndex = position146, tokenIndex146
										if !_rules[ruleHoleValue]() {
											goto l131
										}
									}
								l146:
								l144:
									{
										position145, tokenIndex145 := position, tokenIndex
										if !_rules[ruleWhiteSpacing]() {
											goto l145
										}
										if buffer[position] != rune('+') {
											goto l145
										}
										position++
										if !_rules[ruleWhiteSpacing]() {
											goto l145
										}
										{
											position148, tokenIndex148 := position, tokenIndex
											if !_rules[ruleQuotedStringValue]() {
												goto l149
											}
											goto l148
										l149:
											position, tokenIndex = position148, tokenIndex148
											if !_rules[ruleHoleValue]() {
												goto l145
											}
										}
									l148:
										goto l144
									l145:
										position, tokenIndex = position145, tokenIndex145
									}
									{
										add(ruleAction18, position)
									}
								}
							l133:
								add(ruleConcatenationValue, position132)
							}
							goto l130
						l131:
							position, tokenIndex = position130, tokenIndex130
							{
								position152 := position
								{
									add(ruleAction23, position)
								}
								{
									position154 := position
									if !_rules[ruleHoleValue]() {
										goto l151
									}
									if !_rules[ruleUnquotedParamValue]() {
										goto l151
									}
								l155:
									{
										position156, tokenIndex156 := position, tokenIndex
										if !_rules[ruleUnquotedParamValue]() {
											goto l156
										}
										goto l155
									l156:
										position, tokenIndex = position156, tokenIndex156
									}
								l157:
									{
										position158, tokenIndex158 := position, tokenIndex
										{
											position159, tokenIndex159 := position, tokenIndex
											if !_rules[ruleUnquotedParamValue]() {
												goto l159
											}
											goto l160
										l159:
											position, tokenIndex = position159, tokenIndex159
										}
									l160:
										if !_rules[ruleHoleValue]() {
											goto l158
										}
										{
											position161, tokenIndex161 := position, tokenIndex
											if !_rules[ruleUnquotedParamValue]() {
												goto l161
											}
											goto l162
										l161:
											position, tokenIndex = position161, tokenIndex161
										}
									l162:
										goto l157
									l158:
										position, tokenIndex = position158, tokenIndex158
									}
									add(rulePegText, position154)
								}
								{
									add(ruleAction24, position)
								}
								add(ruleHoleWithSuffixValue, position152)
							}
							goto l130
						l151:
							position, tokenIndex = position130, tokenIndex130
							if !_rules[ruleHoleValue]() {
								goto l164
							}
							goto l130
						l164:
							position, tokenIndex = position130, tokenIndex130
							{
								position166 := position
								{
									add(ruleAction21, position)
								}
								{
									position168 := position
									{
										position171, tokenIndex171 := position, tokenIndex
										if !_rules[ruleUnquotedParamValue]() {
											goto l171
										}
										goto l172
									l171:
										position, tokenIndex = position171, tokenIndex171
									}
								l172:
									if !_rules[ruleHoleValue]() {
										goto l165
									}
									{
										position173, tokenIndex173 := position, tokenIndex
										if !_rules[ruleUnquotedParamValue]() {
											goto l173
										}
										goto l174
									l173:
										position, tokenIndex = position173, tokenIndex173
									}
								l174:
								l169:
									{
										position170, tokenIndex170 := position, tokenIndex
										{
											position175, tokenIndex175 := position, tokenIndex
											if !_rules[ruleUnquotedParamValue]() {
												goto l175
											}
											goto l176
										l175:
											position, tokenIndex = position175, tokenIndex175
										}
									l176:
										if !_rules[ruleHoleValue]() {
											goto l170
										}
										{
											position177, tokenIndex177 := position, tokenIndex
											if !_rules[ruleUnquotedParamValue]() {
												goto l177
											}
											goto l178
										l177:
											position, tokenIndex = position177, tokenIndex177
										}
									l178:
										goto l169
									l170:
										position, tokenIndex = position170, tokenIndex170
									}
									add(rulePegText, position168)
								}
								{
									add(ruleAction22, position)
								}
								add(ruleHolesStringValue, position166)
							}
							goto l130
						l165:
							position, tokenIndex = position130, tokenIndex130
							{
								position181 := position
								{
									position182, tokenIndex182 := position, tokenIndex
									if buffer[position] != rune('@') {
										goto l183
									}
									position++
									{
										position184 := position
										if !_rules[ruleUnquotedParam]() {
											goto l183
										}
										add(rulePegText, position184)
									}
									goto l182
								l183:
									position, tokenIndex = position182, tokenIndex182
									if buffer[position] != rune('@') {
										goto l185
									}
									position++
									if !_rules[ruleDoubleQuotedValue]() {
										goto l185
									}
									goto l182
								l185:
									position, tokenIndex = position182, tokenIndex182
									if buffer[position] != rune('@') {
										goto l180
									}
									position++
									if !_rules[ruleSingleQuotedValue]() {
										goto l180
									}
								}
							l182:
								add(ruleAliasValue, position181)
							}
							{
								add(ruleAction11, position)
							}
							goto l130
						l180:
							position, tokenIndex = position130, tokenIndex130
							if !_rules[ruleDoubleQuote]() {
								goto l187
							}
							if !_rules[ruleCustomTypedValue]() {
								goto l187
							}
							if !_rules[ruleDoubleQuote]() {
								goto l187
							}
							goto l130
						l187:
							position, tokenIndex = position130, tokenIndex130
							if !_rules[ruleSingleQuote]() {
								goto l188
							}
							if !_rules[ruleCustomTypedValue]() {
								goto l188
							}
							if !_rules[ruleSingleQuote]() {
								goto l188
							}
							goto l130
						l188:
							position, tokenIndex = position130, tokenIndex130
							if !_rules[ruleCustomTypedValue]() {
								goto l189
							}
							goto l130
						l189:
							position, tokenIndex = position130, tokenIndex130
							if !_rules[ruleQuotedStringValue]() {
								goto l190
							}
							goto l130
						l190:
							position, tokenIndex = position130, tokenIndex130
							if !_rules[ruleUnquotedParamValue]() {
								goto l122
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
		/* 15 CustomTypedValue <- <(<IntRangeValue> Action13)> */
		func() bool {
			position191, tokenIndex191 := position, tokenIndex
			{
				position192 := position
				{
					position193 := position
					{
						position194 := position
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l191
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
						if buffer[position] != rune('-') {
							goto l191
						}
						position++
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l191
						}
						position++
					l197:
						{
							position198, tokenIndex198 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l198
							}
							position++
							goto l197
						l198:
							position, tokenIndex = position198, tokenIndex198
						}
						add(ruleIntRangeValue, position194)
					}
					add(rulePegText, position193)
				}
				{
					add(ruleAction13, position)
				}
				add(ruleCustomTypedValue, position192)
			}
			return true
		l191:
			position, tokenIndex = position191, tokenIndex191
			return false
		},
		/* 16 UnquotedParamValue <- <(<UnquotedParam> Action14)> */
		func() bool {
			position200, tokenIndex200 := position, tokenIndex
			{
				position201 := position
				{
					position202 := position
					if !_rules[ruleUnquotedParam]() {
						goto l200
					}
					add(rulePegText, position202)
				}
				{
					add(ruleAction14, position)
				}
				add(ruleUnquotedParamValue, position201)
			}
			return true
		l200:
			position, tokenIndex = position200, tokenIndex200
			return false
		},
		/* 17 UnquotedParam <- <((&('*') '*') | (&('>') '>') | (&('<') '<') | (&('@') '@') | (&('~') '~') | (&(';') ';') | (&('+') '+') | (&('/') '/') | (&(':') ':') | (&('_') '_') | (&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position204, tokenIndex204 := position, tokenIndex
			{
				position205 := position
				{
					switch buffer[position] {
					case '*':
						if buffer[position] != rune('*') {
							goto l204
						}
						position++
						break
					case '>':
						if buffer[position] != rune('>') {
							goto l204
						}
						position++
						break
					case '<':
						if buffer[position] != rune('<') {
							goto l204
						}
						position++
						break
					case '@':
						if buffer[position] != rune('@') {
							goto l204
						}
						position++
						break
					case '~':
						if buffer[position] != rune('~') {
							goto l204
						}
						position++
						break
					case ';':
						if buffer[position] != rune(';') {
							goto l204
						}
						position++
						break
					case '+':
						if buffer[position] != rune('+') {
							goto l204
						}
						position++
						break
					case '/':
						if buffer[position] != rune('/') {
							goto l204
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
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
					case '.':
						if buffer[position] != rune('.') {
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
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
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

			l206:
				{
					position207, tokenIndex207 := position, tokenIndex
					{
						switch buffer[position] {
						case '*':
							if buffer[position] != rune('*') {
								goto l207
							}
							position++
							break
						case '>':
							if buffer[position] != rune('>') {
								goto l207
							}
							position++
							break
						case '<':
							if buffer[position] != rune('<') {
								goto l207
							}
							position++
							break
						case '@':
							if buffer[position] != rune('@') {
								goto l207
							}
							position++
							break
						case '~':
							if buffer[position] != rune('~') {
								goto l207
							}
							position++
							break
						case ';':
							if buffer[position] != rune(';') {
								goto l207
							}
							position++
							break
						case '+':
							if buffer[position] != rune('+') {
								goto l207
							}
							position++
							break
						case '/':
							if buffer[position] != rune('/') {
								goto l207
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l207
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l207
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l207
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l207
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l207
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l207
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l207
							}
							position++
							break
						}
					}

					goto l206
				l207:
					position, tokenIndex = position207, tokenIndex207
				}
				add(ruleUnquotedParam, position205)
			}
			return true
		l204:
			position, tokenIndex = position204, tokenIndex204
			return false
		},
		/* 18 ConcatenationValue <- <((Action15 HoleValue (WhiteSpacing '+' WhiteSpacing (QuotedStringValue / HoleValue))+ Action16) / (Action17 QuotedStringValue (WhiteSpacing '+' WhiteSpacing (QuotedStringValue / HoleValue))+ Action18))> */
		nil,
		/* 19 QuotedStringValue <- <(QuotedString Action19)> */
		func() bool {
			position211, tokenIndex211 := position, tokenIndex
			{
				position212 := position
				{
					position213 := position
					{
						position214, tokenIndex214 := position, tokenIndex
						if !_rules[ruleDoubleQuotedValue]() {
							goto l215
						}
						goto l214
					l215:
						position, tokenIndex = position214, tokenIndex214
						if !_rules[ruleSingleQuotedValue]() {
							goto l211
						}
					}
				l214:
					add(ruleQuotedString, position213)
				}
				{
					add(ruleAction19, position)
				}
				add(ruleQuotedStringValue, position212)
			}
			return true
		l211:
			position, tokenIndex = position211, tokenIndex211
			return false
		},
		/* 20 QuotedString <- <(DoubleQuotedValue / SingleQuotedValue)> */
		nil,
		/* 21 DoubleQuotedValue <- <(DoubleQuote <(!'"' .)*> DoubleQuote)> */
		func() bool {
			position218, tokenIndex218 := position, tokenIndex
			{
				position219 := position
				if !_rules[ruleDoubleQuote]() {
					goto l218
				}
				{
					position220 := position
				l221:
					{
						position222, tokenIndex222 := position, tokenIndex
						{
							position223, tokenIndex223 := position, tokenIndex
							if buffer[position] != rune('"') {
								goto l223
							}
							position++
							goto l222
						l223:
							position, tokenIndex = position223, tokenIndex223
						}
						if !matchDot() {
							goto l222
						}
						goto l221
					l222:
						position, tokenIndex = position222, tokenIndex222
					}
					add(rulePegText, position220)
				}
				if !_rules[ruleDoubleQuote]() {
					goto l218
				}
				add(ruleDoubleQuotedValue, position219)
			}
			return true
		l218:
			position, tokenIndex = position218, tokenIndex218
			return false
		},
		/* 22 SingleQuotedValue <- <(SingleQuote <(!'\'' .)*> SingleQuote)> */
		func() bool {
			position224, tokenIndex224 := position, tokenIndex
			{
				position225 := position
				if !_rules[ruleSingleQuote]() {
					goto l224
				}
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
					add(rulePegText, position226)
				}
				if !_rules[ruleSingleQuote]() {
					goto l224
				}
				add(ruleSingleQuotedValue, position225)
			}
			return true
		l224:
			position, tokenIndex = position224, tokenIndex224
			return false
		},
		/* 23 IntRangeValue <- <([0-9]+ '-' [0-9]+)> */
		nil,
		/* 24 RefValue <- <('$' <Identifier>)> */
		nil,
		/* 25 AliasValue <- <(('@' <UnquotedParam>) / ('@' DoubleQuotedValue) / ('@' SingleQuotedValue))> */
		nil,
		/* 26 HoleValue <- <(Hole Action20)> */
		func() bool {
			position233, tokenIndex233 := position, tokenIndex
			{
				position234 := position
				{
					position235 := position
					if buffer[position] != rune('{') {
						goto l233
					}
					position++
					if !_rules[ruleWhiteSpacing]() {
						goto l233
					}
					{
						position236 := position
						if !_rules[ruleIdentifier]() {
							goto l233
						}
						add(rulePegText, position236)
					}
					if !_rules[ruleWhiteSpacing]() {
						goto l233
					}
					if buffer[position] != rune('}') {
						goto l233
					}
					position++
					add(ruleHole, position235)
				}
				{
					add(ruleAction20, position)
				}
				add(ruleHoleValue, position234)
			}
			return true
		l233:
			position, tokenIndex = position233, tokenIndex233
			return false
		},
		/* 27 Hole <- <('{' WhiteSpacing <Identifier> WhiteSpacing '}')> */
		nil,
		/* 28 HolesStringValue <- <(Action21 <(UnquotedParamValue? HoleValue UnquotedParamValue?)+> Action22)> */
		nil,
		/* 29 HoleWithSuffixValue <- <(Action23 <(HoleValue UnquotedParamValue+ (UnquotedParamValue? HoleValue UnquotedParamValue?)*)> Action24)> */
		nil,
		/* 30 Comment <- <(('#' (!EndOfLine .)*) / ('/' '/' (!EndOfLine .)*))> */
		nil,
		/* 31 SingleQuote <- <'\''> */
		func() bool {
			position242, tokenIndex242 := position, tokenIndex
			{
				position243 := position
				if buffer[position] != rune('\'') {
					goto l242
				}
				position++
				add(ruleSingleQuote, position243)
			}
			return true
		l242:
			position, tokenIndex = position242, tokenIndex242
			return false
		},
		/* 32 DoubleQuote <- <'"'> */
		func() bool {
			position244, tokenIndex244 := position, tokenIndex
			{
				position245 := position
				if buffer[position] != rune('"') {
					goto l244
				}
				position++
				add(ruleDoubleQuote, position245)
			}
			return true
		l244:
			position, tokenIndex = position244, tokenIndex244
			return false
		},
		/* 33 WhiteSpacing <- <Whitespace*> */
		func() bool {
			{
				position247 := position
			l248:
				{
					position249, tokenIndex249 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l249
					}
					goto l248
				l249:
					position, tokenIndex = position249, tokenIndex249
				}
				add(ruleWhiteSpacing, position247)
			}
			return true
		},
		/* 34 MustWhiteSpacing <- <Whitespace+> */
		func() bool {
			position250, tokenIndex250 := position, tokenIndex
			{
				position251 := position
				if !_rules[ruleWhitespace]() {
					goto l250
				}
			l252:
				{
					position253, tokenIndex253 := position, tokenIndex
					if !_rules[ruleWhitespace]() {
						goto l253
					}
					goto l252
				l253:
					position, tokenIndex = position253, tokenIndex253
				}
				add(ruleMustWhiteSpacing, position251)
			}
			return true
		l250:
			position, tokenIndex = position250, tokenIndex250
			return false
		},
		/* 35 Equal <- <(WhiteSpacing '=' WhiteSpacing)> */
		func() bool {
			position254, tokenIndex254 := position, tokenIndex
			{
				position255 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l254
				}
				if buffer[position] != rune('=') {
					goto l254
				}
				position++
				if !_rules[ruleWhiteSpacing]() {
					goto l254
				}
				add(ruleEqual, position255)
			}
			return true
		l254:
			position, tokenIndex = position254, tokenIndex254
			return false
		},
		/* 36 BlankLine <- <(WhiteSpacing EndOfLine)> */
		func() bool {
			position256, tokenIndex256 := position, tokenIndex
			{
				position257 := position
				if !_rules[ruleWhiteSpacing]() {
					goto l256
				}
				if !_rules[ruleEndOfLine]() {
					goto l256
				}
				add(ruleBlankLine, position257)
			}
			return true
		l256:
			position, tokenIndex = position256, tokenIndex256
			return false
		},
		/* 37 Whitespace <- <(' ' / '\t')> */
		func() bool {
			position258, tokenIndex258 := position, tokenIndex
			{
				position259 := position
				{
					position260, tokenIndex260 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l261
					}
					position++
					goto l260
				l261:
					position, tokenIndex = position260, tokenIndex260
					if buffer[position] != rune('\t') {
						goto l258
					}
					position++
				}
			l260:
				add(ruleWhitespace, position259)
			}
			return true
		l258:
			position, tokenIndex = position258, tokenIndex258
			return false
		},
		/* 38 EndOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position262, tokenIndex262 := position, tokenIndex
			{
				position263 := position
				{
					position264, tokenIndex264 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l265
					}
					position++
					if buffer[position] != rune('\n') {
						goto l265
					}
					position++
					goto l264
				l265:
					position, tokenIndex = position264, tokenIndex264
					if buffer[position] != rune('\n') {
						goto l266
					}
					position++
					goto l264
				l266:
					position, tokenIndex = position264, tokenIndex264
					if buffer[position] != rune('\r') {
						goto l262
					}
					position++
				}
			l264:
				add(ruleEndOfLine, position263)
			}
			return true
		l262:
			position, tokenIndex = position262, tokenIndex262
			return false
		},
		/* 39 EndOfFile <- <!.> */
		nil,
		/* 41 Action0 <- <{ p.NewStatement() }> */
		nil,
		/* 42 Action1 <- <{ p.StatementDone() }> */
		nil,
		nil,
		/* 44 Action2 <- <{ p.addDeclarationIdentifier(text) }> */
		nil,
		/* 45 Action3 <- <{ p.addValue() }> */
		nil,
		/* 46 Action4 <- <{ p.addAction(text) }> */
		nil,
		/* 47 Action5 <- <{ p.addEntity(text) }> */
		nil,
		/* 48 Action6 <- <{ p.addParamKey(text) }> */
		nil,
		/* 49 Action7 <- <{  p.addFirstValueInList() }> */
		nil,
		/* 50 Action8 <- <{  p.lastValueInList() }> */
		nil,
		/* 51 Action9 <- <{  p.addFirstValueInList() }> */
		nil,
		/* 52 Action10 <- <{  p.lastValueInList() }> */
		nil,
		/* 53 Action11 <- <{  p.addAliasParam(text) }> */
		nil,
		/* 54 Action12 <- <{  p.addParamRefValue(text) }> */
		nil,
		/* 55 Action13 <- <{ p.addParamValue(text) }> */
		nil,
		/* 56 Action14 <- <{ p.addParamValue(text) }> */
		nil,
		/* 57 Action15 <- <{ p.addFirstValueInConcatenation() }> */
		nil,
		/* 58 Action16 <- <{  p.lastValueInConcatenation() }> */
		nil,
		/* 59 Action17 <- <{ p.addFirstValueInConcatenation() }> */
		nil,
		/* 60 Action18 <- <{  p.lastValueInConcatenation() }> */
		nil,
		/* 61 Action19 <- <{ p.addStringValue(text) }> */
		nil,
		/* 62 Action20 <- <{  p.addParamHoleValue(text) }> */
		nil,
		/* 63 Action21 <- <{ p.addFirstValueInConcatenation() }> */
		nil,
		/* 64 Action22 <- <{  p.lastValueInConcatenation() }> */
		nil,
		/* 65 Action23 <- <{ p.addFirstValueInConcatenation() }> */
		nil,
		/* 66 Action24 <- <{  p.lastValueInConcatenation() }> */
		nil,
	}
	p.rules = _rules
}
