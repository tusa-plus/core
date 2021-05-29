package tokenstorage

import (
	"fmt"
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
			switch err {
			case ErrTokenExpired, ErrInvalidToken, ErrInvalidSignature:
				return ctx.SendStatus(401)
			default:
				fmt.Printf("%v\b", err)
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
