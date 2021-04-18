package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/database"
	"github.com/jun2900/online-library/models"
)

func CheckEmailExisted(c *fiber.Ctx) error {
	db := database.DBConn
	user := new(models.User)

	if err := db.Where(&models.User{Email: user.Email}).First(&user).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "user existed"})
	}
	return c.Next()
}
