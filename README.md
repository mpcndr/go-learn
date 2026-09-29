# go-learn — เรียน Go ตั้งแต่ศูนย์ จนคิดเป็นระบบ

เป้าหมายของ repo นี้ไม่ใช่แค่ "เขียน Go ได้" แต่คือ **คิดเป็นระบบ แล้วเขียนออกมาให้ตรงกับที่ออกแบบไว้ในหัว**

Go เป็นภาษาที่ถูกออกแบบมาให้ "เรียบง่ายแบบตั้งใจ" (simple on purpose) — ฟีเจอร์น้อย แต่ทุกฟีเจอร์ต่อกันได้ลงตัว
ถ้าเข้าใจ *เหตุผล* ที่ภาษาเป็นแบบนี้ คุณจะเขียนโค้ดที่อ่านง่าย ทดสอบง่าย และแก้ง่าย โดยไม่ต้องจำ pattern เยอะ

---

## วิธีใช้ repo นี้

แต่ละบทมีโครงสร้างเหมือนกัน:

```
NN-topic/
├── README.md              ← อ่านก่อน: แนวคิด + แก่นแท้ + กับดักที่มักพลาด
├── main.go                ← ตัวอย่างที่รันได้:  go run ./NN-topic
└── exercise/
    ├── exercise.go        ← โจทย์ (มี TODO ให้เติม)
    ├── solution.go        ← เฉลย (ซ่อนด้วย build tag — อย่าเพิ่งเปิดดู!)
    └── exercise_test.go   ← test ที่จะบอกว่าคุณทำถูกหรือยัง
```

ขั้นตอนเรียนแต่ละบท:

```bash
# 1. อ่าน README ของบท
# 2. รันตัวอย่าง แล้วลองแก้เล่นดู
go run ./01-basics

# 3. ทำแบบฝึกหัด — แก้ exercise.go จน test ผ่าน
go test ./01-basics/...

# 4. ติดจริง ๆ ค่อยดูเฉลย (หรือรัน test กับเฉลย)
go test -tags solution ./01-basics/...
```

> **กติกาทอง:** พิมพ์โค้ดเองทุกบรรทัด อย่า copy-paste — มือต้องจำ syntax ไปพร้อมกับสมอง

---

## Roadmap

### ภาค 1 — พื้นฐาน (เข้าใจภาษา)
| บท | หัวข้อ | แก่นที่ต้องได้ |
|----|--------|----------------|
| [01](01-basics) | ตัวแปร, type, zero value, control flow | ทุกอย่างใน Go มี type และมีค่าเริ่มต้นเสมอ |
| [02](02-functions) | ฟังก์ชัน, multiple return, closure, defer | ฟังก์ชันคือค่า (value) ตัวหนึ่ง |
| [03](03-collections) | array, slice, map, string/rune | slice คือ "หน้าต่าง" ที่มองไปยัง array |
| [04](04-pointers) | pointer และ value semantics | Go **copy เสมอ** — pointer คือการ copy "ที่อยู่" |
| [05](05-structs-methods) | struct, method, embedding | ข้อมูลมาก่อน พฤติกรรมตามมา |

### ภาค 2 — แก่นของ Go (คิดแบบ Go)
| บท | หัวข้อ | แก่นที่ต้องได้ |
|----|--------|----------------|
| [06](06-interfaces) | interface | กำหนดจาก "สิ่งที่ต้องการ" ไม่ใช่ "สิ่งที่เป็น" |
| [07](07-errors) | error handling | error คือค่าธรรมดา ต้องจัดการตรงนั้นเลย |
| [08](08-generics) | generics | ใช้เมื่อ logic เหมือนกันทุก type เท่านั้น |
| [09](09-testing) | testing, table-driven test | test คือ spec ที่รันได้ |

### ภาค 3 — Concurrency (จุดแข็งของ Go)
| บท | หัวข้อ | แก่นที่ต้องได้ |
|----|--------|----------------|
| [10](10-concurrency) | goroutine, channel, select, sync | "Don't communicate by sharing memory; share memory by communicating." |
| [11](11-context) | context, timeout, cancel | ทุกงานที่ใช้เวลาต้องหยุดได้ |

### ภาค 4 — ออกแบบระบบ (จากหัว → โค้ด)
| บท | หัวข้อ | แก่นที่ต้องได้ |
|----|--------|----------------|
| [12](12-packages-design) | package, module, การแบ่งชั้น | dependency ไหลทางเดียว |
| [13](13-capstone-library) | **โปรเจกต์: ระบบห้องสมุด** | ออกแบบบนกระดาษ → แปลงเป็นโค้ดทีละชั้น |

และอ่านควบคู่ตลอดทาง: **[docs/thinking-in-go.md](docs/thinking-in-go.md)** — วิธีคิดเป็นระบบก่อนลงมือเขียน

---

## Checklist ความคืบหน้า

- [ ] 01 basics
- [ ] 02 functions
- [ ] 03 collections
- [ ] 04 pointers
- [ ] 05 structs & methods
- [ ] 06 interfaces
- [ ] 07 errors
- [ ] 08 generics
- [ ] 09 testing
- [ ] 10 concurrency
- [ ] 11 context
- [ ] 12 packages & design
- [ ] 13 capstone — ออกแบบเองก่อน แล้วค่อยเทียบกับตัวอย่าง

---

## คำสั่งที่ใช้บ่อย

```bash
go run ./01-basics              # รันโปรแกรม
go test ./...                   # รัน test ทั้งหมด
go test -v -run TestFizzBuzz ./01-basics/...   # รันเฉพาะ test ที่ต้องการ
go test -race ./10-concurrency/...             # ตรวจ data race
go vet ./...                    # ตรวจโค้ดที่น่าจะผิด
gofmt -w .                      # จัด format (Go มี format เดียว — ไม่ต้องเถียงกัน)
go doc strings.Builder          # อ่าน doc ของ standard library
```

## แหล่งอ่านเพิ่ม (ภาษาอังกฤษ)
- [A Tour of Go](https://go.dev/tour) — ทำควบคู่กับภาค 1
- [Effective Go](https://go.dev/doc/effective_go) — อ่านหลังจบภาค 2
- [Go Proverbs](https://go-proverbs.github.io/) — ปรัชญาของภาษาในประโยคสั้น ๆ
- [Go by Example](https://gobyexample.com/) — ตัวอย่างสั้น ๆ ไว้เปิดดูตอนลืม
