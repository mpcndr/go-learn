# 11 — Context: ทุกงานที่ใช้เวลา ต้องหยุดได้

## แก่นแท้
> **`context.Context` คือ "สัญญาณหยุด" ที่ส่งต่อลงไปทุกชั้นของงาน**
> ผู้ใช้ปิดเว็บ → HTTP handler รู้ → query DB ถูกยกเลิก → goroutine ทั้งหมดเลิกทำ
> ไม่มีงานไหนทำต่อไปเปล่า ๆ

ปัญหาที่ context แก้: ในบท 10 goroutine ไม่มีทางรู้เลยว่า "คนที่รอผลลัพธ์ เลิกรอไปแล้ว"

## สร้าง context
```go
ctx := context.Background()                              // รากของทุก context (ใน main, test)
ctx, cancel := context.WithCancel(parent)                 // ยกเลิกเองด้วย cancel()
ctx, cancel := context.WithTimeout(parent, 2*time.Second) // ยกเลิกเมื่อเกินเวลา
ctx, cancel := context.WithDeadline(parent, t)            // ยกเลิกเมื่อถึงเวลา t
defer cancel()   // ⚠️ เรียกเสมอ! ไม่งั้น resource รั่ว (go vet เตือนให้)
```

context เป็น **ต้นไม้**: ยกเลิกพ่อ → ลูกทุกตัวถูกยกเลิก (แต่ยกเลิกลูกไม่กระทบพ่อ)
```
Background
   └── WithTimeout(5s)          ← request ทั้งหมด
          ├── WithTimeout(1s)   ← เรียก service A
          └── WithTimeout(2s)   ← เรียก service B
```

## ฝั่งที่ "ถูกยกเลิก" ต้องคอยฟัง
```go
func worker(ctx context.Context, jobs <-chan Job) error {
    for {
        select {
        case <-ctx.Done():          // channel ที่ถูกปิดเมื่อยกเลิก
            return ctx.Err()        // context.Canceled หรือ context.DeadlineExceeded
        case j, ok := <-jobs:
            if !ok { return nil }
            process(j)
        }
    }
}
```
สำหรับงาน loop ยาว ๆ ที่ไม่มี channel: เช็ค `if err := ctx.Err(); err != nil { return err }` เป็นระยะ

## กฎการใช้ context
1. **เป็น parameter แรกเสมอ** ชื่อ `ctx`: `func Fetch(ctx context.Context, url string) error`
2. **อย่าเก็บ context ไว้ใน struct** — ส่งผ่าน parameter ทุกครั้ง
3. **อย่าส่ง `nil`** — ถ้าไม่รู้จะใช้อะไร ใช้ `context.TODO()`
4. **`WithValue` ใช้เฉพาะข้อมูลที่ผูกกับ request** (request ID, user ที่ login) — ห้ามใช้ส่ง dependency หรือ parameter ปกติ
5. **ผู้สร้าง context เป็นคนเรียก `cancel`**

## Pattern: errgroup (ทำเองแบบง่าย)
เรียกหลายงานพร้อมกัน ถ้าตัวไหนพัง → ยกเลิกตัวอื่นทั้งหมด แล้วคืน error ตัวแรก
(ในงานจริงใช้ `golang.org/x/sync/errgroup` — บทนี้ให้เขียนเองเพื่อเข้าใจกลไก)

## แบบฝึกหัด
1. `SlowSquare` — งานที่เคารพการยกเลิก
2. `WithRetry` — retry แต่หยุดทันทีเมื่อ context ถูกยกเลิก
3. `FetchAll` — เรียกพร้อมกันทั้งหมด ถ้าตัวใดพัง ยกเลิกที่เหลือ (errgroup แบบทำเอง)
4. `Generator` — goroutine ที่ผลิตค่าไปเรื่อย ๆ จนกว่าจะถูกยกเลิก (ไม่ leak)
