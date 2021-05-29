package facebook

import (
	"context"
	"github.com/tusa-plus/core/utils"
	"gopkg.in/ini.v1"
	"testing"
)

func Test_FacebookGetEmail(t *testing.T) {
	cfg, err := ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pool := utils.NewHttpClientPool()
	facebook := Facebook{
		httpClientPool: &pool,
	}
	email, err := facebook.GetEmail(context.Background(), cfg.Section("fb").Key("token").String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedEmail := cfg.Section("fb").Key("email").String()
	if email != expectedEmail {
		t.Fatalf("invalid email: expected %v, got %v", expectedEmail, email)
	}
}

func Test_FacebookGetEmailInvalidToken(t *testing.T) {
	pool := utils.NewHttpClientPool()
	facebook := Facebook{
		httpClientPool: &pool,
	}
	_, err := facebook.GetEmail(context.Background(), "xxttzz")
	if err != ErrValidate {
		t.Fatalf("expected validation error: %v", err)
	}
}
