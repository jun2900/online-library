package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/controllers"
	"github.com/jun2900/online-library/middleware"
)

func PaperRoutes(app *fiber.App) {
	app.Get("/papers", controllers.ReadAllPaper)
	app.Post("/paper", middleware.VerifyToken, controllers.CreatePaper)
}
