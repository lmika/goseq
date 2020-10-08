package seqdiagram

import (
	"errors"
	"sort"

	"github.com/lmika/goseq/seqdiagram/graphbox"
)

// Various position offsets
const (
	posObjectLeftX = 1
	posObjectY     = 1
)

var graphboxArrowStemMapping = map[ArrowStem]graphbox.ActivityArrowStem{
	SolidArrowStem:  graphbox.SolidArrowStem,
	DashedArrowStem: graphbox.DashedArrowStem,
	ThickArrowStem:  graphbox.ThickArrowStem,
}

// Load the internal font
func mustLoadFont() *graphbox.TTFFont {
	font, err := loadInternalFont(dejaVuSansFont)
	if err != nil {
		panic(errors.New("Could not load internal font: " + dejaVuSansFont))
	}

	return font
}

// Information about a particular actor
type actorInfo struct {
	// Extra cols needed on the left or right
	ExtraLeftCol  bool
	ExtraRightCol bool

	// Actor column
	Col int
}

type graphicBuilder struct {
	Diagram *Diagram
	Graphic *graphbox.Graphic
	Style   *DiagramStyles

	actorInfos []actorInfo
}

func newGraphicBuilder(d *Diagram, style *DiagramStyles) (*graphicBuilder, error) {
	return &graphicBuilder{d, nil, style, nil}, nil
}

func (gb *graphicBuilder) buildGraphic() *graphbox.Graphic {
	rows, cols := gb.calcRowsAndCols()
	gb.Graphic = graphbox.NewGraphic(rows, cols)

	gb.Graphic.Margin = gb.Style.Margin
	gb.Graphic.ShowGrid = false

	gb.addActors()

	if len(gb.Diagram.Items) == 0 {
		gb.Graphic.Put(2, 0, &graphbox.Spacer{graphbox.Point{0, 64}})
	} else {
		row := 2
		gb.putItemsInSlice(&row, 0, gb.Diagram.Items)
	}

	// Add a title
	if gb.Diagram.Title != "" {
		gb.Graphic.Put(0, 0, graphbox.NewTitle(cols, gb.Diagram.Title, gb.Style.Title))
	}

	return gb.Graphic
}

// Place items in a slice.  This will update the rows pointer
func (gb *graphicBuilder) putItemsInSlice(row *int, depth int, items []SequenceItem) {
	for _, item := range items {
		switch itemDetails := item.(type) {
		case *Action:
			gb.putAction(*row, itemDetails)
		case *Note:
			gb.putNote(*row, itemDetails)
		case *Divider:
			gb.putDivider(*row, itemDetails)
		case *Block:
			gb.putBlock(row, depth, itemDetails)
		}

		*row += 1
	}
}

// Calculate rows in slice
func (gb *graphicBuilder) calcItemsInSlice(items []SequenceItem) int {
	rows := 0
	for _, item := range items {
		switch itemDetails := item.(type) {
		case *Block:
			if itemDetails.Concurrent() {
				maxRows := 0
				for _, seg := range itemDetails.Segments {
					rowsInSeg := gb.calcItemsInSlice(seg.SubItems) + 1
					if rowsInSeg > maxRows {
						maxRows = rowsInSeg
					}
				}
				rows += maxRows
			} else {
				for _, seg := range itemDetails.Segments {
					rows += gb.calcItemsInSlice(seg.SubItems) + 1
				}
			}
			rows += 1
		default:
			rows++
		}
	}
	return rows
}

// Places a note
func (gb *graphicBuilder) putNote(row int, note *Note) {
	if (note.Actor2 == nil) || (note.Actor1 == note.Actor2) {
		gb.putSingleActorNote(row, note.Actor1, note)
	} else {
		var leftActor, rightActor *Actor
		if gb.colOfActor(note.Actor1) < gb.colOfActor(note.Actor2) {
			leftActor, rightActor = note.Actor1, note.Actor2
		} else {
			leftActor, rightActor = note.Actor2, note.Actor1
		}

		switch note.Align {
		case OverNoteAlignment:
			gb.putMultiActorOverNote(row, leftActor, rightActor, note)
		case LeftNoteAlignment:
			gb.putSingleActorNote(row, leftActor, note)
		case RightNoteAlignment:
			gb.putSingleActorNote(row, rightActor, note)
		}
	}
}

// Places a note over a single actor
func (gb *graphicBuilder) putSingleActorNote(row int, actor *Actor, note *Note) {
	var pos graphbox.NoteBoxPos

	if note.Align == LeftNoteAlignment {
		pos = graphbox.LeftNotePos
	} else if note.Align == OverNoteAlignment {
		pos = graphbox.CenterNotePos
	} else if note.Align == RightNoteAlignment {
		pos = graphbox.RightNotePos
	}

	col := gb.colOfActor(actor)
	gb.Graphic.Put(row, col, graphbox.NewNoteBox(note.Message, gb.Style.NoteBox, pos))
}

