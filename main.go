package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"flag"
	"strings"
)

const (
	LastLine = "└──"
	Line = "├──"
)

func main() {
	showAll := flag.Bool("all", false, "show hidden files and directories")
	maxDepth := flag.Int("depth", 2, "maximum directory depth")
	path := flag.String("path", "./", "path")
	flag.Parse()

	if *maxDepth < 0 {
		log.Fatal("-depth must be non-negative")
	}

	if flag.NArg() != 0 {
		log.Fatal("unexpected arguments; use -path to specify a directory")
	}

	fmt.Println(filepath.Clean(*path))
	err := printDir(*path, 0, "", *maxDepth, *showAll)
	if err != nil {
		log.Fatal(err)
	}
}

func printDir(path string, depth int, prefix string, maxDepth int, showAll bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var visible []os.DirEntry

	for _, e := range entries {
		if !showAll && strings.HasPrefix(e.Name(), ".") {
			continue
		}

		visible = append(visible, e)
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

				if err := printDir(childPath, depth+1, childPrefix, maxDepth, showAll); err != nil {
					return err
				}
			}
		} else {
			fmt.Printf("%v\n", e.Name())
		}
	}

	return nil
}
