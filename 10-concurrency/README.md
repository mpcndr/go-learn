# 10 — Concurrency: Goroutine, Channel, Sync

## แก่นแท้
> **"Don't communicate by sharing memory; share memory by communicating."**
> แทนที่หลาย goroutine จะแย่งกันแก้ตัวแปรเดียว (แล้วต้อง lock) → ให้ส่ง "ความเป็นเจ้าของ" ข้อมูลผ่าน channel

และอีกข้อที่สำคัญพอ ๆ กัน:
> **Concurrency ≠ Parallelism** — concurrency คือ *การออกแบบ* ให้งานหลายอย่างเดินไปพร้อมกันได้
> parallelism คือ *การรัน* พร้อมกันจริงบนหลาย CPU (อาจเกิดหรือไม่เกิดก็ได้)

## Goroutine — ถูกมาก สร้างได้เป็นแสน
```go
go doWork()     // รันแยกไปเลย ไม่รอ
```
- ใช้ memory เริ่มต้นแค่ ~2KB (thread ของ OS ใช้ ~1MB)
- Go runtime จัดตาราง goroutine ลงบน OS thread ให้เอง (M:N scheduler)
- ⚠️ `main` จบ = โปรแกรมจบ ทุก goroutine ตายตาม → ต้อง **รอ** ให้เป็น

## sync.WaitGroup — รอให้ทุกตัวเสร็จ
```go
var wg sync.WaitGroup
for _, url := range urls {
    wg.Add(1)
    go func() {
        defer wg.Done()
        fetch(url)
    }()
}
wg.Wait()
```
(Go 1.25+ มี `wg.Go(func(){...})` ที่รวม Add/Done ให้)

## Channel — ท่อส่งข้อมูลระหว่าง goroutine
```go
ch := make(chan int)      // unbuffered: ส่งกับรับต้อง "เจอกัน" (synchronize)
ch := make(chan int, 10)  // buffered: ส่งได้ 10 ตัวโดยไม่ต้องรอคนรับ

ch <- 42        // ส่ง (block ถ้าไม่มีคนรับ / buffer เต็ม)
v := <-ch       // รับ (block ถ้าไม่มีของ)
v, ok := <-ch   // ok=false เมื่อ channel ถูกปิดและว่างแล้ว
close(ch)       // บอกว่า "ไม่มีของส่งมาอีก" — ผู้ส่งเป็นคนปิดเท่านั้น!
for v := range ch { }   // รับจนกว่าจะถูกปิด
```

**ทิศทาง** ใน signature ช่วยให้ compiler ตรวจ:
```go
func producer(out chan<- int)  // ส่งได้อย่างเดียว
func consumer(in <-chan int)   // รับได้อย่างเดียว
```

## select — รอหลาย channel พร้อมกัน
```go
select {
case msg := <-inbox:
    handle(msg)
case <-time.After(time.Second):
    fmt.Println("timeout")
case <-done:
    return
}
```

## sync.Mutex — เมื่อ "แชร์ state" ง่ายกว่าจริง ๆ
```go
type Counter struct {
    mu sync.Mutex
    n  map[string]int
}
func (c *Counter) Inc(k string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.n[k]++
}
```
**เลือกอะไรดี?**
- ส่งข้อมูล/งาน/ผลลัพธ์ต่อกันเป็นทอด ๆ → **channel**
- ปกป้อง state ที่หลายคนอ่าน/เขียน (cache, counter) → **mutex**

## Pattern ที่ต้องรู้

### Worker pool — จำกัดจำนวนงานที่ทำพร้อมกัน
```
jobs ──▶ [worker 1] ──┐
     ──▶ [worker 2] ──┼──▶ results
     ──▶ [worker 3] ──┘
```

### Pipeline — ต่องานเป็นขั้น ๆ
```
generate ──▶ square ──▶ filter ──▶ sum
```
แต่ละขั้นคือ goroutine ที่รับจาก channel หนึ่งแล้วส่งต่ออีก channel

### Fan-out / Fan-in
แจกงานให้หลายตัวทำ (fan-out) แล้วรวมผลเข้า channel เดียว (fan-in)

## กฎเหล็ก 5 ข้อ
1. **ทุก goroutine ต้องรู้ว่าจะจบเมื่อไหร่** — ไม่งั้นเกิด goroutine leak
2. **ผู้ส่งเป็นคนปิด channel** — ปิดซ้ำหรือส่งเข้า channel ที่ปิดแล้ว = panic
3. **map ไม่ปลอดภัยต่อการเขียนพร้อมกัน** — ต้องใช้ mutex
4. **รัน `go test -race` เสมอ** — มันจับ data race ที่ตาไม่เห็น
5. **เขียนแบบ sequential ให้ถูกก่อน** แล้วค่อยทำ concurrent เฉพาะจุดที่ต้องการ

## แบบฝึกหัด (รันด้วย `go test -race ./10-concurrency/...`)
1. `ParallelSum` — แบ่งงานให้หลาย goroutine แล้วรวมผล
2. `SafeCounter` — mutex
3. `WorkerPool` — จำกัดจำนวน worker, รักษาลำดับผลลัพธ์
4. `Pipeline` — generate → square → sum ด้วย channel
5. `FirstResponse` — select: เอาคำตอบแรกที่มาถึง
