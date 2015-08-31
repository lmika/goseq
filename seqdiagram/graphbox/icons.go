package graphbox

// An icon which can be added to actors
type Icon interface {

    // Return the size of the icon
    Size() (width int, height int)

    // Draw the icon onto the draw context centered at x and y
    Draw(ctx DrawContext, x int, y int)
}



const stickPersonIconHead = 10
const stickPersonTorsoLength = 16
const stickPersonLegLength = 18
const stickPersonLegGap = 8
const stickPersonSholders = 2
const stickPersonArmLength = 12
const stickPersonArmGap = 12

// A stick figure icon
type StickPersonIcon int

func (spi StickPersonIcon) Size() (width int, height int) {
    w := maxInt(maxInt(stickPersonIconHead, stickPersonArmGap), stickPersonLegGap) * 2
    h := stickPersonIconHead * 2 + stickPersonTorsoLength + stickPersonLegLength
    return w, h
}

func (spi StickPersonIcon) Draw(ctx DrawContext, x int, y int) {
    style := "stroke:black;fill:white;stroke-width:2px;"

    _, h := spi.Size()
    ty := y - h / 2

    headX, headY, headR := x, ty + stickPersonIconHead, stickPersonIconHead
    torsoY1 := headY + headR
    sholdersY := torsoY1 + stickPersonSholders
    torsoY2 := torsoY1 + stickPersonTorsoLength
    legY := torsoY2 + stickPersonLegLength

    ctx.Canvas.Circle(headX, headY, headR, style)
    ctx.Canvas.Line(x, headY + headR, x, torsoY2, style)
    ctx.Canvas.Line(x, sholdersY, x - stickPersonArmGap, sholdersY + stickPersonArmLength, style)
    ctx.Canvas.Line(x, sholdersY, x + stickPersonArmGap, sholdersY + stickPersonArmLength, style)
    ctx.Canvas.Line(x, torsoY2, x + stickPersonLegGap, legY, style)
    ctx.Canvas.Line(x, torsoY2, x - stickPersonLegGap, legY, style)
    ctx.Canvas.Line(x, torsoY2, x + stickPersonLegGap, legY, style)
}