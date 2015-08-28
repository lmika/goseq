// Windows font locator implementation

package seqdiagram

import (
    "os"
    "path/filepath"
)


// The font directory
// TODO: Must not hard code to C:\Windows\Font

var winFontDirectory string = "C:\\Windows\\Fonts"

// Desirable fonts
var ttfFonts = []string {
    "calibri.ttf",
    "verdana.ttf",
    "arial.ttf",
}

// Returns the first font found given the directory containing the true
// type fonts.
func locateWinTTFFont(ttfDir string) []string {
    fonts := make([]string, 0)

    for _, fontName := range ttfFonts {
        path := filepath.Join(ttfDir, fontName)
        if stat, _ := os.Stat(path) ; (stat != nil) {
            fonts = append(fonts, path)
        }
    }
    return fonts
}

// Locates an appropriate font on Window
func LocateFont() []string {
    return locateWinTTFFont(winFontDirectory)
}
