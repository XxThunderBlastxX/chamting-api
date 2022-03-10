package router

import (
	"github.com/XxThunderBlastxX/chamting-api/handler"
	"github.com/gofiber/fiber/v2"
)

var api fiber.Router

func Route(a *fiber.App) {
	a.Get("/", handler.Hello)

	//Default /api Route
	api = a.Group("/api")

	//Authentication Route
	authSetup()
}
