// Package showcase แสดงเทคนิคการเขียน test ทั้งหมดในบท 09
// อ่าน price.go ก่อน แล้วค่อยอ่าน price_test.go
package showcase

import (
	"errors"
	"time"
)

var ErrNegativePrice = errors.New("price must not be negative")

// Clock — interface เล็ก ๆ ที่ทำให้ "เวลา" ทดสอบได้
// ถ้าเรียก time.Now() ตรง ๆ ข้างใน จะ test เงื่อนไข "วันหยุด" ไม่ได้เลย
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// Pricer คิดราคาสุดท้าย: ลด 10% ในวันเสาร์-อาทิตย์ และลดเพิ่มตาม coupon
type Pricer struct {
	Clock   Clock
	Coupons map[string]int // code → เปอร์เซ็นต์ส่วนลด
}

func NewPricer(coupons map[string]int) *Pricer {
	return &Pricer{Clock: realClock{}, Coupons: coupons}
}

// Final คืนราคาสุดท้าย (หน่วยสตางค์) — coupon ที่ไม่รู้จักจะถูกเพิกเฉย
func (p *Pricer) Final(price int, coupon string) (int, error) {
	if price < 0 {
		return 0, ErrNegativePrice
	}
	discount := 0
	if wd := p.Clock.Now().Weekday(); wd == time.Saturday || wd == time.Sunday {
		discount += 10
	}
	discount += p.Coupons[coupon]
	if discount > 100 {
		discount = 100
	}
	return price * (100 - discount) / 100, nil
}
