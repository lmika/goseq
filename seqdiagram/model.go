// AST nodes used for by goseq
//

package seqdiagram

import (
    "io"

    "bitbucket.org/lmika/goseq/seqdiagram/parse"
)

// Top level diagram definition
type Diagram struct {
    ProcessingInstructions    []*ProcessingInstruction
    Title           string
    Actors          []*Actor
    Items           []SequenceItem
}

// Creates a new, empty diagram
func NewDiagram() *Diagram {
    return &Diagram{}
}

// Parses a diagram from a reader and returns the diagram or an error
func ParseDiagram(r io.Reader, filename string) (*Diagram, error) {
    //d := NewDiagram()
    nl, err := parse.Parse(r, filename)
    if err != nil {
        return nil, err
    }

    d := NewDiagram()
    tb := &treeBuilder{nl, filename}
    err = tb.buildTree(d)
    if err != nil {
        return nil, err
    }

    return d, nil
}

// Returns an actor by name.  If the actor is undefined, a new actor
// is created and added to the end of the slice.
func (d *Diagram) GetOrAddActor(name string) *Actor {
    return d.GetOrAddActorWithOptions(name, name)
}

func (d *Diagram) GetOrAddActorWithOptions(name string, label string) *Actor {
    for _, a := range d.Actors {
        if a.Name == name {
            return a
        }
    }

    na := &Actor{name, label, nil, len(d.Actors)}
    d.Actors = append(d.Actors, na)
    return na
}

// Adds a new sequence item
func (d *Diagram) AddSequenceItem(item SequenceItem) {
    d.Items = append(d.Items, item)
}

// Write the diagram as an SVG
func (d *Diagram) WriteSVG(w io.Writer) error {
    return d.WriteSVGWithStyle(w, DefaultStyle)
}

// Write the diagram as an SVG using a specific style
func (d *Diagram) WriteSVGWithStyle(w io.Writer, style *DiagramStyles) error {
    gb, err := newGraphicBuilder(d, style)
    if err != nil {
        return err
    }

    gb.buildGraphic().DrawSVG(w)
    return nil
}

// A processing instruction
type ProcessingInstruction struct {
    Prefix          string
    Value           string
}

// A participant
type Actor struct {
    Name            string
    Label           string
    Icon            ActorIcon
    rank            int
}

// Special actors
var LeftOffsideActor *Actor = &Actor{".left", ".left", nil, -1}
var RightOffsideActor *Actor = &Actor{".right", ".right", nil, -2}


// The supported arrow stems
type ArrowStem  int
const (
    SolidArrowStem  ArrowStem = iota
    DashedArrowStem           = iota
    ThickArrowStem            = iota
)

// The supported arrow heads
type ArrowHead  int
const (
    SolidArrowHead  ArrowHead = iota
    OpenArrowHead             = iota
    BarbArrowHead             = iota
    LowerBarbArrowHead        = iota
)


// An arrow
type Arrow struct {
    Stem        ArrowStem
    Head        ArrowHead
}

// Note alignments
type NoteAlignment int
const (
    LeftNoteAlignment       NoteAlignment = iota
    RightNoteAlignment                    = iota
    OverNoteAlignment                     = iota
)

// A sequence item
type SequenceItem interface {
}

// Defines a note
type Note struct {
    // The note's alignment and position
    Actor1      *Actor
    Actor2      *Actor
    
    Align       NoteAlignment

    // The message
    Message     string
}

// Defines an action
type Action struct {
    // The originating actor
    From        *Actor
    
    // The destination actor
    To          *Actor

    // The arrow to use
    Arrow       Arrow

    // The message
    Message     string
}

type DividerType int

const (
    DTSpacer    DividerType = iota
    DTGap   
    DTFrame
    DTLine    
)

// Defines a divider, which spans the diagram.
type Divider struct {
    // The message
    Message     string

    // The divider type
    Type        DividerType
}

// A framed block of sequence items.  Each block can have one or more segments,
// which will appear one after the other.
type Block struct {
    Segments    []*BlockSegment
}

// Returns the maximum number of nested blocks within the segments
func (b *Block) MaxNestDepth() int {
    nestDepth := 0
    for _, seg := range b.Segments {
        nestDepth = maxInt(nestDepth, seg.MaxNestDepth())
    }
    return nestDepth + 1
}

// The type of segment
type SegmentType int
const (
    // The alt segment
    AltSegmentType SegmentType = iota

    // The else segment
    ElseSegmentType
)

// A segment within a block
type BlockSegment struct {
    Type        SegmentType
    Prefix      string
    Message     string
    SubItems    []SequenceItem
}

// Returns the number of nested blocks 
func (bs *BlockSegment) MaxNestDepth() int {
    nestDepth := 0
    for _, subItem := range bs.SubItems {
        if block, isBlock := subItem.(*Block) ; isBlock {
            nestDepth = maxInt(nestDepth, block.MaxNestDepth())
        }
    }
    return nestDepth
}
