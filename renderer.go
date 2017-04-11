package main

import (
	"os"

	"github.com/lmika/goseq/seqdiagram"
)

// Renders the result of the SVG to a destination (e.g. a file)
// If the filename is blank, the result is to go to the "default" destination
// (which is up to the renderer).
type Renderer func(diagram *seqdiagram.Diagram, opts *seqdiagram.ImageOptions, target string) error

// The default renderer: write the diagram to SVG
func SvgRenderer(diagram *seqdiagram.Diagram, opts *seqdiagram.ImageOptions, target string) error {
	if target != "" {
		file, err := os.Create(target)
		if err != nil {
			return err
		}
		defer file.Close()

		return diagram.WriteSVGWithOptions(file, opts)
	} else {
		return diagram.WriteSVGWithOptions(os.Stdout, opts)
	}
}
