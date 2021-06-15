package validators

import (
	"fmt"
	"regexp"
)

const regexr = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@" +
	"[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
const maxEmailLen = 245

var ErrLongEmailLen = fmt.Errorf("length of email is too long")
var ErrInvalidEmailFormat = fmt.Errorf("invalid email format")

func ValidateEmail(email string) error {
	var rxEmail = regexp.MustCompile(regexr)
	if len(email) > maxEmailLen {
		return ErrLongEmailLen
	}
	if !rxEmail.MatchString(email) {
		return ErrInvalidEmailFormat
	}
	return nil
}
