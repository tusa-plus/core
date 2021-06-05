package tokenstorage

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory"
	"go.uber.org/zap"
	"net/http/httptest"
	"testing"
	"time"
)

func createTSApp(t *testing.T) (*fiber.App, *TokenStorage) {
	app := fiber.New()
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	ts, err := NewTokenStorage([]byte("testsecretkey"), logger, memory.New(), time.Second, time.Second*2)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	app.Use("/access", NewCheckTokenMiddleware(ts, "access"))
	app.Use("/refresh", NewCheckTokenMiddleware(ts, "refresh"))
	return app, ts
}

func Test_MiddlewareTestAccessToken(t *testing.T) {
	t.Parallel()
	app, ts := createTSApp(t)
	claims := map[string]interface{}{
		"test_date": "12345678",
	}
	app.Get("/access", func(ctx *fiber.Ctx) error {
		tokenData, ok := ctx.Context().UserValue("token_data").(map[string]interface{})
		if !ok {
			t.Fatalf("failed to convert tokenData to map")
		}
		for key, value := range claims {
			if tokenData[key] != value {
				t.Fatalf("wrong value in claims: got %v, expected %v", tokenData[key], value)
			}
		}
		return ctx.SendStatus(200)
	})
	access, _, err := ts.NewTokenPair(claims)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	request := httptest.NewRequest("GET", "/access", nil)
	request.Header.Add("Authorization", access)
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}

func Test_MiddlewareTestRefreshToken(t *testing.T) {
	t.Parallel()
	app, ts := createTSApp(t)
	claims := map[string]interface{}{
		"test_date": "12345678",
	}
	app.Get("/refresh", func(ctx *fiber.Ctx) error {
		tokenData, ok := ctx.Context().UserValue("token_data").(map[string]interface{})
		if !ok {
			t.Fatalf("failed to convert tokenData to map")
		}
		for key, value := range claims {
			if tokenData[key] != value {
				t.Fatalf("wrong value in claims: got %v, expected %v", tokenData[key], value)
			}
		}
		return ctx.SendStatus(200)
	})
	_, refresh, err := ts.NewTokenPair(claims)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	request := httptest.NewRequest("GET", "/refresh", nil)
	request.Header.Add("Authorization", refresh)
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}

func Test_MiddlewareTestExpired(t *testing.T) {
	t.Parallel()
	app, ts := createTSApp(t)
	claims := map[string]interface{}{
		"test_date": "12345678",
	}
	app.Get("/refresh", func(ctx *fiber.Ctx) error {
		tokenData, ok := ctx.Context().UserValue("token_data").(map[string]interface{})
		if !ok {
			t.Fatalf("failed to convert tokenData to map")
		}
		for key, value := range claims {
			if tokenData[key] != value {
				t.Fatalf("wrong value in claims: got %v, expected %v", tokenData[key], value)
			}
		}
		return ctx.SendStatus(200)
	})
	_, refresh, err := ts.NewTokenPair(claims)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = ts.ExpireToken(refresh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	request := httptest.NewRequest("GET", "/refresh", nil)
	request.Header.Add("Authorization", refresh)
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 401 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}

func Test_MiddlewareTestInvalidToken(t *testing.T) {
	t.Parallel()
	app, _ := createTSApp(t)
	claims := map[string]interface{}{
		"test_date": "12345678",
	}
	app.Get("/access", func(ctx *fiber.Ctx) error {
		tokenData, ok := ctx.Context().UserValue("token_data").(map[string]interface{})
		if !ok {
			t.Fatalf("failed to convert tokenData to map")
		}
		for key, value := range claims {
			if tokenData[key] != value {
				t.Fatalf("wrong value in claims: got %v, expected %v", tokenData[key], value)
			}
		}
		return ctx.SendStatus(200)
	})
	request := httptest.NewRequest("GET", "/access", nil)
	request.Header.Add("Authorization", "12321312321")
	response, err := app.Test(request, 2000)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 401 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}
