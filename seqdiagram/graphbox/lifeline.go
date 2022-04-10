package graphbox

type LifeLineStyle struct {
	Color string
}

// The object lifeline
type LifeLine struct {
	TR, TC int
	Style  LifeLineStyle
}

func (ll *LifeLine) Constraint(r, c int, applier ConstraintApplier) {
}

func (ll *LifeLine) Draw(ctx DrawContext, point Point) {
	s := SvgStyle{}
	s.Set("stroke", ll.Style.Color)
	s.Set("stroke-dasharray", "8,8")
	s.Set("stroke-width", "2px")

	fx, fy := point.X, point.Y
	if point, isPoint := ctx.PointAt(ll.TR, ll.TC); isPoint {
		tx, ty := point.X, point.Y

		ctx.Canvas.Line(fx, fy, tx, ty, s.ToStyle())
	}
}
