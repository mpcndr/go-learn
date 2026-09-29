// Package generics — แบบฝึกหัดบท 08
package generics

// Number — constraint สำหรับตัวเลขที่บวกกันได้
type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

// Set เก็บค่าไม่ซ้ำ — map ที่ value เป็น struct{} ไม่กินหน่วยความจำ
type Set[T comparable] struct {
	m map[T]struct{}
}

// Cache เก็บ key → value แบบจำกัดจำนวน เมื่อเต็มให้ลบตัวที่ถูกใส่ "เก่าที่สุด" ออก (FIFO)
type Cache[K comparable, V any] struct {
	capacity int
	data     map[K]V
	order    []K // ลำดับที่ key ถูกใส่เข้ามา
}
