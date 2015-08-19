package seqdiagram

import (
    "bitbucket.org/lmika/goseq/seqdiagram/graphbox"
)

// Diagram styles
type DiagramStyles struct {
    // Diagram margins
    Margin              graphbox.Point

    // Styling of the actor box
    ActorBox            graphbox.ActorBoxStyle

    // Styling of the note box
    NoteBox             graphbox.NoteBoxStyle

    // Styling of the activity line
    ActivityLine        graphbox.ActivityLineStyle

    // Styling of arrow heads
    ArrowHeads          map[ArrowHead]*graphbox.ArrowHeadStyle

    // Styling of the diagram title
    Title               graphbox.TitleStyle

    // Block styling
    Block               graphbox.BlockStyle

    // Styles of dividers
    Divider             map[DividerType]graphbox.DividerStyle
}

// Fonts
var standardFont = mustLoadFont()


// The Default style
var DefaultStyle = &DiagramStyles {
    Margin: graphbox.Point{8, 8},
    ActorBox: graphbox.ActorBoxStyle {
        Font: standardFont,
        FontSize: 16,
        Padding: graphbox.Point{16, 8},
        Margin: graphbox.Point{8, 8},
    },
    NoteBox: graphbox.NoteBoxStyle {
        Font: standardFont,
        FontSize: 14,
        Padding: graphbox.Point{8, 4},
        Margin: graphbox.Point{8, 8},
    },
    ActivityLine: graphbox.ActivityLineStyle{
        Font: standardFont,
        FontSize: 14,
        Margin: graphbox.Point{16, 8},
        TextGap: 4,
    },
    ArrowHeads: map[ArrowHead]*graphbox.ArrowHeadStyle {
        SolidArrowHead: &graphbox.ArrowHeadStyle {
            Xs: []int { -9, 0, -9 },
            Ys: []int { -5, 0, 5 },
            BaseStyle: "stroke:black;fill:black;stroke-width:2px;",
        },
        OpenArrowHead: &graphbox.ArrowHeadStyle {
            Xs: []int { -9, 0, -9 },
            Ys: []int { -5, 0, 5 },
            BaseStyle: "stroke:black;fill:none;stroke-width:2px;",
        },
        BarbArrowHead: &graphbox.ArrowHeadStyle {
            Xs: []int { -11, 0 },
            Ys: []int { -7, 0 },
            BaseStyle: "stroke:black;fill:black;stroke-width:2px;",
        },
        LowerBarbArrowHead: &graphbox.ArrowHeadStyle {
            Xs: []int { -11, 0 },
            Ys: []int { 7, 0 },
            BaseStyle: "stroke:black;fill:black;stroke-width:2px;",
        },
    },
    Title: graphbox.TitleStyle {
        Font: standardFont,
        FontSize: 20,
        Padding: graphbox.Point{4, 16},
    },
    Block: graphbox.BlockStyle{
        Margin: graphbox.Point{24, 8},
        TextPadding: graphbox.Point{4, 4},
        MessagePadding: graphbox.Point{4, 4},
        GapWidth: 4,
        PrefixExtraWidth: 4,

        Font: standardFont,
        FontSize: 14,
        MidMargin: 4,
    },
    Divider: map[DividerType]graphbox.DividerStyle {
        DTGap: graphbox.DividerStyle {
            Font: standardFont,
            FontSize: 14,
            Padding: graphbox.Point{16, 8},
            Margin: graphbox.Point{8, 8},
            TextPadding: graphbox.Point{0, 0},
            Shape: graphbox.DSFullRect,
        },
        DTFrame: graphbox.DividerStyle {
            Font: standardFont,
            FontSize: 14,
            Padding: graphbox.Point{16, 8},
            Margin: graphbox.Point{8, 8},
            TextPadding: graphbox.Point{0, 0},
            Shape: graphbox.DSFramedRect,
        },
        DTLine: graphbox.DividerStyle {
            Font: standardFont,
            FontSize: 14,
            Padding: graphbox.Point{16, 4},
            Margin: graphbox.Point{8, 16},
            TextPadding: graphbox.Point{4, 2},
            Shape: graphbox.DSFullLine,
        },
        DTSpacer: graphbox.DividerStyle {
            Font: standardFont,
            FontSize: 14,
            Padding: graphbox.Point{16, 4},
            Margin: graphbox.Point{8, 16},
            TextPadding: graphbox.Point{0, 0},
            Shape: graphbox.DSSpacerRect,
        },
    },
}


