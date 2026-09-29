# 05 — Struct และ Method: ข้อมูลมาก่อน พฤติกรรมตามมา

## แก่นแท้
> **Go ไม่มี class ไม่มี inheritance — มีแค่ "ข้อมูล" (struct) กับ "ฟังก์ชันที่ผูกกับข้อมูล" (method)**
> และใช้ **composition** (ประกอบ) แทน inheritance (สืบทอด)

## Struct
```go
type Account struct {
    ID      string
    Owner   string
    balance int      // ตัวเล็ก = private (มองเห็นเฉพาะใน package เดียวกัน)
}

a := Account{ID: "A1", Owner: "Ann"}   // ระบุชื่อ field เสมอ (อ่านง่าย, เพิ่ม field ทีหลังไม่พัง)
```
**ตัวพิมพ์ใหญ่/เล็กคือระบบ public/private ของ Go** — ใช้กับทุกอย่าง: type, func, field, method, const

## Method = ฟังก์ชันที่มี receiver
```go
func (a Account)  Balance() int     { return a.balance }   // value receiver: อ่านอย่างเดียว
func (a *Account) Deposit(n int)    { a.balance += n }     // pointer receiver: แก้ไขได้
```
จริง ๆ แล้ว method ก็คือฟังก์ชันธรรมดา: `a.Deposit(5)` ≈ `(*Account).Deposit(&a, 5)`

### เลือก value หรือ pointer receiver?
- ต้องแก้ข้อมูล → **pointer**
- struct มี mutex หรือใหญ่ → **pointer**
- เล็ก, immutable (เช่น `Point`, `Money`) → **value**
- **ความสม่ำเสมอ:** ถ้ามี method ไหนเป็น pointer แล้ว ให้ทุก method ของ type นั้นเป็น pointer

## Constructor — เป็นแค่ฟังก์ชันธรรมดาที่ชื่อขึ้นต้นด้วย New
```go
func NewAccount(id, owner string) (*Account, error) {
    if id == "" { return nil, errors.New("id required") }
    return &Account{ID: id, Owner: owner}, nil
}
```
ใช้เมื่อต้อง **validate** หรือ **ตั้งค่าเริ่มต้นที่ไม่ใช่ zero value** — ถ้า zero value ใช้ได้อยู่แล้ว ไม่ต้องมี constructor

## Encapsulation — ปกป้อง invariant
**Invariant** = กฎที่ต้องจริงเสมอ เช่น "ยอดเงินห้ามติดลบ"
ทำ field เป็น private แล้วให้แก้ผ่าน method เท่านั้น → กฎถูกบังคับที่จุดเดียว

```go
func (a *Account) Withdraw(n int) error {
    if n <= 0          { return ErrInvalidAmount }
    if n > a.balance   { return ErrInsufficientFunds }
    a.balance -= n
    return nil
}
```
> **คิดแบบนี้:** ถ้าคนนอกทำให้ struct ของคุณอยู่ในสถานะผิดได้ = ออกแบบรั่ว

## Embedding — composition ที่สะดวก
```go
type Timestamps struct { CreatedAt, UpdatedAt time.Time }
func (t *Timestamps) Touch() { t.UpdatedAt = time.Now() }

type Post struct {
    Timestamps           // embed: ไม่มีชื่อ field
    Title string
}
p.Touch()                // method ถูก "ยก" (promote) ขึ้นมา
p.CreatedAt              // field ก็เช่นกัน
```
⚠️ นี่ **ไม่ใช่** inheritance: `Post` ไม่ใช่ `Timestamps` ส่ง `Post` ไปฟังก์ชันที่รับ `Timestamps` ไม่ได้
มันแค่ "มี" Timestamps อยู่ข้างใน และ Go ช่วยให้เรียกสั้นลง

## Method บน type อะไรก็ได้ที่คุณประกาศ
```go
type Celsius float64
func (c Celsius) ToFahrenheit() float64 { return float64(c)*9/5 + 32 }
```
สร้าง type ใหม่จาก type พื้นฐานเพื่อให้ **compiler ช่วยแยกความหมาย** — `Celsius` กับ `Fahrenheit` จะบวกกันผิด ๆ ไม่ได้

## แบบฝึกหัด
ใน `exercise/types.go` มี type ให้แล้ว ให้เขียน method ใน `exercise.go`
1. `Rect.Area`, `Rect.Scale` — value vs pointer receiver
2. `NewAccount`, `Deposit`, `Withdraw`, `Balance` — encapsulation + invariant
3. `Money.String` — แสดงผลเป็น "฿12.50"
4. `Stack` — ใช้ slice เป็น stack (Push, Pop, Len)
