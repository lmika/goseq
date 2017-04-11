package seqdiagram

import (
    "errors"

    "github.com/lmika/goseq/seqdiagram/graphbox"
)

const (
    // DejaVuSans - https://fontlibrary.org/en/font/dejavu-sans
    dejaVuSansFont = "DejaVuSans"
)

// Attempt to load an internal font
func loadInternalFont(fontName string) (*graphbox.TTFFont, error) {
    originalFilename := fontName + ".ttf"
    if fontDataSlice, hasFile := embeddedfiles[originalFilename]; hasFile {
        return graphbox.NewTTFFontFromByteSlice(fontDataSlice, fontName)
    } else {
        return nil, errors.New("No such embedded font: " + fontName)
    }
}