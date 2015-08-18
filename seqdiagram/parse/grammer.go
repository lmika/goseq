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
	"--": DOUBLEDASH,
	"-":  DASH,
	"=":  EQUAL,

	">>":  DOUBLEANGR,
	">":   ANGR,
	"/>":  SLASHANGR,
	"\\>": BACKSLASHANGR,
}

//line grammer.y:31
type yySymType struct {
	yys         int
	nodeList    *NodeList
	node        Node
	arrow       ArrowType
	arrowStem   ArrowStemType
	arrowHead   ArrowHeadType
	noteAlign   NoteAlignment
	dividerType GapType

	sval string
}

const K_TITLE = 57346
const K_PARTICIPANT = 57347
const K_NOTE = 57348
const K_LEFT = 57349
const K_RIGHT = 57350
const K_OVER = 57351
const K_OF = 57352
const K_HORIZONTAL = 57353
const K_SPACER = 57354
const K_GAP = 57355
const K_LINE = 57356
const K_FRAME = 57357
const K_IF = 57358
const K_END = 57359
const DASH = 57360
const DOUBLEDASH = 57361
const EQUAL = 57362
const ANGR = 57363
const DOUBLEANGR = 57364
const BACKSLASHANGR = 57365
const SLASHANGR = 57366
const MESSAGE = 57367
const IDENT = 57368

