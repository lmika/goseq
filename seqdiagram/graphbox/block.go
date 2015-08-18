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

    Style               BlockStyle

    prefixTextBox       *TextBox
    prefixTextBoxRect   Rect
    messageTextBox      *TextBox
    messageTextBoxRect  Rect
}

func NewBlock(toRow int, toCol int, prefix string, text string, style BlockStyle) *Block {
    prefixTextBox := NewTextBox(style.Font, style.FontSize, MiddleTextAlign)
    prefixTextBox.AddText(prefix)
    prefixTextBoxRect := prefixTextBox.BoundingRect()

    messageTextBox := NewTextBox(style.Font, style.FontSize, MiddleTextAlign)
    messageTextBox.AddText(text)
    messageTextBoxRect := messageTextBox.BoundingRect()

    return &Block{toRow, toCol, style, prefixTextBox, prefixTextBoxRect, messageTextBox, messageTextBoxRect}
}

func (block *Block) Constraint(r, c int, applier ConstraintApplier) {
    prefixExtraWidth := block.Style.PrefixExtraWidth + block.Style.TextPadding.X * 2 + block.Style.FontSize / 2

    minWidth := block.prefixTextBoxRect.W + block.messageTextBoxRect.W + 
            prefixExtraWidth + block.Style.GapWidth + block.Style.TextPadding.X * 2 -
            block.Style.Margin.X * 2

    textHeight := maxInt(block.prefixTextBoxRect.H, block.messageTextBoxRect.H)

    
    applier.Apply(TotalSizeConstraint{r, c - 1, r + 1, c, block.Style.Margin.X, 0})
    applier.Apply(TotalSizeConstraint{r, block.TC, r + 1, block.TC + 1, block.Style.Margin.X, 0})
    applier.Apply(TotalSizeConstraint{r, c, r + 1, block.TC, minWidth, 0})

    applier.Apply(AddSizeConstraint{r - 1, c, r, block.TC, 0, block.Style.Margin.Y})
    applier.Apply(AddSizeConstraint{r, c, r + 1, block.TC, 0, textHeight + block.Style.MidMargin})
    applier.Apply(AddSizeConstraint{block.TR - 1, c, block.TR, block.TC, 0, block.Style.Margin.Y})
}

func (block *Block) Draw(ctx DrawContext, point Point) {
    fx, fy := point.X, point.Y
    if point, isPoint := ctx.PointAt(block.TR, block.TC) ; isPoint {
        tx, ty := point.X, point.Y

        fx -= block.Style.Margin.X
        tx += block.Style.Margin.X

        block.drawText(ctx, fx, fy)
        block.drawFrame(ctx, fx, fy, tx, ty)
    }
}

func (block *Block) drawFrame(ctx DrawContext, fx, fy, tx, ty int) {
    w := tx - fx
    h := ty - fy
    ctx.Canvas.Rect(fx, fy, w, h, "stroke:black;stroke-dasharray:4,4;stroke-width:2px;fill:none;")
}

func (block *Block) drawText(ctx DrawContext, fx, fy int) {
    ptr := block.prefixTextBoxRect.BlowOut(block.Style.TextPadding).AddSize(block.Style.PrefixExtraWidth, 0).PositionAt(fx, fy, NorthWestGravity)
    mtr := block.messageTextBoxRect.BlowOut(block.Style.MessagePadding).PositionAt(fx + ptr.W, fy, NorthWestGravity)

    ctx.Canvas.Rect(mtr.X, mtr.Y, mtr.W + block.Style.GapWidth + block.Style.FontSize / 2, mtr.H, "stroke:none;fill:white;")
    block.messageTextBox.Render(ctx.Canvas, mtr.X + block.Style.GapWidth + block.Style.MessagePadding.X, mtr.Y + block.Style.MessagePadding.Y, NorthWestGravity)

    //ctx.Canvas.Rect(ptr.X, ptr.Y, ptr.W, ptr.H, "stroke:black;stroke-width:2px;fill:white;")
    block.drawPageFrame(ctx, ptr.X, ptr.Y, ptr.X + ptr.W, ptr.Y + ptr.H)
    block.prefixTextBox.Render(ctx.Canvas, ptr.X + block.Style.TextPadding.X, ptr.Y + block.Style.TextPadding.Y, NorthWestGravity)
}

func (block *Block) drawPageFrame(ctx DrawContext, fx, fy, tx, ty int) {
    fold := block.Style.FontSize / 2

    xs := []int { fx, fx, tx - fold, tx, tx }
    ys := []int { fy, ty, ty, ty - fold, fy }

    ctx.Canvas.Polygon(xs, ys, "stroke:black;stroke-width:2px;fill:white;")
}