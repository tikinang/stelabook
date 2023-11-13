package main

import (
	"io/fs"
	"log"
	"path/filepath"
)

func main() {
	walkFunc := func(path string, info fs.FileInfo, err error) error {

	}
	if err := filepath.Walk("/", walkFunc); err != nil {
		log.Fatal(err)
	}
}
