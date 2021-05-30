package tokenstorage

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func NewCheckTokenMiddleware(ts *TokenStorage, expectedTokenType string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		inputArray := strings.Split(ctx.Get(fiber.HeaderAuthorization, ""), " ")
		if len(inputArray) != 2 || inputArray[0] != "Bearer" {
			return ctx.Status(401).SendString("{}")
		}
		tokenString := inputArray[1]
		token, err := ts.ParseToken(tokenString)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) || errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrInvalidSignature) {
				return ctx.SendStatus(401)
			} else {
				return ctx.SendStatus(500)
			}
		}
		tokenTypeRaw, ok := token[TokenTypeProperty]
		if !ok {
			return ctx.SendStatus(401)
		}
		tokenType, ok := tokenTypeRaw.(string)
		if !ok || tokenType != expectedTokenType {
			return ctx.SendStatus(401)
		}
		ctx.Context().SetUserValue("token_data", token)
		return ctx.Next()
	}
}
