# 04 — Pointer และ Value Semantics

## แก่นแท้
> **Go copy ค่าเสมอ (pass by value) — ไม่มีข้อยกเว้น**
> pointer ไม่ใช่เวทมนตร์ มันคือการ copy "ที่อยู่" แทนการ copy "ตัวข้อมูล"

ทุกครั้งที่ assign, ส่ง argument, return, หรือ range — Go copy ค่า
คำถามเดียวที่ต้องถามคือ: **"ค่าที่ถูก copy คืออะไร?"**

| สิ่งที่ส่ง | สิ่งที่ถูก copy | แก้แล้วข้างนอกเห็นไหม? |
|-----------|----------------|-----------------------|
| `int`, `string`, `struct` | ข้อมูลทั้งหมด | ❌ |
| `*T` (pointer) | ที่อยู่ (8 byte) | ✅ (ผ่าน `*p`) |
| slice | header (ptr, len, cap) | ✅ แก้ element ได้ / ❌ append ไม่เห็น |
| map, channel | ตัวชี้ภายใน | ✅ |

## Syntax
```go
x := 10
p := &x      // & = "ขอที่อยู่ของ"   → p มี type *int
*p = 20      // * = "ไปที่อยู่นั้น"   → x กลายเป็น 20
```

```
 x ┌────┐        p ┌──────────┐
   │ 20 │ ◀─────── │ 0xc000.. │
   └────┘          └──────────┘
```

## ตัวอย่างที่ต้องเข้าใจให้ได้
```go
func setValue(n int)   { n = 100 }   // แก้ copy — ข้างนอกไม่เปลี่ยน
func setPtr(n *int)    { *n = 100 }  // แก้ผ่านที่อยู่ — ข้างนอกเปลี่ยน

func appendOne(s []int) { s = append(s, 1) }   // header ที่ถูก copy เปลี่ยน — ข้างนอกไม่เห็น len ใหม่
func setFirst(s []int)  { s[0] = 1 }           // แก้ array ที่แชร์กัน — ข้างนอกเห็น
```

## เมื่อไหร่ใช้ pointer? (กฎที่ใช้จริง)
**ใช้ pointer เมื่อ:**
1. ต้องการ **แก้ไข** ค่าของผู้เรียก
2. struct **ใหญ่** และ copy บ่อย (หลายร้อย byte ขึ้นไป)
3. ค่านั้นมี **ตัวตน (identity)** — ต้องมีแค่ชิ้นเดียว เช่น connection, account, cache
4. struct มี `sync.Mutex` อยู่ข้างใน (ห้าม copy mutex!)
5. ต้องแยก "ไม่มีค่า" (`nil`) ออกจาก "zero value"

**ใช้ value เมื่อ:**
- ข้อมูลเล็กและไม่เปลี่ยน เช่น `Point{X, Y}`, `Money`, `time.Time`
- อยากให้ปลอดภัย: คนอื่นแก้ของเราไม่ได้

> **คิดแบบนี้:** value = "สำเนาเอกสาร" (ใครจะขีดเขียนก็ไม่กระทบต้นฉบับ)
> pointer = "ที่อยู่บ้าน" (ใครไปถึงก็เปลี่ยนเฟอร์นิเจอร์ได้)

## ไม่มี pointer arithmetic, มี garbage collector
- `p++` ทำไม่ได้ → ปลอดภัยกว่า C
- คืน pointer ของตัวแปร local ได้สบาย ๆ — compiler จะย้ายไป heap ให้เอง (escape analysis)
```go
func NewUser() *User {
    u := User{Name: "a"}
    return &u     // ✅ ถูกต้องใน Go (ใน C = บั๊ก)
}
```
ดูว่าอะไร escape ไป heap: `go build -gcflags=-m ./04-pointers`

## nil pointer
```go
var p *User
p.Name    // 💥 panic: nil pointer dereference
```
ฟังก์ชันที่คืน pointer ต้องบอกให้ชัดว่า "คืน nil ได้เมื่อไหร่" หรือคืน error คู่ไปด้วย

## แบบฝึกหัด
1. `Swap` — สลับค่าสองตัวแปรผ่าน pointer
2. `Increment` — เพิ่มค่าทุกตัวใน slice (ทำไมไม่ต้องใช้ pointer?)
3. `AppendTo` — append ผ่าน pointer ของ slice
4. `ResetScores` — แก้ค่าใน map ของ struct (กับดักคลาสสิก!)
5. `Ptr` — generic helper สำหรับสร้าง pointer จากค่า
