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

    "github.com/golang/freetype"
    "github.com/golang/freetype/truetype"

    "golang.org/x/image/font"
    "golang.org/x/image/math/fixed"
)

type Font interface {
    // Returns the appropriate name of this font in the SVG
    SvgName() string

    // Measures the size of the particular line of text
    Measure(txt string, size float64) (int, int)
}

// Given a font, font size, points and gravity, returns a rectangle which will contain
// the text centered.  The point and gravity describes the location of the rect.
// The second point is where the text is to start given that it is to be rendered to
// fill the rectangle with default anchoring and alignment
func MeasureFontRect(font Font, size int, text string, x, y int, gravity Gravity) (Rect, Point) {
    w, h := font.Measure(text, float64(size))
    ox, oy := gravity(w, h)
    tp := Point{x - ox, y - oy + h}
    // HACK: May now work in all cases.  Adjust for the hanging measurements
    tp.Y -= size * 1 / 4 - 1
    return Rect{x - ox, y - oy, w, h}, tp
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
// TODO: For long bits of text, the measurement can be slightly off
func (ttf *TTFFont) Measure(txt string, size float64) (int, int) {
    img := nopDrawImage(0)

    ctx := freetype.NewContext()
    ctx.SetDPI(72)
    ctx.SetClip(img.Bounds())
    ctx.SetSrc(img)
    ctx.SetDst(img)
    ctx.SetFont(ttf.font)
    ctx.SetHinting(font.HintingFull)
    ctx.SetFontSize(size)

    np, _ := ctx.DrawString(txt, freetype.Pt(0, 0))

    mx, my := ttf.roundFix32(np.X), int(size) + ttf.roundFix32(np.Y)

    return mx, my
}

// Round a 26.6 fixed number to the nearest integer.
func (ttf *TTFFont) roundFix32(x fixed.Int26_6) int {
    full := int(x >> 6)
    rem := int(x & 0xCF)

    if rem > 0 {
        return full + 1
    } else {
        return full
    }
}

// Return the SVG Name
func (ttf *TTFFont) SvgName() string {
    // !!HACK!!  Need to properly determine how specific font families are defined in SVG
    return ttf.fontName + ",sans-serif"
}


// A no-op drawable image used for measuring the font
type nopDrawImage  int

func (ndi nopDrawImage) ColorModel() color.Model {
    return color.GrayModel
}
func (ndi nopDrawImage) Bounds() image.Rectangle {
    return image.Rectangle{
        Min: image.Point{int(math.MinInt32), int(math.MinInt32)},
        Max: image.Point{int(math.MaxInt32), int(math.MaxInt32)},
    }
}
func (ndi nopDrawImage) At(x, y int) color.Color {
    return color.Black
}
func (ndi nopDrawImage) Set(x, y int, c color.Color) {
    // Do nothing
}