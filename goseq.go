package main

import (
    "strings"
    "fmt"
    "os"
    "./goseq"
//    "./goseq/graphbox"
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
        diagram.WriteSVG(os.Stdout)
    }

        /*
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
        */
    //}

/*
    gb := graphbox.NewGraphic(3, 3)
    gb.Put(0, 0, &graphbox.LifeLine{2, 0})
    gb.Put(0, 1, &graphbox.LifeLine{2, 1})
    gb.Put(0, 0, &graphbox.ActorRect{150, 100, "Hello"})
    gb.Put(2, 1, &graphbox.ActorRect{100, 25, "World"})
    gb.Put(2, 0, &graphbox.ActorRect{100, 50, "Example"})
    gb.Put(0, 1, &graphbox.ActorRect{25, 25, "Here"})
    gb.Put(1, 2, &graphbox.ActorRect{20, 250, "Spacer"})
    gb.DrawSVG(os.Stdout)
*/    
}