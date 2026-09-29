// บท 02: ฟังก์ชัน — รันด้วย: go run ./02-functions
package main

import (
	"errors"
	"fmt"
	"strings"
)

// signature บอกทุกอย่าง: รับ 2 ตัว คืนผลลัพธ์ หรือ error
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("หารด้วยศูนย์ไม่ได้")
	}
	return a / b, nil
}

// named results — ใช้เมื่อชื่อช่วยให้อ่านเข้าใจ
func divmod(a, b int) (quotient, remainder int) {
	return a / b, a % b
}

// variadic
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// ฟังก์ชันรับฟังก์ชัน
func transform(words []string, fn func(string) string) []string {
	out := make([]string, 0, len(words))
	for _, w := range words {
		out = append(out, fn(w))
	}
	return out
}

// ฟังก์ชันคืนฟังก์ชัน (closure)
func multiplier(factor int) func(int) int {
	return func(x int) int { return x * factor }
}

func deferDemo() {
	fmt.Println("  เริ่ม")
	for i := 1; i <= 3; i++ {
		defer fmt.Println("  defer", i) // LIFO: 3, 2, 1
	}
	fmt.Println("  จบเนื้อฟังก์ชัน")
}

func main() {
	if result, err := divide(10, 4); err == nil {
		fmt.Println("10/4 =", result)
	}
	if _, err := divide(1, 0); err != nil {
		fmt.Println("error:", err)
	}

	q, r := divmod(17, 5)
	fmt.Printf("17 = 5*%d + %d\n", q, r)

	fmt.Println("sum() =", sum(), " sum(1,2,3) =", sum(1, 2, 3))
	nums := []int{10, 20, 30}
	fmt.Println("sum(nums...) =", sum(nums...))

	fmt.Println(transform([]string{"go", "is", "fun"}, strings.ToUpper))
	fmt.Println(transform([]string{"go", "is", "fun"}, func(s string) string {
		return s + "!"
	}))

	double, triple := multiplier(2), multiplier(3)
	fmt.Println("double(5) =", double(5), " triple(5) =", triple(5))

	fmt.Println("defer demo:")
	deferDemo()
}
