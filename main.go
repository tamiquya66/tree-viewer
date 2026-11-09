package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	LastLine = "└──"
	Line     = "├──"
)

type options struct {
	maxDepth int
	showAll bool
	dirsFirst bool
	excluded map[string]bool
}

func main() {
	dirsFirst := flag.Bool("dirs-first", false, "shows the directories first")
	excluded := flag.String("exclude", "", "exclude files and directories")
	showAll := flag.Bool("all", false, "show hidden files and directories")
	maxDepth := flag.Int("depth", 2, "maximum directory depth")
	path := flag.String("path", "./", "path")
	flag.Parse()

	names := make(map[string]bool)

	for _, name := range strings.Split(*excluded, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		names[name] = true
	}

	if *maxDepth < 0 {
		log.Fatal("-depth must be non-negative")
	}

	if flag.NArg() != 0 {
		log.Fatal("unexpected arguments; use -path to specify a directory")
	}

	opts := options{
		maxDepth: *maxDepth,
		showAll: *showAll,
		dirsFirst: *dirsFirst,
		excluded: names,
	}

	fmt.Println(filepath.Clean(*path))
	err := printDir(*path, 0, "", opts)
	if err != nil {
		log.Fatal(err)
	}
}

func printDir(path string, depth int, prefix string, opts options) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read directory %q: %w", path, err)
	}

	var visible []os.DirEntry

	for _, e := range entries {
		if !opts.showAll && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if opts.excluded[e.Name()] {
			continue
		}

		visible = append(visible, e)
	}

	if opts.dirsFirst {
		sort.SliceStable(visible, func(i, j int) bool {
			return visible[i].IsDir() && !visible[j].IsDir()
		})
	}

	for i, e := range visible {
		isLast := i == len(visible)-1
		childPath := filepath.Join(path, e.Name())

		fmt.Print(prefix)
		if !isLast {
			fmt.Print(Line)
		} else {
			fmt.Print(LastLine)
		}

		if e.IsDir() {
			fmt.Printf("%v/\n", e.Name())
			if depth < opts.maxDepth {
				childPrefix := prefix

				if isLast {
					childPrefix += "    "
				} else {
					childPrefix += "│   "
				}

				if err := printDir(childPath, depth+1,childPrefix, opts); err != nil {
					log.Printf("warning: %v", err)
				}
			}
		} else {
			fmt.Printf("%v\n", e.Name())
		}
	}

	return nil
}