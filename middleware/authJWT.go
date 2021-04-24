package middleware

import (
	"context"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/database"
)

func VerifyToken(c *fiber.Ctx) error {
	accessToken := c.Get("x-access-token")

	if accessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "No token provided"})
	}

	var ctx = context.Background()
	rdb := database.Rdb
	sismemberReply, _ := rdb.SIsMember(ctx, "blacklist_token", accessToken).Result()
	if sismemberReply {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token blacklisted"})
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
