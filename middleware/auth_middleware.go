package middleware

import (
	"github.com/XxThunderBlastxX/chamting-api/presenter"
	"github.com/XxThunderBlastxX/chamting-api/utils"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"net/http"
)

// AuthRequired is a middleware used before accessing any api endpoint
func AuthRequired(ctx *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:    utils.JwtSecretKey,
		SigningMethod: utils.JwtSignInMethod,
		TokenLookup:   "header:Authorization",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusUnauthorized).JSON(presenter.AuthErr(err))
		},
		AuthScheme: "Bearer",
	})(ctx)
}
