package main

import (
	"flag"
	"fmt"

	"resource-hash/domain"
	"resource-hash/internal/component"
	"resource-hash/pkg/job"
)

func main() {
	var (
		limitArg         uint = 10
		inputFilenameArg string
	)
	flag.UintVar(&limitArg, "concurrency", limitArg, "Specify parallel jobs count")
	flag.StringVar(&inputFilenameArg, "filename", inputFilenameArg, "Specify links filename")
	flag.Parse()

	linksCh, err := component.ReadLinksList(inputFilenameArg, limitArg)
	if err != nil {
		panic(err)
	}

	outputCh := make(chan domain.OutputChunk)
	go func() {
		for chunk := range outputCh {
			fmt.Printf("%s -> %s\n", chunk.Url, chunk.Hash)
		}
	}()

	jobsCh := make(chan job.Job, limitArg)
	go func() {
		for link := range linksCh {
			jobsCh <- job.Job{
				Name: fmt.Sprintf("process %s", link),
				Func: component.NewCheckLinkJob(link, outputCh),
			}
		}
		close(jobsCh)
	}()

	report := job.AsyncLimitedJobRunner(limitArg, jobsCh)
	for link, err := range report {
		fmt.Printf("ERROR: job '%s' failed : %v\n", link, err)
	}
}
