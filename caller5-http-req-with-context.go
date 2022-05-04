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
	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	respSize := 0
	for scanner.Scan() {
		respSize += len(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s - %d", resp.Status, respSize), nil
}

func foo(ctx context.Context) (string, error) {
	return httpGet(ctx, "http://localhost:8080/foo")
}

func bar(ctx context.Context) (string, error) {
	return httpGet(ctx, "http://localhost:8080/bar")
}

// https://pkg.go.dev/net/http#Client
// https://pkg.go.dev/net/http#NewRequestWithContext
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
		fooRes, fooErr = foo(ctx)
		wg.Done()
	}()

	var barRes string
	var barErr error
	wg.Add(1)
	go func() {
		barRes, barErr = bar(ctx)
		wg.Done()
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
