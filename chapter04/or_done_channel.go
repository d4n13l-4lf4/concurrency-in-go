package main

import (
	"fmt"
	"math/rand"
	"time"
)

func orDone(done, c <- chan interface{}) <- chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <- done:
				return
			case v, ok := <- c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <- done:
				}
			}
		}
	}()

	return valStream
}

func fakeMain() {

	done := make(chan interface{})
	myChan := make(chan interface{})

	go func() {
		for {
			select {
			case <- done:
				return
				case myChan <- rand.Int():
			}
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		close(done)
	}()

	for val := range orDone(done, myChan) {
		fmt.Printf("Value: %v\n", val)
	}

}
