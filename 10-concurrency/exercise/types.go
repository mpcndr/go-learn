// Package concurrency — แบบฝึกหัดบท 10
// รัน test ด้วย -race เสมอ:  go test -race ./10-concurrency/...
package concurrency

import "sync"

// SafeCounter นับจำนวนต่อ key — ต้องปลอดภัยเมื่อถูกเรียกจากหลาย goroutine
// zero value ต้องใช้ได้ทันที (var c SafeCounter)
type SafeCounter struct {
	mu     sync.Mutex
	counts map[string]int
}
