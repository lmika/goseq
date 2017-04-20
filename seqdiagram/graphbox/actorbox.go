package graphbox

// ActorBoxPos is used to manage the flags representing the actor boxes position
type ActorBoxPos int

// These flags are for the vertical position
const (
	TopActorBox    ActorBoxPos = iota
	BottomActorBox             = iota
)

// These flags are for the horizontal position
const (
	LeftActorBox   ActorBoxPos = (iota << 8)
	MiddleActorBox             = (iota << 8)
	RightActorBox              = (iota << 8)
)

// ActorBoxStyle defines styling options for the actor boxes
type ActorBoxStyle struct {
	Font      Font
	FontSize  int
	Padding   Point
	Margin    Point
	Color     string
	TextColor string
}

// ActorBox represents an a actor
type ActorBox struct {
	frameRect Rect
	style     ActorBoxStyle
	textBox   *TextBox
	pos       ActorBoxPos
}

// NewActorBox returns a new actor
func NewActorBox(text string, style ActorBoxStyle, pos ActorBoxPos) *ActorBox {
	var textAlign TextAlign = MiddleTextAlign

	textBox := NewTextBox(style.Font, style.FontSize, textAlign)
	textBox.Color = style.TextColor
	textBox.AddText(text)

	trect := textBox.BoundingRect()
	brect := trect.BlowOut(style.Padding)

	return &ActorBox{brect, style, textBox, pos}
}

func (tr *ActorBox) Constraint(r, c int, applier ConstraintApplier) {
	var vertConstraint Constraint
	posHoriz, posVert := tr.pos&0xFF00, tr.pos&0xFF

	if posVert == TopActorBox {
		vertConstraint = SizeConstraint{r, c, 0, 0, tr.frameRect.H / 2, tr.frameRect.H/2 + tr.style.Margin.Y}
	} else {
		vertConstraint = SizeConstraint{r, c, 0, 0, tr.frameRect.H/2 + tr.style.Margin.Y, tr.frameRect.H / 2}
	}

	if posVert == TopActorBox {
		if posHoriz == LeftActorBox {
			applier.Apply(vertConstraint)
			applier.Apply(SizeConstraint{r, c, tr.frameRect.W / 2, 0, 0, 0})
			applier.Apply(AddSizeConstraint{r, c, 0, tr.frameRect.W/2 + tr.style.Margin.X, 0, 0})
		} else if posHoriz == RightActorBox {
			applier.Apply(vertConstraint)
			applier.Apply(SizeConstraint{r, c, 0, tr.frameRect.W / 2, 0, 0})
			applier.Apply(AddSizeConstraint{r, c, tr.frameRect.W/2 + tr.style.Margin.X, 0, 0, 0})
		} else {
			applier.Apply(vertConstraint)
			applier.Apply(AddSizeConstraint{r, c, tr.frameRect.W/2 + tr.style.Margin.X, tr.frameRect.W/2 + tr.style.Margin.X, 0, 0})
		}
	} else {
		applier.Apply(vertConstraint)
	}
}

func (r *ActorBox) Draw(ctx DrawContext, point Point) {
	s := SvgStyle{}
	s.Set("stroke", r.style.Color)
	s.Set("fill", "white")
	s.Set("stroke-width", "2px")

	centerX, centerY := point.X, point.Y

	rect := r.frameRect.PositionAt(centerX, centerY, CenterGravity)
	ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, s.ToStyle())
	r.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
}
