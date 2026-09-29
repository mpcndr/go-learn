// Package errs — แบบฝึกหัดบท 07
package errs

import (
	"errors"
	"fmt"
)

// ErrNotFound — sentinel error สำหรับหาไม่พบ
var ErrNotFound = errors.New("not found")

// ValidationError บอกว่า field ไหนผิดและเพราะอะไร
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

type User struct {
	ID    int
	Name  string
	Email string
}

// users — "ฐานข้อมูล" จำลอง
var users = map[int]User{
	1: {ID: 1, Name: "Ann", Email: "ann@example.com"},
	2: {ID: 2, Name: "Ben", Email: "ben@example.com"},
}
