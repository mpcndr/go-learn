package concurrency

import (
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestParallelSum(t *testing.T) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i + 1
	}
	for _, parts := range []int{0, 1, 3, 7, 1000, 5000} {
		if got := ParallelSum(nums, parts); got != 500500 {
			t.Errorf("ParallelSum(1..1000, %d) = %d, want 500500", parts, got)
		}
	}
	if got := ParallelSum(nil, 4); got != 0 {
		t.Errorf("ParallelSum(nil) = %d, want 0", got)
	}
}

func TestSafeCounter(t *testing.T) {
	var c SafeCounter
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				c.Inc("x")
			}
		}()
	}
	wg.Wait()
	if got := c.Value("x"); got != 10000 {
		t.Errorf("Value(x) = %d, want 10000", got)
	}
	if got := c.Value("missing"); got != 0 {
		t.Errorf("Value(missing) = %d, want 0", got)
	}
}

func TestWorkerPool(t *testing.T) {
	inputs := []int{1, 2, 3, 4, 5, 6, 7, 8}
	var running, maxRunning atomic.Int32
	fn := func(x int) int {
		n := running.Add(1)
		for {
			m := maxRunning.Load()
			if n <= m || maxRunning.CompareAndSwap(m, n) {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
		running.Add(-1)
		return x * x
	}
	start := time.Now()
	got := WorkerPool(inputs, 3, fn)
	elapsed := time.Since(start)

	if !slices.Equal(got, []int{1, 4, 9, 16, 25, 36, 49, 64}) {
		t.Errorf("WorkerPool = %v, results must keep input order", got)
	}
	if m := maxRunning.Load(); m > 3 {
		t.Errorf("max concurrent workers = %d, want <= 3", m)
	}
	if m := maxRunning.Load(); m < 2 {
		t.Errorf("max concurrent workers = %d — are the workers really running in parallel?", m)
	}
	if elapsed > 150*time.Millisecond {
		t.Errorf("took %v — too slow for 3 workers", elapsed)
	}
}

func TestPipeline(t *testing.T) {
	if got := Pipeline([]int{1, 2, 3, 4}); got != 30 {
		t.Errorf("Pipeline(1..4) = %d, want 30", got)
	}
	if got := Pipeline(nil); got != 0 {
		t.Errorf("Pipeline(nil) = %d, want 0", got)
	}
}

func TestFirstResponse(t *testing.T) {
	sleepy := func(d time.Duration, v string) func() string {
		return func() string { time.Sleep(d); return v }
	}
	before := runtime.NumGoroutine()

	got, ok := FirstResponse([]func() string{
		sleepy(80*time.Millisecond, "slow"),
		sleepy(10*time.Millisecond, "fast"),
		sleepy(50*time.Millisecond, "medium"),
	}, time.Second)
	if !ok || got != "fast" {
		t.Errorf("FirstResponse = %q, %v, want fast, true", got, ok)
	}

	if got, ok := FirstResponse([]func() string{sleepy(200*time.Millisecond, "late")}, 20*time.Millisecond); ok {
		t.Errorf("FirstResponse with timeout = %q, %v, want \"\", false", got, ok)
	}

	// ตรวจ goroutine leak: รอให้ fetcher ที่ช้าทำงานเสร็จ แล้ว goroutine ต้องจบหมด
	time.Sleep(300 * time.Millisecond)
	if after := runtime.NumGoroutine(); after > before {
		t.Errorf("goroutine leak: %d before, %d after", before, after)
	}
}
