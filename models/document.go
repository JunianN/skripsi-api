package models

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	UserID                   uint
	TranslatorID             uint // ID of the assigned translator
	Title                    string
	Description              string
	FilePath                 string
	SourceLanguage           string
	TargetLanguage           string
	NumberOfPages            string
	TranslatedFilePath       string    // Path to the translated document
	Status                   string    // e.g., "Pending", "In Progress", "Completed"
	PaymentConfirmed         bool      // Field to check if the payment is confirmed
	ApprovalStatus           string    // e.g., "Pending", "Approved", "Rejected"
	TranslatedApprovalStatus string    // e.g., "Pending", "Approved", "Rejected"
	TranslatorApprovalStatus string    // e.g., "Pending", "Accepted", "Declined"
	PaymentReceiptFilePath   string    // Path to the uploaded payment receipt
	AssignmentTime           time.Time // Time when the document was assigned to the translator
}
