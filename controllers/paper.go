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

type inputPaper struct {
	Title     string   `json:"title"`
	Abstract  string   `json:"abstract"`
	Content   []byte   `json:"content"`
	FacultyID uint     `json:"facultyId"`
	Authors   []string `json:"authors"`
}

func handlingRowError(err error, c *fiber.Ctx) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "messsage": "Paper not found"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error handling query", "err": err})
}

func ReadAllPaper(c *fiber.Ctx) error {
	db := database.DBConn
	var papers []models.Paper
	db.Find(&papers)
	db.Preload("Authors").Find(&papers)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "Papers": papers})
}

func ReadSpecificPaper(c *fiber.Ctx) error {
	db := database.DBConn

	paperId := c.Locals("id")

	var paper models.Paper
	if err := db.First(&paper, paperId).Error; err != nil {
		return handlingRowError(err, c)
	}
	db.Preload("Authors").Find(&paper)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "paper": paper})
}

//Display create paper form
func CreatePaperGet(c *fiber.Ctx) error {
	db := database.DBConn
	var authors []models.Author
	var faculties []models.Faculty

	db.Find(&authors)
	db.Find(&faculties)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"authors": authors, "faculty": faculties})
}

//Handle create paper on POST
func CreatePaperPost(c *fiber.Ctx) error {
	db := database.DBConn
	input := new(inputPaper)

	inputErrors := HandlingInput(input)
	if inputErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "err": inputErrors})
	}

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "review your input"})
	}

	if err := db.Where(&models.Paper{Title: input.Title}).First(&models.Paper{}).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("paper with the title %s already exist", input.Title)})
	}

	fileType := http.DetectContentType(input.Content)
	if fileType != "application/pdf" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "filetype must be application/pdf"})
	}

	var authors []models.Author
	db.Where("name IN (?)", input.Authors).Find(&authors)

	if len(authors) != len(input.Authors) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "author/s does not exist"})
	}

	paper := &models.Paper{Title: input.Title, Abstract: input.Abstract, Content: input.Content, FacultyID: input.FacultyID, Authors: authors}
	db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&paper)

	f, err := os.Create(filepath.Join(uploadPath, fmt.Sprintf("%d-%s.pdf", paper.ID, paper.Title)))
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

//Handle updating paper on PUT
func UpdatePaperPut(c *fiber.Ctx) error {
	db := database.DBConn
	input := new(inputPaper)
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

	paper.Title = input.Title
	paper.Abstract = input.Abstract
	paper.Content = input.Content
	paper.FacultyID = input.FacultyID

	var authors []models.Author
	db.Where("name IN (?)", input.Authors).Find(&authors)

	db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&paper)
	db.Model(&paper).Association("Authors").Replace(authors)

	if input.Title != paper.Title {
		os.Rename(fmt.Sprintf("./uploads/%d-%s.pdf", paper.ID, paper.Title), fmt.Sprintf("./uploads/%d-%s.pdf", paper.ID, input.Title))
	}

	if bytes.Compare(input.Content, paper.Content) != 0 {
		err := ioutil.WriteFile(fmt.Sprintf("./uploads/%d-%s.pdf", paper.ID, paper.Title), input.Content, 0666)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "error reading the file content", "error": err})
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "message": "paper updated"})
}

func DeletePaper(c *fiber.Ctx) error {
	paperId := c.Locals("id")
	db := database.DBConn
	paper := new(models.Paper)

	if err := db.First(&paper, paperId).Error; err != nil {
		return handlingRowError(err, c)
	}

	if err := os.Remove(fmt.Sprintf("./uploads/%d-%s.pdf", paperId, paper.Title)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "error on deleting the paper content"})
	}

	if err := db.Delete(&paper, paperId).Error; err != nil {
		handlingRowError(err, c)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Paper deleted"})
}
