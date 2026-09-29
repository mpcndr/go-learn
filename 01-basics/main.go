// บท 01: พื้นฐาน — รันด้วย: go run ./01-basics
package main

import "fmt"

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// String ทำให้ fmt.Println แสดงชื่อวันแทนตัวเลข (จะเข้าใจลึกขึ้นในบท 06: interfaces)
func (d Weekday) String() string {
	names := [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	if d < Sunday || d > Saturday {
		return "Unknown"
	}
	return names[d]
}

func main() {
	// --- zero values ---
	var (
		i int
		f float64
		s string
		b bool
		p *int
	)
	fmt.Printf("zero values: int=%d float=%v string=%q bool=%v pointer=%v\n", i, f, s, b, p)

	// --- ประกาศแบบสั้น + type inference ---
	name := "Gopher"
	age := 15
	height := 1.75
	fmt.Printf("%s อายุ %d สูง %.2f m (types: %T %T %T)\n", name, age, height, name, age, height)

	// --- ต้องแปลง type เอง ---
	total, count := 7, 2
	fmt.Println("7/2 แบบ int =", total/count)
	fmt.Println("7/2 แบบ float =", float64(total)/float64(count))

	// --- if พร้อมประกาศตัวแปร ---
	if n := len(name); n > 5 {
		fmt.Println(name, "ชื่อยาว", n, "ตัวอักษร")
	}

	// --- for ทุกรูปแบบ ---
	for i := 0; i < 3; i++ {
		fmt.Print(i, " ")
	}
	fmt.Println()

	n := 1
	for n < 100 { // แบบ while
		n *= 2
	}
	fmt.Println("เลขยกกำลัง 2 ตัวแรกที่ >= 100 คือ", n)

	for idx, ch := range "Go!" {
		fmt.Printf("[%d]=%c ", idx, ch)
	}
	fmt.Println()

	// --- switch ไม่ต้อง break ---
	for _, score := range []int{95, 72, 40} {
		switch {
		case score >= 80:
			fmt.Println(score, "→ ดีมาก")
		case score >= 50:
			fmt.Println(score, "→ ผ่าน")
		default:
			fmt.Println(score, "→ ไม่ผ่าน")
		}
	}

	// --- iota + method String() ---
	today := Wednesday
	fmt.Println("วันนี้คือ", today, "(ค่าจริง =", int(today), ")")
}
