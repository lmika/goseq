package graphbox

type NoteBoxPos int
const (
    CenterNotePos   NoteBoxPos     =   iota
    LeftNotePos                     =   iota
    RightNotePos                    =   iota
)

// Styling options for the actor rect
type NoteBoxStyle struct {
    Font        Font
    FontSize    int
    Padding     Point
    Position    NoteBoxPos
}

// Draws an object instance
type NoteBox struct {
    frameRect   Rect
    style       NoteBoxStyle
    textBox     *TextBox
    pos         NoteBoxPos
}

func NewNoteBox(text string, style NoteBoxStyle, pos NoteBoxPos) *NoteBox {
    var textAlign TextAlign = MiddleTextAlign

    textBox := NewTextBox(style.Font, style.FontSize, textAlign)
    textBox.AddText(text)

    trect := textBox.BoundingRect()
    brect := trect.BlowOut(style.Padding)

    return &NoteBox{brect, style, textBox, pos}
}

func (tr *NoteBox) Constraint(r, c int) Constraint {
    var horizConstraint Constraint
    if (tr.pos == LeftNotePos) {
        horizConstraint = SizeConstraint{r, c, tr.frameRect.W + 8, 8, 0, 0}
    } else if (tr.pos == RightNotePos) {
        horizConstraint = SizeConstraint{r, c, 8, tr.frameRect.W + 8, 0, 0}
    } else {
        horizConstraint = SizeConstraint{r, c, tr.frameRect.W, tr.frameRect.W, 0, 0}
    }

    return Constraints([]Constraint {
        horizConstraint,
        AddSizeConstraint{r, c, 0, 0, tr.frameRect.H / 2 + 4, tr.frameRect.H / 2 + 4},
    })
}

func (r *NoteBox) Draw(ctx DrawContext, point Point) {
    centerX, centerY := point.X, point.Y

    if (r.pos == CenterNotePos) {
        rect := r.frameRect.PositionAt(centerX, centerY, CenterGravity)
        ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "stroke:black;fill:white")
        r.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
    } else if (r.pos == LeftNotePos) {
        offsetX := centerX - 8
        textOffsetX := centerX - r.style.Padding.X - 8
        rect := r.frameRect.PositionAt(offsetX, centerY, EastGravity)
        ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "stroke:black;fill:white")
        r.textBox.Render(ctx.Canvas, textOffsetX, centerY, EastGravity)
    } else if (r.pos == RightNotePos) {
        offsetX := centerX + 4 * 2
        textOffsetX := centerX + r.style.Padding.X + 4 * 2
        rect := r.frameRect.PositionAt(offsetX, centerY, WestGravity)
        ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "stroke:black;fill:white")
        r.textBox.Render(ctx.Canvas, textOffsetX, centerY, WestGravity)
    }
}