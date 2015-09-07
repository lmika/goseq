//go:generate go tool yacc -o grammer.go grammer.y
//
package parse

type ArrowStemType int
const (
    SOLID_ARROW_STEM ArrowStemType = iota
    DASHED_ARROW_STEM              = iota
    THICK_ARROW_STEM               = iota
)

type ArrowHeadType int
const (
    SOLID_ARROW_HEAD ArrowHeadType = iota
    OPEN_ARROW_HEAD                = iota
    BARBED_ARROW_HEAD              = iota
    LOWER_BARBED_ARROW_HEAD        = iota
)

type SegmentType int
const (
    ALT_SEGMENT   SegmentType      = iota
    ALT_ELSE_SEGMENT               = iota
    OPT_SEGMENT                    = iota
    LOOP_SEGMENT                   = iota
)

type ArrowType struct {
    Stem        ArrowStemType
    Head        ArrowHeadType
}


// A list of declaration node
type NodeList struct {
    Head        Node
    Tail        *NodeList
}

// A type of declaration node
type Node interface {
}

// A processing instruction node
type ProcessInstructionNode struct {
    Prefix string
    Value string
}

// A title declaration node
type TitleNode struct {
    Title       string
}

// An actor declaration node
type ActorNode struct {
    // Identifier
    Ident       string

    // Returns true if the actor has a separate description
    HasDescr    bool

    // Description
    Descr       string

    // Attributes
    Attributes  *AttributeList
}

// Returns a suitable actor name.  This can either be the description if HasDescr is true
// or the ident value if HasDescr is false.
func (an *ActorNode) ActorName() string {
    if an.HasDescr {
        return an.Descr
    } else {
        return an.Ident
    }
}

// A reference to an actor
type ActorRef interface {
}

// A reference to a normal actor
type NormalActorRef string

// A reference to a pseudo actor
type PseudoActorRef string


// An action node
type ActionNode struct {
    From        ActorRef
    To          ActorRef
    Arrow       ArrowType
    Descr       string
}

// Note node
type NoteAlignment    int
const (
    LEFT_NOTE_ALIGNMENT NoteAlignment = iota
    RIGHT_NOTE_ALIGNMENT              = iota
    OVER_NOTE_ALIGNMENT               = iota
)

type NoteNode struct {
    Actor1      ActorRef
    Actor2      ActorRef        // Can be nil
    
    Position    NoteAlignment
    Descr       string
}

// Gap node
type GapType    int
const (
    SPACER_GAP GapType = iota
    EMPTY_GAP       = iota
    LINE_GAP        = iota
    FRAME_GAP       = iota
)

type GapNode struct {
    Type        GapType
    Descr       string
}

// A block node.  Each block can have one or more segments
type BlockNode struct {
    Segments    *BlockSegmentList
}

type BlockSegmentList struct {
    Head        *BlockSegment
    Tail        *BlockSegmentList
}

type BlockSegment struct {
    Type        SegmentType
    Prefix      string
    Message     string
    SubNodes    *NodeList
}


// Attributes
type Attribute struct {
    Name        string
    Value       string
}

type AttributeList struct {
    Head        *Attribute
    Tail        *AttributeList
}