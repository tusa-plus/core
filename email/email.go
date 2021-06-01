package email

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tusa-plus/core/utils"
)

var ErrInvalidCode = fmt.Errorf("invalid code")
var ErrInvalidFields = fmt.Errorf("invalid fields error")

type CodeSender interface {
	SendCode(to string, code string) error
}

type EmailVerification struct {
	symbols        string
	codeLen        int
	storage        fiber.Storage
	codeExpiration time.Duration
	rndgen         *utils.RandomGenerator
	sender         CodeSender
}

func (ev *EmailVerification) SendCode(to string) error {
	ev.rndgen = utils.NewRandomGenerator(ev.symbols)
	code := ev.rndgen.NextString(ev.codeLen)
	err := ev.storage.Set(to, []byte(code), ev.codeExpiration)
	if err != nil {
		return ErrInvalidFields
	}
	err = ev.sender.SendCode(to, code)
	return err
}

func (ev *EmailVerification) VerifyCode(code string, email string) error {
	data, err := ev.storage.Get(email)
	if data == nil && err == nil {
		return ErrInvalidCode
	}
	if err != nil {
		return ErrInvalidFields
	}
	if string(data) != code {
		return ErrInvalidCode
	}
	return nil
}
