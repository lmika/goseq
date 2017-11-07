package canvas

import (
	"image/color"
	"io"
)

// Canvas is responsible for rendering the image primitives.
// The actual primitives are moderately high-level in order to
// simplify the logic within graphbox.
type Canvas interface {
	io.Closer

	// Line draws a line between two points with a given stroke style
	Line(fx, fy, tx, ty int, stoke StrokeStyle)

	// Rect draws a rectangle with a given stroke and fill style
	Rect(x, y, w, h int, stroke StrokeStyle, fill FillStyle)

	// Circle draws a circle centered at X, Y with the given radius of rad
	Circle(x, y, rad int, stroke StrokeStyle, fill FillStyle)

	// Textbox draws a text-box with font information
	//Textbox() // TODO

	// Text renders a line of text at the given coordinates using the given text style.
	// Temporary until text-boxes have been reimplemented.
	Text(left, bottom int, line string, style FontStyle)

	// Icon draws an "icon" (TODO)
	//Icon() // TODO

	// Polygon draws a polygon, either closed or opened, with a given stroke and fill stype
	Polygon(xs, ys []int, closed bool, stroke StrokeStyle, fill FillStyle)

	// Polyline draws a multi-segmented line with a given stroke style
	Polyline(xs, ys []int, stroke StrokeStyle)

	// Path draws a path defined in SVG's path specification
	Path(path string, stroke StrokeStyle, fill FillStyle)

	// SetSize sets the size of the canvas.  This is to be called before
	// any of the other methods are to be called
	SetSize(width, height int)

	// Close closes the canvas, indicating that drawing is done.
	//Close()
}

// StrokeStyle are the styles for the stroke properties
type StrokeStyle struct {
	// Color is the colour of the stroke.  Use color.Transparent for no stroke
	Color     color.Color
	Width     float64
	DashArray []int
}

// FillStyle contains fill style properties
type FillStyle struct {
	// Color is the fill colour.  Use color.Transparent for no fill
	Color color.Color
}

// FontStyle contains base style information about text properties
type FontStyle struct {
	// Family is the name of the font-family to use.
	Family string

	// Color is the fill colour of the font.
	Color color.Color

	// Size is the size of the font, in points
	Size float64
}

// Predefined strokes and fill styles
var (
	NoStroke    = StrokeStyle{}
	WhiteStroke = StrokeStyle{Color: color.White}

	NoFill    = FillStyle{}
	WhiteFill = FillStyle{Color: color.White}
)
