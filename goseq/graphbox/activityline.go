package graphbox

import (
    "fmt"
)

type ActivityArrowHead int
const (
    SolidArrowHead  ActivityArrowHead   =   iota
    OpenArrowHead                       =   iota
    BarbArrowHead                       =   iota
)

type ActivityArrowStem int
const (
    SolidArrowStem  ActivityArrowStem   =   iota
    DashedArrowStem                     =   iota
)





type ActivityLineStyle struct {
    Font            Font
    FontSize        int
    Margin          Point
    TextGap         int
    ArrowHead       ActivityArrowHead
    ArrowStem       ActivityArrowStem
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

func (al *ActivityLine) Constraint(r, c int, applier ConstraintApplier) {
    /*
    h := al.textBoxRect.H + al.style.Margin.Y + al.style.TextGap
    w := al.textBoxRect.W

    _ = h
    _ = w
    */
    lc, rc := c, al.TC
    if al.TC < c {
        lc, rc = al.TC, c
    }
 
    applier.Apply(AddSizeConstraint{r, c, 0, 0, h, al.style.Margin.Y})
    applier.Apply(TotalSizeConstraint{r - 1, lc, r, rc, w + al.style.Margin.X * 2, 0})
}

func (al *ActivityLine) Draw(ctx DrawContext, point Point) {

    fx, fy := point.X, point.Y
    if point, isPoint := ctx.PointAt(ctx.R, al.TC) ; isPoint {
        tx, ty := point.X, point.Y

        textX := fx + (tx - fx) / 2
        textY := ty - al.style.TextGap
        al.renderMessage(ctx, textX, textY)        

        if al.style.ArrowStem == DashedArrowStem {
            ctx.Canvas.Line(fx, fy, tx, ty, "stroke:black;stroke-dasharray:4,2;stroke-width:2px;")
        } else {
            ctx.Canvas.Line(fx, fy, tx, ty, "stroke:black;stroke-width:2px;")
        }

        al.drawArrow(ctx, tx, ty, al.TC > ctx.C)
    }
}

func (al *ActivityLine) renderMessage(ctx DrawContext, tx, ty int) {
    //rect, textPoint := MeasureFontRect(al.style.Font, al.style.FontSize, al.Text, tx, ty, SouthGravity)
    rect := al.textBoxRect.PositionAt(tx, ty, SouthGravity)

    ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "fill:white;stroke:white;")
    al.textBox.Render(ctx.Canvas, tx, ty, SouthGravity)
}

// Draws the arrow head.
func (al *ActivityLine) drawArrow(ctx DrawContext, x, y int, isRight bool) {
    headStyle := arrowHeadStyles[al.style.ArrowHead]

    var xs, ys = make([]int, len(headStyle.Xs)), make([]int, len(headStyle.Ys))
    if len(xs) != len(ys) {
        panic("length of xs and ys must be the same")
    }

    for i := range headStyle.Xs {
        ox, oy := headStyle.Xs[i], headStyle.Ys[i]
        if isRight {
            xs[i] = x + ox
        } else {
            xs[i] = x - ox
        }
        ys[i] = y + oy
    }

    ctx.Canvas.Polyline(xs, ys, headStyle.BaseStyle.ToStyle())
}

// Style information for arrow heads
type arrowHeadStyle struct {
    // Points from the origin
    Xs              []int
    Ys              []int

    // Base style for the arrow head
    BaseStyle       SvgStyle
}

// Styling of the arrow head
var arrowHeadStyles = map[ActivityArrowHead]*arrowHeadStyle {
    SolidArrowHead: &arrowHeadStyle {
        Xs: []int { -9, 0, -9 },
        Ys: []int { -5, 0, 5 },
        BaseStyle: StyleFromString("stroke:black;fill:black;stroke-width:2px;"),
    },
    OpenArrowHead: &arrowHeadStyle {
        Xs: []int { -9, 0, -9 },
        Ys: []int { -5, 0, 5 },
        BaseStyle: StyleFromString("stroke:black;fill:none;stroke-width:2px;"),
    },
    BarbArrowHead: &arrowHeadStyle {
        Xs: []int { -11, 0 },
        Ys: []int { -7, 0 },
        BaseStyle: StyleFromString("stroke:black;fill:black;stroke-width:2px;"),
    },
}