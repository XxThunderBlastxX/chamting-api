package presenter

import (
	"github.com/XxThunderBlastxX/chamting-api/models"
	"github.com/gofiber/fiber/v2"
)

func AuthErr(err error) *fiber.Map {
	return &fiber.Map{"success": false, "data": "", "error": err.Error()}
}

func AuthPassErr(err error) *fiber.Map {
	return &fiber.Map{"success": false, "data": "", "error": err.Error()}
}

func AuthSuccess(user *models.User, token string) *fiber.Map {
	return &fiber.Map{"success": true, "token": token, "data": user, "error": ""}
}
