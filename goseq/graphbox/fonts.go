// Handle fonts

package graphbox

import (
    "os"
    "bytes"
    "math"
    "path/filepath"
    "image"
    "image/color"
    "strings"

    "code.google.com/p/freetype-go/freetype"
    "code.google.com/p/freetype-go/freetype/raster"
    "code.google.com/p/freetype-go/freetype/truetype"
)

type Font interface {
    // Returns the appropriate name of this font in the SVG
    SvgName() string

    // Measures the size of the particular line of text
    Measure(txt string, size float64) (int, int)
}


// A true-type font
type TTFFont struct {
    font        *truetype.Font
    fontName    string
}

// Returns a new TTFFont struct
func NewTTFFont(path string) (*TTFFont, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    buffer := &bytes.Buffer{}
    _, err = buffer.ReadFrom(file)
    if err != nil {
        return nil, err
    }

    ttfFont, err := freetype.ParseFont(buffer.Bytes())
    if err != nil {
        return nil, err
    }

    return &TTFFont{ttfFont, strings.TrimSuffix(filepath.Base(path), ".ttf")}, nil
}

// Measures the size of a font
func (ttf *TTFFont) Measure(txt string, size float64) (int, int) {
    img := nopDrawImage(0)

    ctx := freetype.NewContext()
    ctx.SetDPI(72)
    ctx.SetClip(img.Bounds())
    ctx.SetSrc(img)
    ctx.SetDst(img)
    ctx.SetFont(ttf.font)
    ctx.SetHinting(freetype.NoHinting)
    ctx.SetFontSize(size)

    np, _ := ctx.DrawString(txt, raster.Point{0.0, 0.0})

    return int(np.X >> 8), int(size) + int(np.Y >> 8)
}

// Return the SVG Name
func (ttf *TTFFont) SvgName() string {
    return ttf.fontName
}


// A no-op drawable image used for measuring the font
type nopDrawImage  int

func (ndi nopDrawImage) ColorModel() color.Model {
    return color.GrayModel
}
func (ndi nopDrawImage) Bounds() image.Rectangle {
    return image.Rectangle{
        Min: image.Point{-int(math.MinInt32), -int(math.MinInt32)},
        Max: image.Point{-int(math.MaxInt32), -int(math.MaxInt32)},
    }
}
func (ndi nopDrawImage) At(x, y int) color.Color {
    return color.Black
}
func (ndi nopDrawImage) Set(x, y int, c color.Color) {
    // Do nothing
}