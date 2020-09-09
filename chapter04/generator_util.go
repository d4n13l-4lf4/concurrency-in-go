package main

import "sync"

func repeat(
	done <- chan interface{},
	values ... interface{},
) <- chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <- done:
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

func take(
	done <- chan interface{},
	valueStream <- chan interface{},
	num int,
) <- chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <- done:
				return
			case takeStream <- <- valueStream:
			}
		}
	}()
	return takeStream
}

func repeatFn (
	done <-chan interface{},
	fn func() interface{},
) <- chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <- done:
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

func toInt(
	done <- chan interface{},
	valueStream <- chan interface{},
) <- chan int {
	newStream := make(chan int)

	go func() {
		defer close(newStream)
		for v := range valueStream {
			select {
			case <- done:
				return
			case newStream <- v.(int):
			}
		}
	}()

	return newStream
}

func fanIn(
	done <- chan interface{},
	channels ...<-chan interface{},
) <- chan interface{} {

	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})
	multiplex := func(c <- chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <- done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream

}