package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/database"
	"github.com/jun2900/online-library/models"
	"gorm.io/gorm"
)

const uploadPath = "./uploads"

func handlingRowError(err error, c *fiber.Ctx) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "messsage": "Paper not found"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error handling query", "err": err})
}

func validatePaperStruct(input models.Paper) []*ErrorResponse {
	return HandlingInput(input)
}

func ReadAllPaper(c *fiber.Ctx) error {
	db := database.DBConn
	var papers []models.Paper
	db.Find(&papers)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "Papers": papers})
}

func ReadSpecificPaper(c *fiber.Ctx) error {
	db := database.DBConn

	paperId := c.Locals("id")

	var paper models.Paper
	if err := db.First(&paper, paperId).Error; err != nil {
		return handlingRowError(err, c)
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "paper": paper})
}

func CreatePaper(c *fiber.Ctx) error {
	db := database.DBConn
	input := new(models.Paper)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "review your input"})
	}

	if err := db.Where(&models.Paper{Title: input.Title}).First(&input).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("paper with the title %s already exist", input.Title)})
	}

	inputErrors := HandlingInput(input)
	if inputErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "err": inputErrors})
	}

	fileType := http.DetectContentType(input.Content)
	if fileType != "application/pdf" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "filetype must be application/pdf"})
	}

	f, err := os.Create(filepath.Join(uploadPath, input.Title+".pdf"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error in proccessing file upload", "err": err})
	}
	defer f.Close()

	if _, err := f.Write(input.Content); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	db.Create(&input)
	db.Save(&input)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "message": "paper created"})
}

func DownloadPaper(c *fiber.Ctx) error {
	paperId := c.Locals("id")

	db := database.DBConn
	paper := new(models.Paper)

	if err := db.Select("title").First(&paper, paperId).Error; err != nil {
		return handlingRowError(err, c)
	}

	return c.Download(fmt.Sprintf("./uploads/%d-%s.pdf", paperId, paper.Title))
}

func UpdatePaper(c *fiber.Ctx) error {
	db := database.DBConn
	input := new(models.Paper)
	paperId := c.Locals("id")
	paper := new(models.Paper)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "review your input"})
	}

	inputErrors := HandlingInput(input)
	if inputErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "err": inputErrors})
	}

	if err := db.First(&paper, paperId).Error; err != nil {
		return handlingRowError(err, c)
	}

	if input.Title != paper.Title {
		os.Rename(fmt.Sprintf("./uploads/%d-%s.pdf", paper.ID, paper.Title), fmt.Sprintf("./uploads/%d-%s.pdf", paper.ID, input.Title))
		paper.Title = input.Title
	}

	if bytes.Compare(input.Content, paper.Content) != 0 {
		err := ioutil.WriteFile(fmt.Sprintf("./uploads/%d-%s.pdf", paper.ID, paper.Title), input.Content, 0666)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "error reading the file content", "error": err})
		}
		paper.Content = input.Content
	}
	paper.Abstract = input.Abstract
	paper.FacultyID = input.FacultyID
	//paper.Authors = input.Authors

	db.Model(&paper).Association("Authors").Replace(input.Authors)

	db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&paper)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "message": "paper updated"})
}

func DeletePaper(c *fiber.Ctx) error {
	paperId := c.Locals("id")
	db := database.DBConn
	paper := new(models.Paper)

	if err := db.Delete(&paper, paperId).Error; err != nil {
		handlingRowError(err, c)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Paper deleted"})
}
