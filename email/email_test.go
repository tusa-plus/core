package email

import (
	"testing"
	"time"

	"github.com/gofiber/storage/memory"
	"github.com/tusa-plus/core/utils"
)

const (
	validSymbols = "0123456789"
	codeLen      = 6
	smallTime    = time.Millisecond * 1000
	bigTime      = time.Millisecond * 2200
)

var sender = &ChannelCodeSender{
	channel: make(chan string, codeLen),
}

var ev = EmailVerification{
	symbols:        validSymbols,
	storage:        memory.New(),
	codeExpiration: time.Second * 2,
	rndgen:         utils.NewRandomGenerator(validSymbols),
	sender:         sender,
}

func Test_Generate(t *testing.T) {
	err := ev.SendCode(ConfigSmtpSender, validSymbols, codeLen)
	code := <-(sender.channel)
	if err != nil || len(code) != codeLen {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
}

func Test_Verify(t *testing.T) {
	err := ev.SendCode(ConfigSmtpSender, validSymbols, codeLen)
	code := <-(sender.channel)
	if err != nil || len(code) != codeLen {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
	err = ev.VerifyCode(code, ConfigSmtpSender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_VerifySmallTimeout(t *testing.T) {
	err := ev.SendCode(ConfigSmtpSender, validSymbols, codeLen)
	code := <-(sender.channel)
	if err != nil || len(code) != codeLen {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
	time.Sleep(smallTime)
	err = ev.VerifyCode(code, ConfigSmtpSender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_VerifyBigTimeout(t *testing.T) {
	err := ev.SendCode(ConfigSmtpSender, validSymbols, codeLen)
	code := <-(sender.channel)
	if err != nil || len(code) != codeLen {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
	time.Sleep(bigTime)
	err = ev.VerifyCode(code, ConfigSmtpSender)
	if err == nil {
		t.Fatalf("Code wasn't delete")
	}
}
