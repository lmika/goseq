package graphbox

// A block stype
type BlockStyle struct {
    Margin      Point

    Font        Font
    FontSize    int

    TextPadding Point
    MessagePadding Point
    
    PrefixExtraWidth int
    GapWidth    int
    MidMargin   int
}

// A block
type Block struct {
    TR                  int
    TC                  int

    MarginMup           int
    IsLast              bool
    ShowPrefix          bool
    ShowMessage         bool
    Style               BlockStyle

    prefixTextBox       *TextBox
    prefixTextBoxRect   Rect
    messageTextBox      *TextBox
    messageTextBoxRect  Rect
}

func NewBlock(toRow int, toCol int, marginMup int, isLast bool, prefix string, showPrefix bool, text string, style BlockStyle) *Block {
    prefixTextBox := NewTextBox(style.Font, style.FontSize, MiddleTextAlign)
    prefixTextBox.AddText(prefix)
    prefixTextBoxRect := prefixTextBox.BoundingRect()

    messageTextBox := NewTextBox(style.Font, style.FontSize, MiddleTextAlign)
    messageTextBox.AddText(text)
    messageTextBoxRect := messageTextBox.BoundingRect()

    return &Block{toRow, toCol, marginMup, isLast, showPrefix, text != "", style, prefixTextBox, prefixTextBoxRect, messageTextBox, messageTextBoxRect}
}

func (block *Block) Constraint(r, c int, applier ConstraintApplier) {
    prefixExtraWidth := block.Style.PrefixExtraWidth + block.Style.TextPadding.X * 2 + block.Style.FontSize / 2
    horizMargin := block.calcHorizMargin()

    minWidth := block.prefixTextBoxRect.W + block.messageTextBoxRect.W + 
            prefixExtraWidth + block.Style.GapWidth + block.Style.TextPadding.X * 2 -
            horizMargin * 2
    textHeight := maxInt(block.prefixTextBoxRect.H, block.messageTextBoxRect.H)

    
    applier.Apply(TotalSizeConstraint{r, c - 1, r + 1, c, horizMargin, 0})
    applier.Apply(TotalSizeConstraint{r, block.TC, r + 1, block.TC + 1, horizMargin, 0})
    applier.Apply(TotalSizeConstraint{r, c, r + 1, block.TC, minWidth, 0})

    applier.Apply(AddSizeConstraint{r, c, 0, 0, block.Style.Margin.Y, textHeight + block.Style.MidMargin})
    applier.Apply(AddSizeConstraint{block.TR, c, 0, 0, block.Style.Margin.Y, 0})
}

func (block *Block) Draw(ctx DrawContext, point Point) {
    fx, fy := point.X, point.Y
    if point, isPoint := ctx.PointAt(block.TR, block.TC) ; isPoint {
        tx, ty := point.X, point.Y

        fx -= block.calcHorizMargin()
        tx += block.calcHorizMargin()

        block.drawText(ctx, fx, fy)
        block.drawFrame(ctx, fx, fy, tx, ty)
    }
}

// Calculate the horizontal margin based on the configured style margin and depth
func (block *Block) calcHorizMargin() int {
    return block.Style.Margin.X * block.MarginMup
}

func (block *Block) drawFrame(ctx DrawContext, fx, fy, tx, ty int) {
    //w := tx - fx
    //h := ty - fy

    xs := []int { fx, fx, tx, tx }
    ys := []int { ty, fy, fy, ty }

    lineStyle := "stroke:black;stroke-dasharray:4,4;stroke-width:2px;fill:none;"
    if block.IsLast {
        //ctx.Canvas.Rect(fx, fy, w, h, lineStyle)
        ctx.Canvas.Polygon(xs, ys, lineStyle)
    } else {
        ctx.Canvas.Polyline(xs, ys, lineStyle)
        /*
        ctx.Canvas.Polyline(
            []int { fx, fx, tx, tx },
            []int { ty, fy, fy, ty },
            lineStyle)
        */
    }
}

func (block *Block) drawText(ctx DrawContext, fx, fy int) {
    ptr := block.prefixTextBoxRect.BlowOut(block.Style.TextPadding).AddSize(block.Style.PrefixExtraWidth, 0).PositionAt(fx, fy, NorthWestGravity)
    mtr := block.messageTextBoxRect.BlowOut(block.Style.MessagePadding).PositionAt(fx + ptr.W, fy, NorthWestGravity)

    if block.ShowMessage {
        ctx.Canvas.Rect(mtr.X, mtr.Y, mtr.W + block.Style.GapWidth + block.Style.FontSize / 2, mtr.H, "stroke:none;fill:white;")
        block.messageTextBox.Render(ctx.Canvas, mtr.X + block.Style.GapWidth + block.Style.MessagePadding.X, mtr.Y + block.Style.MessagePadding.Y, NorthWestGravity)
    }        

    if block.ShowPrefix {
        block.drawPrefixFrame(ctx, ptr.X, ptr.Y, ptr.X + ptr.W, ptr.Y + ptr.H)
        block.prefixTextBox.Render(ctx.Canvas, ptr.X + block.Style.TextPadding.X, ptr.Y + block.Style.TextPadding.Y, NorthWestGravity)
    }
}

func (block *Block) drawPrefixFrame(ctx DrawContext, fx, fy, tx, ty int) {
    fold := block.Style.FontSize / 2

    xs := []int { fx, fx, tx - fold, tx, tx }
    ys := []int { fy, ty, ty, ty - fold, fy }

    ctx.Canvas.Polygon(xs, ys, "stroke:black;stroke-width:2px;fill:white;")
}