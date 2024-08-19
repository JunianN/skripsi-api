package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// func Protected() fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         userToken := c.Get("Authorization")[7:] // Bearer token
//         token, err := jwt.ParseWithClaims(userToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
//             return []byte("secret"), nil // Use your secret key
//         })

//         if err != nil || !token.Valid {
//             return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
//         }

//         claims := token.Claims.(*jwt.RegisteredClaims)
//         c.Locals("userID", claims.Subject) // Store userID from token in Locals

//         return c.Next()
//     }
// }

func Authenticated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or missing Authorization header"})
		}

		tokenString := authHeader[7:]

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - Invalid token"})
		}

		if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
			userID := (*claims)["user_id"]
			userRole := (*claims)["userRole"]
			if userID == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - User ID not found in token"})
			}
			c.Locals("userID", userID)
			c.Locals("userRole", userRole)
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
}
