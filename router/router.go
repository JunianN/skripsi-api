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

	// auth routes
	app.Post("/register", handlers.Register(db))
	app.Post("api/login", handlers.Login(db))

	// user routes
	api := app.Group("/api")
	api.Use(middleware.Authenticated())

	api.Post("/upload", middleware.Authenticated(), handlers.UploadDocument)
	api.Get("/documents", middleware.Authenticated(), handlers.GetDocuments)
	api.Get("/documents/:id", middleware.Authenticated(), handlers.GetDocument)
	api.Get("/documents/:id/discussions", middleware.Authenticated(), handlers.GetDiscussions)
    api.Post("/documents/:id/discussions", middleware.Authenticated(), handlers.PostDiscussion)
	api.Get("/documents/:id/download", middleware.Authenticated(), handlers.DownloadTranslatedDocument)
	api.Post("/documents/:id/upload-receipt", handlers.UploadPaymentReceipt(db))

	api.Get("/translations", handlers.ListTranslations(db))
	api.Post("/translations", handlers.AddTranslation(db))

	api.Post("/notifications", handlers.SendNotification(db))

	// Admin routes
    admin := app.Group("/api/admin")
    admin.Use(middleware.Authenticated())
    admin.Use(middleware.AdminRequired())

	admin.Get("/documents", handlers.GetAllDocuments(db))
	admin.Get("/documents/:id", handlers.GetDocumentDetails(db))
	admin.Get("/documents/:id/download", handlers.DownloadUserDocument(db))
	admin.Post("/documents/:id/approve", handlers.ApproveDocument(db))
	admin.Post("/documents/:id/reject", handlers.RejectDocument(db))
	admin.Post("/documents/:id/assign", handlers.AssignDocument(db))
	admin.Get("/translators", handlers.GetTranslators(db))
	admin.Get("/documents/:id/translated/download", handlers.DownloadTranslatedFile(db))
	admin.Post("/documents/:id/translated/approve", handlers.ApproveTranslatedDocument(db))
	admin.Post("/documents/:id/translated/reject", handlers.RejectTranslatedDocument(db))
	admin.Get("/documents/:id/payment-receipt", handlers.DownloadPaymentReceipt(db))
	admin.Post("/documents/:id/payment-approve", handlers.ApprovePayment(db))


	// Group routes for translators
    translators := app.Group("/api/translators")
    translators.Use(middleware.Authenticated())

    translators.Patch("/documents/:id/status", handlers.UpdateDocumentStatus)
	translators.Get("/documents", handlers.GetTranslatorDocuments)
	translators.Post("/documents/:id/upload_translation", handlers.UploadTranslatedDocument)
}
