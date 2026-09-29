//go:build !solution

// Package ctxex — แบบฝึกหัดบท 11
package ctxex

import (
	"context"
	"time"
)

// SlowSquare คืน n*n หลังจากรอ delay
// ถ้า ctx ถูกยกเลิกก่อนครบ delay ให้คืน (0, ctx.Err()) ทันที
func SlowSquare(ctx context.Context, n int, delay time.Duration) (int, error) {
	// TODO: select ระหว่าง time.After(delay) กับ ctx.Done()
	return 0, nil
}

// WithRetry เรียก fn สูงสุด attempts ครั้ง จนกว่าจะสำเร็จ (คืน nil)
// ระหว่างแต่ละครั้งให้รอ wait
// ถ้า ctx ถูกยกเลิกระหว่างรอ ให้หยุดทันทีและคืน ctx.Err()
// ถ้าครบทุกครั้งแล้วยังพัง ให้คืน error ครั้งสุดท้าย
func WithRetry(ctx context.Context, attempts int, wait time.Duration, fn func(context.Context) error) error {
	// TODO
	return nil
}

// FetchAll เรียก fetch(ctx, id) สำหรับทุก id "พร้อมกัน"
//   - ถ้าสำเร็จทั้งหมด: คืน map[id]ผลลัพธ์ และ nil
//   - ถ้ามีตัวใดพัง: ยกเลิก context ที่ส่งให้ตัวอื่น ๆ ทันที แล้วคืน (nil, error ตัวแรก)
//
// ห้าม return ก่อนที่ goroutine ทุกตัวจะจบ (ไม่ leak)
// คำใบ้: context.WithCancel + sync.WaitGroup + sync.Once (เก็บ error ตัวแรก) + sync.Mutex (ป้องกัน map)
func FetchAll(ctx context.Context, ids []int, fetch func(context.Context, int) (string, error)) (map[int]string, error) {
	// TODO
	return nil, nil
}

// Generator ส่ง 0, 1, 2, 3, ... ออกทาง channel ไปเรื่อย ๆ
// เมื่อ ctx ถูกยกเลิก goroutine ต้องจบและ close channel
// ระวัง: ตอนกำลัง "ส่ง" ก็ต้องฟัง ctx.Done() ด้วย ไม่งั้นจะค้างถ้าไม่มีใครรับแล้ว
func Generator(ctx context.Context) <-chan int {
	// TODO
	return nil
}
