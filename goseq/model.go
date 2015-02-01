// AST nodes used for by goseq
//

package goseq

import (
    "io"
)

// Top level diagram definition
type Diagram struct {
    Title           string
    Actors          []*Actor
    Items           []SequenceItem
}

// Returns an actor by name.  If the actor is undefined, a new actor
// is created and added to the end of the slice.
func (d *Diagram) GetOrAddActor(name string) *Actor {
    for _, a := range d.Actors {
        if a.Name == name {
            return a
        }
    }

    na := &Actor{name, len(d.Actors)}
    d.Actors = append(d.Actors, na)
    return na
}

// Adds a new sequence item
func (d *Diagram) AddSequenceItem(item SequenceItem) {
    d.Items = append(d.Items, item)
}

// Write the diagram as an SVG
func (d *Diagram) WriteSVG(w io.Writer) error {
    gb, err := NewGraphicBuilder(d)
    if err != nil {
        return err
    }

    gb.BuildGraphic().DrawSVG(w)
    return nil
}

// A participant
type Actor struct {
    Name            string
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