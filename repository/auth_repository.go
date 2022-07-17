package repository

import (
	"context"
	"github.com/XxThunderBlastxX/chamting-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

//AuthRepo interface allows us to access the CRUD Operations in mongo here
type AuthRepo interface {
	AddUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (user *models.User, err error)
}

type auth struct {
	Collection *mongo.Collection
}

//NewAuthRepo is the single instance repo that is being created.
func NewAuthRepo(collection *mongo.Collection) AuthRepo {
	return &auth{Collection: collection}
}

//AddUser is a mongo repo to add new users
func (a *auth) AddUser(user *models.User) (*models.User, error) {
	user.CreatedAt = time.Now()

	_, err := a.Collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//GetUserByEmail is a mongo repo to check whether user with same email exist or not
func (a *auth) GetUserByEmail(email string) (user *models.User, err error) {
	err = a.Collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
