package graphbox


// TEMP
type Icon struct {
}

func (icon *Icon) Size() (w int, h int) {
    return 48, 48
}


// Styling options for the actor rect
type ActorIconBoxStyle struct {
    Font        Font
    FontSize    int
    Padding     Point
    Margin      Point
    IconGap     int
}


// The actor icon box
type ActorIconBox struct {
    //Caption     string
    textBox     *TextBox
    Icon        *Icon
    style       ActorIconBoxStyle
    pos         ActorBoxPos
}

func NewActorIconBox(text string, icon *Icon, style ActorIconBoxStyle, pos ActorBoxPos) *ActorIconBox {
    textBox := NewTextBox(style.Font, style.FontSize, MiddleTextAlign)
    textBox.AddText(text)

    //trect := textBox.BoundingRect()
    //brect := trect.BlowOut(style.Padding)

    return &ActorIconBox{textBox, icon, style, pos}
}

func (tr *ActorIconBox) Constraint(r, c int, applier ConstraintApplier) {
    posHoriz, posVert := tr.pos & 0xFF00, tr.pos & 0xFF
    iconW, iconH := tr.Icon.Size()

    brect := tr.textBox.BoundingRect()
    w := maxInt(iconW, brect.W) + tr.style.Padding.X

    topH := iconH / 2
    bottomH := iconH / 2 + brect.H + tr.style.IconGap + tr.style.Padding.Y
    marginX := tr.style.Margin.X

    if posVert == TopActorBox {
        if posHoriz == LeftActorBox {
            applier.Apply(SizeConstraint{r, c, w / 2, marginX / 2, 0, 0})
            applier.Apply(AddSizeConstraint{r, c, 0, w / 2, 0, 0})
        } else if posHoriz == RightActorBox {
            applier.Apply(SizeConstraint{r, c, marginX / 2, w / 2, 0, 0})
            applier.Apply(AddSizeConstraint{r, c, w / 2, 0, 0, 0})
        } else {
            applier.Apply(SizeConstraint{r, c, marginX / 2, marginX / 2, 0, 0})
            applier.Apply(AddSizeConstraint{r, c, w / 2, w / 2, 0, 0})
        }
        applier.Apply(SizeConstraint{r, c, 0, 0, topH, bottomH})
    }
}

func (tr *ActorIconBox) Draw(ctx DrawContext, point Point) {
    centerX, centerY := point.X, point.Y

    // Draw the icon
    iconW, iconH := tr.Icon.Size()
    iconX, iconY := centerX - iconW / 2, centerY - iconH / 2

    // Draw the text
    brect := tr.textBox.BoundingRect()
    textY := iconY + iconH + tr.style.IconGap
    rect := brect.PositionAt(centerX, textY, NorthGravity)

    ctx.Canvas.Rect(rect.X, rect.Y - tr.style.IconGap, rect.W, rect.H + tr.style.IconGap, "stroke:white;fill:white;stroke-width:2px;")
    tr.textBox.Render(ctx.Canvas, centerX, textY, NorthGravity)

    // Draw the icon
    // TEMP
    ctx.Canvas.Rect(iconX, iconY, iconW, iconH, "stroke:black;fill:white;stroke-width:2px;")
}