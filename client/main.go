package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.Println("ahoy!")
	log.Println(os.Environ())
	log.Println(os.Getwd())
	walkFunc := func(dir string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			log.Println(dir)
		}
		return nil
	}
	if err := filepath.Walk("/mnt/ext1", walkFunc); err != nil {
		log.Fatal(err)
	}
}
