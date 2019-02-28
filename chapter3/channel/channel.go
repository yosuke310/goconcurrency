package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

func main() {
	bufferCh()
	randomCh()
	timeoutCh()
}

func bufferCh() {
	var stdoutBuf bytes.Buffer
	defer stdoutBuf.WriteTo(os.Stdout)

	ch := make(chan int, 1)
	go func() {
		defer close(ch)
		defer fmt.Fprintln(&stdoutBuf, "Producer done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuf, "Sending: %d\n", i)
			ch <- i
		}
	}()

	for integer := range ch {
		fmt.Fprintf(&stdoutBuf, "Received: %d\n", integer)
	}
}

func randomCh() {
	ch1 := make(chan interface{})
	close(ch1)
	ch2 := make(chan interface{})
	close(ch2)

	var ch1cnt, ch2cnt int
	for i := 0; i < 1000; i++ {
		select {
		case <-ch1:
			ch1cnt++
		case <-ch2:
			ch2cnt++
		}
	}

	fmt.Printf("ch1cnt: %d\nch2cnt: %d\n", ch1cnt, ch2cnt)
}

func timeoutCh() {
	var c <-chan int
	select {
	case <-c:
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}
}
