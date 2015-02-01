package goseq

import (
    "./graphbox"
)

type GraphicBuilder struct {
    Diagram           *Diagram
    Font              graphbox.Font
}


func NewGraphicBuilder(d *Diagram) (*GraphicBuilder, error) {
    font, err := graphbox.NewTTFFont("/usr/share/fonts/truetype/freefont/FreeSans.ttf")
    if err != nil {
        return nil, err
    }

    return &GraphicBuilder{d, font}, nil
}

func (gb *GraphicBuilder) BuildGraphic() *graphbox.Graphic {
    rows, cols := gb.calcRowsAndCols()
    g := graphbox.NewGraphic(rows, cols)

    g.Margin = graphbox.Point{16, 8}
    g.Padding = graphbox.Point{16, 8}
    gb.addObjects(g)

    // TEMP
    for i, item := range gb.Diagram.Items {
        row := i + 2
        switch itemDetails := item.(type) {
        case *Action:
            g.Put(row, gb.colOfActor(itemDetails.From), 
                &graphbox.ActivityLine{gb.colOfActor(itemDetails.To), itemDetails.Message})
        }
    }

    return g
}

// Count the number of rows needed in the graphic
func (gb *GraphicBuilder) calcRowsAndCols() (int, int) {
    // 1 for the title, object header and object footer
    return len(gb.Diagram.Items) + 2 + 1, len(gb.Diagram.Actors)
}

// Add the object headers and footers
func (gb *GraphicBuilder) addObjects(g *graphbox.Graphic) {
    // TODO: Proper styling
    bottomRow := g.Rows() - 1
    for rank, actor := range gb.Diagram.Actors {
        g.Put(1, rank, &graphbox.LifeLine{bottomRow, rank})

        g.Put(1, rank, graphbox.NewActorRect(actor.Name, gb.Font))
        g.Put(bottomRow, rank, graphbox.NewActorRect(actor.Name, gb.Font))
    }
}

// Returns the column position of an actor
func (gb *GraphicBuilder) colOfActor(actor *Actor) int {
    return actor.rank
}