// Places a note over a multiple actors.  This actually uses the divider graphics object
// with the style adopted from the note style
func (gb *graphicBuilder) putMultiActorOverNote(row int, leftActor *Actor, rightActor *Actor, note *Note) {
	dividerBox := graphbox.DividerStyle{
		Font:        gb.Style.NoteBox.Font,
		FontSize:    gb.Style.NoteBox.FontSize,
		Padding:     gb.Style.NoteBox.Padding,
		Margin:      gb.Style.NoteBox.Margin,
		TextPadding: graphbox.Point{0, 0},
		Shape:       graphbox.DSFramedRect,
		Overlap:     gb.Style.MultiNoteOverlap,
	}

	fromCol := gb.colOfActor(leftActor)
	toCol := gb.colOfActor(rightActor)

	// TODO: This was just to avoid bad styling of notes which reference 'left' and 'right'
	// This needs to be fixed using proper styling, instead of this hack.
	if fromCol == 0 {
		fromCol = 1
	}
	if toCol == gb.Graphic.Cols()-1 {
		toCol = gb.Graphic.Cols() - 2
	}

	gb.Graphic.Put(row, fromCol, graphbox.NewDivider(toCol, note.Message, dividerBox))
}

// Places an action
func (gb *graphicBuilder) putAction(row int, action *Action) {
	fromCol := gb.colOfActor(action.From)
	toCol := gb.colOfActor(action.To)

	style := gb.Style.ActivityLine

	style.ArrowHead = gb.Style.ArrowHeads[action.Arrow.Head] //graphboxArrowHeadMapping[action.Arrow.Head]
	style.ArrowStem = graphboxArrowStemMapping[action.Arrow.Stem]

	gb.Graphic.Put(row, fromCol, graphbox.NewActivityLine(toCol, fromCol == toCol, action.Message, style))
}

// Places a divider
func (gb *graphicBuilder) putDivider(row int, action *Divider) {
	fromCol := 0
	toCol := gb.Graphic.Cols() - 1
	style := gb.Style.Divider[action.Type]

	gb.Graphic.Put(row, fromCol, graphbox.NewDivider(toCol, action.Message, style))
}

// Places a block
func (gb *graphicBuilder) putBlock(row *int, depth int, action *Block) {
	if action.Concurrent() {
		gb.putBlockSegmentsConcurrently(row, depth, action)
	} else {
		gb.putBlockSegmentsSequentially(row, depth, action)
	}
}

func (gb *graphicBuilder) putBlockSegmentsConcurrently(row *int, depth int, action *Block) {
	startRow := *row

	for _, seg := range action.Segments {
		thisRow := startRow
		gb.putItemsInSlice(&thisRow, depth+1, seg.SubItems)
		if thisRow > *row {
			*row = thisRow
		}
	}
	*row++
}

func getInnerRanksRecursive(subItems []SequenceItem) []int {
	ranks := []int{}
	for _, subItem := range subItems {
		if action, isAction := subItem.(*Action); isAction {
			ranks = append(ranks, action.From.rank)
			ranks = append(ranks, action.To.rank)
		} else if block, isBlock := subItem.(*Block); isBlock {
			for _, segment := range block.Segments {
				ranks = append(ranks, getInnerRanksRecursive(segment.SubItems)...)
			}
		}
	}
	return ranks
}

func (gb *graphicBuilder) getStartAndEndColsBasedOnContent(subItems []SequenceItem) (int, int) {
	startCol := 0
	endCol := gb.Graphic.Cols() - 1 // This needs to be the column of the last actor

	innerRanks := getInnerRanksRecursive(subItems)
	sort.Ints(innerRanks)
	// log.Println(innerRanks)
	if len(innerRanks) > 0 {
		// +1 because actor rank and column values are offset by one
		startCol = innerRanks[0] + 1
		endCol = innerRanks[len(innerRanks)-1] + 1
	}

	return startCol, endCol
}

