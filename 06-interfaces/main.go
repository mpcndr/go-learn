// บท 06: Interface — รันด้วย: go run ./06-interfaces
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// ---------- 1. implicit implementation ----------

type Speaker interface {
	Speak() string
}

type Dog struct{ Name string }
type Robot struct{ Model int }

func (d Dog) Speak() string   { return d.Name + ": โฮ่ง!" }
func (r Robot) Speak() string { return fmt.Sprintf("RX-%d: บี๊บ", r.Model) }

// ไม่มีคำว่า implements — Dog และ Robot เป็น Speaker อัตโนมัติ
func chorus(ss []Speaker) {
	for _, s := range ss {
		fmt.Println(" ", s.Speak())
	}
}

// ---------- 2. interface ประกาศฝั่งผู้ใช้ + ทดแทนกันได้ ----------

// ReportPrinter ต้องการแค่ io.Writer — จะเขียนลงจอ ไฟล์ หรือ buffer ก็ได้
func printReport(w io.Writer, items map[string]int) {
	for name, qty := range items {
		fmt.Fprintf(w, "%-8s %3d\n", name, qty)
	}
}

// ---------- 3. decorator: ห่อ interface ด้วย interface เดียวกัน ----------

type upperWriter struct{ w io.Writer }

func (u upperWriter) Write(p []byte) (int, error) {
	return u.w.Write([]byte(strings.ToUpper(string(p))))
}

// ---------- 4. กับดัก nil interface ----------

type MyErr struct{}

func (*MyErr) Error() string { return "my error" }

func mayFailBad(fail bool) error {
	var e *MyErr // nil pointer
	if fail {
		e = &MyErr{}
	}
	return e // ❌ ถึง e จะเป็น nil แต่ interface ที่คืนไปมี type = *MyErr → ไม่เท่ากับ nil
}

func mayFailGood(fail bool) error {
	if fail {
		return &MyErr{}
	}
	return nil // ✅
}

func main() {
	fmt.Println("=== implicit interface ===")
	chorus([]Speaker{Dog{"ด่าง"}, Robot{42}})

	fmt.Println("\n=== io.Writer ตัวเดียว ใช้ได้หลายปลายทาง ===")
	stock := map[string]int{"apple": 5}
	printReport(os.Stdout, stock)

	var buf strings.Builder // *strings.Builder ก็เป็น io.Writer
	printReport(&buf, stock)
	fmt.Printf("เก็บใน buffer: %q\n", buf.String())

	printReport(upperWriter{os.Stdout}, stock) // decorator

	fmt.Println("\n=== type switch ===")
	for _, v := range []any{42, "go", 3.14, Dog{"ขาว"}, nil} {
		switch x := v.(type) {
		case int:
			fmt.Println("  int, x*2 =", x*2)
		case string:
			fmt.Println("  string, len =", len(x))
		case Speaker:
			fmt.Println("  Speaker →", x.Speak())
		case nil:
			fmt.Println("  nil")
		default:
			fmt.Printf("  อื่น ๆ: %T\n", x)
		}
	}

	fmt.Println("\n=== กับดัก nil interface ===")
	fmt.Println("mayFailBad(false) == nil ?", mayFailBad(false) == nil, " ← ผิดคาด!")
	fmt.Println("mayFailGood(false) == nil ?", mayFailGood(false) == nil)

	var target *MyErr
	fmt.Println("errors.As ใช้ได้:", errors.As(mayFailGood(true), &target))
}
