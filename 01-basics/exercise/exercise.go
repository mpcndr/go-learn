//go:build !solution

// Package basics — แบบฝึกหัดบท 01
// ทำให้ test ผ่าน:  go test ./01-basics/...
package basics

// FizzBuzz คืน slice ของ string สำหรับเลข 1..n
//   - หาร 3 และ 5 ลงตัว → "FizzBuzz"
//   - หาร 3 ลงตัว      → "Fizz"
//   - หาร 5 ลงตัว      → "Buzz"
//   - นอกนั้น          → ตัวเลขเป็น string (ใช้ strconv.Itoa)
//
// ถ้า n <= 0 ให้คืน slice ว่าง
func FizzBuzz(n int) []string {
	// TODO
	return nil
}

// Grade แปลงคะแนน 0–100 เป็นเกรด
//
//	80–100 → "A", 70–79 → "B", 60–69 → "C", 50–59 → "D", 0–49 → "F"
//	นอกช่วง 0–100 → "invalid"
func Grade(score int) string {
	// TODO: ลองใช้ switch แบบไม่มีเงื่อนไข (switch { case ...: })
	return ""
}

// SumDigits คืนผลรวมของเลขแต่ละหลัก เช่น 1234 → 10
// ถ้าเป็นเลขลบให้คิดจากค่าสัมบูรณ์ เช่น -45 → 9
func SumDigits(n int) int {
	// TODO: ใช้ n % 10 เพื่อเอาหลักสุดท้าย และ n / 10 เพื่อตัดหลักสุดท้ายทิ้ง
	return 0
}

// IsPrime ตรวจว่า n เป็นจำนวนเฉพาะหรือไม่ (n < 2 ไม่ใช่จำนวนเฉพาะ)
func IsPrime(n int) bool {
	// TODO: คำใบ้ — ลองหารแค่ถึง i*i <= n ก็พอ
	return false
}
