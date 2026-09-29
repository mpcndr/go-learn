//go:build solution

package errs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseAge(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, &ValidationError{Field: "age", Reason: "not a number"}
	}
	if n < 0 || n > 150 {
		return 0, &ValidationError{Field: "age", Reason: "out of range"}
	}
	return n, nil
}

func FindUser(id int) (User, error) {
	u, ok := users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func LoadProfile(id int) (string, error) {
	u, err := FindUser(id)
	if err != nil {
		return "", fmt.Errorf("load profile %d: %w", id, err)
	}
	return fmt.Sprintf("%s <%s>", u.Name, u.Email), nil
}

func ValidateAll(u User) error {
	var errs []error
	if u.Name == "" {
		errs = append(errs, &ValidationError{Field: "name", Reason: "required"})
	}
	if !strings.Contains(u.Email, "@") {
		errs = append(errs, &ValidationError{Field: "email", Reason: "must contain @"})
	}
	return errors.Join(errs...)
}

func SafeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("safe divide: %v", r)
		}
	}()
	return a / b, nil
}
