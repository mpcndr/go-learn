// Package memory เก็บข้อมูลห้องสมุดไว้ใน memory
//
// Store implement service.Books, service.Members และ service.Loans โดยไม่ต้อง import service เลย
// (interface ของ Go เป็นแบบ implicit — บท 06) ถ้าวันหนึ่งจะเปลี่ยนเป็น database
// ก็แค่เขียน package ใหม่ที่มี method ชุดเดียวกัน
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/mpcndr/go-learn/13-capstone-library/domain"
)

// Store ปลอดภัยต่อการใช้จากหลาย goroutine
type Store struct {
	mu      sync.RWMutex // อ่านพร้อมกันได้หลายคน, เขียนได้ทีละคน
	books   map[domain.BookID]domain.Book
	members map[domain.MemberID]domain.Member
	loans   map[domain.BookID]domain.Loan // หนังสือ 1 เล่มมีการยืมที่ active ได้แค่ 1 รายการ
}

func NewStore() *Store {
	return &Store{
		books:   make(map[domain.BookID]domain.Book),
		members: make(map[domain.MemberID]domain.Member),
		loans:   make(map[domain.BookID]domain.Loan),
	}
}

// ---------- Books ----------

func (s *Store) Book(_ context.Context, id domain.BookID) (domain.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	if !ok {
		return domain.Book{}, fmt.Errorf("%w: %s", domain.ErrBookNotFound, id)
	}
	return b, nil
}

func (s *Store) AddBook(_ context.Context, b domain.Book) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.books[b.ID]; ok {
		return fmt.Errorf("book %s: %w", b.ID, domain.ErrAlreadyExists)
	}
	s.books[b.ID] = b
	return nil
}

func (s *Store) ListBooks(_ context.Context) ([]domain.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Book, 0, len(s.books))
	for _, b := range s.books {
		out = append(out, b)
	}
	return out, nil
}

// ---------- Members ----------

func (s *Store) Member(_ context.Context, id domain.MemberID) (domain.Member, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.members[id]
	if !ok {
		return domain.Member{}, fmt.Errorf("%w: %s", domain.ErrMemberNotFound, id)
	}
	return m, nil
}

func (s *Store) AddMember(_ context.Context, m domain.Member) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.members[m.ID]; ok {
		return fmt.Errorf("member %s: %w", m.ID, domain.ErrAlreadyExists)
	}
	s.members[m.ID] = m
	return nil
}

// ---------- Loans ----------

func (s *Store) ActiveLoan(_ context.Context, book domain.BookID) (domain.Loan, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.loans[book]
	if !ok {
		return domain.Loan{}, domain.ErrLoanNotFound
	}
	return l, nil
}

func (s *Store) LoansByMember(_ context.Context, member domain.MemberID) ([]domain.Loan, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []domain.Loan
	for _, l := range s.loans {
		if l.MemberID == member {
			out = append(out, l)
		}
	}
	return out, nil
}

func (s *Store) SaveLoan(_ context.Context, l domain.Loan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.loans[l.BookID] = l
	return nil
}

func (s *Store) DeleteLoan(_ context.Context, book domain.BookID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.loans[book]; !ok {
		return domain.ErrLoanNotFound
	}
	delete(s.loans, book)
	return nil
}
