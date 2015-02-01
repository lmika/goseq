package graphbox

import (
    "fmt"
)

type TextRectPos int
const (
    CenterPos   TextRectPos     =   iota
    LeftPos                     =   iota
    RightPos                    =   iota
)

// Styling options for the actor rect
type TextRectStyle struct {
    Font        Font
    FontSize    int
    Padding     Point
    Position    TextRectPos
}

// Draws an object instance
type TextRect struct {
    frameRect   Rect
    style       TextRectStyle
    textBox     *TextBox
    pos         TextRectPos
}

func NewTextRect(text string, style TextRectStyle, pos TextRectPos) *TextRect {
    var textAlign TextAlign = MiddleTextAlign

    textBox := NewTextBox(style.Font, style.FontSize, textAlign)
    textBox.AddText(text)

    trect := textBox.BoundingRect()
    brect := trect.BlowOut(style.Padding)

    return &TextRect{brect, style, textBox, pos}
}

/*
func (r *TextRect) Size() (int, int) {
    if (r.pos == CenterPos) {
        return r.frameRect.W, r.frameRect.H
    } else {
        return 0, r.frameRect.H
    }
}

func (r *TextRect) Margin() (int, int, int, int) {
    if (r.pos == LeftPos) {
        return r.frameRect.W + 8, 0, 0, 0
    } else if (r.pos == RightPos) {
        return 0, r.frameRect.W + 8, 0, 0
    } else {
        return 0, 0, 0, 0
    }
}
*/

func (tr *TextRect) Constraint(r, c int) Constraint {
    return SizeConstraint{r, c, tr.frameRect.W, tr.frameRect.W, tr.frameRect.H / 2, tr.frameRect.H / 2}
}

func (r *TextRect) Draw(ctx DrawContext, point Point) {
    centerX, centerY := point.X, point.Y

    if (r.pos == CenterPos) {
        rect := r.frameRect.PositionAt(centerX, centerY, CenterGravity)
        ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "stroke:black;fill:white")
        r.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
    } else if (r.pos == LeftPos) {
        offsetX := centerX - 8
        textOffsetX := centerX - r.style.Padding.X - 8
        rect := r.frameRect.PositionAt(offsetX, centerY, EastGravity)
        ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "stroke:black;fill:white")
        r.textBox.Render(ctx.Canvas, textOffsetX, centerY, EastGravity)
    } else if (r.pos == RightPos) {
        offsetX := centerX + 4 * 2
        textOffsetX := centerX + r.style.Padding.X + 4 * 2
        rect := r.frameRect.PositionAt(offsetX, centerY, WestGravity)
        ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "stroke:black;fill:white")
        r.textBox.Render(ctx.Canvas, textOffsetX, centerY, WestGravity)
    }
}


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


type ActivityLineStyle struct {
    Font            Font
    FontSize        int
    PaddingTop      int
    PaddingBottom   int
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
//    r, _ := MeasureFontRect(style.Font, style.FontSize, text, 0, 0, NorthWestGravity)

    textBox := NewTextBox(style.Font, style.FontSize, MiddleTextAlign)
    textBox.AddText(text)

    brect := textBox.BoundingRect()
    return &ActivityLine{toCol, style, textBox, brect}
}

/*
func (al *ActivityLine) Size() (int, int) {
    return 50, al.textBoxRect.H + al.style.PaddingTop + al.style.PaddingBottom + al.style.TextGap
}
*/

func (al *ActivityLine) Constraint(r, c int) Constraint {
    h := al.textBoxRect.H + al.style.PaddingTop + al.style.PaddingBottom + al.style.TextGap
    w := al.textBoxRect.W

    return TotalSizeConstraint{r - 1, c, r, al.TC, w, h}
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