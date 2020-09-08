package main

import (
	"fmt"
	"time"
)

func main() {

	letterStream := make(chan string, 10)

	letterStreamer := func () {
		for {
			select {
			case letterStream <- "a":
			case <- time.After(1 * time.Second):
				close(letterStream)
				fmt.Println("Finishing streamer...")
				return
			}
		}
	}

	go letterStreamer()

	for letter := range letterStream {
		fmt.Printf("The letter was %v\n", letter)
	}

	fmt.Println("Finishing program")
}