// The small style.  This has narrower margins and font sizes and
// is used to produce smaller diagrams.
var SmallStyle = &DiagramStyles {
    Margin: graphbox.Point{4, 4},
    ActorBox: graphbox.ActorBoxStyle {
        Font: standardFont,
        FontSize: 14,
        Padding: graphbox.Point{12, 6},
        Margin: graphbox.Point{8, 8},
    },
    NoteBox: graphbox.NoteBoxStyle {
        Font: standardFont,
        FontSize: 12,
        Padding: graphbox.Point{6, 3},
        Margin: graphbox.Point{6, 6},
    },
    ActivityLine: graphbox.ActivityLineStyle{
        Font: standardFont,
        FontSize: 12,
        Margin: graphbox.Point{8, 6},
        TextGap: 3,
    },
    ArrowHeads: map[ArrowHead]*graphbox.ArrowHeadStyle {
        SolidArrowHead: &graphbox.ArrowHeadStyle {
            Xs: []int { -7, 0, -7 },
            Ys: []int { -4, 0, 4 },
            BaseStyle: "stroke:black;fill:black;stroke-width:2px;",
        },
        OpenArrowHead: &graphbox.ArrowHeadStyle {
            Xs: []int { -7, 0, -7 },
            Ys: []int { -4, 0, 4 },
            BaseStyle: "stroke:black;fill:none;stroke-width:2px;",
        },
        BarbArrowHead: &graphbox.ArrowHeadStyle {
            Xs: []int { -9, 0 },
            Ys: []int { -5, 0 },
            BaseStyle: "stroke:black;fill:black;stroke-width:2px;",
        },
        LowerBarbArrowHead: &graphbox.ArrowHeadStyle {
            Xs: []int { -9, 0 },
            Ys: []int { 5, 0 },
            BaseStyle: "stroke:black;fill:black;stroke-width:2px;",
        },
    },
    Title: graphbox.TitleStyle {
        Font: standardFont,
        FontSize: 18,
        Padding: graphbox.Point{2, 8},
    },
    Block: graphbox.BlockStyle{
        Margin: graphbox.Point{18, 4},
        TextPadding: graphbox.Point{3, 3},
        MessagePadding: graphbox.Point{3, 3},
        GapWidth: 3,
        PrefixExtraWidth: 3,

        Font: standardFont,
        FontSize: 12,
        MidMargin: 2,
    },
    Divider: map[DividerType]graphbox.DividerStyle {
        DTGap: graphbox.DividerStyle {
            Font: standardFont,
            FontSize: 12,
            Padding: graphbox.Point{12, 6},
            Margin: graphbox.Point{6, 6},
            TextPadding: graphbox.Point{0, 0},
            Shape: graphbox.DSFullRect,
        },
        DTFrame: graphbox.DividerStyle {
            Font: standardFont,
            FontSize: 12,
            Padding: graphbox.Point{12, 6},
            Margin: graphbox.Point{6, 6},
            TextPadding: graphbox.Point{0, 0},
            Shape: graphbox.DSFramedRect,
        },
        DTLine: graphbox.DividerStyle {
            Font: standardFont,
            FontSize: 12,
            Padding: graphbox.Point{12, 6},
            Margin: graphbox.Point{6, 12},
            TextPadding: graphbox.Point{2, 1},
            Shape: graphbox.DSFullLine,
        },
        DTSpacer: graphbox.DividerStyle {
            Font: standardFont,
            FontSize: 12,
            Padding: graphbox.Point{12, 6},
            Margin: graphbox.Point{6, 12},
            TextPadding: graphbox.Point{2, 1},
            Shape: graphbox.DSSpacerRect,
        },
    },
}

var StyleNames = map[string]*DiagramStyles {
    "default": DefaultStyle,
    "small": SmallStyle,
}