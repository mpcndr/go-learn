# 07 — Error: ค่าธรรมดาที่ต้องใส่ใจ

## แก่นแท้
> **Error คือค่า (value) ธรรมดา ไม่ใช่ exception — ไม่มีการกระโดดข้ามฟังก์ชันแบบมองไม่เห็น**
> ทุกจุดที่พังได้ จะเห็นชัดในโค้ด และต้องตัดสินใจตรงนั้นว่าจะทำอะไร

```go
type error interface {
    Error() string
}
```
แค่นั้นเอง — error คือ interface ที่มี method เดียว

## Pattern พื้นฐาน
```go
f, err := os.Open(name)
if err != nil {
    return fmt.Errorf("load config: %w", err)
}
defer f.Close()
```
`if err != nil` ดูซ้ำ ๆ แต่มันคือ **ความตั้งใจ**: ทุกทางที่พังได้ถูกมองเห็นและจัดการ

## เมื่อได้ error มา ทำได้ 3 อย่าง
1. **ส่งต่อพร้อมบริบท** (บ่อยที่สุด) — `fmt.Errorf("กำลังทำอะไร: %w", err)`
2. **จัดการเอง** — retry, ใช้ค่า default, แสดงข้อความให้ผู้ใช้
3. **log แล้วจบ** — ที่ขอบนอกสุดเท่านั้น (main, HTTP handler)

❌ อย่า log แล้ว return ด้วย — error เดียวจะถูก log ซ้ำหลายรอบ
❌ อย่าทิ้ง error เงียบ ๆ (`_ = doSomething()`) ยกเว้นตั้งใจจริงและมีคอมเมนต์

## การเติมบริบท (wrapping) — ทำให้ error เล่าเรื่องได้
```
checkout order 42: charge card: call payment API: connection refused
   ↑ service         ↑ payment     ↑ http client    ↑ ต้นเหตุ
```
แต่ละชั้นเติม "กำลังทำอะไร" — อ่านบรรทัดเดียวรู้เส้นทางทั้งหมด ไม่ต้องมี stack trace

- `%w` = wrap (ห่อ) → ผู้เรียกยังตรวจต้นเหตุได้ด้วย `errors.Is` / `errors.As`
- `%v` = แค่ข้อความ → ตัดความสัมพันธ์ (ใช้เมื่อไม่อยากให้ผู้เรียกผูกกับ error ภายใน)

## 3 แบบของ error

### 1) Sentinel error — ค่าคงที่ให้เปรียบเทียบ
```go
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) { ... }   // ✅ ทะลุผ่าน wrap ได้
if err == ErrNotFound { ... }            // ❌ ใช้ไม่ได้ถ้าถูก wrap
```

### 2) Custom error type — เมื่อต้องการข้อมูลเพิ่ม
```go
type ValidationError struct {
    Field  string
    Reason string
}
func (e *ValidationError) Error() string { return e.Field + ": " + e.Reason }

var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Println("field ที่ผิด:", ve.Field)
}
```

### 3) Opaque error — ไม่บอกอะไรเลย แค่ "พัง"
`errors.New("...")` หรือ `fmt.Errorf(...)` ธรรมดา — ใช้เมื่อผู้เรียกไม่ต้องแยกประเภท

> **หลักการเลือก:** ถามว่า *"ผู้เรียกจำเป็นต้องทำอะไรต่างกันตามชนิดของ error ไหม?"*
> ไม่ → opaque, ต้องแยกประเภท → sentinel, ต้องการรายละเอียด → custom type

## รวมหลาย error
```go
err := errors.Join(err1, err2)   // errors.Is ใช้กับทุกตัวข้างในได้
```

## panic ใช้เมื่อไหร่?
**เกือบไม่เคย** — panic = "โปรแกรมมีบั๊ก สถานะเสียหาย ไปต่อไม่ได้"
- ✅ index เกิน, nil pointer, `MustCompile` ตอนเริ่มโปรแกรมที่ input เป็นค่าคงที่
- ❌ ไฟล์ไม่เจอ, input ผู้ใช้ผิด, network ล่ม → เป็นเรื่องปกติ ต้องคืน error

`recover()` ใช้ที่ขอบนอก (เช่น HTTP server) เพื่อกันทั้งระบบล่มจาก request เดียว

## ออกแบบด้วย error
ตอนเขียน signature ให้ถามเสมอ: **"ฟังก์ชันนี้ล้มเหลวได้แบบไหนบ้าง?"**
แล้วตั้งชื่อความล้มเหลวเหล่านั้นให้เป็น sentinel/type — มันคือ **ส่วนหนึ่งของ API** ไม่ใช่เรื่องรอง

## แบบฝึกหัด
1. `ParseAge` — validate input คืน `*ValidationError`
2. `FindUser` — คืน sentinel `ErrNotFound`
3. `LoadProfile` — wrap error พร้อมบริบท, `errors.Is` ต้องยังทำงาน
4. `ValidateAll` — รวมหลาย error ด้วย `errors.Join`
5. `SafeDivide` — แปลง panic เป็น error ด้วย `recover`
