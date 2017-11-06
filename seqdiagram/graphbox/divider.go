package graphbox

import (
	"image/color"

	"github.com/lmika/goseq/seqdiagram/canvas"
)

// DividerShape determines which shape to use for the divider
type DividerShape int

const (
	// DSFullRect is a rectangle which will span the entire graphic from end to end.
	// The text will be centered in front of it.
	DSFullRect DividerShape = iota

	// DSFramedRect is like FullRect but using a framed rectangle
	DSFramedRect

	// DSSpacerRect is like FullRect but "transparent".  If there is any text, it will be
	// blocked out.
	DSSpacerRect

	// DSFullLine is a line which will span the entire grapic.  The text will be
	// centered in front of it.
	DSFullLine
)

// DividerStyle defines the style of the divider
type DividerStyle struct {
	Font        Font
	FontSize    int
	Padding     Point
	Margin      Point
	TextPadding Point
	Overlap     int
	Shape       DividerShape
}

// Divider is a divider graphics object.  This spans the entire diagram.
type Divider struct {
	TC int

	leftOverlap  int
	rightOverlap int

	style       DividerStyle
	hasText     bool
	textBox     *TextBox
	textBoxRect Rect
	marginRect  Rect
}

// NewDivider creates a new divider
func NewDivider(toCol int, text string, style DividerStyle) *Divider {
	textBox := NewTextBox(style.Font, style.FontSize, MiddleTextAlign)
	textBox.AddText(text)
	textBoxRect := textBox.BoundingRect()
	marginRect := textBoxRect.BlowOut(style.Padding)

	return &Divider{toCol, 0, 0, style, text != "", textBox, textBoxRect, marginRect}
}

// Constraint returns the constraints of the graphics object
func (div *Divider) Constraint(r, c int, applier ConstraintApplier) {
	// There must be enought horizontal space to accommodate the text
	// and vertical space to display the divider
	requiredHeight := div.marginRect.H + div.style.Margin.Y*2
	requiredWidth := div.marginRect.W + div.style.Margin.X*2

	applier.Apply(AddSizeConstraint{r, c, 0, 0, requiredHeight / 2, requiredHeight / 2})

	if div.style.Overlap > 0 {
		div.leftOverlap = div.style.Overlap
		div.rightOverlap = div.style.Overlap

		if c == 0 {
			div.leftOverlap = 0
		}
		if div.TC == applier.Cols()-1 {
			div.rightOverlap = 0
		}

		applier.Apply(SizeConstraint{r, c, div.leftOverlap, 0, 0, 0})
		applier.Apply(SizeConstraint{r, div.TC, 0, div.rightOverlap, 0, 0})
		applier.Apply(TotalSizeConstraint{r - 1, c, r, div.TC, requiredWidth - (div.leftOverlap + div.rightOverlap), 0})
	} else {
		applier.Apply(TotalSizeConstraint{r - 1, c, r, div.TC, requiredWidth, 0})
	}
}

// Draw draws the graphics object
func (div *Divider) Draw(ctx DrawContext, point Point) {
	fx, fy := point.X, point.Y
	if point, isPoint := ctx.PointAt(ctx.R, div.TC); isPoint {
		fx -= div.leftOverlap
		tx, _ := point.X+div.rightOverlap, point.Y

		centerX := fx + (tx-fx)/2
		centerY := fy

		borderRect := Rect{fx, fy - div.marginRect.H/2, tx - fx, div.marginRect.H}
		textBoxRect := div.textBoxRect.PositionAt(centerX, centerY, CenterGravity).BlowOut(div.style.TextPadding)

		// Draw the shape and text
		switch div.style.Shape {
		case DSFullRect:
			//ctx.Canvas.Rect(borderRect.X, borderRect.Y, borderRect.W, borderRect.H, "fill:white;stroke:white;")
			ctx.Canvas.Rect(borderRect.X, borderRect.Y, borderRect.W, borderRect.H, canvas.WhiteStroke, canvas.WhiteFill)
			div.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
		case DSFramedRect:
			//ctx.Canvas.Rect(borderRect.X, borderRect.Y, borderRect.W, borderRect.H, "fill:white;stroke:black;stroke-width:2px")
			ctx.Canvas.Rect(borderRect.X, borderRect.Y, borderRect.W, borderRect.H, canvas.StrokeStyle{Color: color.Black, Width: 2}, canvas.WhiteFill)
			div.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
		case DSSpacerRect:
			//ctx.Canvas.Rect(textBoxRect.X, textBoxRect.Y, textBoxRect.W, textBoxRect.H, "fill:white;stroke:white;")
			ctx.Canvas.Rect(textBoxRect.X, textBoxRect.Y, textBoxRect.W, textBoxRect.H, canvas.WhiteStroke, canvas.WhiteFill)
			div.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
		case DSFullLine:
			// Draw the rectangle for clearing the image
			/*
				ctx.Canvas.Rect(borderRect.X, borderRect.Y, borderRect.W, borderRect.H, "fill:white;stroke:white;")
				ctx.Canvas.Line(borderRect.X, centerY, borderRect.W, centerY, "fill:white;stroke:black;stroke-width:2px;") //stroke-dasharray:16,8")
			*/
			ctx.Canvas.Rect(borderRect.X, borderRect.Y, borderRect.W, borderRect.H, canvas.WhiteStroke, canvas.WhiteFill)
			ctx.Canvas.Line(borderRect.X, centerY, borderRect.W, centerY, canvas.StrokeStyle{Color: color.Black, Width: 2}) //stroke-dasharray:16,8")

			if div.hasText {
				//ctx.Canvas.Rect(textBoxRect.X, textBoxRect.Y, textBoxRect.W, textBoxRect.H, "fill:white;stroke:white;")
				ctx.Canvas.Rect(textBoxRect.X, textBoxRect.Y, textBoxRect.W, textBoxRect.H, canvas.WhiteStroke, canvas.WhiteFill)
				div.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
			}
		}
	}
}
