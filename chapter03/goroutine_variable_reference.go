package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		/* Should pass salutation var to prevent taking the last value "good day" */
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()
}
