package graphbox


// Draws a rectangle around the inner rectangle frame
type Rectangle struct {
    // Width and height of the rectangle
    W, H        int
}

func (r *Rectangle) Size() (int, int) {
    return r.W, r.H
}

func (r *Rectangle) Draw(ctx *DrawContext, frame BoxFrame) {
    ctx.Canvas.Rect(frame.InnerRect.X, frame.InnerRect.Y, frame.InnerRect.W, frame.InnerRect.H, 
            "fill:none;stroke:black")
}
