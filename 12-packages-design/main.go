// บท 12: Package — รันด้วย: go run ./12-packages-design
package main

import (
	"fmt"

	"github.com/mpcndr/go-learn/12-packages-design/greeting"
	// "github.com/mpcndr/go-learn/12-packages-design/greeting/internal/format"
	// ↑ ลองเปิดคอมเมนต์ดู: "use of internal package ... not allowed"
)

func main() {
	g := greeting.New("สวัสดี")
	fmt.Println(g.Greet("gopher"))
	fmt.Println(g.Greet(""))

	// g.prefix = "hi" ← compile error: prefix เป็น field ที่ไม่ได้ export
	fmt.Println("จำนวนครั้งที่ทักทาย:", g.Count())
}
