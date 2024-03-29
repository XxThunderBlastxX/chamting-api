package main

import (
	"fmt"
	"github.com/XxThunderBlastxX/chamting-api/database"
	"github.com/XxThunderBlastxX/chamting-api/repository"
	"github.com/XxThunderBlastxX/chamting-api/routes"
	"github.com/XxThunderBlastxX/chamting-api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	//Loads variables from .env
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	//Connect to mongo-database
	db, cancel, dbErr := database.DBConnect()
	if dbErr != nil {
		log.Fatal("Database Connection Error $s ", dbErr)
	}
	fmt.Println("Database Connection Successful 🙌")

	//Instance of authentication handler/service/repository
	authCollection := db.Collection("auth")
	authRepo := repository.NewAuthRepo(authCollection)
	authService := service.NewAuthService(authRepo)

	//Init New Fiber App
	app := fiber.New()

	//Enable CORS
	app.Use(cors.New())

	//Upgrade to Websocket
	app.Use("/ws", func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("allowed", true)
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	//Application Route
	routes.Router(app, authService)

	//closes the database connection
	defer cancel()

	//Listen Application at desired port from .env
	log.Fatal(app.Listen(os.Getenv("PORT")))
}
