package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
)

func main() {
	new(gff).run()
}

type gff struct {
	Yes       bool   `short:"y" help:"Make real changes."`
	Separator string `short:"s" default:"_" help:"Separator for file renaming."`
	About     bool
	Patterns  []string `arg:"" help:"for example: *.jpg"`
}

func (c *gff) run() {
	kong.Parse(c,
		kong.Name("gff"),
		kong.Description("Extract files from nested directory."),
		kong.UsageOnError(),
	)
	if c.About {
		fmt.Println("Visit https://github.com/gonejack/gff")
		return
	}
	if !c.Yes {
		log.Println("changes preview:")
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot read directory %s", err)
		return
	}
	if c.walk(dir) {
		log.Printf("done.")
	} else {
		log.Printf("no matches.")
	}
}
func (c *gff) walk(dir string) (any bool) {
	if !filepath.IsAbs(dir) {
		dir, _ = filepath.Abs(dir)
	}
	var n int
	for _, patten := range c.Patterns {
		filepath.Walk(dir, func(name string, info os.FileInfo, e error) (err error) {
			if !info.Mode().IsRegular() {
				return
			}
			match, _ := filepath.Match(patten, filepath.Base(name))
			if !match {
				return
			}
			n += 1
			rename := strings.TrimPrefix(strings.TrimPrefix(name, dir), string(filepath.Separator))
			rename = strings.Replace(rename, string(filepath.Separator), c.Separator, -1)
			rename = filepath.Join(dir, rename)
			if name == rename {
				log.Printf("skip %s", rename)
				return
			}
			log.Printf("rename %s => %s", name, rename)
			if c.Yes {
				err = os.Rename(name, rename)
				if err != nil {
					log.Printf("rename %s => %s failed: %s", name, rename, err)
					return
				}
			}
			return
		})
	}
	return n > 0
}
