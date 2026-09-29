//go:build !solution

package generics

// Map แปลงทุกตัวใน s ด้วย fn
func Map[T, U any](s []T, fn func(T) U) []U {
	// TODO
	return nil
}

// Filter คืนเฉพาะตัวที่ keep(v) == true
func Filter[T any](s []T, keep func(T) bool) []T {
	// TODO
	return nil
}

// Reduce ยุบ s ให้เหลือค่าเดียว เริ่มจาก initial แล้วสะสมด้วย fn(acc, v)
// เช่น Reduce([1,2,3], 0, func(acc, v int) int { return acc + v }) == 6
func Reduce[T, A any](s []T, initial A, fn func(A, T) A) A {
	// TODO
	return initial
}

// Sum รวมตัวเลขทั้งหมด ใช้ได้กับทุก type ใน Number
func Sum[T Number](nums []T) T {
	// TODO
	return 0
}

// GroupBy จัดกลุ่ม s ตาม key ที่ได้จาก keyFn โดยรักษาลำดับเดิมในแต่ละกลุ่ม
// เช่น GroupBy(["apple","avocado","banana"], ตัวอักษรแรก) → {"a": [apple avocado], "b": [banana]}
func GroupBy[T any, K comparable](s []T, keyFn func(T) K) map[K][]T {
	// TODO
	return nil
}

// NewSet สร้าง Set พร้อมค่าเริ่มต้น
func NewSet[T comparable](items ...T) *Set[T] {
	// TODO: อย่าลืม make map
	return nil
}

// Add เพิ่มค่า
func (s *Set[T]) Add(v T) {
	// TODO
}

// Has ตรวจว่ามีค่าหรือไม่
func (s *Set[T]) Has(v T) bool {
	// TODO
	return false
}

// Len จำนวนสมาชิก
func (s *Set[T]) Len() int {
	// TODO
	return 0
}

// Intersect คืน Set ใหม่ที่มีเฉพาะค่าที่อยู่ในทั้ง s และ other
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	// TODO
	return nil
}

// NewCache สร้าง Cache ที่เก็บได้สูงสุด capacity ตัว
func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	// TODO
	return nil
}

// Put ใส่ค่า ถ้า key มีอยู่แล้วให้อัปเดตค่า (ไม่เปลี่ยนลำดับ)
// ถ้า key ใหม่และ cache เต็ม ให้ลบ key ที่เก่าที่สุดออกก่อน
func (c *Cache[K, V]) Put(key K, value V) {
	// TODO
}

// Get คืนค่าและ ok=true ถ้ามี, ถ้าไม่มีคืน zero value ของ V และ ok=false
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// TODO
	var zero V
	return zero, false
}
