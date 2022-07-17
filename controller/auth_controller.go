package controller

import (
	"github.com/XxThunderBlastxX/chamting-api/models"
	"github.com/XxThunderBlastxX/chamting-api/presenter"
	"github.com/XxThunderBlastxX/chamting-api/service"
	"github.com/XxThunderBlastxX/chamting-api/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SignUp handler/controller for adding new users
// @Summary Post request to auth collection
// @Description Post user data to auth collection
// @Tags Authentication
// @Accept json
// @Produce json
// @Param api query string true "API Key"
// @Param Authorization header string true "Authorization token"
// @Param Body body models.SignUpInput true "Request send to auth collection"
// @Success 200 {object} models.SignUpResponse
// @Failure 422 {object} models.SignUpResponseFailure
// @Failure 409 {object} models.SignUpResponseFailure
// @Failure 500 {object} models.SignUpResponseFailure
// @Failure 503 {object} models.SignUpResponseFailure
// @Router /auth/signup [post]
func SignUp(authService service.AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var user models.User

		err := ctx.BodyParser(&user)
		if err != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(presenter.AuthErr(err))
		}

		email := user.Email
		emailExist, _ := authService.GetUserByEmail(email)
		if emailExist != nil {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"success": false, "data": "", "error": "email already exist"})
		}

		user.Email = utils.TrimString(user.Email)
		user.Password, err = utils.HashPassword(user.Password)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.AuthErr(err))
		}

		user.Id = primitive.NewObjectID()

		result, userErr := authService.AddUser(&user)
		if userErr != nil {
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(presenter.AuthErr(userErr))
		}

		token, tokenErr := utils.GenerateToken(user.Id.Hex())
		if tokenErr != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.AuthErr(tokenErr))
		}

		return ctx.Status(fiber.StatusOK).JSON(presenter.AuthSuccess(result, token))
	}
}

//SignIn handler/controller for Signing In users
// @Summary Post request to auth collection
// @Description Post user email and password to auth collection
// @Tags Authentication
// @Accept json
// @Produce json
// @Param api query string true "API Key"
// @Param Authorization header string true "Authorization token"
// @Param Body body models.SignInInput true "Request send to auth collection"
// @Success 200 {object} models.SignInResponse
// @Failure 422 {object} models.SignInResponseFailure
// @Failure 409 {object} models.SignInResponseFailure
// @Failure 500 {object} models.SignInResponseFailure
// @Failure 503 {object} models.SignInResponseFailure
// @Router /auth/sigin [post]
func SignIn(authService service.AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var inputData models.User

		err := ctx.BodyParser(&inputData)
		if err != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(presenter.AuthErr(err))
		}

		inputData.Email = utils.TrimString(inputData.Email)

		user, userErr := authService.GetUserByEmail(inputData.Email)
		if userErr != nil {
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(presenter.AuthErr(userErr))
		}

		verifyPassErr := utils.VerifyPassword(inputData.Password, user.Password)
		if verifyPassErr != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.AuthPassErr(verifyPassErr))
		}

		token, tokenErr := utils.GenerateToken(user.Id.Hex())
		if tokenErr != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.AuthErr(tokenErr))
		}

		return ctx.Status(fiber.StatusOK).JSON(presenter.AuthSuccess(user, token))
	}
}
