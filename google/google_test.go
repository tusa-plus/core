package google

import (
	"github.com/tusa-plus/core/common"
	"gopkg.in/ini.v1"
	"testing"
)

var google = Google{
	httpClientPool: common.NewHttpClientPool(),
	tokenType:      "Bearer",
}

func Test_ManualGoogleGetEmail(t *testing.T) {
	cfg, err := ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	email, err := google.GetEmail(cfg.Section("google").Key("token").String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedEmail := cfg.Section("google").Key("email").String()
	if email != expectedEmail {
		t.Fatalf("invalid email: expected %v, got %v", expectedEmail, email)
	}
}
