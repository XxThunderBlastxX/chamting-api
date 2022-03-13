package handler

import (
	"github.com/XxThunderBlastxX/chamting-api/helpers"
	"github.com/XxThunderBlastxX/chamting-api/model"
	"github.com/gofiber/fiber/v2"
)

// Login handler
func Login(c *fiber.Ctx) error {
	//Input Structure
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}
	var input LoginInput
	var ud *model.User

	//Getting InputData
	err := c.BodyParser(&input)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "something went wrong on login", "data": err})
	}

	identity := input.Identity
	pass := input.Password
	email, _ := helpers.GetUserByEmail(identity, c)
	uname, _ := helpers.GetUserByUsername(identity, c)

	//Checking if the identity entered is false
	if email != nil && uname != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "cant find any email or username with matching identity"})
	}

	if email == nil {
		ud = &model.User{
			ID:       uname.ID,
			UserName: uname.UserName,
			Email:    uname.Email,
			Password: uname.Password,
		}
	} else {
		ud = &model.User{
			ID:       email.ID,
			UserName: email.UserName,
			Email:    email.Email,
			Password: email.Password,
		}
	}
	//Checking if the entered password matches with the DB
	if !ud.CheckPass(pass) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "password does not match"})
	}
	t, err := ud.GetJWT()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}
