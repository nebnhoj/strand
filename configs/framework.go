package configs

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func SetFiberConfig() fiber.Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
	return fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  os.Getenv("APP_HEADER"),
		AppName:       os.Getenv("APP_NAME"),
	}
}
