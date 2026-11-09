package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"path/filepath"
)

func main() {
	err := printDir("./", 0)
	if err != nil {
		log.Fatal(err)
	}
}

func printDir(path string, depth int) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		indent := strings.Repeat(" ", depth)
		childPath := filepath.Join(path, e.Name())

		if e.IsDir(){
			fmt.Printf("%v[D] %v\n", indent, e.Name())
			if depth < 1	{
				if err := printDir(childPath, depth + 1); err != nil {
					return err
				}
			}
		} else {
			fmt.Printf("%v[F] %v\n", indent, e.Name())
		}
	}

	return nil
}