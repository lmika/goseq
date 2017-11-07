package svgcanvas

import (
	"fmt"
	"image/color"
	"io"

	svg "github.com/ajstarks/svgo"
	"github.com/lmika/goseq/seqdiagram/canvas"
)

// New creates a new canvas which will write to an SVG
func New(w io.Writer) canvas.Canvas {
	return &svgCanvas{svg.New(w)}
}

type svgCanvas struct {
	svg *svg.SVG
}

func (c *svgCanvas) Line(fx int, fy int, tx int, ty int, stroke canvas.StrokeStyle) {
	c.svg.Line(fx, fy, tx, ty, c.strokeStyleToString(stroke))
}

func (c *svgCanvas) Rect(x int, y int, w int, h int, stroke canvas.StrokeStyle, fill canvas.FillStyle) {
	c.svg.Rect(x, y, w, h, c.strokeStyleToString(stroke)+";"+c.fillStyleToString(fill))
}

func (c *svgCanvas) Circle(x int, y int, rad int, stroke canvas.StrokeStyle, fill canvas.FillStyle) {
	c.svg.Circle(x, y, rad, c.strokeStyleToString(stroke)+";"+c.fillStyleToString(fill))
}

func (c *svgCanvas) Text(left int, bottom int, line string, style canvas.FontStyle) {
	c.svg.Text(left, bottom, line, fmt.Sprintf("font-family:%s;font-size:%fpx;fill:%s",
		style.Family, style.Size, c.colorToString(style.Color)))
}

func (c *svgCanvas) Polygon(xs []int, ys []int, closed bool, stroke canvas.StrokeStyle, fill canvas.FillStyle) {
	if closed {
		c.svg.Polygon(xs, ys, c.strokeStyleToString(stroke)+";"+c.fillStyleToString(fill))
	} else {
		c.svg.Polyline(xs, ys, c.strokeStyleToString(stroke)+";"+c.fillStyleToString(fill))
	}
}

func (c *svgCanvas) Polyline(xs []int, ys []int, stroke canvas.StrokeStyle) {
	c.svg.Polyline(xs, ys, c.strokeStyleToString(stroke)+";fill:transparent")
}

func (c *svgCanvas) Path(path string, stroke canvas.StrokeStyle, fill canvas.FillStyle) {
	c.svg.Path(path, c.strokeStyleToString(stroke)+";"+c.fillStyleToString(fill))
}

func (c *svgCanvas) SetSize(width int, height int) {
	c.svg.Start(width, height)

	// Add styles
	c.svg.Def()
	c.addStyles()
	c.svg.DefEnd()
}

func (c *svgCanvas) Close() error {
	c.svg.End()
	return nil
}

func (c *svgCanvas) addStyles() {
	fmt.Fprintln(c.svg.Writer, "<style>")

	// !!TEMP!!
	fmt.Fprintln(c.svg.Writer, "@font-face {")
	fmt.Fprintln(c.svg.Writer, "  font-family: 'DejaVuSans';")
	fmt.Fprintln(c.svg.Writer, "  src: url('https://fontlibrary.org/assets/fonts/dejavu-sans/f5ec8426554a3a67ebcdd39f9c3fee83/49c0f03ec2fa354df7002bcb6331e106/DejaVuSansBook.ttf') format('truetype');")
	fmt.Fprintln(c.svg.Writer, "  font-weight: normal;")
	fmt.Fprintln(c.svg.Writer, "  font-style: normal;")
	fmt.Fprintln(c.svg.Writer, "}")
	// !!END TEMP!!

	fmt.Fprintln(c.svg.Writer, "</style>")
}

func (c *svgCanvas) strokeStyleToString(ss canvas.StrokeStyle) string {
	// TODO: Use string buffers as they're more efficient
	s := "stroke:transparent"
	if ss.Color != nil {
		s = fmt.Sprintf("stroke:%s", c.colorToString(ss.Color))
	}
	if ss.Width > 0.0 {
		s += fmt.Sprintf(";stroke-width:%fpx", ss.Width)
	}
	if len(ss.DashArray) > 0 {
		var da string = ""
		for i, a := range ss.DashArray {
			if i > 0 {
				da += ","
			}
			da += fmt.Sprint(a)
		}
		s += ";stroke-dasharray:" + da
	}

	return s
}

func (c *svgCanvas) fillStyleToString(fs canvas.FillStyle) string {
	if fs.Color == nil {
		return fmt.Sprintf("fill:transparent")
	}

	return fmt.Sprintf("fill:%s", c.colorToString(fs.Color))
}

func (c *svgCanvas) colorToString(col color.Color) string {
	r, g, b, a := col.RGBA()
	return fmt.Sprintf("rgba(%d,%d,%d,%f)", r>>8, g>>8, b>>8, float64(a)/float64(0xFFFF))
}
