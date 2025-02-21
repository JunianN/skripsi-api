package routes

import (
	"translation-app-backend/internal/handlers"
	"translation-app-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Ini API untuk web app sistem layanan penerjemahan dokumen.")
	})

	// auth routes
	app.Post("api/register", handlers.Register(db))
	app.Post("api/login", handlers.Login(db))
	app.Post("api/mail", handlers.Mail(db))

	// user routes
	api := app.Group("/api")
	api.Use(middleware.Authenticated())

	api.Post("/upload", middleware.Authenticated(), handlers.UploadDocument(db))
	api.Get("/documents", middleware.Authenticated(), handlers.GetDocuments(db))
	api.Get("/documents/:id", middleware.Authenticated(), handlers.GetDocument(db))
	api.Get("/documents/:id/discussions", middleware.Authenticated(), handlers.GetDiscussions(db))
	api.Post("/documents/:id/discussions", middleware.Authenticated(), handlers.PostDiscussion(db))
	api.Post("/documents/:id/upload-receipt", handlers.UploadPaymentReceipt(db))
	api.Get("/documents/:id/download", middleware.Authenticated(), handlers.DownloadTranslatedDocument(db))
	api.Post("/ratings", handlers.SubmitRating(db))
	api.Get("/:id/average-rating", handlers.GetTranslatorAverageRating(db))
	api.Get("/documents/:id/rating", handlers.GetRatings(db))

	api.Get("/notifications", handlers.FetchNotifications(db))
	api.Post("/notifications/read", handlers.MarkNotificationsAsRead(db))

	// Admin routes
	admin := app.Group("/api/admin")
	admin.Use(middleware.Authenticated())
	admin.Use(middleware.AdminRequired())

	admin.Post("/register", handlers.RegisterAdmin(db))
	admin.Get("/documents", handlers.GetAllDocuments(db))
	admin.Get("/documents/:id", handlers.GetDocumentDetails(db))
	admin.Get("/documents/:id/download", handlers.DownloadUserDocument(db))
	admin.Post("/documents/:id/approve", handlers.ApproveDocument(db))
	admin.Post("/documents/:id/reject", handlers.RejectDocument(db))
	admin.Get("/translators", handlers.GetTranslators(db))
	admin.Get("/translators/by-language", handlers.GetTranslatorsByLanguage(db))
	admin.Post("/documents/:id/assign", handlers.AssignDocument(db))
	admin.Delete("/translators/:id", handlers.DeleteTranslator(db))
	admin.Get("/documents/:id/translated/download", handlers.DownloadTranslatedFile(db))
	admin.Post("/documents/:id/translated/approve", handlers.ApproveTranslatedDocument(db))
	admin.Post("/documents/:id/translated/reject", handlers.RejectTranslatedDocument(db))
	admin.Get("/documents/:id/payment-receipt", handlers.DownloadPaymentReceipt(db))
	admin.Post("/documents/:id/payment-approve", handlers.ApprovePayment(db))
	admin.Get("/mails", handlers.GetMailSubmissions(db))
	admin.Put("/settings/price", handlers.UpdatePricePerWord(db))

	// Group routes for translators
	translators := app.Group("/api/translator")
	translators.Use(middleware.Authenticated())
	translators.Use(middleware.TranslatorRequired())

	translators.Get("/assigned-documents", handlers.GetAssignedDocuments(db))
	translators.Get("/documents/:id", handlers.GetAssignedDocument(db))
	translators.Post("/documents/:id/approve", handlers.ApproveAssignedDocument(db))
	translators.Post("/documents/:id/decline", handlers.DeclineAssignedDocument(db))
	translators.Get("/documents/:id/download", handlers.DownloadAssignedDocument(db))
	translators.Post("/documents/:id/upload", handlers.UploadTranslatedDocument(db))
}
