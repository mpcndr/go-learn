# 03 — Collections: Array, Slice, Map, String

## แก่นแท้
> **Slice ไม่ได้เก็บข้อมูล — slice คือ "หน้าต่าง" ที่มองไปยัง array ข้างหลัง**

ถ้าเข้าใจข้อนี้ข้อเดียว บั๊กเรื่อง slice จะหายไป 90%

## Array — ขนาดตายตัว เป็นส่วนหนึ่งของ type
```go
var a [3]int          // [0 0 0]
b := [...]int{1,2,3}  // ให้ compiler นับขนาด
```
`[3]int` กับ `[4]int` เป็นคนละ type กัน ส่ง array ไปฟังก์ชัน = **copy ทั้งก้อน**
ในทางปฏิบัติแทบไม่ใช้ array ตรง ๆ — ใช้ slice แทน

## Slice — โครงสร้างจริงมีแค่ 3 ช่อง
```
slice header:  ┌─────┬─────┬─────┐
               │ ptr │ len │ cap │
               └──┬──┴─────┴─────┘
                  ↓
array ข้างหลัง: [ 10 | 20 | 30 | 40 | 50 ]
```
- `len` = จำนวนที่มองเห็น, `cap` = พื้นที่ที่ขยายได้โดยไม่ต้องย้าย array
- `s[1:3]` สร้าง header ใหม่ที่ **ชี้ array เดิม** → แก้ตัวหนึ่ง อีกตัวเปลี่ยนตาม!

```go
nums := []int{10, 20, 30, 40, 50}
part := nums[1:3]    // [20 30]
part[0] = 99
fmt.Println(nums)    // [10 99 30 40 50]  ← array เดียวกัน!
```

### append — ต้องรับค่ากลับเสมอ
```go
s = append(s, x)     // ✅
append(s, x)         // ❌ compile ไม่ผ่าน (ผลลัพธ์ไม่ได้ใช้)
```
ถ้า `cap` พอ → เขียนลง array เดิม, ถ้าไม่พอ → จอง array ใหม่ (ใหญ่ขึ้น ~2 เท่า) แล้ว copy
นี่คือเหตุผลที่ต้องรับค่าคืน: header อาจชี้ไป array ใหม่แล้ว

### จองพื้นที่ล่วงหน้าเมื่อรู้ขนาด
```go
out := make([]int, 0, len(input))   // len 0, cap เท่า input — ไม่ต้องขยายระหว่างทาง
```

### copy เพื่อตัดความสัมพันธ์
```go
clone := make([]int, len(src))
copy(clone, src)          // หรือ slices.Clone(src)
```

## Map — hash table
```go
ages := map[string]int{"alice": 30}
ages["bob"] = 25
age, ok := ages["carol"]   // ok=false, age=0 (zero value) — "comma ok" idiom
delete(ages, "bob")
```
- ลำดับตอน `range` **สุ่มโดยตั้งใจ** — อย่าพึ่งลำดับ (ถ้าต้องการให้ sort key ก่อน)
- `var m map[string]int` คือ nil map → อ่านได้ แต่ **เขียนแล้ว panic** → ใช้ `make` หรือ literal
- map เป็น reference-like — ส่งเข้าฟังก์ชันแล้วแก้ ข้างนอกเห็นด้วย
- map **ไม่ปลอดภัย** เมื่อหลาย goroutine เขียนพร้อมกัน (ดูบท 10)

## String, byte, rune
```go
s := "สวัสดี"
len(s)                     // 18  ← byte (UTF-8, ภาษาไทยตัวละ 3 byte)
utf8.RuneCountInString(s)  // 6   ← ตัวอักษร
for i, r := range s { }    // range วนทีละ rune (ตัวอักษร) ✅
s[0]                       // byte ตัวแรก — ไม่ใช่ "ส" ❌
[]rune(s)                  // แปลงเป็น slice ของตัวอักษร
```
- string เป็น **immutable** — ต่อ string ใน loop เยอะ ๆ ให้ใช้ `strings.Builder`
- `byte` = `uint8`, `rune` = `int32` (Unicode code point)

## Package ที่ควรรู้จัก
- `slices`: `Sort`, `Contains`, `Index`, `Reverse`, `Clone`, `Equal`
- `maps`: `Keys`, `Values`, `Clone`
- `strings`: `Fields`, `Split`, `Join`, `ToLower`, `Builder`

## แบบฝึกหัด
1. `WordCount` — นับคำ (map)
2. `Reverse` — กลับ string แบบรองรับภาษาไทย (rune)
3. `Unique` — ลบตัวซ้ำโดยรักษาลำดับ (map + slice)
4. `Chunk` — แบ่ง slice เป็นก้อน ๆ
5. `TopN` — หาคำที่เจอบ่อยที่สุด N อันดับ (map + sort)
