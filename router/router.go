package router

import (
	"translation-app-backend/handlers"
	"translation-app-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Ini API untuk web app sistem layanan penerjemahan dokumen.")
	})

	app.Post("/api/register", handlers.Register(db))
	app.Post("/api/login", handlers.Login(db))

	app.Post("/api/upload", middleware.Authenticated(), handlers.UploadDocument)
	app.Get("/api/documents", middleware.Authenticated(), handlers.GetDocuments)

	app.Post("/api/documents/:id/messages", middleware.Authenticated(), handlers.AddMessage)
	app.Get("/api/documents/:id/messages", middleware.Authenticated(), handlers.GetMessages)

	app.Get("/api/translations", handlers.ListTranslations(db))
	app.Post("/api/translations", handlers.AddTranslation(db))

	app.Post("/api/notifications", handlers.SendNotification(db))

	// Admin routes
    admin := app.Group("/admin")
    admin.Use(middleware.Authenticated())
    // admin.Use(middleware.AdminRequired())
	
    admin.Get("/documents", handlers.GetAllDocuments)
	admin.Patch("/documents/:id/assign", handlers.AssignDocument)

	// Group routes for translators
    translators := app.Group("/translators")
    translators.Use(middleware.Authenticated())
    translators.Patch("/documents/:id/status", handlers.UpdateDocumentStatus)
	translators.Get("/documents", handlers.GetTranslatorDocuments)
	translators.Post("/documents/:id/upload_translation", handlers.UploadTranslatedDocument)
}
