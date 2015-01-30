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

    // The margin between items
    Margin      Point
}

func NewGraphic(rows, cols int) *Graphic {
    g := &Graphic{}
    g.resizeTo(rows, cols)
    return g
}

// Returns the number of items within the matrix
func (g *Graphic) rows() int {
    return len(g.matrix)
}

// Returns the number of columns within the matrix
func (g *Graphic) cols() int {
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

// Calculate the maximum height of elements within the particular row
func (g *Graphic) maxHeight(r int) int {
    maxHeight := 0
    if (r >= 0) && (r < len(g.matrix)) {
        for _, item := range g.matrix[r] {
            if (item.Item != nil) {
                _, itemHeight := item.Item.Size()
                maxHeight = maxInt(maxHeight, itemHeight)
            }
        }
    }
    return maxHeight
}

// Calculate the maximum width of elements within the particular column
func (g *Graphic) maxWidth(c int) int {
    maxWidth := 0
    for _, row := range g.matrix {
        item := row[c]
        if (c >= 0) && (c < len(row)) {
            if (item.Item != nil) {
                itemWidth, _ := item.Item.Size()
                maxWidth = maxInt(maxWidth, itemWidth)
            }
        }
    }
    return maxWidth
}

// Remeasure the entire drawing.  Returns a rect containing the size of the image
func (g *Graphic) remeasure() Rect {
    // TODO: Margin and padding
    y := 0
    maxX := 0       // Maximum right most point
    for r, row := range g.matrix {
        x := 0
        maxHeight := g.maxHeight(r)

        // TODO: Padding
        cellHeight := maxHeight

        for c, _ := range row {
            // TODO: Caching
            maxWidth := g.maxWidth(c)

            // TODO: Inner item padding
            cellWidth := maxWidth

            g.matrix[r][c].OuterRect = Rect{x, y, cellWidth, cellHeight}
            x += cellWidth
        }

        maxX = maxInt(maxX, x)
        y += cellHeight
    }

    return Rect{0, 0, maxX, y}
}

// Sets a point in the matrix.  If the point is beyond the scope of the matrix,
// returns false.
func (g *Graphic) Put(r, c int, item GraphboxItem) bool {
    if (r >= 0) && (c >= 0) && (r < len(g.matrix)) && (c < len(g.matrix[r])) {
        g.matrix[r][c].Item = item
        //g.matrix[r][c].OuterRect = Rect{c * 150, r * 100, 150, 100}     // TEMP
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

    ctx := &DrawContext{canvas}
    for r, rs := range g.matrix {
        for c := range rs {
            g.drawItem(ctx, r, c)
        }
    }
}

// Draws the item
func (g *Graphic) drawItem(ctx *DrawContext, r, c int) {
    if !((r >= 0) && (c >= 0) && (r < len(g.matrix)) && (c < len(g.matrix[r]))) {
        // Do nothing
        return
    }

    item := g.matrix[r][c]
    if (item.Item == nil) {
        return
    }

    frame := BoxFrame{item.OuterRect, item.OuterRect}
    item.Item.Draw(ctx, frame)
}


// A matrix cell item
type matrixItem struct {
    Item        GraphboxItem
    OuterRect   Rect
}