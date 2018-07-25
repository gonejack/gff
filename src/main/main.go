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
	dry      bool
	patterns []string
})

func init() {
	log.SetFlags(log.Ltime)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
`Usage: gff [-b basedir] patterns..
Example: gff '*.txt'
`, )

		flag.PrintDefaults()
	}
	argument.dir, _ = os.Getwd()
	dry := flag.Bool("dry", false, "Dry running")
	dir := flag.String("dir", argument.dir, "Working directory")
	sep := flag.String("sep", "_", "Separator")

	flag.Parse()

	argument.dir, _ = filepath.Abs(*dir)
	argument.sep = *sep
	argument.dry = *dry
	argument.patterns = flag.Args()

	if !strings.HasSuffix(argument.dir, osSep) {
		argument.dir += osSep
	}
}

func check() (pass bool) {
	if len(argument.patterns) == 0 {
		return
	}
	if _, e := os.Stat(argument.dir); os.IsNotExist(e) {
		log.Printf("Dir not exist: %s", argument.dir)
		return
	}

	return true
}

func walk() (hasFound bool) {
	for _, pt := range argument.patterns {
		filepath.Walk(argument.dir, func(p string, info os.FileInfo, e error) error {
			if info.Mode().IsRegular() {
				match, _ := filepath.Match(pt, filepath.Base(p))

				if match {
					hasFound = true

					n := strings.TrimPrefix(p, argument.dir)
					n = strings.Replace(n, osSep, argument.sep, -1)
					n = argument.dir + n

					if p == n {
						log.Printf("Skip: %s", n)
					} else {
						log.Printf("Rename: %s => %s", p, n)

						if !argument.dry {
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

func main() {
	if check() {
		if walk() {
			log.Printf("Done.")
		} else {
			log.Printf("No matches.")
		}
	} else {
		flag.Usage()
	}
}
