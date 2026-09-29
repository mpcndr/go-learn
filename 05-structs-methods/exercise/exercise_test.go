package structs

import (
	"fmt"
	"testing"
)

func TestRect(t *testing.T) {
	r := Rect{Width: 3, Height: 4}
	if got := r.Area(); got != 12 {
		t.Errorf("Area = %v, want 12", got)
	}
	r.Scale(2)
	if r.Width != 6 || r.Height != 8 {
		t.Errorf("after Scale(2): %+v, want {Width:6 Height:8}", r)
	}
}

func TestAccount(t *testing.T) {
	if _, err := NewAccount(""); err != ErrEmptyOwner {
		t.Errorf("NewAccount(\"\") err = %v, want ErrEmptyOwner", err)
	}
	a, err := NewAccount("Ann")
	if err != nil || a == nil {
		t.Fatalf("NewAccount(\"Ann\") = %v, %v", a, err)
	}
	if a.Owner != "Ann" {
		t.Errorf("Owner = %q, want Ann", a.Owner)
	}
	if err := a.Deposit(0); err != ErrInvalidAmount {
		t.Errorf("Deposit(0) err = %v, want ErrInvalidAmount", err)
	}
	if err := a.Deposit(100); err != nil {
		t.Fatalf("Deposit(100) err = %v", err)
	}
	if err := a.Withdraw(-1); err != ErrInvalidAmount {
		t.Errorf("Withdraw(-1) err = %v, want ErrInvalidAmount", err)
	}
	if err := a.Withdraw(150); err != ErrInsufficientFunds {
		t.Errorf("Withdraw(150) err = %v, want ErrInsufficientFunds", err)
	}
	if got := a.Balance(); got != 100 {
		t.Errorf("balance after failed withdraw = %d, want 100 (unchanged)", got)
	}
	if err := a.Withdraw(40); err != nil {
		t.Fatalf("Withdraw(40) err = %v", err)
	}
	if got := a.Balance(); got != 60 {
		t.Errorf("Balance = %d, want 60", got)
	}
}

func TestMoneyString(t *testing.T) {
	tests := map[Money]string{
		0:      "฿0.00",
		5:      "฿0.05",
		1250:   "฿12.50",
		100000: "฿1000.00",
		-5:     "-฿0.05",
		-1999:  "-฿19.99",
	}
	for m, want := range tests {
		if got := fmt.Sprint(m); got != want {
			t.Errorf("Money(%d) prints %q, want %q", int64(m), got, want)
		}
	}
}

func TestStack(t *testing.T) {
	var s Stack // zero value ต้องใช้ได้
	if _, ok := s.Pop(); ok {
		t.Error("Pop on empty stack should return ok=false")
	}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if s.Len() != 3 {
		t.Errorf("Len = %d, want 3", s.Len())
	}
	for _, want := range []int{3, 2, 1} {
		if got, ok := s.Pop(); !ok || got != want {
			t.Errorf("Pop = (%d, %v), want (%d, true)", got, ok, want)
		}
	}
	if s.Len() != 0 {
		t.Errorf("Len = %d, want 0", s.Len())
	}
}
