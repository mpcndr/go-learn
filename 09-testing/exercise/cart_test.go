package testingex

import "testing"

// แบบฝึกหัด: เขียน test ให้ Cart (ดูกฎใน cart.go)
//
// ขั้นตอนแนะนำ:
//  1. ลบ t.Skip ออก
//  2. เติม case ลงใน tests ให้ครอบคลุมกฎทุกข้อ — โดยเฉพาะ "ขอบ" ของกฎ (9, 10, 11 ชิ้น)
//  3. รัน go test -v -run TestCartTotal ./09-testing/exercise
//  4. เมื่อ test เจอบั๊ก → แก้ใน cart.go → รันจนผ่าน
//  5. เขียน TestCartRemove เพิ่มเองด้วย (ทดสอบ error ทั้งสองแบบ)
func TestCartTotal(t *testing.T) {
	t.Skip("TODO: ลบบรรทัดนี้แล้วเขียน test ของคุณ")

	type item struct {
		sku   string
		price int
		qty   int
	}
	tests := []struct {
		name  string
		items []item
		want  int
	}{
		{"empty cart", nil, 0},
		{"single item", []item{{"pen", 1000, 2}}, 2000},
		// TODO: สินค้าเดิม Add ซ้ำ
		// TODO: ส่วนลดที่ 10 ชิ้นพอดี
		// TODO: ...
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Cart
			for _, it := range tt.items {
				if err := c.Add(it.sku, it.price, it.qty); err != nil {
					t.Fatalf("Add(%q, %d, %d): %v", it.sku, it.price, it.qty, err)
				}
			}
			if got := c.Total(); got != tt.want {
				t.Errorf("Total() = %d, want %d", got, tt.want)
			}
		})
	}
}
