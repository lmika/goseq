// AST nodes used for by goseq
//

package seqdiagram

import (
    "io"

    "bitbucket.org/lmika/goseq/seqdiagram/parse"
)

// Top level diagram definition
type Diagram struct {
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

    na := &Actor{name, label, len(d.Actors)}
    d.Actors = append(d.Actors, na)
    return na
}

// Adds a new sequence item
func (d *Diagram) AddSequenceItem(item SequenceItem) {
    d.Items = append(d.Items, item)
}

// Write the diagram as an SVG
func (d *Diagram) WriteSVG(w io.Writer) error {
    gb, err := newGraphicBuilder(d, DefaultStyle)
    if err != nil {
        return err
    }

    gb.buildGraphic().DrawSVG(w)
    return nil
}

// A participant
type Actor struct {
    Name            string
    Label           string
    rank            int
}


// The supported arrow stems
type ArrowStem  int
const (
    SolidArrowStem  ArrowStem = iota
    DashedArrowStem           = iota
)

// The supported arrow heads
type ArrowHead  int
const (
    SolidArrowHead  ArrowHead = iota
    OpenArrowHead             = iota
    BarbArrowHead             = iota
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
    Actor       *Actor
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
    DTGap   DividerType = iota
    DTLine
)

// Defines a divider, which spans the diagram.
type Divider struct {
    // The message
    Message     string

    // The divider type
    Type        DividerType
}