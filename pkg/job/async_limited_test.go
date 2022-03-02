package job

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestAsyncLimitedJobRunner(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		jobs := make(chan Job)
		var counter uint32
		go func() {
			jobs <- Job{
				Name: "fake 1 - no errors",
				Func: func() error {
					atomic.AddUint32(&counter, 1)
					return nil
				},
			}
			jobs <- Job{
				Name: "fake 2 - no errors",
				Func: func() error {
					atomic.AddUint32(&counter, 1)
					return nil
				},
			}
			jobs <- Job{
				Name: "fake 3 - no errors",
				Func: func() error {
					atomic.AddUint32(&counter, 1)
					return nil
				},
			}
			close(jobs)
		}()
		if err := AsyncLimitedJobRunner(1, jobs); err != nil {
			t.Errorf("should not emit error")
		}
		if counter != 3 {
			t.Errorf("all jobs should be finished")
		}
	})

	t.Run("errors collected", func(t *testing.T) {
		jobs := make(chan Job)
		go func() {
			jobs <- Job{
				Name: "fake 1",
				Func: func() error {
					return errors.New("fake 1 error")
				},
			}
			jobs <- Job{
				Name: "fake 2",
				Func: func() error {
					return errors.New("fake 2 error")
				},
			}
			jobs <- Job{
				Name: "fake 3 - no errors",
				Func: func() error {
					return nil
				},
			}
			close(jobs)
		}()
		if err := AsyncLimitedJobRunner(1, jobs); err == nil {
			t.Errorf("should collect error")
		}
	})

	t.Run("limited works count", func(t *testing.T) {
		total := 10

		terminators := make([]chan struct{}, total)
		for i := 0; i < total; i++ {
			terminators[i] = make(chan struct{})
		}

		var startedJobsCounter int64

		jobs := make(chan Job, total)
		go func() {
			for ji := 0; ji < total; ji++ {
				jobs <- Job{
					Name: fmt.Sprintf("%d", ji),
					Func: func(index int) func() error {
						return func() error {
							atomic.AddInt64(&startedJobsCounter, 1)
							<-terminators[index] // lock job, wait for termination
							return nil
						}
					}(ji),
				}
			}
		}()

		go func() {
			// detach to prevent deadlock
			_ = AsyncLimitedJobRunner(2, jobs)
		}()
		time.Sleep(time.Second)

		if startedJobsCounter != 2 {
			t.Errorf("wrong started job count, got = %d, want = %d", startedJobsCounter, 2)
		}

		terminators[0] <- struct{}{} // set free 0-th goroutine
		time.Sleep(time.Second)

		if startedJobsCounter != 3 {
			t.Errorf("wrong started job count, got = %d, want = %d", startedJobsCounter, 3)
		}
	})
}
