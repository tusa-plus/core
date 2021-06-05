package google

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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
	google, err := NewGoogle(logger, &pool)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	app.Use(NewGoogleEmailMiddleware(google))
	return app
}

func Test_MiddlewareGetValidEmailJSON(t *testing.T) {
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
	params := LoginWithGoogleRequest{
		Token: googleToken,
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	request := httptest.NewRequest("GET", "/", bytes.NewBuffer(paramsJSON))
	request.Header.Add("Content-Type", "application/json")
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}

func Test_MiddlewareGetValidEmailXML(t *testing.T) {
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
	params := LoginWithGoogleRequest{
		Token: googleToken,
	}
	paramsXML, err := xml.Marshal(params)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	request := httptest.NewRequest("GET", "/", bytes.NewBuffer(paramsXML))
	request.Header.Add("Content-Type", "application/xml")
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
	if response.StatusCode != 400 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}

func Test_MiddlewareGetInvalidToken(t *testing.T) {
	t.Parallel()
	app := createGoogleApp(t)
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("{}")
	})
	params := LoginWithGoogleRequest{
		Token: "xxxxxxx",
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	request := httptest.NewRequest("GET", "/", bytes.NewBuffer(paramsJSON))
	request.Header.Add("Content-Type", "application/json")
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 401 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}
