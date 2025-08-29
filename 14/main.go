package main

import (
	"fmt"
	"sync"
	"time"
)

var or func(channels ...<-chan interface{}) <-chan interface{} = orFunc

func main() {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Second),
		sig(1*time.Second),
		sig(1*time.Second),
		sig(1*time.Second),
		sig(2*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done in %v\n", time.Since(start))
}

func orFunc[T any](
	channels ...<-chan T,
) <-chan T {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	out := make(chan T)

	once := sync.Once{}
	closeFunc := func() {
		once.Do(func() {
			close(out)
		})
	}

	for _, channel := range channels {
		go func(ch <-chan T) {
			defer closeFunc()

			for {
				select {
				case <-ch:
					return
				case <-out:
					return
				}
			}
		}(channel)
	}

	return out
}
