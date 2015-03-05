//line goseq/grammer.y:6
package goseq

import __yyfmt__ "fmt"

//line goseq/grammer.y:6
import (
	"bytes"
	"errors"
	"io"
	"strings"
	"text/scanner"
)

var DualRunes = map[string]int{
	"--": DOUBLEDASH,
	"-":  DASH,

	">>": DOUBLEANGR,
	">":  ANGR,
	"*>": STARANGR,
}

//line goseq/grammer.y:28
type yySymType struct {
	yys       int
	seqItem   SequenceItem
	arrow     Arrow
	arrowStem ArrowStem
	arrowHead ArrowHead
	noteAlign NoteAlignment

	sval string
}

const K_TITLE = 57346
const K_PARTICIPANT = 57347
const K_NOTE = 57348
const K_LEFT = 57349
const K_RIGHT = 57350
const K_OVER = 57351
const K_OF = 57352
const K_GAP = 57353
const DASH = 57354
const DOUBLEDASH = 57355
const ANGR = 57356
const DOUBLEANGR = 57357
const STARANGR = 57358
const MESSAGE = 57359
const IDENT = 57360

var yyToknames = []string{
	"K_TITLE",
	"K_PARTICIPANT",
	"K_NOTE",
	"K_LEFT",
	"K_RIGHT",
	"K_OVER",
	"K_OF",
	"K_GAP",
	"DASH",
	"DOUBLEDASH",
	"ANGR",
	"DOUBLEANGR",
	"STARANGR",
	"MESSAGE",
	"IDENT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line goseq/grammer.y:169

// Manages the lexer as well as the current diagram being parsed
type parseState struct {
	S       scanner.Scanner
	err     error
	atEof   bool
	diagram *Diagram
}

func newParseState(src io.Reader) *parseState {
	ps := &parseState{}
	ps.S.Init(src)
	ps.diagram = &Diagram{}

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
		case '-', '>', '*':
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
	case "gap":
		return K_GAP
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
	r := ps.NextRune()
	for (r != '\n') && (r != scanner.EOF) {
		r = ps.NextRune()
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
	ps.err = errors.New(err)
}

func Parse(reader io.Reader) (*Diagram, error) {
	ps := newParseState(reader)
	yyParse(ps)

	if ps.err != nil {
		return nil, ps.err
	} else {
		return ps.diagram, nil
	}
}

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 25
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 37

var yyAct = []int{

	7, 8, 13, 33, 30, 31, 32, 14, 5, 28,
	17, 37, 36, 35, 12, 27, 26, 16, 20, 21,
	23, 24, 25, 2, 34, 4, 3, 15, 1, 22,
	29, 19, 18, 11, 10, 9, 6,
}
var yyPact = []int{

	-4, -1000, -1000, -4, -1000, -1000, -1000, 0, -8, -1000,
	-1000, -1000, 6, 13, -1, -1000, -1000, -2, -9, -10,
	-1000, -1000, -15, 14, 3, -1000, -1000, -1000, -5, -1000,
	-1000, -1000, -1000, -6, -1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 36, 35, 34, 33, 32, 31, 30, 29, 28,
	23, 26, 25, 8,
}
var yyR1 = []int{

	0, 9, 10, 10, 11, 11, 11, 12, 13, 13,
	1, 1, 1, 2, 3, 4, 8, 8, 8, 5,
	6, 6, 7, 7, 7,
}
var yyR2 = []int{

	0, 1, 0, 2, 1, 1, 1, 2, 2, 3,
	1, 1, 1, 4, 4, 2, 2, 2, 1, 2,
	1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -9, -10, -11, -12, -13, -1, 4, 5, -2,
	-3, -4, 18, 6, 11, -10, 17, 18, -5, -6,
	12, 13, -8, 7, 8, 9, 17, 17, 18, -7,
	14, 15, 16, 18, 10, 10, 17, 17,
}
var yyDef = []int{

	2, -2, 1, 2, 4, 5, 6, 0, 0, 10,
	11, 12, 0, 0, 0, 3, 7, 8, 0, 0,
	20, 21, 0, 0, 0, 18, 15, 9, 0, 19,
	22, 23, 24, 0, 16, 17, 13, 14,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18,
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

	case 6:
		//line goseq/grammer.y:70
		{
			yylex.(*parseState).diagram.AddSequenceItem(yyS[yypt-0].seqItem)
		}
	case 7:
		//line goseq/grammer.y:77
		{
			yylex.(*parseState).diagram.Title = yyS[yypt-0].sval
		}
	case 8:
		//line goseq/grammer.y:84
		{
			yylex.(*parseState).diagram.GetOrAddActor(yyS[yypt-0].sval)
		}
	case 9:
		//line goseq/grammer.y:88
		{
			yylex.(*parseState).diagram.GetOrAddActorWithOptions(yyS[yypt-1].sval, yyS[yypt-0].sval)
		}
	case 10:
		yyVAL.seqItem = yyS[yypt-0].seqItem
	case 11:
		yyVAL.seqItem = yyS[yypt-0].seqItem
	case 12:
		yyVAL.seqItem = yyS[yypt-0].seqItem
	case 13:
		//line goseq/grammer.y:101
		{
			d := yylex.(*parseState).diagram
			yyVAL.seqItem = &Action{d.GetOrAddActor(yyS[yypt-3].sval), d.GetOrAddActor(yyS[yypt-1].sval), yyS[yypt-2].arrow, yyS[yypt-0].sval}
		}
	case 14:
		//line goseq/grammer.y:109
		{
			d := yylex.(*parseState).diagram
			yyVAL.seqItem = &Note{d.GetOrAddActor(yyS[yypt-1].sval), yyS[yypt-2].noteAlign, yyS[yypt-0].sval}
		}
	case 15:
		//line goseq/grammer.y:117
		{
			yyVAL.seqItem = &Divider{yyS[yypt-0].sval}
		}
	case 16:
		//line goseq/grammer.y:124
		{
			yyVAL.noteAlign = LeftNoteAlignment
		}
	case 17:
		//line goseq/grammer.y:128
		{
			yyVAL.noteAlign = RightNoteAlignment
		}
	case 18:
		//line goseq/grammer.y:132
		{
			yyVAL.noteAlign = OverNoteAlignment
		}
	case 19:
		//line goseq/grammer.y:139
		{
			yyVAL.arrow = Arrow{yyS[yypt-1].arrowStem, yyS[yypt-0].arrowHead}
		}
	case 20:
		//line goseq/grammer.y:146
		{
			yyVAL.arrowStem = SolidArrowStem
		}
	case 21:
		//line goseq/grammer.y:150
		{
			yyVAL.arrowStem = DashedArrowStem
		}
	case 22:
		//line goseq/grammer.y:157
		{
			yyVAL.arrowHead = SolidArrowHead
		}
	case 23:
		//line goseq/grammer.y:161
		{
			yyVAL.arrowHead = OpenArrowHead
		}
	case 24:
		//line goseq/grammer.y:165
		{
			yyVAL.arrowHead = BarbArrowHead
		}
	}
	goto yystack /* stack new state and value */
}
