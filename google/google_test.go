package google

import (
	"context"
	"errors"
	"github.com/tusa-plus/core/utils"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"testing"
)

func Test_NewGoogleEmptyLogger(t *testing.T) {
	t.Parallel()
	pool := utils.NewHTTPClientPool()
	_, err := NewGoogle(nil, &pool)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

func Test_NewGoogleEmptyPool(t *testing.T) {
	t.Parallel()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	_, err = NewGoogle(logger, nil)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

func Test_GoogleGetEmail(t *testing.T) {
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
	google, err := NewGoogle(logger, &pool)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	email, err := google.GetEmail(context.Background(), cfg.Section("google").Key("token").String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedEmail := cfg.Section("google").Key("email").String()
	if email != expectedEmail {
		t.Fatalf("invalid email: expected %v, got %v", expectedEmail, email)
	}
}

func Test_GoogleGetEmailInvalidToken(t *testing.T) {
	t.Parallel()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	pool := utils.NewHTTPClientPool()
	google, err := NewGoogle(logger, &pool)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	invalidTokens := []string{"", "abcefre", "32424++_!>?|~`"}
	for index := range invalidTokens {
		if _, err = google.GetEmail(context.Background(), invalidTokens[index]); !errors.Is(err, ErrValidate) {
			t.Fatalf("expected validation error: %v", err)
		}
	}
}
