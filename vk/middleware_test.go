package vk

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tusa-plus/core/utils"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"net/http/httptest"
	"testing"
)

func createVkApp(t *testing.T) *fiber.App {
	app := fiber.New()
	pool := utils.NewHTTPClientPool()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	vk := Vk{
		logger:         logger,
		httpClientPool: &pool,
	}
	app.Use(NewVkEmailMiddleware(&vk))
	return app
}

func Test_MiddlewareGetValidEmail(t *testing.T) {
	t.Parallel()
	cfg, err := ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := createVkApp(t)
	app.Get("/", func(ctx *fiber.Ctx) error {
		id, err := ctx.Context().UserValue("id").(string)
		if !err {
			t.Fatalf("unexpected error parsing email")
		}
		expectedID := cfg.Section("vk").Key("id").String()
		if id != expectedID {
			t.Fatalf("invalid ID: expected %v, got %v", expectedID, id)
		}
		return ctx.Status(200).SendString("{}")
	})
	vkToken := cfg.Section("vk").Key("token").String()
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("access_token", vkToken)
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}

func Test_MiddlewareGetNoToken(t *testing.T) {
	t.Parallel()
	app := createVkApp(t)
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("{}")
	})
	request := httptest.NewRequest("GET", "/", nil)
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 401 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}

func Test_MiddlewareGetInvalidToken(t *testing.T) {
	t.Parallel()
	app := createVkApp(t)
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("{}")
	})
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Bearer 1234324234")
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 401 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}
