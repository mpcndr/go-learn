# 08 — Generics: เขียนครั้งเดียว ใช้ได้หลาย type

## แก่นแท้
> **ใช้ generics เมื่อ "logic เหมือนกันทุกบรรทัด ต่างกันแค่ type"**
> ถ้าพฤติกรรมต่างกันตาม type → ใช้ interface

| สถานการณ์ | ใช้อะไร |
|-----------|---------|
| `Max` ของ int, float, string — วิธีหาเหมือนกัน | **generics** |
| `Area()` ของวงกลม, สี่เหลี่ยม — วิธีคำนวณต่างกัน | **interface** |
| container: Stack, Queue, Set, Cache | **generics** |

## Syntax
```go
//        type parameter   constraint
//              ↓            ↓
func Map[T, U any](s []T, fn func(T) U) []U {
    out := make([]U, 0, len(s))
    for _, v := range s {
        out = append(out, fn(v))
    }
    return out
}

lengths := Map([]string{"go", "rust"}, func(s string) int { return len(s) })
// compiler เดา T=string, U=int ให้เอง
```

## Constraint — "type นี้ต้องทำอะไรได้บ้าง"
constraint ก็คือ interface นั่นเอง
```go
any                    // อะไรก็ได้ (แต่ทำได้แค่ assign/ส่งต่อ)
comparable             // ใช้ == และ != ได้ (จำเป็นสำหรับ key ของ map)
cmp.Ordered            // ใช้ < > <= >= ได้ (ตัวเลข, string)

// constraint ที่กำหนดเอง
type Number interface {
    ~int | ~int64 | ~float64    // ~ = รวม type ที่สร้างจาก type นี้ด้วย เช่น type Celsius float64
}
```

## Generic type
```go
type Stack[T any] struct {
    items []T
}
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

var s Stack[string]
```
⚠️ method **ไม่สามารถ** มี type parameter เพิ่มเองได้ — ต้องใช้ของ type เท่านั้น
(ถ้าต้องการ ให้เขียนเป็นฟังก์ชันธรรมดาแทน)

## Standard library ที่เป็น generic แล้ว — ใช้ก่อนเขียนเอง!
- `slices`: `Contains`, `Index`, `Sort`, `SortFunc`, `Max`, `Min`, `Reverse`, `Compact`
- `maps`: `Keys`, `Values`, `Clone`, `Equal`
- `cmp`: `Compare`, `Or`
- `sync.OnceValue`, `atomic.Pointer[T]`

## ข้อควรระวัง
- **อย่าเริ่มจาก generics** — เขียนแบบ type เดียวให้ใช้ได้ก่อน เมื่อต้องการ type ที่สองค่อยแปลง
- Generic ที่ซับซ้อนเกินไปทำให้อ่านยาก — ถ้าต้องอธิบาย constraint ยาวเป็นบรรทัด ลองคิดใหม่
- ถ้าฟังก์ชันรับ interface แล้วเรียกแค่ method ของมัน → ไม่ต้องใช้ generics

## แบบฝึกหัด
1. `Map`, `Filter`, `Reduce` — ฟังก์ชันพื้นฐานของ functional programming
2. `Sum` — ใช้ custom constraint `Number`
3. `GroupBy` — จัดกลุ่มด้วย key (ต้องการ `comparable`)
4. `Set[T]` — generic type พร้อม method
5. `Cache[K, V]` — generic type ที่มีสอง type parameter
