package canvas

import "image/color"

// Canvas is responsible for rendering the image primitives.
// The actual primitives are moderately high-level in order to
// simplify the logic within graphbox.
type Canvas interface {
	// Line draws a line between two points with a given stroke style
	Line(fx, fy, tx, ty float64, stoke StrokeStyle)

	// Rect draws a rectangle with a given stroke and fill style
	Rect(x, y, w, h float64, stroke StrokeStyle, fill FillStyle)

	// Textbox draws a text-box with font information
	Textbox() // TODO

	// Icon draws an "icon" (TODO)
	Icon() // TODO

	// Polygon draws a polygon with a given stroke and fill stype
	Polygon()

	// Polyline draws a polyline with a given stroke style
	Polyline()

	// SvgPath draws a path defined in SVG's path specification
	SvgPath()
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

type FontStyle struct {
	// TODO
}
