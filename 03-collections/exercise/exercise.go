//go:build !solution

// Package collections — แบบฝึกหัดบท 03
package collections

// WordCount นับจำนวนครั้งที่แต่ละคำปรากฏ
// ไม่สนตัวพิมพ์เล็ก/ใหญ่ ("Go" กับ "go" คือคำเดียวกัน) และแยกคำด้วยช่องว่าง
// คำใบ้: strings.Fields, strings.ToLower
func WordCount(text string) map[string]int {
	// TODO
	return nil
}

// Reverse กลับลำดับตัวอักษร ต้องรองรับ Unicode เช่น "สวัสดี" → "ีดสัวส"
// คำใบ้: แปลงเป็น []rune ก่อน
func Reverse(s string) string {
	// TODO
	return ""
}

// Unique คืน slice ใหม่ที่ไม่มีตัวซ้ำ โดยรักษาลำดับการปรากฏครั้งแรก
// เช่น [3 1 3 2 1] → [3 1 2]
func Unique(nums []int) []int {
	// TODO: ใช้ map[int]bool (หรือ map[int]struct{}) เป็น set
	return nil
}

// Chunk แบ่ง nums เป็นก้อนละ size ตัว ก้อนสุดท้ายอาจน้อยกว่า
// เช่น Chunk([1 2 3 4 5], 2) → [[1 2] [3 4] [5]]
// ถ้า size <= 0 ให้คืน nil
func Chunk(nums []int, size int) [][]int {
	// TODO
	return nil
}

// TopN คืนคำที่พบบ่อยที่สุด n คำ เรียงจากมากไปน้อย
// ถ้าจำนวนเท่ากัน ให้เรียงตามตัวอักษร (a → z)
// ใช้ WordCount ที่เขียนไว้ด้านบนได้เลย
// คำใบ้: slices.SortFunc
func TopN(text string, n int) []string {
	// TODO
	return nil
}
