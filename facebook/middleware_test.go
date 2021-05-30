package facebook

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tusa-plus/core/utils"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"net/http/httptest"
	"testing"
)

func createFacebookApp(t *testing.T) *fiber.App {
	t.Parallel()
	app := fiber.New()
	pool := utils.NewHttpClientPool()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	facebook := Facebook{
		logger:         logger,
		httpClientPool: &pool,
	}
	app.Use(NewFacebookEmailMiddleware(&facebook))
	return app
}

func Test_MiddlewareGetValidEmail(t *testing.T) {
	t.Parallel()
	cfg, err := ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := createFacebookApp(t)
	app.Get("/", func(ctx *fiber.Ctx) error {
		email, err := ctx.Context().UserValue("email").(string)
		if !err {
			t.Fatalf("unexpected error parsing email")
		}
		expectedEmail := cfg.Section("fb").Key("email").String()
		if email != expectedEmail {
			t.Fatalf("invalid email: expected %v, got %v", expectedEmail, email)
		}
		return ctx.Status(200).SendString("{}")
	})
	facebookToken := cfg.Section("fb").Key("token").String()
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Bearer "+facebookToken)
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}

func Test_MiddlewareGetNoToken(t *testing.T) {
	app := createFacebookApp(t)
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
	app := createFacebookApp(t)
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
