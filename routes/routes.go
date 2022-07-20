package routes

import (
	"github.com/XxThunderBlastxX/chamting-api/controller"
	"github.com/XxThunderBlastxX/chamting-api/middleware"
	"github.com/XxThunderBlastxX/chamting-api/service"
	webs "github.com/XxThunderBlastxX/chamting-api/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/websocket/v2"
	"log"
)

func Router(app *fiber.App, authService service.AuthService) {
	//Initial Route
	//app.Get("/", middleware.RateLimiting(), controller.InitialRoute())
	app.Static("/", "./template/home.html")

	//Fiber Monitor
	app.Get("/monitor", middleware.RateLimiting(), monitor.New(monitor.Config{Title: "Chamting-API"}))

	//Authentication group route
	auth := app.Group("/auth")
	auth.Post("/signup", middleware.RateLimiting(), controller.SignUp(authService))
	auth.Post("/signin", middleware.RateLimiting(), controller.SignIn(authService))

	//Websocket group route
	ws := app.Group("/ws")
	go webs.RunHub()
	ws.Get("/", websocket.New(func(c *websocket.Conn) {
		// When the function returns, unregister the client and close the connection
		defer func() {
			webs.Unregister <- c
			c.Close()
		}()

		// Register the client
		webs.Register <- c

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				// Broadcast the received message
				webs.Broadcast <- string(message)
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))

}
