// Graphics display model
//

package graphbox

import (
    "github.com/ajstarks/svgo"
)


type GraphboxItem interface {

    // Call to draw this box 
    Draw(ctx *DrawContext, frame BoxFrame)
}

// An item that takes up space within a cell
type Graphbox2DItem interface {

    // The width and height of this particlar item
    Size()      (int, int)
}


// A drawing context
type DrawContext struct {
    Canvas          *svg.SVG
    Graphic         *Graphic
}

// Returns the outer rectangle of a particular cell
func (dc *DrawContext) GridRect(r, c int) (Rect, bool) {
    return dc.Graphic.outerRectAtCell(r, c)
}


// An anchor point located in a rectangle at 0, 0 with the w, h passed in
type Gravity         func(w, h int) (int, int)

var CenterGravity Gravity = func(w, h int) (int, int) { return w / 2, h / 2 }

// A rectangle
type Rect struct {
    X, Y            int
    W, H            int
}

// Returns a point located at a specific gravity within the rectangle
func (r Rect) PointAt(gravity Gravity) (int, int) {
    lx, ly := gravity(r.W, r.H)
    return r.X + lx, r.Y + ly
}

// Returns a new rect which will be a rectangle with the 
// given dimensions centered in this rect
func (r Rect) CenteredRect(w, h int) Rect {
    x := r.X + (r.W / 2) - w / 2
    y := r.Y + (r.H / 2) - h / 2
    return Rect{x, y, w, h}
}

// A point
type Point struct {
    X, Y            int
}

// A box frame
type BoxFrame struct {
    // The outer rectangle.  This encompasses margins, etc.
    OuterRect       Rect

    // The inner rectangle.
    InnerRect       Rect
}