package utils

import (
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

// generateCost a random cost
func generateCost() int {
	max := 8
	min := 5
	randCost := rand.Intn(max-min) + min
	return randCost
}

//TrimString trims leading and trailing whitespaces and also lowers the string
func TrimString(str string) string {
	return strings.TrimSpace(strings.ToLower(str))
}

//HashPassword generates hashed password from given password
func HashPassword(pass string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), generateCost())
	if err != nil {
		return "", err
	}
	return string(hashPass), nil
}

//VerifyPassword is a method to compare the hash and entered password
func VerifyPassword(pass string, hashPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass))
	if err != nil {
		return err
	}
	return nil
}

//GenerateToken is a method to generate new JWT token
func GenerateToken(uid string) (string, error) {
	claims := jwt.StandardClaims{
		Id:        uid,
		Issuer:    uid,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30 * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}
