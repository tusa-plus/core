package facebook

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func NewFacebookEmailMiddleware(facebook *Facebook) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		inputArray := strings.Split(ctx.Get(fiber.HeaderAuthorization, ""), " ")
		if len(inputArray) != 2 || inputArray[0] != "Bearer" {
			return ctx.Status(401).SendString("{}")
		}
		tokenString := inputArray[1]
		email, err := facebook.GetEmail(ctx.Context(), tokenString)
		if err != nil {
			switch err {
			case ErrValidate:
				return ctx.SendStatus(401)
			default:
				return ctx.SendStatus(500)
			}
		}
		ctx.Context().SetUserValue("facebook_email", email)
		return ctx.Next()
	}
}
