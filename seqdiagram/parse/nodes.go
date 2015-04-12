package parse

type ArrowStemType int
const (
    SOLID_ARROW_STEM ArrowStemType = iota
    DASHED_ARROW_STEM              = iota
)

type ArrowHeadType int
const (
    SOLID_ARROW_HEAD ArrowHeadType = iota
    OPEN_ARROW_HEAD                = iota
    BARBED_ARROW_HEAD              = iota
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

// An action node
type ActionNode struct {
    From        string
    To          string
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
    Actor       string
    Position    NoteAlignment
    Descr       string
}

// Gap node
type GapType    int
const (
    EMPTY_GAP GapType = iota
    LINE_GAP        = iota
)

type GapNode struct {
    Type        GapType
    Descr       string
}