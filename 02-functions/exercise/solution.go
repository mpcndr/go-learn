//go:build solution

package functions

func MinMax(nums []int) (min, max int, ok bool) {
	if len(nums) == 0 {
		return 0, 0, false
	}
	min, max = nums[0], nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return min, max, true
}

func MakeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func Apply(nums []int, fn func(int) int) []int {
	out := make([]int, len(nums))
	for i, n := range nums {
		out[i] = fn(n)
	}
	return out
}

func Compose(f, g func(int) int) func(int) int {
	return func(x int) int { return g(f(x)) }
}

func Memoize(fn func(int) int) func(int) int {
	cache := make(map[int]int)
	return func(x int) int {
		if v, ok := cache[x]; ok {
			return v
		}
		v := fn(x)
		cache[x] = v
		return v
	}
}
