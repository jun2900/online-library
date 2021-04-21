package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/database"
	"github.com/jun2900/online-library/models"
)

func CreateAuthorPost(c *fiber.Ctx) error {
	db := database.DBConn
	author := new(models.Author)

	if err := c.BodyParser(author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "review your input"})
	}

	db.Create(&author)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "author created"})
}
