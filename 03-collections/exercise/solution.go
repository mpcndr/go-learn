//go:build solution

package collections

import (
	"cmp"
	"maps"
	"slices"
	"strings"
)

func WordCount(text string) map[string]int {
	counts := make(map[string]int)
	for _, w := range strings.Fields(text) {
		counts[strings.ToLower(w)]++
	}
	return counts
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func Unique(nums []int) []int {
	seen := make(map[int]struct{}, len(nums))
	out := make([]int, 0, len(nums))
	for _, n := range nums {
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	return out
}

func Chunk(nums []int, size int) [][]int {
	if size <= 0 {
		return nil
	}
	var out [][]int
	for size < len(nums) {
		// full slice expression (low:high:max) จำกัด cap ไม่ให้ append ของผู้ใช้ไปทับก้อนถัดไป
		nums, out = nums[size:], append(out, nums[:size:size])
	}
	if len(nums) > 0 {
		out = append(out, nums)
	}
	return out
}

func TopN(text string, n int) []string {
	counts := WordCount(text)
	words := slices.Collect(maps.Keys(counts))
	slices.SortFunc(words, func(a, b string) int {
		if c := cmp.Compare(counts[b], counts[a]); c != 0 {
			return c // มากไปน้อย
		}
		return cmp.Compare(a, b)
	})
	if n < len(words) {
		words = words[:n]
	}
	return words
}
