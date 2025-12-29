package main

import (
	"log"

	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/database"
	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()

	app := fiber.New()

	// Setup Routes
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
