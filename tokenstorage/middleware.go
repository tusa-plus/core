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
			validationErrors := []error{ErrTokenExpired, ErrInvalidToken, ErrInvalidSignature, ErrInvalidFields}
			for index := range validationErrors {
				if errors.Is(err, validationErrors[index]) {
					return ctx.Status(401).SendString("{}")
				}
			}
			ts.logger.Error("unknown error parsing token",
				zap.String("token_string", tokenString),
				zap.Error(err),
			)
			return ctx.Status(500).SendString("{}")
		}
		tokenTypeRaw, ok := token[TokenTypeProperty]
		if !ok {
			ts.logger.Warn("token doesn't contain token_type",
				zap.String("token_string", tokenString),
			)
			return ctx.Status(401).SendString("{}")
		}
		tokenType, ok := tokenTypeRaw.(string)
		if !ok || tokenType != expectedTokenType {
			ts.logger.Warn("token_type is not string",
				zap.String("token_string", tokenString),
			)
			return ctx.Status(401).SendString("{}")
		}
		ctx.Context().SetUserValue("token_string", tokenString)
		ctx.Context().SetUserValue("token_data", token)
		return ctx.Next()
	}
}
