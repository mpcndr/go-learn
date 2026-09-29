// บท 03: Collections — รันด้วย: go run ./03-collections
package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"unicode/utf8"
)

func main() {
	fmt.Println("=== slice คือหน้าต่างที่มองไปยัง array ===")
	nums := []int{10, 20, 30, 40, 50}
	part := nums[1:3]
	fmt.Printf("nums=%v part=%v len=%d cap=%d\n", nums, part, len(part), cap(part))
	part[0] = 99
	fmt.Println("หลังแก้ part[0]=99 → nums =", nums)

	fmt.Println("\n=== append อาจเขียนทับ array เดิม! ===")
	part = append(part, 777) // cap ของ part ยังพอ → เขียนลง nums[3]
	fmt.Println("หลัง append(part, 777) → nums =", nums)

	fmt.Println("\n=== ดู cap โตขึ้นเมื่อ append ===")
	var s []int
	prevCap := -1
	for i := range 20 {
		s = append(s, i)
		if cap(s) != prevCap {
			fmt.Printf("len=%2d cap=%2d\n", len(s), cap(s))
			prevCap = cap(s)
		}
	}

	fmt.Println("\n=== clone เพื่อตัดความสัมพันธ์ ===")
	clone := slices.Clone(nums)
	clone[0] = -1
	fmt.Println("nums =", nums, " clone =", clone)

	fmt.Println("\n=== map + comma-ok ===")
	stock := map[string]int{"apple": 5, "banana": 0}
	for _, fruit := range []string{"apple", "banana", "cherry"} {
		qty, ok := stock[fruit]
		fmt.Printf("%-7s qty=%d มีใน map=%v\n", fruit, qty, ok)
	}
	// range map ลำดับสุ่ม → sort key ก่อนถ้าต้องการลำดับแน่นอน
	for _, k := range slices.Sorted(maps.Keys(stock)) {
		fmt.Print(k, " ")
	}
	fmt.Println()

	fmt.Println("\n=== string / byte / rune ===")
	greet := "สวัสดี Go"
	fmt.Println("len (bytes) =", len(greet), " runes =", utf8.RuneCountInString(greet))
	for i, r := range greet {
		fmt.Printf("%d:%c ", i, r) // index กระโดดทีละ 3 สำหรับภาษาไทย
	}
	fmt.Println()

	var b strings.Builder
	for i := range 5 {
		fmt.Fprintf(&b, "%d,", i)
	}
	fmt.Println("builder:", strings.TrimSuffix(b.String(), ","))
}
