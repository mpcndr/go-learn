//go:build solution

package pointers

func Swap(a, b *int) {
	*a, *b = *b, *a
}

// slice header ถูก copy แต่ ptr ในนั้นชี้ array เดียวกับของผู้เรียก จึงแก้ element ได้
func Increment(nums []int) {
	for i := range nums {
		nums[i]++
	}
}

func AppendTo(s *[]int, values ...int) {
	*s = append(*s, values...)
}

func ResetScores(m map[string]Student) {
	for k, st := range m {
		st.Score = 0
		m[k] = st
	}
}

func Ptr[T any](v T) *T {
	return &v
}
