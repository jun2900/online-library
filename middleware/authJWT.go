package middleware

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func VerifyToken(c *fiber.Ctx) error {
	accessToken := c.Get("x-access-token")

	if accessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "No token provided"})
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, c.JSON(fiber.Map{"status": "error", "message": "Unexpected signing method"})
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token expired"})
	}

	if _, err := token.Claims.(jwt.Claims); err && !token.Valid {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token is not valid"})
	}

	return c.Next()
}
