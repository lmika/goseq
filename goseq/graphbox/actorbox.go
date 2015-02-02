package graphbox

type ActorBoxPos int
const (
    TopActorBox       ActorBoxPos     =   iota
    BottomActorBox                    =   iota
)

// Styling options for the actor rect
type ActorBoxStyle struct {
    Font        Font
    FontSize    int
    Padding     Point
    Margin      Point
}

// Draws an object instance
type ActorBox struct {
    frameRect   Rect
    style       ActorBoxStyle
    textBox     *TextBox
    pos         ActorBoxPos
}

func NewActorBox(text string, style ActorBoxStyle, pos ActorBoxPos) *ActorBox {
    var textAlign TextAlign = MiddleTextAlign

    textBox := NewTextBox(style.Font, style.FontSize, textAlign)
    textBox.AddText(text)

    trect := textBox.BoundingRect()
    brect := trect.BlowOut(style.Padding)

    return &ActorBox{brect, style, textBox, pos}
}

func (tr *ActorBox) Constraint(r, c int) Constraint {
    var vertConstraint Constraint

    if (tr.pos == TopActorBox) {
        vertConstraint = SizeConstraint{r, c, 0, 0, tr.frameRect.H / 2, tr.frameRect.H / 2 + tr.style.Margin.Y}
    } else {
        vertConstraint = SizeConstraint{r, c, 0, 0, tr.frameRect.H / 2 + tr.style.Margin.Y, tr.frameRect.H / 2}
    }

    if (tr.pos == TopActorBox) {
        if (c == 1) {
            return Constraints([]Constraint {
                vertConstraint,
                SizeConstraint{r, c, tr.frameRect.W / 2, 0, 0, 0},
                AddSizeConstraint{r, c, 0, tr.frameRect.W / 2 + tr.style.Margin.X, 0, 0},
            })
        } else {
            return Constraints([]Constraint {
                vertConstraint,
                AddSizeConstraint{r, c, tr.frameRect.W / 2 + tr.style.Margin.X, tr.frameRect.W / 2 + tr.style.Margin.X, 0, 0},
            })
        }
    } else {
        return vertConstraint
    }
}

func (r *ActorBox) Draw(ctx DrawContext, point Point) {
    centerX, centerY := point.X, point.Y

    rect := r.frameRect.PositionAt(centerX, centerY, CenterGravity)
    ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "stroke:black;fill:white")
    r.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
}