package yandexgames

import (
	"strings"
	"testing"

	"github.com/tusa-plus/core/utils"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

func Test_YandexGamesValidateSignature(t *testing.T) {
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
	yandexGames := YandexGames{
		httpClientPool: &pool,
		logger:         logger,
		secret:         []byte(cfg.Section("yg").Key("key").String()),
	}
	inputArray := strings.Split(cfg.Section("yg").Key("signature").String(), ".")
	if len(inputArray) != 2 {
		t.Fatalf("invalid signature")
	}
	err = yandexGames.ValidateSignature(inputArray[0], inputArray[1])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_YandexGamesInvalidKey(t *testing.T) {
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
	yandexGames := YandexGames{
		httpClientPool: &pool,
		logger:         logger,
		secret:         []byte("1234567890"),
	}
	inputArray := strings.Split(cfg.Section("yg").Key("signature").String(), ".")
	if len(inputArray) != 2 {
		t.Fatalf("invalid signature")
	}
	err = yandexGames.ValidateSignature(inputArray[0], inputArray[1])
	if err == nil {
		t.Fatalf("Invalid key is valid")
	}
}
