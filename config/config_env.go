package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

//Env func for retrieving variables from .env
func Env(key string) string {
	//Load variable from file .env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error in loading .env files!! ")
	}
	return os.Getenv(key)
}
