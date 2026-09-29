# 13 — Capstone: ระบบห้องสมุด (จากหัว → โค้ด)

บทนี้ไม่ได้สอน syntax ใหม่ แต่สาธิต **กระบวนการคิด** จาก [docs/thinking-in-go.md](../docs/thinking-in-go.md)
ตั้งแต่ประโยคโจทย์ ไปจนถึงโค้ดที่รันและทดสอบได้

> ⚠️ **ทำบท 12 ส่วนแบบฝึกหัดก่อน** — ออกแบบบนกระดาษเองก่อน แล้วค่อยอ่านต่อ
> การเทียบ "ที่เราคิด" กับ "ตัวอย่าง" คือจุดที่เรียนรู้มากที่สุด

```bash
go run ./13-capstone-library/cmd/library        # ลองเล่น
go test -race ./13-capstone-library/...          # รัน test ทั้งหมด
```

---

## ① เข้าใจโจทย์ — เขียนเป็นประโยค

> ห้องสมุดมี**หนังสือ**และ**สมาชิก** สมาชิก**ยืม**หนังสือได้ไม่เกิน **3 เล่ม**
> หนังสือเล่มหนึ่งถูกยืมได้**ทีละคน** ต้อง**คืน**ภายใน **14 วัน**
> ถ้าคืนช้ามี**ค่าปรับ**วันละ **5 บาท** (เกินแม้ไม่ถึงวันก็นับเป็นวัน)
> บรรณารักษ์**เพิ่มหนังสือ** **ลงทะเบียนสมาชิก** และ**ดูรายการ**หนังสือพร้อมสถานะได้

## ② คำนาม → ข้อมูล (`domain/`)

| คำนาม | ตัดสินใจ | เหตุผล |
|-------|----------|--------|
| หนังสือ | `Book{ID, Title, Author}` | |
| สมาชิก | `Member{ID, Name}` | |
| การยืม | `Loan{BookID, MemberID, BorrowedAt, DueAt}` | เป็น "ความสัมพันธ์" ระหว่าง Book กับ Member ที่มีข้อมูลของตัวเอง |
| ค่าปรับ | `Money` (int64 สตางค์) | **ห้ามใช้ float กับเงิน** |
| ID | `BookID`, `MemberID` (type แยก) | compiler กันการส่ง ID สลับกัน |
| 3 เล่ม, 14 วัน, 5 บาท | `const` | กฎอยู่ที่เดียว เปลี่ยนง่าย |

**คำถามที่ต้องตอบตอนออกแบบข้อมูล:**
- *"สถานะว่าง/ถูกยืม" เก็บที่ไหน?* → **ไม่เก็บใน Book** แต่ดูจากว่ามี `Loan` ที่ active ไหม
  (ถ้าเก็บสองที่ สักวันมันจะไม่ตรงกัน — *single source of truth*)
- *คำนวณค่าปรับที่ไหน?* → `Loan.FineAt(returnedAt)` เป็น **pure function** รับเวลาเข้ามา
  ไม่เรียก `time.Now()` เอง → ทดสอบได้ทุกกรณีในไม่กี่บรรทัด
- *ข้อมูลผิดตั้งแต่เกิดได้ไหม?* → ใช้ `NewBook`/`NewMember` ที่ validate → Book ทุกตัวในระบบถูกต้องเสมอ

## ③ คำกริยา → signature + ความล้มเหลว (`service/`)

เขียน signature และ **error ทุกตัวที่เป็นไปได้** ก่อนเขียนเนื้อใน:

| กริยา | signature | ล้มเหลวได้เมื่อ |
|-------|-----------|----------------|
| เพิ่มหนังสือ | `AddBook(ctx, id, title, author) (Book, error)` | `ValidationError`, `ErrAlreadyExists` |
| ลงทะเบียน | `RegisterMember(ctx, id, name) (Member, error)` | `ValidationError`, `ErrAlreadyExists` |
| ยืม | `Borrow(ctx, memberID, bookID) (Loan, error)` | `ErrMemberNotFound`, `ErrBookNotFound`, `ErrBookUnavailable`, `ErrLoanLimitReached` |
| คืน | `Return(ctx, bookID) (Money, error)` | `ErrLoanNotFound` |
| ดูรายการ | `Catalog(ctx) ([]BookStatus, error)` | error จาก storage |

error เหล่านี้ถูกประกาศใน `domain` เป็น **sentinel** เพราะผู้เรียก (CLI) ต้องแสดงข้อความต่างกันตามชนิด
และ `ValidationError` เป็น **custom type** เพราะต้องบอกได้ว่า field ไหนผิด (บท 07)

## ④ ขีดเส้นแบ่ง → interface

service ต้องการอะไรจาก storage? เขียน **เท่าที่ต้องใช้จริง** ไว้ใน `service/library.go`:

```go
type Loans interface {
    ActiveLoan(ctx, bookID) (Loan, error)
    LoansByMember(ctx, memberID) ([]Loan, error)
    SaveLoan(ctx, Loan) error
    DeleteLoan(ctx, bookID) error
}
```

