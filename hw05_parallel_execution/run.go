package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func consumer(wg *sync.WaitGroup, m int, counter *int32, ch chan Task) {
	defer wg.Done()

	for t := range ch {
		res := t()
		if res != nil {
			atomic.AddInt32(counter, 1)
		}

		if int(atomic.LoadInt32(counter)) > m {
			return
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var counter int32
	var result error
	wg := &sync.WaitGroup{}
	ch := make(chan Task, len(tasks)+1)

	if m < 1 {
		return ErrErrorsLimitExceeded
	}

	wg.Add(n)

	for i := 0; i < len(tasks); i++ {
		ch <- tasks[i]
	}
	close(ch)

	for i := 0; i < n; i++ {
		go consumer(wg, m, &counter, ch)
	}

	wg.Wait()

	if int(counter) > m {
		result = ErrErrorsLimitExceeded
	}

	return result
}
