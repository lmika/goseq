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