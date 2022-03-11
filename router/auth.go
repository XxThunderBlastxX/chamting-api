package router

import auth2 "github.com/XxThunderBlastxX/chamting-api/handler/auth"

//Authentication Route
func authSetup() {
	auth := api.Group("/auth")
	auth.Post("/login", auth2.Login)
	auth.Post("/signup", auth2.Signup)
}
