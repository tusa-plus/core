package facebook

import (
	"github.com/tusa-plus/core/common"
	"gopkg.in/ini.v1"
	"testing"
)

var facebook = Facebook{
	httpClientPool: common.NewHttpClientPool(),
}

func Test_ManualFbGetEmail(t *testing.T) {
	cfg, err := ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	email, err := facebook.GetEmail(cfg.Section("fb").Key("token").String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedEmail := cfg.Section("fb").Key("email").String()
	if email != expectedEmail {
		t.Fatalf("invalid email: expected %v, got %v", expectedEmail, email)
	}
}
