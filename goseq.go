package main

import (
    "os"
    "log"
    "./goseq"
//    "./goseq/graphbox"
)

func main() {
    diagram, err := goseq.Parse(os.Stdin)

    if err != nil {
        log.Fatal(err)
    }
    err = diagram.WriteSVG(os.Stdout)
    if err != nil {
        log.Fatal(err)
    }
}