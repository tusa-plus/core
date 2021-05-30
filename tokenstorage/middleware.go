package tokenstorage

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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
			if errors.Is(err, ErrTokenExpired) || errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrInvalidSignature) || errors.Is(err, ErrUnknown) {
				return ctx.SendStatus(401)
			} else {
				return ctx.SendStatus(500)
			}
		}
		tokenTypeRaw, ok := token[TokenTypeProperty]
		if !ok {
			ts.logger.Warn("token doesn't contain token_type",
				zap.String("token_string", tokenString),
			)
			return ctx.SendStatus(401)
		}
		tokenType, ok := tokenTypeRaw.(string)
		if !ok || tokenType != expectedTokenType {
			ts.logger.Warn("token_type is not string",
				zap.String("token_string", tokenString),
			)
			return ctx.SendStatus(401)
		}
		ctx.Context().SetUserValue("token_data", token)
		return ctx.Next()
	}
}
