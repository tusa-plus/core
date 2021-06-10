package email

import (
	"testing"
	"time"
)

const (
	smallTime = time.Millisecond * 1000
)

func Test_Generate(t *testing.T) {
	userEmail, ev, sender := NewMockEmailVerification()
	err := ev.SendCode(userEmail)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.Channel)
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
}

func Test_Verify(t *testing.T) {
	userEmail, ev, sender := NewMockEmailVerification()
	err := ev.SendCode(userEmail)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.Channel)
	if len(code) != codeLen {
		t.Fatalf("invalid code len")
	}
	err = ev.VerifyCode(code, userEmail)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_VerifySmallTimeout(t *testing.T) {
	userEmail, ev, sender := NewMockEmailVerification()
	err := ev.SendCode(userEmail)
	now := time.Now()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.Channel)
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
	userEmail, ev, sender := NewMockEmailVerification()
	err := ev.SendCode(userEmail)
	now := time.Now()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code := <-(sender.Channel)
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
	userEmail, ev, _ := NewMockEmailVerification()
	code := "Some invalid code"
	err := ev.VerifyCode(code, userEmail)
	if err == nil {
		t.Fatalf("Code wasn't delete")
	}
}
