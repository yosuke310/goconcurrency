package main

import (
	"fmt"
	"math/rand"
)

func main() {
	done := make(chan interface{})
	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}

func doWork(done <-chan interface{}) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	results := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(results)

		for i := 0; i < 10; i++ {
			select {
			case heartbeat <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case results <- rand.Intn(10):
			}
		}
	}()

	return heartbeat, results
}
