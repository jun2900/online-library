package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/database"
	"github.com/jun2900/online-library/models"
)

func ReadAllPaper(c *fiber.Ctx) error {
	db := database.DBConn
	var papers []models.Paper
	if err := db.Find(&papers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error on retreiving papers", "error": err})
	}
	return c.JSON(fiber.Map{"status": "success", "papers": papers})
}

func CreatePaper(c *fiber.Ctx) error {
	db := database.DBConn
	paper := new(models.Paper)
	if err := c.BodyParser(paper); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error processing input", "error": err})
	}

	db.Create(&paper)
	db.Save(&paper)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "message": "Paper created"})
}
