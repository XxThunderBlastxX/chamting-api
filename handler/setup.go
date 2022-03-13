package handler

import (
	"github.com/XxThunderBlastxX/chamting-api/config"
	"go.mongodb.org/mongo-driver/mongo"
)

var app *config.HandlerConfig
var userDb *mongo.Collection

func Setup(config *config.HandlerConfig) {
	app = config
	userDb = app.DB.Collection("user")
}
