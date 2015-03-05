package graphbox

// Divider style
type DividerStyle struct {
    Font        Font
    FontSize    int
    Padding     Point
    Margin      Point
}


// A divider.  This spans the entire diagram.
type Divider struct {
    TC              int

    style           DividerStyle
    textBox         *TextBox
    textBoxRect     Rect
    marginRect      Rect
}

func NewDivider(toCol int, text string, style DividerStyle) *Divider {
    textBox := NewTextBox(style.Font, style.FontSize, MiddleTextAlign)
    textBox.AddText(text)
    textBoxRect := textBox.BoundingRect()
    marginRect := textBoxRect.BlowOut(style.Padding)

    return &Divider{toCol, style, textBox, textBoxRect, marginRect}
}

func (div *Divider) Constraint(r, c int, applier ConstraintApplier) {
    // There must be enought horizontal space to accomodate the text
    // and vertical space to display the divider
    requiredHeight := div.marginRect.H + div.style.Margin.Y * 2
    requiredWidth := div.marginRect.W + div.style.Margin.X * 2

    applier.Apply(AddSizeConstraint{r, c, 0, 0, requiredHeight / 2, requiredHeight / 2})
    applier.Apply(TotalSizeConstraint{r - 1, c, r, div.TC, requiredWidth, 0})    
}

func (div *Divider) Draw(ctx DrawContext, point Point) {
    fx, fy := point.X, point.Y
    if point, isPoint := ctx.PointAt(ctx.R, div.TC) ; isPoint {
        tx, _ := point.X, point.Y
        centerX := fx + (tx - fx) / 2
        centerY := fy

        // Draw the boundary
        borderRect := Rect{fx, fy - div.marginRect.H / 2, tx - fx, div.marginRect.H}
        ctx.Canvas.Rect(borderRect.X, borderRect.Y, borderRect.W, borderRect.H, "fill:white;stroke:white;")

        // Center text
        div.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
    }    
}