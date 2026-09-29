# 02 — ฟังก์ชัน: หน่วยพื้นฐานของการคิด

## แก่นแท้
> **ฟังก์ชันคือ "สัญญา": รับอะไร → คืนอะไร → พังได้อย่างไร**
> และใน Go ฟังก์ชันเป็น *ค่า* ตัวหนึ่ง — เก็บในตัวแปร ส่งเป็น argument คืนจากฟังก์ชันได้

การออกแบบที่ดีเริ่มจากการเขียน **signature** ก่อนเขียนเนื้อใน:
```go
func Divide(a, b float64) (float64, error)
```
อ่านบรรทัดเดียวรู้เลยว่า "หารได้ แต่อาจพัง" — นี่คือเอกสารที่ compiler ตรวจให้

## Multiple return values
```go
func MinMax(nums []int) (min, max int) { ... }   // named results
q, r := divmod(7, 2)
_, r = divmod(7, 2)   // ไม่ใช้ค่าไหนก็ทิ้งด้วย _
```
Pattern ที่เห็นทุกที่: `(result, error)` และ `(value, ok bool)`

## Variadic
```go
func Sum(nums ...int) int { ... }   // nums คือ []int
Sum(1, 2, 3)
Sum(list...)                        // กระจาย slice
```

## ฟังก์ชันเป็นค่า + Closure
```go
func MakeCounter() func() int {
    count := 0                 // ตัวแปรนี้ "ถูกจับ" (captured) ไว้
    return func() int {
        count++                // closure แก้ตัวแปรข้างนอกได้
        return count
    }
}
c := MakeCounter()
c(); c(); c()  // 1, 2, 3 — count ยังอยู่ แม้ MakeCounter จบไปแล้ว
```
Closure = ฟังก์ชัน + ตัวแปรที่มันจำได้ ใช้ทำ state เล็ก ๆ, callback, middleware

## defer — "ทำตอนจบฟังก์ชัน"
```go
f, err := os.Open(path)
if err != nil { return err }
defer f.Close()     // วางไว้ติดกับตอนเปิด — ไม่มีวันลืมปิด
```
- ทำงานตอนฟังก์ชัน return (ทั้งปกติและ panic)
- หลาย defer ทำงานแบบ **LIFO** (เข้าหลังออกก่อน)
- argument ของ defer ถูกประเมิน **ตอนเจอ defer** ไม่ใช่ตอนรัน

## หลักการออกแบบฟังก์ชัน
1. **ทำเรื่องเดียว** — ถ้าต้องใช้คำว่า "และ" อธิบาย ให้แยก
2. **ชื่อบอกว่าทำอะไร** ไม่ใช่ทำอย่างไร (`SortByPrice` ไม่ใช่ `BubbleLoop`)
3. **Pure function ก่อน** — รับ input คืน output ไม่แตะของข้างนอก → ทดสอบง่ายที่สุด
4. **Early return** — เช็ค error/กรณีพิเศษแล้ว return ออกก่อน ให้ happy path ชิดซ้าย

```go
// ❌ ซ้อนลึก
if ok { if valid { if ready { doWork() } } }

// ✅ early return
if !ok    { return errNotOK }
if !valid { return errInvalid }
if !ready { return errNotReady }
doWork()
```

## กับดักที่มักพลาด
- Closure ใน loop: ตั้งแต่ Go 1.22 ตัวแปร loop ถูกสร้างใหม่ทุกรอบแล้ว (ปลอดภัย) แต่โค้ดเก่าจะเห็น `v := v` ซึ่งไม่จำเป็นอีกต่อไป
- `defer` ใน loop ยาว ๆ จะสะสมจนจบฟังก์ชัน — ถ้าเปิดไฟล์ใน loop ให้แยกเป็นฟังก์ชันย่อย

## แบบฝึกหัด
1. `MinMax` — multiple return
2. `MakeCounter` — closure
3. `Apply` — ส่งฟังก์ชันเป็น argument
4. `Compose` — คืนฟังก์ชันที่รวมสองฟังก์ชัน
5. `Memoize` — closure + map เพื่อ cache ผลลัพธ์
