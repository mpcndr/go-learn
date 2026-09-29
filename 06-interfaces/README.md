# 06 — Interface: หัวใจของการออกแบบใน Go

## แก่นแท้
> **Interface บอกว่า "ฉันต้องการอะไร" ไม่ใช่ "คุณเป็นอะไร"**
> และ Go implement interface **โดยอัตโนมัติ** (implicit) — ไม่มีคำว่า `implements`

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```
type ไหนมี method `Write` ตรง signature = เป็น `Writer` ทันที ไม่ต้องประกาศ
ผลลัพธ์: **ผู้ใช้เป็นคนกำหนด interface** ไม่ใช่ผู้สร้าง → decoupling ของจริง

## ทำไมเรื่องนี้ถึงสำคัญมาก
ในภาษาอื่น: สร้าง class → ประกาศว่า implements อะไร → ผูกกันแน่น
ใน Go:
1. เขียน struct ตรง ๆ ก่อน (ไม่ต้องคิดถึง interface เลย)
2. ตอนที่ **ผู้ใช้** ต้องการความยืดหยุ่น (ทดสอบ, สลับ implementation) → ผู้ใช้ประกาศ interface เล็ก ๆ ของตัวเอง
3. struct เดิมเข้ากับ interface นั้นได้ทันทีโดยไม่ต้องแก้

```go
// package notification (ผู้ใช้) — ต้องการแค่ "ส่งได้"
type Sender interface {
    Send(to, msg string) error
}

func NotifyAll(s Sender, users []string) error { ... }
```
ตอน production ส่ง `*SMTPClient`, ตอน test ส่ง `fakeSender` — `NotifyAll` ไม่ต้องรู้เลย

## หลักการออกแบบ interface
1. **ยิ่งเล็กยิ่งดี** — 1–3 method (`io.Reader`, `io.Writer`, `fmt.Stringer`, `error` มี method เดียว)
2. **ประกาศฝั่งผู้ใช้** ไม่ใช่ฝั่งผู้สร้าง
3. **Accept interfaces, return structs** — รับ interface ให้ยืดหยุ่น, คืน type จริงให้ผู้ใช้ได้ครบทุกความสามารถ
4. **อย่าสร้าง interface ก่อนมีคนใช้** — ถ้ามี implementation เดียวและไม่มีใครต้องสลับ ไม่ต้องมี interface
5. ชื่อ interface method เดียว มักลงท้ายด้วย `-er`: `Reader`, `Sender`, `Validator`

## Interface ประกอบกันได้
```go
type ReadWriter interface {
    Reader
    Writer
}
```

## ข้างในของ interface value
```
interface value = ( type , pointer ไปยังค่า )
```
- `var s Shape` → `(nil, nil)` = nil interface
- ⚠️ **กับดักใหญ่:** `var p *Circle = nil; var s Shape = p` → `(*Circle, nil)` ซึ่ง **ไม่เท่ากับ nil**!
  เพราะฉะนั้น: ฟังก์ชันที่คืน `error` ให้ `return nil` ตรง ๆ อย่าคืนตัวแปร pointer ที่เป็น nil

## Type assertion และ type switch
```go
var s Shape = Circle{R: 2}
c, ok := s.(Circle)          // ถ้าไม่ใช่ → ok=false (ไม่ panic)

switch v := s.(type) {
case Circle: fmt.Println("วงกลม r =", v.R)
case Rect:   fmt.Println("สี่เหลี่ยม")
default:     fmt.Println("ไม่รู้จัก")
}
```
ใช้เท่าที่จำเป็น — ถ้าต้อง switch type บ่อย ๆ แปลว่าควรเพิ่ม method เข้าไปใน interface แทน

## `any` (= `interface{}`)
รับได้ทุกอย่าง แต่ **สูญเสีย type safety** — ใช้เมื่อจำเป็นจริง ๆ (เช่น JSON, fmt) ถ้าใช้บ่อย ลองดู generics (บท 08)

## Interface ที่ควรรู้จักใน standard library
| interface | method | ใช้ทำอะไร |
|-----------|--------|-----------|
| `fmt.Stringer` | `String() string` | กำหนดวิธีพิมพ์ |
| `error` | `Error() string` | error ทุกตัว |
| `io.Reader` | `Read([]byte) (int, error)` | อ่านข้อมูลจากอะไรก็ได้ |
| `io.Writer` | `Write([]byte) (int, error)` | เขียนข้อมูลไปที่ไหนก็ได้ |
| `sort.Interface` | `Len, Less, Swap` | เรียงลำดับ |
| `http.Handler` | `ServeHTTP(w, r)` | จัดการ HTTP request |

## แบบฝึกหัด
1. ทำให้ `Circle`, `Rectangle`, `Triangle` เป็น `Shape` แล้วเขียน `TotalArea`, `Largest`
2. `Describe` — type switch
3. `CountingWriter` — ห่อ `io.Writer` แล้วนับ byte (decorator pattern)
4. `NotifyAll` — ใช้ interface ที่ประกาศฝั่งผู้ใช้ + fake ใน test
