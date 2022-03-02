package job

import "sync"

type Job struct {
	Name string
	Func func() error
}

func AsyncLimitedJobRunner(limit uint, jobs <-chan Job) (report map[string]error) {
	report = make(map[string]error, len(jobs))
	limiter := make(chan struct{}, limit)
	wg := sync.WaitGroup{}

	for job := range jobs {
		limiter <- struct{}{} // reserve
		wg.Add(1)

		go func(job Job) {
			defer func() {
				<-limiter // release
				wg.Done()
			}()

			if err := job.Func(); err != nil {
				report[job.Name] = err
			}
		}(job)
	}

	wg.Wait()
	return report
}
