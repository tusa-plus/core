package vk

import (
	"errors"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"testing"
)

func Test_VkGetEmail(t *testing.T) {
	t.Parallel()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	cfg, err := ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vk := Vk{
		logger:         logger,
	}
	email, err := vk.GetID(cfg.Section("vk").Key("token").String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedEmail, _ := cfg.Section("vk").Key("id").Uint64()
	if email != expectedEmail {
		t.Fatalf("invalid email: expected %v, got %v", expectedEmail, email)
	}
}

func Test_VkGetEmailInvalidToken(t *testing.T) {
	t.Parallel()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	vk := Vk{
		logger:         logger,
	}
	invalidTokens := []string{"", "abcefre", "32424++_!>?|~`"}
	for index := range invalidTokens {
		if _, err = vk.GetID(invalidTokens[index]); !errors.Is(err, ErrValidate) {
			t.Fatalf("expected validation error: %v", err)
		}
	}
}
