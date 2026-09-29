// Command library — CLI ของระบบห้องสมุด
//
//	go run ./13-capstone-library/cmd/library
//
// main มีหน้าที่แค่ 2 อย่าง:
//  1. ประกอบ (wire) ชิ้นส่วน: storage → service
//  2. แปลงสิ่งที่ผู้ใช้พิมพ์ ↔ การเรียก service (นี่คือ "ขอบนอก" ของระบบ)
//
// กฎธุรกิจทั้งหมดอยู่ใน domain/service — ถ้าวันหนึ่งเปลี่ยนจาก CLI เป็น HTTP API
// จะเขียนแค่ขอบนอกใหม่ โดยไม่ต้องแตะกฎเลย
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mpcndr/go-learn/13-capstone-library/domain"
	"github.com/mpcndr/go-learn/13-capstone-library/service"
	"github.com/mpcndr/go-learn/13-capstone-library/storage/memory"
)

// clock นาฬิกาที่หยุดนิ่ง เลื่อนได้ด้วยคำสั่ง "skip" เท่านั้น (ไว้ลองเรื่องค่าปรับให้ได้ผลแน่นอน)
// ใน production จะส่ง time.Now เข้า service.New แทน — service ไม่ต้องรู้ความต่างนี้เลย
type clock struct{ now time.Time }

func (c *clock) Now() time.Time { return c.now }

const help = `คำสั่ง:
  book   <id> <title>        เพิ่มหนังสือ
  member <id> <name>         ลงทะเบียนสมาชิก
  borrow <member-id> <book-id>
  return <book-id>
  list                       แสดงหนังสือทั้งหมดและสถานะ
  skip   <days>              เลื่อนเวลาไปข้างหน้า (ทดลองค่าปรับ)
  help
  quit`

func main() {
	// ---- 1. wiring ----
	store := memory.NewStore()
	clk := &clock{now: time.Now()}
	lib := service.New(store, store, store, clk.Now)

	ctx := context.Background()
	if err := seed(ctx, lib); err != nil {
		fmt.Fprintln(os.Stderr, "seed:", err)
		os.Exit(1)
	}

	// ---- 2. REPL ----
	fmt.Println("📚 ระบบห้องสมุด — พิมพ์ help เพื่อดูคำสั่ง")
	app := &app{lib: lib, clock: clk, out: os.Stdout}
	sc := bufio.NewScanner(os.Stdin)
	for fmt.Print("> "); sc.Scan(); fmt.Print("> ") {
		if quit := app.run(ctx, sc.Text()); quit {
			return
		}
	}
}

type app struct {
	lib   *service.Library
	clock *clock
	out   io.Writer
}

// run ประมวลผล 1 บรรทัด คืน true เมื่อต้องการออก
func (a *app) run(ctx context.Context, line string) (quit bool) {
	cmd, rest, _ := strings.Cut(strings.TrimSpace(line), " ")
	arg1, arg2, _ := strings.Cut(strings.TrimSpace(rest), " ")

	var err error
	switch cmd {
	case "":
	case "help":
		fmt.Fprintln(a.out, help)
	case "quit", "exit":
		return true
	case "book":
		var b domain.Book
		if b, err = a.lib.AddBook(ctx, domain.BookID(arg1), arg2, ""); err == nil {
			fmt.Fprintf(a.out, "✅ เพิ่ม %s %q\n", b.ID, b.Title)
		}
	case "member":
		var m domain.Member
		if m, err = a.lib.RegisterMember(ctx, domain.MemberID(arg1), arg2); err == nil {
			fmt.Fprintf(a.out, "✅ สมาชิก %s (%s)\n", m.ID, m.Name)
		}
	case "borrow":
		var l domain.Loan
		if l, err = a.lib.Borrow(ctx, domain.MemberID(arg1), domain.BookID(arg2)); err == nil {
			fmt.Fprintf(a.out, "✅ %s ยืม %s — คืนภายใน %s\n", l.MemberID, l.BookID, l.DueAt.Format("2006-01-02"))
		}
	case "return":
		var fine domain.Money
		if fine, err = a.lib.Return(ctx, domain.BookID(arg1)); err == nil {
			if fine > 0 {
				fmt.Fprintf(a.out, "⚠️  คืนช้า ค่าปรับ %v\n", fine)
			} else {
				fmt.Fprintln(a.out, "✅ คืนเรียบร้อย ไม่มีค่าปรับ")
			}
		}
	case "list":
		err = a.list(ctx)
	case "skip":
		var days int
		if days, err = strconv.Atoi(arg1); err == nil {
			a.clock.now = a.clock.now.AddDate(0, 0, days)
			fmt.Fprintf(a.out, "⏩ วันนี้คือ %s\n", a.clock.Now().Format("2006-01-02"))
		}
	default:
		fmt.Fprintf(a.out, "ไม่รู้จักคำสั่ง %q — พิมพ์ help\n", cmd)
	}

	if err != nil {
		fmt.Fprintln(a.out, "❌", describe(err))
	}
	return false
}

func (a *app) list(ctx context.Context) error {
	cat, err := a.lib.Catalog(ctx)
	if err != nil {
		return err
	}
	for _, st := range cat {
		status := "ว่าง"
		if st.Loan != nil {
			status = fmt.Sprintf("ยืมโดย %s (กำหนดคืน %s)", st.Loan.MemberID, st.Loan.DueAt.Format("2006-01-02"))
		}
		fmt.Fprintf(a.out, "  %-4s %-36q %s\n", st.Book.ID, st.Book.Title, status)
	}
	return nil
}

// describe แปลง error เป็นข้อความสำหรับผู้ใช้ — ขอบนอกคือที่เดียวที่ "ตัดสินใจ" เรื่องนี้
func describe(err error) string {
	var ve *domain.ValidationError
	switch {
	case errors.As(err, &ve):
		return fmt.Sprintf("ข้อมูลไม่ถูกต้อง: %s (%s)", ve.Field, ve.Reason)
	case errors.Is(err, domain.ErrBookUnavailable):
		return "หนังสือเล่มนี้ถูกยืมอยู่"
	case errors.Is(err, domain.ErrLoanLimitReached):
		return fmt.Sprintf("ยืมได้ไม่เกิน %d เล่ม", domain.MaxLoansPerMember)
	case errors.Is(err, domain.ErrBookNotFound):
		return "ไม่พบหนังสือ"
	case errors.Is(err, domain.ErrMemberNotFound):
		return "ไม่พบสมาชิก"
	case errors.Is(err, domain.ErrLoanNotFound):
		return "หนังสือเล่มนี้ไม่ได้ถูกยืม"
	case errors.Is(err, domain.ErrAlreadyExists):
		return "มีอยู่แล้ว"
	default:
		return err.Error()
	}
}

func seed(ctx context.Context, lib *service.Library) error {
	books := []struct{ id, title, author string }{
		{"B1", "The Go Programming Language", "Donovan & Kernighan"},
		{"B2", "Concurrency in Go", "Katherine Cox-Buday"},
		{"B3", "Learning Go", "Jon Bodner"},
		{"B4", "100 Go Mistakes", "Teiva Harsanyi"},
	}
	for _, b := range books {
		if _, err := lib.AddBook(ctx, domain.BookID(b.id), b.title, b.author); err != nil {
			return err
		}
	}
	for _, m := range []struct{ id, name string }{{"ann", "Ann"}, {"ben", "Ben"}} {
		if _, err := lib.RegisterMember(ctx, domain.MemberID(m.id), m.name); err != nil {
			return err
		}
	}
	return nil
}
