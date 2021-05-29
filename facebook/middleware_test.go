package facebook

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/tusa-plus/core/common"
	"gopkg.in/ini.v1"
	"net/http"
	"testing"
)

var app = fiber.New()
var httpClientPool = common.NewHttpClientPool()
var middlewareTestCfg *ini.File

func Test_loadMiddlewareTestConfig(t *testing.T) {
	var err error
	middlewareTestCfg, err = ini.Load("./config_test.ini")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_CreateAppWithMiddleware(t *testing.T) {
	port := middlewareTestCfg.Section("middleware").Key("port").String()
	facebook := Facebook{
		httpClientPool: &httpClientPool,
	}
	app.Use(NewFacebookEmailMiddleware(&facebook))
	go func() {
		t.Fatal(app.Listen(fmt.Sprintf(":%v", port)))
	}()
}

func Test_MiddlewareGetValidEmail(t *testing.T) {
	app.Get("/MiddlewareGetValidEmail", func(ctx *fiber.Ctx) error {
		email, err := ctx.Context().UserValue("facebook_email").(string)
		if !err {
			t.Fatalf("unexpected error parsing email")
			return ctx.Status(500).SendString("{}")
		}
		expectedEmail := facebookTestCfg.Section("fb").Key("email").String()
		if email != expectedEmail {
			t.Fatalf("invalid email: expected %v, got %v", expectedEmail, email)
		}
		return ctx.Status(200).SendString("{}")
	})
	port := middlewareTestCfg.Section("middleware").Key("port").String()
	request, err := http.NewRequest("GET", fmt.Sprintf("http://0.0.0.0:%v/MiddlewareGetValidEmail", port), nil)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	facebookToken := facebookTestCfg.Section("fb").Key("token").String()
	request.Header.Add("Authorization", "Bearer "+facebookToken)
	client := httpClientPool.Get()
	defer httpClientPool.Put(client)
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}

func Test_MiddlewareGetNoToken(t *testing.T) {
	app.Get("/MiddlewareGetNoToken", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("{}")
	})
	port := middlewareTestCfg.Section("middleware").Key("port").String()
	request, err := http.NewRequest("GET", fmt.Sprintf("http://0.0.0.0:%v/MiddlewareGetNoToken", port), nil)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	client := httpClientPool.Get()
	defer httpClientPool.Put(client)
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 401 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}

func Test_MiddlewareGetInvalidToken(t *testing.T) {
	app.Get("/MiddlewareGetInvalidToken", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("{}")
	})
	port := middlewareTestCfg.Section("middleware").Key("port").String()
	request, err := http.NewRequest("GET", fmt.Sprintf("http://0.0.0.0:%v/MiddlewareGetInvalidToken", port), nil)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	request.Header.Add("Authorization", "Bearer 32432432423423423")
	client := httpClientPool.Get()
	defer httpClientPool.Put(client)
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if response.StatusCode != 401 {
		t.Fatalf("unexpected status code %v", response.StatusCode)
	}
}
