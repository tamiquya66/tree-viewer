package main

import (
    "bytes"
    "os"
    "path/filepath"
    "testing"
)

func TestPrintDir(t *testing.T) {
    dir := t.TempDir()

    for _, name := range []string{"a.txt", "b.txt", ".hidden"} {
        path := filepath.Join(dir, name)
        if err := os.WriteFile(path, nil, 0644); err != nil {
            t.Fatal(err)
        }
    }

    var out bytes.Buffer

    err := printDir(&out, dir, 0, "", options{})
    if err != nil {
        t.Fatal(err)
    }

    want := "├──a.txt\n└──b.txt\n"
    got := out.String()

    if got != want {
        t.Errorf("unexpected tree:\ngot:  %q\nwant: %q", got, want)
    }
}