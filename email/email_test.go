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
	userEmail := "xxxxxxx@gmail.com"
	return userEmail, ev, sender
}

func Test_Generate(t *testing.T) {
	userEmail, ev, sender := createEVWithSender()
	err := ev.SendCode(userEmail)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.channel)
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
}

func Test_Verify(t *testing.T) {
	userEmail, ev, sender := createEVWithSender()
	err := ev.SendCode(userEmail)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.channel)
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
	err = ev.VerifyCode(code, userEmail)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_VerifySmallTimeout(t *testing.T) {
	userEmail, ev, sender := createEVWithSender()
	err := ev.SendCode(userEmail)
	now := time.Now()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.channel)
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
	time.Sleep(time.Until(now.Add(smallTime)))
	err = ev.VerifyCode(code, userEmail)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_VerifyBigTimeout(t *testing.T) {
	userEmail, ev, sender := createEVWithSender()
	err := ev.SendCode(userEmail)
	now := time.Now()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.channel)
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
	time.Sleep(time.Until(now.Add(ev.codeExpiration)))
	err = ev.VerifyCode(code, userEmail)
	if err == nil {
		t.Fatalf("Code wasn't delete")
	}
}

func Test_VerifyWithoutSendCode(t *testing.T) {
	userEmail, ev, _ := createEVWithSender()
	code := "Some invalid code"
	err := ev.VerifyCode(code, userEmail)
	if err == nil {
		t.Fatalf("Code wasn't delete")
	}
}
