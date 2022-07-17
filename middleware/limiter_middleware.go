package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

//RateLimiting middleware used for rate limiting at each end point
func RateLimiting() fiber.Handler {
	return limiter.New(limiter.Config{Expiration: time.Second * 60, Max: 10 * 2, LimiterMiddleware: limiter.SlidingWindow{}, LimitReached: func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"success": false, "data": "", "error": "too many requests!!"})
	}})
}
