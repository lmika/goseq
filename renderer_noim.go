// Renderers disabled if noim is specified
//
// +build noim

package main

import (
    "errors"
    "bitbucket.org/lmika/goseq/goseq"
)

func PngRenderer(diagram *goseq.Diagram, target string) error {
    return errors.New("PNG renderer not available")
}
