// Gramma for goseq
//
// Based on the gramma used for js-sequence-diagram

%{
package parse

import (
    "io"
    "bytes"
    "errors"
    "strings"
    "fmt"
    "text/scanner"
)

var DualRunes = map[string]int {
    "--":   DOUBLEDASH,
    "-":    DASH,
    "=":    EQUAL,

    ">>":   DOUBLEANGR,
    ">":    ANGR,
    "/>":   SLASHANGR,
    "\\>":  BACKSLASHANGR,
}


%}

%union {
    nodeList        *NodeList
    node            Node
    arrow           ArrowType
    arrowStem       ArrowStemType
    arrowHead       ArrowHeadType
    noteAlign       NoteAlignment
    dividerType     GapType

    sval            string
}

%token  K_TITLE K_PARTICIPANT K_NOTE
%token  K_LEFT  K_RIGHT  K_OVER  K_OF
%token  K_HORIZONTAL K_GAP K_LINE K_FRAME

%token  DASH    DOUBLEDASH      EQUAL
%token  ANGR    DOUBLEANGR      BACKSLASHANGR       SLASHANGR

%token  <sval>  MESSAGE
%token  <sval>  IDENT

%type   <nodeList>      top decls
%type   <node>          decl
%type   <node>          title actor action note gap
%type   <arrow>         arrow
%type   <arrowStem>     arrowStem
%type   <arrowHead>     arrowHead
%type   <noteAlign>     noteplace
%type   <dividerType>   dividerType

%%

top         
    :   decls
    {
        yylex.(*parseState).nodeList = $1
    }
    ;

decls       
    :   /* empty */
    {
        $$ = nil
    }
    |   decl decls
    {
        $$ = &NodeList{$1, $2}
    }
    ;

decl
    :   title
    |   actor
    |   action
    |   note
    |   gap
    ;

title
    :   K_TITLE MESSAGE
    {
        $$ = &TitleNode{$2}
    }
    ;

actor
    :   K_PARTICIPANT IDENT
    {
        $$ = &ActorNode{$2, false, ""}
    }
    |   K_PARTICIPANT IDENT MESSAGE
    {
        $$ = &ActorNode{$2, true, $3}
    }
    ;

action
    :   IDENT arrow IDENT MESSAGE
    {
        $$ = &ActionNode{$1, $3, $2, $4}
    }
    ;

note
    :   K_NOTE noteplace IDENT MESSAGE
    {
        $$ = &NoteNode{$3, $2, $4}
    }
    ;

gap
    :   K_HORIZONTAL dividerType
    {
        $$ = &GapNode{$2, ""}
    }
    |   K_HORIZONTAL dividerType MESSAGE
    {
        $$ = &GapNode{$2, $3}
    }
    ;

dividerType
    :   K_GAP               { $$ = EMPTY_GAP }
    |   K_LINE              { $$ = LINE_GAP }
    |   K_FRAME             { $$ = FRAME_GAP }
    ;

noteplace
    :   K_LEFT K_OF         { $$ = LEFT_NOTE_ALIGNMENT }
    |   K_RIGHT K_OF        { $$ = RIGHT_NOTE_ALIGNMENT }
    |   K_OVER              { $$ = OVER_NOTE_ALIGNMENT }
    ;

arrow
    :   arrowStem   arrowHead
    {
        $$ = ArrowType{$1, $2}
    }
    ;

arrowStem
    :   DASH                { $$ = SOLID_ARROW_STEM }
    |   DOUBLEDASH          { $$ = DASHED_ARROW_STEM }    
    |   EQUAL               { $$ = THICK_ARROW_STEM }
    ;

arrowHead
    :   ANGR                { $$ = SOLID_ARROW_HEAD }
    |   DOUBLEANGR          { $$ = OPEN_ARROW_HEAD }
    |   BACKSLASHANGR       { $$ = BARBED_ARROW_HEAD }
    |   SLASHANGR           { $$ = LOWER_BARBED_ARROW_HEAD }
    ;
%%

// Manages the lexer as well as the current diagram being parsed
type parseState struct {
    S           scanner.Scanner
    err         error
    atEof       bool
    //diagram     *Diagram
    procInstrs  []string
    nodeList    *NodeList
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
    case "frame":
        return K_FRAME
    case "line":
        return K_LINE
    case "horizontal":
        return K_HORIZONTAL
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
    var buf *bytes.Buffer

    r := ps.NextRune()
    if (r == '!') {
        // This starts a processor instruction
        buf = new(bytes.Buffer)
        r = ps.NextRune()
    }

    for ((r != '\n') && (r != scanner.EOF)) {
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