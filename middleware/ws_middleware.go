package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func WsGetRoomId() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		return ctx.Next()
	}

}
