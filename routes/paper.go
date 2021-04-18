package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/controllers"
)

func PaperRoutes(app *fiber.App) {
	app.Get("/paper", controllers.ReadAllPaper)
	app.Post("/paper", controllers.CreatePaper)
}
