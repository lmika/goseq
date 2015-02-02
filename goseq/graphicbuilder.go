package goseq

import (
    "./graphbox"
)


// Various position offsets
const (
    ObjectLeftX     =   1
    TitleY          =   0
    ObjectY         =   1
)

// This arrow head to arrow head
var ArrowHeadMapping = map[ArrowHead]graphbox.ActivityArrowHead {
    SolidArrowHead: graphbox.SolidArrowHead,
    OpenArrowHead: graphbox.OpenArrowHead,
}

var ArrowStemMapping = map[ArrowStem]graphbox.ActivityArrowStem {
    SolidArrowStem: graphbox.SolidArrowStem,
    DashedArrowStem: graphbox.DashedArrowStem,
}

type DiagramStyles struct {
    ActorBox            graphbox.ActorBoxStyle
    NoteBox             graphbox.NoteBoxStyle
    ActivityLine        graphbox.ActivityLineStyle
    Title               graphbox.TitleStyle
}

var DefaultStyle DiagramStyles

func init() {
    font, err := graphbox.NewTTFFont("/usr/share/fonts/truetype/freefont/FreeSans.ttf")
    if err != nil { 
        panic(err)
    }

    DefaultStyle = DiagramStyles {
        ActorBox: graphbox.ActorBoxStyle {
            Font: font,
            FontSize: 16,
            Padding: graphbox.Point{16, 8},
            Margin: graphbox.Point{8, 8},
        },
        NoteBox: graphbox.NoteBoxStyle {
            Font: font,
            FontSize: 14,
            Padding: graphbox.Point{8, 4},
            Margin: graphbox.Point{8, 8},
        },
        ActivityLine: graphbox.ActivityLineStyle{
            Font: font,
            FontSize: 14,
            Margin: graphbox.Point{16, 8},
            TextGap: 4,
        },
        Title: graphbox.TitleStyle {
            Font: font,
            FontSize: 20,
            Padding: graphbox.Point{16, 8},
        },        
    }
}


// Information about a particular actor
type actorInfo struct {
    // Extra cols needed on the left or right
    ExtraLeftCol    bool
    ExtraRightCol   bool

    // Actor column
    Col             int
}


type GraphicBuilder struct {
    Diagram           *Diagram
    Graphic           *graphbox.Graphic
    Font              graphbox.Font
    Style             DiagramStyles

    actorInfos        []actorInfo
}


func NewGraphicBuilder(d *Diagram) (*GraphicBuilder, error) {
    font, err := graphbox.NewTTFFont("/usr/share/fonts/truetype/freefont/FreeSans.ttf")
    if err != nil {
        return nil, err
    }

    return &GraphicBuilder{d, nil, font, DefaultStyle, nil}, nil
}

func (gb *GraphicBuilder) BuildGraphic() *graphbox.Graphic {
    rows, cols := gb.calcRowsAndCols()
    gb.Graphic = graphbox.NewGraphic(rows, cols)

    gb.Graphic.Margin = graphbox.Point{16, 8}

    // Add a title
    if gb.Diagram.Title != "" {
        gb.Graphic.Put(0, 0, graphbox.NewTitle(cols, gb.Diagram.Title, gb.Style.Title))
    }

    gb.addObjects()

    // TEMP
    for i, item := range gb.Diagram.Items {
        row := i + 2
        switch itemDetails := item.(type) {
        case *Action:
            gb.putAction(row, itemDetails)
        case *Note:
            gb.putNote(row, itemDetails)
        }
    }

    return gb.Graphic
}

// Places a note
func (gb *GraphicBuilder) putNote(row int, note *Note) {
    var pos graphbox.NoteBoxPos

    if note.Align == LeftNoteAlignment {
        pos = graphbox.LeftNotePos
    } else if note.Align == OverNoteAlignment {
        pos = graphbox.CenterNotePos
    } else if note.Align == RightNoteAlignment {
        pos = graphbox.RightNotePos
    }

    col := gb.colOfActor(note.Actor)
    gb.Graphic.Put(row, col, graphbox.NewNoteBox(note.Message, gb.Style.NoteBox, pos))    
}

// Places an action
func (gb *GraphicBuilder) putAction(row int, action *Action) {
    fromCol := gb.colOfActor(action.From)
    toCol := gb.colOfActor(action.To)
    style := gb.Style.ActivityLine

    style.ArrowHead = ArrowHeadMapping[action.Arrow.Head]
    style.ArrowStem = ArrowStemMapping[action.Arrow.Stem]

    gb.Graphic.Put(row, fromCol, graphbox.NewActivityLine(toCol, action.Message, style))
}

// Count the number of rows needed in the graphic
func (gb *GraphicBuilder) calcRowsAndCols() (int, int) {
    cols := gb.determineActorInfo()

    // 1 for the title, object header and object footer
    return len(gb.Diagram.Items) + ObjectY + 2, cols
}

// Determine actor information.  Returns the number of colums required
func (gb *GraphicBuilder) determineActorInfo() int {
    gb.actorInfos = make([]actorInfo, len(gb.Diagram.Actors))

    // Allocate the columns
    cols := ObjectLeftX
    for _, actor := range gb.Diagram.Actors {
        colsRequiredByActor := 1
        actorCol := cols

        if (gb.actorInfos[actor.rank].ExtraLeftCol) {
            colsRequiredByActor++
            actorCol++
        }
        if (gb.actorInfos[actor.rank].ExtraRightCol) {
            colsRequiredByActor++
        }

        gb.actorInfos[actor.rank].Col = actorCol
        cols += colsRequiredByActor
    }

    return cols
}

// Add the object headers and footers
func (gb *GraphicBuilder) addObjects() {
    // TODO: Proper styling
    bottomRow := gb.Graphic.Rows() - 1
    for _, actor := range gb.Diagram.Actors {
        col := gb.colOfActor(actor)
        gb.Graphic.Put(ObjectY, col, &graphbox.LifeLine{bottomRow, col})

        gb.Graphic.Put(ObjectY, col, graphbox.NewActorBox(actor.Name, gb.Style.ActorBox, graphbox.TopActorBox))
        gb.Graphic.Put(bottomRow, col, graphbox.NewActorBox(actor.Name, gb.Style.ActorBox, graphbox.BottomActorBox))
    }
}

// Returns the column position of an actor
func (gb *GraphicBuilder) colOfActor(actor *Actor) int {
    return gb.actorInfos[actor.rank].Col
}