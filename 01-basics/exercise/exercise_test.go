package basics

import (
	"slices"
	"testing"
)

func TestFizzBuzz(t *testing.T) {
	got := FizzBuzz(15)
	want := []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"}
	if !slices.Equal(got, want) {
		t.Errorf("FizzBuzz(15)\n got: %v\nwant: %v", got, want)
	}
	if got := FizzBuzz(0); got == nil || len(got) != 0 {
		t.Errorf("FizzBuzz(0) = %#v, want empty non-nil slice []string{}", got)
	}
}

func TestGrade(t *testing.T) {
	tests := []struct {
		score int
		want  string
	}{
		{100, "A"}, {80, "A"}, {79, "B"}, {70, "B"}, {65, "C"},
		{50, "D"}, {49, "F"}, {0, "F"}, {-1, "invalid"}, {101, "invalid"},
	}
	for _, tt := range tests {
		if got := Grade(tt.score); got != tt.want {
			t.Errorf("Grade(%d) = %q, want %q", tt.score, got, tt.want)
		}
	}
}

func TestSumDigits(t *testing.T) {
	tests := map[int]int{0: 0, 7: 7, 1234: 10, 9999: 36, -45: 9}
	for in, want := range tests {
		if got := SumDigits(in); got != want {
			t.Errorf("SumDigits(%d) = %d, want %d", in, got, want)
		}
	}
}

func TestIsPrime(t *testing.T) {
	primes := []int{2, 3, 5, 7, 11, 13, 97, 7919}
	notPrimes := []int{-7, 0, 1, 4, 9, 15, 100, 7917}
	for _, n := range primes {
		if !IsPrime(n) {
			t.Errorf("IsPrime(%d) = false, want true", n)
		}
	}
	for _, n := range notPrimes {
		if IsPrime(n) {
			t.Errorf("IsPrime(%d) = true, want false", n)
		}
	}
}
