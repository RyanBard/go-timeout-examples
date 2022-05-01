package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func httpGet(url string, resChan chan string, errChan chan error) {
	defer close(resChan)
	defer close(errChan)

	resp, err := http.Get(url)
	if err != nil {
		errChan <- err
		return
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	respSize := 0
	for scanner.Scan() {
		respSize += len(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		errChan <- err
		return
	}
	resChan <- fmt.Sprintf("%s - %d", resp.Status, respSize)
}

func foo(resChan chan string, errChan chan error) {
	httpGet("http://localhost:8080/foo", resChan, errChan)
}

func bar(resChan chan string, errChan chan error) {
	httpGet("http://localhost:8080/bar", resChan, errChan)
}

// https://gobyexample.com/
// https://devhints.io/go
func main() {
	fooResChan := make(chan string)
	fooErrChan := make(chan error)
	barResChan := make(chan string)
	barErrChan := make(chan error)

	go foo(fooResChan, fooErrChan)
	go bar(barResChan, barErrChan)

	var fooRes string
	var fooErr error
	select {
	case fooRes = <-fooResChan:
	case fooErr = <-fooErrChan:
	case <-time.After(2 * time.Second):
		fooErr = errors.New("timeout!")
	}

	var barRes string
	var barErr error
	select {
	case barRes = <-barResChan:
	case barErr = <-barErrChan:
	case <-time.After(2 * time.Second):
		barErr = errors.New("timeout!")
	}

	if fooErr == nil {
		fmt.Println("foo succeeded:", fooRes)
	} else {
		fmt.Println("foo failed:", fooErr)
	}

	if barErr == nil {
		fmt.Println("bar succeeded:", barRes)
	} else {
		fmt.Println("bar failed:", barErr)
	}
}
