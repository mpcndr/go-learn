//go:build solution

package structs

import "fmt"

func (r Rect) Area() float64 { return r.Width * r.Height }

func (r *Rect) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

func NewAccount(owner string) (*Account, error) {
	if owner == "" {
		return nil, ErrEmptyOwner
	}
	return &Account{Owner: owner}, nil
}

func (a *Account) Deposit(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.balance += amount
	return nil
}

func (a *Account) Withdraw(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if amount > a.balance {
		return ErrInsufficientFunds
	}
	a.balance -= amount
	return nil
}

func (a *Account) Balance() int { return a.balance }

func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign, m = "-", -m
	}
	return fmt.Sprintf("%s฿%d.%02d", sign, m/100, m%100)
}

func (s *Stack) Push(v int) { s.items = append(s.items, v) }

func (s *Stack) Pop() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	last := len(s.items) - 1
	v := s.items[last]
	s.items = s.items[:last]
	return v, true
}

func (s *Stack) Len() int { return len(s.items) }
