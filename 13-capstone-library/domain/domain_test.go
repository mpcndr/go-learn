package domain

import (
	"errors"
	"testing"
	"time"
)

func TestLoanFineAt(t *testing.T) {
	borrowed := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	loan := NewLoan("B1", "M1", borrowed)
	due := borrowed.Add(LoanPeriod)

	if !loan.DueAt.Equal(due) {
		t.Fatalf("DueAt = %v, want %v", loan.DueAt, due)
	}

	tests := []struct {
		name     string
		returned time.Time
		want     Money
	}{
		{"early", borrowed.Add(24 * time.Hour), 0},
		{"exactly on due", due, 0},
		{"1 minute late counts as 1 day", due.Add(time.Minute), 500},
		{"exactly 1 day late", due.Add(24 * time.Hour), 500},
		{"1 day and 1 hour late", due.Add(25 * time.Hour), 1000},
		{"10 days late", due.Add(10 * 24 * time.Hour), 5000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loan.FineAt(tt.returned); got != tt.want {
				t.Errorf("FineAt(%v) = %v, want %v", tt.returned.Sub(due), got, tt.want)
			}
		})
	}
}

func TestNewBookValidation(t *testing.T) {
	tests := []struct {
		id, title string
		field     string
	}{
		{"", "Go", "id"},
		{"  ", "Go", "id"},
		{"B1", "", "title"},
		{"B1", "   ", "title"},
	}
	for _, tt := range tests {
		_, err := NewBook(BookID(tt.id), tt.title, "someone")
		var ve *ValidationError
		if !errors.As(err, &ve) || ve.Field != tt.field {
			t.Errorf("NewBook(%q, %q) err = %v, want ValidationError on %q", tt.id, tt.title, err, tt.field)
		}
	}
	b, err := NewBook("B1", "  The Go Programming Language ", "Donovan")
	if err != nil || b.Title != "The Go Programming Language" {
		t.Errorf("NewBook valid = %+v, %v", b, err)
	}
}

func TestMoneyString(t *testing.T) {
	for m, want := range map[Money]string{0: "฿0.00", 500: "฿5.00", 12345: "฿123.45", -50: "-฿0.50"} {
		if got := m.String(); got != want {
			t.Errorf("Money(%d).String() = %q, want %q", int64(m), got, want)
		}
	}
}
