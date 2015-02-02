package graphbox


// An array of constraints
type Constraints    []Constraint

func (cs Constraints) Apply(cm ConstraintChanger) {
    for _, c := range cs {
        c.Apply(cm)
    }
}


// Apply a size constraint which requests a minimum gap between
// points
type SizeConstraint struct {
    R, C            int
    Left, Right     int
    Top, Bottom     int
}

func (sc SizeConstraint) Apply(cm ConstraintChanger) {
    cm.EnsureLeftIsAtleast(sc.C, sc.Left)
    cm.EnsureLeftIsAtleast(sc.C + 1, sc.Right)
    cm.EnsureTopIsAtLeast(sc.R, sc.Top)
    cm.EnsureTopIsAtLeast(sc.R + 1, sc.Bottom)
}


// Adds a size constraint which requests a minimum gap between
// points
type AddSizeConstraint struct {
    R, C            int
    Left, Right     int
    Top, Bottom     int
}

func (sc AddSizeConstraint) Apply(cm ConstraintChanger) {
    cm.AddLeftToCol(sc.C, sc.Left)
    cm.AddLeftToCol(sc.C + 1, sc.Right)
    cm.AddTopToRow(sc.R, sc.Top)
    cm.AddTopToRow(sc.R + 1, sc.Bottom)
}

// Ensures that the total size between the two points is big enough for
// the rectangle.  If not, resize the grid points uniformally.
type TotalSizeConstraint struct {
    FR, FC          int
    TR, TC          int
    Width, Height   int
}

func (sc TotalSizeConstraint) Apply(cm ConstraintChanger) {
    w, h := cm.GridPointRect(sc.FR, sc.FC, sc.TR, sc.TC)

    if w < sc.Width {
        widthOfEachCell := sc.Width / (sc.TC - sc.FC)
        for c := sc.FC + 1; c <= sc.TC; c++ {
            cm.EnsureLeftIsAtleast(c, widthOfEachCell)
        }
    }

    if h < sc.Height {
        heightOfEachCell := sc.Height / (sc.TR - sc.FR)
        for r := sc.FR + 1; r <= sc.TR; r++ {
            cm.EnsureTopIsAtLeast(r, heightOfEachCell)
        }
    }
}