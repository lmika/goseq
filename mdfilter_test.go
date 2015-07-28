package main

import (
    "strings"
    "bytes"
    "testing"
    "io"

    "github.com/seanpont/assert"
)

func TestFilter1(t *testing.T) {
    assert := assert.Assert(t)
    md := `
This is normal markdown

    #!goseq
    This is some seq diagram

This is normal markdown again`

    expMd := `
This is normal markdown

This is normal markdown again
`

    blocks, actual := runFilter(md)

    assert.Equal("[" + actual + "]", "[" + expMd + "]")
    assert.Equal(len(blocks), 1)
    assert.Equal(blocks[0], "    #!goseq\n    This is some seq diagram\n\n")
}


func TestFilter2(t *testing.T) {
    assert := assert.Assert(t)
    md := `
This is normal markdown

    This is a standard code block that does nothing.

Back to normal markdown

    #!goseq
    This is some seq diagram
    
    More sequence diagram

This is normal markdown again
asdasdasdasdasdasdasd

    #!goseq
    Seq diagram again
`

    expMd := `
This is normal markdown

    This is a standard code block that does nothing.

Back to normal markdown

This is normal markdown again
asdasdasdasdasdasdasd

`

    blocks, actual := runFilter(md)

    assert.Equal("[" + actual + "]", "[" + expMd + "]")
    assert.Equal(len(blocks), 2)
    assert.Equal(blocks[0], "    #!goseq\n    This is some seq diagram\n    \n    More sequence diagram\n\n")
    assert.Equal(blocks[1], "    #!goseq\n    Seq diagram again\n")
}


func runFilter(input string) (blocks []string, output string) {
    bufout := new(bytes.Buffer)
    mf := &MarkdownFilter{strings.NewReader(input), bufout, func(codeblock string, output io.Writer) error {
        blocks = append(blocks, codeblock)
        return nil
    }}
    mf.Scan()

    output = bufout.String()
    return
}