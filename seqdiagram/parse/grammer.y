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
    ".":    DOT,
    ",":    COMMA,

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
    actorRef        ActorRef
    noteAlign       NoteAlignment
    dividerType     GapType
    blockSegList    *BlockSegmentList
    attrList        *AttributeList
    attr            *Attribute

    sval            string
}

%token  K_TITLE K_PARTICIPANT K_NOTE K_STYLE
%token  K_LEFT  K_RIGHT  K_OVER  K_OF
%token  K_HORIZONTAL K_SPACER   K_GAP K_LINE K_FRAME
%token  K_ALT   K_ELSEALT   K_ELSE   K_END  K_LOOP K_OPT

%token  DASH    DOUBLEDASH      DOT                 EQUAL       COMMA
%token  ANGR    DOUBLEANGR      BACKSLASHANGR       SLASHANGR
%token  PARL    PARR

%token  <sval>  MESSAGE
%token  <sval>  IDENT

%type   <nodeList>      top decls
%type   <node>          decl
%type   <node>          title style actor action note gap altblock optblock loopblock
%type   <arrow>         arrow
%type   <actorRef>      actorref
%type   <arrowStem>     arrowStem
%type   <arrowHead>     arrowHead
%type   <noteAlign>     noteplace
%type   <dividerType>   dividerType
%type   <blockSegList>  altblocklist
%type   <attrList>      maybeattrs attrs attrset
%type   <attr>          attr

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
    |   style
    |   actor
    |   action
    |   note
    |   gap
    |   altblock
    |   optblock
    |   loopblock
    ;

title
    :   K_TITLE MESSAGE
    {
        $$ = &TitleNode{$2}
    }
    ;

style
    :   K_STYLE IDENT attrset
    {
        $$ = &StyleNode{$2, $3}
    }
    ;

maybeattrs
    :   /* empty */
    {
        $$ = nil
    }
    |   attrset
    {
        $$ = $1;
    }
    ;

attrset
    :   PARL attrs PARR
    {
        $$ = $2;
    }
    ;

attrs
    :   /* empty */
    {
        $$ = nil
    }
    |   attr
    {
        $$ = &AttributeList{$1, nil}
    }
    |   attr COMMA attrs
    {
        $$ = &AttributeList{$1, $3}
    }
    ;

attr
    :   IDENT EQUAL IDENT
    {
        $$ = &Attribute{$1, $3}
    }
    ;

actor
    :   K_PARTICIPANT IDENT maybeattrs
    {
        $$ = &ActorNode{$2, false, "", $3}
    }
    |   K_PARTICIPANT IDENT maybeattrs MESSAGE
    {
        $$ = &ActorNode{$2, true, $4, $3}
    }
    ;

action
    :   actorref arrow actorref MESSAGE
    {
        $$ = &ActionNode{$1, $3, $2, $4}
    }
    ;

note
    :   K_NOTE noteplace actorref MESSAGE
    {
        $$ = &NoteNode{$3, nil, $2, $4}
    }
    |   K_NOTE noteplace actorref COMMA actorref MESSAGE
    {
        $$ = &NoteNode{$3, $5, $2, $6}
    }
    ;

actorref
    :   IDENT
    {
        $$ = NormalActorRef($1)
    }
    |   K_LEFT
    {
        $$ = PseudoActorRef("left")
    }
    |   K_RIGHT
    {
        $$ = PseudoActorRef("right")
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

altblock
    :   K_ALT MESSAGE decls altblocklist K_END
    {
        $$ = &BlockNode{&BlockSegmentList{&BlockSegment{ALT_SEGMENT, "", $2, $3}, $4}}
    }
    ;

altblocklist
    :   /* empty */
    {
        $$ = nil
    }
    |   K_ELSE MESSAGE decls
    {
        $$ = &BlockSegmentList{&BlockSegment{ALT_ELSE_SEGMENT, "", $2, $3}, nil}
    }
    |   K_ELSEALT MESSAGE decls altblocklist
    {
        $$ = &BlockSegmentList{&BlockSegment{ALT_SEGMENT, "", $2, $3}, $4}
    }
    ;

optblock
    :   K_OPT MESSAGE decls K_END
    {
        $$ = &BlockNode{&BlockSegmentList{&BlockSegment{OPT_SEGMENT, "", $2, $3}, nil}}
    }
    ;

loopblock
    :   K_LOOP MESSAGE decls K_END
    {
        $$ = &BlockNode{&BlockSegmentList{&BlockSegment{LOOP_SEGMENT, "", $2, $3}, nil}}
    }
    ;

dividerType
    :   K_SPACER            { $$ = SPACER_GAP }
    |   K_GAP               { $$ = EMPTY_GAP }
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
        case '(':
            return PARL
        case ')':
            return PARR
        case '-', '>', '*', '=', '/', '\\', '.', ',':
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
