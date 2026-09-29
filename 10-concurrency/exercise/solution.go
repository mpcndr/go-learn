//go:build solution

package concurrency

import (
	"sync"
	"time"
)

func ParallelSum(nums []int, parts int) int {
	if parts <= 0 {
		parts = 1
	}
	size := (len(nums) + parts - 1) / parts
	sums := make(chan int, parts)
	var wg sync.WaitGroup
	for start := 0; start < len(nums); start += size {
		end := min(start+size, len(nums))
		wg.Add(1)
		go func(chunk []int) {
			defer wg.Done()
			s := 0
			for _, n := range chunk {
				s += n
			}
			sums <- s
		}(nums[start:end])
	}
	wg.Wait()
	close(sums)

	total := 0
	for s := range sums {
		total += s
	}
	return total
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.counts == nil {
		c.counts = make(map[string]int)
	}
	c.counts[key]++
}

func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[key]
}

func WorkerPool(inputs []int, workers int, fn func(int) int) []int {
	if workers <= 0 {
		workers = 1
	}
	results := make([]int, len(inputs))
	jobs := make(chan int)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				results[i] = fn(inputs[i]) // แต่ละ index มีคนเขียนคนเดียว → ไม่ race
			}
		}()
	}
	for i := range inputs {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return results
}

func Pipeline(nums []int) int {
	gen := func() <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				out <- n
			}
		}()
		return out
	}
	square := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * n
			}
		}()
		return out
	}
	total := 0
	for v := range square(gen()) {
		total += v
	}
	return total
}

func FirstResponse(fetchers []func() string, timeout time.Duration) (string, bool) {
	results := make(chan string, len(fetchers))
	for _, f := range fetchers {
		go func() { results <- f() }()
	}
	select {
	case r := <-results:
		return r, true
	case <-time.After(timeout):
		return "", false
	}
}
