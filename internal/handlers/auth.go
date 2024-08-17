package handlers

import (
	"os"
	"time"
	"translation-app-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Username            string   `json:"username"`
			Email               string   `json:"email"`
			Password            string   `json:"password"`
			Role                string   `json:"role"`
			ProficientLanguages []string `json:"proficient_languages"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})

		}

		if input.Role == "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Cannot register as admin"})
		}

		// Validate inputs
		if input.Email == "" || input.Password == "" || input.Username == "" || input.Role == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "please provide all required fields"})
		}

		// Check for existing user
		var exists models.User
		db.Where("email = ?", input.Email).First(&exists)
		if exists.ID != 0 {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already in use"})
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not hash password"})
		}

		user := models.User{
			Username: input.Username,
			Email:    input.Email,
			Password: string(hashedPassword),
			Role:     input.Role,
		}

		if input.Role == "translator" {
			user.ProficientLanguages = input.ProficientLanguages
		}
		
		// Save the user to the database
		result := db.Create(&user)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
		}

		return c.JSON(user)
	}
}

func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input models.LoginInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		// Validate input
		if input.Email == "" || input.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "please provide email and password"})
		}

		// Retrieve user from database
		var user models.User
		db.Where("email = ?", input.Email).First(&user)
		if user.ID == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}

		// Compare password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "incorrect password"})
		}

		// Generate JWT token
		token, err := generateJWT(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
		}

		return c.JSON(fiber.Map{"message": "login successful", "token": token})
	}
}

func generateJWT(user models.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"userRole": user.Role,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	token, err := claims.SignedString([]byte(os.Getenv("SECRET"))) // Use a secret from env variable in production
	return token, err
}
