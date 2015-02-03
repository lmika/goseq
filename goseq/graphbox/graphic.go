// The main graphbox type
//

package graphbox

import (
    "io"

    "github.com/ajstarks/svgo"
)


// A graphbox diagram.  This diagram is made up of uniform points.
// Each item on the diagram has the option of describing a constraint such
// as the amount of spacing around the point it requires.
type Graphic struct {
    matrix      [][]matrixItem
    items       []itemInstance

    // The margin between items
    Margin      Point
}

func NewGraphic(rows, cols int) *Graphic {
    g := &Graphic{}
    g.resizeTo(rows + 1, cols + 1)
    return g
}

// Returns the number of items within the matrix
func (g *Graphic) Rows() int {
    return len(g.matrix) - 1
}

// Returns the number of columns within the matrix
func (g *Graphic) Cols() int {
    if (len(g.matrix) > 0) {
        return len(g.matrix[0]) - 1
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
func (g *Graphic) remeasure() (int, int) {
    // Reinitializes the matrix
    g.reinitMatrix()

    // Run through the constraints
    for _, item := range g.items {
        constraint := item.Item.Constraint(item.R, item.C)
        if constraint != nil {
            constraint.Apply(g)
        }
    }

    // Reposition the grid points
    return g.repositionGridPoints()
}

func (g *Graphic) reinitMatrix() {
    for r, row := range g.matrix {
        for c, _ := range row {
            g.matrix[r][c] = matrixItem{}
        }
    }
}

func (g *Graphic) repositionGridPoints() (int, int) {
    maxX := 0

    py := g.Margin.Y
    for r, row := range g.matrix {
        px := g.Margin.X
        py += g.matrix[r][0].Delta.Y

        for c := range row {
            px += g.matrix[r][c].Delta.X
            g.matrix[r][c].Point.X = px
            g.matrix[r][c].Point.Y = py
        }
        //px += g.matrix[r][len(row) - 1].Delta.X + g.Margin.X
        px += g.Margin.X
        maxX = maxInt(px, maxX)
    }
    //py += g.matrix[len(g.matrix) - 1][0].Delta.Y + g.Margin.Y
    py += g.Margin.Y

    return maxX, py
}


// Implementation of ConstrantModifier
//
// TODO: Instead of using loops, add cheats like only using
// the first row/column
func (g *Graphic) GridPointRect(fr, fc, tr, tc int) (int, int) {
    w, h := 0, 0
    for r := fr + 1; r < tr; r++ {
        for c := fc + 1; c < tc; c++ {
            w += g.matrix[r][c].Delta.X
            h += g.matrix[r][c].Delta.Y
        }
    }

    return w, h
}

func (g *Graphic) EnsureLeftIsAtleast(col, newLeft int) {
    for r := 0; r < g.Rows(); r++ {
        g.matrix[r][col].Delta.X = maxInt(g.matrix[r][col].Delta.X, newLeft)
    }
}

func (g *Graphic) EnsureTopIsAtLeast(row, newTop int) {
    for c := 0; c < g.Cols(); c++ {
        g.matrix[row][c].Delta.Y = maxInt(g.matrix[row][c].Delta.Y, newTop)
    }
}

func (g *Graphic) AddLeftToCol(col, newLeft int) {
    for r := 0; r < g.Rows(); r++ {
        g.matrix[r][col].Delta.X += newLeft
    }
}

func (g *Graphic) AddTopToRow(row, newTop int) {
    for c := 0; c < g.Cols(); c++ {
        g.matrix[row][c].Delta.Y += newTop
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
    sizeW, sizeH := g.remeasure()

    canvas := svg.New(w)
    canvas.Start(sizeW, sizeH)
    defer canvas.End()

    for _, item := range g.items {
        g.drawItem(canvas, item)
    }

    // DEBUG: Draw the grid
    for _, row := range g.matrix {
        for _, cell := range row {
            canvas.Circle(cell.Point.X, cell.Point.Y, 2, "brush:blue")
        }
    }
}

// Draws the item
func (g *Graphic) drawItem(canvas *svg.SVG, item itemInstance) {
    if !((item.R >= 0) && (item.C >= 0) && (item.R < len(g.matrix)) && (item.C < len(g.matrix[item.R]))) {
        // Do nothing
        return
    }

    ctx := DrawContext{canvas, g, item.R, item.C}
    point := g.matrix[item.R][item.C].Point
    item.Item.Draw(ctx, point)
}

func (g *Graphic) PointAt(r, c int) (Point, bool) {
    if (r >= 0) && (c >= 0) && (r < len(g.matrix)) && (c < len(g.matrix[r])) {
        return g.matrix[r][c].Point, true
    } else {
        return Point{}, false
    }
}


// A matrix cell item
type matrixItem struct {
    // The location of the item.
    Point       Point

    // Deltas between this point and the point to the left (X) and above (Y)
    Delta       Point
}

type itemInstance struct {
    R, C        int
    Item        GraphboxItem
}