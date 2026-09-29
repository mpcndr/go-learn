//go:build !solution

package concurrency

import "time"

// ParallelSum หาผลรวมของ nums โดยแบ่งเป็น parts ส่วน แต่ละส่วนให้ goroutine ของตัวเองรวม
// แล้วค่อยรวมผลทุกส่วนเข้าด้วยกัน (ถ้า parts <= 0 ให้ถือว่า parts = 1)
// คำใบ้: ให้แต่ละ goroutine ส่งผลรวมย่อยเข้า channel หรือเขียนลง slice คนละ index
func ParallelSum(nums []int, parts int) int {
	// TODO
	return 0
}

// Inc เพิ่มค่าของ key ขึ้น 1
func (c *SafeCounter) Inc(key string) {
	// TODO: Lock, สร้าง map ถ้ายังเป็น nil, เพิ่มค่า, Unlock
}

// Value คืนค่าปัจจุบันของ key
func (c *SafeCounter) Value(key string) int {
	// TODO
	return 0
}

// WorkerPool ประมวลผล inputs ด้วย fn โดยใช้ goroutine ไม่เกิน workers ตัวพร้อมกัน
// ผลลัพธ์ต้องเรียงตามลำดับของ inputs (results[i] = fn(inputs[i]))
// คำใบ้: ส่ง "index" ของงานผ่าน channel แล้วให้ worker เขียนผลลง results[i]
func WorkerPool(inputs []int, workers int, fn func(int) int) []int {
	// TODO
	return nil
}

// Pipeline สร้าง 3 ขั้นต่อกันด้วย channel:
//
//	gen(nums) → square → sum
//
// โดย gen และ square ต้องเป็น goroutine ที่ close channel ขาออกของตัวเองเมื่อเสร็จ
// Pipeline คืนผลรวมของกำลังสองของทุกตัว
func Pipeline(nums []int) int {
	// TODO
	return 0
}

// FirstResponse เรียก fetchers ทุกตัวพร้อมกัน แล้วคืนผลลัพธ์แรกที่ได้
// ถ้าเกิน timeout ให้คืน ("", false)
//
// ⚠️ ระวัง goroutine leak: fetcher ที่ช้าจะยังพยายามส่งผลหลังเรา return ไปแล้ว
// ใช้ channel ที่มี buffer = len(fetchers) เพื่อให้พวกมันส่งได้โดยไม่ block ค้างตลอดไป
func FirstResponse(fetchers []func() string, timeout time.Duration) (string, bool) {
	// TODO: select ระหว่าง channel ผลลัพธ์ กับ time.After(timeout)
	return "", false
}
