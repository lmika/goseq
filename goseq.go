package main

import (
    "io"
    "io/ioutil"
    "os"
    "fmt"
    "flag"
    "path/filepath"
    "strings"

    "bitbucket.org/lmika/goseq/seqdiagram"
)

// Name of the output file
var flagOut = flag.String("o", "", "Output file")

// The style to use
var flagStyle = flag.String("s", "default", "The style to use")

// Die with error
func die(msg string) {
    fmt.Fprintf(os.Stderr, "goseq: %s\n", msg)
    os.Exit(1)
}

// Processes a md file
func processMdFile(inFilename string, outFilename string, renderer Renderer) error {
    srcFile, err := openSourceFile(inFilename)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    targetFile := ioutil.Discard

    mf := &MarkdownFilter{srcFile, targetFile, func(codeblock string, output io.Writer) error {
        fmt.Fprint(output, codeblock)
        err := processSeqDiagram(strings.NewReader(codeblock), inFilename, "/dev/null", nil)
        if err != nil {
            fmt.Fprintf(os.Stderr, "goseq: %s:embedded block - %s\n", inFilename, err.Error())
        }
        return nil
    }}
    return mf.Scan()
}

// Processes a seq file
func processSeqFile(inFilename string, outFilename string, renderer Renderer) error {
    srcFile, err := openSourceFile(inFilename)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    return processSeqDiagram(srcFile, inFilename, outFilename, renderer)
}

// Processes a sequence diagram
func processSeqDiagram(infile io.Reader, inFilename string, outFilename string, renderer Renderer) error {
    diagram, err := seqdiagram.ParseDiagram(infile, inFilename)
    if err != nil {
        return err
    }

    style := seqdiagram.DefaultStyle
    if altStyle, hasStyle := seqdiagram.StyleNames[*flagStyle] ; hasStyle {
        style = altStyle
    }

    // If there's a process instruction, use it as the target of the diagram
    // TODO: be a little smarter with the process instructions
    for _, pr := range diagram.ProcessingInstructions {
        if pr.Prefix == "goseq" {
            outFilename = pr.Value
        }
    }

    if renderer == nil {
        renderer, err = chooseRendererBaseOnOutfile(outFilename)
        if err != nil {
            return err
        }
    }

    err = renderer(diagram, style, outFilename)
    if err != nil {
        return err
    }

    return nil
}

// Processes a file.  This switches based on the file extension
func processFile(inFilename string, outFilename string, renderer Renderer) error {
    ext := filepath.Ext(inFilename)
    if ext == ".md" {
        return processMdFile(inFilename, outFilename, renderer)
    } else {
        return processSeqFile(inFilename, outFilename, renderer)
    }
}

func main() {    
    var err error

    renderer := SvgRenderer
    outFile := ""

    flag.Parse()

    // Select a suitable renderer (based on the suffix of the output file, if there is one)
    if *flagOut != "" {
        renderer, err = chooseRendererBaseOnOutfile(*flagOut)
        if err != nil {
            die(err.Error())
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