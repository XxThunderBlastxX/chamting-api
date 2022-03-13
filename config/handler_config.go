package config

import "go.mongodb.org/mongo-driver/mongo"

type HandlerConfig struct {
	DB *mongo.Database
}
