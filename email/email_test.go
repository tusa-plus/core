package email

import (
	"os"
	"testing"
	"time"

	"github.com/gofiber/storage/memory"
	"github.com/tusa-plus/core/utils"
)

const (
	validSymbols = "0123456789"
	codeLen      = 6
	smallTime    = time.Millisecond * 1000
)

func createEVWithSender() (string, EmailVerification, *ChannelCodeSender) {
	sender := &ChannelCodeSender{
		channel: make(chan string, codeLen),
	}
	ev := EmailVerification{
		symbols:        validSymbols,
		codeLen:        codeLen,
		storage:        memory.New(),
		codeExpiration: time.Second * 2,
		rndgen:         utils.NewRandomGenerator(validSymbols),
		sender:         sender,
	}
	ConfigSMTPSender := os.Getenv("EMAIL_SENDER")
	return ConfigSMTPSender, ev, sender
}

func Test_Generate(t *testing.T) {
	ConfigSMTPSender, ev, sender := createEVWithSender()
	err := ev.SendCode(ConfigSMTPSender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.channel)
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
}

func Test_Verify(t *testing.T) {
	ConfigSMTPSender, ev, sender := createEVWithSender()
	err := ev.SendCode(ConfigSMTPSender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.channel)
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
	err = ev.VerifyCode(code, ConfigSMTPSender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_VerifySmallTimeout(t *testing.T) {
	ConfigSMTPSender, ev, sender := createEVWithSender()
	err := ev.SendCode(ConfigSMTPSender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.channel)
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
	time.Sleep(smallTime)
	err = ev.VerifyCode(code, ConfigSMTPSender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_VerifyBigTimeout(t *testing.T) {
	ConfigSMTPSender, ev, sender := createEVWithSender()
	err := ev.SendCode(ConfigSMTPSender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.channel)
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
	time.Sleep(ev.codeExpiration)
	err = ev.VerifyCode(code, ConfigSMTPSender)
	if err == nil {
		t.Fatalf("Code wasn't delete")
	}
}

func Test_VerifyWithoutSendCode(t *testing.T) {
	ConfigSMTPSender, ev, _ := createEVWithSender()
	code := "Some invalid code"
	err := ev.VerifyCode(code, ConfigSMTPSender)
	if err == nil {
		t.Fatalf("Code wasn't delete")
	}
}
