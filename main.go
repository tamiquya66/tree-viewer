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

func main() {
	dirsFirst := flag.Bool("dirs-first", false, "shows the directories first")
	exclude := flag.String("exclude", "", "exclude files and directories")
	showAll := flag.Bool("all", false, "show hidden files and directories")
	maxDepth := flag.Int("depth", 2, "maximum directory depth")
	path := flag.String("path", "./", "path")
	flag.Parse()

	names := make(map[string]bool)

	for _, name := range strings.Split(*exclude, ",") {
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

	fmt.Println(filepath.Clean(*path))
	err := printDir(*path, 0, "", *maxDepth, *showAll, names, *dirsFirst)
	if err != nil {
		log.Fatal(err)
	}
}

func printDir(
	path string, depth int,
	prefix string, maxDepth int,
	showAll bool, names map[string]bool,
	dirsFirst bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read directory %q: %w", path, err)
	}

	var visible []os.DirEntry

	for _, e := range entries {
		if !showAll && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if names[e.Name()] {
			continue
		}

		visible = append(visible, e)
	}

	if dirsFirst {
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
			if depth < maxDepth {
				childPrefix := prefix

				if isLast {
					childPrefix += "    "
				} else {
					childPrefix += "│   "
				}

				if err := printDir(childPath, depth+1,
					childPrefix, maxDepth,
					showAll, names,
					dirsFirst); err != nil {
					log.Printf("warning: %v", err)
				}
			}
		} else {
			fmt.Printf("%v\n", e.Name())
		}
	}

	return nil
}