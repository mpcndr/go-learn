// Package format เป็น internal: import ได้เฉพาะจาก package ที่อยู่ใต้ greeting/
package format

import "strings"

// Join ต่อคำด้วยช่องว่าง แล้วใส่ "!" ท้าย
func Join(parts ...string) string {
	return strings.Join(parts, " ") + "!"
}
