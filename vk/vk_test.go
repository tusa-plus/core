package vk

import (
	"errors"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"testing"
)

func Test_VkGetAccount(t *testing.T) {
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
	account, err := vk.GetAccount(cfg.Section("vk").Key("token").String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedId, _ := cfg.Section("vk").Key("id").Uint64()
	if account.Id != expectedId {
		t.Fatalf("id: expected %v, got %v", expectedId, account.Id)
	}
	expectedSex, _ := cfg.Section("vk").Key("sex").Int()
	if account.Sex != expectedSex {
		t.Fatalf("id: expected %v, got %v", expectedSex, account.Sex)
	}
	expectedName := cfg.Section("vk").Key("name").String()
	if account.Name != expectedName {
		t.Fatalf("name: expected %v, got %v", expectedName, account.Name)
	}
	expectedSurname := cfg.Section("vk").Key("surname").String()
	if account.Surname != expectedSurname {
		t.Fatalf("surname: expected %v, got %v", expectedSurname, account.Surname)
	}
	expectedPhoto := cfg.Section("vk").Key("photo").String()
	if account.Photo != expectedPhoto {
		t.Fatalf("photo: expected %v, got %v", expectedPhoto, account.Photo)
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
		if _, err = vk.GetAccount(invalidTokens[index]); !errors.Is(err, ErrValidate) {
			t.Fatalf("expected validation error: %v", err)
		}
	}
}
