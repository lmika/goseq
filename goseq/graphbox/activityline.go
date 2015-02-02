package graphbox

import (
    "fmt"
)

type ActivityLineStyle struct {
    Font            Font
    FontSize        int
    Margin          Point
    TextGap         int
}

// Returns the text style
func (as ActivityLineStyle) textStyle() string {
    s := SvgStyle{}

    s.Set("font-family", as.Font.SvgName())
    s.Set("font-size", fmt.Sprintf("%dpx", as.FontSize))

    return s.ToStyle()
}


// An activity arrow
type ActivityLine struct {
    TC              int
    style           ActivityLineStyle
    textBox         *TextBox
    textBoxRect     Rect
}

func NewActivityLine(toCol int, text string, style ActivityLineStyle) *ActivityLine {
    textBox := NewTextBox(style.Font, style.FontSize, MiddleTextAlign)
    textBox.AddText(text)

    brect := textBox.BoundingRect()
    return &ActivityLine{toCol, style, textBox, brect}
}

func (al *ActivityLine) Constraint(r, c int) Constraint {
    h := al.textBoxRect.H + al.style.Margin.Y + al.style.TextGap
    w := al.textBoxRect.W

    _ = h
    _ = w
    return Constraints([]Constraint{ 
        AddSizeConstraint{r, c, 0, 0, h, al.style.Margin.Y},
        TotalSizeConstraint{r - 1, c, r, al.TC, w + al.style.Margin.X * 2, 0},
    })
}

// func (al *ActivityLine) Draw(ctx DrawContext, frame BoxFrame) {
func (al *ActivityLine) Draw(ctx DrawContext, point Point) {

    fx, fy := point.X, point.Y
    if point, isPoint := ctx.PointAt(ctx.R, al.TC) ; isPoint {
        tx, ty := point.X, point.Y

        ctx.Canvas.Line(fx, fy, tx, ty, "stroke:black")
        al.drawArrow(ctx, tx, ty, al.TC > ctx.C)

        textX := fx + (tx - fx) / 2
        textY := ty - al.style.TextGap
        al.renderMessage(ctx, textX, textY)
    }
}

func (al *ActivityLine) renderMessage(ctx DrawContext, tx, ty int) {
    //rect, textPoint := MeasureFontRect(al.style.Font, al.style.FontSize, al.Text, tx, ty, SouthGravity)
    rect := al.textBoxRect.PositionAt(tx, ty, SouthGravity)

    ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "fill:white;stroke:white;")
    al.textBox.Render(ctx.Canvas, tx, ty, SouthGravity)
    //ctx.Canvas.Text(textPoint.X, textPoint.Y, al.Text, al.style.textStyle())
}

// TODO: Type of arrow
func (al *ActivityLine) drawArrow(ctx DrawContext, x, y int, isRight bool) {
    var xs, ys []int

    ys = []int { y - 5, y, y + 5 }
    if isRight {
        xs = []int { x - 8, x, x - 8 }
    } else {
        xs = []int { x + 8, x, x + 8 }
    }

    ctx.Canvas.Polyline(xs, ys, "stroke:black;fill:none")
}