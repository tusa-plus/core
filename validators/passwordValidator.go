package validators

import (
	"fmt"
	"unicode"
)

const maxPasswordLen = 32
const minPasswordLen = 8

var ErrLongPasswordLen = fmt.Errorf("length of password must be at most %v symbols", maxPasswordLen)
var ErrShortPasswordLen = fmt.Errorf("length of password must be at lest %v symbols", minPasswordLen)
var ErrNotDigits = fmt.Errorf("password should have at least one numeral")
var ErrNotLowercaseLatters = fmt.Errorf("password should have at least one lowercase letter")
var ErrNotUppercaseLatters = fmt.Errorf("password should have at least one uppercase letter")

func ValidatePassword(pwd string) error {
	if len(pwd) < minPasswordLen {
		return ErrShortPasswordLen
	}
	if len(pwd) > maxPasswordLen {
		return ErrLongPasswordLen
	}
	var hasUpper = false
	var hasLower = false
	var hasDigit = false
	for _, symb := range pwd {
		if unicode.IsDigit(symb) {
			hasDigit = true
		}
		if unicode.IsLower(symb) {
			hasLower = true
		}
		if unicode.IsUpper(symb) {
			hasUpper = true
		}
	}
	if !hasDigit {
		return ErrNotDigits
	}
	if !hasLower {
		return ErrNotLowercaseLatters
	}
	if !hasUpper {
		return ErrNotUppercaseLatters
	}
	return nil
}
