//go:build !solution

// Package functions — แบบฝึกหัดบท 02
package functions

// MinMax คืนค่าน้อยสุดและมากสุดใน nums
// ถ้า nums ว่าง ให้คืน ok = false
func MinMax(nums []int) (min, max int, ok bool) {
	// TODO
	return 0, 0, false
}

// MakeCounter คืนฟังก์ชันที่แต่ละครั้งที่ถูกเรียกจะคืนค่าเพิ่มขึ้นทีละ 1 เริ่มจาก 1
// counter แต่ละตัวต้องนับแยกกัน
func MakeCounter() func() int {
	// TODO
	return nil
}

// Apply เรียก fn กับทุกตัวใน nums แล้วคืน slice ผลลัพธ์ใหม่ (ห้ามแก้ nums เดิม)
func Apply(nums []int, fn func(int) int) []int {
	// TODO
	return nil
}

// Compose คืนฟังก์ชันที่ทำ f ก่อน แล้วนำผลไปเข้า g — คือ g(f(x))
func Compose(f, g func(int) int) func(int) int {
	// TODO
	return nil
}

// Memoize ห่อ fn ด้วย cache: ถ้าเคยคำนวณ x แล้ว ให้คืนค่าเดิมโดยไม่เรียก fn ซ้ำ
func Memoize(fn func(int) int) func(int) int {
	// TODO: สร้าง map[int]int ไว้นอก closure
	return nil
}
