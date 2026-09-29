package testingex

import "errors"

var (
	ErrInvalidQty   = errors.New("quantity must be positive")
	ErrItemNotFound = errors.New("item not in cart")
)

// Cart ตะกร้าสินค้า — ราคาเป็นหน่วยสตางค์
//
// กฎ:
//  1. Add สินค้าเดิมซ้ำ → รวมจำนวน
//  2. Remove ลดจำนวน ถ้าเหลือ 0 ให้เอาออกจากตะกร้า, ลบเกินที่มี → ErrInvalidQty
//  3. Total = ผลรวม ราคา × จำนวน
//  4. ซื้อสินค้าชิ้นเดียวกันตั้งแต่ 10 ชิ้นขึ้นไป ลด 5% เฉพาะรายการนั้น (ปัดเศษทิ้ง)
//
// ⚠️ โค้ดนี้มีบั๊ก 1 จุด — เขียน test ใน cart_test.go ให้ครอบคลุมกฎทุกข้อ แล้วคุณจะเจอเอง
type Cart struct {
	items map[string]line
}

type line struct {
	price int
	qty   int
}

func (c *Cart) Add(sku string, price, qty int) error {
	if qty <= 0 {
		return ErrInvalidQty
	}
	if c.items == nil {
		c.items = make(map[string]line)
	}
	l := c.items[sku]
	l.price = price
	l.qty += qty
	c.items[sku] = l
	return nil
}

func (c *Cart) Remove(sku string, qty int) error {
	l, ok := c.items[sku]
	if !ok {
		return ErrItemNotFound
	}
	if qty <= 0 || qty > l.qty {
		return ErrInvalidQty
	}
	l.qty -= qty
	if l.qty == 0 {
		delete(c.items, sku)
		return nil
	}
	c.items[sku] = l
	return nil
}

func (c *Cart) Count() int {
	n := 0
	for _, l := range c.items {
		n += l.qty
	}
	return n
}

func (c *Cart) Total() int {
	total := 0
	for _, l := range c.items {
		sub := l.price * l.qty
		if l.qty > 10 {
			sub = sub * 95 / 100
		}
		total += sub
	}
	return total
}
