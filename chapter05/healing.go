package main

import (
	"log"
	"os"
	"time"
)

type startGoroutineFn func(
	done <- chan interface{},
	pulseInterval time.Duration,
) (heartbeat <- chan interface{})

func steward(
	timeout time.Duration,
	startGoroutine startGoroutineFn,
) startGoroutineFn {
	return func(
		done <- chan interface{},
		pulseInterval time.Duration,
	) (<- chan interface{}) {
		heartbeat := make(chan interface{})
		go func() {
			defer close(heartbeat)

			var wardDone chan interface{}
			var wardHeartbeat <- chan interface{}
			startWard := func() {
				wardDone = make(chan interface{})
				wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2) /* import or */
			}
			startWard()
			pulse := time.Tick(pulseInterval)

		monitorLoop:
			for {
				timeoutSignal := time.After(timeout)

				for {
					select{
					case <- pulse:
						select {
						case heartbeat <- struct {}{}:
						default:
						}
					case <- wardHeartbeat:
						continue monitorLoop
					case <- timeoutSignal:
						log.Println("steward: ward unhealthy; restarting")
						close(wardDone)
						startWard()
						continue monitorLoop
					case <- done:
						return
					}
				}
			}
		}()
		return heartbeat
	}
}

func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWork := func(done <- chan interface{}, _ time.Duration) <- chan interface{} {
		log.Println("ward: Hello, I'm irresponsible!")
		go func() {
			<- done
			log.Println("ward: I am halting.")
		}()
		return nil
	}

	doWorkWithSteward := steward(4 * time.Second, doWork)

	done := make(chan interface{})
	time.AfterFunc(9 * time.Second, func() {
		log.Println("main: halting steward and ward.")
		close(done)
	})

	for range doWorkWithSteward(done, 4 * time.Second) {}
	log.Println("Done")
}
