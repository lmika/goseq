package graphbox

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/lmika/goseq/seqdiagram/canvas"
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
	//style := lineStyle.ToStyle()
	strokeStyle := canvas.StrokeStyle{Color: color.Black, Width: 2} // TODO Styling
	fillStyle := canvas.FillStyle{Color: color.White}

	_, h := spi.Size()
	ty := y - h/2

	headX, headY, headR := x, ty+stickPersonIconHead, stickPersonIconHead
	torsoY1 := headY + headR
	sholdersY := torsoY1 + stickPersonSholders
	torsoY2 := torsoY1 + stickPersonTorsoLength
	legY := torsoY2 + stickPersonLegLength

	ctx.Canvas.Circle(headX, headY, headR, strokeStyle, fillStyle)
	ctx.Canvas.Line(x, headY+headR, x, torsoY2, strokeStyle)
	ctx.Canvas.Line(x, sholdersY, x-stickPersonArmGap, sholdersY+stickPersonArmLength, strokeStyle)
	ctx.Canvas.Line(x, sholdersY, x+stickPersonArmGap, sholdersY+stickPersonArmLength, strokeStyle)
	ctx.Canvas.Line(x, torsoY2, x+stickPersonLegGap, legY, strokeStyle)
	ctx.Canvas.Line(x, torsoY2, x-stickPersonLegGap, legY, strokeStyle)
	ctx.Canvas.Line(x, torsoY2, x+stickPersonLegGap, legY, strokeStyle)
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
	//style := lineStyle.ToStyle()
	storkeStyle := canvas.StrokeStyle{Color: color.Black, Width: 2} // TODO Styling
	fillStyle := canvas.FillStyle{Color: color.White}

	leftX, rightX := x-cylinderLargeRadius, x+cylinderLargeRadius
	upperEllipseY := y - cylinderHeight/2
	lowerEllipseY := y + cylinderHeight/2

	ci.drawCurve(ctx, leftX, upperEllipseY, rightX, upperEllipseY, cylinderSmallRadius, storkeStyle, fillStyle)
	ci.drawCurve(ctx, leftX, upperEllipseY, rightX, upperEllipseY, -cylinderSmallRadius, storkeStyle, fillStyle)

	ctx.Canvas.Line(leftX, upperEllipseY, leftX, lowerEllipseY, storkeStyle)
	ctx.Canvas.Line(rightX, upperEllipseY, rightX, lowerEllipseY, storkeStyle)

	ci.drawCurve(ctx, leftX, lowerEllipseY, rightX, lowerEllipseY, -cylinderSmallRadius, storkeStyle, fillStyle)

	//ctx.Canvas.Path("M", ...)
}

func (ci CylinderIcon) drawCurve(ctx DrawContext, fx, fy, tx, ty, mag int, strokeStyle canvas.StrokeStyle, fillStyle canvas.FillStyle) {
	pathCmds := new(bytes.Buffer)

	fmt.Fprint(pathCmds, "M", fx, fy, " ")
	fmt.Fprint(pathCmds, "C", fx, fy-mag*2, ",", tx, ty-mag*2, ",", tx, ty)

	ctx.Canvas.Path(pathCmds.String(), strokeStyle, fillStyle)
}
