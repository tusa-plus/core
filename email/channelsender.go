package email

import (
	"time"

	"github.com/gofiber/storage/memory"
	"github.com/tusa-plus/core/utils"
)

type ChannelCodeSender struct {
	Channel chan string
}

func (ccs *ChannelCodeSender) SendCode(to string, code string) error {
	ccs.Channel <- code
	return nil
}

const (
	validSymbols = "0123456789"
	codeLen      = 6
)

func NewMockEmailVerification() (string, EmailVerification, *ChannelCodeSender) {
	sender := &ChannelCodeSender{
		Channel: make(chan string, codeLen),
	}
	ev := EmailVerification{
		symbols:        validSymbols,
		codeLen:        codeLen,
		storage:        memory.New(),
		codeExpiration: time.Second * 2,
		rndgen:         utils.NewRandomGenerator(validSymbols),
		sender:         sender,
	}
	userEmail := "xxxxxxx@gmail.com"
	return userEmail, ev, sender
}
