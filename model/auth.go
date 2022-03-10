package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Auth struct {
	Id       primitive.ObjectID `bson:"id"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	UserName string             `json:"userName" bson:"userName"`
}
