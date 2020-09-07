package main

import (
	"fmt"
	"sync"
)

func main() {

	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating a new instance")
			return struct{}{}
		},
	}

	myPool.Get()

	instance := myPool.Get()
	/* You should put the new instance back in the pool, otherwise the pool does not make sense */
	myPool.Put(instance)
	myPool.Get()
}
