package router

import (
	"github.com/XxThunderBlastxX/chamting-api/handler"
	auth2 "github.com/XxThunderBlastxX/chamting-api/handler/auth"
)

//Authentication Route
func authSetup() {
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/signup", auth2.Signup)
}
