package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/jun2900/online-library/database"
	"github.com/jun2900/online-library/routes"
)

func initDatabase() {
	var err error
	dsn := os.Getenv("DSN")
	database.DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("connection open")
	//database.DBConn.AutoMigrate(&models.User{}, &models.Paper{}, &models.Author{}, &models.Faculty{})
}

func main() {
	app := fiber.New()

	//Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initDatabase()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	//Main routes
	routes.AuthRoutes(app)
	routes.PaperRoutes(app)

	port := os.Getenv("PORT")
	app.Listen(port)
}
