package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	sayHello := func() {
		defer wg.Done()
		fmt.Println("Something cool!")
	}
	wg.Add(1)
	go sayHello()
	wg.Wait()
}
