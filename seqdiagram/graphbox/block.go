package graphbox

// A block stype
type BlockStyle struct {
    Margin      Point
    Padding     Point
}

// A block
type Block struct {
    TR              int
    TC              int

    Style           BlockStyle
}

func NewBlock(toRow int, toCol int, style BlockStyle) *Block {
    return &Block{toRow, toCol, style}
}

func (block *Block) Constraint(r, c int, applier ConstraintApplier) {
    // There must be enought horizontal space to accomodate the text
    // and vertical space to display the divider
    applier.Apply(AddSizeConstraint{r, c, r + 1, block.TC, 0, block.Style.Margin.Y + block.Style.Padding.Y})
    applier.Apply(AddSizeConstraint{block.TR - 1, 0, block.TR, block.TC, 0, block.Style.Margin.Y + block.Style.Padding.Y})
}

func (block *Block) Draw(ctx DrawContext, point Point) {
    fx, fy := point.X, point.Y
    if point, isPoint := ctx.PointAt(block.TR, block.TC) ; isPoint {
        tx, ty := point.X, point.Y

        frameFy := fy + block.Style.Margin.Y + block.Style.Padding.Y
        w, h := tx - fx, ty - frameFy - block.Style.Margin.Y + block.Style.Padding.Y

        ctx.Canvas.Rect(fx, frameFy, w, h, "stroke:black;stroke-dasharray:14,6;stroke-width:2px;fill:none;")
    }
}
