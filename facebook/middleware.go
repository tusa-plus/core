package facebook

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

type LoginWithFacebookRequest struct {
	Token string `json:"token" xml:"token"`
}

func NewFacebookEmailMiddleware(facebook Facebook) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		type EmptyResponse struct{}
		var request LoginWithFacebookRequest
		if err := ctx.BodyParser(&request); err != nil {
			return ctx.Status(400).JSON(EmptyResponse{})
		}
		email, err := facebook.GetEmail(ctx.Context(), request.Token)
		if err != nil {
			if errors.Is(err, ErrValidate) {
				return ctx.Status(401).JSON(EmptyResponse{})
			} else {
				return ctx.Status(500).JSON(EmptyResponse{})
			}
		}
		ctx.Context().SetUserValue("email", email)
		return ctx.Next()
	}
}
