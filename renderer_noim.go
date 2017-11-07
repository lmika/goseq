// Renderers disabled if noim is specified
//

//+build !im

package main

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/lmika/goseq/seqdiagram"
)

func PngRenderer(diagram *seqdiagram.Diagram, opts *seqdiagram.ImageOptions, target string) error {
	if target == "" {
		target = "out.png"
	}

	image, err := diagram.Draw(opts)
	if err != nil {
		return err
	}

	// TEMP
	return draw2dimg.SaveToPngFile(target, image)
}
