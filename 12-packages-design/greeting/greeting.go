// Package greeting สาธิต API ที่ export เฉพาะสิ่งที่จำเป็น
package greeting

import "github.com/mpcndr/go-learn/12-packages-design/greeting/internal/format"

// Greeter ทักทายด้วยคำนำหน้าที่กำหนด — field ทั้งหมดเป็น private
type Greeter struct {
	prefix string
	count  int
}

// New สร้าง Greeter (export — คนนอกเรียกได้)
func New(prefix string) *Greeter {
	return &Greeter{prefix: prefix}
}

// Greet คืนข้อความทักทาย ถ้าไม่มีชื่อใช้ "คุณ" แทน
func (g *Greeter) Greet(name string) string {
	g.count++
	return format.Join(g.prefix, defaultName(name))
}

// Count จำนวนครั้งที่ทักทาย — ให้อ่านได้ แต่ไม่ให้แก้
func (g *Greeter) Count() int { return g.count }

// defaultName ไม่ export — รายละเอียดภายใน เปลี่ยนได้โดยไม่กระทบใคร
func defaultName(name string) string {
	if name == "" {
		return "คุณ"
	}
	return name
}
