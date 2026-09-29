package testingex

import (
	"fmt"
	"strings"
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"simple", "Hello World", "hello-world"},
		{"trim and punctuation", "  Hello, World!  ", "hello-world"},
		{"digits and dots", "Go 1.24 Release", "go-1-24-release"},
		{"collapse separators", "a -- b __ c", "a-b-c"},
		{"thai", "เรียน Go", "เรียน-go"},
		{"thai with marks", "สวัสดี ชาวโลก", "สวัสดี-ชาวโลก"},
		{"empty", "", ""},
		{"only symbols", "!!! ???", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.in); got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func ExampleSlugify() {
	fmt.Println(Slugify("  Learn Go, The Right Way!  "))
	// Output: learn-go-the-right-way
}

// FuzzSlugify ตรวจ "คุณสมบัติ" ที่ต้องจริงเสมอ ไม่ว่า input จะเป็นอะไร
// รัน: go test -fuzz=FuzzSlugify -fuzztime=10s ./09-testing/exercise
func FuzzSlugify(f *testing.F) {
	for _, seed := range []string{"Hello World", "  --  ", "เรียน Go", "a.b.c"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, in string) {
		got := Slugify(in)
		if strings.HasPrefix(got, "-") || strings.HasSuffix(got, "-") {
			t.Errorf("Slugify(%q) = %q has leading/trailing dash", in, got)
		}
		if strings.Contains(got, "--") {
			t.Errorf("Slugify(%q) = %q has double dash", in, got)
		}
		if again := Slugify(got); again != got {
			t.Errorf("not idempotent: Slugify(%q) = %q, Slugify(that) = %q", in, got, again)
		}
	})
}

func BenchmarkSlugify(b *testing.B) {
	in := strings.Repeat("Hello, World! ", 20)
	for b.Loop() {
		Slugify(in)
	}
}
