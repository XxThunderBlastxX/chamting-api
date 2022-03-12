package router

import (
	"github.com/XxThunderBlastxX/chamting-api/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

var api fiber.Router

func Route(a *fiber.App) {
	//Initial Route
	a.Get("/", handler.Hello)

	//Monitor Dashboard
	a.Get("/monitor", monitor.New())

	//Default /api Route
	api = a.Group("/api")

	//Authentication Route
	authSetup()
}
