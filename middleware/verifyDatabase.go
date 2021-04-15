package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/database"
	"github.com/jun2900/online-library/models"
)

func CheckEmailExist(c *fiber.Ctx) error {
	db := database.DBConn

	type InputLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input InputLogin
	user := new(models.User)

	c.BodyParser(input)

	if err := db.Where(&models.User{Email: input.Email}).First(&user).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Email is already exist"})
	}

	return c.Next()
}
