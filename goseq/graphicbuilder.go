package goseq

import (
    "./graphbox"
)


type DiagramStyles struct {
    ActorBox            graphbox.TextRectStyle
    NoteBox             graphbox.TextRectStyle
    ActivityLine        graphbox.ActivityLineStyle
}

var DefaultStyle DiagramStyles

func init() {
    font, err := graphbox.NewTTFFont("/usr/share/fonts/truetype/freefont/FreeSans.ttf")
    if err != nil { 
        panic(err)
    }

    DefaultStyle = DiagramStyles {
        ActorBox: graphbox.TextRectStyle {
            Font: font,
            FontSize: 16,
            Padding: graphbox.Point{16, 8},
        },
        NoteBox: graphbox.TextRectStyle {
            Font: font,
            FontSize: 14,
            Padding: graphbox.Point{8, 4},
        },
        ActivityLine: graphbox.ActivityLineStyle{
            Font:           font,
            FontSize:       14,
            PaddingTop:     4,
            PaddingBottom:  8,
            TextGap:        8,
        },
    }
}

type GraphicBuilder struct {
    Diagram           *Diagram
    Graphic           *graphbox.Graphic
    Font              graphbox.Font
    Style             DiagramStyles
}


func NewGraphicBuilder(d *Diagram) (*GraphicBuilder, error) {
    font, err := graphbox.NewTTFFont("/usr/share/fonts/truetype/freefont/FreeSans.ttf")
    if err != nil {
        return nil, err
    }

    return &GraphicBuilder{d, nil, font, DefaultStyle}, nil
}

func (gb *GraphicBuilder) BuildGraphic() *graphbox.Graphic {
    rows, cols := gb.calcRowsAndCols()
    gb.Graphic = graphbox.NewGraphic(rows, cols)

    gb.Graphic.Margin = graphbox.Point{16, 8}
    gb.Graphic.Padding = graphbox.Point{16, 8}
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
    col := gb.colOfActor(note.Actor)
    gb.Graphic.Put(row, col, graphbox.NewTextRect(note.Message, gb.Style.NoteBox))    
}

// Places an action
func (gb *GraphicBuilder) putAction(row int, action *Action) {
    fromCol := gb.colOfActor(action.From)
    toCol := gb.colOfActor(action.To)
    style := gb.Style.ActivityLine

    gb.Graphic.Put(row, fromCol, graphbox.NewActivityLine(toCol, action.Message, style))
}

// Count the number of rows needed in the graphic
func (gb *GraphicBuilder) calcRowsAndCols() (int, int) {
    // 1 for the title, object header and object footer
    return len(gb.Diagram.Items) + 2 + 1, len(gb.Diagram.Actors)
}

// Add the object headers and footers
func (gb *GraphicBuilder) addObjects() {
    // TODO: Proper styling
    bottomRow := gb.Graphic.Rows() - 1
    for rank, actor := range gb.Diagram.Actors {
        gb.Graphic.Put(1, rank, &graphbox.LifeLine{bottomRow, rank})

        gb.Graphic.Put(1, rank, graphbox.NewTextRect(actor.Name, gb.Style.ActorBox))
        gb.Graphic.Put(bottomRow, rank, graphbox.NewTextRect(actor.Name, gb.Style.ActorBox))
    }
}

// Returns the column position of an actor
func (gb *GraphicBuilder) colOfActor(actor *Actor) int {
    return actor.rank
}