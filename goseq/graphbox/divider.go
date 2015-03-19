package graphbox

// Divider shape
type DividerShape int

const (
    // A rectangle which will span the entire graphic from end to end.
    // The text will be centered in front of it.    
    DSFullRect  DividerShape    = iota

    // A line which will span the entire grapic.  The text will be
    // centered in front of it.
    DSFullLine
)

// Divider style
type DividerStyle struct {
    Font        Font
    FontSize    int
    Padding     Point
    Margin      Point
    TextPadding Point
    Shape       DividerShape
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

        borderRect := Rect{fx, fy - div.marginRect.H / 2, tx - fx, div.marginRect.H}
        textBoxRect := div.textBoxRect.PositionAt(centerX, centerY, CenterGravity).BlowOut(div.style.TextPadding)

        // Draw the shape and text
        switch div.style.Shape {
        case DSFullRect:
            ctx.Canvas.Rect(borderRect.X, borderRect.Y, borderRect.W, borderRect.H, "fill:white;stroke:white;")
            div.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
        case DSFullLine:
            // Draw the rectangle for clearing the image
            ctx.Canvas.Rect(borderRect.X, borderRect.Y, borderRect.W, borderRect.H, "fill:white;stroke:white;")
            ctx.Canvas.Line(borderRect.X, centerY, borderRect.W, centerY, "fill:white;stroke:black;stroke-width:2px;stroke-dasharray:16,8")

            ctx.Canvas.Rect(textBoxRect.X, textBoxRect.Y, textBoxRect.W, textBoxRect.H, "fill:white;stroke:white;")
            div.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
        }
    }    
}

