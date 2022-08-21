package utils

import (
	"github.com/golang-jwt/jwt"
	"os"
)

var (
	JwtSecretKey    = []byte(os.Getenv("SECRET_KEY"))
	JwtSignInMethod = jwt.SigningMethodHS256.Name
)
