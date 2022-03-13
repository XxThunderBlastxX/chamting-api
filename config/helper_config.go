package config

import "go.mongodb.org/mongo-driver/mongo"

type HelperConfig struct {
	DB *mongo.Database
}
