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

type workingRemeasurements struct {
    size        int
    before      int
    after       int
}

// Remeasure the entire drawing.  Returns a rect containing the size of the image
func (g *Graphic) remeasure() Rect {

    cols, rows := g.Cols(), g.Rows()
    colWidths := make([]workingRemeasurements, g.Cols())
    rowHeights := make([]workingRemeasurements, g.Rows())

    // Resize the cells
    for _, item := range g.items {
        if (item.R >= 0) && (item.C >= 0) && (item.R < len(g.matrix)) && (item.C < len(g.matrix[item.R])) {
            if item2d, is2dItem := item.Item.(Graphbox2DItem) ; is2dItem {
                itemWidth, itemHeight := item2d.Size()
                rowHeights[item.R].size = maxInt(rowHeights[item.R].size, itemHeight)
                colWidths[item.C].size = maxInt(colWidths[item.C].size, itemWidth)
            }

            // If the item requires a margin.
            if marginItem, isMarginItem := item.Item.(MarginItem) ; isMarginItem {
                l, r, t, b := marginItem.Margin()
                rowHeights[item.R].before = maxInt(rowHeights[item.R].before, t)
                rowHeights[item.R].after = maxInt(rowHeights[item.R].after, b)
                colWidths[item.C].before = maxInt(colWidths[item.C].before, l)
                colWidths[item.C].after = maxInt(colWidths[item.C].after, r)
            }
        }
    }

    // Recalculate cell rectanges
    y := g.Margin.Y
    for r, row := range g.matrix {
        x := g.Margin.X
        largestH := 0
        for c, _ := range row {
            cr, rr := colWidths[c], rowHeights[r]
            innerRect := Rect {
                X: x,
                Y: y,
                W: cr.size,
                H: rr.size,
            }

            ox := x - maxInt(cr.before, g.Padding.X)
            oy := y - maxInt(rr.before, g.Padding.Y)
            outerRect := Rect {
                X: ox,
                Y: oy,
                W: maxInt(cr.before, g.Padding.X) + innerRect.W + maxInt(cr.after, g.Padding.X),
                H: maxInt(rr.before, g.Padding.Y) + innerRect.H + maxInt(rr.after, g.Padding.Y),
            }

            g.matrix[r][c].Frame = BoxFrame{outerRect, innerRect}
            x += outerRect.W
            largestH = maxInt(largestH, outerRect.H)
        }
        //y += rowHeights[r] + g.Padding.Y * 2
        y += largestH
    }    

    lastRect := g.matrix[rows - 1][cols - 1].Frame.OuterRect
    return Rect{
        W: lastRect.X + lastRect.W - g.Padding.X + g.Margin.X, 
        H: lastRect.Y + lastRect.H - g.Padding.Y + g.Margin.Y,
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