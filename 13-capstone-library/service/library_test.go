// ใช้ package service_test (black-box) — ทดสอบผ่าน API ที่ export เท่านั้น เหมือนผู้ใช้จริง
package service_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/mpcndr/go-learn/13-capstone-library/domain"
	"github.com/mpcndr/go-learn/13-capstone-library/service"
	"github.com/mpcndr/go-learn/13-capstone-library/storage/memory"
)

// ตรวจตอน compile ว่า memory.Store เข้ากับ interface ที่ service ต้องการ
var (
	_ service.Books   = (*memory.Store)(nil)
	_ service.Members = (*memory.Store)(nil)
	_ service.Loans   = (*memory.Store)(nil)
)

// fakeClock เวลาที่เราควบคุมได้ — เลื่อนไปข้างหน้าได้ตามใจ
type fakeClock struct{ t time.Time }

func (c *fakeClock) Now() time.Time          { return c.t }
func (c *fakeClock) Advance(d time.Duration) { c.t = c.t.Add(d) }

const day = 24 * time.Hour

// setup สร้าง Library ที่มีหนังสือ B1..B5 และสมาชิก ann, ben
func setup(t *testing.T) (*service.Library, *fakeClock) {
	t.Helper()
	ctx := context.Background()
	clock := &fakeClock{t: time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)}
	store := memory.NewStore()
	lib := service.New(store, store, store, clock.Now)

	for _, id := range []domain.BookID{"B1", "B2", "B3", "B4", "B5"} {
		if _, err := lib.AddBook(ctx, id, "Book "+string(id), "Author"); err != nil {
			t.Fatalf("AddBook(%s): %v", id, err)
		}
	}
	for _, id := range []domain.MemberID{"ann", "ben"} {
		if _, err := lib.RegisterMember(ctx, id, string(id)); err != nil {
			t.Fatalf("RegisterMember(%s): %v", id, err)
		}
	}
	return lib, clock
}

func TestBorrowRules(t *testing.T) {
	ctx := context.Background()

	t.Run("borrow available book", func(t *testing.T) {
		lib, clock := setup(t)
		loan, err := lib.Borrow(ctx, "ann", "B1")
		if err != nil {
			t.Fatalf("Borrow: %v", err)
		}
		if want := clock.Now().Add(14 * day); !loan.DueAt.Equal(want) {
			t.Errorf("DueAt = %v, want %v", loan.DueAt, want)
		}
	})

	t.Run("a book can be borrowed by one member at a time", func(t *testing.T) {
		lib, _ := setup(t)
		mustBorrow(t, lib, "ann", "B1")
		_, err := lib.Borrow(ctx, "ben", "B1")
		if !errors.Is(err, domain.ErrBookUnavailable) {
			t.Errorf("second Borrow err = %v, want ErrBookUnavailable", err)
		}
	})

	t.Run("member can borrow at most 3 books", func(t *testing.T) {
		lib, _ := setup(t)
		mustBorrow(t, lib, "ann", "B1")
		mustBorrow(t, lib, "ann", "B2")
		mustBorrow(t, lib, "ann", "B3")
		_, err := lib.Borrow(ctx, "ann", "B4")
		if !errors.Is(err, domain.ErrLoanLimitReached) {
			t.Errorf("4th Borrow err = %v, want ErrLoanLimitReached", err)
		}
		// คนอื่นยังยืมได้ตามปกติ
		mustBorrow(t, lib, "ben", "B4")
	})

	t.Run("returning frees a slot", func(t *testing.T) {
		lib, _ := setup(t)
		mustBorrow(t, lib, "ann", "B1")
		mustBorrow(t, lib, "ann", "B2")
		mustBorrow(t, lib, "ann", "B3")
		if _, err := lib.Return(ctx, "B2"); err != nil {
			t.Fatalf("Return: %v", err)
		}
		mustBorrow(t, lib, "ann", "B4")
	})

	t.Run("unknown member or book", func(t *testing.T) {
		lib, _ := setup(t)
		if _, err := lib.Borrow(ctx, "nobody", "B1"); !errors.Is(err, domain.ErrMemberNotFound) {
			t.Errorf("unknown member err = %v, want ErrMemberNotFound", err)
		}
		if _, err := lib.Borrow(ctx, "ann", "B999"); !errors.Is(err, domain.ErrBookNotFound) {
			t.Errorf("unknown book err = %v, want ErrBookNotFound", err)
		}
	})
}

