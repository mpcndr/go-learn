//go:build solution

package interfaces

import (
	"fmt"
	"math"
)

func (c Circle) Area() float64 { return math.Pi * c.Radius * c.Radius }
func (c Circle) Name() string  { return "circle" }

func (r Rectangle) Area() float64 { return r.Width * r.Height }
func (r Rectangle) Name() string  { return "rectangle" }

func (t Triangle) Area() float64 { return 0.5 * t.Base * t.Height }
func (t Triangle) Name() string  { return "triangle" }

func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

func Largest(shapes []Shape) Shape {
	var best Shape
	for _, s := range shapes {
		if best == nil || s.Area() > best.Area() {
			best = s
		}
	}
	return best
}

func Describe(v any) string {
	switch x := v.(type) {
	case int:
		return fmt.Sprintf("int: %d", x)
	case string:
		return fmt.Sprintf("string of length %d", len(x))
	case Shape:
		return "shape: " + x.Name()
	case error:
		return "error: " + x.Error()
	default:
		return "unknown"
	}
}

func (c *CountingWriter) Write(p []byte) (int, error) {
	n, err := c.W.Write(p)
	c.Count += n
	return n, err
}

func NotifyAll(s Sender, users []string, msg string) (sent int, err error) {
	for _, u := range users {
		if e := s.Send(u, msg); e != nil {
			if err == nil {
				err = e
			}
			continue
		}
		sent++
	}
	return sent, err
}
