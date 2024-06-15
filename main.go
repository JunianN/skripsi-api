package main

import (
	"log"
	"translation-app-backend/database"
	"translation-app-backend/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // set the max body size to 10MB
	})

	// Middleware
	app.Use(cors.New()) // Enable CORS

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	router.SetupRoutes(app, db)

	log.Fatal(app.Listen(":3001"))
}
