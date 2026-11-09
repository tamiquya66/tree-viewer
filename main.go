package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const LastLine = "└──"
const Line = "├──"
const maxDepth = 2

func main() {
	err := printDir("./", 0, "")
	if err != nil {
		log.Fatal(err)
	}
}

func printDir(path string, depth int, prefix string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i, e := range entries {
		isLast := i == len(entries)-1
		childPath := filepath.Join(path, e.Name())

		fmt.Print(prefix)
		if !isLast {
			fmt.Print(Line)
		} else {
			fmt.Print(LastLine)
		}
		// TODO: ignore .* files
		if e.IsDir() {
			fmt.Printf("%v/\n", e.Name())
			if depth < maxDepth {
				childPrefix := prefix

				if isLast {
					childPrefix += "	"
				} else {
					childPrefix += "│   "
				}

				if err := printDir(childPath, depth+1, childPrefix); err != nil {
					return err
				}
			}
		} else {
			fmt.Printf("%v\n", e.Name())
		}
	}

	return nil
}
