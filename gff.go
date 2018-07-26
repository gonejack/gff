package main

import (
    "log"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

type gff struct {
    dir      string
    sep      string
    yes      bool
    osSep    string
    patterns []string
}

func (g *gff) check() (pass bool) {
    st, e := os.Stat(g.dir)
    switch true {
    case !st.IsDir():
        log.Printf("Error: %s is not a directory", g.dir)
    case os.IsNotExist(e):
        log.Printf("Error: directory %s do not exist", g.dir)
    case len(g.patterns) == 0:
        log.Println("Error: not pattern assigned")
    default:
        pass = true
    }

    return
}

func (g *gff) walk() (hasAny bool) {
    for _, pt := range g.patterns {
        filepath.Walk(g.dir, func(p string, info os.FileInfo, e error) error {
            if info.Mode().IsRegular() {
                match, _ := filepath.Match(pt, filepath.Base(p))

                if match {
                    hasAny = true

                    n := strings.TrimPrefix(p, g.dir)
                    n = strings.Replace(n, g.osSep, g.sep, -1)
                    n = g.dir + n

                    if p == n {
                        log.Printf("Skip: %s", n)
                    } else {
                        log.Printf("Rename: %s => %s", p, n)

                        if g.yes {
                            if e := os.Rename(p, n); e != nil {
                                log.Printf("Error: %s", e)

                                return e
                            }
                        }
                    }
                }
            }

            return nil
        })
    }

    return
}

func NewGff() *gff {
    flag.Usage = func() {
        fmt.Fprintf(flag.CommandLine.Output(), "Usage: gff [options] patterns..\nExample: gff -yes '*.txt'\n")
        flag.PrintDefaults()
    }

    wd, _ := os.Getwd()
    dir := flag.String("dir", wd, "Work Directory")
    sep := flag.String("sep", "_", "Separator")
    yes := flag.Bool("yes", false, "Make Real Changes")

    flag.Parse()

    g := new(gff)
    g.dir, _ = filepath.Abs(*dir)
    g.sep = *sep
    g.yes = *yes
    g.osSep = string(filepath.Separator)
    g.patterns = flag.Args()

    if !strings.HasSuffix(g.dir, g.osSep) {
        g.dir += g.osSep
    }

    return g
}
