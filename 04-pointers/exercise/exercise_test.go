package pointers

import (
	"slices"
	"testing"
)

func TestSwap(t *testing.T) {
	a, b := 1, 2
	Swap(&a, &b)
	if a != 2 || b != 1 {
		t.Errorf("after Swap: a=%d b=%d, want a=2 b=1", a, b)
	}
}

func TestIncrement(t *testing.T) {
	nums := []int{1, 2, 3}
	Increment(nums)
	if !slices.Equal(nums, []int{2, 3, 4}) {
		t.Errorf("Increment = %v, want [2 3 4]", nums)
	}
}

func TestAppendTo(t *testing.T) {
	var s []int
	AppendTo(&s, 1)
	AppendTo(&s, 2, 3)
	if !slices.Equal(s, []int{1, 2, 3}) {
		t.Errorf("AppendTo = %v, want [1 2 3]", s)
	}
}

func TestResetScores(t *testing.T) {
	m := map[string]Student{
		"a": {"Ann", 90},
		"b": {"Ben", 75},
	}
	ResetScores(m)
	for k, st := range m {
		if st.Score != 0 {
			t.Errorf("m[%q].Score = %d, want 0", k, st.Score)
		}
	}
	if m["a"].Name != "Ann" {
		t.Errorf("Name changed: %q", m["a"].Name)
	}
}

func TestPtr(t *testing.T) {
	p := Ptr(42)
	if p == nil || *p != 42 {
		t.Fatalf("Ptr(42) = %v, want pointer to 42", p)
	}
	x := 7
	q := Ptr(x)
	*q = 100
	if x != 7 {
		t.Error("Ptr must point to a copy, not the original variable")
	}
	if s := Ptr("go"); s == nil || *s != "go" {
		t.Error("Ptr should work with any type")
	}
}
