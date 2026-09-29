package ctxex

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestSlowSquare(t *testing.T) {
	ctx := context.Background()
	if v, err := SlowSquare(ctx, 4, 5*time.Millisecond); err != nil || v != 16 {
		t.Errorf("SlowSquare = %d, %v, want 16, nil", v, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := SlowSquare(ctx, 4, time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want context.DeadlineExceeded", err)
	}
	if time.Since(start) > 200*time.Millisecond {
		t.Error("SlowSquare did not stop promptly on cancel")
	}
}

func TestWithRetry(t *testing.T) {
	var calls int
	err := WithRetry(context.Background(), 5, time.Millisecond, func(context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("flaky")
		}
		return nil
	})
	if err != nil || calls != 3 {
		t.Errorf("succeed on 3rd: err=%v calls=%d, want nil, 3", err, calls)
	}

	calls = 0
	err = WithRetry(context.Background(), 3, time.Millisecond, func(context.Context) error {
		calls++
		return fmt.Errorf("fail #%d", calls)
	})
	if err == nil || err.Error() != "fail #3" || calls != 3 {
		t.Errorf("always fail: err=%v calls=%d, want 'fail #3', 3", err, calls)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	calls = 0
	start := time.Now()
	err = WithRetry(ctx, 100, 20*time.Millisecond, func(context.Context) error {
		calls++
		return errors.New("down")
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("cancelled retry err = %v, want DeadlineExceeded", err)
	}
	if time.Since(start) > 200*time.Millisecond || calls > 5 {
		t.Errorf("retry did not stop on cancel (calls=%d)", calls)
	}
}

func TestFetchAllSuccess(t *testing.T) {
	fetch := func(ctx context.Context, id int) (string, error) {
		time.Sleep(20 * time.Millisecond)
		return fmt.Sprintf("item-%d", id), nil
	}
	start := time.Now()
	got, err := FetchAll(context.Background(), []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, fetch)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(got) != 10 || got[3] != "item-3" {
		t.Errorf("got %v", got)
	}
	if time.Since(start) > 120*time.Millisecond {
		t.Errorf("took %v — fetches should run concurrently", time.Since(start))
	}
}

func TestFetchAllCancelsOthers(t *testing.T) {
	var cancelled atomic.Int32
	boom := errors.New("id 2 failed")
	fetch := func(ctx context.Context, id int) (string, error) {
		if id == 2 {
			time.Sleep(10 * time.Millisecond)
			return "", boom
		}
		select {
		case <-time.After(time.Second):
			return "slow", nil
		case <-ctx.Done():
			cancelled.Add(1)
			return "", ctx.Err()
		}
	}
	start := time.Now()
	got, err := FetchAll(context.Background(), []int{1, 2, 3, 4}, fetch)
	if !errors.Is(err, boom) {
		t.Errorf("err = %v, want the first error (boom)", err)
	}
	if got != nil {
		t.Errorf("results = %v, want nil on error", got)
	}
	if time.Since(start) > 300*time.Millisecond {
		t.Error("FetchAll did not cancel the other fetches")
	}
	if n := cancelled.Load(); n != 3 {
		t.Errorf("cancelled fetches = %d, want 3 (FetchAll must wait for all goroutines)", n)
	}
}

func TestGenerator(t *testing.T) {
	before := runtime.NumGoroutine()
	ctx, cancel := context.WithCancel(context.Background())
	ch := Generator(ctx)
	if ch == nil {
		t.Fatal("Generator returned nil channel")
	}
	for want := range 5 {
		if got := <-ch; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
	}
	cancel() // หยุดโดยไม่อ่านต่อ — goroutine ต้องไม่ค้าง

	deadline := time.After(time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				time.Sleep(10 * time.Millisecond)
				if after := runtime.NumGoroutine(); after > before {
					t.Errorf("goroutine leak: %d before, %d after", before, after)
				}
				return
			}
		case <-deadline:
			t.Fatal("channel was not closed after cancel")
		}
	}
}
