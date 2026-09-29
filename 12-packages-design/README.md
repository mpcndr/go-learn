# 12 — Package, Module และการแบ่งชั้น

## แก่นแท้
> **Package คือหน่วยของ "ความรับผิดชอบ" ไม่ใช่แค่โฟลเดอร์เก็บไฟล์**
> และ **dependency ต้องไหลทางเดียว** — ชั้นในไม่รู้จักชั้นนอกเลย

## พื้นฐาน
- 1 โฟลเดอร์ = 1 package, ทุกไฟล์ในโฟลเดอร์ต้องเป็น package เดียวกัน
- **Module** = กลุ่มของ package ที่ version ไปด้วยกัน กำหนดด้วย `go.mod`
- import path = module path + โฟลเดอร์: `github.com/mpcndr/go-learn/12-packages-design/greeting`
- ตัวพิมพ์ใหญ่ = export (public), ตัวเล็ก = ใช้ได้เฉพาะใน package
- **ห้าม import วน** (A → B → A) — compiler ไม่ยอม นี่คือสิ่งที่บังคับให้ออกแบบชั้นดี ๆ
- โฟลเดอร์ชื่อ `internal/` → import ได้เฉพาะจาก package ที่อยู่ใต้ parent ของมัน (ซ่อนจากโลกภายนอก)

## ตั้งชื่อ package
- สั้น ตัวเล็ก คำเดียว ไม่มี `_`: `http`, `json`, `library`, `memory`
- **ชื่อ package เป็นส่วนหนึ่งของชื่อ** → อย่าพูดซ้ำ
  - ✅ `http.Client`, `memory.NewStore()`
  - ❌ `http.HTTPClient`, `memory.NewMemoryStore()`
- ❌ อย่าตั้ง `utils`, `common`, `helpers`, `models` — ไม่บอกว่ารับผิดชอบอะไร กลายเป็นถังขยะ

## การแบ่งชั้น (Layered / Hexagonal แบบย่อ)

```
            cmd/library (main)          ← ประกอบทุกอย่างเข้าด้วยกัน (wiring)
             │        │
             ▼        ▼
     storage/memory   service           ← use case: กฎธุรกิจที่ต้องใช้หลาย entity
             │        │   ▲
             │        │   └── interface ที่ service "ต้องการ" (Books, Loans...)
             ▼        ▼       storage/memory implement โดยไม่ต้อง import service
              domain                    ← ข้อมูล + กฎพื้นฐาน, ไม่ import อะไรในโปรเจกต์เลย
```

**กฎ:**
1. `domain` ไม่ import ใคร (มีแค่ standard library)
2. `service` import แค่ `domain` และประกาศ interface ของสิ่งที่ต้องการเอง
3. `storage` implement interface เหล่านั้น (implicit — บท 06!) import แค่ `domain`
4. `main` เป็นที่เดียวที่รู้จักทุกคน และ "ต่อสาย" (dependency injection ด้วยมือ)

ผลลัพธ์:
- เปลี่ยน memory → PostgreSQL = เขียน package ใหม่ 1 ตัว + แก้ `main` 1 บรรทัด
- test `service` ได้โดยไม่ต้องมี database
- อ่าน `domain` อย่างเดียวก็เข้าใจว่าระบบทำอะไร

## โครงสร้างโปรเจกต์ที่พบบ่อย
```
myapp/
├── go.mod
├── cmd/
│   ├── api/main.go        ← แต่ละโปรแกรม (binary) มี main ของตัวเอง
│   └── worker/main.go
├── internal/              ← โค้ดที่ไม่อยากให้โปรเจกต์อื่น import
│   ├── domain/
│   ├── service/
│   └── storage/postgres/
└── README.md
```
**อย่าเริ่มจากโครงใหญ่** — โปรเจกต์เล็กเริ่มจาก `main.go` ไฟล์เดียวได้เลย
แยก package เมื่อ "เห็นเส้นแบ่งความรับผิดชอบชัด" แล้วเท่านั้น

## ตัวอย่างในบทนี้
```bash
go run ./12-packages-design
```
- `greeting/` — package ที่ export API เล็ก ๆ
- `greeting/internal/format/` — package ที่ถูกซ่อน ใช้ได้เฉพาะใน `greeting`
  (ลอง import จาก `main.go` ดู — compiler จะไม่ยอม)

## แบบฝึกหัด (บนกระดาษ ก่อนไปบท 13)
อ่านโจทย์ใน [บท 13](../13-capstone-library) ส่วน "① เข้าใจโจทย์" แล้ว **ออกแบบเองก่อน** โดยยังไม่เปิดดูโค้ด:
1. ขีดเส้นใต้คำนาม → เขียน struct ของ domain
2. ขีดเส้นใต้คำกริยา → เขียน signature ของ service พร้อม error ที่เป็นไปได้
3. service ต้องการอะไรจาก storage? → เขียน interface
4. เขียนชื่อ test case อย่างน้อย 8 ข้อจากกฎในโจทย์

แล้วค่อยเปิดโค้ดบท 13 มาเทียบ — ต่างกันตรงไหน? ทำไม?
