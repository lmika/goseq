package graphbox

import (
	"bytes"
	"fmt"
)

// An icon which can be added to actors
type Icon interface {

	// Return the size of the icon
	Size() (width int, height int)

	// Draw the icon onto the draw context centered at x and y
	Draw(ctx DrawContext, x int, y int, lineStyle *SvgStyle)
}

// A stick figure icon
//

const stickPersonIconHead = 10
const stickPersonTorsoLength = 16
const stickPersonLegLength = 18
const stickPersonLegGap = 8
const stickPersonSholders = 2
const stickPersonArmLength = 12
const stickPersonArmGap = 12

type StickPersonIcon int

func (spi StickPersonIcon) Size() (width int, height int) {
	width = maxInt(maxInt(stickPersonIconHead, stickPersonArmGap), stickPersonLegGap) * 2
	height = stickPersonIconHead*2 + stickPersonTorsoLength + stickPersonLegLength
	return
}

func (spi StickPersonIcon) Draw(ctx DrawContext, x int, y int, lineStyle *SvgStyle) {
	style := lineStyle.ToStyle()

	_, h := spi.Size()
	ty := y - h/2

	headX, headY, headR := x, ty+stickPersonIconHead, stickPersonIconHead
	torsoY1 := headY + headR
	sholdersY := torsoY1 + stickPersonSholders
	torsoY2 := torsoY1 + stickPersonTorsoLength
	legY := torsoY2 + stickPersonLegLength

	ctx.Canvas.Circle(headX, headY, headR, style)
	ctx.Canvas.Line(x, headY+headR, x, torsoY2, style)
	ctx.Canvas.Line(x, sholdersY, x-stickPersonArmGap, sholdersY+stickPersonArmLength, style)
	ctx.Canvas.Line(x, sholdersY, x+stickPersonArmGap, sholdersY+stickPersonArmLength, style)
	ctx.Canvas.Line(x, torsoY2, x+stickPersonLegGap, legY, style)
	ctx.Canvas.Line(x, torsoY2, x-stickPersonLegGap, legY, style)
	ctx.Canvas.Line(x, torsoY2, x+stickPersonLegGap, legY, style)
}

// A cylinder suggesting a data source
//

const cylinderSmallRadius = 5
const cylinderLargeRadius = 18
const cylinderHeight = 28

type CylinderIcon int

func (ci CylinderIcon) Size() (width int, height int) {
	width = cylinderLargeRadius * 2
	height = cylinderHeight + cylinderSmallRadius*3
	return
}

func (ci CylinderIcon) Draw(ctx DrawContext, x int, y int, lineStyle *SvgStyle) {
	//	style := "stroke:black;fill:white;stroke-width:2px;"
	style := lineStyle.ToStyle()

	leftX, rightX := x-cylinderLargeRadius, x+cylinderLargeRadius
	upperEllipseY := y - cylinderHeight/2
	lowerEllipseY := y + cylinderHeight/2

	ci.drawCurve(ctx, leftX, upperEllipseY, rightX, upperEllipseY, cylinderSmallRadius, style)
	ci.drawCurve(ctx, leftX, upperEllipseY, rightX, upperEllipseY, -cylinderSmallRadius, style)

	ctx.Canvas.Line(leftX, upperEllipseY, leftX, lowerEllipseY, style)
	ctx.Canvas.Line(rightX, upperEllipseY, rightX, lowerEllipseY, style)

	ci.drawCurve(ctx, leftX, lowerEllipseY, rightX, lowerEllipseY, -cylinderSmallRadius, style)

	//ctx.Canvas.Path("M", ...)
}

func (ci CylinderIcon) drawCurve(ctx DrawContext, fx, fy, tx, ty, mag int, style string) {
	pathCmds := new(bytes.Buffer)

	fmt.Fprint(pathCmds, "M", fx, fy, " ")
	fmt.Fprint(pathCmds, "C", fx, fy-mag*2, ",", tx, ty-mag*2, ",", tx, ty)

	ctx.Canvas.Path(pathCmds.String(), style)
}

// A cloud
//

type PathIcon struct {
	Data PathIconData
}

func (pi PathIcon) Size() (width int, height int) {
	width = int(pi.Data.TargetIconSize * pi.Data.Width / pi.Data.Height)
	height = int(pi.Data.TargetIconSize)
	return
}

func (pi PathIcon) Draw(ctx DrawContext, x int, y int, lineStyle *SvgStyle) {
	scaleFactor := (pi.Data.TargetIconSize - pi.Data.IconPadding*2) / pi.Data.Width

	style := lineStyle.ToStyle()

	tx, ty := float64(x)/scaleFactor-pi.Data.Width/2.0, float64(y)/scaleFactor-pi.Data.Height/2.0
	transformations := fmt.Sprintf("scale(%f) translate(%f %f)", scaleFactor, tx, ty)

	ctx.Canvas.Group(style)
	ctx.Canvas.Path(pi.Data.Path, "transform=\""+transformations+"\"", "style=\"stroke-width:10px\"")
	ctx.Canvas.Gend()
}

type PathIconData struct {
	Width          float64
	Height         float64
	TargetIconSize float64
	IconPadding    float64
	Path           string
}


// CloudPathData represents the path data for the cloud icon.
// The image was adapted from the one at https://freesvg.org/simple-white-cloud-icon-vector-graphics
var CloudPathData = PathIconData{
	Width:          365.0,
	Height:         185.0,
	TargetIconSize: 44.0,
	IconPadding:    -18.0,
	Path:           `
		m299.75 60.587c-4.108 0-8.123 0.411-12.001 1.188-11.95-35.875-46.03-61.775-86.23-61.775-33.761
		0-63.196 18.27-78.847 45.363-0.919-0.036-1.84-0.07-2.769-0.07-23.494 0-44.202 11.816-56.409
		29.777-3.263-0.617-6.627-0.953-10.071-0.953-29.503 0-53.42 23.702-53.42 52.943 0 29.24 23.917
		52.94 53.421 52.94h66.435 0.045 0.045 81.525 0.045 0.047 98.145 0.046c33.28 0 60.25-26.73
		60.25-59.71 0-32.972-26.97-59.703-60.25-59.703z
	`,
}
