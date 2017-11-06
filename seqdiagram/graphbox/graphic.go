// The main graphbox type
//

package graphbox

import (
	"fmt"
	"io"

	"github.com/lmika/goseq/seqdiagram/canvas"
	"github.com/lmika/goseq/seqdiagram/canvas/svgcanvas"

	"github.com/ajstarks/svgo"
)

// // Options for the SVG images
// type SvgOptions struct {
//     // If true, the viewport attribute for the SVG diagram will be set and the
//     // width and height will be converted to percentages.
//     Embedded            bool
// }

// A graphbox diagram.  This diagram is made up of uniform points.
// Each item on the diagram has the option of describing a constraint such
// as the amount of spacing around the point it requires.
type Graphic struct {
	matrix [][]matrixItem
	items  []itemInstance

	// The margin between items
	Margin Point

	// Show the grid
	ShowGrid bool

	// If true, generate a 'viewport' attribute with the image size and
	// use percentages for the original image size
	Viewport bool
}

func NewGraphic(rows, cols int) *Graphic {
	g := &Graphic{}
	g.resizeTo(rows+1, cols+1)
	return g
}

// Returns the number of items within the matrix
func (g *Graphic) Rows() int {
	return len(g.matrix) - 1
}

// Returns the number of columns within the matrix
func (g *Graphic) Cols() int {
	if len(g.matrix) > 0 {
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
		item.Item.Constraint(item.R, item.C, ConstraintApplier{g})
	}

	g.propogateDeltas()

	// Reposition the grid points
	return g.repositionGridPoints()
}

func (g *Graphic) reinitMatrix() {
	for r, row := range g.matrix {
		for c := range row {
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

// Implementation of ConstrantModifier.
//
// While updating constraints:
//
//      - The first column is used to maintain Y deltas
//      - The first row is used to maintain X deltas
//
// TODO: Instead of using loops, add cheats like only using
// the first row/column
func (g *Graphic) GridPointRect(fr, fc, tr, tc int) (int, int) {
	w, h := 0, 0
	for r := fr + 1; r <= tr; r++ {
		h += g.matrix[r][0].Delta.Y
	}
	for c := fc + 1; c <= tc; c++ {
		w += g.matrix[0][c].Delta.X
	}

	return w, h
}

func (g *Graphic) EnsureLeftIsAtleast(col, newLeft int) {
	g.matrix[0][col].Delta.X = maxInt(g.matrix[0][col].Delta.X, newLeft)
}

func (g *Graphic) EnsureTopIsAtLeast(row, newTop int) {
	g.matrix[row][0].Delta.Y = maxInt(g.matrix[row][0].Delta.Y, newTop)
}

func (g *Graphic) AddLeftToCol(col, newLeft int) {
	g.matrix[0][col].Delta.X += newLeft
}

func (g *Graphic) AddTopToRow(row, newTop int) {
	g.matrix[row][0].Delta.Y += newTop
}

// Deltas are store at 0 row (for Y deltas) and 0 col (for X deltas) for
// speed reasons.  Change all deltas to match the
func (g *Graphic) propogateDeltas() {
	for r, row := range g.matrix {
		for c := range row {
			g.matrix[r][c].Delta.X = g.matrix[0][c].Delta.X
			g.matrix[r][c].Delta.Y = g.matrix[r][0].Delta.Y
		}
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

	/*
		canvas := svg.New(w)

		if g.Viewport {
			canvas.StartviewUnit(100, 100, "%", 0, 0, sizeW, sizeH)
		} else {
			canvas.Start(sizeW, sizeH)
		}
		defer canvas.End()

		// Add styles
		canvas.Def()
		g.addStyles(canvas)
		canvas.DefEnd()

		for _, item := range g.items {
			g.drawItem(canvas, item)
		}
	*/
	canvas := svgcanvas.New(w)
	defer canvas.Close()

	canvas.SetSize(sizeW, sizeH)

	for _, item := range g.items {
		g.drawItem(canvas, item)
	}

	// Draw the grid.  Used manily for debugging
	/*
		if g.ShowGrid {
			for _, row := range g.matrix {
				for _, cell := range row {
					canvas.Circle(cell.Point.X, cell.Point.Y, 2, "brush:red;stroke:red;")
				}
			}
		}
	*/
}

// Add the style definitions, including font faces
func (g *Graphic) addStyles(canvas *svg.SVG) {
	fmt.Fprintln(canvas.Writer, "<style>")

	// !!TEMP!!
	fmt.Fprintln(canvas.Writer, "@font-face {")
	fmt.Fprintln(canvas.Writer, "  font-family: 'DejaVuSans';")
	fmt.Fprintln(canvas.Writer, "  src: url('https://fontlibrary.org/assets/fonts/dejavu-sans/f5ec8426554a3a67ebcdd39f9c3fee83/49c0f03ec2fa354df7002bcb6331e106/DejaVuSansBook.ttf') format('truetype');")
	fmt.Fprintln(canvas.Writer, "  font-weight: normal;")
	fmt.Fprintln(canvas.Writer, "  font-style: normal;")
	fmt.Fprintln(canvas.Writer, "}")
	// !!END TEMP!!

	fmt.Fprintln(canvas.Writer, "</style>")
}

// Draws the item
func (g *Graphic) drawItem(canvas canvas.Canvas, item itemInstance) {
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
	Point Point

	// Deltas between this point and the point to the left (X) and above (Y)
	Delta Point
}

type itemInstance struct {
	R, C int
	Item GraphboxItem
}
