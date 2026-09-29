# 09 — Testing: Test คือ Spec ที่รันได้

## แก่นแท้
> **Test ไม่ได้มีไว้ "จับบั๊ก" อย่างเดียว — มันคือเครื่องมือออกแบบ**
> ถ้าโค้ดทดสอบยาก = ออกแบบผูกกันแน่นเกินไป (test ที่ยาก คือสัญญาณเตือนที่ดีที่สุด)

Go มี testing มาในตัว ไม่ต้องลง framework:
```bash
go test ./...                     # ทุก package
go test -v -run TestSlugify/space ./09-testing/...   # เฉพาะ subtest
go test -cover ./...              # ดู coverage
go test -race ./...               # ตรวจ data race
go test -bench=. -benchmem ./09-testing/...          # benchmark
go test -fuzz=FuzzSlugify -fuzztime=10s ./09-testing/exercise   # fuzz
```

## กติกาพื้นฐาน
- ไฟล์ลงท้าย `_test.go` (ไม่ถูก build เข้าโปรแกรมจริง)
- ฟังก์ชัน `func TestXxx(t *testing.T)`
- `t.Errorf` = แจ้งผิดแล้วทำต่อ, `t.Fatalf` = แจ้งผิดแล้วหยุด test นี้

## Table-driven test — pattern หลักของ Go
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```
**ข้อดี:** เพิ่ม case = เพิ่ม 1 บรรทัด, เห็น "spec" ทั้งหมดในที่เดียว, ชื่อ subtest บอกว่าพังตรงไหน

## ข้อความ error ที่ดี
รูปแบบ: `Func(input) = got, want expected`
อ่านแล้วต้องรู้ทันทีว่า **เรียกอะไร ด้วยอะไร ได้อะไร ควรได้อะไร** โดยไม่ต้องเปิดโค้ด test

## เครื่องมือที่ควรรู้
| เครื่องมือ | ใช้ทำอะไร |
|-----------|-----------|
| `t.Run` | subtest — แยกชื่อ, รันเฉพาะตัวได้ |
| `t.Helper()` | ใส่ในฟังก์ชันช่วย ให้ error ชี้บรรทัดของ test จริง |
| `t.Cleanup(fn)` | เก็บกวาดหลัง test จบ (แทน defer ในฟังก์ชันช่วย) |
| `t.TempDir()` | โฟลเดอร์ชั่วคราว ลบให้อัตโนมัติ |
| `t.Parallel()` | รัน test พร้อมกัน |
| `ExampleXxx` | ตัวอย่างใน doc ที่ถูกรันและตรวจ output จริง |
| `BenchmarkXxx` | วัดความเร็ว |
| `FuzzXxx` | สุ่ม input หาเคสที่คุณคิดไม่ถึง |

## Test double: fake ดีกว่า mock
ถ้าโค้ดรับ interface (บท 06) คุณสามารถส่ง **fake** ที่เขียนเองง่าย ๆ เข้าไปได้
```go
type fakeClock struct{ now time.Time }
func (f fakeClock) Now() time.Time { return f.now }
```
ไม่ต้องใช้ library mock — fake ธรรมดาอ่านง่ายกว่าและพังยากกว่า

## ทดสอบอะไร?
- ✅ **พฤติกรรม** ที่เห็นจากภายนอก (input → output, error)
- ✅ **ขอบ** (edge cases): ว่าง, ศูนย์, ติดลบ, ค่าสูงสุด, Unicode, nil
- ✅ **กฎธุรกิจ** ทุกข้อจากโจทย์
- ❌ รายละเอียดภายใน (private field, ลำดับการเรียก) — ทำให้ refactor ยาก

## ดูตัวอย่างครบทุกเทคนิค
[`showcase/`](showcase) — โค้ดคิดราคาสินค้า พร้อม test ที่ใช้ทุกเทคนิคข้างบน

## แบบฝึกหัด
1. **เขียนโค้ดให้ผ่าน test:** `Slugify` ใน `exercise/exercise.go`
   (มี table test, example test และ fuzz test รอไว้แล้ว)
2. **เขียน test เอง:** เปิด `exercise/cart_test.go` แล้วเขียน table-driven test ให้ `Cart`
   — `Cart` มี **บั๊กซ่อนอยู่ 1 จุด** หาให้เจอด้วย test ของคุณ แล้วแก้ใน `cart.go`!
