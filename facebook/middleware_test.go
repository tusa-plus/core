package facebook

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

func createFacebookApp(t *testing.T) *fiber.App {
	app := fiber.New()
	pool := utils.NewHTTPClientPool()
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

func Test_MiddlewareGetValidEmailJSON(t *testing.T) {
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
	params := LoginWithFacebookRequest{
		Token: facebookToken,
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
	params := LoginWithFacebookRequest{
		Token: facebookToken,
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
	app := createFacebookApp(t)
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("{}")
	})
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Content-Type", "application/json")
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
	app := createFacebookApp(t)
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("{}")
	})
	params := LoginWithFacebookRequest{
		Token: "1234324234",
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
