package main

import (
	"fmt"
	"time"
)

func DoWork(done <-chan interface{}, vals <-chan string) <-chan interface{} {
	stopped := make(chan interface{})
	go func() {
		defer close(stopped)
		for {
			select {
			case val := <-vals:
				fmt.Println("Processing val: " + val)
			case  <-done:
				fmt.Println("Worker is shutting down")
				return
			}
		}
	}()
	return stopped
}


func main() {
	fmt.Println("Starting...")

	done := make(chan interface{})
	vals := make(chan string)

	stopped := DoWork(done, vals)

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Sending shutdown...")
		close(done)
	}()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Sending s2...")
		vals <- "s2"
	}()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Sending s1...")
		vals <- "s1"
	}()

	<-stopped
	fmt.Println("Finished, exiting")
}
