package main

import (
//    "log"
    "io"
    "fmt"
    "bufio"
    "bytes"
    "strings"
    "unicode"
)

// A markdown filter.  This will read a markdown file and will search for codeblocks
// starting with "#!goseq".  It will interpret this as a call to goseq, parse the content
// and build an image.
type MarkdownFilter struct {
    input       io.Reader
    output      io.Writer
    handler     CodeblockHandler
}

type CodeblockHandler func(codeblock string, output io.Writer) error


func (cb *MarkdownFilter) Scan() error {
    scanner := bufio.NewScanner(cb.input)
    inblock := false
    currentIndent := 0
    blockcontent := new(bytes.Buffer)

    for scanner.Scan() {
        line := scanner.Text()
        trimmedLine := strings.TrimSpace(line)
        indent := cb.lineIndent(line)

        if (!inblock) && (indent >= currentIndent + 4) && (strings.HasPrefix(trimmedLine, "#!goseq")) {
            inblock = true
            currentIndent = indent
        } else if (inblock) && (indent <= currentIndent - 4) && (trimmedLine != "") {
            inblock = false
            cb.handler(blockcontent.String(), cb.output)
            blockcontent.Reset()
            currentIndent = indent
        }

        if inblock {
            fmt.Fprintln(blockcontent, line)
        } else {
            fmt.Fprintln(cb.output, line)
        }

        //log.Println(inblock, line)
    }

    if inblock {
        cb.handler(blockcontent.String(), cb.output)
    }

    return scanner.Err()
}

// Determine the line indent
func (cb *MarkdownFilter) lineIndent(line string) int {
    indent := 0
    for _, r := range line {
        if (unicode.IsSpace(r)) {
            indent++
        } else {
            break
        }
    }
    return indent
}
    