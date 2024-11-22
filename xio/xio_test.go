package xio

import (
	"sync"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	const (
		size       = 10024
		goroutines = 100
	)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	defer wg.Wait()

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			var buf = Get(size)
			defer Put(buf)

			if len(buf) != size {
				t.Error("buf len not is", size)
				return
			}

			time.Sleep(time.Millisecond * 10)
		}()
	}
}
