// Adds the parse tree to the model

package seqdiagram

import (
    "fmt"

    "bitbucket.org/lmika/goseq/seqdiagram/parse"
)

var arrowStemMap = map[parse.ArrowStemType]ArrowStem {
    parse.SOLID_ARROW_STEM: SolidArrowStem,
    parse.DASHED_ARROW_STEM: DashedArrowStem,
    parse.THICK_ARROW_STEM: ThickArrowStem,
}

var arrowHeadMap = map[parse.ArrowHeadType]ArrowHead {
    parse.SOLID_ARROW_HEAD: SolidArrowHead,
    parse.OPEN_ARROW_HEAD: OpenArrowHead,
    parse.BARBED_ARROW_HEAD: BarbArrowHead,
    parse.LOWER_BARBED_ARROW_HEAD: LowerBarbArrowHead,
}

var noteAlignmentMap = map[parse.NoteAlignment]NoteAlignment {
    parse.LEFT_NOTE_ALIGNMENT: LeftNoteAlignment,
    parse.RIGHT_NOTE_ALIGNMENT: RightNoteAlignment,
    parse.OVER_NOTE_ALIGNMENT: OverNoteAlignment,
}

var dividerTypeMap = map[parse.GapType]DividerType {
    parse.EMPTY_GAP: DTGap,
    parse.FRAME_GAP: DTFrame,
    parse.LINE_GAP: DTLine,
}

type treeBuilder struct {
    nodeList        *parse.NodeList
    filename        string
}

func (tb *treeBuilder) buildTree(d *Diagram) error {
    for nodeList := tb.nodeList; nodeList != nil; nodeList = nodeList.Tail {
        err := tb.addNode(nodeList.Head, d)
        if err != nil {
            return err
        }
    }

    return nil
}

func (tb *treeBuilder) makeError(msg string) error {
    return fmt.Errorf("%s:%s", tb.filename, msg)
}

func (tb *treeBuilder) addNode(node parse.Node, d *Diagram) error {
    switch n := node.(type) {
    case *parse.ProcessInstructionNode:
        d.ProcessingInstructions = append(d.ProcessingInstructions, &ProcessingInstruction{
            Prefix: n.Prefix,
            Value: n.Value,
        })
        return nil
    case *parse.TitleNode:
        d.Title = n.Title
        return nil
    case *parse.ActorNode:
        d.GetOrAddActorWithOptions(n.Ident, n.ActorName())
        return nil
    case *parse.ActionNode:
        return tb.addAction(n, d)
    case *parse.NoteNode:
        return tb.addNote(n, d)
    case *parse.GapNode:
        return tb.addGap(n, d)
    default:
        return tb.makeError("Unrecognised declaration")
    }
}

func (tb *treeBuilder) addAction(an *parse.ActionNode, d *Diagram) error {
    arrow := Arrow{arrowStemMap[an.Arrow.Stem], arrowHeadMap[an.Arrow.Head]}
    action := &Action{d.GetOrAddActor(an.From), d.GetOrAddActor(an.To), arrow, an.Descr}
    d.AddSequenceItem(action)
    return nil
}

func (tb *treeBuilder) addNote(nn *parse.NoteNode, d *Diagram) error {
    note := &Note{d.GetOrAddActor(nn.Actor), noteAlignmentMap[nn.Position], nn.Descr}
    d.AddSequenceItem(note)
    return nil
}

func (tb *treeBuilder) addGap(gn *parse.GapNode, d *Diagram) error {
    divider := &Divider{gn.Descr, dividerTypeMap[gn.Type]}
    d.AddSequenceItem(divider)
    return nil
}