//line grammer.y:6
package parse

import __yyfmt__ "fmt"

//line grammer.y:6
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

var DualRunes = map[string]int{
	".": DOT,
	",": COMMA,

	"--": DOUBLEDASH,
	"-":  DASH,
	"=":  EQUAL,

	">>":  DOUBLEANGR,
	">":   ANGR,
	"/>":  SLASHANGR,
	"\\>": BACKSLASHANGR,
}

//line grammer.y:34
type yySymType struct {
	yys          int
	nodeList     *NodeList
	node         Node
	arrow        ArrowType
	arrowStem    ArrowStemType
	arrowHead    ArrowHeadType
	actorRef     ActorRef
	noteAlign    NoteAlignment
	dividerType  GapType
	blockSegList *BlockSegmentList
	attrList     *AttributeList
	attr         *Attribute

	sval string
}

const K_TITLE = 57346
const K_PARTICIPANT = 57347
const K_NOTE = 57348
const K_STYLE = 57349
const K_LEFT = 57350
const K_RIGHT = 57351
const K_OVER = 57352
const K_OF = 57353
const K_HORIZONTAL = 57354
const K_SPACER = 57355
const K_GAP = 57356
const K_LINE = 57357
const K_FRAME = 57358
const K_ALT = 57359
const K_ELSEALT = 57360
const K_ELSE = 57361
const K_END = 57362
const K_LOOP = 57363
const K_OPT = 57364
const DASH = 57365
const DOUBLEDASH = 57366
const DOT = 57367
const EQUAL = 57368
const COMMA = 57369
const ANGR = 57370
const DOUBLEANGR = 57371
const BACKSLASHANGR = 57372
const SLASHANGR = 57373
const PARL = 57374
const PARR = 57375
const MESSAGE = 57376
const IDENT = 57377

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"K_TITLE",
	"K_PARTICIPANT",
	"K_NOTE",
	"K_STYLE",
	"K_LEFT",
	"K_RIGHT",
	"K_OVER",
	"K_OF",
	"K_HORIZONTAL",
	"K_SPACER",
	"K_GAP",
	"K_LINE",
	"K_FRAME",
	"K_ALT",
	"K_ELSEALT",
	"K_ELSE",
	"K_END",
	"K_LOOP",
	"K_OPT",
	"DASH",
	"DOUBLEDASH",
	"DOT",
	"EQUAL",
	"COMMA",
	"ANGR",
	"DOUBLEANGR",
	"BACKSLASHANGR",
	"SLASHANGR",
	"PARL",
	"PARR",
	"MESSAGE",
	"IDENT",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line grammer.y:290

// Manages the lexer as well as the current diagram being parsed
type parseState struct {
	S     scanner.Scanner
	err   error
	atEof bool
	//diagram     *Diagram
	procInstrs []string
	nodeList   *NodeList
}

func newParseState(src io.Reader, filename string) *parseState {
	ps := &parseState{}
	ps.S.Init(src)
	ps.S.Position.Filename = filename
	//    ps.diagram = &Diagram{}

	return ps
}

func (ps *parseState) Lex(lval *yySymType) int {
	if ps.atEof {
		return 0
	}
	for {
		tok := ps.S.Scan()
		switch tok {
		case scanner.EOF:
			ps.atEof = true
			return 0
		case '#':
			ps.scanComment()
		case ':':
			return ps.scanMessage(lval)
		case '(':
			return PARL
		case ')':
			return PARR
		case '-', '>', '*', '=', '/', '\\', '.', ',':
			if res, isTok := ps.handleDoubleRune(tok); isTok {
				return res
			} else {
				ps.Error("Invalid token: " + scanner.TokenString(tok))
			}
		case scanner.Ident:
			return ps.scanKeywordOrIdent(lval)
		default:
			ps.Error("Invalid token: " + scanner.TokenString(tok))
		}
	}
}

func (ps *parseState) handleDoubleRune(firstRune rune) (int, bool) {
	nextRune := ps.S.Peek()

	// Try the double rune
	if nextRune != scanner.EOF {
		tokStr := string(firstRune) + string(nextRune)
		if tok, hasTok := DualRunes[tokStr]; hasTok {
			ps.NextRune()
			return tok, true
		}
	}

	// Try the single rune
	tokStr := string(firstRune)
	if tok, hasTok := DualRunes[tokStr]; hasTok {
		return tok, true
	}

	return 0, false
}

