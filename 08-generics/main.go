// บท 08: Generics — รันด้วย: go run ./08-generics
package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// ก่อนมี generics ต้องเขียนซ้ำแบบนี้:
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// มี generics — เขียนครั้งเดียว
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// custom constraint + ~
type Number interface {
	~int | ~int64 | ~float64
}

type Celsius float64 // underlying type = float64 → เข้ากับ ~float64

func Average[T Number](nums []T) T {
	if len(nums) == 0 {
		return 0
	}
	var total T
	for _, n := range nums {
		total += n
	}
	return total / T(len(nums))
}

// generic type
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func (p Pair[K, V]) String() string { return fmt.Sprintf("%v=%v", p.Key, p.Value) }

// Queue แบบ generic
type Queue[T any] struct{ items []T }

func (q *Queue[T]) Enqueue(v T) { q.items = append(q.items, v) }
func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T // zero value ของ T ใด ๆ
	if len(q.items) == 0 {
		return zero, false
	}
	v := q.items[0]
	q.items = q.items[1:]
	return v, true
}

func main() {
	fmt.Println("maxInt(3, 7)      =", maxInt(3, 7))
	fmt.Println("Max(3, 7)         =", Max(3, 7))
	fmt.Println("Max(2.5, 1.5)     =", Max(2.5, 1.5))
	fmt.Println(`Max("go", "rust") =`, Max("go", "rust"))

	temps := []Celsius{30, 32, 35}
	fmt.Printf("Average(Celsius) = %.2f (type %T)\n", Average(temps), Average(temps))

	pairs := []Pair[string, int]{{"b", 2}, {"a", 1}, {"c", 3}}
	slices.SortFunc(pairs, func(x, y Pair[string, int]) int { return cmp.Compare(x.Key, y.Key) })
	fmt.Println("sorted pairs:", pairs)

	var q Queue[string]
	q.Enqueue("first")
	q.Enqueue("second")
	v, _ := q.Dequeue()
	fmt.Println("dequeue:", v)

	// standard library ที่เป็น generic
	words := []string{"Go", "is", "simple", "go", "IS"}
	fmt.Println("slices.Contains:", slices.Contains(words, "simple"))
	fmt.Println("slices.MaxFunc by len:", slices.MaxFunc(words, func(a, b string) int { return cmp.Compare(len(a), len(b)) }))
	lower := make([]string, len(words))
	for i, w := range words {
		lower[i] = strings.ToLower(w)
	}
	slices.Sort(lower)
	fmt.Println("unique (Compact after Sort):", slices.Compact(lower))
}
