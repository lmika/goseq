package graphbox

import (
	"fmt"
	"image/color"

	"github.com/lmika/goseq/seqdiagram/canvas"
)

// ActivityArrowStem is the type of arrow stem to use for activity arrows
type ActivityArrowStem int

const (
	// SolidArrowStem draws a solid arrow stem
	SolidArrowStem ActivityArrowStem = iota

	// DashedArrowStem draws a dashed arrow stem
	DashedArrowStem = iota

	// ThickArrowStem draws a thick arrow stem
	ThickArrowStem = iota
)

// ActivityLineStyle defines the style to use for an activity line
type ActivityLineStyle struct {
	Font          Font
	FontSize      int
	Margin        Point
	TextGap       int
	SelfRefWidth  int
	SelfRefHeight int
	//ArrowHead       ActivityArrowHead
	ArrowHead *ArrowHeadStyle
	ArrowStem ActivityArrowStem
}

// Returns the text style
func (as ActivityLineStyle) textStyle() string {
	s := SvgStyle{}

	s.Set("font-family", as.Font.SvgName())
	s.Set("font-size", fmt.Sprintf("%dpx", as.FontSize))

	return s.ToStyle()
}

// ActivityLine is an activity line graphical object
type ActivityLine struct {
	TC          int
	style       ActivityLineStyle
	textBox     *TextBox
	textBoxRect Rect
}

// NewActivityLine constructs a new ActivityLine
func NewActivityLine(toCol int, selfRef bool, text string, style ActivityLineStyle) *ActivityLine {
	var textBoxAlign TextAlign = MiddleTextAlign
	if selfRef {
		textBoxAlign = LeftTextAlign
	}

	textBox := NewTextBox(style.Font, style.FontSize, textBoxAlign)
	textBox.AddText(text)

	brect := textBox.BoundingRect()
	return &ActivityLine{toCol, style, textBox, brect}
}

// Constraint returns the constraints of the graphics object
func (al *ActivityLine) Constraint(r, c int, applier ConstraintApplier) {
	h := al.textBoxRect.H + al.style.Margin.Y + al.style.TextGap
	w := al.textBoxRect.W

	lc, rc := c, al.TC
	if al.TC < c {
		lc, rc = al.TC, c
	}

	if al.TC == c {
		// An arrow referring to itself
		w = maxInt(w, al.style.SelfRefWidth) + al.style.TextGap*3
		h += al.style.TextGap / 2

		applier.Apply(AddSizeConstraint{r, c, 0, 0, h, al.style.Margin.Y + al.style.SelfRefHeight})
		applier.Apply(TotalSizeConstraint{r - 1, lc, r, lc + 1, w, 0})
	} else {
		applier.Apply(AddSizeConstraint{r, c, 0, 0, h, al.style.Margin.Y})
		applier.Apply(TotalSizeConstraint{r - 1, lc, r, rc, w + al.style.Margin.X*2, 0})
	}
}

// Draw draws the graphics object
func (al *ActivityLine) Draw(ctx DrawContext, point Point) {
	fx, fy := point.X, point.Y

	if ctx.C == al.TC {
		// A self reference arrow
		if point, isPoint := ctx.PointAt(ctx.R, ctx.C+1); isPoint {
			// Draw an arrow referencing itself
			ty := point.Y
			stemX, stemY := fx+al.style.SelfRefWidth, ty+al.style.SelfRefHeight

			textX := fx + al.style.TextGap*2
			textY := ty - al.style.TextGap - al.style.TextGap/2
			al.renderMessage(ctx, textX, textY, true)

			al.drawArrowStemPath(ctx,
				[]int{fx, stemX, stemX, fx},
				[]int{fy, fy, stemY, stemY})
			al.drawArrow(ctx, fx, stemY, false)
		}
	} else {

		if point, isPoint := ctx.PointAt(ctx.R, al.TC); isPoint {
			tx, ty := point.X, point.Y

			textX := fx + (tx-fx)/2
			textY := ty - al.style.TextGap
			al.renderMessage(ctx, textX, textY, false)
			al.drawArrowStem(ctx, fx, fy, tx, ty)
			al.drawArrow(ctx, tx, ty, al.TC > ctx.C)
		}
	}
}

// Draws the arrow stem
func (al *ActivityLine) drawArrowStem(ctx DrawContext, fx, fy, tx, ty int) {
	switch al.style.ArrowStem {
	case SolidArrowStem:
		//ctx.Canvas.Line(fx, fy, tx, ty, "stroke:black;stroke-width:2px;")
		ctx.Canvas.Line(fx, fy, tx, ty, canvas.StrokeStyle{Color: color.Black, Width: 2.0})
	case DashedArrowStem:
		//ctx.Canvas.Line(fx, fy, tx, ty, "stroke:black;stroke-dasharray:4,2;stroke-width:2px;")
		ctx.Canvas.Line(fx, fy, tx, ty, canvas.StrokeStyle{Color: color.Black, Width: 2.0, DashArray: []int{4, 2}})
	case ThickArrowStem:
		//ctx.Canvas.Line(fx, fy, tx, ty, "stroke:black;stroke-width:4px;")
		ctx.Canvas.Line(fx, fy, tx, ty, canvas.StrokeStyle{Color: color.Black, Width: 4.0})
	}
}

// Draws the arrow stem path
func (al *ActivityLine) drawArrowStemPath(ctx DrawContext, xs, ys []int) {
	switch al.style.ArrowStem {
	case SolidArrowStem:
		//ctx.Canvas.Polyline(xs, ys, "fill:none;stroke:black;stroke-width:2px;")
		ctx.Canvas.Polyline(xs, ys, canvas.StrokeStyle{Color: color.Black, Width: 2.0})
	case DashedArrowStem:
		//ctx.Canvas.Polyline(xs, ys, "fill:none;stroke:black;stroke-dasharray:4,2;stroke-width:2px;")
		ctx.Canvas.Polyline(xs, ys, canvas.StrokeStyle{Color: color.Black, Width: 2.0, DashArray: []int{4, 2}})
	case ThickArrowStem:
		//ctx.Canvas.Polyline(xs, ys, "fill:none;stroke:black;stroke-width:4px;")
		ctx.Canvas.Polyline(xs, ys, canvas.StrokeStyle{Color: color.Black, Width: 4.0})
	}
}

func (al *ActivityLine) renderMessage(ctx DrawContext, tx, ty int, anchorLeft bool) {
	//rect, textPoint := MeasureFontRect(al.style.Font, al.style.FontSize, al.Text, tx, ty, SouthGravity)
	anchor := SouthGravity
	if anchorLeft {
		anchor = SouthWestGravity
	}

	rect := al.textBoxRect.PositionAt(tx, ty, anchor)

	ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, canvas.StrokeStyle{}, canvas.FillStyle{})
	al.textBox.Render(ctx.Canvas, tx, ty, anchor)
	/*
		ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, "fill:white;stroke:white;")
		al.textBox.Render(ctx.Canvas, tx, ty, anchor)
	*/
}

// Draws the arrow head.
func (al *ActivityLine) drawArrow(ctx DrawContext, x, y int, isRight bool) {
	headStyle := al.style.ArrowHead

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

	//ctx.Canvas.Polyline(xs, ys, StyleFromString(headStyle.BaseStyle).ToStyle())
	// TODO: StyleFromString
	ctx.Canvas.Polyline(xs, ys, canvas.StrokeStyle{Color: color.Black})
}

// ArrowHeadStyle defines style information for the arrow heads
type ArrowHeadStyle struct {
	// Points from the origin
	Xs []int
	Ys []int

	// Base style for the arrow head
	BaseStyle string
}
