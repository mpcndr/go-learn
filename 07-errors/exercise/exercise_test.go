package errs

import (
	"errors"
	"testing"
)

func TestParseAge(t *testing.T) {
	if n, err := ParseAge("30"); err != nil || n != 30 {
		t.Errorf("ParseAge(30) = %d, %v", n, err)
	}
	tests := map[string]string{
		"abc": "not a number",
		"":    "not a number",
		"-1":  "out of range",
		"200": "out of range",
	}
	for in, reason := range tests {
		_, err := ParseAge(in)
		var ve *ValidationError
		if !errors.As(err, &ve) {
			t.Errorf("ParseAge(%q) err = %v, want *ValidationError", in, err)
			continue
		}
		if ve.Field != "age" || ve.Reason != reason {
			t.Errorf("ParseAge(%q) = %+v, want Field=age Reason=%q", in, *ve, reason)
		}
	}
}

func TestFindUser(t *testing.T) {
	if u, err := FindUser(1); err != nil || u.Name != "Ann" {
		t.Errorf("FindUser(1) = %+v, %v", u, err)
	}
	if _, err := FindUser(99); err != ErrNotFound {
		t.Errorf("FindUser(99) err = %v, want ErrNotFound", err)
	}
}

func TestLoadProfile(t *testing.T) {
	if s, err := LoadProfile(2); err != nil || s != "Ben <ben@example.com>" {
		t.Errorf("LoadProfile(2) = %q, %v", s, err)
	}
	_, err := LoadProfile(99)
	if err == nil {
		t.Fatal("LoadProfile(99) err = nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("errors.Is(err, ErrNotFound) = false — did you wrap with %%w?")
	}
	if err.Error() != "load profile 99: not found" {
		t.Errorf("message = %q, want %q", err.Error(), "load profile 99: not found")
	}
}

func TestValidateAll(t *testing.T) {
	if err := ValidateAll(User{Name: "Ann", Email: "a@b.c"}); err != nil {
		t.Errorf("valid user err = %v, want nil", err)
	}
	err := ValidateAll(User{})
	if err == nil {
		t.Fatal("ValidateAll(empty) = nil")
	}
	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("err = %v, want errors combined with errors.Join", err)
	}
	var fields []string
	for _, e := range joined.Unwrap() {
		var ve *ValidationError
		if errors.As(e, &ve) {
			fields = append(fields, ve.Field)
		}
	}
	if len(fields) != 2 || fields[0] != "name" || fields[1] != "email" {
		t.Errorf("fields with errors = %v, want [name email]", fields)
	}
}

func TestSafeDivide(t *testing.T) {
	if r, err := SafeDivide(10, 3); err != nil || r != 3 {
		t.Errorf("SafeDivide(10,3) = %d, %v", r, err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SafeDivide let a panic escape: %v", r)
		}
	}()
	if _, err := SafeDivide(1, 0); err == nil {
		t.Error("SafeDivide(1,0) err = nil, want error")
	}
}
