// Package interfaces — แบบฝึกหัดบท 06
package interfaces

import "io"

// Shape รูปทรงอะไรก็ได้ที่หาพื้นที่และชื่อได้
type Shape interface {
	Area() float64
	Name() string
}

type Circle struct{ Radius float64 }
type Rectangle struct{ Width, Height float64 }
type Triangle struct{ Base, Height float64 }

// CountingWriter ห่อ io.Writer อีกตัว และนับจำนวน byte ที่เขียนสำเร็จทั้งหมด
type CountingWriter struct {
	W     io.Writer
	Count int
}

// Sender คือ interface ที่ประกาศ "ฝั่งผู้ใช้" — NotifyAll ต้องการแค่ความสามารถนี้
type Sender interface {
	Send(to, msg string) error
}
