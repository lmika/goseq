package main

import (
    "strings"
    "fmt"
    "os"
    "./goseq"
    "./goseq/graphbox"
)

func main() {
    diagram, err := goseq.Parse(strings.NewReader(`
        title: ABC123    is the    best  
        participant c
        participant b
        participant a

        a->b: Does something
        b->c: Does something as well
        c -> d: Another thing
        `))
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Printf("Diagram: [%s]\n", diagram.Title)
        for _, p := range diagram.Actors {
            fmt.Printf("Participant: [%s]\n", p.Name)
        }
        for _, i := range diagram.Items {
            switch ip := i.(type) {
            case *goseq.Action:
                fmt.Printf("Action from %s to %s: %s\n", ip.From.Name, ip.To.Name, ip.Message)
            }
        }
    }


    gb := graphbox.NewGraphic(2, 2)
    gb.Put(0, 0, &graphbox.Rectangle{100, 100})
    gb.Put(1, 1, &graphbox.Rectangle{100, 100})
    gb.Put(1, 0, &graphbox.Rectangle{100, 100})
    gb.DrawSVG(os.Stdout)
}