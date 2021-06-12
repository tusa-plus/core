package email

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory"
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

const (
	validSymbolsMock = "0123456789"
	codeLenMock      = 6
)

func NewMockEmailVerification() (string, EmailVerification, *ChannelCodeSender) {
	codeSender := ChannelCodeSender{
		Channel: make(chan string, codeLenMock),
	}
	ev := EmailVerification{
		symbols:        validSymbolsMock,
		codeLen:        codeLenMock,
		storage:        memory.New(),
		codeExpiration: time.Second * 2,
		rndgen:         utils.NewRandomGenerator(validSymbolsMock),
		sender:         &codeSender,
	}
	userEmail := "xxxxxxx@gmail.com"
	return userEmail, ev, &codeSender
}

func NewEmailVerification(codeLen int, codeExpiration time.Duration,
	sender, login, password, host, port, validCodeSymbols string) EmailVerification {
	codeSender := NewEmailCodeSender(sender, login, password, host, port)
	ev := EmailVerification{
		symbols:        validCodeSymbols,
		codeLen:        codeLen,
		storage:        memory.New(),
		codeExpiration: codeExpiration,
		rndgen:         utils.NewRandomGenerator(validCodeSymbols),
		sender:         &codeSender,
	}
	return ev
}
