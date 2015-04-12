package main

import (
    "os"
    "fmt"
    "flag"
    "path/filepath"

    "bitbucket.org/lmika/goseq/seqdiagram"
)

// Name of the output file
var flagOut = flag.String("o", "", "Output file")

// Die with error
func die(msg string) {
    fmt.Fprintf(os.Stderr, "goseq: %s\n", msg)
    os.Exit(1)
}

// Processes a file
func processFile(inFilename string, outFilename string, renderer Renderer) error {
    var infile *os.File
    var err error

    if inFilename == "-" {
        infile = os.Stdin
    } else {
        infile, err = os.Open(inFilename)
        if err != nil {
            return err
        }
        defer infile.Close()
    }

    diagram, err := seqdiagram.ParseDiagram(infile, inFilename)
    if err != nil {
        return err
    }

    err = renderer(diagram, outFilename)
    if err != nil {
        return err
    }

    return nil
}

func main() {    
    renderer := SvgRenderer
    outFile := ""

    flag.Parse()

    // Select a suitable renderer (based on the suffix of the output file, if there is one)
    if *flagOut != "" {
        ext := filepath.Ext(*flagOut)
        if ext == ".png" {
            renderer = PngRenderer
        } else if ext != ".svg" {
            die("Unsupported extension: " + ext)
        }

        outFile = *flagOut
    }

    // Process each file (or stdin)
    if flag.NArg() == 0 {
        err := processFile("-", outFile, renderer)
        if err != nil {
            die("stdin - " + err.Error())
        }
    } else {
        for _, inFile := range flag.Args() {
            err := processFile(inFile, outFile, renderer)
            if err != nil {
                die(inFile + " - " + err.Error())
            }
        }
    }
}