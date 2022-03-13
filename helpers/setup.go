package helpers

import (
	"github.com/XxThunderBlastxX/chamting-api/config"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database
var userDb *mongo.Collection

func Setup(config *config.HelperConfig) {
	db = config.DB
	userDb = db.Collection("user")
}
