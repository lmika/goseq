// Graphics display model
//

package graphbox

import (
    "github.com/ajstarks/svgo"
)


type GraphboxItem interface {

    // The width and height of this particlar item
    Size()      (int, int)

    // Call to draw this box 
    Draw(ctx *DrawContext, frame BoxFrame)
}


// A drawing context
type DrawContext struct {
    Canvas          *svg.SVG
}


// A rectangle
type Rect struct {
    X, Y            int
    W, H            int
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