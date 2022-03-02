package main

import (
	"flag"
	"resource-hash/internal/component"
)

func main() {
	var (
		limitArg         uint = 10
		inputFilenameArg string
	)
	flag.UintVar(&limitArg, "concurrency", limitArg, "Specify parallel jobs count")
	flag.StringVar(&inputFilenameArg, "filename", inputFilenameArg, "Specify links filename")
	flag.Parse()

	application := component.NewApplication(inputFilenameArg, limitArg)
	application.Run()
}
