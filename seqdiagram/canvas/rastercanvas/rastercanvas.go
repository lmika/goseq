package rastercanvas

import (
	"image"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/lmika/goseq/seqdiagram/canvas"
	"github.com/llgcode/draw2d"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"fmt"
	"image/color"
	"bitbucket.org/lmika/vecdraw/paths/svgpath"
	"bitbucket.org/lmika/vecdraw/renderer/d2renderer"
	"log"
)

var goSansRegular = draw2d.FontData{
	Name:   "gofont",
	Family: draw2d.FontFamilySans,
	Style:  draw2d.FontStyleNormal,
}

func init() {
	regularTtf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(fmt.Errorf("error parsing builtin ttf regular", err))
	}

	// Register the builtin go fonts
	draw2d.RegisterFont(goSansRegular, regularTtf)
}

type rasterCanvas struct {
	outFile string
	dest    *image.Image
	gc      *draw2dimg.GraphicContext
}

func New(image *image.Image) canvas.Canvas {
	return &rasterCanvas{
		dest: image,
	}
}

func (rc *rasterCanvas) SetSize(width int, height int) {
	dest := image.NewRGBA(image.Rect(0, 0, width, height))
	*rc.dest = dest

	rc.gc = draw2dimg.NewGraphicContext(dest)
	rc.gc.SetFontData(goSansRegular)
}

func (rc *rasterCanvas) setStrokeStyle(stroke canvas.StrokeStyle) {
	if stroke.Color != nil {
		rc.gc.SetStrokeColor(stroke.Color)
	} else {
		rc.gc.SetStrokeColor(color.Transparent)
	}

	if stroke.DashArray != nil {
		fdash := make([]float64, len(stroke.DashArray))
		for i, d := range stroke.DashArray {
			fdash[i] = float64(d)
		}
		rc.gc.SetLineDash(fdash, 0.0)
	} else {
		rc.gc.SetLineDash(nil, 0.0)
	}

	rc.gc.SetLineWidth(stroke.Width)
}

func (rc *rasterCanvas) setFilleStyle(fill canvas.FillStyle) {
	if fill.Color != nil {
		rc.gc.SetFillColor(fill.Color)
	} else {
		rc.gc.SetFillColor(color.Transparent)
	}
}

func (rc *rasterCanvas) Line(fx int, fy int, tx int, ty int, stroke canvas.StrokeStyle) {
	rc.setStrokeStyle(stroke)

	rc.gc.BeginPath()
	rc.gc.MoveTo(float64(fx), float64(fy))
	rc.gc.LineTo(float64(tx), float64(ty))
	rc.gc.Stroke()
}



func (rc *rasterCanvas) Rect(x int, y int, w int, h int, stroke canvas.StrokeStyle, fill canvas.FillStyle) {
	rc.setStrokeStyle(stroke)
	rc.setFilleStyle(fill)

	rc.gc.BeginPath()
	rc.gc.MoveTo(float64(x), float64(y))
	rc.gc.LineTo(float64(x+w), float64(y))
	rc.gc.LineTo(float64(x+w), float64(y+h))
	rc.gc.LineTo(float64(x), float64(y+h))
	rc.gc.Close()
	rc.gc.FillStroke()
}

func (rc *rasterCanvas) Circle(x int, y int, rad int, stroke canvas.StrokeStyle, fill canvas.FillStyle) {
	rc.setStrokeStyle(stroke)
	rc.setFilleStyle(fill)

	rc.gc.BeginPath()
	rc.gc.MoveTo(float64(x+rad), float64(y))
	rc.gc.ArcTo(float64(x), float64(y), float64(rad), float64(rad), 0, 360)
	rc.gc.FillStroke()
}

func (rc *rasterCanvas) Text(left int, bottom int, line string, style canvas.FontStyle) {
	// TEMP
	rc.gc.SetStrokeColor(style.Color)
	rc.gc.SetFillColor(style.Color)
	rc.gc.SetFontSize(style.Size)

	rc.gc.FillStringAt(line, float64(left), float64(bottom))
	// TODO
}

func (rc *rasterCanvas) Polygon(xs []int, ys []int, closed bool, stroke canvas.StrokeStyle, fill canvas.FillStyle) {
	rc.setStrokeStyle(stroke)
	rc.setFilleStyle(fill)

	maxI := len(xs)
	if len(ys) < maxI {
		maxI = len(ys)
	}

	rc.gc.BeginPath()
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
	rc.setStrokeStyle(stroke)

	maxI := len(xs)
	if len(ys) < maxI {
		maxI = len(ys)
	}

	rc.gc.BeginPath()
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
	p, err := svgpath.Parse(path)
	if err != nil {
		log.Printf("Bad path: [%s]: %v", path, err)
		return
	}

	rc.setStrokeStyle(stroke)
	rc.setFilleStyle(fill)

	d2d := &d2renderer.Draw2DRenderer{ GC: rc.gc }
	d2d.FillStroke(p)
}

func (rc *rasterCanvas) Close() error {
	return nil
}
