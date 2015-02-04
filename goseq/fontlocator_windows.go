// Windows font locator implementation

package goseq

import (
    "os"
    "path/filepath"
)


// The font directory
// TODO: Must not hard code to C:\Windows\Font

var winFontDirectory string = "C:\\Windows\\Fonts"

// Desirable fonts
var ttfFonts = []string {
    "verdana.ttf",
    "arial.ttf",
}

// Returns the first font found given the directory containing the true
// type fonts.
func locateWinTTFFont(ttfDir string) string {
    for _, fontName := range ttfFonts {
        path := filepath.Join(ttfDir, fontName)
        if stat, _ := os.Stat(path) ; (stat != nil) {
            return path
        }
    }
    return ""
}

// Locates an appropriate font on Window
func locateWindowsFont() string {
    return locateWinTTFFont(winFontDirectory)
}

func init() {
    fontLocatorFn = locateWindowsFont
}
