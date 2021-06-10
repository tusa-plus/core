package yandexgames

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tusa-plus/core/utils"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

func createYandexGamesApp(t *testing.T, key string) *fiber.App {
	app := fiber.New()
	pool := utils.NewHTTPClientPool()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	yandexGames := YandexGames{
		logger:         logger,
		httpClientPool: &pool,
		secret:         []byte(key),
	}
	app.Use(NewYandexGamesMiddleware(&yandexGames))
	return app
}

func Test_MiddlewareValidVerification(t *testing.T) {
	t.Parallel()
	cfg, err := ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	key := cfg.Section("yg").Key("key").String()
	app := createYandexGamesApp(t, key)
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("{}")
	})
	signature := cfg.Section("yg").Key("signature").String()
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", signature)
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}

func Test_MiddlewareInvalidKey(t *testing.T) {
	t.Parallel()
	cfg, err := ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	key := "1234567890"
	app := createYandexGamesApp(t, key)
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("{}")
	})
	signature := cfg.Section("yg").Key("signature").String()
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", signature)
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 401 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}
