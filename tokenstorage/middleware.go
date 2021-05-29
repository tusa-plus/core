package tokenstorage

import "github.com/gofiber/fiber/v2"

func NewCheckAccessTokenMiddleware(ts *TokenStorage) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Get(fiber.HeaderAuthorization, "")
		if tokenString == "" {
			return ctx.Status(401).SendString("{}")
		}
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
