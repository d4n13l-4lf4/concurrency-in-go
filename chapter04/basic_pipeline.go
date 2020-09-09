package main

import "fmt"

func main() {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(done <- chan interface{}, intStream <- chan int, multiplier int,) <- chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range(intStream) {
				select {
				case <- done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4, 5)
	pipeline := multiply(done, intStream, 4)

	for v := range pipeline {
		fmt.Println(v)
	}
}
