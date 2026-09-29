// Package domain คือหัวใจของระบบห้องสมุด: ข้อมูล + กฎพื้นฐาน
//
// package นี้ไม่ import อะไรในโปรเจกต์เลย และไม่รู้ว่าข้อมูลถูกเก็บที่ไหน
// หรือถูกเรียกจาก CLI หรือ HTTP — อ่านไฟล์นี้ไฟล์เดียวต้องเข้าใจว่าระบบทำอะไร
package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ---------- กฎของระบบ (จากโจทย์ตรง ๆ) ----------

const (
	MaxLoansPerMember       = 3
	LoanPeriod              = 14 * 24 * time.Hour
	FinePerDay        Money = 500 // 5 บาท (หน่วยสตางค์)
)

// ---------- ความล้มเหลวที่เป็นไปได้ (ส่วนหนึ่งของ API) ----------

var (
	ErrBookNotFound     = errors.New("book not found")
	ErrMemberNotFound   = errors.New("member not found")
	ErrLoanNotFound     = errors.New("loan not found")
	ErrBookUnavailable  = errors.New("book is already on loan")
	ErrLoanLimitReached = errors.New("member has reached the loan limit")
	ErrAlreadyExists    = errors.New("already exists")
)

// ValidationError ข้อมูลเข้าไม่ถูกต้อง พร้อมบอกว่า field ไหน
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string { return e.Field + ": " + e.Reason }

// ---------- ข้อมูล (คำนามในโจทย์) ----------

// ใช้ type แยกสำหรับ ID แต่ละชนิด: compiler จะไม่ยอมให้ส่ง MemberID ไปแทน BookID
type (
	BookID   string
	MemberID string
)

type Book struct {
	ID     BookID
	Title  string
	Author string
}

// NewBook สร้าง Book ที่ถูกต้องเสมอ — ทุกที่ที่มี Book ในระบบ จึงไม่ต้อง validate ซ้ำ
func NewBook(id BookID, title, author string) (Book, error) {
	b := Book{ID: id, Title: strings.TrimSpace(title), Author: strings.TrimSpace(author)}
	switch {
	case strings.TrimSpace(string(id)) == "":
		return Book{}, &ValidationError{"id", "required"}
	case b.Title == "":
		return Book{}, &ValidationError{"title", "required"}
	}
	return b, nil
}

type Member struct {
	ID   MemberID
	Name string
}

func NewMember(id MemberID, name string) (Member, error) {
	m := Member{ID: id, Name: strings.TrimSpace(name)}
	switch {
	case strings.TrimSpace(string(id)) == "":
		return Member{}, &ValidationError{"id", "required"}
	case m.Name == "":
		return Member{}, &ValidationError{"name", "required"}
	}
	return m, nil
}

// Loan การยืมหนังสือ 1 เล่ม โดยสมาชิก 1 คน
type Loan struct {
	BookID     BookID
	MemberID   MemberID
	BorrowedAt time.Time
	DueAt      time.Time
}

// NewLoan สร้างการยืมที่กำหนดวันคืนตามกฎ LoanPeriod
func NewLoan(book BookID, member MemberID, now time.Time) Loan {
	return Loan{BookID: book, MemberID: member, BorrowedAt: now, DueAt: now.Add(LoanPeriod)}
}

// FineAt คิดค่าปรับถ้าคืน ณ เวลา returnedAt
// เกินกำหนดแม้ไม่ถึงวัน ก็นับเป็น 1 วัน
//
// สังเกต: เป็น pure function — รับเวลาเข้ามาแทนการเรียก time.Now() เอง จึงทดสอบง่ายมาก
func (l Loan) FineAt(returnedAt time.Time) Money {
	late := returnedAt.Sub(l.DueAt)
	if late <= 0 {
		return 0
	}
	days := (late + 24*time.Hour - 1) / (24 * time.Hour) // ปัดขึ้น
	return FinePerDay * Money(days)
}

// ---------- Value object ----------

// Money หน่วยสตางค์ — ใช้ int เสมอสำหรับเงิน (float มีปัญหาความคลาดเคลื่อน)
type Money int64

func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign, m = "-", -m
	}
	return fmt.Sprintf("%s฿%d.%02d", sign, m/100, m%100)
}
