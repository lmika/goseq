package goseq

import (
    "./graphbox"
)

type graphicBuilder struct {
    Diagram           *Diagram
}


func (gb *graphicBuilder) BuildGraphic() *graphbox.Graphic {
    rows, cols := gb.calcRowsAndCols()
    g := graphbox.NewGraphic(rows, cols)

    gb.addObjects(g)

    // TEMP
    for i, _ := range gb.Diagram.Items {
        g.Put(i + 2, 0, &graphbox.ActorRect{5, 50, "|"})
    }

    return g
}

// Count the number of rows needed in the graphic
func (gb *graphicBuilder) calcRowsAndCols() (int, int) {
    // 1 for the title, object header and object footer
    return len(gb.Diagram.Items) + 2 + 1, len(gb.Diagram.Actors)
}

// Add the object headers and footers
func (gb *graphicBuilder) addObjects(g *graphbox.Graphic) {
    // TODO: Proper styling
    bottomRow := g.Rows() - 1
    for rank, actor := range gb.Diagram.Actors {
        g.Put(1, rank, &graphbox.LifeLine{bottomRow, rank})

        g.Put(1, rank, &graphbox.ActorRect{100, 35, actor.Name})
        g.Put(bottomRow, rank, &graphbox.ActorRect{100, 35, actor.Name})
    }
}