package graphbox


// The object lifeline
type LifeLine struct {
    TR, TC      int
}

func (ll *LifeLine) Constraint(r, c int) Constraint {
    return nil
}

//func (ll *LifeLine) Draw(ctx DrawContext, frame BoxFrame) {
func (ll *LifeLine) Draw(ctx DrawContext, point Point) {
    fx, fy := point.X, point.Y
    if point, isPoint := ctx.PointAt(ll.TR, ll.TC) ; isPoint {
        tx, ty := point.X, point.Y

        ctx.Canvas.Line(fx, fy, tx, ty, "stroke:black;stroke-dasharray:8,8")
    }
}