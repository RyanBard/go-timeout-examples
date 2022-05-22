package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func httpGet(parentCtx context.Context, url string) (string, error) {
	resChan := make(chan string)
	errChan := make(chan error)

	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer cancel()

	go func() {
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
	}()

	select {
	case res := <-resChan:
		return res, nil
	case err := <-errChan:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func foo(ctx context.Context) (string, error) {
	return httpGet(ctx, "http://localhost:8080/foo")
}

func bar(ctx context.Context) (string, error) {
	return httpGet(ctx, "http://localhost:8080/bar")
}

// https://gobyexample.com/
// https://devhints.io/go
// https://www.sohamkamani.com/golang/context-cancellation-and-values/
func main() {
	var wg sync.WaitGroup

	ctx := context.Background()

	var fooRes string
	var fooErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		fooRes, fooErr = foo(ctx)
	}()

	var barRes string
	var barErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		barRes, barErr = bar(ctx)
	}()

	wg.Wait()

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
