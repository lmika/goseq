// Locates font on the system.  This is depending on the GOOS
//

package goseq


// The locateFont function is the function that the appropriate fontlocator_* file
// is to implement.
//type FontLocator func() string


// This global function is set by the apropriate fontlocator_* file depending on
// the GOOS.
//var fontLocatorFn FontLocator


// The locateFont function returns a suitble filename of a font that can be
// used by the graphics builder.  If the empty string is returned, no font can
// be found.
//func LocateFont() string {
    //if fontLocatorFn != nil {
//        return fontLocatorFn()
    //} else {
    //    return ""
    //}
//}