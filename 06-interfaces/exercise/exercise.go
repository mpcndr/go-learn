//go:build !solution

package interfaces

// --- 1. ทำให้ทั้ง 3 type เป็น Shape ---
// Name() ให้คืน "circle", "rectangle", "triangle"
// คำใบ้: math.Pi

// TODO: func (c Circle) Area() float64 { ... }
// TODO: func (c Circle) Name() string  { ... }
// TODO: Rectangle, Triangle

// TotalArea รวมพื้นที่ของทุก shape
func TotalArea(shapes []Shape) float64 {
	// TODO
	return 0
}

// Largest คืน shape ที่พื้นที่มากที่สุด ถ้า slice ว่างคืน nil
func Largest(shapes []Shape) Shape {
	// TODO
	return nil
}

// Describe ใช้ type switch อธิบายค่า v:
//
//	int        → "int: <ค่า>"          เช่น "int: 42"
//	string     → "string of length N"  เช่น "string of length 2"
//	Shape      → "shape: <Name()>"     เช่น "shape: circle"
//	error      → "error: <Error()>"
//	อื่น ๆ     → "unknown"
func Describe(v any) string {
	// TODO
	return ""
}

// Write ทำให้ *CountingWriter เป็น io.Writer:
// ส่งต่อไปให้ W แล้วบวกจำนวน byte ที่เขียนได้เข้า Count แล้วคืนค่าเดิมของ W
// (ลบคอมเมนต์ข้างล่างออกแล้วเขียนเนื้อใน)
//
// func (c *CountingWriter) Write(p []byte) (int, error) {
// }

// NotifyAll ส่ง msg ให้ทุกคนใน users ผ่าน s
// ถ้ามีคนส่งไม่สำเร็จ ให้ส่งคนที่เหลือต่อจนครบ แล้วคืนจำนวนที่ส่งสำเร็จ และ error ตัวแรกที่เจอ
func NotifyAll(s Sender, users []string, msg string) (sent int, err error) {
	// TODO
	return 0, nil
}
