package main

import (
    "log"
    "flag"
)

func init() {
    log.SetFlags(0)
}

func main() {
    g := NewGff()

    if g.check() {
        if !g.yes {
            log.Println("Dry Run:")
        }
        if g.walk() {
            log.Printf("Done.")
        } else {
            log.Printf("No matches.")
        }
    } else {
        flag.Usage()
    }
}