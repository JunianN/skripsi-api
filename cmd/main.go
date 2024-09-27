package main

import (
	"log"
	"translation-app-backend/internal/database"
	"translation-app-backend/internal/handlers"
	"translation-app-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // set the max body size to 10MB
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://doc-translation.vercel.app/",
		// AllowOrigins: "http://localhost:3000",
		ExposeHeaders: "Content-Disposition",
	}))

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	routes.SetupRoutes(app, db)

	// Set up the cron job
	c := cron.New()
	_, err2 := c.AddFunc("@hourly", handlers.CheckAndDeclineUnconfirmedDocuments)
	if err2 != nil {
		panic(err)
	}
	c.Start()
	defer c.Stop()

	log.Fatal(app.Listen(":3001"))
}
