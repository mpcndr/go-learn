package generics

import (
	"reflect"
	"slices"
	"strconv"
	"testing"
)

func TestMapFilterReduce(t *testing.T) {
	strs := Map([]int{1, 2, 3}, strconv.Itoa)
	if !slices.Equal(strs, []string{"1", "2", "3"}) {
		t.Errorf("Map = %v", strs)
	}
	evens := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
	if !slices.Equal(evens, []int{2, 4, 6}) {
		t.Errorf("Filter = %v", evens)
	}
	sum := Reduce([]int{1, 2, 3, 4}, 0, func(acc, v int) int { return acc + v })
	if sum != 10 {
		t.Errorf("Reduce sum = %d, want 10", sum)
	}
	joined := Reduce([]string{"a", "b", "c"}, ">", func(acc string, v string) string { return acc + v })
	if joined != ">abc" {
		t.Errorf("Reduce join = %q, want >abc", joined)
	}
}

type Baht float64

func TestSum(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Errorf("Sum(int) = %d", got)
	}
	if got := Sum([]float64{0.5, 0.25}); got != 0.75 {
		t.Errorf("Sum(float64) = %v", got)
	}
	if got := Sum([]Baht{10, 20.5}); got != 30.5 {
		t.Errorf("Sum(Baht) = %v", got)
	}
}

func TestGroupBy(t *testing.T) {
	got := GroupBy([]string{"apple", "banana", "avocado", "blueberry", "cherry"},
		func(s string) byte { return s[0] })
	want := map[byte][]string{
		'a': {"apple", "avocado"},
		'b': {"banana", "blueberry"},
		'c': {"cherry"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GroupBy = %v, want %v", got, want)
	}
}

func TestSet(t *testing.T) {
	a := NewSet(1, 2, 3, 3)
	if a == nil {
		t.Fatal("NewSet returned nil")
	}
	if a.Len() != 3 {
		t.Errorf("Len = %d, want 3", a.Len())
	}
	a.Add(4)
	if !a.Has(4) || a.Has(99) {
		t.Error("Has is wrong")
	}
	b := NewSet(3, 4, 5)
	in := a.Intersect(b)
	if in == nil || in.Len() != 2 || !in.Has(3) || !in.Has(4) {
		t.Errorf("Intersect wrong: %+v", in)
	}
	words := NewSet("go", "go", "rust")
	if words.Len() != 2 {
		t.Errorf("string set Len = %d, want 2", words.Len())
	}
}

func TestCache(t *testing.T) {
	c := NewCache[string, int](2)
	if c == nil {
		t.Fatal("NewCache returned nil")
	}
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("a", 10) // update — ไม่เปลี่ยนลำดับ
	if v, ok := c.Get("a"); !ok || v != 10 {
		t.Errorf("Get(a) = %d, %v, want 10, true", v, ok)
	}
	c.Put("c", 3) // เต็ม → ลบ a (เก่าสุด)
	if _, ok := c.Get("a"); ok {
		t.Error("a should have been evicted")
	}
	for k, want := range map[string]int{"b": 2, "c": 3} {
		if v, ok := c.Get(k); !ok || v != want {
			t.Errorf("Get(%s) = %d, %v, want %d, true", k, v, ok, want)
		}
	}
	if v, ok := c.Get("zzz"); ok || v != 0 {
		t.Errorf("Get(missing) = %d, %v, want 0, false", v, ok)
	}
}
