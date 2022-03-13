package model

import (
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	UserName string             `json:"username" bson:"username"`
}

type NewUser struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	UserName string `json:"username" bson:"username"`
}

func (user *User) GetJWT() (string, error) {
	//Generating JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.UserName
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	//Getting encoded JWT token
	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	return t, err
}

func (user *User) CheckPass(pass string) bool {
	hash := user.Password

	//Comparing Hashed and received password is correct
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}