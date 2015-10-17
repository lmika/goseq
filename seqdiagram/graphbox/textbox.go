package graphbox

import (
	"fmt"
	"strings"

	"github.com/ajstarks/svgo"
)

const (
	LINE_GAP = 2
)

type TextAlign int

const (
	LeftTextAlign   TextAlign = iota
	MiddleTextAlign           = iota
	RightTextAlign            = iota
)

// A block of prose
type TextBox struct {
	Lines    []string
	Font     Font
	FontSize int
	Align    TextAlign

	Color string
}

// Returns a new text box
func NewTextBox(font Font, fontSize int, align TextAlign) *TextBox {
	return &TextBox{
		Lines:    make([]string, 0),
		Font:     font,
		FontSize: fontSize,
		Align:    align,
	}
}

// Adds some text
func (tb *TextBox) AddText(text string) {
	tb.Lines = append(tb.Lines, strings.Split(text, "\n")...)
}

// Returns the width and height of the text box.
func (tb *TextBox) Measure() (int, int) {
	w := 0
	h := 0

	for _, line := range tb.Lines {
		lw, lh := tb.measureLine(line)
		w = maxInt(w, lw)
		h += lh + LINE_GAP
	}

	return w, h - LINE_GAP
}

// Measures a line
func (tb *TextBox) measureLine(line string) (int, int) {
	fs := float64(tb.FontSize)
	return tb.Font.Measure(line, fs)
}

// Given a font, font size, points and gravity, returns a rectangle which will contain
// the text centered.  The point and gravity describes the location of the rect.
// The second point is where the text is to start given that it is to be rendered to
// fill the rectangle with default anchoring and alignment
func (tb *TextBox) BoundingRect() Rect {
	w, h := tb.Measure()
	//ox, oy := gravity(w, h)

	return Rect{0, 0, w, h}
}

// Renders the text from the given point and gravity
func (tb *TextBox) Render(svg *svg.SVG, x, y int, gravity Gravity) {
	rect := tb.BoundingRect().PositionAt(x, y, gravity)
	left := rect.X
	currY := rect.Y
	style := tb.textStyle()

	for _, line := range tb.Lines {
		var textLeft int

		lineW, lineH := tb.measureLine(line)

		switch tb.Align {
		case LeftTextAlign:
			textLeft = left
		case MiddleTextAlign:
			textLeft = left + (rect.W / 2) - (lineW / 2)
		case RightTextAlign:
			textLeft = left + rect.W - lineW
		}

		textBottom := currY + lineH - (tb.FontSize*1/4 - 1)

		if line != "" {
			svg.Text(textLeft, textBottom, line, style)
		}

		currY += lineH + LINE_GAP
	}
}

// Returns the text styling
func (tb *TextBox) textStyle() string {
	s := SvgStyle{}

	s.Set("font-family", tb.Font.SvgName())
	s.Set("font-size", fmt.Sprintf("%dpx", tb.FontSize))

	if tb.Color != "" {
		s.Set("fill", tb.Color)
	}

	return s.ToStyle()
}
