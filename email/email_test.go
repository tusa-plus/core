package email

import (
	"testing"
	"time"

	"github.com/gofiber/storage/memory"
	"github.com/tusa-plus/core/common"
)

const (
	validSymbols = "0123456789"
)

var channelSender = ChannelCodeSender{
	channel: make(chan string),
}

var ev = EmailVerification{
	symbols:        validSymbols,
	storage:        memory.New(),
	codeExpiration: time.Second,
	rndgen:         *common.NewRandomGenerator(validSymbols),
	sender:         channelSender.(*CodeSender),
}

func Test_Generate(t *testing.T) {
	err := ev.SendCode(ConfigSmtpSender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_GenerateCheck(t *testing.T) {
	err := ev.SendCode(ConfigSmtpSender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s, ok := ev.(sender); !ok {
		t.Fatalf("unexpected error: %v", err)
	}
}
