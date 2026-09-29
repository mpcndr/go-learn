package collections

import (
	"maps"
	"reflect"
	"slices"
	"testing"
)

func TestWordCount(t *testing.T) {
	got := WordCount("Go is fun  and go is fast")
	want := map[string]int{"go": 2, "is": 2, "fun": 1, "and": 1, "fast": 1}
	if !maps.Equal(got, want) {
		t.Errorf("WordCount\n got: %v\nwant: %v", got, want)
	}
}

func TestReverse(t *testing.T) {
	tests := map[string]string{
		"":       "",
		"abc":    "cba",
		"Hello!": "!olleH",
		"สวัสดี": "ีดสัวส",
	}
	for in, want := range tests {
		if got := Reverse(in); got != want {
			t.Errorf("Reverse(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestUnique(t *testing.T) {
	got := Unique([]int{3, 1, 3, 2, 1, 3})
	if !slices.Equal(got, []int{3, 1, 2}) {
		t.Errorf("Unique = %v, want [3 1 2]", got)
	}
	if got := Unique([]int{}); len(got) != 0 {
		t.Errorf("Unique(empty) = %v, want empty", got)
	}
}

func TestChunk(t *testing.T) {
	tests := []struct {
		nums []int
		size int
		want [][]int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}},
		{[]int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {3, 4}}},
		{[]int{1, 2}, 5, [][]int{{1, 2}}},
		{[]int{}, 3, nil},
		{[]int{1, 2}, 0, nil},
	}
	for _, tt := range tests {
		got := Chunk(tt.nums, tt.size)
		if len(got) == 0 && len(tt.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Chunk(%v, %d) = %v, want %v", tt.nums, tt.size, got, tt.want)
		}
	}
}

func TestTopN(t *testing.T) {
	text := "the cat and the dog and the bird"
	got := TopN(text, 2)
	if !slices.Equal(got, []string{"the", "and"}) {
		t.Errorf("TopN(2) = %v, want [the and]", got)
	}
	// เท่ากัน → เรียงตามตัวอักษร
	got = TopN(text, 4)
	if !slices.Equal(got, []string{"the", "and", "bird", "cat"}) {
		t.Errorf("TopN(4) = %v, want [the and bird cat]", got)
	}
	if got := TopN(text, 100); len(got) != 5 {
		t.Errorf("TopN(100) len = %d, want 5", len(got))
	}
}
