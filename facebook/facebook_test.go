package facebook

import (
	"context"
	"errors"
	"github.com/tusa-plus/core/utils"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"testing"
)

func Test_FacebookGetEmail(t *testing.T) {
	t.Parallel()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	cfg, err := ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pool := utils.NewHTTPClientPool()
	facebook := Facebook{
		httpClientPool: &pool,
		logger:         logger,
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
	t.Parallel()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	pool := utils.NewHTTPClientPool()
	facebook := Facebook{
		httpClientPool: &pool,
		logger:         logger,
	}
	invalidTokens := []string{"", "abcefre", "32424++_!>?|~`"}
	for index := range invalidTokens {
		if _, err = facebook.GetEmail(context.Background(), invalidTokens[index]); !errors.Is(err, ErrValidate) {
			t.Fatalf("expected validation error: %v", err)
		}
	}
}
