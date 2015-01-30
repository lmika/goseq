package graphbox


// Draws a rectangle centered within the inner rectangle
type Rectangle struct {
    // Width and height of the rectangle
    W, H        int
}

func (r *Rectangle) Size() (int, int) {
    return r.W, r.H
}

func (r *Rectangle) Draw(ctx *DrawContext, frame BoxFrame) {
    centeredRect := frame.InnerRect.CenteredRect(r.W, r.H)
    ctx.Canvas.Rect(centeredRect.X, centeredRect.Y, centeredRect.W, centeredRect.H, 
            "fill:none;stroke:black")
}
