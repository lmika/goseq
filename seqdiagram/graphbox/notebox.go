package graphbox

import (
	"image/color"

	"github.com/lmika/goseq/seqdiagram/canvas"
)

type NoteBoxPos int

const (
	CenterNotePos NoteBoxPos = iota
	LeftNotePos              = iota
	RightNotePos             = iota
)

// Styling options for the actor rect
type NoteBoxStyle struct {
	Font     Font
	FontSize int
	Padding  Point
	Margin   Point
	Position NoteBoxPos
}

// Draws an object instance
type NoteBox struct {
	frameRect Rect
	style     NoteBoxStyle
	textBox   *TextBox
	pos       NoteBoxPos
}

func NewNoteBox(text string, style NoteBoxStyle, pos NoteBoxPos) *NoteBox {
	var textAlign TextAlign = MiddleTextAlign

	textBox := NewTextBox(style.Font, style.FontSize, textAlign)
	textBox.AddText(text)

	trect := textBox.BoundingRect()
	brect := trect.BlowOut(style.Padding)

	return &NoteBox{brect, style, textBox, pos}
}

func (tr *NoteBox) Constraint(r, c int, applier ConstraintApplier) {
	var horizConstraint Constraint

	marginX := tr.style.Margin.X
	marginY := tr.style.Margin.Y
	if tr.pos == LeftNotePos {
		horizConstraint = SizeConstraint{r, c, tr.frameRect.W + marginX*2, marginX, 0, 0}
	} else if tr.pos == RightNotePos {
		horizConstraint = SizeConstraint{r, c, marginX, tr.frameRect.W + marginX*2, 0, 0}
	} else {
		horizConstraint = SizeConstraint{r, c, tr.frameRect.W/2 + marginX, tr.frameRect.W/2 + marginX, 0, 0}
	}

	applier.Apply(horizConstraint)
	applier.Apply(AddSizeConstraint{r, c, 0, 0, tr.frameRect.H/2 + marginY, tr.frameRect.H/2 + marginY})
}

func (r *NoteBox) Draw(ctx DrawContext, point Point) {
	centerX, centerY := point.X, point.Y
	marginX := r.style.Margin.X

	strokeStyle := canvas.StrokeStyle{Color: color.Black, Width: 2.0}
	fillStyle := canvas.FillStyle{Color: color.White}

	if r.pos == CenterNotePos {
		rect := r.frameRect.PositionAt(centerX, centerY, CenterGravity)
		ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, strokeStyle, fillStyle)
		r.textBox.Render(ctx.Canvas, centerX, centerY, CenterGravity)
	} else if r.pos == LeftNotePos {
		offsetX := centerX - marginX
		textOffsetX := centerX - r.style.Padding.X - marginX
		rect := r.frameRect.PositionAt(offsetX, centerY, EastGravity)
		ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, strokeStyle, fillStyle)
		r.textBox.Render(ctx.Canvas, textOffsetX, centerY, EastGravity)
	} else if r.pos == RightNotePos {
		offsetX := centerX + marginX
		textOffsetX := centerX + r.style.Padding.X + marginX
		rect := r.frameRect.PositionAt(offsetX, centerY, WestGravity)
		ctx.Canvas.Rect(rect.X, rect.Y, rect.W, rect.H, strokeStyle, fillStyle)
		r.textBox.Render(ctx.Canvas, textOffsetX, centerY, WestGravity)
	}
}
