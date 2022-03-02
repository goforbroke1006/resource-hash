package component

import (
	"fmt"

	"resource-hash/domain"
	"resource-hash/pkg/job"
	"resource-hash/pkg/links"
)

func NewApplication(filename string, concurrency uint) *application {
	return &application{
		filename:    filename,
		concurrency: concurrency,
	}
}

type application struct {
	filename    string
	concurrency uint
}

// Run starts 3-step pipeline
//   1 - read AsyncLimitedJobRunner's jobs output
//   2 - transform link to job
//   3 - show error list
func (app application) Run() {
	linksCh, err := links.ReadLinksList(app.filename, app.concurrency)
	if err != nil {
		panic(err)
	}

	outputCh := make(chan domain.OutputChunk)
	go func() {
		for chunk := range outputCh {
			if chunk.Warn != nil {
				fmt.Printf("WARN: job '%s' rejected: %v\n", chunk.Url, chunk.Warn)
			} else {
				fmt.Printf("%s -> %s\n", chunk.Url, chunk.Hash)
			}
		}
	}()

	jobsCh := make(chan job.Job, app.concurrency)
	go func() {
		for link := range linksCh {
			jobsCh <- job.Job{
				Name: fmt.Sprintf("process %s", link),
				Func: NewCheckLinkJob(link, outputCh),
			}
		}
		close(jobsCh)
	}()

	if err := job.AsyncLimitedJobRunner(app.concurrency, jobsCh); err != nil {
		fmt.Printf("ERROR: job runner stop with err: %v", err)
	}
	close(outputCh)
}
