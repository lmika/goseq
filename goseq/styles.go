package goseq

import (
    "bitbucket.org/lmika/goseq/goseq/graphbox"
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

    // Styling of the diagram title
    Title               graphbox.TitleStyle

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
    Title: graphbox.TitleStyle {
        Font: standardFont,
        FontSize: 20,
        Padding: graphbox.Point{4, 16},
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
        DTLine: graphbox.DividerStyle {
            Font: standardFont,
            FontSize: 14,
            Padding: graphbox.Point{16, 4},
            Margin: graphbox.Point{8, 16},
            TextPadding: graphbox.Point{0, 0},
            Shape: graphbox.DSFullLine,
        },
    },
}
