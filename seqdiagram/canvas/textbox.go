package canvas

import (
	"image"
)

// TextBox contains information about a block of text.  The text block
// is divided into "runs", with each run containing attributes.  This,
// mixed with the font information, will combine to produce text-layout
// information.
type TextBox struct {
	runs []textRun
}

// textRun is a run of text with some attributes
type textRun struct {
	text  string
	attrs int // TODO.  Attributes include "bold", "italic", "underline", "fixed".  Could include colour?
}

// Layout lays out the font without performing any wrapping.
func (tb *TextBox) Layout(fontInfo FontInfo) TextLayout {
	return TextLayout{}
}

// LayoutWrap lays out the contents of the text box such that no single
// line is larger than maxWidth in points.
func (tb *TextBox) LayoutWrap(fontInfo FontInfo, maxWidth float64) TextLayout {
	return TextLayout{}
}

// TextLayout contains layout information about the font.  This includes
// the individual lines of font and bounding information.
type TextLayout struct {
}

// Bounds returns the bounding rect that will contain all the text.
func (tl TextLayout) Bounds() image.Rectangle {
	return image.Rectangle{}
}

// FontInfo contains information about the font used for
type FontInfo struct {
}
