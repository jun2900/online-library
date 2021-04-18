package controllers

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/database"
	"github.com/jun2900/online-library/models"
	"gorm.io/gorm"
)

func ReadAllPaper(c *fiber.Ctx) error {
	db := database.DBConn
	var papers []models.Paper
	db.Find(&papers)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "Papers": papers})
}

func CreatePaper(c *fiber.Ctx) error {
	const uploadPath = "./uploads"
	db := database.DBConn
	paper := new(models.Paper)

	if err := c.BodyParser(paper); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "review your input"})
	}

	fileType := http.DetectContentType(paper.Content)
	if fileType != "application/pdf" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "filetype must be application/pdf"})
	}

	f, err := os.Create(filepath.Join(uploadPath, paper.Title+".pdf"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error in proccessing file upload", "err": err})
	}
	defer f.Close()

	if _, err := f.Write(paper.Content); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	db.Create(&paper)
	db.Save(&paper)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "message": "paper created"})
}

func DownloadPaper(c *fiber.Ctx) error {
	paperId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error processing param id", "err": err})
	}

	db := database.DBConn
	paper := new(models.Paper)

	if err := db.Select("content").First(&paper, paperId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "messsage": "Paper not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error handling query", "err": err})
	}

	return c.SendString("Hello")
}
