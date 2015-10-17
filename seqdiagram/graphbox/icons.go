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
