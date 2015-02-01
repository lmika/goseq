// The main graphbox type
//

package graphbox

import (
    "io"

    "github.com/ajstarks/svgo"
)


// A graphbox diagram.
type Graphic struct {
    matrix      [][]matrixItem
    items       []itemInstance

    // The margin between items
    Margin      Point
    Padding     Point
}

func NewGraphic(rows, cols int) *Graphic {
    g := &Graphic{}
    g.resizeTo(rows, cols)
    return g
}

// Returns the number of items within the matrix
func (g *Graphic) Rows() int {
    return len(g.matrix)
}

// Returns the number of columns within the matrix
func (g *Graphic) Cols() int {
    if (len(g.matrix) > 0) {
        return len(g.matrix[0])
    } else {
        return 0
    }
}

// Resize the matrix
func (g *Graphic) resizeTo(rows, cols int) {
    newRows := make([][]matrixItem, rows)
    copy(newRows, g.matrix)

    for i := range newRows {
        newCols := make([]matrixItem, cols)
        if i < len(g.matrix) {
            copy(newCols, g.matrix[i])
        }
        newRows[i] = newCols
    }

    g.matrix = newRows
}

// Remeasure the entire drawing.  Returns a rect containing the size of the image
// TODO: Chean this UP!!!
func (g *Graphic) remeasure() Rect {

    cols, rows := g.Cols(), g.Rows()
    colWidths := make([]int, g.Cols())
    rowHeights := make([]int, g.Rows())

    // Gaps between rows and columns (to the left or or top of)
    colGaps := make([]int, g.Cols() + 1)
    rowGaps := make([]int, g.Rows() + 1)

    for r := 1; r < len(rowGaps); r++ {
        rowGaps[r] = g.Padding.Y
    }
    for c := 1; c < len(colGaps); c++ {
        colGaps[c] = g.Padding.X
    }

    // Resize the cells
    for _, item := range g.items {
        if (item.R >= 0) && (item.C >= 0) && (item.R < len(g.matrix)) && (item.C < len(g.matrix[item.R])) {
            if item2d, is2dItem := item.Item.(Graphbox2DItem) ; is2dItem {
                itemWidth, itemHeight := item2d.Size()
                rowHeights[item.R] = maxInt(rowHeights[item.R], itemHeight)
                colWidths[item.C] = maxInt(colWidths[item.C], itemWidth)
            }
        }
    }

    // Determine padding.  The total gap between two colums is the total amount of gap requested by
    // the two columns PLUS the widths of each of the two columns MINUS the column widths.
    //
    // This is to make it possible to interlink overlaps
    //
    // EG:
    //
    //          [  |      ]   [ |  ]
    //            [| ] [        | ]
    for _, item := range g.items {
        if (item.R >= 0) && (item.C >= 0) && (item.R < len(g.matrix)) && (item.C < len(g.matrix[item.R])) {
            requestW, requestH := 0, 0
            if item2d, is2dItem := item.Item.(Graphbox2DItem) ; is2dItem {
                requestW, requestH = item2d.Size()
            }

            // If the item requires a margin.
            if marginItem, isMarginItem := item.Item.(MarginItem) ; isMarginItem {
                l, r, t, b := marginItem.Margin()

                rowGaps[item.R] = maxInt(rowGaps[item.R], t + requestH / 2 - rowHeights[item.R] / 2)
                colGaps[item.C] = maxInt(colGaps[item.C], l + requestW / 2 - colWidths[item.C] / 2)

                if (item.R < len(g.matrix) - 1) {
                    rowGaps[item.R + 1] = maxInt(rowGaps[item.R + 1], b + requestH / 2 - rowHeights[item.R + 1] / 2)
                } else {
                    rowGaps[item.R + 1] = maxInt(rowGaps[item.R + 1], b + requestH / 2)
                }
                if (item.C < len(g.matrix[item.R]) - 1) {
                    colGaps[item.C + 1] = maxInt(colGaps[item.C + 1], r + requestW / 2 - colWidths[item.C + 1] / 2)
                } else {
                    colGaps[item.C + 1] = maxInt(colGaps[item.C + 1], r + requestW / 2)
                }
            }
        }
    }

    // Recalculate cell rectanges
    y := g.Margin.Y - rowGaps[0]
    for r, row := range g.matrix {
        x := g.Margin.X
        largestH := 0
        for c, _ := range row {
            cr, rr := colWidths[c], rowHeights[r]
            innerRect := Rect {
                X: colGaps[c] + x,
                Y: rowGaps[r] + y,
                W: cr,
                H: rr,
            }

            ox := x
            oy := y
            outerRect := Rect {
                X: ox,
                Y: oy,
                W: colGaps[c] + cr + colGaps[c + 1],
                H: rowGaps[r] + rr + rowGaps[r + 1],
            }

            g.matrix[r][c].Frame = BoxFrame{outerRect, innerRect}

            x += colGaps[c] + cr
            largestH = maxInt(largestH, outerRect.H)
        }
        //y += rowHeights[r] + g.Padding.Y * 2
        y += largestH
    }    

    lastRect := g.matrix[rows - 1][cols - 1].Frame.OuterRect
    return Rect{
        W: lastRect.X + lastRect.W + g.Margin.X,
        H: lastRect.Y + lastRect.H + g.Margin.Y,
    }
}

// Sets a point in the matrix.  If the point is beyond the scope of the matrix,
// returns false.
func (g *Graphic) Put(r, c int, item GraphboxItem) bool {
    if (r >= 0) && (c >= 0) && (r < len(g.matrix)) && (c < len(g.matrix[r])) {
        //g.matrix[r][c].Item = item
        g.items = append(g.items, itemInstance{r, c, item})
        return true
    } else {
        return false
    }    
}

// Draws the graphics as an SVG
func (g *Graphic) DrawSVG(w io.Writer) {
    size := g.remeasure()

    canvas := svg.New(w)
    canvas.Start(size.W, size.H)
    defer canvas.End()

    for _, item := range g.items {
        g.drawItem(canvas, item)
    }
}

// Draws the item
func (g *Graphic) drawItem(canvas *svg.SVG, item itemInstance) {
    if !((item.R >= 0) && (item.C >= 0) && (item.R < len(g.matrix)) && (item.C < len(g.matrix[item.R]))) {
        // Do nothing
        return
    }

    ctx := DrawContext{canvas, g, item.R, item.C}
    frame := g.matrix[item.R][item.C].Frame
    item.Item.Draw(ctx, frame)
}

// Gets the outer rectangle of a particular cell
func (g *Graphic) frameAtCell(r, c int) (BoxFrame, bool) {
    if (r >= 0) && (c >= 0) && (r < len(g.matrix)) && (c < len(g.matrix[r])) {
        return g.matrix[r][c].Frame, true
    } else {
        return BoxFrame{}, false
    }
}


// A matrix cell item
type matrixItem struct {
    Frame       BoxFrame
}

type itemInstance struct {
    R, C        int
    Item        GraphboxItem
}