package graphbox

import (
    "fmt"
)

// Returns the maximum of two integer.
func maxInt(x, y int) int {
    if (x > y) {
        return x
    } else {
        return y
    }
}


// A SVG style
type SvgStyle       map[string]string

func (ss SvgStyle) Set(key, value string) {
    ss[key] = value
}

func (ss SvgStyle) ToStyle() string {
    s := ""
    for k, v := range ss {
        s += fmt.Sprintf("%s:%s;", k, v)
    }
    return s
}