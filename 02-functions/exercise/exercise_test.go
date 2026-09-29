package functions

import (
	"slices"
	"testing"
)

func TestMinMax(t *testing.T) {
	min, max, ok := MinMax([]int{3, -2, 9, 4})
	if !ok || min != -2 || max != 9 {
		t.Errorf("MinMax = (%d, %d, %v), want (-2, 9, true)", min, max, ok)
	}
	if _, _, ok := MinMax(nil); ok {
		t.Error("MinMax(nil) ok = true, want false")
	}
}

func TestMakeCounter(t *testing.T) {
	a, b := MakeCounter(), MakeCounter()
	if a == nil || b == nil {
		t.Fatal("MakeCounter returned nil")
	}
	a()
	a()
	if got := a(); got != 3 {
		t.Errorf("third call of a = %d, want 3", got)
	}
	if got := b(); got != 1 {
		t.Errorf("first call of b = %d, want 1 (counters must be independent)", got)
	}
}

func TestApply(t *testing.T) {
	in := []int{1, 2, 3}
	got := Apply(in, func(x int) int { return x * x })
	if !slices.Equal(got, []int{1, 4, 9}) {
		t.Errorf("Apply = %v, want [1 4 9]", got)
	}
	if !slices.Equal(in, []int{1, 2, 3}) {
		t.Errorf("Apply modified input: %v", in)
	}
}

func TestCompose(t *testing.T) {
	addOne := func(x int) int { return x + 1 }
	double := func(x int) int { return x * 2 }
	h := Compose(addOne, double)
	if h == nil {
		t.Fatal("Compose returned nil")
	}
	if got := h(5); got != 12 {
		t.Errorf("Compose(addOne, double)(5) = %d, want 12", got)
	}
}

func TestMemoize(t *testing.T) {
	calls := 0
	slowSquare := func(x int) int { calls++; return x * x }
	fast := Memoize(slowSquare)
	if fast == nil {
		t.Fatal("Memoize returned nil")
	}
	for range 3 {
		if got := fast(4); got != 16 {
			t.Fatalf("fast(4) = %d, want 16", got)
		}
	}
	fast(5)
	if calls != 2 {
		t.Errorf("underlying fn called %d times, want 2", calls)
	}
}
