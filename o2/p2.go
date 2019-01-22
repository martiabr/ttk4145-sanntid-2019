// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
)

// Control signals
const (
	GetNumber = iota  // 0
	Exit  // 1
)

func number_server(add_number <-chan int, control <-chan int, number chan<- int) {
	var i = 0

	// This for-select pattern is one you will become familiar with if you're using go "correctly".
	for {
		select {
		// Receive different messages and handle them correctly
		// You will at least need to update the number and handle control signals.

		case j := <-add_number:
			i += j
		case c := <-control:
			if c == GetNumber {
				number <- i
			} else if c == Exit {
				return
			}
		}
	}
}

func incrementing(add_number chan<-int, finished chan<- bool) {
	for j := 0; j<1000420; j++ {
		add_number <- 1
	}
	//Signal that the goroutine is finished
	finished <- true
}

func decrementing(add_number chan<- int, finished chan<- bool) {
	for j := 0; j<1000000; j++ {
		add_number <- -1
	}
	//Signal that the goroutine is finished
	finished <- true
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Construct the required channels
	// Think about whether the receptions of the number should be unbuffered, or buffered with a fixed queue size. ???
	num_ch := make(chan int)
	add_ch := make(chan int)
	control_ch := make(chan int)
	finished_inc_ch := make(chan bool)
	finished_dec_ch := make(chan bool)


	// Spawn the required goroutines
	go number_server(add_ch, control_ch, num_ch)
	go incrementing(add_ch, finished_inc_ch);
	go decrementing(add_ch, finished_dec_ch);

	// Block on finished from both "worker" goroutines
	sum := 0
	for sum < 2 {
		select {
		case <- finished_inc_ch:
			sum++
		case <- finished_dec_ch:
			sum++
		}
	}

	control_ch<-GetNumber
	Println("The magic number is:", <- num_ch)
	control_ch<-Exit
}