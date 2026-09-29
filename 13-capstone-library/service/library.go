// Package service คือ use case ของระบบห้องสมุด: กฎที่ต้องใช้ข้อมูลหลายชิ้นร่วมกัน
// เช่น "ยืมได้ไม่เกิน 3 เล่ม" ต้องรู้ทั้งสมาชิกและรายการยืมทั้งหมด
//
// service ไม่รู้ว่าข้อมูลเก็บที่ไหน — มันประกาศ interface ของสิ่งที่ "ต้องการ" ไว้ที่นี่
// แล้วให้ main เป็นคนส่ง implementation จริงเข้ามา
package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/mpcndr/go-learn/13-capstone-library/domain"
)

// ---------- สิ่งที่ service ต้องการจากโลกภายนอก (interface เล็ก ๆ ฝั่งผู้ใช้) ----------

type Books interface {
	Book(ctx context.Context, id domain.BookID) (domain.Book, error) // ไม่พบ → domain.ErrBookNotFound
	AddBook(ctx context.Context, b domain.Book) error                // ซ้ำ → domain.ErrAlreadyExists
	ListBooks(ctx context.Context) ([]domain.Book, error)
}

type Members interface {
	Member(ctx context.Context, id domain.MemberID) (domain.Member, error) // ไม่พบ → domain.ErrMemberNotFound
	AddMember(ctx context.Context, m domain.Member) error                  // ซ้ำ → domain.ErrAlreadyExists
}

type Loans interface {
	ActiveLoan(ctx context.Context, book domain.BookID) (domain.Loan, error) // ไม่พบ → domain.ErrLoanNotFound
	LoansByMember(ctx context.Context, member domain.MemberID) ([]domain.Loan, error)
	SaveLoan(ctx context.Context, l domain.Loan) error
	DeleteLoan(ctx context.Context, book domain.BookID) error
}

// ---------- Library ----------

// Library รวม use case ทั้งหมด
type Library struct {
	books   Books
	members Members
	loans   Loans
	now     func() time.Time // inject เวลาเข้ามา → test เรื่องค่าปรับได้โดยไม่ต้องรอ 14 วันจริง

	// mu ทำให้ขั้นตอน "ตรวจ แล้วค่อยบันทึก" ของ Borrow/Return เป็น atomic
	// ถ้าไม่มี: สองคนยืมเล่มเดียวกันพร้อมกัน อาจผ่านการตรวจทั้งคู่ (race)
	// (ถ้าใช้ database จริง หน้าที่นี้จะย้ายไปอยู่ใน transaction แทน)
	mu sync.Mutex
}

// New ประกอบ Library จาก dependency ที่ส่งเข้ามา — ถ้า now เป็น nil ใช้ time.Now
func New(books Books, members Members, loans Loans, now func() time.Time) *Library {
	if now == nil {
		now = time.Now
	}
	return &Library{books: books, members: members, loans: loans, now: now}
}

// AddBook เพิ่มหนังสือใหม่เข้าระบบ
func (l *Library) AddBook(ctx context.Context, id domain.BookID, title, author string) (domain.Book, error) {
	b, err := domain.NewBook(id, title, author)
	if err != nil {
		return domain.Book{}, fmt.Errorf("add book: %w", err)
	}
	if err := l.books.AddBook(ctx, b); err != nil {
		return domain.Book{}, fmt.Errorf("add book %s: %w", id, err)
	}
	return b, nil
}

// RegisterMember ลงทะเบียนสมาชิกใหม่
func (l *Library) RegisterMember(ctx context.Context, id domain.MemberID, name string) (domain.Member, error) {
	m, err := domain.NewMember(id, name)
	if err != nil {
		return domain.Member{}, fmt.Errorf("register member: %w", err)
	}
	if err := l.members.AddMember(ctx, m); err != nil {
		return domain.Member{}, fmt.Errorf("register member %s: %w", id, err)
	}
	return m, nil
}

// Borrow ให้สมาชิกยืมหนังสือ
//
// ล้มเหลวเมื่อ: ไม่พบสมาชิก/หนังสือ, หนังสือถูกยืมอยู่, สมาชิกยืมครบ MaxLoansPerMember แล้ว
func (l *Library) Borrow(ctx context.Context, memberID domain.MemberID, bookID domain.BookID) (domain.Loan, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 1. ทุกคนต้องมีตัวตน
	if _, err := l.members.Member(ctx, memberID); err != nil {
		return domain.Loan{}, fmt.Errorf("borrow: %w", err)
	}
	if _, err := l.books.Book(ctx, bookID); err != nil {
		return domain.Loan{}, fmt.Errorf("borrow: %w", err)
	}

	// 2. หนังสือต้องว่าง
	_, err := l.loans.ActiveLoan(ctx, bookID)
	switch {
	case err == nil:
		return domain.Loan{}, fmt.Errorf("borrow %s: %w", bookID, domain.ErrBookUnavailable)
	case !errors.Is(err, domain.ErrLoanNotFound):
		return domain.Loan{}, fmt.Errorf("borrow %s: check availability: %w", bookID, err)
	}

	// 3. สมาชิกต้องยังไม่เกินโควตา
	current, err := l.loans.LoansByMember(ctx, memberID)
	if err != nil {
		return domain.Loan{}, fmt.Errorf("borrow: list loans of %s: %w", memberID, err)
	}
	if len(current) >= domain.MaxLoansPerMember {
		return domain.Loan{}, fmt.Errorf("borrow: member %s: %w", memberID, domain.ErrLoanLimitReached)
	}

	// 4. ทุกเงื่อนไขผ่าน → บันทึก
	loan := domain.NewLoan(bookID, memberID, l.now())
	if err := l.loans.SaveLoan(ctx, loan); err != nil {
		return domain.Loan{}, fmt.Errorf("borrow: save loan: %w", err)
	}
	return loan, nil
}

// Return คืนหนังสือ และคืนค่าปรับ (0 ถ้าไม่ช้า)
func (l *Library) Return(ctx context.Context, bookID domain.BookID) (domain.Money, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	loan, err := l.loans.ActiveLoan(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("return %s: %w", bookID, err)
	}
	fine := loan.FineAt(l.now())
	if err := l.loans.DeleteLoan(ctx, bookID); err != nil {
		return 0, fmt.Errorf("return %s: %w", bookID, err)
	}
	return fine, nil
}

// BookStatus หนังสือพร้อมสถานะการยืม — เป็น "view" สำหรับแสดงผล
type BookStatus struct {
	Book domain.Book
	Loan *domain.Loan // nil = ว่าง
}

// Catalog คืนหนังสือทั้งหมดพร้อมสถานะ เรียงตาม ID
func (l *Library) Catalog(ctx context.Context) ([]BookStatus, error) {
	books, err := l.books.ListBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("catalog: %w", err)
	}
	slices.SortFunc(books, func(a, b domain.Book) int { return strings.Compare(string(a.ID), string(b.ID)) })

	out := make([]BookStatus, 0, len(books))
	for _, b := range books {
		st := BookStatus{Book: b}
		loan, err := l.loans.ActiveLoan(ctx, b.ID)
		switch {
		case err == nil:
			st.Loan = &loan
		case !errors.Is(err, domain.ErrLoanNotFound):
			return nil, fmt.Errorf("catalog: %w", err)
		}
		out = append(out, st)
	}
	return out, nil
}
