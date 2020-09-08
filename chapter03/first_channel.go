package main

import (
	"fmt"
	"sync"
)

/* A channel is blocking by itself */
var waitGroup sync.WaitGroup

func main() {

	stringStream := make(chan string)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		stringStream <- "Hello channel!"
	}()

	fmt.Println(<-stringStream)
	waitGroup.Wait()
	close(stringStream)
}
