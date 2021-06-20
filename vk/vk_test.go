package vk

import (
	"errors"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"testing"
)

func Test_VkGetId(t *testing.T) {
	t.Parallel()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	cfg, err := ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vk := NewVk(logger)
	id, err := vk.GetID(cfg.Section("vk").Key("token").String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedId, _ := cfg.Section("vk").Key("id").Uint64()
	if id != expectedId {
		t.Fatalf("id email: expected %v, got %v", expectedId, id)
	}
}

func Test_VkGetEmailInvalidToken(t *testing.T) {
	t.Parallel()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	vk := NewVk(logger)
	invalidTokens := []string{"", "abcefre", "32424++_!>?|~`"}
	for index := range invalidTokens {
		if _, err = vk.GetID(invalidTokens[index]); !errors.Is(err, ErrValidate) {
			t.Fatalf("expected validation error: %v", err)
		}
	}
}

func Test_VkMockOk(t *testing.T) {
	t.Parallel()
	vk := NewMockVk()
	result, err := vk.GetID("123")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if result != 123 {
		t.Fatalf("id: expected %v, got %v", 123, result)
	}
}

func Test_VkMockInvalid(t *testing.T) {
	t.Parallel()
	vk := NewMockVk()
	_, err := vk.GetID("aafw")
	if err == nil {
		t.Fatalf("expected validation error")
	}
}
