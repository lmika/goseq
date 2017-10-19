// AST nodes used for by goseq
//

package seqdiagram

import (
	"io"

	"github.com/lmika/goseq/seqdiagram/parse"
)

// Top level diagram definition
type Diagram struct {
	ProcessingInstructions []*ProcessingInstruction
	Title                  string
	Actors                 []*Actor
	Items                  []SequenceItem
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
	tb := newTreeBuilder(nl, filename)
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

	na := &Actor{
		Name:      name,
		Label:     label,
		InHeader:  true,
		InFooter:  true,
		Lifeline:  true,
		Color:     "black",
		TextColor: "black",
		rank:      len(d.Actors),
	}
	d.Actors = append(d.Actors, na)
	return na
}

// Adds a new sequence item
func (d *Diagram) AddSequenceItem(item SequenceItem) {
	d.Items = append(d.Items, item)
}

// Write the diagram as an SVG
func (d *Diagram) WriteSVG(w io.Writer) error {
	return d.WriteSVGWithOptions(w, DefaultOptions)
}

// Write the diagram as an SVG using a specific style
func (d *Diagram) WriteSVGWithOptions(w io.Writer, options *ImageOptions) error {
	gb, err := newGraphicBuilder(d, options.Style)
	if err != nil {
		return err
	}

	// Generate the SVG file
	graphics := gb.buildGraphic()
	graphics.Viewport = options.Embedded
	graphics.DrawSVG(w)

	return nil
}

// Options for SVG image generation
type ImageOptions struct {
	// The diagram style
	Style *DiagramStyles

	// If true, generate attributes to make the SVG suitable for embedding
	// in other documents (e.g. HTML).
	Embedded bool
}

// The default options
var DefaultOptions = &ImageOptions{
	Style:    DefaultStyle,
	Embedded: false,
}

// A processing instruction
type ProcessingInstruction struct {
	Prefix string
	Value  string
}

// A participant
type Actor struct {
	Name  string
	Label string

	Icon      ActorIcon
	InHeader  bool
	InFooter  bool
	Lifeline  bool
	Color     string
	TextColor string

	rank int
}

// Special actors
var LeftOffsideActor *Actor = &Actor{rank: -1}
var RightOffsideActor *Actor = &Actor{rank: -2}

// The supported arrow stems
type ArrowStem int

const (
	SolidArrowStem  ArrowStem = iota
	DashedArrowStem           = iota
	ThickArrowStem            = iota
)

// The supported arrow heads
type ArrowHead int

const (
	SolidArrowHead     ArrowHead = iota
	OpenArrowHead                = iota
	BarbArrowHead                = iota
	LowerBarbArrowHead           = iota
)

// An arrow
type Arrow struct {
	Stem ArrowStem
	Head ArrowHead
}

// Note alignments
type NoteAlignment int

const (
	LeftNoteAlignment  NoteAlignment = iota
	RightNoteAlignment               = iota
	OverNoteAlignment                = iota
)

// A sequence item
type SequenceItem interface {
}

// Defines a note
type Note struct {
	// The note's alignment and position
	Actor1 *Actor
	Actor2 *Actor

	Align NoteAlignment

	// The message
	Message string
}

// Defines an action
type Action struct {
	// The originating actor
	From *Actor

	// The destination actor
	To *Actor

	// The arrow to use
	Arrow Arrow

	// The message
	Message string
}

type DividerType int

const (
	DTSpacer DividerType = iota
	DTGap
	DTFrame
	DTLine
)

// Defines a divider, which spans the diagram.
type Divider struct {
	// The message
	Message string

	// The divider type
	Type DividerType
}

// A framed block of sequence items.  Each block can have one or more segments,
// which will appear one after the other.
type Block struct {
	Segments []*BlockSegment
}

// Concurrent returns true if the block is a concurrent block segment.
func (b *Block) Concurrent() bool {
	return len(b.Segments) > 0 && b.Segments[0].Type == ConcurrentSegmentType
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

	// ParSegmentType is for the "par" blocks
	ParSegmentType

	// ParElseSegmentType is for the "parelse" blocks
	ParElseSegmentType

	// The opt segment
	OptSegmentType

	// The loop segment
	LoopSegmentType

	// ConcurrentSegmentType is for the first segment of a concurrent block.
	ConcurrentSegmentType

	// ConcurrentWhilstSegmentType is for the subsequent segments (i.e. the "whilst" segments)
	// of a concurrent block.
	ConcurrentWhilstSegmentType
)

// A segment within a block
type BlockSegment struct {
	Type     SegmentType
	Prefix   string
	Message  string
	SubItems []SequenceItem
}

// Returns the number of nested blocks
func (bs *BlockSegment) MaxNestDepth() int {
	nestDepth := 0
	for _, subItem := range bs.SubItems {
		if block, isBlock := subItem.(*Block); isBlock {
			nestDepth = maxInt(nestDepth, block.MaxNestDepth())
		}
	}
	return nestDepth
}
