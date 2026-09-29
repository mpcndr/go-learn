package interfaces

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"testing"
)

func almostEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

// asShape ตรวจตอน runtime ว่า v เป็น Shape หรือยัง
// (ใช้แทนการเขียน var _ Shape = Circle{} เพื่อให้ test compile ได้แม้ยังไม่ทำโจทย์)
func asShape(t *testing.T, v any) Shape {
	t.Helper()
	s, ok := v.(Shape)
	if !ok {
		t.Fatalf("%T does not implement Shape yet", v)
	}
	return s
}

func TestShapes(t *testing.T) {
	c := asShape(t, Circle{Radius: 1})
	r := asShape(t, Rectangle{Width: 2, Height: 3})
	tr := asShape(t, Triangle{Base: 4, Height: 5})

	if !almostEqual(c.Area(), math.Pi) || c.Name() != "circle" {
		t.Errorf("Circle: area=%v name=%q", c.Area(), c.Name())
	}
	if r.Area() != 6 || r.Name() != "rectangle" {
		t.Errorf("Rectangle: area=%v name=%q", r.Area(), r.Name())
	}
	if tr.Area() != 10 || tr.Name() != "triangle" {
		t.Errorf("Triangle: area=%v name=%q", tr.Area(), tr.Name())
	}

	shapes := []Shape{c, r, tr}
	if got, want := TotalArea(shapes), math.Pi+16; !almostEqual(got, want) {
		t.Errorf("TotalArea = %v, want %v", got, want)
	}
	if got := Largest(shapes); got == nil || got.Name() != "triangle" {
		t.Errorf("Largest = %v, want triangle", got)
	}
	if got := Largest(nil); got != nil {
		t.Errorf("Largest(nil) = %v, want nil", got)
	}
}

func TestDescribe(t *testing.T) {
	var circle any = Circle{Radius: 1}
	if _, ok := circle.(Shape); !ok {
		t.Skip("Circle is not a Shape yet — finish TestShapes first")
	}
	tests := []struct {
		in   any
		want string
	}{
		{42, "int: 42"},
		{"go", "string of length 2"},
		{circle, "shape: circle"},
		{errors.New("boom"), "error: boom"},
		{3.14, "unknown"},
		{nil, "unknown"},
	}
	for _, tt := range tests {
		if got := Describe(tt.in); got != tt.want {
			t.Errorf("Describe(%#v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCountingWriter(t *testing.T) {
	var sb strings.Builder
	var w any = &CountingWriter{W: &sb}
	iw, ok := w.(io.Writer)
	if !ok {
		t.Fatal("*CountingWriter does not implement io.Writer yet")
	}
	fmt.Fprint(iw, "hello ")
	fmt.Fprint(iw, "สวัสดี")
	cw := w.(*CountingWriter)
	if sb.String() != "hello สวัสดี" {
		t.Errorf("underlying writer got %q", sb.String())
	}
	if cw.Count != 6+18 {
		t.Errorf("Count = %d, want 24", cw.Count)
	}
}

// fakeSender — ตัวอย่างการทำ test double: ไม่ต้องมีระบบส่งจริง
type fakeSender struct {
	sent []string
	fail map[string]bool
}

func (f *fakeSender) Send(to, msg string) error {
	if f.fail[to] {
		return fmt.Errorf("cannot reach %s", to)
	}
	f.sent = append(f.sent, to+":"+msg)
	return nil
}

func TestNotifyAll(t *testing.T) {
	f := &fakeSender{fail: map[string]bool{"bob": true, "dan": true}}
	n, err := NotifyAll(f, []string{"ann", "bob", "cat", "dan"}, "hi")
	if n != 2 {
		t.Errorf("sent = %d, want 2", n)
	}
	if err == nil || !strings.Contains(err.Error(), "bob") {
		t.Errorf("err = %v, want first error (about bob)", err)
	}
	if strings.Join(f.sent, ",") != "ann:hi,cat:hi" {
		t.Errorf("fake received %v", f.sent)
	}
}
