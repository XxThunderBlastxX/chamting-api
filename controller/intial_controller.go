package controller

import "github.com/gofiber/fiber/v2"

func InitialRoute() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		Map := fiber.Map{"Name": "Malay Optic Care - API", "Created By": "Koustav Mondal", "Status": "Running", "Version": "1.0.0"}

		return ctx.Status(fiber.StatusOK).JSON(Map)
	}
}
