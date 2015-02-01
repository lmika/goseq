package graphbox

import (
    "fmt"
    //"log"
)


// Styling options for the actor rect
type TextRectStyle struct {
    Font        Font
    FontSize    int
    Padding     Point
}

// Draws an object instance
type TextRect struct {
    // Width and height of the rectangle
    Text        string

    frameRect   Rect
    style       TextRectStyle
    textBox     *TextBox
}

func NewTextRect(text string, style TextRectStyle) *TextRect {
    textBox := NewTextBox(style.Font, style.FontSize, MiddleTextAlign)
    textBox.AddText(text)

    trect := textBox.BoundingRect(0, 0, NorthWestGravity)
    brect := trect.BlowOut(style.Padding)


    return &TextRect{text, brect, style, textBox}
}

func (r *TextRect) Size() (int, int) {
    return r.frameRect.W, r.frameRect.H
}

func (r *TextRect) Draw(ctx DrawContext, frame BoxFrame) {
    centerX, centerY := frame.InnerRect.PointAt(CenterGravity)

    rect := r.frameRect.PositionTo(centerX, centerY, CenterGravity)

    ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "stroke:black;fill:white")
    r.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
}


// The object lifeline
type LifeLine struct {
    TR, TC      int
}

func (ll *LifeLine) Draw(ctx DrawContext, frame BoxFrame) {
    fx, fy := frame.InnerRect.PointAt(CenterGravity)
    if toOuterRect, isCell := ctx.GridRect(ll.TR, ll.TC) ; isCell {
        tx, ty := toOuterRect.PointAt(CenterGravity)

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
    TC           int
    Text         string
    style        ActivityLineStyle

    height       int
}

func NewActivityLine(toCol int, text string, style ActivityLineStyle) *ActivityLine {
    r, _ := MeasureFontRect(style.Font, style.FontSize, text, 0, 0, NorthWestGravity)
    height := r.H
    return &ActivityLine{toCol, text, style, height}
}

func (al *ActivityLine) Size() (int, int) {
    return 50, al.height + al.style.PaddingTop + al.style.PaddingBottom + al.style.TextGap
}

func (al *ActivityLine) Draw(ctx DrawContext, frame BoxFrame) {
    lineGravity := SouthGravity

    fx, fy := frame.InnerRect.PointAt(lineGravity)
    fy -= al.style.PaddingBottom
    if toOuterRect, isCell := ctx.GridRect(ctx.R, al.TC) ; isCell {
        tx, ty := toOuterRect.PointAt(lineGravity)
        ty -= al.style.PaddingBottom

        ctx.Canvas.Line(fx, fy, tx, ty, "stroke:black")
        al.drawArrow(ctx, tx, ty, al.TC > ctx.C)

        textX := fx + (tx - fx) / 2
        textY := ty - al.style.TextGap
        al.renderMessage(ctx, textX, textY)
    }
}

func (al *ActivityLine) renderMessage(ctx DrawContext, tx, ty int) {
    rect, textPoint := MeasureFontRect(al.style.Font, al.style.FontSize, al.Text, tx, ty, SouthGravity)

    ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "fill:white;stroke:white;")
    ctx.Canvas.Text(textPoint.X, textPoint.Y, al.Text, al.style.textStyle())
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