func (ps *parseState) scanKeywordOrIdent(lval *yySymType) int {
	tokVal := ps.S.TokenText()
	switch strings.ToLower(tokVal) {
	case "title":
		return K_TITLE
	case "participant":
		return K_PARTICIPANT
	case "note":
		return K_NOTE
	case "left":
		return K_LEFT
	case "right":
		return K_RIGHT
	case "over":
		return K_OVER
	case "of":
		return K_OF
	case "spacer":
		return K_SPACER
	case "gap":
		return K_GAP
	case "frame":
		return K_FRAME
	case "line":
		return K_LINE
	case "style":
		return K_STYLE
	case "horizontal":
		return K_HORIZONTAL
	case "alt":
		return K_ALT
	case "elsealt":
		return K_ELSEALT
	case "else":
		return K_ELSE
	case "end":
		return K_END
	case "loop":
		return K_LOOP
	case "opt":
		return K_OPT
	default:
		lval.sval = tokVal
		return IDENT
	}
}

// Scans a message.  A message is all characters up to the new line
func (ps *parseState) scanMessage(lval *yySymType) int {
	buf := new(bytes.Buffer)
	r := ps.NextRune()
	for (r != '\n') && (r != scanner.EOF) {
		if r == '\\' {
			nr := ps.NextRune()
			switch nr {
			case 'n':
				buf.WriteRune('\n')
			case '\\':
				buf.WriteRune('\\')
			default:
				ps.Error("Invalid backslash escape: \\" + string(nr))
			}
		} else {
			buf.WriteRune(r)
		}
		r = ps.NextRune()
	}

	lval.sval = strings.TrimSpace(buf.String())
	return MESSAGE
}

// Scans a comment.  This ignores all characters up to the new line.
func (ps *parseState) scanComment() {
	var buf *bytes.Buffer

	r := ps.NextRune()
	if r == '!' {
		// This starts a processor instruction
		buf = new(bytes.Buffer)
		r = ps.NextRune()
	}

	for (r != '\n') && (r != scanner.EOF) {
		if buf != nil {
			buf.WriteRune(r)
		}
		r = ps.NextRune()
	}

	if buf != nil {
		ps.procInstrs = append(ps.procInstrs, strings.TrimSpace(buf.String()))
	}
}

func (ps *parseState) NextRune() rune {
	if ps.atEof {
		return scanner.EOF
	}

	r := ps.S.Next()
	if r == scanner.EOF {
		ps.atEof = true
	}

	return r
}

func (ps *parseState) Error(err string) {
	errMsg := fmt.Sprintf("%s:%d: %s", ps.S.Position.Filename, ps.S.Position.Line, err)
	ps.err = errors.New(errMsg)
}

