// บท 11: Context — รันด้วย: go run ./11-context
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "request-id"

// queryDB จำลองงานที่ใช้เวลา — เคารพการยกเลิกด้วย select
func queryDB(ctx context.Context, query string, took time.Duration) (string, error) {
	reqID, _ := ctx.Value(requestIDKey).(string)
	select {
	case <-time.After(took):
		return fmt.Sprintf("[%s] ผลของ %q", reqID, query), nil
	case <-ctx.Done():
		return "", fmt.Errorf("[%s] query %q: %w", reqID, query, ctx.Err())
	}
}

func handleRequest(parent context.Context, id string, dbTime time.Duration) {
	ctx, cancel := context.WithTimeout(parent, 100*time.Millisecond)
	defer cancel()
	ctx = context.WithValue(ctx, requestIDKey, id)

	res, err := queryDB(ctx, "SELECT * FROM books", dbTime)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Println("  ⏱  timeout:", err)
	case errors.Is(err, context.Canceled):
		fmt.Println("  🛑 ถูกยกเลิก:", err)
	case err != nil:
		fmt.Println("  ❌", err)
	default:
		fmt.Println("  ✅", res)
	}
}

// ticker ที่หยุดเมื่อ context ถูกยกเลิก — ไม่มี goroutine leak
func heartbeat(ctx context.Context) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		t := time.NewTicker(20 * time.Millisecond)
		defer t.Stop()
		for i := 1; ; i++ {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				select {
				case out <- i:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

func main() {
	root := context.Background()

	fmt.Println("=== timeout ===")
	handleRequest(root, "req-1", 30*time.Millisecond)  // ทันเวลา
	handleRequest(root, "req-2", 500*time.Millisecond) // เกิน 100ms

	fmt.Println("\n=== ยกเลิกจากพ่อ → ลูกหยุดด้วย ===")
	parent, cancelParent := context.WithCancel(root)
	go func() {
		time.Sleep(20 * time.Millisecond)
		fmt.Println("  (ผู้ใช้ปิดหน้าเว็บ → cancelParent())")
		cancelParent()
	}()
	handleRequest(parent, "req-3", 500*time.Millisecond)

	fmt.Println("\n=== generator ที่หยุดได้ ===")
	ctx, cancel := context.WithTimeout(root, 110*time.Millisecond)
	defer cancel()
	for beat := range heartbeat(ctx) {
		fmt.Print("  ♥", beat)
	}
	fmt.Println("\n  generator จบเอง เพราะ context หมดเวลา:", ctx.Err())
}
