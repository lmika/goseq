package graphbox

import (
	"image/color"

	"github.com/lmika/goseq/seqdiagram/canvas"
)

// ActorIconBoxStyle defines styling options for an actor icon
type ActorIconBoxStyle struct {
	Font      Font
	FontSize  int
	Padding   Point
	Margin    Point
	IconGap   int
	Color     string
	TextColor string
}

// ActorIconBox represents an actor icon
type ActorIconBox struct {
	//Caption     string
	textBox *TextBox
	Icon    Icon
	style   ActorIconBoxStyle
	pos     ActorBoxPos
}

// NewActorIconBox constructs a new actor icon
func NewActorIconBox(text string, icon Icon, style ActorIconBoxStyle, pos ActorBoxPos) *ActorIconBox {
	textBox := NewTextBox(style.Font, style.FontSize, MiddleTextAlign)
	textBox.Color = style.TextColor
	textBox.AddText(text)

	return &ActorIconBox{textBox, icon, style, pos}
}

func (tr *ActorIconBox) Constraint(r, c int, applier ConstraintApplier) {
	posHoriz, posVert := tr.pos&0xFF00, tr.pos&0xFF
	iconW, iconH := tr.Icon.Size()

	brect := tr.textBox.BoundingRect()
	w := maxInt(iconW, brect.W) + tr.style.Padding.X

	topH := iconH / 2
	bottomH := iconH/2 + brect.H + tr.style.IconGap + tr.style.Padding.Y
	marginX := tr.style.Margin.X

	if posVert == TopActorBox {
		if posHoriz == LeftActorBox {
			applier.Apply(SizeConstraint{r, c, w / 2, marginX / 2, 0, 0})
			applier.Apply(AddSizeConstraint{r, c, 0, w / 2, 0, 0})
		} else if posHoriz == RightActorBox {
			applier.Apply(SizeConstraint{r, c, marginX / 2, w / 2, 0, 0})
			applier.Apply(AddSizeConstraint{r, c, w / 2, 0, 0, 0})
		} else {
			applier.Apply(SizeConstraint{r, c, marginX / 2, marginX / 2, 0, 0})
			applier.Apply(AddSizeConstraint{r, c, w / 2, w / 2, 0, 0})
		}
		applier.Apply(SizeConstraint{r, c, 0, 0, topH, bottomH})
	}
}

func (tr *ActorIconBox) Draw(ctx DrawContext, point Point) {
	centerX, centerY := point.X, point.Y

	iconW, iconH := tr.Icon.Size()
	iconX, iconY := centerX, centerY

	// Draw the text
	brect := tr.textBox.BoundingRect()
	textY := iconY + iconH/2 + tr.style.IconGap
	rect := brect.PositionAt(centerX, textY, NorthGravity)

	// Draw the icon
	iconStyle := SvgStyle{}
	iconStyle.Set("stroke", tr.style.Color)
	iconStyle.Set("fill", "white")
	iconStyle.Set("stroke-width", "2px")

	//ctx.Canvas.Rect(rect.X, rect.Y-tr.style.IconGap, rect.W, rect.H+tr.style.IconGap, "stroke:white;fill:white;stroke-width:2px;")
	ctx.Canvas.Rect(rect.X, rect.Y-tr.style.IconGap, rect.W, rect.H+tr.style.IconGap,
		canvas.StrokeStyle{Color: color.White}, canvas.FillStyle{Color: color.White})
	tr.textBox.Render(ctx.Canvas, centerX, textY, NorthGravity)

	// ctx.Canvas.Rect(centerX-iconW/2, centerY-iconH/2, iconW, iconH, "stroke:white;fill:white;stroke-width:1px;")
	ctx.Canvas.Rect(centerX-iconW/2, centerY-iconH/2, iconW, iconH, canvas.StrokeStyle{Color: color.White}, canvas.FillStyle{Color: color.White})
	tr.Icon.Draw(ctx, iconX, iconY, &iconStyle)
}
