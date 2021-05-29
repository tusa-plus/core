package tokenstorage

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func NewCheckAccessTokenMiddleware(ts *TokenStorage) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		inputArray := strings.Split(ctx.Get(fiber.HeaderAuthorization, ""), " ")
		if len(inputArray) != 2 || inputArray[0] != "Bearer" {
			return ctx.Status(401).SendString("{}")
		}
		tokenString := inputArray[1]
		token, err := ts.ParseToken(tokenString)
		if err != nil {
			switch err {
			case ErrTokenExpired:
				return ctx.Status(401).SendString("{}")
			case ErrInvalidToken, ErrInvalidSignature:
				return ctx.Status(400).SendString("{}")
			default:
				return ctx.Status(500).SendString("{}")
			}
		}
		ctx.Context().SetUserValue("token_data", token)
		return ctx.Next()
	}
}
