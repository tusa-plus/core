package facebook

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func NewFacebookEmailMiddleware(facebook *Facebook) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		inputArray := strings.Split(ctx.Get(fiber.HeaderAuthorization, ""), " ")
		if len(inputArray) != 2 || inputArray[0] != "Bearer" {
			return ctx.SendStatus(401)
		}
		tokenString := inputArray[1]
		email, err := facebook.GetEmail(ctx.Context(), tokenString)
		if err != nil {
			if errors.Is(err, ErrValidate) {
				return ctx.SendStatus(401)
			} else {
				return ctx.SendStatus(500)
			}
		}
		ctx.Context().SetUserValue("email", email)
		return ctx.Next()
	}
}
