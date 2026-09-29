// บท 07: Error — รันด้วย: go run ./07-errors
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
)

// --- sentinel ---
var ErrOutOfStock = errors.New("out of stock")

// --- custom type ---
type PaymentError struct {
	Code   int
	Amount int
}

func (e *PaymentError) Error() string {
	return fmt.Sprintf("payment declined (code %d) for %d baht", e.Code, e.Amount)
}

// แต่ละชั้นเติมบริบท "กำลังทำอะไร"
func reserveStock(item string) error {
	if item == "rare-card" {
		return ErrOutOfStock
	}
	return nil
}

func charge(amount int) error {
	if amount > 1000 {
		return &PaymentError{Code: 51, Amount: amount}
	}
	return nil
}

func checkout(orderID int, item string, amount int) error {
	if err := reserveStock(item); err != nil {
		return fmt.Errorf("checkout order %d: reserve %q: %w", orderID, item, err)
	}
	if err := charge(amount); err != nil {
		return fmt.Errorf("checkout order %d: %w", orderID, err)
	}
	return nil
}

func handle(err error) {
	var pe *PaymentError
	switch {
	case err == nil:
		fmt.Println("  ✅ สำเร็จ")
	case errors.Is(err, ErrOutOfStock):
		fmt.Println("  📦 สินค้าหมด → แจ้งผู้ใช้ให้รอ |", err)
	case errors.As(err, &pe):
		fmt.Printf("  💳 จ่ายไม่ผ่าน code=%d → ให้ลองบัตรอื่น | %v\n", pe.Code, err)
	default:
		fmt.Println("  ❌ ไม่รู้จัก:", err)
	}
}

func mustAtoi(s string) (n int, err error) {
	// recover แปลง panic กลับเป็น error — ใช้ที่ขอบนอกเท่านั้น
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	v, e := strconv.Atoi(s)
	if e != nil {
		panic(e) // แค่สาธิต — ปกติให้ return error
	}
	return v, nil
}

func main() {
	fmt.Println("=== wrapping + Is/As ===")
	handle(checkout(1, "book", 300))
	handle(checkout(2, "rare-card", 300))
	handle(checkout(3, "tv", 25000))

	fmt.Println("\n=== error จาก standard library ก็ใช้แบบเดียวกัน ===")
	_, err := os.Open("/no/such/file")
	fmt.Println("  errors.Is(err, fs.ErrNotExist) =", errors.Is(err, fs.ErrNotExist))
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		fmt.Println("  path ที่พัง:", pathErr.Path, "| op:", pathErr.Op)
	}

	fmt.Println("\n=== errors.Join ===")
	joined := errors.Join(errors.New("name is empty"), ErrOutOfStock)
	fmt.Printf("  %q (Join ต่อข้อความด้วย newline)\n  Is(ErrOutOfStock)=%v\n", joined.Error(), errors.Is(joined, ErrOutOfStock))

	fmt.Println("\n=== panic → recover ===")
	fmt.Println(mustAtoi("12"))
	fmt.Println(mustAtoi("abc"))
}
