package rastercanvas

import (
	"image"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/lmika/goseq/seqdiagram/canvas"
)

func New(image *image.Image) canvas.Canvas {
	return &rasterCanvas{
		dest: image,
	}
}

type rasterCanvas struct {
	outFile string
	dest    *image.Image
	gc      *draw2dimg.GraphicContext
}

func (rc *rasterCanvas) Line(fx int, fy int, tx int, ty int, stoke canvas.StrokeStyle) {
	rc.gc.MoveTo(float64(fx), float64(fy))
	rc.gc.LineTo(float64(tx), float64(ty))
	rc.gc.Stroke()
}

func (rc *rasterCanvas) Rect(x int, y int, w int, h int, stroke canvas.StrokeStyle, fill canvas.FillStyle) {
	rc.gc.MoveTo(float64(x), float64(y))
	rc.gc.LineTo(float64(x+w), float64(y))
	rc.gc.LineTo(float64(x+w), float64(y+h))
	rc.gc.LineTo(float64(x), float64(y+h))
	rc.gc.Close()
	rc.gc.FillStroke()
}

func (rc *rasterCanvas) Circle(x int, y int, rad int, stroke canvas.StrokeStyle, fill canvas.FillStyle) {
	rc.gc.MoveTo(float64(x+rad), float64(y))
	rc.gc.ArcTo(float64(x), float64(y), float64(rad), float64(rad), 0, 360)
	rc.gc.FillStroke()
}

func (rc *rasterCanvas) Text(left int, bottom int, line string, style canvas.FontStyle) {
	// TODO
}

func (rc *rasterCanvas) Polygon(xs []int, ys []int, closed bool, stroke canvas.StrokeStyle, fill canvas.FillStyle) {
	maxI := len(xs)
	if len(ys) < maxI {
		maxI = len(ys)
	}
	for i := 0; i < maxI; i++ {
		if i == 0 {
			rc.gc.MoveTo(float64(xs[i]), float64(ys[i]))
		} else {
			rc.gc.LineTo(float64(xs[i]), float64(ys[i]))
		}
	}
	if closed {
		rc.gc.Close()
	}
	rc.gc.FillStroke()
}

func (rc *rasterCanvas) Polyline(xs []int, ys []int, stroke canvas.StrokeStyle) {
	maxI := len(xs)
	if len(ys) < maxI {
		maxI = len(ys)
	}
	for i := 0; i < maxI; i++ {
		if i == 0 {
			rc.gc.MoveTo(float64(xs[i]), float64(ys[i]))
		} else {
			rc.gc.LineTo(float64(xs[i]), float64(ys[i]))
		}
	}
	rc.gc.Stroke()
}

func (rc *rasterCanvas) Path(path string, stroke canvas.StrokeStyle, fill canvas.FillStyle) {
	// TODO
}

func (rc *rasterCanvas) SetSize(width int, height int) {
	dest := image.NewRGBA(image.Rect(0, 0, width, height))
	*rc.dest = dest
	rc.gc = draw2dimg.NewGraphicContext(dest)
}

func (rc *rasterCanvas) Close() error {
	//return draw2dimg.SaveToPngFile(rc.outFile, rc.dest)
	return nil
}
