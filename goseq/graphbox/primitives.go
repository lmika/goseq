package graphbox

import (
    //"log"
)


// Draws an object instance
type ActorRect struct {
    // Width and height of the rectangle
    W, H        int
    Text        string
    Font        Font
    Padding     Point
}

func NewActorRect(text string, font Font) *ActorRect {
    w, h := font.Measure(text, 18.0)
    padding := Point{16, 8}
    //log.Printf("W,H = %d,%d", w +  * 2, h)
    
    return &ActorRect{w, h, text, font, padding}
}

func (r *ActorRect) Size() (int, int) {
    return r.W + r.Padding.X * 2, r.H + r.Padding.Y * 2
}

func (r *ActorRect) Draw(ctx DrawContext, frame BoxFrame) {
    centeredRect := frame.InnerRect.CenteredRect(r.W, r.H)
    tx, ty := centeredRect.PointAt(CenterGravity)

    ctx.Canvas.Rect(centeredRect.X - r.Padding.X, centeredRect.Y - r.Padding.Y,
            centeredRect.W + r.Padding.X * 2, centeredRect.H + r.Padding.Y * 2, 
            "fill:none;stroke:black;fill:white")

    ctx.Canvas.Text(tx, ty, r.Text, "text-anchor:middle;dominant-baseline:middle;font-size:18px;fill:black;" +
        "font-family:" + r.Font.SvgName() + "")
}


// The object lifeline
type LifeLine struct {
    TR, TC      int
}

func (ll *LifeLine) Draw(ctx DrawContext, frame BoxFrame) {
    fx, fy := frame.InnerRect.PointAt(CenterGravity)
    if toOuterRect, isCell := ctx.GridRect(ll.TR, ll.TC) ; isCell {
        tx, ty := toOuterRect.PointAt(CenterGravity)

        ctx.Canvas.Line(fx, fy, tx, ty, "stroke:black;stroke-dasharray:8,8")
    }
}


// An activity arrow
type ActivityLine struct {
    TC           int
    Text         string
}

func (r *ActivityLine) Size() (int, int) {
    return 50, 70
}

func (al *ActivityLine) Draw(ctx DrawContext, frame BoxFrame) {
    lineGravity := AtSpecificGravity(0.5, 0.7)

    fx, fy := frame.InnerRect.PointAt(lineGravity)
    if toOuterRect, isCell := ctx.GridRect(ctx.R, al.TC) ; isCell {
        tx, ty := toOuterRect.PointAt(lineGravity)
        ctx.Canvas.Line(fx, fy, tx, ty, "stroke:black")

        al.drawArrow(ctx, tx, ty, al.TC > ctx.C)

        textX := fx + (tx - fx) / 2
        textY := ty - 15
        ctx.Canvas.Text(textX, textY, al.Text, "text-anchor:middle;dominant-baseline:bottom;font-size:14px;fill:black")
    }
}

// TODO: Type of arrow
func (al *ActivityLine) drawArrow(ctx DrawContext, x, y int, isRight bool) {
    var xs, ys []int

    ys = []int { y - 5, y, y + 5 }
    if isRight {
        xs = []int { x - 8, x, x - 8 }
    } else {
        xs = []int { x + 8, x, x + 8 }
    }

    ctx.Canvas.Polyline(xs, ys, "stroke:black;fill:none")
}