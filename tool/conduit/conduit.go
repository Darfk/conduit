package main

import (
	"github.com/darfk/conduit"
	"os"
	"strings"
	"log"
)

func main() {
	path := os.Getenv("GOFILE")
	if path == "" {
		log.Fatal("GOFILE is empty")
	}

	var outputPath string

	index := strings.LastIndex(path, ".")
	if index == -1 {
		outputPath = path + "_conduit"
	} else {
		outputPath = path[0:index] + "_conduit" + path[index:]
	}

	err := conduit.CreateConduitFile(outputPath, path)
	if err != nil {
		log.Fatal(err)
	}
}
