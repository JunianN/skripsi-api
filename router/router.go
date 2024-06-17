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

	api := app.Group("/api")

	api.Post("/register", handlers.Register(db))
	api.Post("/login", handlers.Login(db))

	api.Post("/upload", middleware.Authenticated(), handlers.UploadDocument)
	api.Get("/documents", middleware.Authenticated(), handlers.GetDocuments)
	api.Get("/documents/:id", middleware.Authenticated(), handlers.GetDocument)

	api.Get("/documents/:id/discussions", middleware.Authenticated(), handlers.GetDiscussions)
    api.Post("/documents/:id/discussions", middleware.Authenticated(), handlers.PostDiscussion)

	api.Post("/documents/:id/messages", middleware.Authenticated(), handlers.AddMessage)
	api.Get("/documents/:id/messages", middleware.Authenticated(), handlers.GetMessages)

	api.Get("/translations", handlers.ListTranslations(db))
	api.Post("/translations", handlers.AddTranslation(db))

	api.Post("/notifications", handlers.SendNotification(db))

	// Admin routes
    admin := app.Group("/api/admin")
    admin.Use(middleware.Authenticated())
    // admin.Use(middleware.AdminRequired())
	
    admin.Get("/documents", handlers.GetAllDocuments)
	admin.Patch("/documents/:id/assign", handlers.AssignDocument)

	// Group routes for translators
    translators := app.Group("/api/translators")
    translators.Use(middleware.Authenticated())

    translators.Patch("/documents/:id/status", handlers.UpdateDocumentStatus)
	translators.Get("/documents", handlers.GetTranslatorDocuments)
	translators.Post("/documents/:id/upload_translation", handlers.UploadTranslatedDocument)
}
