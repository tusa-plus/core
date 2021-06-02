package vk

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func NewVkEmailMiddleware(vk *Vk) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		inputArray := strings.Split(ctx.Get(fiber.HeaderAuthorization, ""), " ")
		if len(inputArray) != 2 || (inputArray[0] != "access_token" && inputArray[0] != "Access_token") {
			return ctx.Status(401).SendString("{}")
		}
		tokenString := inputArray[1]
		id, err := vk.GetID(ctx.Context(), tokenString)
		if err != nil {
			if errors.Is(err, ErrValidate) {
				return ctx.SendStatus(401)
			} else {
				return ctx.SendStatus(500)
			}
		}
		ctx.Context().SetUserValue("id", id)
		return ctx.Next()
	}
}
