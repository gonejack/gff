package main

import (
    "os"
    "log"
    "path/filepath"
    "strings"
    "flag"
    "fmt"
)

var osSep = string(filepath.Separator)

var argument = new(struct {
    dir      string
    sep      string
    yes      bool
    patterns []string
})

func init() {
    log.SetFlags(log.Ltime)

    flag.Usage = func() {
        fmt.Fprintf(flag.CommandLine.Output(), "Usage: gff [options] patterns..\nExample: gff -yes '*.txt'\n")

        flag.PrintDefaults()
    }
    argument.dir, _ = os.Getwd()
    dir := flag.String("dir", argument.dir, "Work Directory")
    sep := flag.String("sep", "_", "Separator")
    yes := flag.Bool("yes", false, "Make Real Changes")

    flag.Parse()

    argument.dir, _ = filepath.Abs(*dir)
    argument.sep = *sep
    argument.yes = *yes
    argument.patterns = flag.Args()

    if !strings.HasSuffix(argument.dir, osSep) {
        argument.dir += osSep
    }
}

func check() (pass bool) {
    pass = true

    if len(argument.patterns) == 0 {
        pass = false
    }
    if _, e := os.Stat(argument.dir); os.IsNotExist(e) {
        log.Printf("Dir not exist: %s", argument.dir)

        pass = false
    }

    if pass {
        if !argument.yes {
            log.Println("Dry Running")
        }
    } else {
        flag.Usage()
    }

    return pass
}

func walk() (hasAny bool) {
    for _, pt := range argument.patterns {
        filepath.Walk(argument.dir, func(p string, info os.FileInfo, e error) error {
            if info.Mode().IsRegular() {
                match, _ := filepath.Match(pt, filepath.Base(p))

                if match {
                    hasAny = true

                    n := strings.TrimPrefix(p, argument.dir)
                    n = strings.Replace(n, osSep, argument.sep, -1)
                    n = argument.dir + n

                    if p == n {
                        log.Printf("Skip: %s", n)
                    } else {
                        log.Printf("Rename: %s => %s", p, n)

                        if argument.yes {
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

    if hasAny {
        log.Printf("Done.")
    } else {
        log.Printf("No matches.")
    }

    return
}

func main() {
    if check() {
        walk()
    }
}
