// $ go test -bench=. -cpu=1 mesurecontextswitch_test.go                                                                            (master) mesure
// goos: darwin
// goarch: amd64
// BenchmarkContextSwitch  10000000               186 ns/op
// PASS
// ok      command-line-arguments  2.073s

package main

import (
	"sync"
	"testing"
)

func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}
	receiver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}

	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer()
	close(begin)
	wg.Wait()
}
