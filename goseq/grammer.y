// Gramma for goseq
//
// Based on the gramma used for js-sequence-diagram

%{
package goseq

import (
    "io"
    "bytes"
    "errors"
    "strings"
    "text/scanner"
)

var DualRunes = map[string]int {
    "--":   DOUBLEDASH,
    "-":    DASH,

    ">>":   DOUBLEANGR,
    ">":    ANGR,
    "*>":   STARANGR,
}


%}

%union {
    seqItem     SequenceItem
    arrow       Arrow
    arrowStem   ArrowStem
    arrowHead   ArrowHead
    noteAlign   NoteAlignment

    sval        string
}

%token  K_TITLE K_PARTICIPANT K_NOTE
%token  K_LEFT  K_RIGHT  K_OVER  K_OF
%token  K_GAP

%token  DASH    DOUBLEDASH
%token  ANGR    DOUBLEANGR      STARANGR

%token  <sval>  MESSAGE
%token  <sval>  IDENT

%type   <seqItem>   seqitem
%type   <seqItem>   action      note    gap
%type   <arrow>     arrow
%type   <arrowStem> arrowStem
%type   <arrowHead> arrowHead
%type   <noteAlign> noteplace

%%

top         
    :   decls
    ;

decls       
    :   /* empty */
    |   decl decls
    ;

decl
    :   title
    |   actor
    |   seqitem
    {
        yylex.(*parseState).diagram.AddSequenceItem($1)
    }
    ;

title
    :   K_TITLE MESSAGE
    {
        yylex.(*parseState).diagram.Title = $2
    }
    ;

actor
    :   K_PARTICIPANT IDENT
    {
        yylex.(*parseState).diagram.GetOrAddActor($2)
    }
    |   K_PARTICIPANT IDENT MESSAGE
    {
        yylex.(*parseState).diagram.GetOrAddActorWithOptions($2, $3)
    }
    ;

seqitem
    :   action
    |   note
    |   gap
    ;

action
    :   IDENT arrow IDENT MESSAGE
    {
        d := yylex.(*parseState).diagram
        $$ = &Action{d.GetOrAddActor($1), d.GetOrAddActor($3), $2, $4}
    }
    ;

note
    :   K_NOTE noteplace IDENT MESSAGE
    {
        d := yylex.(*parseState).diagram
        $$ = &Note{d.GetOrAddActor($3), $2, $4}
    }
    ;

gap
    :   K_GAP MESSAGE
    {
        $$ = &Divider{$2}
    }
    ;

noteplace
    :   K_LEFT K_OF
    {
        $$ = LeftNoteAlignment
    }
    |   K_RIGHT K_OF
    {
        $$ = RightNoteAlignment
    }
    |   K_OVER
    {
        $$ = OverNoteAlignment
    }
    ;

arrow
    :   arrowStem   arrowHead
    {
        $$ = Arrow{$1, $2}
    }
    ;

arrowStem
    :   DASH
    {
        $$ = SolidArrowStem
    }
    |   DOUBLEDASH
    {
        $$ = DashedArrowStem
    }    
    ;

arrowHead
    :   ANGR
    {
        $$ = SolidArrowHead
    }
    |   DOUBLEANGR
    {
        $$ = OpenArrowHead
    }
    |   STARANGR
    {
        $$ = BarbArrowHead
    }
    ;
%%

// Manages the lexer as well as the current diagram being parsed
type parseState struct {
    S           scanner.Scanner
    err         error
    atEof       bool
    diagram     *Diagram
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
            if res, isTok := ps.handleDoubleRune(tok) ; isTok {
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
        if tok, hasTok := DualRunes[tokStr] ; hasTok {
            ps.NextRune()
            return tok, true
        }
    }

    // Try the single rune
    tokStr := string(firstRune)
    if tok, hasTok := DualRunes[tokStr] ; hasTok {
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
    for ((r != '\n') && (r != scanner.EOF)) {
        if (r == '\\') {
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
    for ((r != '\n') && (r != scanner.EOF)) {
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