func Parse(reader io.Reader, filename string) (*NodeList, error) {
	ps := newParseState(reader, filename)
	yyParse(ps)

	// Add processing instructions to the start of the node list
	for i := len(ps.procInstrs) - 1; i >= 0; i-- {
		instrParts := strings.SplitN(ps.procInstrs[i], " ", 2)
		name, value := strings.TrimSpace(instrParts[0]), strings.TrimSpace(instrParts[1])
		ps.nodeList = &NodeList{&ProcessInstructionNode{name, value}, ps.nodeList}
	}

	if ps.err != nil {
		return nil, ps.err
	} else {
		return ps.nodeList, nil
	}
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 55
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 92

var yyAct = [...]int{

	72, 2, 65, 16, 28, 25, 13, 15, 17, 14,
	23, 24, 23, 24, 18, 85, 67, 30, 71, 19,
	86, 83, 82, 21, 20, 70, 69, 68, 61, 54,
	55, 56, 57, 77, 29, 52, 48, 22, 49, 22,
	58, 47, 46, 45, 26, 78, 79, 62, 63, 64,
	33, 34, 81, 35, 74, 73, 60, 76, 75, 41,
	42, 43, 44, 37, 38, 39, 27, 51, 59, 66,
	50, 40, 36, 53, 32, 80, 31, 12, 11, 10,
	9, 84, 8, 7, 87, 88, 6, 5, 4, 89,
	3, 1,
}
var yyPact = [...]int{

	2, -1000, -1000, 2, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 10, -1, -18, 27, 55, 46, 9,
	8, 7, -1000, -1000, -1000, -1000, -1000, 6, -1000, -1000,
	6, 4, 1, -1000, -1000, -1000, 4, 57, 45, -1000,
	-6, -1000, -1000, -1000, -1000, 2, 2, 2, -1000, -19,
	-7, -1000, -8, -1000, -1000, -1000, -1000, -1000, -9, -1000,
	-1000, -1000, 36, 38, 37, 0, 18, 20, -1000, -1000,
	-1000, 4, 32, -12, -13, -1000, -1000, -1000, -19, -20,
	-14, -1000, 2, 2, -1000, -1000, -1000, -1000, 36, -1000,
}
var yyPgo = [...]int{

	0, 91, 1, 90, 88, 87, 86, 83, 82, 80,
	79, 78, 77, 76, 3, 74, 73, 72, 71, 0,
	70, 2, 36, 69, 66,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 4, 5, 24, 24, 20, 20, 22,
	21, 21, 21, 23, 6, 6, 7, 8, 8, 14,
	14, 14, 9, 9, 10, 19, 19, 19, 11, 12,
	18, 18, 18, 18, 17, 17, 17, 13, 15, 15,
	15, 16, 16, 16, 16,
}
var yyR2 = [...]int{

	0, 1, 0, 2, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 2, 3, 1, 1, 0, 1, 3,
	0, 1, 3, 3, 3, 4, 4, 4, 6, 1,
	1, 1, 2, 3, 5, 0, 3, 4, 4, 4,
	1, 1, 1, 1, 2, 2, 1, 2, 1, 1,
	1, 1, 1, 1, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -5, -6, -7, -8, -9,
	-10, -11, -12, 4, 7, 5, -14, 6, 12, 17,
	22, 21, 35, 8, 9, -2, 34, -24, 5, 35,
	35, -13, -15, 23, 24, 26, -17, 8, 9, 10,
	-18, 13, 14, 15, 16, 34, 34, 34, -22, 32,
	-20, -22, -14, -16, 28, 29, 30, 31, -14, 11,
	11, 34, -2, -2, -2, -21, -23, 35, 34, 34,
	34, 27, -19, 19, 18, 20, 20, 33, 27, 26,
	-14, 20, 34, 34, -21, 35, 34, -2, -2, -19,
}
var yyDef = [...]int{

	2, -2, 1, 2, 4, 5, 6, 7, 8, 9,
	10, 11, 12, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 29, 30, 31, 3, 13, 0, 15, 16,
	17, 0, 0, 48, 49, 50, 0, 0, 0, 46,
	32, 40, 41, 42, 43, 2, 2, 2, 14, 20,
	24, 18, 0, 47, 51, 52, 53, 54, 0, 44,
	45, 33, 35, 0, 0, 0, 21, 0, 25, 26,
	27, 0, 0, 0, 0, 38, 39, 19, 20, 0,
	0, 34, 2, 2, 22, 23, 28, 36, 35, 37,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:80
		{
			yylex.(*parseState).nodeList = yyDollar[1].nodeList
		}
	case 2:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line grammer.y:87
		{
			yyVAL.nodeList = nil
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:91
		{
			yyVAL.nodeList = &NodeList{yyDollar[1].node, yyDollar[2].nodeList}
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:110
		{
			yyVAL.node = &TitleNode{yyDollar[2].sval}
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:117
		{
			yyVAL.node = &StyleNode{yyDollar[2].sval, yyDollar[3].attrList}
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:123
		{
			yyVAL.sval = "participant"
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:124
		{
			yyVAL.sval = yyDollar[1].sval
		}
	case 17:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line grammer.y:129
		{
			yyVAL.attrList = nil
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:133
		{
			yyVAL.attrList = yyDollar[1].attrList
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:140
		{
			yyVAL.attrList = yyDollar[2].attrList
		}
	case 20:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line grammer.y:147
		{
			yyVAL.attrList = nil
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:151
		{
			yyVAL.attrList = &AttributeList{yyDollar[1].attr, nil}
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:155
		{
			yyVAL.attrList = &AttributeList{yyDollar[1].attr, yyDollar[3].attrList}
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:162
		{
			yyVAL.attr = &Attribute{yyDollar[1].sval, yyDollar[3].sval}
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:169
		{
			yyVAL.node = &ActorNode{yyDollar[2].sval, false, "", yyDollar[3].attrList}
		}
	case 25:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:173
		{
			yyVAL.node = &ActorNode{yyDollar[2].sval, true, yyDollar[4].sval, yyDollar[3].attrList}
		}
	case 26:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:180
		{
			yyVAL.node = &ActionNode{yyDollar[1].actorRef, yyDollar[3].actorRef, yyDollar[2].arrow, yyDollar[4].sval}
		}
	case 27:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:187
		{
			yyVAL.node = &NoteNode{yyDollar[3].actorRef, nil, yyDollar[2].noteAlign, yyDollar[4].sval}
		}
	case 28:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line grammer.y:191
		{
			yyVAL.node = &NoteNode{yyDollar[3].actorRef, yyDollar[5].actorRef, yyDollar[2].noteAlign, yyDollar[6].sval}
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:198
		{
			yyVAL.actorRef = NormalActorRef(yyDollar[1].sval)
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:202
		{
			yyVAL.actorRef = PseudoActorRef("left")
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:206
		{
			yyVAL.actorRef = PseudoActorRef("right")
		}
	case 32:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:213
		{
			yyVAL.node = &GapNode{yyDollar[2].dividerType, ""}
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:217
		{
			yyVAL.node = &GapNode{yyDollar[2].dividerType, yyDollar[3].sval}
		}
	case 34:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line grammer.y:224
		{
			yyVAL.node = &BlockNode{&BlockSegmentList{&BlockSegment{ALT_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, yyDollar[4].blockSegList}}
		}
	case 35:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line grammer.y:231
		{
			yyVAL.blockSegList = nil
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:235
		{
			yyVAL.blockSegList = &BlockSegmentList{&BlockSegment{ALT_ELSE_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, nil}
		}
	case 37:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:239
		{
			yyVAL.blockSegList = &BlockSegmentList{&BlockSegment{ALT_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, yyDollar[4].blockSegList}
		}
	case 38:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:246
		{
			yyVAL.node = &BlockNode{&BlockSegmentList{&BlockSegment{OPT_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, nil}}
		}
	case 39:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:253
		{
			yyVAL.node = &BlockNode{&BlockSegmentList{&BlockSegment{LOOP_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, nil}}
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:259
		{
			yyVAL.dividerType = SPACER_GAP
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:260
		{
			yyVAL.dividerType = EMPTY_GAP
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:261
		{
			yyVAL.dividerType = LINE_GAP
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:262
		{
			yyVAL.dividerType = FRAME_GAP
		}
	case 44:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:266
		{
			yyVAL.noteAlign = LEFT_NOTE_ALIGNMENT
		}
	case 45:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:267
		{
			yyVAL.noteAlign = RIGHT_NOTE_ALIGNMENT
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:268
		{
			yyVAL.noteAlign = OVER_NOTE_ALIGNMENT
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:273
		{
			yyVAL.arrow = ArrowType{yyDollar[1].arrowStem, yyDollar[2].arrowHead}
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:279
		{
			yyVAL.arrowStem = SOLID_ARROW_STEM
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:280
		{
			yyVAL.arrowStem = DASHED_ARROW_STEM
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:281
		{
			yyVAL.arrowStem = THICK_ARROW_STEM
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:285
		{
			yyVAL.arrowHead = SOLID_ARROW_HEAD
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:286
		{
			yyVAL.arrowHead = OPEN_ARROW_HEAD
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:287
		{
			yyVAL.arrowHead = BARBED_ARROW_HEAD
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:288
		{
			yyVAL.arrowHead = LOWER_BARBED_ARROW_HEAD
		}
	}
	goto yystack /* stack new state and value */
}
