package main

import (
    "os"

    "bitbucket.org/lmika/goseq/seqdiagram"
)

// Renders the result of the SVG to a destination (e.g. a file)
// If the filename is blank, the result is to go to the "default" destination
// (which is up to the renderer).
type Renderer func(diagram *seqdiagram.Diagram, target string) error


// The default renderer: write the diagram to SVG
func SvgRenderer(diagram *seqdiagram.Diagram, target string) error {
    if target != "" {
        file, err := os.Create(target)
        if err != nil {
            return err
        }
        defer file.Close()

        return diagram.WriteSVG(file)
    } else {
        return diagram.WriteSVG(os.Stdout)
    }
}


/*
// The internal PNG renderer.  This is set if the tags 'im' is set.
var internalPngRenderer Renderer = nil

// The PNG renderer
func PngRenderer(diagram *goseq.Diagram, target string) error {
    if internalPngRenderer != nil {
        return internalPngRenderer(diagram, target)
    } else {
        return errors.New("PNG rendering not supported")
    }
}
*/