var yyToknames = []string{
	"K_TITLE",
	"K_PARTICIPANT",
	"K_NOTE",
	"K_LEFT",
	"K_RIGHT",
	"K_OVER",
	"K_OF",
	"K_HORIZONTAL",
	"K_SPACER",
	"K_GAP",
	"K_LINE",
	"K_FRAME",
	"K_IF",
	"K_END",
	"DASH",
	"DOUBLEDASH",
	"EQUAL",
	"ANGR",
	"DOUBLEANGR",
	"BACKSLASHANGR",
	"SLASHANGR",
	"MESSAGE",
	"IDENT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line grammer.y:174

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
		case '-', '>', '*', '=', '/', '\\':
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
	case "horizontal":
		return K_HORIZONTAL
	case "if":
		return K_IF
	case "end":
		return K_END
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
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 33
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 48

var yyAct = []int{

	2, 10, 11, 13, 16, 41, 35, 18, 14, 37,
	38, 39, 40, 15, 47, 46, 44, 34, 33, 17,
	21, 22, 23, 12, 29, 30, 31, 32, 48, 25,
	26, 27, 28, 43, 45, 42, 24, 36, 20, 19,
	9, 8, 7, 6, 5, 4, 3, 1,
}
var yyPact = []int{

	-3, -1000, -1000, -3, -1000, -1000, -1000, -1000, -1000, -1000,
	-6, -19, 2, 22, 12, -7, -1000, -1000, -8, -20,
	-12, -1000, -1000, -1000, -21, 25, 23, -1000, -9, -1000,
	-1000, -1000, -1000, -3, -1000, -10, -1000, -1000, -1000, -1000,
	-1000, -11, -1000, -1000, -1000, 11, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 47, 0, 46, 45, 44, 43, 42, 41, 40,
	39, 38, 37, 36, 32,
}
var yyR1 = []int{

	0, 1, 2, 2, 3, 3, 3, 3, 3, 3,
	4, 5, 5, 6, 7, 8, 8, 9, 14, 14,
	14, 14, 13, 13, 13, 10, 11, 11, 11, 12,
	12, 12, 12,
}
var yyR2 = []int{

	0, 1, 0, 2, 1, 1, 1, 1, 1, 1,
	2, 2, 3, 4, 4, 2, 3, 4, 1, 1,
	1, 1, 2, 2, 1, 2, 1, 1, 1, 1,
	1, 1, 1,
}
var yyChk = []int{

	-1000, -1, -2, -3, -4, -5, -6, -7, -8, -9,
	4, 5, 26, 6, 11, 16, -2, 25, 26, -10,
	-11, 18, 19, 20, -13, 7, 8, 9, -14, 12,
	13, 14, 15, 25, 25, 26, -12, 21, 22, 23,
	24, 26, 10, 10, 25, -2, 25, 25, 17,
}
var yyDef = []int{

	2, -2, 1, 2, 4, 5, 6, 7, 8, 9,
	0, 0, 0, 0, 0, 0, 3, 10, 11, 0,
	0, 26, 27, 28, 0, 0, 0, 24, 15, 18,
	19, 20, 21, 2, 12, 0, 25, 29, 30, 31,
	32, 0, 22, 23, 16, 0, 13, 14, 17,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
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

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
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
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
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
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
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
			if yyn < 0 || yyn == yychar {
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
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
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
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
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
		//line grammer.y:67
		{
			yylex.(*parseState).nodeList = yyS[yypt-0].nodeList
		}
	case 2:
		//line grammer.y:74
		{
			yyVAL.nodeList = nil
		}
	case 3:
		//line grammer.y:78
		{
			yyVAL.nodeList = &NodeList{yyS[yypt-1].node, yyS[yypt-0].nodeList}
		}
	case 4:
		yyVAL.node = yyS[yypt-0].node
	case 5:
		yyVAL.node = yyS[yypt-0].node
	case 6:
		yyVAL.node = yyS[yypt-0].node
	case 7:
		yyVAL.node = yyS[yypt-0].node
	case 8:
		yyVAL.node = yyS[yypt-0].node
	case 9:
		yyVAL.node = yyS[yypt-0].node
	case 10:
		//line grammer.y:94
		{
			yyVAL.node = &TitleNode{yyS[yypt-0].sval}
		}
	case 11:
		//line grammer.y:101
		{
			yyVAL.node = &ActorNode{yyS[yypt-0].sval, false, ""}
		}
	case 12:
		//line grammer.y:105
		{
			yyVAL.node = &ActorNode{yyS[yypt-1].sval, true, yyS[yypt-0].sval}
		}
	case 13:
		//line grammer.y:112
		{
			yyVAL.node = &ActionNode{yyS[yypt-3].sval, yyS[yypt-1].sval, yyS[yypt-2].arrow, yyS[yypt-0].sval}
		}
	case 14:
		//line grammer.y:119
		{
			yyVAL.node = &NoteNode{yyS[yypt-1].sval, yyS[yypt-2].noteAlign, yyS[yypt-0].sval}
		}
	case 15:
		//line grammer.y:126
		{
			yyVAL.node = &GapNode{yyS[yypt-0].dividerType, ""}
		}
	case 16:
		//line grammer.y:130
		{
			yyVAL.node = &GapNode{yyS[yypt-1].dividerType, yyS[yypt-0].sval}
		}
	case 17:
		//line grammer.y:137
		{
			yyVAL.node = &BlockNode{yyS[yypt-2].sval, yyS[yypt-1].nodeList}
		}
	case 18:
		//line grammer.y:143
		{
			yyVAL.dividerType = SPACER_GAP
		}
	case 19:
		//line grammer.y:144
		{
			yyVAL.dividerType = EMPTY_GAP
		}
	case 20:
		//line grammer.y:145
		{
			yyVAL.dividerType = LINE_GAP
		}
	case 21:
		//line grammer.y:146
		{
			yyVAL.dividerType = FRAME_GAP
		}
	case 22:
		//line grammer.y:150
		{
			yyVAL.noteAlign = LEFT_NOTE_ALIGNMENT
		}
	case 23:
		//line grammer.y:151
		{
			yyVAL.noteAlign = RIGHT_NOTE_ALIGNMENT
		}
	case 24:
		//line grammer.y:152
		{
			yyVAL.noteAlign = OVER_NOTE_ALIGNMENT
		}
	case 25:
		//line grammer.y:157
		{
			yyVAL.arrow = ArrowType{yyS[yypt-1].arrowStem, yyS[yypt-0].arrowHead}
		}
	case 26:
		//line grammer.y:163
		{
			yyVAL.arrowStem = SOLID_ARROW_STEM
		}
	case 27:
		//line grammer.y:164
		{
			yyVAL.arrowStem = DASHED_ARROW_STEM
		}
	case 28:
		//line grammer.y:165
		{
			yyVAL.arrowStem = THICK_ARROW_STEM
		}
	case 29:
		//line grammer.y:169
		{
			yyVAL.arrowHead = SOLID_ARROW_HEAD
		}
	case 30:
		//line grammer.y:170
		{
			yyVAL.arrowHead = OPEN_ARROW_HEAD
		}
	case 31:
		//line grammer.y:171
		{
			yyVAL.arrowHead = BARBED_ARROW_HEAD
		}
	case 32:
		//line grammer.y:172
		{
			yyVAL.arrowHead = LOWER_BARBED_ARROW_HEAD
		}
	}
	goto yystack /* stack new state and value */
}
