package helpers

import (
	"github.com/XxThunderBlastxX/chamting-api/model"
	"github.com/gofiber/fiber/v2"
)

func GetUserByEmail(email string, c *fiber.Ctx) (*model.User, error) {
	//Instance of model.User
	var user model.User

	//Finding the received email
	err := userDb.FindOne(c.Context(), fiber.Map{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername Get the user credential by username
func GetUserByUsername(uname string, c *fiber.Ctx) (*model.User, error) {
	//Instance of model.User
	var user *model.User
	//Finding the received username
	err := userDb.FindOne(c.Context(), fiber.Map{"username": uname}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, err
}
