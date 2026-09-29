// Package structs — แบบฝึกหัดบท 05
// type ทั้งหมดอยู่ในไฟล์นี้ ให้เขียน method ใน exercise.go
package structs

import "errors"

// Rect สี่เหลี่ยมผืนผ้า
type Rect struct {
	Width, Height float64
}

// Error ที่ Account ใช้ — ประกาศเป็นตัวแปรเพื่อให้ผู้เรียกเปรียบเทียบได้ (ดูบท 07)
var (
	ErrEmptyOwner        = errors.New("owner must not be empty")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// Account บัญชีเงินฝาก — invariant: balance ต้องไม่ติดลบเสมอ
// balance เป็น private: แก้ได้ผ่าน method เท่านั้น
type Account struct {
	Owner   string
	balance int
}

// Money จำนวนเงินหน่วย "สตางค์" (1 บาท = 100 สตางค์)
// เก็บเป็น int เพื่อเลี่ยงปัญหาความคลาดเคลื่อนของ float — อย่าใช้ float กับเงิน!
type Money int64

// Stack โครงสร้างข้อมูลแบบ LIFO — zero value ต้องใช้งานได้ทันที (var s Stack)
type Stack struct {
	items []int
}
