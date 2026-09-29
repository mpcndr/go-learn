package showcase

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// ---------- fake: แทน dependency จริงด้วยของที่ควบคุมได้ ----------

type fakeClock struct{ now time.Time }

func (f fakeClock) Now() time.Time { return f.now }

var (
	monday   = time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	saturday = time.Date(2026, 10, 3, 12, 0, 0, 0, time.UTC)
)

// ---------- helper: t.Helper() ทำให้ error ชี้ไปที่บรรทัดที่เรียก ----------

func newTestPricer(t *testing.T, now time.Time) *Pricer {
	t.Helper()
	return &Pricer{
		Clock:   fakeClock{now},
		Coupons: map[string]int{"SAVE20": 20, "ALL": 100},
	}
}

// ---------- table-driven + subtests ----------

func TestPricerFinal(t *testing.T) {
	tests := []struct {
		name   string
		now    time.Time
		price  int
		coupon string
		want   int
	}{
		{"weekday no coupon", monday, 1000, "", 1000},
		{"weekend discount", saturday, 1000, "", 900},
		{"coupon on weekday", monday, 1000, "SAVE20", 800},
		{"coupon stacks with weekend", saturday, 1000, "SAVE20", 700},
		{"unknown coupon ignored", monday, 1000, "BOGUS", 1000},
		{"discount capped at 100%", saturday, 1000, "ALL", 0},
		{"zero price", monday, 0, "SAVE20", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // subtest อิสระต่อกัน → รันพร้อมกันได้
			p := newTestPricer(t, tt.now)
			got, err := p.Final(tt.price, tt.coupon)
			if err != nil {
				t.Fatalf("Final(%d, %q) unexpected error: %v", tt.price, tt.coupon, err)
			}
			if got != tt.want {
				t.Errorf("Final(%d, %q) on %s = %d, want %d",
					tt.price, tt.coupon, tt.now.Weekday(), got, tt.want)
			}
		})
	}
}

// ---------- ทดสอบ error path แยกให้ชัด ----------

func TestPricerFinalNegative(t *testing.T) {
	p := newTestPricer(t, monday)
	_, err := p.Final(-1, "")
	if !errors.Is(err, ErrNegativePrice) {
		t.Errorf("Final(-1) err = %v, want ErrNegativePrice", err)
	}
}

// ---------- Example: เป็นทั้ง doc และ test (ตรวจ // Output:) ----------

func ExamplePricer_Final() {
	p := &Pricer{
		Clock:   fakeClock{saturday},
		Coupons: map[string]int{"SAVE20": 20},
	}
	price, _ := p.Final(1000, "SAVE20")
	fmt.Println(price)
	// Output: 700
}

// ---------- Benchmark: go test -bench=. -benchmem ./09-testing/showcase ----------

func BenchmarkPricerFinal(b *testing.B) {
	p := &Pricer{Clock: fakeClock{saturday}, Coupons: map[string]int{"SAVE20": 20}}
	for b.Loop() {
		_, _ = p.Final(1000, "SAVE20")
	}
}
