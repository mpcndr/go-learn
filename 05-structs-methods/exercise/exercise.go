//go:build !solution

package structs

// Area คืนพื้นที่ (value receiver — ไม่แก้ข้อมูล)
func (r Rect) Area() float64 {
	// TODO
	return 0
}

// Scale ขยายขนาดด้วย factor (ต้องแก้ข้อมูลจริง → receiver ควรเป็นแบบไหน?)
func (r *Rect) Scale(factor float64) {
	// TODO
}

// NewAccount สร้างบัญชีใหม่ — owner ว่างไม่ได้ (คืน ErrEmptyOwner)
func NewAccount(owner string) (*Account, error) {
	// TODO
	return nil, nil
}

// Deposit ฝากเงิน — amount <= 0 คืน ErrInvalidAmount
func (a *Account) Deposit(amount int) error {
	// TODO
	return nil
}

// Withdraw ถอนเงิน — amount <= 0 คืน ErrInvalidAmount, เงินไม่พอคืน ErrInsufficientFunds
// ถ้า error ต้องไม่เปลี่ยน balance
func (a *Account) Withdraw(amount int) error {
	// TODO
	return nil
}

// Balance คืนยอดเงินปัจจุบัน
func (a *Account) Balance() int {
	// TODO
	return 0
}

// String แสดงเงินเป็นบาท ทศนิยม 2 ตำแหน่ง เช่น Money(1250) → "฿12.50", Money(-5) → "-฿0.05"
// เมื่อ type มี method String() string, fmt.Println จะเรียกใช้อัตโนมัติ (interface fmt.Stringer — บท 06)
func (m Money) String() string {
	// TODO: fmt.Sprintf("%d.%02d", ...)
	return ""
}

// Push ใส่ค่าลงบนสุด
func (s *Stack) Push(v int) {
	// TODO
}

// Pop เอาค่าบนสุดออก ถ้าว่างคืน (0, false)
func (s *Stack) Pop() (int, bool) {
	// TODO
	return 0, false
}

// Len จำนวนสมาชิก
func (s *Stack) Len() int {
	// TODO
	return 0
}