แยกเป็น `Books`, `Members`, `Loans` (เล็ก ๆ) แทน `Repository` ก้อนใหญ่ → อนาคตถ้าจะเก็บ Loans ใน Redis แต่ Books ใน Postgres ก็ทำได้

**ผังความสัมพันธ์:**
```
cmd/library ──▶ service ──▶ domain
     │                        ▲
     └────▶ storage/memory ───┘
```
- `domain` ไม่ import ใคร
- `storage/memory` **ไม่ import service** — แต่ implement interface ของ service ได้ (implicit!)
- `cmd/library` คนเดียวที่รู้จักทุกคน และต่อสายเข้าด้วยกัน

## ⑤ เขียน test จากประโยคโจทย์

ดู `service/library_test.go` — ชื่อ test แทบจะ copy มาจากโจทย์:
```
TestBorrowRules/a_book_can_be_borrowed_by_one_member_at_a_time
TestBorrowRules/member_can_borrow_at_most_3_books
TestReturnFines/a_few_hours_late_counts_as_a_day
TestConcurrentBorrowSameBook      ← 50 คนแย่งเล่มเดียว ต้องได้คนเดียว
TestBorrowStorageFailure          ← storage พัง error ต้องเล่าเรื่องได้
```
ใช้ `fakeClock` แทนเวลาจริง → ทดสอบ "คืนช้า 16 วัน" ได้ในเสี้ยววินาที

## ⑥ เขียนโค้ดจากในออกนอก
1. `domain/` + `domain_test.go` — รัน test ผ่าน
2. `service/` + `library_test.go` — ใช้ `storage/memory` เป็นตัวเก็บข้อมูลใน test
3. `storage/memory/` — map + `sync.RWMutex`
4. `cmd/library/` — ขอบนอก: แปลงข้อความ ↔ การเรียก service และแปลง error เป็นภาษาคน

---

## จุดที่ควรสังเกตในโค้ด

| จุด | ไฟล์ | บทที่เกี่ยวข้อง |
|-----|------|----------------|
| type แยกสำหรับ ID | `domain/domain.go` | 05 |
| constructor ที่ validate | `domain.NewBook` | 05, 07 |
| pure function สำหรับกฎ | `Loan.FineAt` | 02 |
| interface ฝั่งผู้ใช้ | `service.Books/Members/Loans` | 06 |
| inject เวลา | `service.New(..., now func() time.Time)` | 06, 09 |
| wrap error ทุกชั้น | `fmt.Errorf("borrow: %w", err)` | 07 |
| mutex ทำให้ check-then-act เป็น atomic | `Library.mu` | 10 |
| `RWMutex` ใน storage | `memory.Store` | 10 |
| `context.Context` เป็น parameter แรก | ทุก method | 11 |
| embed interface เพื่อทำ fake บางส่วน | `brokenLoans` ใน test | 06, 09 |
| แปลง error เป็นข้อความที่ขอบนอก | `describe()` ใน main | 07 |

---

## โจทย์ต่อยอด (เรียงจากง่ายไปยาก)

ทุกข้อให้ทำตามกระบวนการ ①–⑥: **เขียนประโยค → แก้ domain → signature + error → test → code**

1. **ต่ออายุการยืม** `Renew(ctx, bookID)` — ต่อได้ 1 ครั้ง, ต่อไม่ได้ถ้าเลยกำหนดแล้ว
2. **ประวัติการยืม** — เก็บ Loan ที่คืนแล้วไว้ด้วย พร้อม `ReturnedAt` และ `Fine`
3. **ค่าปรับค้าง** — สมาชิกที่มีค่าปรับค้างเกิน 50 บาท ยืมใหม่ไม่ได้
4. **จองหนังสือ** — ถ้าถูกยืมอยู่ ให้เข้าคิวจองได้ คนที่จองก่อนได้สิทธิ์ยืมก่อนเมื่อหนังสือถูกคืน
5. **เก็บลงไฟล์ JSON** — เขียน `storage/jsonfile` ที่ implement interface เดิม
   แก้ `main.go` แค่บรรทัดเดียว (`encoding/json`, `os`) — service ต้องไม่ถูกแก้เลย
6. **HTTP API** — เขียน `cmd/libraryd` ด้วย `net/http`
   (`GET /books`, `POST /loans`, `DELETE /loans/{bookID}`) ใช้ service ตัวเดิม
   แปลง error เป็น HTTP status (404, 409, 422) — คล้าย `describe()` ใน CLI
7. **Graceful shutdown** — HTTP server หยุดรับ request ใหม่เมื่อกด Ctrl+C แล้วรอ request ที่ค้างอยู่ให้จบ
   (`signal.NotifyContext`, `http.Server.Shutdown`)

> ถ้าข้อ 5–6 ทำได้โดย **ไม่แก้ `domain/` และ `service/` เลย** แปลว่าคุณเข้าใจการแบ่งชั้นแล้ว 🎉
