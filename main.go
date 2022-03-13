package main

import (
	"context"
	"fmt"
	"github.com/XxThunderBlastxX/chamting-api/config"
	"github.com/XxThunderBlastxX/chamting-api/db"
	"github.com/XxThunderBlastxX/chamting-api/handler"
	"github.com/XxThunderBlastxX/chamting-api/helpers"
	"github.com/XxThunderBlastxX/chamting-api/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

var database *mongo.Database

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error in loading .env files!! ")
	}
	//Config and Connect to Mongo
	dbClient := db.ConnectDb()
	database = dbClient.Database("chamting_app")
	defer func() {
		if err := dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	helpers.Setup(&config.HelperConfig{
		DB: database,
	})
	handler.Setup()
	//Initiated the fiber instance
	app := fiber.New()
	app.Use(cors.New())

	//Setup Router
	router.Route(app)

	//Application listening to the port
	log.Fatal(app.Listen(":8080"))
}
