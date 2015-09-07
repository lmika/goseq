// Renderers disabled if noim is specified
//

//+build !im

package main

import (
    "errors"
    "bitbucket.org/lmika/goseq/seqdiagram"
)

func PngRenderer(diagram *seqdiagram.Diagram, style *seqdiagram.DiagramStyles, target string) error {
    return errors.New("PNG renderer not available")
}
