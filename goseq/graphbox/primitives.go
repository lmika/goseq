package graphbox

import (
    "fmt"
)

type ActorBoxPos int
const (
    TopActorBox       ActorBoxPos     =   iota
    BottomActorBox                    =   iota
)

// Styling options for the actor rect
type ActorBoxStyle struct {
    Font        Font
    FontSize    int
    Padding     Point
}

// Draws an object instance
type ActorBox struct {
    frameRect   Rect
    style       ActorBoxStyle
    textBox     *TextBox
    pos         ActorBoxPos
}

func NewActorBox(text string, style ActorBoxStyle, pos ActorBoxPos) *ActorBox {
    var textAlign TextAlign = MiddleTextAlign

    textBox := NewTextBox(style.Font, style.FontSize, textAlign)
    textBox.AddText(text)

    trect := textBox.BoundingRect()
    brect := trect.BlowOut(style.Padding)

    return &ActorBox{brect, style, textBox, pos}
}

func (tr *ActorBox) Constraint(r, c int) Constraint {
    var vertConstraint Constraint

    if (tr.pos == TopActorBox) {
        vertConstraint = SizeConstraint{r, c, 0, 0, tr.frameRect.H / 2, tr.frameRect.H / 2 + 4}
    } else {
        vertConstraint = SizeConstraint{r, c, 0, 0, tr.frameRect.H / 2 + 8, tr.frameRect.H / 2}
    }

    if (tr.pos == TopActorBox) {
        if (c == 0) {
            return Constraints([]Constraint {
                vertConstraint,
                SizeConstraint{r, c, tr.frameRect.W / 2, 0, 0, 0},
                AddSizeConstraint{r, c, 0, tr.frameRect.W / 2, 0, 0},
            })
        } else {
            return Constraints([]Constraint {
                vertConstraint,
                AddSizeConstraint{r, c, tr.frameRect.W / 2, 0, 0, 0},
                AddSizeConstraint{r, c, 0, tr.frameRect.W / 2, 0, 0},
            })
        }
    } else {
        return vertConstraint
    }
}

func (r *ActorBox) Draw(ctx DrawContext, point Point) {
    centerX, centerY := point.X, point.Y

    rect := r.frameRect.PositionAt(centerX, centerY, CenterGravity)
    ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "stroke:black;fill:white")
    r.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
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
    h := al.textBoxRect.H + al.style.PaddingTop + al.style.TextGap
    w := al.textBoxRect.W

    _ = h
    _ = w
    return Constraints([]Constraint{ 
        AddSizeConstraint{r, c, 0, 0, h, al.style.PaddingBottom},
        TotalSizeConstraint{r - 1, c, r, al.TC, w + 32, 0},
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