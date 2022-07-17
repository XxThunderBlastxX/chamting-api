package routes

import (
	"github.com/XxThunderBlastxX/chamting-api/controller"
	"github.com/XxThunderBlastxX/chamting-api/middleware"
	"github.com/XxThunderBlastxX/chamting-api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Router(app *fiber.App, authService service.AuthService) {
	//Initial Route
	app.Get("/", middleware.RateLimiting(), controller.InitialRoute())

	//Fiber Monitor
	app.Get("/monitor", middleware.RateLimiting(), monitor.New(monitor.Config{Title: "Chamting-API"}))

	//Authentication group route
	auth := app.Group("/auth")
	auth.Post("/signup", middleware.RateLimiting(), controller.SignUp(authService))
	auth.Post("/signin", middleware.RateLimiting(), controller.SignIn(authService))
}
