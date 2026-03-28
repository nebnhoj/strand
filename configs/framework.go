package configs

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func SetFiberConfig() fiber.Config {
	godotenv.Load() //nolint: ignore missing .env in container (vars set via env_file)
	return fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  os.Getenv("APP_HEADER"),
		AppName:       os.Getenv("APP_NAME"),
	}
}
