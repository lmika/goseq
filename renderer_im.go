// Extra renderers available if ImageMagick is available.
//

//+build im

package main

import (
	"bytes"

	"github.com/lmika/goseq/seqdiagram"
	"github.com/quirkey/magick"
)

func PngRenderer(diagram *seqdiagram.Diagram, opts *seqdiagram.ImageOptions, target string) error {
	if target == "" {
		target = "out.png"
	}

	svgbufr := new(bytes.Buffer)
	err := diagram.WriteSVGWithOptions(svgbufr, opts)
	if err != nil {
		return err
	}

	img, err := magick.NewFromBlob(svgbufr.Bytes(), "svg")
	if err != nil {
		return err
	}
	defer img.Destroy()

	return img.ToFile(target)
}