func (gb *graphicBuilder) putBlockSegmentsSequentially(row *int, depth int, action *Block) {
	style := gb.Style.Block

	var startRow, endRow int
	startRow = *row
	nestDepth := action.MaxNestDepth()

	for i, seg := range action.Segments {
		// To outline only the inner actors of the block we need to set...
		//  - startCol to the leftmost inner actor of the block
		//  - endCol to the rightmost inner actor of the block
		startCol, endCol := gb.getStartAndEndColsBasedOnContent(seg.SubItems)

		*row++
		gb.putItemsInSlice(row, depth+1, seg.SubItems)
		endRow = *row

		segPrefix := ""
		showPrefix := true

		switch seg.Type {
		case AltSegmentType:
			segPrefix = "alt"
		case ElseSegmentType:
			segPrefix = "alt"
			showPrefix = false
		case ParSegmentType:
			segPrefix = "par"
		case ParElseSegmentType:
			segPrefix = "par"
			showPrefix = false
		case OptSegmentType:
			segPrefix = "opt"
		case LoopSegmentType:
			segPrefix = "loop"
		}

		if seg.Prefix != "" {
			segPrefix = seg.Prefix
		}

		block := graphbox.NewBlock(endRow, endCol, nestDepth, i == len(action.Segments)-1,
			segPrefix, showPrefix, seg.Message, style)
		gb.Graphic.Put(startRow, startCol, block)

		startRow = endRow
	}
}

// Count the number of rows needed in the graphic
func (gb *graphicBuilder) calcRowsAndCols() (int, int) {
	cols := gb.determineActorInfo() + 1

	// 1 for the title, object header and object footer
	if len(gb.Diagram.Items) == 0 {
		return posObjectY + 3, cols
	} else {
		return gb.calcItemsInSlice(gb.Diagram.Items) + posObjectY + 2, cols
	}
}

// Determine actor information.  Returns the number of colums required
func (gb *graphicBuilder) determineActorInfo() int {
	gb.actorInfos = make([]actorInfo, len(gb.Diagram.Actors))

	// Allocate the columns
	cols := posObjectLeftX
	for _, actor := range gb.Diagram.Actors {
		colsRequiredByActor := 1
		actorCol := cols

		if actor.rank != -1 {
			if gb.actorInfos[actor.rank].ExtraLeftCol {
				colsRequiredByActor++
				actorCol++
			}
			if gb.actorInfos[actor.rank].ExtraRightCol {
				colsRequiredByActor++
			}

			gb.actorInfos[actor.rank].Col = actorCol
			cols += colsRequiredByActor
		}
	}

	return cols
}

// Add the object headers and footers
func (gb *graphicBuilder) addActors() {
	// TODO: Proper styling
	bottomRow := gb.Graphic.Rows() - 1
	for rank, actor := range gb.Diagram.Actors {
		var actorBoxPos graphbox.ActorBoxPos

		if rank == 0 {
			actorBoxPos = graphbox.LeftActorBox
		} else if rank == len(gb.Diagram.Actors)-1 {
			actorBoxPos = graphbox.RightActorBox
		} else {
			actorBoxPos = graphbox.MiddleActorBox
		}

		col := gb.colOfActor(actor)

		if actor.Lifeline {
			gb.Graphic.Put(posObjectY, col, &graphbox.LifeLine{
				TR: bottomRow,
				TC: col,
				Style: graphbox.LifeLineStyle{
					Color: actor.Color,
				},
			})
		}

		if actor.Icon != nil {
			actorIconStyle := gb.Style.ActorIconBox
			actorIconStyle.Color = actor.Color
			actorIconStyle.TextColor = actor.TextColor

			if actor.InHeader {
				gb.Graphic.Put(posObjectY, col, graphbox.NewActorIconBox(actor.Label, actor.Icon.graphboxIcon(), actorIconStyle, actorBoxPos|graphbox.TopActorBox))
			}
		} else {
			// Configure the style
			actorStyle := gb.Style.ActorBox
			actorStyle.Color = actor.Color
			actorStyle.TextColor = actor.TextColor

			if actor.InHeader {
				gb.Graphic.Put(posObjectY, col, graphbox.NewActorBox(actor.Label, actorStyle, actorBoxPos|graphbox.TopActorBox))
				if actor.InFooter {
					gb.Graphic.Put(bottomRow, col, graphbox.NewActorBox(actor.Label, actorStyle, actorBoxPos|graphbox.BottomActorBox))
				}
			} else {
				if actor.InFooter {
					// Use the TopActorBox as that performs the layout
					gb.Graphic.Put(bottomRow, col, graphbox.NewActorBox(actor.Label, actorStyle, actorBoxPos|graphbox.TopActorBox))
				}
			}
		}
	}
}

// Returns the column position of an actor
func (gb *graphicBuilder) colOfActor(actor *Actor) int {
	if actor == LeftOffsideActor {
		return 0
	} else if actor == RightOffsideActor {
		return gb.Graphic.Cols() - 1
	} else {
		return gb.actorInfos[actor.rank].Col
	}
}
