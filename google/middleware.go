package google

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

type LoginWithGoogleRequest struct {
	Token string `json:"token"`
}

func NewGoogleEmailMiddleware(google *Google) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		type EmptyResponse struct{}
		var request LoginWithGoogleRequest
		if err := ctx.BodyParser(&request); err != nil {
			return ctx.Status(401).JSON(EmptyResponse{})
		}
		email, err := google.GetEmail(ctx.Context(), request.Token)
		if err != nil {
			if errors.Is(err, ErrValidate) {
				return ctx.Status(401).SendString("{}")
			} else {
				return ctx.Status(500).SendString("{}")
			}
		}
		ctx.Context().SetUserValue("email", email)
		return ctx.Next()
	}
}