func TestReturnFines(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name  string
		after time.Duration
		want  domain.Money
	}{
		{"on time", 10 * day, 0},
		{"exactly due", 14 * day, 0},
		{"2 days late", 16 * day, 1000},
		{"a few hours late counts as a day", 14*day + 3*time.Hour, 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lib, clock := setup(t)
			mustBorrow(t, lib, "ann", "B1")
			clock.Advance(tt.after)
			fine, err := lib.Return(ctx, "B1")
			if err != nil {
				t.Fatalf("Return: %v", err)
			}
			if fine != tt.want {
				t.Errorf("fine after %v = %v, want %v", tt.after, fine, tt.want)
			}
		})
	}

	t.Run("return a book that is not borrowed", func(t *testing.T) {
		lib, _ := setup(t)
		if _, err := lib.Return(ctx, "B1"); !errors.Is(err, domain.ErrLoanNotFound) {
			t.Errorf("err = %v, want ErrLoanNotFound", err)
		}
	})
}

func TestAddValidation(t *testing.T) {
	ctx := context.Background()
	lib, _ := setup(t)

	var ve *domain.ValidationError
	if _, err := lib.AddBook(ctx, "B9", "   ", "x"); !errors.As(err, &ve) || ve.Field != "title" {
		t.Errorf("empty title err = %v, want ValidationError(title)", err)
	}
	if _, err := lib.AddBook(ctx, "B1", "Dup", "x"); !errors.Is(err, domain.ErrAlreadyExists) {
		t.Errorf("duplicate book err = %v, want ErrAlreadyExists", err)
	}
	if _, err := lib.RegisterMember(ctx, "ann", "Ann again"); !errors.Is(err, domain.ErrAlreadyExists) {
		t.Errorf("duplicate member err = %v, want ErrAlreadyExists", err)
	}
}

func TestCatalog(t *testing.T) {
	ctx := context.Background()
	lib, _ := setup(t)
	mustBorrow(t, lib, "ben", "B3")

	cat, err := lib.Catalog(ctx)
	if err != nil {
		t.Fatalf("Catalog: %v", err)
	}
	if len(cat) != 5 || cat[0].Book.ID != "B1" || cat[4].Book.ID != "B5" {
		t.Fatalf("Catalog not sorted or wrong size: %+v", cat)
	}
	for _, st := range cat {
		onLoan := st.Loan != nil
		if want := st.Book.ID == "B3"; onLoan != want {
			t.Errorf("%s on loan = %v, want %v", st.Book.ID, onLoan, want)
		}
	}
}

// หลายคนแย่งยืมเล่มเดียวกันพร้อมกัน — ต้องได้แค่คนเดียว (รันด้วย -race)
func TestConcurrentBorrowSameBook(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	lib := service.New(store, store, store, nil)
	if _, err := lib.AddBook(ctx, "HOT", "Popular Book", "Famous"); err != nil {
		t.Fatal(err)
	}
	const n = 50
	for i := range n {
		if _, err := lib.RegisterMember(ctx, memberID(i), "m"); err != nil {
			t.Fatal(err)
		}
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		success int
	)
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := lib.Borrow(ctx, memberID(i), "HOT"); err == nil {
				mu.Lock()
				success++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if success != 1 {
		t.Errorf("successful borrows = %d, want exactly 1", success)
	}
}

// ---------- เมื่อ storage พัง: error ต้องถูกส่งต่อพร้อมบริบท ----------

type brokenLoans struct{ service.Loans }

var errDiskFull = errors.New("disk full")

func (brokenLoans) SaveLoan(context.Context, domain.Loan) error { return errDiskFull }

func TestBorrowStorageFailure(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	// embed interface แล้ว override เฉพาะ method ที่อยากให้พัง
	lib := service.New(store, store, brokenLoans{store}, nil)
	if _, err := lib.AddBook(ctx, "B1", "Go", "x"); err != nil {
		t.Fatal(err)
	}
	if _, err := lib.RegisterMember(ctx, "ann", "Ann"); err != nil {
		t.Fatal(err)
	}
	_, err := lib.Borrow(ctx, "ann", "B1")
	if !errors.Is(err, errDiskFull) {
		t.Fatalf("err = %v, want wrapped errDiskFull", err)
	}
	t.Logf("error message tells the whole story: %v", err)
}

func memberID(i int) domain.MemberID { return domain.MemberID(fmt.Sprintf("m%d", i)) }

func mustBorrow(t *testing.T, lib *service.Library, m domain.MemberID, b domain.BookID) {
	t.Helper()
	if _, err := lib.Borrow(context.Background(), m, b); err != nil {
		t.Fatalf("Borrow(%s, %s): %v", m, b, err)
	}
}
