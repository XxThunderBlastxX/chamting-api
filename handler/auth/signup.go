package auth

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

// GenerateCost a random cost
func GenerateCost() int {
	max := 17
	min := 7
	randCrypto := rand.Intn(max-min) + min
	return randCrypto
}

//Returns the password as hash
func hashPassword(password string) (string, error) {
	cost := GenerateCost()
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// Signup handler for setting new user credentials
func Signup(c *fiber.Ctx) error {
	type NewUser struct {
		UserName string `json:"username" bson:"userName"`
		Email    string `json:"email" bson:"email"`
		Password string `json:"password" bson:"password"`
	}
	u := new(NewUser)

	//Get newUser data
	err := c.BodyParser(u)
	if err != nil {
		return c.JSON(fiber.Map{"error": "Something went wrong", "message": err})
	}

	//Hash the received password
	u.Password, err = hashPassword(u.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Couldn't hash password", "message": err})
	}

	uname := u.UserName
	email := u.Email
	password := u.Password

	return c.JSON(fiber.Map{"UserName": uname, "email": email, "password": password})
}
