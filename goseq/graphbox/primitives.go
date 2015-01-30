package graphbox


// Draws an object instance
type ActorRect struct {
    // Width and height of the rectangle
    W, H        int
    Text        string
}

func (r *ActorRect) Size() (int, int) {
    return r.W, r.H
}

func (r *ActorRect) Draw(ctx *DrawContext, frame BoxFrame) {
    centeredRect := frame.InnerRect.CenteredRect(r.W, r.H)
    tx, ty := centeredRect.PointAt(CenterGravity)

    ctx.Canvas.Rect(centeredRect.X, centeredRect.Y, centeredRect.W, centeredRect.H, 
            "fill:none;stroke:black;fill:white")
    ctx.Canvas.Text(tx, ty, r.Text, "text-anchor:middle;dominant-baseline:middle;font-size:18px;fill:black")
}


// The object lifeline
type LifeLine struct {
    TR, TC      int
}

func (ll *LifeLine) Draw(ctx *DrawContext, frame BoxFrame) {
    fx, fy := frame.InnerRect.PointAt(CenterGravity)
    if toOuterRect, isCell := ctx.GridRect(ll.TR, ll.TC) ; isCell {
        tx, ty := toOuterRect.PointAt(CenterGravity)

        ctx.Canvas.Line(fx, fy, tx, ty, "stroke:black;stroke-dasharray:8,8")
    }
}
