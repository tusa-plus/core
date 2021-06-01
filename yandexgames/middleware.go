package yandexgames

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewYandexGamesMiddleware(yandexGames *YandexGames) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		inputArray := strings.Split(ctx.Get(fiber.HeaderAuthorization, ""), ".")
		if len(inputArray) != 2 {
			return ctx.SendStatus(401)
		}
		err := yandexGames.ValidateSignature(inputArray[0], inputArray[1])
		if err != nil {
			if errors.Is(err, ErrValidate) {
				return ctx.SendStatus(401)
			} else {
				return ctx.SendStatus(500)
			}
		}
		return ctx.Next()
	}
}
