package main

import (
    "os"
    "log"
    "flag"

    "bitbucket.org/lmika/goseq/goseq"
)

var flagPng = flag.Bool("p", false, "Render the image as a PNG")

func main() {    
    renderer := SvgRenderer

    flag.Parse()

    if *flagPng {
        renderer = PngRenderer
    }

    diagram, err := goseq.Parse(os.Stdin)
    if err != nil {
        log.Fatal(err)
    }

    //err = diagram.WriteSVG(os.Stdout)
    err = renderer(diagram, "")
    if err != nil {
        log.Fatal(err)
    }
}