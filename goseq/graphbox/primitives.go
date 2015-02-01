package graphbox

import (
    "fmt"
    //"log"
)


// Styling options for the actor rect
type ActorRectStyle struct {
    Font        Font
    FontSize    int
    Padding     Point
}

// Returns the text style
func (rs ActorRectStyle) textStyle() string {
    s := SvgStyle{}

    s.Set("font-family", rs.Font.SvgName())
    s.Set("font-size", fmt.Sprintf("%dpx", rs.FontSize))

    return s.ToStyle()
}

// Measure a string based on the font settings
func (rs ActorRectStyle) measure(s string) (int, int) {
    return rs.Font.Measure(s, float64(rs.FontSize))
}



// Draws an object instance
type ActorRect struct {
    // Width and height of the rectangle
    Text        string
    w, h        int
    style       ActorRectStyle
}

func NewActorRect(text string, font Font) *ActorRect {
    style := ActorRectStyle{
        Font:       font,
        FontSize:   16,
        Padding:    Point{16, 8},
    }

    trect, _ := MeasureFontRect(style.Font, style.FontSize, text, 0, 0, NorthWestGravity)
    brect := trect.BlowOut(style.Padding)

    return &ActorRect{text, brect.W, brect.H, style}
}

func (r *ActorRect) Size() (int, int) {
    return r.w, r.h
}

func (r *ActorRect) Draw(ctx DrawContext, frame BoxFrame) {
    centerX, centerY := frame.InnerRect.PointAt(CenterGravity)
    trect, tp := MeasureFontRect(r.style.Font, r.style.FontSize, r.Text, centerX, centerY, CenterGravity)
    brect := trect.BlowOut(r.style.Padding)

    ctx.Canvas.Rect(brect.X, brect.Y, brect.W, brect.H, "stroke:black;fill:white")
    ctx.Canvas.Text(tp.X, tp.Y, r.Text, r.style.textStyle())
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

func NewActivityLine(toCol int, text string, font Font) *ActivityLine {
    style := ActivityLineStyle{
        Font:           font,
        FontSize:       14,
        PaddingTop:     8,
        PaddingBottom:  8,
        TextGap:        8,
    }

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