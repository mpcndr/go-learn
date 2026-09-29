//go:build solution

package ctxex

import (
	"context"
	"sync"
	"time"
)

func SlowSquare(ctx context.Context, n int, delay time.Duration) (int, error) {
	select {
	case <-time.After(delay):
		return n * n, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func WithRetry(ctx context.Context, attempts int, wait time.Duration, fn func(context.Context) error) error {
	var err error
	for i := range attempts {
		if err = fn(ctx); err == nil {
			return nil
		}
		if i == attempts-1 {
			break
		}
		t := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
		}
	}
	return err
}

func FetchAll(ctx context.Context, ids []int, fetch func(context.Context, int) (string, error)) (map[int]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		once     sync.Once
		firstErr error
		results  = make(map[int]string, len(ids))
	)
	for _, id := range ids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := fetch(ctx, id)
			if err != nil {
				once.Do(func() {
					firstErr = err
					cancel()
				})
				return
			}
			mu.Lock()
			results[id] = v
			mu.Unlock()
		}()
	}
	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	return results, nil
}

func Generator(ctx context.Context) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; ; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}
