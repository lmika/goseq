package main

import (
    "os"
    "log"
    "flag"
    "path/filepath"

    "bitbucket.org/lmika/goseq/goseq"
)

// Name of the output file
var flagOut = flag.String("o", "", "Output file")

func main() {    
    renderer := SvgRenderer

    flag.Parse()

    // Select a suitable renderer (based on the suffix of the output file, if there is one)
    if *flagOut != "" {
        ext := filepath.Ext(*flagOut)
        if ext == ".png" {
            renderer = PngRenderer
        } else if ext != ".svg" {
            log.Fatal("Unrecognised extension: " + ext)
        }
    }

    diagram, err := goseq.Parse(os.Stdin)
    if err != nil {
        log.Fatal(err)
    }

    err = renderer(diagram, *flagOut)
    if err != nil {
        log.Fatal(err)
    }
}