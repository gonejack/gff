package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
)

type command struct {
	Separator string   `short:"s" default:"_" help:"Filename separator."`
	Yes       bool     `short:"y" help:"Make real changes."`
	Patterns  []string `arg:"" optional:""`
}

func (c *command) run() {
	kong.Parse(c,
		kong.Name("gff"),
		kong.Description("Extract files from nested directory https://github.com/gonejack/gff"),
		kong.UsageOnError(),
	)

	if !c.Yes {
		log.Println("changes preview:")
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot read working directory: %s", err)
		return
	}

	if c.walk(wd) {
		log.Printf("done.")
	} else {
		log.Printf("no matches.")
	}
}
func (c *command) walk(dir string) (hasMatches bool) {
	if !filepath.IsAbs(dir) {
		dir, _ = filepath.Abs(dir)
	}

	for _, pattern := range c.Patterns {
		filepath.Walk(dir, func(file string, info os.FileInfo, e error) (err error) {
			if !info.Mode().IsRegular() {
				return
			}

			match, _ := filepath.Match(pattern, filepath.Base(file))
			if !match {
				return
			}

			hasMatches = true

			rename := strings.TrimPrefix(file, dir)
			rename = strings.TrimPrefix(rename, string(filepath.Separator))
			rename = strings.Replace(rename, string(filepath.Separator), c.Separator, -1)
			rename = filepath.Join(dir, rename)

			if file == rename {
				log.Printf("skip %s", rename)
				return nil
			}

			log.Printf("rename %s => %s", file, rename)

			if c.Yes {
				err = os.Rename(file, rename)
				if err != nil {
					log.Printf("error: %s", err)
					return err
				}
			}

			return
		})
	}

	return
}

func main() {
	new(command).run()
}
