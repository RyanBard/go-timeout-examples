package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Config struct {
	Foo          FooConfig
	Bar          BarConfig
	TotalTimeout time.Duration
}

type FooConfig struct {
	NumIterations         int64
	MaxStepsBeforeCheckIn int
	Timeout               time.Duration
}

type FooService struct {
	Config FooConfig
}

func newFooService(config FooConfig) *FooService {
	return &FooService{
		Config: config,
	}
}

func (s *FooService) DoWork(parentCtx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(parentCtx, s.Config.Timeout)
	defer cancel()
	var i int64
	steps := 0
	ans := 0
	for i < s.Config.NumIterations {
		if steps > s.Config.MaxStepsBeforeCheckIn {
			select {
			case <-time.After(1 * time.Millisecond):
				steps = 0
			case <-ctx.Done():
				return -1, ctx.Err()
			}
		}
		i++
		steps++
		ans += 10
		ans *= 123567
		ans /= 123567
		ans *= 123567
	}
	return ans, nil
}

type BarConfig struct {
	SleepTime time.Duration
	Timeout   time.Duration
}

type BarService struct {
	Config BarConfig
}

func newBarService(config BarConfig) *BarService {
	return &BarService{
		Config: config,
	}
}

func (s *BarService) DoWork(parentCtx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(parentCtx, s.Config.Timeout)
	defer cancel()
	ansChan := make(chan string)
	go func() {
		time.Sleep(s.Config.SleepTime)
		ansChan <- "success"
		close(ansChan)
	}()
	select {
	case ans := <-ansChan:
		return ans, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	config := Config{
		TotalTimeout: 30 * time.Second,
		Foo: FooConfig{
			NumIterations:         10000000,
			MaxStepsBeforeCheckIn: 1000,
			Timeout:               20 * time.Second,
		},
		Bar: BarConfig{
			SleepTime: 1 * time.Second,
			Timeout:   10 * time.Second,
		},
	}
	fooSvc := newFooService(config.Foo)
	barSvc := newBarService(config.Bar)

	ctx, cancel := context.WithTimeout(context.Background(), config.TotalTimeout)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	var fooRes int
	var fooErr error
	go func() {
		defer wg.Done()
		fooRes, fooErr = fooSvc.DoWork(ctx)
	}()

	wg.Add(1)
	var barRes string
	var barErr error
	go func() {
		defer wg.Done()
		barRes, barErr = barSvc.DoWork(ctx)
	}()

	wg.Wait()
	fmt.Printf("done: fooRes=%v, fooErr=%v, barRes=%v, barErr=%v\n", fooRes, fooErr, barRes, barErr)
}
