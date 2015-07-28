package graphbox

type TitleStyle struct {
    Font            Font
    FontSize        int
    Padding         Point
}


// A title
type Title struct {
    TC              int
    style           TitleStyle
    textBox         *TextBox
    textBoxRect     Rect
}

func NewTitle(toCol int, text string, style TitleStyle) *Title {
    textBox := NewTextBox(style.Font, style.FontSize, LeftTextAlign)
    textBox.AddText(text)

    brect := textBox.BoundingRect()
    return &Title{toCol, style, textBox, brect}
}

func (al *Title) Constraint(r, c int, applier ConstraintApplier) {
    h := al.textBoxRect.H + al.style.Padding.Y
    w := al.textBoxRect.W

    applier.Apply(SizeConstraint{r, c, 0, 0, h, al.style.Padding.Y})
    applier.Apply(TotalSizeConstraint{r, c, r + 1, al.TC, w + al.style.Padding.X, 0})
}

func (al *Title) Draw(ctx DrawContext, point Point) {
    fx, fy := point.X, point.Y
    if point, isPoint := ctx.PointAt(ctx.R, al.TC) ; isPoint {
        tx, _ := point.X, point.Y
        _ = tx

        textX := fx + al.style.Padding.X // + (tx - fx) / 2
        textY := fy - al.style.Padding.Y
        al.renderMessage(ctx, textX, textY)
    }
}

func (al *Title) renderMessage(ctx DrawContext, tx, ty int) {
    rect := al.textBoxRect.PositionAt(tx, ty, SouthWestGravity)

    ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "fill:white;stroke:white;")
    al.textBox.Render(ctx.Canvas, tx, ty, SouthWestGravity)
}
