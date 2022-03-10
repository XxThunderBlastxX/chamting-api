package router

import (
	"github.com/XxThunderBlastxX/chamting-api/handler"
	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/", handler.Hello)
}
