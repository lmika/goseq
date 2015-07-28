// Extra renderers available if ImageMagick is available.
//

//+build !noim

package main

import (
    "bytes"

    "github.com/quirkey/magick"
    "bitbucket.org/lmika/goseq/seqdiagram"
)

func PngRenderer(diagram *seqdiagram.Diagram, target string) error {
    if target == "" {
        target = "out.png"
    }

    svgbufr := new(bytes.Buffer)
    err := diagram.WriteSVG(svgbufr)
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

