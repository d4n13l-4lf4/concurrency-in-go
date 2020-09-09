package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func primeFinder(done chan interface{}, randIntStream <- chan int) <- chan interface{} {
	intStream := make(chan interface{})
	go func() {
		defer close(intStream)
		for v := range randIntStream {
			select {
			case <- done:
				return
			case intStream <- v:
			}
		}
	}()
	return intStream
}

func main() {
	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	rand := func() interface{} { return rand.Intn(50000000) }
	randIntStream := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([] <- chan interface{}, numFinders)
	fmt.Println("Primes:")

	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))
}