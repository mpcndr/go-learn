# 01 — พื้นฐาน: ตัวแปร, Type, Control Flow

## แก่นแท้
> **ทุกค่าใน Go มี type ที่รู้ตั้งแต่ compile และมี "zero value" เสมอ**

ไม่มี `undefined` ไม่มีตัวแปรที่ "ยังไม่มีค่า" — ประกาศปุ๊บมีค่าปั๊บ:

| type | zero value |
|------|-----------|
| `int`, `float64` | `0` |
| `string` | `""` |
| `bool` | `false` |
| pointer, slice, map, channel, func, interface | `nil` |
| struct | ทุก field เป็น zero value ของตัวเอง |

นี่คือเหตุผลที่โค้ด Go ไม่ค่อยต้องเช็ค "มีค่าหรือยัง" — ออกแบบให้ zero value ใช้งานได้เลย

## ประกาศตัวแปร

```go
var age int          // zero value = 0
var name = "Go"      // ให้ compiler เดา type (string)
count := 10          // แบบสั้น — ใช้ได้เฉพาะในฟังก์ชัน (ใช้บ่อยที่สุด)
const Pi = 3.14159   // ค่าคงที่ — คำนวณตอน compile
```

**กฎที่ควรจำ:** ประกาศตัวแปรแล้วไม่ใช้ = compile ไม่ผ่าน, import แล้วไม่ใช้ = compile ไม่ผ่าน
Go บังคับความสะอาดตั้งแต่ต้น

## Type ต้องแปลงเองเสมอ
```go
var i int = 10
var f float64 = float64(i)  // ต้องแปลงเอง — Go ไม่แปลงให้อัตโนมัติ
```
ดูเหมือนน่ารำคาญ แต่ป้องกันบั๊กแบบ "นึกว่าเป็น int แต่จริง ๆ เป็น float" ได้หมด

## Control flow — มีแค่ 3 แบบ

```go
// if — ประกาศตัวแปรในเงื่อนไขได้ (scope อยู่แค่ใน if)
if n := len(s); n > 10 { ... }

// for — loop เดียวของ Go (ไม่มี while)
for i := 0; i < 10; i++ { }   // แบบ C
for cond { }                   // แบบ while
for { }                        // วนไม่รู้จบ
for i, v := range items { }    // วนใน collection

// switch — ไม่ต้อง break (ไม่ตกไป case ถัดไป)
switch {
case score >= 80: grade = "A"
case score >= 70: grade = "B"
default:          grade = "F"
}
```

## `iota` — สร้างค่าคงที่แบบลำดับ (ใช้แทน enum)
```go
type Weekday int
const (
    Sunday Weekday = iota  // 0
    Monday                 // 1
    Tuesday                // 2
)
```

## กับดักที่มักพลาด
- `:=` ใน scope ใหม่จะสร้างตัวแปร **ใหม่** ที่บังตัวเก่า (shadowing) — ระวังใน `if`/`for`
- หาร int ด้วย int ได้ int (`7/2 == 3`) ถ้าอยากได้ทศนิยมต้องแปลงเป็น float ก่อน
- `len("สวัสดี")` ได้ 18 ไม่ใช่ 6 — เพราะนับเป็น byte (ดูบท 03)

## แบบฝึกหัด
เปิด `exercise/exercise.go` แล้วทำให้ `go test ./01-basics/...` ผ่าน
1. `FizzBuzz(n)` — คืน slice ของ string ตั้งแต่ 1 ถึง n
2. `Grade(score)` — แปลงคะแนนเป็นเกรด
3. `SumDigits(n)` — ผลรวมของเลขแต่ละหลัก
4. `IsPrime(n)` — ตรวจจำนวนเฉพาะ
