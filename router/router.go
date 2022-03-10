package router

import (
	"github.com/XxThunderBlastxX/chamting-api/handler"
	"github.com/gofiber/fiber/v2"
)

var api fiber.Router

func Route(a *fiber.App) {
	//Default /api Route
	api = a.Group("/api")
	api.Get("/", handler.Hello)

	//Authentication Route
	authSetup()
}
