// บท 04: Pointer และ Value Semantics — รันด้วย: go run ./04-pointers
package main

import "fmt"

type Player struct {
	Name  string
	Score int
}

func bonusByValue(p Player)    { p.Score += 10 } // แก้ copy
func bonusByPointer(p *Player) { p.Score += 10 } // แก้ต้นฉบับ (Go ให้เขียน p.Score แทน (*p).Score)

func appendOne(s []int)     { s = append(s, 1) }
func appendOnePtr(s *[]int) { *s = append(*s, 1) }
func setFirst(s []int)      { s[0] = 999 }

func newPlayer(name string) *Player {
	p := Player{Name: name} // ตัวแปร local...
	return &p               // ...คืนที่อยู่ได้ปลอดภัย (escape ไป heap)
}

func main() {
	fmt.Println("=== & และ * ===")
	x := 10
	p := &x
	fmt.Printf("x=%d  p=%p  *p=%d\n", x, p, *p)
	*p = 20
	fmt.Println("หลัง *p = 20 → x =", x)

	fmt.Println("\n=== struct: value vs pointer ===")
	alice := Player{Name: "Alice", Score: 50}
	bonusByValue(alice)
	fmt.Println("หลัง bonusByValue:  ", alice.Score) // 50
	bonusByPointer(&alice)
	fmt.Println("หลัง bonusByPointer:", alice.Score) // 60

	fmt.Println("\n=== assign ก็คือ copy ===")
	bob := alice
	bob.Name = "Bob"
	fmt.Println("alice:", alice.Name, " bob:", bob.Name)
	ref := &alice
	ref.Name = "Alicia"
	fmt.Println("หลังแก้ผ่าน ref → alice:", alice.Name)

	fmt.Println("\n=== slice: header ถูก copy แต่ array แชร์กัน ===")
	s := make([]int, 1, 10)
	setFirst(s)
	fmt.Println("หลัง setFirst:", s) // เห็นการเปลี่ยน
	appendOne(s)
	fmt.Println("หลัง appendOne:", s, "len =", len(s)) // ไม่เห็น (len ใน header ของเรายังเป็น 1)
	appendOnePtr(&s)
	fmt.Println("หลัง appendOnePtr:", s, "len =", len(s))

	fmt.Println("\n=== range copy ค่า! ===")
	team := []Player{{"A", 1}, {"B", 2}}
	for _, pl := range team {
		pl.Score = 100 // แก้ copy — ไม่มีผล
	}
	fmt.Println("range แบบ value:", team)
	for i := range team {
		team[i].Score = 100 // แก้ผ่าน index — มีผล
	}
	fmt.Println("range แบบ index:", team)

	fmt.Println("\n=== คืน pointer ของตัวแปร local ===")
	np := newPlayer("Gopher")
	fmt.Printf("%+v\n", *np)

	fmt.Println("\n=== nil pointer ===")
	var nobody *Player
	if nobody == nil {
		fmt.Println("เช็ค nil ก่อนใช้เสมอ ไม่งั้น panic")
	}
}
