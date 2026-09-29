//go:build solution

package testingex

import (
	"strings"
	"unicode"
)

func Slugify(s string) string {
	var b strings.Builder
	pendingDash := false
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r) {
			if pendingDash && b.Len() > 0 {
				b.WriteByte('-')
			}
			pendingDash = false
			b.WriteRune(r)
			continue
		}
		pendingDash = true
	}
	return b.String()
}

// เฉลย Cart (อย่าเพิ่งอ่านถ้ายังไม่ได้เขียน test เอง!):
// บั๊กอยู่ใน Total — เงื่อนไข `l.qty > 10` ต้องเป็น `l.qty >= 10`
// test ที่จับได้คือ case "ซื้อ 10 ชิ้นพอดี" — บั๊กส่วนใหญ่อยู่ที่ "ขอบ" ของเงื่อนไขเสมอ
