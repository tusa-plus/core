package facebook

import (
	"github.com/tusa-plus/core/common"
	"gopkg.in/ini.v1"
	"testing"
)

var facebookTestCfg *ini.File

func Test_loadFacebookTestConfig(t *testing.T) {
	var err error
	facebookTestCfg, err = ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_FacebookGetEmail(t *testing.T) {
	pool := common.NewHttpClientPool()
	facebook := Facebook{
		httpClientPool: &pool,
	}
	email, err := facebook.GetEmail(facebookTestCfg.Section("fb").Key("token").String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedEmail := facebookTestCfg.Section("fb").Key("email").String()
	if email != expectedEmail {
		t.Fatalf("invalid email: expected %v, got %v", expectedEmail, email)
	}
}

func Test_FacebookGetEmailInvalidToken(t *testing.T) {
	pool := common.NewHttpClientPool()
	facebook := Facebook{
		httpClientPool: &pool,
	}
	_, err := facebook.GetEmail("xxttzz")
	if err != ErrValidate {
		t.Fatalf("expected validation error: %v", err)
	}
}
