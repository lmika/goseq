// Linux font locator implementation

package seqdiagram

import (
    "os"
    "path/filepath"
)


// TODO: This directory is not definitive and will not work for
// all systems.  Will need to know exactly how each dist. organises
// fonts
var ttfFontDirs = []string {
    "/usr/share/fonts/truetype",
    "/usr/share/fonts",
}

// Desirable fonts
var ttfFonts = []string {
    "freefont/FreeSans.ttf",
    "dejavu/DejaVuSans.ttf",
}


// Returns the directory which contains all the truetype fonts.  First
// directory found will be returned.
func locateTTFDirectory() string {
    for _, potentialDir := range ttfFontDirs {
        if stat, _ := os.Stat(potentialDir) ; (stat != nil) && (stat.IsDir()) {
            return potentialDir
        }
    }
    return ""
}

// Returns the first font found given the directory containing the true
// type fonts.
func locateTTFFonts(ttfDir string) []string {
    fonts := make([]string, 0)

    for _, fontName := range ttfFonts {
        path := filepath.Join(ttfDir, fontName)
        if stat, _ := os.Stat(path) ; (stat != nil) {
            fonts = append(fonts, path)
        }
    }
    return fonts
}

// Locates an appropriate font on the linux file system
func LocateFonts() []string {
    fontDir := locateTTFDirectory()
    if fontDir != "" {
        return locateTTFFonts(fontDir)
    } else {
        return nil
    }
}
