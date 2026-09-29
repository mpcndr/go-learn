//go:build !solution

// Package pointers — แบบฝึกหัดบท 04
package pointers

// Swap สลับค่าที่ a และ b ชี้อยู่
func Swap(a, b *int) {
	// TODO
}

// Increment เพิ่มค่าทุกตัวใน nums ขึ้น 1 "ในที่เดิม" (in place)
// คำถามให้คิด: ทำไมฟังก์ชันนี้รับ []int ไม่ใช่ *[]int แต่ยังแก้ข้อมูลของผู้เรียกได้?
func Increment(nums []int) {
	// TODO
}

// AppendTo เพิ่ม values ต่อท้าย slice ที่ s ชี้อยู่ โดยผู้เรียกต้องเห็น len ใหม่
func AppendTo(s *[]int, values ...int) {
	// TODO
}

// ResetScores ตั้ง Score ของนักเรียนทุกคนใน map ให้เป็น 0
//
// กับดัก: m[k].Score = 0 จะ compile ไม่ผ่าน เพราะค่าใน map "ไม่สามารถอ้างที่อยู่ได้"
// ต้องดึงออกมา แก้ แล้วใส่กลับ — หรือเปลี่ยน map ให้เก็บ *Student (แต่ในข้อนี้ห้ามเปลี่ยน signature)
func ResetScores(m map[string]Student) {
	// TODO
}

// Ptr คืน pointer ที่ชี้ไปยังสำเนาของ v
// มีประโยชน์เวลาต้องการ pointer ของค่าคงที่ เช่น Ptr(42) — เพราะเขียน &42 ไม่ได้
func Ptr[T any](v T) *T {
	// TODO
	return nil
}
