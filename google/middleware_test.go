package google

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tusa-plus/core/utils"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"net/http/httptest"
	"testing"
)

func createGoogleApp(t *testing.T) *fiber.App {
	app := fiber.New()
	pool := utils.NewHTTPClientPool()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	google := Google{
		logger:         logger,
		httpClientPool: &pool,
	}
	app.Use(NewGoogleEmailMiddleware(&google))
	return app
}

func Test_MiddlewareGetValidEmail(t *testing.T) {
	t.Parallel()
	cfg, err := ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := createGoogleApp(t)
	app.Get("/", func(ctx *fiber.Ctx) error {
		email, err := ctx.Context().UserValue("email").(string)
		if !err {
			t.Fatalf("unexpected error parsing email")
		}
		expectedEmail := cfg.Section("google").Key("email").String()
		if email != expectedEmail {
			t.Fatalf("invalid email: expected %v, got %v", expectedEmail, email)
		}
		return ctx.Status(200).SendString("{}")
	})
	googleToken := cfg.Section("google").Key("token").String()
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Bearer "+googleToken)
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
	app := createGoogleApp(t)
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
	app := createGoogleApp(t)
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
