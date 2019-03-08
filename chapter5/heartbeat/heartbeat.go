package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan interface{})
	time.AfterFunc(15*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}

func doWork(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)

	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.NewTicker(pulseInterval)
		workGen := time.NewTicker(pulseInterval * 2)
		defer func() {
			fmt.Println("tickers stopping.")
			pulse.Stop()
			workGen.Stop()
		}()

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}
		sendResult := func(r time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse.C:
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		for {
			select {
			case <-done:
				return
			case <-pulse.C:
				sendPulse()
			case r := <-workGen.C:
				sendResult(r)
			}
		}
	}()

	return heartbeat, results
}
