package job

import (
	"context"
	"fmt"
	"sync"
)

func AsyncLimitedJobRunner(limit uint, jobs <-chan Job) error {
	limiter := make(chan struct{}, limit)
	wg := sync.WaitGroup{}

	resultErr := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for job := range jobs {
			limiter <- struct{}{} // reserve
			wg.Add(1)

			go func(ctx context.Context, job Job) {
				defer func() {
					<-limiter // release
					wg.Done()
				}()

				if err := job.Func(); err != nil {
					resultErr <- fmt.Errorf("job '%s' failed: %v", job.Name, err)
				} else {
					resultErr <- nil
				}
			}(ctx, job)
		}

		wg.Wait()
		close(resultErr)
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case lErr := <-resultErr:
			if lErr != nil {
				cancel()
				return lErr
			}
		}
	}
}
