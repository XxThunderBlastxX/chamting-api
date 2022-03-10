package router

import auth2 "github.com/XxThunderBlastxX/chamting-api/handler/auth"

//Authentication Route
func authSetup() {
	auth := api.Group("/auth")
	auth.Get("/login", auth2.Login)
	auth.Get("/signup", auth2.Signup)
}
