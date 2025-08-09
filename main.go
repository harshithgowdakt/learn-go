package main

import "fmt"

func main() {
	ch := make(chan int, 2)

	ch <- 1
	ch <- 2

	close(ch)
	// Reading from closed channel returns zero value + false
	v1, ok1 := <-ch // v1=1, ok1=true
	v2, ok2 := <-ch // v2=2, ok2=true
	v3, ok3 := <-ch // v3=0, ok3=false (zero value)

	fmt.Println(v1)
	fmt.Println(v2)
	fmt.Println(v3)

	fmt.Println(ok1)
	fmt.Println(ok2)
	fmt.Println(ok3)

	// Sending to closed channel panics!
	// ch <- 3  // panic: send on closed channel

	// Closing already closed channel panics!
	// close(ch)  // panic: close of closed channel
}
