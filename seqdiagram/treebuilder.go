// Adds the parse tree to the model

package seqdiagram

import (
    "fmt"
    "strings"

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
    parse.SPACER_GAP: DTSpacer,
    parse.EMPTY_GAP: DTGap,
    parse.FRAME_GAP: DTFrame,
    parse.LINE_GAP: DTLine,
}

var segmentTypeMap = map[parse.SegmentType]SegmentType {
    parse.ALT_SEGMENT: AltSegmentType,
    parse.ALT_ELSE_SEGMENT: ElseSegmentType,
    parse.LOOP_SEGMENT: LoopSegmentType,
}

type treeBuilder struct {
    nodeList        *parse.NodeList
    filename        string
}

func (tb *treeBuilder) buildTree(d *Diagram) error {
    for nodeList := tb.nodeList; nodeList != nil; nodeList = nodeList.Tail {
        seqItem, err := tb.toSequenceItem(nodeList.Head, d)
        if err != nil {
            return err
        } else if seqItem != nil {
            d.AddSequenceItem(seqItem)
        }
    }

    return nil
}

func (tb *treeBuilder) nodesToSlice(nodeList *parse.NodeList, d *Diagram) ([]SequenceItem, error) {
    seq := make([]SequenceItem, 0)

    for ; nodeList != nil; nodeList = nodeList.Tail {
        seqItem, err := tb.toSequenceItem(nodeList.Head, d)
        if err != nil {
            return nil, err
        } else if seqItem != nil {
            seq = append(seq, seqItem)
        }
    }

    return seq, nil
}

func (tb *treeBuilder) makeError(msg string) error {
    return fmt.Errorf("%s:%s", tb.filename, msg)
}

func (tb *treeBuilder) toSequenceItem(node parse.Node, d *Diagram) (SequenceItem, error) {
    switch n := node.(type) {
    case *parse.ProcessInstructionNode:
        d.ProcessingInstructions = append(d.ProcessingInstructions, &ProcessingInstruction{
            Prefix: n.Prefix,
            Value: n.Value,
        })
        return nil, nil
    case *parse.TitleNode:
        d.Title = n.Title
        return nil, nil
    case *parse.ActorNode:
        err := tb.addActor(n, d) // d.GetOrAddActorWithOptions(n.Ident, n.ActorName())
        return nil, err
    case *parse.ActionNode:
        return tb.addAction(n, d)
    case *parse.NoteNode:
        return tb.addNote(n, d)
    case *parse.GapNode:
        return tb.addGap(n, d)
    case *parse.BlockNode:
        return tb.addBlock(n, d)
    default:
        return nil, tb.makeError("Unrecognised declaration")
    }
}

func (tb *treeBuilder) addActor(an *parse.ActorNode, d *Diagram) error {
    actor := d.GetOrAddActorWithOptions(an.Ident, an.ActorName())

    if attrs := an.Attributes ; attrs != nil {
        attrMap, err := tb.attrsToMap(attrs, d)
        if err != nil {
            return err
        }

        // Configure the attributes
        if iconName, hasIconName := attrMap.Get("icon") ; hasIconName {
            if icon, err := LookupActorIcon(iconName); err == nil {
                actor.Icon = icon
            } else {
                return fmt.Errorf("error loading icon '%s': %s", iconName, err.Error())
            }
        }

        headerAttr, _ := attrMap.Get("header")
        footerAttr, _ := attrMap.Get("footer")
        
        actor.InHeader = (headerAttr != "none")
        actor.InFooter = (footerAttr != "none")
    }

    return nil
}

func (tb *treeBuilder) addAction(an *parse.ActionNode, d *Diagram) (SequenceItem, error) {
    from, err := tb.getOrAddActor(an.From, d)
    if err != nil {
        return nil, err
    }

    to, err := tb.getOrAddActor(an.To, d)
    if err != nil {
        return nil, err
    }

    arrow := Arrow{arrowStemMap[an.Arrow.Stem], arrowHeadMap[an.Arrow.Head]}
    action := &Action{from, to, arrow, an.Descr}
    return action, nil
}

func (tb *treeBuilder) addNote(nn *parse.NoteNode, d *Diagram) (SequenceItem, error) {
    actor1, err := tb.getOrAddActor(nn.Actor1, d)
    if err != nil {
        return nil, err
    }

    var actor2 *Actor = nil
    if nn.Actor2 != nil {
        actor2, err = tb.getOrAddActor(nn.Actor2, d)
        if err != nil {
            return nil, err
        }        
    }

    note := &Note{actor1, actor2, noteAlignmentMap[nn.Position], nn.Descr}
    return note, nil
}

func (tb *treeBuilder) getOrAddActor(ar parse.ActorRef, d *Diagram) (*Actor, error) {
    switch a := ar.(type) {
    case parse.NormalActorRef:
        return d.GetOrAddActor(string(a)), nil
    case parse.PseudoActorRef:
        pn := string(a)
        switch pn {
        case "left":
            return LeftOffsideActor, nil
        case "right":
            return RightOffsideActor, nil
        default:
            return nil, fmt.Errorf("Invalid pseudo actor: ", pn)
        }
    default:
        return nil, fmt.Errorf("Unknown actor reference")
    }
}

func (tb *treeBuilder) addGap(gn *parse.GapNode, d *Diagram) (SequenceItem, error) {
    divider := &Divider{gn.Descr, dividerTypeMap[gn.Type]}
    return divider, nil
}

func (tb *treeBuilder) addBlock(bn *parse.BlockNode, d *Diagram) (SequenceItem, error) {
    segs := make([]*BlockSegment, 0)
    for sn := bn.Segments; sn != nil; sn = sn.Tail {
        seg, err := tb.buildSegment(sn.Head, d)
        if err != nil {
            return nil, err
        }

        segs = append(segs, seg)
    }

    return &Block{segs}, nil
}

func (tb *treeBuilder) buildSegment(sn *parse.BlockSegment, d *Diagram) (*BlockSegment, error) {
    slice, err := tb.nodesToSlice(sn.SubNodes, d)
    if err != nil {
        return nil, err
    }

    return &BlockSegment{
        Type: segmentTypeMap[sn.Type],
        Prefix: sn.Prefix,
        Message: sn.Message,
        SubItems: slice,
    }, nil
}

func (tb *treeBuilder) attrsToMap(attrs *parse.AttributeList, d *Diagram) (AttributeSet, error) {
    attrMaps := make(map[string]string)

    for ; attrs != nil ; attrs = attrs.Tail {
        attr := attrs.Head
        attrMaps[attr.Name] = attr.Value
    }

    return AttributeSet(attrMaps), nil
}


// An attribute set
type AttributeSet map[string]string

// Get an attribute value and if the attribute is 
func (as AttributeSet) Get(name string) (value string, hasValue bool) {
    value, hasValue = as[name]
    return
}

// Gets a boolean value.  If the value is undefined, returns the default.
func (as AttributeSet) GetBool(name string, def bool) bool {
    if value, hasValue := as.Get(name) ; hasValue {
        value = strings.ToLower(value)
        if (value == "true") || (value == "yes") || (value == "on") || (value == "1") {
            return true
        } else {
            return false
        }
    } else {
        return def
    }
}