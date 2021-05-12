package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/database"
	"github.com/jun2900/online-library/models"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

const uploadPath = "./uploads"

type inputPaper struct {
	Title     string   `json:"title" validate:"required"`
	Abstract  string   `json:"abstract" validate:"required"`
	Content   []byte   `json:"content" validate:"required"`
	FacultyID uint     `json:"facultyId" validate:"required"`
	Authors   []string `json:"authors" validate:"required"`
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

func ReadAllAcceptedPaper(c *fiber.Ctx) error {
	db := database.DBConn
	var papers []models.Paper
	db.Where("status = ?", "accepted").Preload("Authors").Find(&papers)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "Papers": papers})
}

func ReadSpecificPaper(c *fiber.Ctx) error {
	db := database.DBConn

	paperId := c.Locals("paperId")

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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		db.Find(&authors)
		wg.Done()
	}()
	go func() {
		db.Find(&faculties)
		wg.Done()
	}()

	wg.Wait()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"authors": authors, "faculty": faculties})
}

//Handle create paper on POST
func CreatePaperPost(c *fiber.Ctx) error {
	db := database.DBConn
	input := new(inputPaper)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "review your input"})
	}

	inputErrors := HandlingInput(input)
	if inputErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "err": inputErrors})
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

	//RabbitMq connection
	conn, err := amqp.Dial(database.RabbitMqUrl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to connect to RabbitMQ"})
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to open a channel"})
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"insert_paper_content", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to open a channel"})
	}

	body, _ := json.Marshal(paper)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "success", "message": "paper created"})
}

func DownloadPaper(c *fiber.Ctx) error {
	paperId := c.Locals("paperId")

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
	paperId := c.Locals("paperId")
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

func UpdatePaperStatus(c *fiber.Ctx) error {
	paperId := c.Locals("paperId")
	db := database.DBConn

	if err := db.First(&models.Paper{}, paperId).Update("status", c.Query("status")).Error; err != nil {
		return handlingRowError(err, c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "paper updated"})
}

func DeletePaper(c *fiber.Ctx) error {
	paperId := c.Locals("paperId")
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
