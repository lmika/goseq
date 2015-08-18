package graphbox

// A block stype
type BlockStyle struct {
    Margin      Point
    Padding     Point

    Font        Font
    FontSize    int

    TextPadding Point
    MessagePadding Point
    
    PrefixExtraWidth int
    GapWidth    int
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
    // There must be enought horizontal space to accomodate the text
    // and vertical space to display the divider
    topMargin := block.Style.Margin.Y + block.Style.Padding.Y + block.Style.TextPadding.Y * 2 +
            maxInt(block.prefixTextBoxRect.H, block.messageTextBoxRect.H)
    minWidth := block.prefixTextBoxRect.W + block.messageTextBoxRect.W + 
            block.Style.PrefixExtraWidth + block.Style.GapWidth + block.Style.TextPadding.X * 2

    applier.Apply(TotalSizeConstraint{r, c, r + 1, block.TC, minWidth, 0})
    applier.Apply(AddSizeConstraint{r, c, r + 1, block.TC, 0, topMargin})
    applier.Apply(AddSizeConstraint{block.TR - 1, 0, block.TR, block.TC, 0, block.Style.Margin.Y + block.Style.Padding.Y})
}

func (block *Block) Draw(ctx DrawContext, point Point) {
    fx, fy := point.X, point.Y
    if point, isPoint := ctx.PointAt(block.TR, block.TC) ; isPoint {
        tx, ty := point.X, point.Y

        frameFy := fy + block.Style.Margin.Y + block.Style.Padding.Y
        w, h := tx - fx, ty - frameFy - block.Style.Margin.Y + block.Style.Padding.Y

        // The prefix and message text rectangles
        ptr := block.prefixTextBoxRect.BlowOut(block.Style.TextPadding).PositionAt(fx, frameFy, NorthWestGravity)
        ptr.W += block.Style.PrefixExtraWidth

        mtr := block.messageTextBoxRect.BlowOut(block.Style.MessagePadding).
                PositionAt(fx + ptr.W, frameFy, NorthWestGravity)

        ctx.Canvas.Rect(mtr.X, mtr.Y, mtr.W + block.Style.GapWidth, mtr.H, "stroke:none;fill:red;")
        block.messageTextBox.Render(ctx.Canvas, mtr.X + block.Style.GapWidth + block.Style.MessagePadding.X,
                mtr.Y + block.Style.MessagePadding.Y, NorthWestGravity)

        ctx.Canvas.Rect(ptr.X, ptr.Y, ptr.W, ptr.H, "stroke:black;stroke-width:2px;fill:white;")
        block.prefixTextBox.Render(ctx.Canvas, ptr.X + block.Style.TextPadding.X, ptr.Y + block.Style.TextPadding.Y, NorthWestGravity)

        ctx.Canvas.Rect(fx, frameFy, w, h, "stroke:black;stroke-dasharray:4,4;stroke-width:2px;fill:none;")
    }
}
