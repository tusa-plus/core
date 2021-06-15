package validators

import (
	"testing"
)

func Test_ValidateEmailOk(t *testing.T) {
	t.Parallel()
	email := "example@mail.ru"
	err := ValidateEmail(email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_ValidateEmailFail(t *testing.T) {
	t.Parallel()
	email := "examplefail@mail.ru"
	longEmail := make([]byte, 0, maxEmailLen+1)
	for _, symb := range email {
		longEmail = append(longEmail, byte(symb))
	}
	for len(longEmail) <= maxEmailLen {
		longEmail = append(longEmail, 'u')
	}
	emails := []string{"aa", "examplemail.ru", "example@", "@", "@example", string(longEmail)}
	for _, email := range emails {
		err := ValidateEmail(email)
		if err == nil {
			t.Fatalf("Email passes validation, but it is incorrect\nemail: %v", email)
		}
	}
}

func Test_ValidatePasswordOk(t *testing.T) {
	t.Parallel()
	password := "Password12345fail"
	err := ValidatePassword(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_ValidatePasswordFail(t *testing.T) {
	t.Parallel()
	password := "Password12345"
	longPassword := make([]byte, 0, maxPasswordLen+1)
	for _, symb := range password {
		longPassword = append(longPassword, byte(symb))
	}
	for len(longPassword) <= maxPasswordLen {
		longPassword = append(longPassword, 'u')
	}
	passwords := []string{"aB1", string(longPassword), "123456789ABC", "12345bca", "AbcAbcAbc", "12345678", "AAAAAAAA", "aaaaaaaaa"}
	for _, password := range passwords {
		err := ValidatePassword(password)
		if err == nil {
			t.Fatalf("Password passes validation, but it is incorrect\npassword: %v", password)
		}
	}
}
