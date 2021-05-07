package middleware

import (
	"context"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/database"
	"github.com/jun2900/online-library/models"
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

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token is not valid"})
	}

	roleId, _ := token.Claims.(jwt.MapClaims)["role_id"]
	c.Locals("roleId", roleId)

	return c.Next()
}

func IsAdmin(c *fiber.Ctx) error {
	db := database.DBConn
	role := new(models.Role)

	db.First(&role, c.Locals("roleId"))
	if role.Name != "admin" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "user is not an admin"})
	}
	return c.Next()
}
