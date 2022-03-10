package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

//Config func for Global config for the app
func Config(key string) string {
	//Load variable from file .env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error in loading .env files!! ")
	}
	return os.Getenv(key)
}
