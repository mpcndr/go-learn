//go:build !solution

// Package testingex — แบบฝึกหัดบท 09
package testingex

// Slugify แปลงข้อความเป็น slug สำหรับ URL
//   - ตัวอักษรภาษาอังกฤษแปลงเป็นตัวเล็ก
//   - เก็บเฉพาะ a-z, 0-9 และตัวอักษรภาษาอื่น (unicode.IsLetter / unicode.IsDigit)
//   - ทุกอย่างที่เหลือ (ช่องว่าง, เครื่องหมาย) ถือเป็นตัวคั่น → แทนด้วย "-" ตัวเดียว
//   - ไม่มี "-" ที่ต้นหรือท้าย และไม่มี "--" ติดกัน
//
// ตัวอย่าง: "  Hello, World!  " → "hello-world"
//
//	"Go 1.24 Release" → "go-1-24-release"
//	"เรียน Go"         → "เรียน-go"
//
// คำใบ้: วนทีละ rune, ใช้ strings.Builder, และจำไว้ว่า "ตัวก่อนหน้าเป็นตัวคั่นหรือเปล่า"
// ระวัง: สระและวรรณยุกต์ไทย (ั ี ่ ้) ไม่ใช่ Letter แต่เป็น Mark → ใช้ unicode.IsMark ด้วย
func Slugify(s string) string {
	// TODO
	return ""
}
