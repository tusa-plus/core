package email

import (
	"github.com/gofiber/storage/memory"
	"github.com/tusa-plus/core/utils"
	"testing"
	"time"
)

const (
	validSymbols = "0123456789"
)

var sender = &ChannelCodeSender{
	channel: make(chan string),
}

var ev = EmailVerification{
	symbols:        validSymbols,
	storage:        memory.New(),
	codeExpiration: time.Second,
	rndgen:         utils.NewRandomGenerator(validSymbols),
	sender:         sender,
}

func Test_Generate(t *testing.T) {
	err := ev.SendCode(ConfigSmtpSender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
