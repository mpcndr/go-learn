// บท 05: Struct และ Method — รันด้วย: go run ./05-structs-methods
package main

import (
	"errors"
	"fmt"
	"time"
)

// --- type ใหม่จาก type พื้นฐาน: ให้ compiler ช่วยแยกความหมาย ---
type Celsius float64

func (c Celsius) ToFahrenheit() float64 { return float64(c)*9/5 + 32 }

// --- struct ขนาดเล็ก ไม่เปลี่ยนแปลง → value receiver ---
type Point struct{ X, Y float64 }

func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }

// --- struct ที่มี invariant → field private + pointer receiver ---
var ErrInsufficient = errors.New("ยอดเงินไม่พอ")

type Wallet struct {
	owner   string
	balance int
}

func NewWallet(owner string) *Wallet { return &Wallet{owner: owner} }

func (w *Wallet) Deposit(n int) { w.balance += n }

func (w *Wallet) Pay(n int) error {
	if n > w.balance {
		return ErrInsufficient
	}
	w.balance -= n
	return nil
}

func (w *Wallet) Balance() int { return w.balance }

// --- embedding ---
type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Timestamps) Touch() { t.UpdatedAt = time.Now() }

type Post struct {
	Timestamps // embedded
	Title      string
}

func main() {
	fmt.Println("=== type ใหม่ ===")
	temp := Celsius(37)
	fmt.Printf("%.1f°C = %.1f°F\n", temp, temp.ToFahrenheit())

	fmt.Println("\n=== value receiver ===")
	a := Point{1, 2}
	b := a.Add(Point{10, 20})
	fmt.Println("a =", a, " b =", b, "(a ไม่เปลี่ยน)")

	fmt.Println("\n=== invariant ถูกปกป้องด้วย method ===")
	w := NewWallet("Ann")
	w.Deposit(100)
	if err := w.Pay(30); err != nil {
		fmt.Println("error:", err)
	}
	if err := w.Pay(500); err != nil {
		fmt.Println("จ่าย 500 →", err)
	}
	fmt.Println("คงเหลือ:", w.Balance())
	// w.balance = -999  ← ถ้าอยู่คนละ package จะทำไม่ได้ (compile error)

	fmt.Println("\n=== embedding ===")
	p := Post{Title: "Hello Go", Timestamps: Timestamps{CreatedAt: time.Now()}}
	time.Sleep(10 * time.Millisecond)
	p.Touch() // promoted method
	fmt.Printf("%q แก้ไขล่าสุดหลังสร้าง %v\n", p.Title, p.UpdatedAt.Sub(p.CreatedAt).Round(time.Millisecond))

	fmt.Println("\n=== รูปแบบการพิมพ์ struct: v, +v, #v ===")
	fmt.Printf("%v\n%+v\n%#v\n", a, a, a)
}
